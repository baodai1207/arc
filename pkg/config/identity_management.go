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

package config

import "github.com/cisco/arc/pkg/msg"

type IdentityManagement struct {
	Region_  string    `json:"region"`
	Policies []*Policy `json:"policies"`
	Roles    []*Role   `json:"roles"`
	Groups   []*Group  `json:"groups"`
	Users    []*User   `json:"users"`
}

func (i *IdentityManagement) Region() string {
	return i.Region_
}

// Print provides a user friendly way to view the entire storage configuration.
func (i *IdentityManagement) Print() {
	msg.Info("Identity Management Config")
	msg.IndentInc()
	for _, p := range i.Policies {
		p.Print()
	}
	for _, r := range i.Roles {
		r.Print()
	}
	for _, g := range i.Groups {
		g.Print()
	}
	for _, u := range i.Users {
		u.Print()
	}
	msg.IndentDec()
}