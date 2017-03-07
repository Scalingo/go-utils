package nsqconsumer

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nsqio/go-nsq"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
)

type NsqMessageDeserialize struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type nsqConsumer struct {
	NsqConfig      *nsq.Config
	NsqLookupdURLs []string
	Topic          string
	Channel        string
	MessageHandler func(*NsqMessageDeserialize) error
	MaxInFlight    int
	count          uint64
}

type ConsumerOpts struct {
	NsqConfig      *nsq.Config
	NsqLookupdURLs []string
	Topic          string
	Channel        string
	MaxInFlight    int
	MessageHandler func(*NsqMessageDeserialize) error
}

type Consumer interface {
	Start() func()
}

func New(opts ConsumerOpts) (Consumer, error) {
	consumer := &nsqConsumer{
		NsqConfig:      opts.NsqConfig,
		NsqLookupdURLs: opts.NsqLookupdURLs,
		Topic:          opts.Topic,
		Channel:        opts.Channel,
		MessageHandler: opts.MessageHandler,
		MaxInFlight:    opts.MaxInFlight,
	}
	if consumer.MaxInFlight == 0 {
		consumer.MaxInFlight = opts.NsqConfig.MaxInFlight
	}
	if opts.Topic == "" {
		return nil, errgo.New("topic can't be blank")
	}
	if opts.MessageHandler == nil {
		return nil, errgo.New("message handler can't be blank")
	}
	if opts.Channel == "" {
		consumer.Channel = "default"
	}
	return consumer, nil
}

func (c *nsqConsumer) Start() func() {
	logger := log.New(os.Stdout, "[nsq-consumer] ", log.Flags())
	logger.Println("Start the worker")

	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, c.NsqConfig)
	if err != nil {
		rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
		logger.Fatalf("fail to create new NSQ consumer: %+v\n", err)
	}
	consumer.SetLogger(logger, nsq.LogLevelWarning)

	consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) (err error) {
		defer func() {
			if errRecovered := recover(); err != nil {
				err = errgo.Newf("recover panic from nsq consumer: %+v", errRecovered)
			}
		}()
		if len(message.Body) == 0 {
			errMsg := "body is blank, re-enqueued message"
			logger.Printf("%s\n", errMsg)
			err := errgo.New(errMsg)
			rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			return err
		}
		var msg NsqMessageDeserialize
		err = json.Unmarshal(message.Body, &msg)
		if err != nil {
			rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			logger.Printf("Failed to unmarshal message: %+v\n", err)
			return errgo.Mask(err, errgo.Any)
		}

		logger.Printf("New message: '%s'", msg.Type)
		err = c.MessageHandler(&msg)
		if err != nil {
			rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			logger.Printf("Failed to handle the message: %+v\n", err)
			return errgo.Mask(err, errgo.Any)
		}
		return nil
	}), c.MaxInFlight)

	if err = consumer.ConnectToNSQLookupds(c.NsqLookupdURLs); err != nil {
		rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "authenticator"})
		logger.Fatalf("Fail to connect to NSQ lookupd: %+v\n", err)
	}

	return func() {
		consumer.Stop()
		// block until stop process is complete
		<-consumer.StopChan
	}
}
