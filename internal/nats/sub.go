package nats

import (
	"casstm-dashboard/internal/iface"
	"casstm-dashboard/pkg/utils"
	"context"
	"errors"

	"casstm-dashboard/internal/config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	config     *config.Config
	subscriber *nats.Subscriber
}

func NewSub(config *config.Config) *NatsSubscriber {
	logger := watermill.NewStdLogger(false, false)
	marshaler := &PlainTextMarshaler{}
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(config.Timeouts.Server),
		nc.ReconnectWait(config.Timeouts.ReconnectWait),
	}
	jsConfig := nats.JetStreamConfig{Disabled: true}
	subscriber, _ := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:            config.Nats.URL,
			CloseTimeout:   config.Timeouts.Close,
			AckWaitTimeout: config.Timeouts.AckWait,
			NatsOptions:    options,
			Unmarshaler:    marshaler,
			JetStream:      jsConfig,
		},
		logger,
	)
	return &NatsSubscriber{
		config:     config,
		subscriber: subscriber,
	}
}

func (n *NatsSubscriber) Subscribe(handlers iface.MessageHandler) {
	logger := utils.GetSugaredLogger()
	logger.Infof("Connecting to NATS server: %s", n.config.Nats.URL)
	logger.Infof("Subscribing to topic: %s", n.config.Subscription.Topic)
	defer n.subscriber.Close()
	messages, err := n.subscriber.Subscribe(context.Background(), n.config.Subscription.Topic)
	if err != nil {
		logger.Errorf("Failed to subscribe to topic: %v", err)
		return
	}

	for msg := range messages {
		if err := handlers.HandleMessage(msg.Payload, msg.UUID); err == nil {
			logger.Infof("Message handled: %s", msg.UUID)
			msg.Ack()
		} else {
			logger.Errorf("Failed to handle message [%s]: %v", msg.UUID, err)
			msg.Nack()
		}
	}
}

type PlainTextMarshaler struct{}

func (m *PlainTextMarshaler) Marshal(topic string, msg nc.Msg) ([]byte, error) {
	return msg.Data, nil
}

func (m *PlainTextMarshaler) Unmarshal(newMsg *nc.Msg) (*message.Message, error) {
	if newMsg == nil {
		return nil, errors.New("empty message")
	}
	msg := message.NewMessage(watermill.NewUUID(), newMsg.Data)
	return msg, nil
}
