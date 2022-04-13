package service

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mills.io/saltyim"
)

// Service Always call NewService function to create Service
type Service struct {
	loaded       bool
	currentID    *saltyim.Identity
	isRegistered bool
	addresses    []*saltyim.Addr
	saltyClient  *saltyim.Client
}

// NewService Always call this function to create Service
func NewService() *Service {
	s := Service{addresses: make([]*saltyim.Addr, 0)}
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
	s.saltyClient = nil
	s.addresses = nil
}

func (s Service) IsCurrentIDAddr(addressStr string) bool {
	return s.CurrentIdentity() != nil &&
		s.CurrentIdentity().Addr().String() == addressStr
}

func (s *Service) CreateIdentity(addressStr string) (err error) {
	// clear Credentials if only address changes
	isCurrent := s.IsCurrentIDAddr(addressStr)
	if !isCurrent {
		s.clearCredentials()
	}
	addr, err := saltyim.ParseAddr(addressStr)
	if err != nil {
		s.clearCredentials()
		return err
	} else {
		s.currentID, err = saltyim.CreateIdentity(saltyim.WithIdentityAddr(addr))
		if err != nil {
			s.clearCredentials()
		} else {
			err = s.saveCurrentIdentity()
			if err != nil {
				log.Error(err)
				err = nil // this is intentional
			}
			err = s.createSaltyClient()
			if err != nil {
				log.Error(err)
				err = nil
			}
		}
	}
	return err // err is expected to be nil here
}

func (s *Service) CurrentIdentity() *saltyim.Identity {
	return s.currentID
}

func (s *Service) Register(addrStr string) (err error) {
	if s.currentID == nil {
		err = s.CreateIdentity(addrStr)
		if err != nil {
			s.isRegistered = false
			return err
		}
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

func (s *Service) GetAddr(addr string) *saltyim.Addr {
	if len(s.addresses) == 0 {
		return nil
	}
	for _, address := range s.addresses {
		if address != nil && address.String() == addr {
			return address
		}
	}
	return nil
}

func (s *Service) NewChat(addrStr string) (err error) {
	var addr *saltyim.Addr
	if s.addresses == nil {
		s.addresses = make([]*saltyim.Addr, 0)
	} else {
		if cl := s.GetAddr(addrStr); cl != nil {
			return errors.New("client already exist")
		}
	}
	addr, err = saltyim.LookupAddr(addrStr)
	if err != nil {
		return err
	}
	s.addresses = append(s.addresses, addr)
	return err
}
func (s *Service) Addresses() []*saltyim.Addr {
	if s.addresses == nil {
		s.addresses = make([]*saltyim.Addr, 0)
	}
	return s.addresses
}

func (s *Service) SendMessage(address string, msg string) (err error) {
	addr := s.GetAddr(address)
	if addr == nil {
		return errors.New("address not found")
	}
	if s.saltyClient == nil {
		err = s.createSaltyClient()
		if err != nil {
			return err
		}
	}
	err = s.saltyClient.SendToAddr(addr, msg)
	return err
}

func (s *Service) createSaltyClient() (err error) {
	currentID := s.CurrentIdentity()
	if currentID == nil {
		return errors.New("current id is nil")
	}
	contents := currentID.Contents()
	idOption := saltyim.WithIdentityBytes(contents)
	clientOptions := saltyim.WithClientIdentity(idOption)
	s.saltyClient, err = saltyim.NewClient(s.CurrentIdentity().Addr(), clientOptions)
	if err != nil {
		return err
	}
	s.runClientService()
	return err
}

func (s *Service) runClientService() {
	if s.currentID == nil || s.saltyClient == nil {
		return
	}
	s.saltyClient.SetSend(&saltyim.ProxySend{SendEndpoint: AppSendEndPoint})
	s.saltyClient.SetLookup(&saltyim.ProxyLookup{LookupEndpoint: AppLookupEndPoint})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("recovered from error ", r)
			}
			log.Debugln("returning from runClientService's goroutine")
		}()
		var ctx context.Context
		ctx = context.Background()
		inboxCh := s.saltyClient.Subscribe(ctx, "", "", "")
		outboxCh := s.saltyClient.OutboxClient(nil).Subscribe(ctx, "", "", "")
		for {
			select {
			case <-ctx.Done():
				close(inboxCh)
				close(outboxCh)
				return
			case msg := <-inboxCh:
				fmt.Println(msg)
			case msg := <-outboxCh:
				fmt.Println(msg)
			}
		}
	}()
}
