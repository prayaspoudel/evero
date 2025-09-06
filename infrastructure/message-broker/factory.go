package messagebroker

import (
	"errors"
)

var (
	errInvalidBrokerInstance = errors.New("invalid message broker instance")
	errBrokerNotConnected    = errors.New("message broker not connected")
	errTopicNotFound         = errors.New("topic not found")
	errSubscriptionNotFound  = errors.New("subscription not found")
	errInvalidMessage        = errors.New("invalid message format")
	errPublishFailed         = errors.New("failed to publish message")
	errSubscribeFailed       = errors.New("failed to subscribe to topic")
)

const (
	InstanceRabbitMQ int = iota
	InstanceNATS
	InstanceKafka
)

// NewMessageBrokerFactory creates a new message broker instance based on the specified type
func NewMessageBrokerFactory(instance int, config *BrokerConfig) (MessageBroker, error) {
	switch instance {
	case InstanceRabbitMQ:
		return NewRabbitMQBroker(config)
	case InstanceNATS:
		return NewNATSBroker(config)
	case InstanceKafka:
		return NewKafkaBroker(config)
	default:
		return nil, errInvalidBrokerInstance
	}
}
