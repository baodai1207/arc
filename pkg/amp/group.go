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

type group struct {
	*config.Group
	identityManagement *identityManagement
	providerGroup      resource.ProviderGroup
}

func newGroup(cfg *config.Group, identityManagement *identityManagement, prov provider.IdentityManagement) (*group, error) {
	log.Debug("Initializing Group, %q", cfg.Name())
	g := &group{
		Group:              cfg,
		identityManagement: identityManagement,
	}

	var err error
	g.providerGroup, err = prov.NewGroup(g, cfg)
	if err != nil {
		return nil, err
	}

	return g, nil
}

// Route satisfies the embedded resource.Resource interface in resource.Group.
// Group handles load, create, destroy, provision, audit, config and info requests by delegating them
// to the providerGroup.
func (g *group) Route(req *route.Request) route.Response {
	log.Route(req, "Group %q", g.Name())
	switch req.Command() {
	case route.Load:
		if err := g.Load(); err != nil {
			return route.FAIL
		}
		return route.OK
	case route.Create:
		if err := g.Create(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Destroy:
		if err := g.Destroy(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Provision:
		if err := g.Provision(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Audit:
		if err := g.Audit(req.Flags().Get()...); err != nil {
			msg.Error(err.Error())
			return route.FAIL
		}
		return route.OK
	case route.Info:
		g.Info()
		return route.OK
	case route.Config:
		g.Print()
		return route.OK
	case route.Help:
		g.Help()
		return route.OK
	default:
		msg.Error("Internal Error: Unknown command " + req.Command().String())
		g.Help()
		return route.FAIL
	}
}

// Created satisfies the embedded resource.Creator interface in resource.Group.
// It delegates the call to the provider's group.
func (g *group) Created() bool {
	return g.providerGroup.Created()
}

// Destroyed satisfies the embedded resource.Destroyer interaface in resource.Group.
// It delegates the call to the provider's group.
func (g *group) Destroyed() bool {
	return g.providerGroup.Destroyed()
}

func (g *group) IdentityManagement() resource.IdentityManagement {
	return g.identityManagement
}

func (g *group) ProviderGroup() resource.ProviderGroup {
	return g.providerGroup
}

func (g *group) Provision(flags ...string) error {
	return g.providerGroup.Provision(flags...)
}

func (g *group) Load() error {
	return g.providerGroup.Load()
}

func (g *group) Create(flags ...string) error {
	if g.Created() {
		msg.Detail("Group exists, skipping...")
		return nil
	}
	return g.ProviderGroup().Create(flags...)
}

func (g *group) Destroy(flags ...string) error {
	if g.Destroyed() {
		msg.Detail("Group does not exist, skipping...")
		return nil
	}
	return g.ProviderGroup().Destroy(flags...)
}

func (g *group) Audit(flags ...string) error {
	if len(flags) == 0 || flags[0] == "" {
		return fmt.Errorf("No flag set to find the audit object")
	}
	return g.ProviderGroup().Audit(flags...)
}

func (g *group) Info() {
	if g.Destroyed() {
		return
	}
	g.ProviderGroup().Info()
}

func (g *group) Help() {
	commands := []help.Command{
		{Name: route.Create.String(), Desc: fmt.Sprintf("create group %s", g.Name())},
		{Name: route.Destroy.String(), Desc: fmt.Sprintf("destroy group %s", g.Name())},
		{Name: route.Provision.String(), Desc: fmt.Sprintf("update group %s", g.Name())},
		{Name: route.Audit.String(), Desc: fmt.Sprintf("audit group %s", g.Name())},
		{Name: route.Info.String(), Desc: "show information about allocated group"},
		{Name: route.Config.String(), Desc: "show the configuration for the given group"},
		{Name: route.Help.String(), Desc: "show this help"},
	}
	help.Print("group "+g.Name(), commands)
}
