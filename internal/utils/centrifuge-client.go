package utils

import (
	"fmt"
	"sync"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/rs/zerolog"
)

const (
	DefaultWebsocketURL  string = "ws://localhost:8000/connection/websocket"
	ServerPublishChannel string = "server-general"
)

type ClientOptions struct {
	MessageHandler     centrifuge.MessageHandler
	PublicationHandler centrifuge.ServerPublicationHandler
}

type ClientOption func(options *ClientOptions)

func WithMessageHandler(handler centrifuge.MessageHandler) ClientOption {
	return func(options *ClientOptions) {
		options.MessageHandler = handler
	}
}

func WithServerPublicationHandler(handler centrifuge.ServerPublicationHandler) ClientOption {
	return func(options *ClientOptions) {
		options.PublicationHandler = handler
	}
}

// Note: the waitgroup introduces a bug when server disconnect and reconnect.
// Automated reconnection of the client will try to decrement null waitgroup counter (wg.Done) and
// raise an exception.
func NewClient(log *zerolog.Logger, wsURL string, wg *sync.WaitGroup, opts ...ClientOption) *centrifuge.Client {
	clientOpts := &ClientOptions{}
	for _, opt := range opts {
		opt(clientOpts)
	}

	c := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Name:    "go-client",
		Version: "0.0.1",
	})

	c.OnConnected(func(e centrifuge.ConnectedEvent) {
		log.Debug().Msgf("connected %#v", e)
		wg.Done()
	})

	c.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Debug().Msgf("disconnected event: %d %s", e.Code, e.Reason)
	})

	c.OnError(func(e centrifuge.ErrorEvent) {
		log.Debug().Msgf("error: %s", e.Error.Error())
	})

	var messageHandler = func(e centrifuge.MessageEvent) {
		log.Debug().Msgf("message received from server %s", string(e.Data))
	}

	if clientOpts.MessageHandler != nil {
		messageHandler = clientOpts.MessageHandler
	}

	c.OnMessage(messageHandler)

	var publicationHandler = func(e centrifuge.ServerPublicationEvent) {
		log.Debug().Msgf("publication received from server %s", string(e.Data))
	}

	if clientOpts.PublicationHandler != nil {
		publicationHandler = clientOpts.PublicationHandler
	}

	c.OnPublication(publicationHandler)

	return c
}

type SubscriptionOptions struct {
	PublicationHandler centrifuge.PublicationHandler
	Config             centrifuge.SubscriptionConfig
}

type SubscriptionOption func(options *SubscriptionOptions)

func WithPublicationHandler(handler centrifuge.PublicationHandler) SubscriptionOption {
	return func(options *SubscriptionOptions) {
		options.PublicationHandler = handler
	}
}

func WithSubscriptionConfig(config centrifuge.SubscriptionConfig) SubscriptionOption {
	return func(options *SubscriptionOptions) {
		options.Config = config
	}
}

func Subscribe(log *zerolog.Logger, c *centrifuge.Client, topicName string, opts ...SubscriptionOption) error {
	var wg sync.WaitGroup

	subscriptionOpts := &SubscriptionOptions{}
	for _, opt := range opts {
		opt(subscriptionOpts)
	}

	subscription, err := c.NewSubscription(topicName, subscriptionOpts.Config)
	if err != nil {
		return fmt.Errorf("new subscription to %s error: %s", topicName, err.Error())
	}

	subscription.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Debug().Msgf("[%s] subscription error event: %s", topicName, e.Error.Error())
	})

	var publicationHandler = func(e centrifuge.PublicationEvent) {
		log.Debug().Msgf("[%s] publication event: %s", topicName, string(e.Data))
	}

	if subscriptionOpts.PublicationHandler != nil {
		publicationHandler = subscriptionOpts.PublicationHandler
	}

	subscription.OnPublication(publicationHandler)

	subscription.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Debug().Msgf("[%s] subscribed event", topicName)
		wg.Done()
	})

	wg.Add(1)
	err = subscription.Subscribe()
	if err != nil {
		log.Error().Msgf("subscription error: %s", err.Error())
	}

	wg.Wait()

	return nil
}
