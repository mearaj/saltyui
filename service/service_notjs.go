//go:build !js

package service

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"github.com/mearaj/saltyui/alog"
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
	identity          *saltyim.Identity
	identities        []*saltyim.Identity
	isRegistered      bool
	userIdentities    []*saltyim.Identity
	contactsAddresses []*saltyim.Addr
	saltyClient       *saltyim.Client
	initialized       bool
	configJSON        *ConfigJSON
	messages          map[string][]Message
	isClientRunning   bool
}

func (s *service) IsRegistered() bool {
	return s.Identity() != nil && s.isRegistered
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
	s.identity, err = saltyim.GetOrCreateIdentity(ops)
	if err != nil {
		return err
	}
	err = <-s.createSaltyClient()
	if err != nil {
		log.Error(err)
		err = nil
	}

	return err // err is expected to be nil here
}

func (s *service) saveIdentity(identity *saltyim.Identity) <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil || identity == nil {
		err = errors.New("current id or provided id is nil")
		errCh <- err
		close(errCh)
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		var configDir string
		configDir, err = app.DataDir()
		if err != nil {
			return
		}
		var usr *user.User
		usr, err = user.Current()
		if err != nil {
			return
		}
		userFileName := fmt.Sprintf("%s.key", usr.Username)
		userFilePath := filepath.Join(configDir, DBPathCfgName, userFileName)
		contentsBytes := identity.Contents()
		err = os.WriteFile(userFilePath, contentsBytes, 0644)
		if err != nil {
			return
		}
		userID := identity.Addr().String()
		userIDFileName := fmt.Sprintf("%s.key", userID)
		userIDFilePath := filepath.Join(configDir, DBPathCfgName, userIDFileName)
		err = os.WriteFile(userIDFilePath, contentsBytes, 0644)
		return
	}()
	return errCh
}
func (s *service) saveIdentities() <-chan error {
	errCh := make(chan error, 1)
	errCh <- nil
	close(errCh)
	return errCh
}
func (s *service) saveMessage(message Message) <-chan error {
	errCh := make(chan error, 1)
	errCh <- nil
	close(errCh)
	return errCh
}

func (s *service) saveContacts() <-chan error {
	errCh := make(chan error, 1)
	errCh <- nil
	close(errCh)
	return errCh
}
