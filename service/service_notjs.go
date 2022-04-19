//go:build !js

package service

import (
	"errors"
	"fmt"
	"gioui.org/app"
	log "github.com/sirupsen/logrus"
	"go.mills.io/saltyim"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Service Always call NewService function to create Service
type service struct {
	loaded            bool
	currentIDState    IdentityState
	isRegistered      bool
	userIdentities    []*saltyim.Identity
	contactsAddresses []*saltyim.Addr
	saltyClient       *saltyim.Client
	initialized       bool
}

func (s *service) IsRegistered() bool {
	return s.CurrentIdentityState() != nil && s.isRegistered
}

func (s *service) GetAddr(addr string) *saltyim.Addr {
	if len(s.ContactsAddresses()) == 0 {
		return nil
	}
	for _, address := range s.ContactsAddresses() {
		if address != nil && address.String() == addr {
			return address
		}
	}
	return nil
}

func (s *service) Addresses() []*saltyim.Addr {
	if s.contactsAddresses == nil {
		s.contactsAddresses = make([]*saltyim.Addr, 0)
	}
	return s.contactsAddresses
}

func (s *service) Loaded() bool {
	return s.loaded
}

func (s *service) init() {
	_ = s.loadCurrentIdentity()
	s.loaded = true
}

func (s *service) loadCurrentIdentity() (err error) {
	configDir, err := app.DataDir()
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	userKeyFileName := fmt.Sprintf("%s.key", usr.Username)
	appDir := filepath.Join(configDir, DBPathCfgName, userKeyFileName)
	contents, err := ioutil.ReadFile(appDir)
	if err != nil {
		return err
	}
	ops := saltyim.WithIdentityBytes(contents)
	s.currentID, err = saltyim.GetOrCreateIdentity(ops)
	if err != nil {
		return err
	}
	err = s.createSaltyClient()
	if err != nil {
		log.Error(err)
		err = nil
	}

	return err // err is expected to be nil here
}

func (s *service) saveCurrentIdentity() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				log.Error("recovered from panic", r)
			}
			errCh <- err
			close(errCh)
		}()
		currentId := s.currentID
		if currentId == nil {
			err = errors.New("current identity is nil")
		}
		configDir, err := app.DataDir()
		if err != nil {
			return
		}
		usr, err := user.Current()
		if err != nil {
			return
		}
		userFileName := fmt.Sprintf("%s.key", usr.Username)
		userFilePath := filepath.Join(configDir, DBPathCfgName, userFileName)
		contentsBytes := currentId.Contents()
		err = os.WriteFile(userFilePath, contentsBytes, 0644)
		if err != nil {
			return
		}
		userID := currentId.Addr().String()
		userIDFileName := fmt.Sprintf("%s.key", userID)
		userIDFilePath := filepath.Join(configDir, DBPathCfgName, userIDFileName)
		err = os.WriteFile(userIDFilePath, contentsBytes, 0644)
		return
	}()
	return errCh
}
