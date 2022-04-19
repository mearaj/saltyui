package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mearaj/saltyui/alog"
	"github.com/sirupsen/logrus"
	"go.mills.io/saltyim"
	"strings"
	"time"
)

type Service interface {
	Loaded() bool
	IsCurrIDAddr(addrStr string) bool
	Identity() *saltyim.Identity
	Register(addrStr string) <-chan error
	CreateID(addrStr string) <-chan error
	SendMessage(addrStr string, msg string) <-chan error
	NewChat(addrStr string) <-chan error
	GetContactAddr(addrStr string) *saltyim.Addr
	ContactsAddresses() []*saltyim.Addr
	Identities() []*saltyim.Identity
	ConfigJSON() (*ConfigJSON, error)
	Messages(contactAddr string) []Message
}

type ConfigJSON struct {
	Config saltyim.Config
	Hash   string
	User   string
	Domain string
}

// NewService Always call this function to create Service
func NewService() Service {
	s := service{
		identities:        make([]*saltyim.Identity, 0),
		contactsAddresses: make([]*saltyim.Addr, 0),
	}
	go s.init()
	return &s
}

func (s *service) IsCurrIDAddr(addrStr string) bool {
	return s.Identity() != nil &&
		s.Identity().Addr().String() == addrStr
}

func (s *service) Identity() *saltyim.Identity {
	return s.identity
}

func (s *service) Identities() []*saltyim.Identity {
	return s.identities
}

func (s *service) setIdentity(identity *saltyim.Identity) {
	s.identity = identity
}

func (s *service) GetContactAddr(addrStr string) *saltyim.Addr {
	if len(s.contactsAddresses) == 0 {
		return nil
	}
	for _, addr := range s.contactsAddresses {
		if addr != nil && addr.String() == addrStr {
			return addr
		}
	}
	return nil
}
func (s *service) ContactsAddresses() []*saltyim.Addr {
	if s.contactsAddresses == nil {
		s.contactsAddresses = make([]*saltyim.Addr, 0)
	}
	return s.contactsAddresses
}

func (s *service) CreateID(addressStr string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		// clear Credentials if only address changes
		isCurrent := s.IsCurrIDAddr(addressStr)
		if !isCurrent {
			s.clearCredentials()
		}
		addr, err := saltyim.ParseAddr(addressStr)
		if err != nil {
			s.clearCredentials()
			s.setIdentity(nil)
			return
		} else {
			var currID *saltyim.Identity
			currID, err = saltyim.CreateIdentity(saltyim.WithIdentityAddr(addr))
			if err != nil {
				s.clearCredentials()
				s.setIdentity(currID)
				return
			} else {
				s.setIdentity(currID)
				s.addIdentity(currID)
				err = <-s.saveIdentity(currID)
				if err != nil {
					return
				}
				err = <-s.createSaltyClient()
				if err != nil {
					return
				}
				err = <-s.saveIdentities()
				if err != nil {
					return
				}
				if !s.IsClientRunning() {
					s.runClientService()
				}
			}
		}
	}()
	return errCh
}
func (s *service) Register(addrStr string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		if s.Identity() == nil {
			select {
			case err = <-s.CreateID(addrStr):
				if err != nil {
					s.isRegistered = false
					return
				}
			}
		}
		ops := saltyim.WithClientIdentity(saltyim.WithIdentityBytes(s.Identity().Contents()))
		cl, err := saltyim.NewClient(s.Identity().Addr(), ops)
		if err != nil {
			s.isRegistered = false
			return
		}
		err = cl.Register(AppRegisterEndPoint)
		if err != nil {
			s.isRegistered = false
			return
		}
		s.isRegistered = true
		return
	}()
	return errCh
}
func (s *service) ConfigJSON() (*ConfigJSON, error) {
	if s.configJSON != nil {
		return s.configJSON, nil
	}
	var err error
	defer func() {
		recoverPanic(alog.Logger())
	}()
	if s.Identity() == nil {
		return nil, errors.New("current id is nil")
	}
	if s.saltyClient == nil {
		err = <-s.createSaltyClient()
		if err != nil {
			return nil, err
		}
	}
	// url can be nil and error is unrecoverable, esp in wasm
	url := s.Identity().Addr().Endpoint()
	if url == nil {
		return nil, errors.New("endpoint url is nil")
	}

	endPointSlice := strings.Split(url.String(), "/")
	endPoint := fmt.Sprintf("/%s", endPointSlice[len(endPointSlice)-1])
	key := s.Identity().Key().PublicKey().String()
	hashSlice := strings.Split(s.Identity().Addr().HashURI(), "/")
	hash := hashSlice[len(hashSlice)-1]
	configJson := ConfigJSON{
		Config: saltyim.Config{
			Endpoint: endPoint,
			Key:      key,
		},
		Hash:   hash,
		User:   s.Identity().Addr().User,
		Domain: s.Identity().Addr().Domain,
	}
	// Todo: Just for debug, remove it
	st, err := json.MarshalIndent(&configJson.Config, "", "  ")
	if err != nil {
		alog.Logger().Println(err)
		err = nil
	}
	fmt.Println(string(st))
	fmt.Println(hash)

	s.configJSON = &configJson
	return s.configJSON, nil
}

