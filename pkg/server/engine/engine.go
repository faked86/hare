package engine

import (
	"hare/pkg/transport"
	"sync"
)

const minTopicLen int = 10

type Topic struct {
	// messages  [][]byte
	consumers []chan []byte
}

type Engine struct {
	mu     sync.RWMutex
	topics map[string]*Topic
}

func New() *Engine {
	return &Engine{topics: make(map[string]*Topic)}
}

func (e *Engine) Publish(msg transport.Message) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, exists := e.topics[string(msg.Topic)]
	if !exists {
		e.topics[string(msg.Topic)] = &Topic{
			// messages:  make([][]byte, 0, minTopicLen),
			consumers: make([]chan []byte, 0, minTopicLen),
		}
	}
	// e.topics[topicName].messages = append(e.topics[string(msg.Topic)].messages, payload)
	for _, ch := range e.topics[string(msg.Topic)].consumers {
		ch <- transport.EncodeMessage(msg)
	}
}

func (e *Engine) Subscribe(topicName string) chan []byte {
	e.mu.Lock()
	defer e.mu.Unlock()

	ch := make(chan []byte, 100)

	_, exists := e.topics[topicName]
	if !exists {
		e.topics[topicName] = &Topic{
			// messages:  make([][]byte, 0, minTopicLen),
			consumers: make([]chan []byte, 0, minTopicLen),
		}
	}
	e.topics[topicName].consumers = append(e.topics[topicName].consumers, ch)

	return ch
}

func (e *Engine) Unsubscribe(topicName string, ch chan []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delId := -1
	for i, tCh := range e.topics[topicName].consumers {
		if tCh == ch {
			delId = i
			break
		}
	}
	if delId > -1 {
		e.topics[topicName].consumers = append(e.topics[topicName].consumers[:delId], e.topics[topicName].consumers[delId+1:]...)
	}
}

// func (e *Engine) GetMessages(topicName string) [][]byte {
// 	e.mu.RLock()
// 	defer e.mu.RUnlock()

// 	r, exists := e.topics[topicName]
// 	if !exists {
// 		return nil
// 	}
// 	return r.messages
// }
