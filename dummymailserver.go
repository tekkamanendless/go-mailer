package mailer

import (
	"fmt"
	"sync"
)

// DummyMailServer is our "dummy" mail server.
// It is indexed by the API key and contains a list of messages in ascending time order.
type DummyMailServer struct {
	messagesByAPIKey map[string][]*Message
	mutex            sync.Mutex
}

func (d *DummyMailServer) Send(apiKey string, message Message) error {
	if apiKey == "" {
		return fmt.Errorf("invalid API key")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.messagesByAPIKey == nil {
		d.messagesByAPIKey = map[string][]*Message{}
	}
	d.messagesByAPIKey[apiKey] = append(d.messagesByAPIKey[apiKey], &message)

	return nil
}

func (d *DummyMailServer) Outbox(apiKey string, address string) []*Message {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.messagesByAPIKey == nil {
		d.messagesByAPIKey = map[string][]*Message{}
	}

	messages := []*Message{}
	for _, message := range d.messagesByAPIKey[apiKey] {
		if message.From.Address == address {
			messages = append(messages, message)
		}
	}
	return messages
}

func (d *DummyMailServer) Inbox(apiKey string, address string) []*Message {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.messagesByAPIKey == nil {
		d.messagesByAPIKey = map[string][]*Message{}
	}

	messages := []*Message{}
	for _, message := range d.messagesByAPIKey[apiKey] {
		if message.To.Address == address {
			messages = append(messages, message)
		}
	}
	return messages
}

func (d *DummyMailServer) LastMessageInOutbox(apiKey string, address string) *Message {
	messages := d.Outbox(apiKey, address)
	if len(messages) > 0 {
		return messages[len(messages)-1]
	}
	return nil
}

func (d *DummyMailServer) LastMessageInInbox(apiKey string, address string) *Message {
	messages := d.Inbox(apiKey, address)
	if len(messages) > 0 {
		return messages[len(messages)-1]
	}
	return nil
}

var dummyMailServer DummyMailServer

// Dummy returns the dummy mail server.
func Dummy() *DummyMailServer {
	return &dummyMailServer
}
