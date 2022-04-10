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
	clients      []*saltyim.Addr
}

// NewService Always call this function to create Service
func NewService() *Service {
	s := Service{clients: make([]*saltyim.Addr, 0)}
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
				alog.Logger().Println(err)
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
func (s *Service) GetClient(addr string) *saltyim.Addr {
	if len(s.clients) == 0 {
		return nil
	}
	for _, cl := range s.clients {
		if cl != nil && cl.String() == addr {
			return cl
		}
	}
	return nil
}

func (s *Service) NewChat(addrStr string) (err error) {
	var addr *saltyim.Addr
	if s.clients == nil {
		s.clients = make([]*saltyim.Addr, 0)
	} else {
		if cl := s.GetClient(addrStr); cl != nil {
			return errors.New("client already exist")
		}
	}
	addr, err = saltyim.LookupAddr(addrStr)
	if err != nil {
		return err
	}
	s.clients = append(s.clients, addr)
	return err
}
func (s *Service) Clients() []*saltyim.Addr {
	if s.clients == nil {
		s.clients = make([]*saltyim.Addr, 0)
	}
	return s.clients
}
