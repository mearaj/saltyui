package service

import (
	"go.mills.io/saltyim"
	"time"
)

// Service Always call NewService function to create Service
type Service struct {
	loaded    bool
	currentID *saltyim.Identity
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
		return err
	} else {
		s.currentID, err = saltyim.CreateIdentity(saltyim.WithIdentityAddr(addr))
	}
	return err // err is expected to be nil here
}

func (s *Service) CurrentIdentity() *saltyim.Identity {
	return s.currentID
}
