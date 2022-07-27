package events

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/DarioRoman01/cqrs/models"
	"github.com/nats-io/nats.go"
)

// NatsEventStore is an event store that uses NATS to publish events and implements the EventStore interface
type NatsEventStore struct {
	// conn is the NATS connection used to publish events
	conn *nats.Conn
	// feedCreatedSubject is the subject used to publish CreatedFeedMessage
	feedCreatedSub *nats.Subscription
	// feedUpdatedSubject is the subject used to publish UpdatedFeedMessage
	feedUpdatedSub *nats.Subscription
	// feedDeletedSubject is the subject used to publish DeletedFeedMessage
	feedDeletedSub *nats.Subscription
	// feedCreatedChan is the channel used to receive CreatedFeedMessage
	feedCreatedChan chan CreatedFeedMessage
	// feedUpdatedChan is the channel used to receive UpdatedFeedMessage
	feedUpdatedChan chan UpdatedFeedMessage
	// feedDeletedChan is the channel used to receive DeletedFeedMessage
	feedDeletedChan chan DeletedFeedMessage
}

// NewNatsEventStore creates a new NatsEventStore
func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{conn: conn}, nil
}

// Close closes the event store
func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}

	if n.feedUpdatedSub != nil {
		n.feedUpdatedSub.Unsubscribe()
	}

	if n.feedDeletedSub != nil {
		n.feedDeletedSub.Unsubscribe()
	}

	if n.feedCreatedChan != nil {
		close(n.feedCreatedChan)
	}

	if n.feedUpdatedChan != nil {
		close(n.feedUpdatedChan)
	}

	if n.feedDeletedChan != nil {
		close(n.feedDeletedChan)
	}
}

// encodeMessage encodes the message into a byte array
func (n *NatsEventStore) encodeMessage(msg Message) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(msg)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// decodeMessage decodes the message from a byte array
func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	buff := bytes.NewBuffer(data)
	return gob.NewDecoder(buff).Decode(m)
}

// PublishCreatedFeed publishes a CreatedFeedMessage
func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

// OnCreatedFeed registers a handler for CreatedFeedMessage
func (n *NatsEventStore) OnCreatedFeed(handler func(*CreatedFeedMessage)) error {
	msg := CreatedFeedMessage{}
	var err error
	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		handler(&msg)
	})

	if err != nil {
		return err
	}

	return nil
}

// SubscribeCreatedFeed subscribes to the CreatedFeedMessage subject
func (n *NatsEventStore) SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	msg := CreatedFeedMessage{}
	n.feedCreatedChan = make(chan CreatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error

	n.feedCreatedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for m := range ch {
			n.decodeMessage(m.Data, &msg)
			n.feedCreatedChan <- msg
		}
		// for {
		// 	select {
		// 	case m := <-ch:
		// 		n.decodeMessage(m.Data, &msg)
		// 		n.feedCreatedChan <- msg
		// 	}
		// }
	}()

	return (<-chan CreatedFeedMessage)(n.feedCreatedChan), nil
}

// PublishUpdatedFeed publishes a UpdatedFeedMessage
func (n *NatsEventStore) PublishUpdatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := UpdatedFeedMessage{
		ID:          feed.ID,
		Title:       feed.Title,
		Description: feed.Description,
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

// OnUpdatedFeed registers a handler for UpdatedFeedMessage
func (n *NatsEventStore) OnUpdatedFeed(handler func(*UpdatedFeedMessage)) error {
	msg := UpdatedFeedMessage{}
	var err error
	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		handler(&msg)
	})

	if err != nil {
		return err
	}

	return nil
}

// SubscribeUpdatedFeed subscribes to the UpdatedFeedMessage subject
func (n *NatsEventStore) SubscribeUpdatedFeed(ctx context.Context) (<-chan UpdatedFeedMessage, error) {
	msg := UpdatedFeedMessage{}
	n.feedUpdatedChan = make(chan UpdatedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error

	n.feedUpdatedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for m := range ch {
			n.decodeMessage(m.Data, &msg)
			n.feedUpdatedChan <- msg
		}
		// for {
		// 	select {
		// 	case m := <-ch:
		// 		n.decodeMessage(m.Data, &msg)
		// 		n.feedUpdatedChan <- msg
		// 	}
		// }
	}()

	return (<-chan UpdatedFeedMessage)(n.feedUpdatedChan), nil
}

// PublishDeletedFeed publishes a DeletedFeedMessage
func (n *NatsEventStore) PublishDeletedFeed(ctx context.Context, feed *models.Feed) error {
	msg := DeletedFeedMessage{
		ID:        feed.ID,
		DeletedAt: time.Now(),
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

// OnDeletedFeed registers a handler for DeletedFeedMessage
func (n *NatsEventStore) OnDeletedFeed(handler func(*DeletedFeedMessage)) error {
	msg := DeletedFeedMessage{}
	var err error
	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		handler(&msg)
	})

	if err != nil {
		return err
	}

	return nil
}

// SubscribeDeletedFeed subscribes to the DeletedFeedMessage subject
func (n *NatsEventStore) SubscribeDeletedFeed(ctx context.Context) (<-chan DeletedFeedMessage, error) {
	msg := DeletedFeedMessage{}
	n.feedDeletedChan = make(chan DeletedFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error

	n.feedDeletedSub, err = n.conn.ChanSubscribe(msg.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for m := range ch {
			n.decodeMessage(m.Data, &msg)
			n.feedDeletedChan <- msg
		}
		// for {
		// 	select {
		// 	case m := <-ch:
		// 		n.decodeMessage(m.Data, &msg)
		// 		n.feedDeletedChan <- msg
		// 	}
		// }
	}()

	return (<-chan DeletedFeedMessage)(n.feedDeletedChan), nil
}
