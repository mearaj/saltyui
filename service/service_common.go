//go:build !js

package service

import (
	"errors"
	"fmt"
	"gioui.org/app"
	"go.mills.io/saltyim"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func (s *Service) loadCurrentIdentity() (err error) {
	configDir, err := app.DataDir()
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	userFileName := fmt.Sprintf("%s.key", usr.Username)
	appDir := filepath.Join(configDir, AppName, userFileName)
	contents, err := ioutil.ReadFile(appDir)
	if err != nil {
		return err
	}
	ops := saltyim.WithIdentityBytes(contents)
	s.currentID, err = saltyim.GetOrCreateIdentity(ops)
	if err != nil {
		return err
	}
	return err // err is expected to be nil here
}

func (s *Service) saveCurrentIdentity() error {
	currentId := s.currentID
	if currentId == nil {
		return errors.New("current identity is nil")
	}
	configDir, err := app.DataDir()
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	userFileName := fmt.Sprintf("%s.key", usr.Username)
	appDir := filepath.Join(configDir, AppName, userFileName)
	contentsBytes := currentId.Contents()
	err = os.WriteFile(appDir, contentsBytes, 0644)
	if err != nil {
		return err
	}
	return err // err is expected to be nil here
}
