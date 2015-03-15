// Copyright 2015 lessOS.com, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package srvmgr

import (
	"errors"
	"fmt"
	"sync"
)

var (
	lock            sync.Mutex
	services        = map[string]Service{}
	errServNotFound = errors.New("No Service Found")
)

type Service interface {
	Error() error
	Start() error
	Stop() error
	Status() (Status, error)
}

func ServiceFetch(name string) (Service, bool) {
	srv, ok := services[name]
	return srv, ok
}

func ServiceRegister(name string, srv Service) error {

	if srv.Error() != nil {
		return srv.Error()
	}

	lock.Lock()
	defer lock.Unlock()

	services[name] = srv

	return nil
}

func ServiceStart(name string) error {

	lock.Lock()
	defer lock.Unlock()

	srv, ok := services[name]
	if !ok {
		return errServNotFound
	}

	return srv.Start()
}

func ServiceStop(name string) error {

	lock.Lock()
	defer lock.Unlock()

	if srv, ok := services[name]; ok {
		return srv.Stop()
	}

	return nil
}

func ServiceRestart(name string) error {

	lock.Lock()
	defer lock.Unlock()

	srv, ok := services[name]
	if !ok {
		return errServNotFound
	}

	if err := srv.Stop(); err != nil {
		return err
	}

	return srv.Start()
}

func ServiceRemove(name string) {

	lock.Lock()
	defer lock.Unlock()

	if srv, ok := services[name]; ok {
		srv.Stop()
		delete(services, name)
	}
}

func ServiceStatus(name string) (Status, error) {

	lock.Lock()
	defer lock.Unlock()

	if srv, ok := services[name]; ok {
		return srv.Status()
	}

	return Status{}, fmt.Errorf("Service Not Found")
}
