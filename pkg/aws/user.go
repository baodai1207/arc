//
// Copyright (c) 2018, Cisco Systems
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice, this
//   list of conditions and the following disclaimer in the documentation and/or
//   other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
// ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package aws

import (
	"fmt"
	// "strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/cisco/arc/pkg/aaa"
	"github.com/cisco/arc/pkg/config"
	"github.com/cisco/arc/pkg/log"
	"github.com/cisco/arc/pkg/msg"
	"github.com/cisco/arc/pkg/resource"
)

type user struct {
	*config.User
	provider *identityManagementProvider
	iam      *iam.IAM

	identityManagement *identityManagement

	user     *iam.User
	policies []string
}

func newUser(usr resource.User, cfg *config.User, prov *identityManagementProvider) (resource.ProviderUser, error) {
	log.Debug("Initializing AWS User %q", cfg.Name())

	u := &user{
		User:               cfg,
		provider:           prov,
		iam:                prov.iam,
		identityManagement: usr.IdentityManagement().ProviderIdentityManagement().(*identityManagement),
	}

	for _, p := range u.Policies() {
		policy := newIamPolicy(prov.number, p)
		u.policies = append(u.policies, policy.String())
	}

	return u, nil
}

func (u *user) Audit(flags ...string) error {
	if len(flags) == 0 || flags[0] == "" {
		return fmt.Errorf("No flag set to find audit object")
	}
	a := aaa.AuditBuffer[flags[0]]
	if a == nil {
		return fmt.Errorf("Audit Object does not exist")
	}
	if u.user == nil {
		a.Audit(aaa.Configured, "%s", u.Name())
		return nil
	}
	return nil
}

func (u *user) set(user *iam.User) {
	u.user = user
}

func (u *user) clear() {
	u.user = nil
}

func (u *user) Created() bool {
	return u.user != nil
}

func (u *user) Destroyed() bool {
	return u.user == nil
}

func (u *user) Create(flags ...string) error {
	if err := u.createUser(); err != nil {
		return err
	}
	if err := u.attachPolicies(); err != nil {
		return err
	}
	if err := u.attachToGroups(); err != nil {
		return err
	}
	return nil
}

func (u *user) Load() error {
	if err := u.loadPolicies(); err != nil {
		return err
	}
	return nil
}

func (u *user) Destroy(flags ...string) error {
	if err := u.detachPolicies(); err != nil {
		return err
	}
	if err := u.deleteUser(); err != nil {
		return err
	}
	u.clear()
	return nil
}

func (u *user) Provision(flags ...string) error {
	return nil
}

func (u *user) Info() {
	if u.Destroyed() {
		return
	}
	msg.Info("User")
	msg.Detail("%-20s\t%s", "name", aws.StringValue(u.user.UserName))
}

func (u *user) loadPolicies() error {
	return nil
}

func (u *user) createUser() error {
	msg.Info("User Creation: %s", u.Name())

	params := &iam.CreateUserInput{
		UserName: aws.String(u.Name()),
	}

	resp, err := u.iam.CreateUser(params)
	if err != nil {
		return err
	}
	u.user = resp.User
	if err := u.Load(); err != nil {
		return err
	}
	msg.Detail("User created: %s", u.Name())
	return nil
}

func (u *user) attachPolicies() error {
	msg.Info("Attach Policies")
	for _, policy := range u.policies {
		params := &iam.AttachUserPolicyInput{
			PolicyArn: aws.String(policy),
			UserName:  aws.String(u.Name()),
		}
		_, err := u.iam.AttachUserPolicy(params)
		if err != nil {
			return err
		}
	}
	msg.Detail("Attached Policies")
	return nil
}

func (u *user) attachToGroups() error {
	msg.Info("Attach to Groups")
	for _, g := range u.Groups() {
		params := &iam.AddUserToGroupInput{
			GroupName: aws.String(g),
			UserName:  aws.String(u.Name()),
		}

		_, err := u.iam.AddUserToGroup(params)
		if err != nil {
			return err
		}
		msg.Detail("Attached to Group %s", g)
	}
	return nil
}

func (u *user) deleteUser() error {
	msg.Info("User Deletion: %s", u.Name())
	params := &iam.DeleteUserInput{
		UserName: aws.String(u.Name()),
	}

	_, err := u.iam.DeleteUser(params)
	if err != nil {
		return err
	}
	msg.Detail("User deleted: %s", u.Name())
	return nil
}

func (u *user) detachPolicies() error {
	msg.Info("Detach Policies")
	for _, policy := range u.policies {
		params := &iam.DetachUserPolicyInput{
			PolicyArn: aws.String(policy),
			UserName:  aws.String(u.Name()),
		}
		_, err := u.iam.DetachUserPolicy(params)
		if err != nil {
			return err
		}
	}
	msg.Detail("Detached Policies")
	return nil
}
