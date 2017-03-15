package nsqconsumer

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/Soulou/errgo-rollbar"
	"github.com/nsqio/go-nsq"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
)

type NsqMessageDeserialize struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	NsqMsg  *nsq.Message
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
	// How long can the consumer keep the message before the message is considered as 'Timed Out'
	MsgTimeout     time.Duration
	MessageHandler func(*NsqMessageDeserialize) error
}

type Consumer interface {
	Start() func()
}

func New(opts ConsumerOpts) (Consumer, error) {
	if opts.MsgTimeout != 0 {
		opts.NsqConfig.MsgTimeout = opts.MsgTimeout
	}

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
			if errRecovered := recover(); errRecovered != nil {
				err = errgo.Newf("recover panic from nsq consumer: %+v", errRecovered)
				logger.Printf("[PANIC] %v: %v", err, errgo.Details(errRecovered.(error)))
				rollbar.Error(rollbar.ERR, errRecovered.(error), &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
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

		logger.Printf("[%s] BEGIN Message: '%s'", message.ID, msg.Type)
		msg.NsqMsg = message
		err = c.MessageHandler(&msg)
		if err != nil {
			rollbar.ErrorWithStack(rollbar.ERR, err, errgorollbar.BuildStack(err), &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			logger.Printf("[%s] ERROR: %+v\n", message.ID, errgo.Details(err))
			return errgo.Mask(err, errgo.Any)
		}
		logger.Printf("[%s] END Message: '%s'", message.ID, msg.Type)
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