func (s *service) SendMessage(address string, msg string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		addr := s.GetContactAddr(address)
		if addr == nil {
			err = errors.New("address not found")
			return
		}
		id := s.Identity()
		if id == nil {
			err = errors.New("current id is nil")
			return
		}
		message := Message{
			UserAddr:    s.Identity().Addr().String(),
			ContactAddr: address,
			From:        s.Identity().Addr().String(),
			To:          address,
			Created:     time.Now().String(),
			Text:        msg,
			Key:         s.Identity().Key().String(),
		}
		s.addMessage(message)
		s.saveMessage(message)
		if s.saltyClient == nil {
			err = <-s.createSaltyClient()
			if err != nil {
				return
			}
		}
		if !s.IsClientRunning() {
			<-s.runClientService()
		}
		_, err = saltyim.LookupAddr(s.Identity().Addr().String())
		if err != nil {
			alog.Logger().Errorln(err)
			return
		}
		if s.saltyClient != nil {
			//err = s.saltyClient.SendToAddr(addr, msg) // <-Todo it panics
			err = s.saltyClient.Send(address, msg)
			if err != nil {
				return
			}
		}
	}()
	return errCh
}

func (s *service) Messages(addr string) []Message {
	if s.messages == nil {
		s.messages = map[string][]Message{}
	}
	return s.messages[addr]
}

func (s *service) addMessage(message Message) bool {
	if s.messages == nil {
		s.messages = map[string][]Message{}
	}
	if s.Identity() == nil {
		return false
	}
	contactAddr := message.ContactAddr
	if _, ok := s.messages[contactAddr]; !ok {
		s.messages[contactAddr] = make([]Message, 1)
	}
	if !s.isMessageDuplicate(message) {
		s.messages[contactAddr] = append(s.messages[contactAddr], message)
	}
	return true
}

func (s *service) isMessageDuplicate(message Message) bool {
	if s.messages == nil {
		s.messages = map[string][]Message{}
	}
	contactAddr := message.ContactAddr
	if _, ok := s.messages[contactAddr]; !ok {
		s.messages[contactAddr] = make([]Message, 1)
		return false
	}
	for _, msg := range s.messages[contactAddr] {
		isDuplicate := msg.isDuplicate(message)
		if isDuplicate {
			return isDuplicate
		}
	}
	return false
}

func (s *service) clearCredentials() {
	s.setIdentity(nil)
	s.isRegistered = false
	s.saltyClient = nil
	s.contactsAddresses = nil
	s.configJSON = nil
}

func (s *service) createSaltyClient() <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil {
		err = errors.New("current id is nil")
		errCh <- err
		close(errCh)
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		contents := s.Identity().Contents()
		idOption := saltyim.WithIdentityBytes(contents)
		clientOptions := saltyim.WithClientIdentity(idOption)
		s.saltyClient, err = saltyim.NewClient(s.Identity().Addr(), clientOptions)
		if err != nil {
			return
		}
	}()
	return errCh
}

