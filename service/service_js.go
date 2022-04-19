package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mearaj/saltyui/alog"
	"go.mills.io/saltyim"
	"syscall/js"
)

var DBPaths = map[string]string{
	DBPathIDsDir:      DBPathIDsDir,
	DBPathCurrIDDir:   DBPathCurrIDDir,
	DBPathContactsDir: DBPathContactsDir,
	DBPathMessagesDir: DBPathMessagesDir,
}

// Service Always call NewService function to create Service
type service struct {
	identity          *saltyim.Identity
	isRegistered      bool
	identities        []*saltyim.Identity
	contactsAddresses []*saltyim.Addr
	saltyClient       *saltyim.Client
	indexedDB         js.Value
	initialized       bool
	loaded            bool
	isClientRunning   bool
	configJSON        *ConfigJSON
	// string refers to contact addr
	messages map[string][]Message
}

type KeyPath struct {
	KeyPath string
}

func init() {
	saltyim.SetResolver(&saltyim.DNSOverHTTPResolver{})
}

// config used by identities, identity and contacts
var keyConfigCommon = js.ValueOf(map[string]interface{}{"keyPath": "Id"})
var keyConfigMessages = js.ValueOf(map[string]interface{}{
	"keyPath": []interface{}{"UserAddr", "ContactAddr"},
})

func (s *service) IsRegistered() bool {
	return s.Identity() != nil && s.isRegistered
}

func (s *service) Loaded() bool {
	return s.loaded && s.initialized
}
func (s *service) init() {
	defer func() {
		recoverPanic(alog.Logger())
	}()
	indexedDBReq := js.Global().Get("indexedDB")
	req := indexedDBReq.Call("open", js.ValueOf(DBPathCfgName), js.ValueOf(AppIndexedDBVersion))
	req.Set("onsuccess", js.FuncOf(s.onInitSuccess))
	req.Set("onupgradeneeded", js.FuncOf(s.onUpgradeNeeded))
	req.Set("onerror", js.FuncOf(s.onInitError))
}

func (s *service) onInitSuccess(this js.Value, args []js.Value) interface{} {
	s.indexedDB = args[0].Get("target").Get("result")
	s.initialized = true
	s.loadDatabase()
	return nil
}

func (s *service) onUpgradeNeeded(this js.Value, args []js.Value) interface{} {
	defer func() {
		recoverPanic(alog.Logger())
	}()
	s.indexedDB = args[0].Get("target").Get("result")
	currentVersion := s.indexedDB.Get("version").Int()
	oldVersion := args[0].Get("oldVersion").Int()
	newVersion := args[0].Get("newVersion").Int() // equivalent to currentVersion
	if oldVersion < 1 {
		// The app is installed first time
		_ = currentVersion
		_ = newVersion
	}
	// Currently, regardless of the version, check if store items exists and if not then create one
	objectStoreNames := s.indexedDB.Get("objectStoreNames")
	for k := range DBPaths {
		if !objectStoreNames.Call("contains", js.ValueOf(k)).Bool() {
			if k == DBPathMessagesDir {
				s.indexedDB.Call("createObjectStore", js.ValueOf(k), keyConfigMessages)
			} else {
				s.indexedDB.Call("createObjectStore", js.ValueOf(k), keyConfigCommon)
			}
		}
	}
	return nil
}
func (s *service) onInitError(this js.Value, args []js.Value) interface{} {
	errorJs := args[0].Get("target").Get("errorCode")
	alog.Logger().Error(errorJs)
	return nil
}

func (s *service) loadDatabase() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		<-s.loadIdentities()
		<-s.loadCurrentIdentity()
		<-s.loadContacts()
		<-s.loadMessages()
		s.loaded = true
	}()
	return errCh
}

func (s *service) loadIdentities() <-chan error {
	errCh := make(chan error, 1)
	var err error
	go func() {
		defer func() { recoverPanic(alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathIDsDir), "readonly")
		objStore := txn.Call("objectStore", js.ValueOf(DBPathIDsDir))
		req := objStore.Call("getAll")
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
			accountsJS := args[0].Get("target").Get("result")
			if accountsJS.Truthy() {
				sizeOFArray := accountsJS.Get("length").Int()
				for i := 0; i < sizeOFArray; i++ {
					eachAccountJS := accountsJS.Index(i)
					addrContents := eachAccountJS.Get("Contents").String()
					idOption := saltyim.WithIdentityBytes([]byte(addrContents))
					var saltyId *saltyim.Identity
					saltyId, err = saltyim.GetIdentity(idOption)
					if err != nil {
						alog.Logger().Errorln(err)
					}
					s.addIdentity(saltyId)
				}
			}
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() {
				recoverPanicCloseCh(errCh, err, alog.Logger())
			}()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error in retrieving current id, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
	}()
	return errCh
}

