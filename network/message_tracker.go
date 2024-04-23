package network

import (
	"errors"
)

// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker, deleting the oldest message if necessary
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Message returns a message for a given ID.  Message is retained in tracker
	Message(id string) (message *Message, err error)
	// Messages returns messages in FIFO order
	Messages() (messages []*Message)
}

// ErrMessageNotFound is an error returned by MessageTracker when a message with specified id is not found
var ErrMessageNotFound = errors.New("message not found")

func NewMessageTracker(length int) MessageTracker {
	return &MessageTrackerImpl{
		messages:   make([]*Message, 0, length),
		messageMap: make(map[string]*Message),
		capacity:   length,
	}
}

// MessageTrackerImpl implements MessageTracker interface.
type MessageTrackerImpl struct {
	// The message queue as an array.
	messages []*Message
	// The hash table for random access operations.
	messageMap map[string]*Message
	capacity   int
}

func (mt *MessageTrackerImpl) Add(message *Message) (err error) {
	// If message with the ID exists then ignore the request.
	if _, ok := mt.messageMap[message.ID]; ok {
		return nil
	}

	// If the buffer is full then remove the first (oldest) message.
	if len(mt.messages) == mt.capacity {
		delete(mt.messageMap, mt.messages[0].ID)
		mt.messages = mt.messages[1:]
	}

	// Add the new message
	mt.messages = append(mt.messages, message)
	mt.messageMap[message.ID] = message

	return nil
}

func (mt *MessageTrackerImpl) Delete(id string) (err error) {
	// Return error if message does not exist.
	if _, ok := mt.messageMap[id]; !ok {
		return ErrMessageNotFound
	}

	// Delete the message from the map.
	delete(mt.messageMap, id)
	// Find the message, delete it and shrink the queue.
	for i, v := range mt.messages {
		if v.ID == id {
			mt.messages = append(mt.messages[:i], mt.messages[i+1:]...)
			break
		}
	}

	return nil
}

func (mt *MessageTrackerImpl) Message(id string) (message *Message, err error) {
	// Find the message in the map.
	message, ok := mt.messageMap[id]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return message, nil
}

func (mt *MessageTrackerImpl) Messages() (messages []*Message) {
	// The message queue is always up-to-date, just return it.
	return mt.messages
}
