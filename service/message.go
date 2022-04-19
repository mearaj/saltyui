package service

// Message properties UserAddr and ContactAddr keyConfigMessages
type Message struct {
	UserAddr    string
	ContactAddr string
	From        string
	To          string
	Created     string
	Text        string
	Key         string
}

func (m *Message) isDuplicate(msg Message) bool {
	return m.UserAddr == msg.UserAddr &&
		m.ContactAddr == msg.ContactAddr &&
		m.Key == msg.Key &&
		m.Text == msg.Text &&
		m.Created == msg.Created &&
		m.ContactAddr == msg.ContactAddr
}
