package service

import (
	"errors"
	"github.com/mearaj/saltyui/alog"
	"go.mills.io/saltyim"
)

// Service Always call NewService function to create Service
type Service struct {
	loaded       bool
	currentID    *saltyim.Identity
	isRegistered bool
	address      string
}

// NewService Always call this function to create Service
func NewService() *Service {
	s := Service{}
	s.init()
	return &s
}

func (s *Service) init() {
	go func() {
		_ = s.loadCurrentIdentity()
		s.loaded = true
	}()
}

func (s *Service) Loaded() bool {
	return s.loaded
}
func (s *Service) clearCredentials() {
	s.currentID = nil
	s.isRegistered = false
	s.address = ""
}

func (s *Service) CreateIdentity(address string) (err error) {
	addr, err := saltyim.ParseAddr(address)
	if err != nil {
		s.clearCredentials()
		return err
	} else {
		s.currentID, err = saltyim.CreateIdentity(saltyim.WithIdentityAddr(addr))
		if err != nil {
			s.clearCredentials()
		} else {
			s.address = address
			err = s.saveCurrentIdentity()
			if err != nil {
				alog.Println(err)
				err = nil //
			}
		}
	}
	return err // err is expected to be nil here
}

func (s *Service) CurrentIdentity() *saltyim.Identity {
	return s.currentID
}

func (s *Service) Register() (err error) {
	if s.currentID == nil {
		err = errors.New("current id is nil")
		s.isRegistered = false
		return err
	}
	ops := saltyim.WithClientIdentity(saltyim.WithIdentityBytes(s.currentID.Contents()))
	cl, err := saltyim.NewClient(s.currentID.Addr(), ops)
	if err != nil {
		s.isRegistered = false
		return err
	}
	err = cl.Register("https://salty.mills.io/")
	if err != nil {
		s.isRegistered = false
		return err
	}
	s.isRegistered = true
	return err // err is expected to be nil here
}

func (s *Service) IsRegistered() bool {
	return s.currentID != nil && s.isRegistered
}
func (s *Service) Address() string {
	return s.address
}
