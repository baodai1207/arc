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

type storage struct {
	*config.Storage
	opt options
}

func newStorage(cfg *config.Storage, p *storageProvider) (resource.ProviderStorage, error) {
	log.Info("Initializing Mock Storage")
	s := &storage{
		Storage: cfg,
		opt:     options{p.Provider.Data},
	}
	if s.opt.err("stor.New") {
		return nil, err{"stor.New"}
	}
	return s, nil
}

func (s *storage) Load() error {
	log.Info("Loading Mock Storage")
	if s.opt.err("stor.Load") {
		return err{"stor.Load"}
	}
	return nil
}

func (s *storage) Provision(flags ...string) error {
	msg.Info("Provisioning Mock Storage")
	if s.opt.err("stor.Provision") {
		return err{"stor.Provision"}
	}
	return nil
}

func (s *storage) Audit(flags ...string) error {
	msg.Info("Auditing Mock Storage")
	if s.opt.err("stor.Audit") {
		return err{"stor.Audit"}
	}
	return nil
}

func (s *storage) Info() {
	msg.Info("Mock Storage")
	msg.Detail("...")
}
