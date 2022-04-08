package service

import (
	"errors"
	"go.mills.io/saltyim"
	"time"
)

// Service Always call NewService function to create Service
type Service struct {
	loaded       bool
	currentID    *saltyim.Identity
	isRegistered bool
}

// NewService Always call this function to create Service
func NewService() *Service {
	s := Service{}
	s.init()
	return &s
}

func (s *Service) init() {
	go func() {
		time.Sleep(time.Millisecond * 250)
		s.loaded = true
	}()
}

func (s *Service) Loaded() bool {
	return s.loaded
}

func (s *Service) CreateIdentity(address string) (err error) {
	addr, err := saltyim.ParseAddr(address)
	if err != nil {
		s.currentID = nil
		s.isRegistered = false
		return err
	} else {
		s.currentID, err = saltyim.CreateIdentity(saltyim.WithIdentityAddr(addr))
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
