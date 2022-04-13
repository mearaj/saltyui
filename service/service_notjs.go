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

func (s *Service) loadCurrentIdentity() (err error) {
	configDir, err := app.DataDir()
	if err != nil {
		return err
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	userKeyFileName := fmt.Sprintf("%s.key", usr.Username)
	appDir := filepath.Join(configDir, AppCfgDirName, userKeyFileName)
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
	userFilePath := filepath.Join(configDir, AppCfgDirName, userFileName)
	contentsBytes := currentId.Contents()
	err = os.WriteFile(userFilePath, contentsBytes, 0644)
	if err != nil {
		return err
	}
	userID := currentId.Addr().String()
	userIDFileName := fmt.Sprintf("%s.key", userID)
	userIDFilePath := filepath.Join(configDir, AppCfgDirName, userIDFileName)
	err = os.WriteFile(userIDFilePath, contentsBytes, 0644)
	if err != nil {
		return err
	}
	return err // err is expected to be nil here
}
