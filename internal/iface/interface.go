package iface

import "context"

type MessageHandler interface {
	HandleMessage(msg []byte, id string) error
}

type MessagePublisher interface {
	Publish(message interface{}) error
}

type MessageSubscriber interface {
	Subscribe(ctx context.Context, handler MessageHandler) error
}

// type MessageRepository interface {
// 	CreateNew(message *domain.ParsedMessage) error
// }
