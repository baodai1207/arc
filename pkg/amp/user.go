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

package amp

import (
	"fmt"

	"github.com/cisco/arc/pkg/config"
	"github.com/cisco/arc/pkg/help"
	"github.com/cisco/arc/pkg/log"
	"github.com/cisco/arc/pkg/msg"
	"github.com/cisco/arc/pkg/provider"
	"github.com/cisco/arc/pkg/resource"
	"github.com/cisco/arc/pkg/route"
)

type user struct {
	*config.User
	identityManagement *identityManagement
	providerUser       resource.ProviderUser
}

func newUser(cfg *config.User, identityManagement *identityManagement, prov provider.IdentityManagement) (*user, error) {
	log.Debug("Initializing User, %q", cfg.Name())
	u := &user{
		User:               cfg,
		identityManagement: identityManagement,
	}

	var err error
	u.providerUser, err = prov.NewUser(u, cfg)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Route satisfies the embedded resource.Resource interface in resource.User.
// User handles load, create, destroy, provision, audit, config and info requests by delegating them
// to the providerUser.
func (u *user) Route(req *route.Request) route.Response {
	log.Route(req, "User %q", u.Name())
	switch req.Command() {
	case route.Load:
		if err := u.Load(); err != nil {
			return route.FAIL
		}
		return route.OK
	case route.Create:
		if err := u.Create(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Destroy:
		if err := u.Destroy(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Provision:
		if err := u.Provision(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Audit:
		if err := u.Audit(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Info:
		u.Info()
		return route.OK
	case route.Config:
		u.Print()
		return route.OK
	case route.Help:
		u.Help()
		return route.OK
	default:
		msg.Error("Internal Error: Unknown command " + req.Command().String())
		u.Help()
		return route.FAIL
	}
}

// Created satisfies the embedded resource.Creator interface in resource.User.
// It delegates the call to the provider's user.
func (u *user) Created() bool {
	return u.providerUser.Created()
}

// Destroyed satisfies the embedded resource.Destroyer interaface in resource.User.
// It delegates the call to the provider's user.
func (u *user) Destroyed() bool {
	return u.providerUser.Destroyed()
}

func (u *user) IdentityManagement() resource.IdentityManagement {
	return u.identityManagement
}

func (u *user) ProviderUser() resource.ProviderUser {
	return u.providerUser
}

func (u *user) Provision(flags ...string) error {
	return u.providerUser.Provision(flags...)
}

func (u *user) Load() error {
	return u.providerUser.Load()
}

func (u *user) Create(flags ...string) error {
	if u.Created() {
		msg.Detail("User exists, skipping...")
		return nil
	}
	return u.ProviderUser().Create(flags...)
}

func (u *user) Destroy(flags ...string) error {
	if u.Destroyed() {
		msg.Detail("User does not exist, skipping...")
		return nil
	}
	return u.ProviderUser().Destroy(flags...)
}

func (u *user) Audit(flags ...string) error {
	if len(flags) == 0 || flags[0] == "" {
		return fmt.Errorf("No flag set to find the audit object")
	}
	return u.ProviderUser().Audit(flags...)
}

func (u *user) Info() {
	if u.Destroyed() {
		return
	}
	u.ProviderUser().Info()
}

func (u *user) Help() {
	commands := []help.Command{
		{Name: route.Create.String(), Desc: fmt.Sprintf("create user %s", u.Name())},
		{Name: route.Destroy.String(), Desc: fmt.Sprintf("destroy user %s", u.Name())},
		{Name: route.Provision.String(), Desc: fmt.Sprintf("update user %s", u.Name())},
		{Name: route.Audit.String(), Desc: fmt.Sprintf("audit user %s", u.Name())},
		{Name: route.Info.String(), Desc: "show information about allocated user"},
		{Name: route.Config.String(), Desc: "show the configuration for the given user"},
		{Name: route.Help.String(), Desc: "show this help"},
	}
	help.Print("user "+u.Name(), commands)
}
