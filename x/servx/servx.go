package servx

import (
	"errors"
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
}

func ServiceFetch(name string) (Service, bool) {
	sv, ok := services[name]
	return sv, ok
}

func ServiceRegister(name string, sv Service) error {

	if sv.Error() != nil {
		return sv.Error()
	}

	lock.Lock()
	defer lock.Unlock()

	services[name] = sv

	return nil
}

func ServiceStart(name string) error {

	lock.Lock()
	defer lock.Unlock()

	sv, ok := services[name]
	if !ok {
		return errServNotFound
	}

	return sv.Start()
}

func ServiceStop(name string) error {

	lock.Lock()
	defer lock.Unlock()

	if sv, ok := services[name]; ok {
		return sv.Stop()
	}

	return nil
}

func ServiceRestart(name string) error {

	lock.Lock()
	defer lock.Unlock()

	sv, ok := services[name]
	if !ok {
		return errServNotFound
	}

	if err := sv.Stop(); err != nil {
		return err
	}

	return sv.Start()
}
