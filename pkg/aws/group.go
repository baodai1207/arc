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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"

	"github.com/cisco/arc/pkg/aaa"
	"github.com/cisco/arc/pkg/config"
	"github.com/cisco/arc/pkg/log"
	"github.com/cisco/arc/pkg/msg"
	"github.com/cisco/arc/pkg/resource"
)

type group struct {
	*config.Group
	provider *identityManagementProvider
	iam      *iam.IAM

	identityManagement *identityManagement

	group    *iam.Group
	policies []string
}

func newGroup(grp resource.Group, cfg *config.Group, prov *identityManagementProvider) (resource.ProviderGroup, error) {
	log.Debug("Initializing AWS Group %q", cfg.Name())

	g := &group{
		Group:              cfg,
		provider:           prov,
		iam:                prov.iam,
		identityManagement: grp.IdentityManagement().ProviderIdentityManagement().(*identityManagement),
	}

	for _, p := range g.Policies() {
		policy := newIamPolicy(prov.number, p)
		g.policies = append(g.policies, policy.String())
	}

	return g, nil
}

func (g *group) Audit(flags ...string) error {
	if len(flags) == 0 || flags[0] == "" {
		return fmt.Errorf("No flag set to find audit object")
	}
	a := aaa.AuditBuffer[flags[0]]
	if a == nil {
		return fmt.Errorf("Audit Object does not exist")
	}
	if g.group == nil {
		a.Audit(aaa.Configured, "%s", g.Name())
		return nil
	}
	return nil
}

func (g *group) set(group *iam.Group) {
	g.group = group
}

func (g *group) clear() {
	g.group = nil
}

func (g *group) Created() bool {
	return g.group != nil
}

func (g *group) Destroyed() bool {
	return g.group == nil
}

func (g *group) Create(flags ...string) error {
	if err := g.createGroup(); err != nil {
		return err
	}
	if err := g.attachPolicies(); err != nil {
		return err
	}
	return nil
}

func (g *group) Load() error {
	if err := g.loadPolicies(); err != nil {
		return err
	}
	if group := g.identityManagement.groupCache.find(g); group != nil {
		log.Debug("Skipping group load, cached...")
		g.set(group)
		return nil
	}
	params := &iam.GetGroupInput{
		GroupName: aws.String(g.Name()),
	}
	resp, err := g.iam.GetGroup(params)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchEntity") {
			log.Debug("No Such Entity: Group %q", g.Name())
			return nil
		}
		return err
	}
	g.group = resp.Group
	return nil
}

func (g *group) Destroy(flags ...string) error {
	if err := g.detachPolicies(); err != nil {
		return err
	}
	if err := g.deleteGroup(); err != nil {
		return err
	}
	g.clear()
	return nil
}

func (g *group) Provision(flags ...string) error {
	return nil
}

func (g *group) Info() {
	if g.Destroyed() {
		return
	}
	msg.Info("Group")
	msg.Detail("%-20s\t%s", "name", aws.StringValue(g.group.GroupName))
}

func (g *group) loadPolicies() error {
	return nil
}

func (g *group) createGroup() error {
	msg.Info("Group Creation: %s", g.Name())

	params := &iam.CreateGroupInput{
		GroupName: aws.String(g.Name()),
	}

	resp, err := g.iam.CreateGroup(params)
	if err != nil {
		return err
	}
	g.group = resp.Group
	if err := g.Load(); err != nil {
		return err
	}
	msg.Detail("Group created: %s", g.Name())
	return nil
}

func (g *group) attachPolicies() error {
	msg.Info("Attach Policies")
	for _, policy := range g.policies {
		params := &iam.AttachGroupPolicyInput{
			GroupName: aws.String(g.Name()),
			PolicyArn: aws.String(policy),
		}
		_, err := g.iam.AttachGroupPolicy(params)
		if err != nil {
			return err
		}
	}
	msg.Detail("Attached Policies")
	return nil
}

func (g *group) deleteGroup() error {
	msg.Info("Group Deletion: %s", g.Name())
	params := &iam.DeleteGroupInput{
		GroupName: aws.String(g.Name()),
	}

	_, err := g.iam.DeleteGroup(params)
	if err != nil {
		return err
	}
	msg.Detail("Group deleted: %s", g.Name())
	return nil
}

func (g *group) detachPolicies() error {
	return nil
}