func (s *service) runClientService() <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.IsClientRunning() {
		err = errors.New("client already running")
		errCh <- err
		close(errCh)
		return errCh
	}
	if s.Identity() == nil {
		err = errors.New("current id is nil")
		errCh <- err
		close(errCh)
		return errCh
	}
	if s.saltyClient == nil {
		err = errors.New("client is nil")
		errCh <- err
		close(errCh)
		return errCh
	}
	s.saltyClient.SetSend(&saltyim.ProxySend{SendEndpoint: AppSendEndPoint})
	s.saltyClient.SetLookup(&saltyim.ProxyLookup{LookupEndpoint: AppLookupEndPoint})

	// FixMe: The saltyim.Client.OutboxClient() panics and
	//  is unrecoverable if this method is not called
	//  due to nil pointer reference at saltyim.Client.Outbox
	_, err = saltyim.LookupAddr(s.Identity().Addr().String())
	if err != nil {
		alog.Logger().Errorln(err)
		errCh <- err
		close(errCh)
		return errCh
	}

	s.isClientRunning = true
	go func() {
		defer func() {
			recoverPanic(alog.Logger())
			s.isClientRunning = false
		}()
		var ctx context.Context
		ctx = context.Background()
		inboxCh := s.saltyClient.Subscribe(ctx, "", "", "")
		outboxCh := s.saltyClient.OutboxClient(s.Identity().Addr()).Subscribe(ctx, "", "", "")
		for s.IsClientRunning() {
			select {
			case <-ctx.Done():
				close(inboxCh)
				close(outboxCh)
				return
			case msg := <-inboxCh:
				sl := strings.Fields(msg.Text)
				if len(sl) >= 3 {
					created := sl[0]
					contactAddr := sl[1][1 : len(sl[1])-1]
					dbMsg := Message{
						UserAddr:    s.Identity().Addr().String(),
						ContactAddr: contactAddr,
						From:        contactAddr,
						To:          s.Identity().Addr().String(),
						Created:     created,
						Text:        strings.Join(sl[2:], " "),
						Key:         msg.Key.String(),
					}
					if !s.isMessageDuplicate(dbMsg) {
						s.addMessage(dbMsg)
						s.saveMessage(dbMsg)
					}
				}
			case msg := <-outboxCh:
				alog.Logger().Debugln(msg)
			}
		}
	}()
	return errCh
}

func (s *service) NewChat(addrStr string) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			recoverPanicCloseCh(errCh, err, alog.Logger())
		}()
		var addr *saltyim.Addr
		if s.ContactsAddresses() == nil {
			s.contactsAddresses = make([]*saltyim.Addr, 0)
		} else {
			if cl := s.GetContactAddr(addrStr); cl != nil {
				err = errors.New("client already exist")
				return
			}
		}
		addr, err = saltyim.LookupAddr(addrStr)
		if err != nil {
			alog.Logger().Errorln(err)
			return
		}
		s.contactsAddresses = append(s.contactsAddresses, addr)
		<-s.saveContacts()
		return
	}()
	return errCh
}

func (s *service) IsClientRunning() bool {
	return s.isClientRunning
}

func (s *service) addIdentity(i *saltyim.Identity) {
	if cap(s.Identities()) == 0 {
		s.identities = make([]*saltyim.Identity, 0, 1)
	}
	if i == nil {
		return
	}
	s.identities = append(s.identities, i)
}

func (s *service) addContact(i *saltyim.Addr) {
	if cap(s.Identities()) == 0 {
		s.contactsAddresses = make([]*saltyim.Addr, 1)
	}
	if i == nil {
		return
	}
	s.contactsAddresses = append(s.contactsAddresses, i)
}

func recoverPanic(entry *logrus.Entry) {
	if r := recover(); r != nil {
		entry.Errorln("recovered from panic", r)
	}
}
func recoverPanicCloseCh[S any](stateChan chan<- S, state S, entry *logrus.Entry) {
	recoverPanic(entry)
	stateChan <- state
	close(stateChan)
}
