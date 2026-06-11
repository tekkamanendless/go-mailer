package mailer

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

// dummyMailServer is the dummy mail server.
var dummyMailServer DummyMailServer

// Dummy returns the dummy mail server.
//
// This is a singleton instance of the dummy mail server.
func Dummy() *DummyMailServer {
	return &dummyMailServer
}

// DummyMailServer is our "dummy" mail server.
// It is indexed by the API key and contains a list of messages in ascending time order.
type DummyMailServer struct {
	messagesByAPIKey map[string][]*Message
	mutex            sync.Mutex
}

// Send a message via the dummy mail server.
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

// APIKeys returns the API keys used by the dummy mail server so far.
func (d *DummyMailServer) APIKeys() []string {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	keys := slices.Collect(maps.Keys(d.messagesByAPIKey))
	if keys == nil {
		return []string{}
	}
	slices.Sort(keys)
	return keys
}

// Addresses returns the list of e-mail addresses available for the given API key.
func (d *DummyMailServer) Addresses(apiKey string) []string {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	addressMap := map[string]bool{}
	for _, message := range d.messagesByAPIKey[apiKey] {
		addressMap[message.From.Address] = true
		addressMap[message.To.Address] = true
	}
	addresses := slices.Collect(maps.Keys(addressMap))
	if addresses == nil {
		return []string{}
	}
	slices.Sort(addresses)
	return addresses
}

// Outbox returns the messages in the outbox for the given API key and address.
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

// Inbox returns the messages in the inbox for the given API key and address.
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

// LastMessageInOutbox returns the last message in the outbox for the given API key and address.
// If there are no messages in the outbox, it returns nil.
func (d *DummyMailServer) LastMessageInOutbox(apiKey string, address string) *Message {
	messages := d.Outbox(apiKey, address)
	if len(messages) > 0 {
		return messages[len(messages)-1]
	}
	return nil
}

// LastMessageInInbox returns the last message in the inbox for the given API key and address.
// If there are no messages in the inbox, it returns nil.
func (d *DummyMailServer) LastMessageInInbox(apiKey string, address string) *Message {
	messages := d.Inbox(apiKey, address)
	if len(messages) > 0 {
		return messages[len(messages)-1]
	}
	return nil
}