func (s *service) saveIdentities() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathIDsDir), "readwrite")
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathIDsDir))
		for _, addr := range s.Identities() {
			req := objectStore.Call("put", map[string]interface{}{
				"Id":       addr.Addr().String(),
				"Contents": string(addr.Contents()),
			})
			waitCh := make(chan error, 1)
			req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
				defer func() { recoverPanicCloseCh(waitCh, err, alog.Logger()) }()
				err = nil
				return nil
			}))
			req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
				defer func() { recoverPanicCloseCh(waitCh, err, alog.Logger()) }()
				errorJs := args[0].Get("target").Get("errorCode")
				errStr := fmt.Sprintf("error in retrieving current id, errCode is %s", errorJs.String())
				err = errors.New(errStr)
				alog.Logger().Println(errStr)
				return nil
			}))
			<-waitCh
		}
	}()
	return errCh
}

//
func (s *service) loadCurrentIdentity() <-chan error {
	errCh := make(chan error, 1)
	var err error
	go func() {
		defer func() { recoverPanic(alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathCurrIDDir), "readonly")
		objStore := txn.Call("objectStore", js.ValueOf(DBPathCurrIDDir))
		req := objStore.Call("getAll")
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() {
				recoverPanicCloseCh(errCh, err, alog.Logger())
			}()
			accountsJS := args[0].Get("target").Get("result")
			sizeOFArray := accountsJS.Get("length").Int()
			if sizeOFArray > 0 {
				eachAccountJS := accountsJS.Index(0)
				addrContents := eachAccountJS.Get("Contents").String()
				idOption := saltyim.WithIdentityBytes([]byte(addrContents))
				var identity *saltyim.Identity
				identity, err = saltyim.GetIdentity(idOption)
				s.setIdentity(identity)
			}
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() {
				recoverPanicCloseCh(errCh, err, alog.Logger())
			}()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error in retrieving current id, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			s.setIdentity(nil)
			return nil
		}))
	}()
	return errCh
}

//
func (s *service) loadIdentity(iDAddr string) (<-chan *saltyim.Identity, <-chan error) {
	errCh := make(chan error, 1)
	iDCh := make(chan *saltyim.Identity, 1)
	var err error
	var identity *saltyim.Identity
	go func() {
		defer func() { recoverPanic(alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathCurrIDDir), "readonly")
		objStore := txn.Call("objectStore", js.ValueOf(DBPathCurrIDDir))
		req := objStore.Call("get", iDAddr)
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() {
				recoverPanicCloseCh(errCh, err, alog.Logger())
				recoverPanicCloseCh(iDCh, identity, alog.Logger())
			}()
			accountsJS := args[0].Get("target").Get("result")
			sizeOFArray := accountsJS.Get("length").Int()
			if sizeOFArray > 0 {
				eachAccountJS := accountsJS.Index(0)
				addrContents := eachAccountJS.Get("Contents").String()
				idOption := saltyim.WithIdentityBytes([]byte(addrContents))
				identity, err = saltyim.GetIdentity(idOption)
			}
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() {
				recoverPanicCloseCh(errCh, err, alog.Logger())
				recoverPanicCloseCh(iDCh, identity, alog.Logger())
			}()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error in retrieving current id, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
	}()
	return iDCh, errCh
}

// saveIdentity
func (s *service) saveIdentity(identity *saltyim.Identity) <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil || identity == nil {
		err = errors.New("current id or provided id is nil")
		recoverPanicCloseCh(errCh, err, alog.Logger())
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathCurrIDDir), js.ValueOf("readwrite"))
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathCurrIDDir))
		req := objectStore.Call("clear")
		val := map[string]interface{}{
			"Id":       identity.Addr().String(),
			"Contents": string(identity.Contents()),
		}
		req = objectStore.Call("put", val)
		errCh2 := make(chan error, 1)
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error saving current id, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
		err = <-errCh2
	}()
	return errCh
}

// saveContacts replaces all previous contacts with current Contacts
func (s *service) saveContacts() <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil {
		err = errors.New("cannot save contacts, current id is nil")
	}
	if len(s.ContactsAddresses()) == 0 {
		err = errors.New("there are no contacts to save")
	}
	if err != nil {
		recoverPanicCloseCh(errCh, err, alog.Logger())
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathContactsDir), "readwrite")
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathContactsDir))
		userAddr := s.Identity().Addr().String()
		contactsJsArr := js.ValueOf([]interface{}{})
		for _, addr := range s.ContactsAddresses() {
			contactsJsArr.Call("push", addr.String())
		}
		req := objectStore.Call("put", js.ValueOf(map[string]interface{}{
			"Id":       userAddr,
			"Contacts": contactsJsArr,
		}))
		errCh2 := make(chan error, 1)
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error saving contacts, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
		<-errCh2
	}()
	return errCh
}

