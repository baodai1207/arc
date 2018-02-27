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

package mock

import (
	"github.com/cisco/arc/pkg/config"
	"github.com/cisco/arc/pkg/log"
	"github.com/cisco/arc/pkg/msg"
	"github.com/cisco/arc/pkg/resource"
)

type bucket struct {
	*config.Bucket
	opt options
}

func newBucket(cfg *config.Bucket, p *storageProvider) (resource.ProviderBucket, error) {
	log.Info("Initializing Mock Bucket %q", cfg.Name())
	b := &bucket{
		Bucket: cfg,
		opt:    options{p.Provider.Data},
	}
	if b.opt.err("bkt.New") {
		return nil, err{"bkt.New"}
	}
	return b, nil
}

func (b *bucket) Load() error {
	log.Info("Loading Mock Bucket %q", b.Name())
	if b.opt.err("bkt.Load") {
		return err{"bkt.Load"}
	}
	return nil
}

func (b *bucket) Create(flags ...string) error {
	msg.Info("Creating Mock Bucket %q", b.Name())
	if b.opt.err("bkt.Create") {
		return err{"bkt.Create"}
	}
	return nil
}

func (b *bucket) Created() bool {
	if b.opt.err("bkt.Created") {
		return false
	}
	return true
}

func (b *bucket) Destroy(flags ...string) error {
	msg.Info("Destroying Mock Bucket %q", b.Name())
	if b.opt.err("bkt.Destroy") {
		return err{"bkt.Destroy"}
	}
	return nil
}

func (b *bucket) Destroyed() bool {
	if b.opt.err("bkt.Destroyed") {
		return false
	}
	return true
}

func (b *bucket) Provision(flags ...string) error {
	msg.Info("Provisioning Mock Bucket %q", b.Name())
	if b.opt.err("bkt.Provision") {
		return err{"bkt.Provision"}
	}
	return nil
}

func (b *bucket) Audit(flags ...string) error {
	msg.Info("Auditing Mock Bucket %q", b.Name())
	if b.opt.err("bkt.Audit") {
		return err{"bkt.Audit"}
	}
	return nil
}

func (b *bucket) Info() {
	msg.Info("Mock Bucket")
	msg.Detail("%-20s\t%s", "Name", b.Name())
	msg.Detail("%-20s\t%s", "Region", b.Region())
}

func (b *bucket) SetTags(map[string]string) error {
	msg.Info("Set Tags for Mock Bucket %q", b.Name())
	if b.opt.err("bkt.SetTags") {
		return err{"bkt.SetTags"}
	}
	return nil
}

func (b *bucket) EnableReplication() error {
	msg.Info("Enabling Replication for Mock Bucket %q", b.Name())
	if b.opt.err("bkt.EnableReplcation") {
		return err{"bkt.EnableReplication"}
	}
	return nil
}
