package nats

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"casstm-dashboard/internal/config"
	"casstm-dashboard/internal/iface"
	"casstm-dashboard/pkg/utils"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	nc "github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	config     *config.Config
	subscriber *nats.Subscriber
	workers    int
}

const defaultWorkerCount = 4

func NewSub(config *config.Config) (*NatsSubscriber, error) {
	logger := watermill.NewStdLogger(false, false)
	marshaler := &PlainTextMarshaler{}
	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(config.Timeouts.Server),
		nc.ReconnectWait(config.Timeouts.ReconnectWait),
	}
	jsConfig := nats.JetStreamConfig{Disabled: true}
	subscriber, err := nats.NewSubscriber(
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
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS subscriber: %w", err)
	}
	return &NatsSubscriber{
		config:     config,
		subscriber: subscriber,
		workers:    defaultWorkerCount,
	}, nil
}

func (n *NatsSubscriber) Subscribe(ctx context.Context, handler iface.MessageHandler) error {
	if n.subscriber == nil {
		return errors.New("nats subscriber not initialised")
	}
	logger := utils.GetSugaredLogger()
	logger.Infof("Connecting to NATS server: %s", n.config.Nats.URL)
	logger.Infof("Subscribing to topic: %s", n.config.Subscription.Topic)
	defer func() {
		if err := n.subscriber.Close(); err != nil {
			logger.Warnf("error closing NATS subscriber: %v", err)
		}
	}()

	messages, err := n.subscriber.Subscribe(ctx, n.config.Subscription.Topic)
	if err != nil {
		logger.Errorf("Failed to subscribe to topic: %v", err)
		return err
	}

	workerCount := n.workers
	if workerCount <= 0 {
		workerCount = 1
	}

	messageCh := make(chan *message.Message)
	var wg sync.WaitGroup
	workerCtx, cancelWorkers := context.WithCancel(ctx)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-workerCtx.Done():
					return
				case msg, ok := <-messageCh:
					if !ok {
						return
					}
					if err := handler.HandleMessage(msg.Payload, msg.UUID); err == nil {
						logger.Infof("Message handled: %s", msg.UUID)
						msg.Ack()
					} else {
						logger.Errorf("Failed to handle message [%s]: %v", msg.UUID, err)
						msg.Nack()
					}
				}
			}
		}(i)
	}

	for {
		select {
		case <-ctx.Done():
			cancelWorkers()
			close(messageCh)
			wg.Wait()
			return ctx.Err()
		case msg, ok := <-messages:
			if !ok {
				cancelWorkers()
				close(messageCh)
				wg.Wait()
				return nil
			}

			select {
			case messageCh <- msg:
			case <-ctx.Done():
				cancelWorkers()
				close(messageCh)
				wg.Wait()
				return ctx.Err()
			}
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
