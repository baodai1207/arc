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
	"github.com/cisco/arc/pkg/log"
)

type userCacheEntry struct {
	deployed   *iam.User
	configured *user
}

type userCache struct {
	cache   map[string]*userCacheEntry
	unnamed []*iam.User
}

func newUserCache(i *identityManagement) (*userCache, error) {
	log.Debug("Initializing AWS User Cache")

	c := &userCache{
		cache: map[string]*userCacheEntry{},
	}

	next := ""
	for {
		params := &iam.ListUsersInput{}

		if next != "" {
			params.Marker = aws.String(next)
		}

		resp, err := i.iam.ListUsers(params)
		if err != nil {
			return nil, err
		}
		truncated := false
		if resp.IsTruncated != nil {
			truncated = *resp.IsTruncated
		}
		next = ""
		if resp.Marker != nil {
			next = *resp.Marker
		}

		for _, user := range resp.Users {
			if user.UserName == nil {
				log.Verbose("Unnamed user")
				c.unnamed = append(c.unnamed, user)
				continue
			}
			log.Debug("Caching %s", aws.StringValue(user.UserName))
			c.cache[aws.StringValue(user.UserName)] = &userCacheEntry{deployed: user}
		}
		if truncated == false {
			break
		}
	}

	return c, nil
}

func (c *userCache) find(u *user) *iam.User {
	e := c.cache[u.Name()]
	if e == nil {
		return nil
	}
	e.configured = u
	return e.deployed
}

func (c *userCache) remove(u *user) {
	log.Debug("Deleting %s from userCache", u.Name())
	delete(c.cache, u.Name())
}

func (c *userCache) audit(flags ...string) error {
	if len(flags) == 0 || flags[0] == "" {
		return fmt.Errorf("No flag set to find audit object")
	}
	a := aaa.AuditBuffer[flags[0]]
	if a == nil {
		return fmt.Errorf("Audit Object does not exist")
	}
	for k, v := range c.cache {
		if v.configured == nil {
			a.Audit(aaa.Deployed, "%s", k)
		}
	}
	if c.unnamed != nil {
		a.Audit(aaa.Deployed, "\r")
		for i, v := range c.unnamed {
			u := "\t" + strings.Replace(fmt.Sprintf("%+v", v), "\n", "\n\t", -1)
			m := fmt.Sprintf("Unnamed User %d - User ID: %q %s", i+1, *v.UserId, u)
			a.Audit(aaa.Deployed, m)
		}
	}
	return nil
}