func (s *service) loadContacts() <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil {
		err = errors.New("cannot load contacts, current id is nil")
		recoverPanicCloseCh(errCh, err, alog.Logger())
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathContactsDir), "readwrite")
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathContactsDir))
		userAddr := s.Identity().Addr().String()
		req := objectStore.Call("get", userAddr)
		errCh2 := make(chan error, 1)
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			contactsJs := args[0].Get("target").Get("result")
			if contactsJs.Truthy() {
				contactsJsArr := contactsJs.Get("Contacts")
				sizeOFArray := contactsJsArr.Get("length").Int()
				for i := 0; i < sizeOFArray; i++ {
					addrJs := contactsJsArr.Index(i)
					var addr *saltyim.Addr
					addr, err = saltyim.ParseAddr(addrJs.String())
					if err != nil {
						alog.Logger().Errorln(err)
					}
					s.addContact(addr)
				}
			}
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error loading contacts, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
		<-errCh2
	}()
	return errCh
}
func (s *service) saveMessage(message Message) <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil {
		err = errors.New("cannot save message, current id is nil")
		recoverPanicCloseCh(errCh, err, alog.Logger())
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathMessagesDir), "readwrite")
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathMessagesDir))
		messagesArr := js.ValueOf([]interface{}{})
		var contactAddr = message.ContactAddr
		for _, msg := range s.Messages(contactAddr) {
			var msgStruct map[string]interface{}
			var data []byte
			data, err = json.Marshal(msg)
			if err != nil {
				alog.Logger().Errorln(err)
			}
			err = json.Unmarshal(data, &msgStruct)
			if err != nil {
				alog.Logger().Errorln(err)
			}
			messagesArr.Call("push", msgStruct)
		}
		req := objectStore.Call("put", js.ValueOf(map[string]interface{}{
			"UserAddr":    message.UserAddr,
			"ContactAddr": message.ContactAddr,
			"From":        message.From,
			"To":          message.To,
			"Messages":    messagesArr,
		}))
		errCh2 := make(chan error, 1)
		req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			alog.Logger().Debugln("message saved successfully")
			return nil
		}))
		req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
			defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
			errorJs := args[0].Get("target").Get("errorCode")
			errStr := fmt.Sprintf("error saving contacts, errCode is %s", errorJs.String())
			err = errors.New(errStr)
			alog.Logger().Println(errStr)
			return nil
		}))
		<-errCh2
	}()
	return errCh
}

func (s *service) loadMessages() <-chan error {
	errCh := make(chan error, 1)
	var err error
	if s.Identity() == nil {
		err = errors.New("cannot save message, current id is nil")
		recoverPanicCloseCh(errCh, err, alog.Logger())
		return errCh
	}
	go func() {
		defer func() { recoverPanicCloseCh(errCh, err, alog.Logger()) }()
		txn := s.indexedDB.Call("transaction", js.ValueOf(DBPathMessagesDir), "readwrite")
		objectStore := txn.Call("objectStore", js.ValueOf(DBPathMessagesDir))
		userAddr := s.Identity().Addr().String()
		for _, eachAddr := range s.ContactsAddresses() {
			req := objectStore.Call("get", []interface{}{userAddr, eachAddr.String()})
			errCh2 := make(chan error, 1)
			req.Set("onsuccess", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
				messagesJS := args[0].Get("target").Get("result")
				if messagesJS.Truthy() {
					messagesArr := messagesJS.Get("Messages")
					if messagesArr.Truthy() {
						sizeOFArray := messagesArr.Get("length").Int()
						for i := 0; i < sizeOFArray; i++ {
							msg := messagesArr.Index(i)
							s.addMessage(Message{
								UserAddr:    userAddr,
								ContactAddr: eachAddr.String(),
								From:        msg.Get("From").String(),
								To:          msg.Get("To").String(),
								Created:     msg.Get("Created").String(),
								Text:        msg.Get("Text").String(),
								Key:         msg.Get("Key").String(),
							})
						}
					}
				}
				return nil
			}))
			req.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) any {
				defer func() { recoverPanicCloseCh(errCh2, err, alog.Logger()) }()
				errorJs := args[0].Get("target").Get("errorCode")
				errStr := fmt.Sprintf("error saving contacts, errCode is %s", errorJs.String())
				err = errors.New(errStr)
				alog.Logger().Println(errStr)
				return nil
			}))
			<-errCh2
		}
	}()
	return errCh
}
