package service

const (
	DBPathCfgName = "salty"
	// DBPathIDsDir contains all the user accounts
	DBPathIDsDir = "identities"
	// DBPathCurrIDDir contains the current account in the ui
	DBPathCurrIDDir = "identity"
	// DBPathContactsDir contains account directory to which it belongs to and
	// then all the contacts inside that directory
	DBPathContactsDir = "contacts"
	// DBPathMessagesDir contains account directory which in turn
	// contains contact directory to which it belongs to and
	// then all the messages inside that directory
	DBPathMessagesDir = "messages"
)
const AppRegisterEndPoint = "https://salty.mills.io"
const BaseURL = "https://salty.mills.io"
const AppSendEndPoint = BaseURL + "/api/v1/send"
const AppLookupEndPoint = BaseURL + "/api/v1/lookup"

const AppIndexedDBVersion int64 = 1
