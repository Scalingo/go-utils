package nsqconsumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Scalingo/go-internal-tools/nsqproducer"
	"github.com/Sirupsen/logrus"
	"github.com/Soulou/errgo-rollbar"
	"github.com/nsqio/go-nsq"
	"github.com/stvp/rollbar"
	"gopkg.in/errgo.v1"
)

const (
	// defaultChannel is the name of the channel we're using when we want the
	// message to be receive only by 1 consumer, but no matter which one
	defaultChannel = "default"
)

var (
	maxPostponeDelay int64 = 3600
)

type NsqMessageDeserialize struct {
	Type    string          `json:"type"`
	At      int64           `json:"at"`
	Payload json.RawMessage `json:"payload"`
	NsqMsg  *nsq.Message
	logger  logrus.FieldLogger
}

func (msg *NsqMessageDeserialize) Logger() logrus.FieldLogger {
	if msg.logger == nil {
		msg.logger = logrus.New().WithFields(logrus.Fields{"message-type": msg.Type})
	}
	return msg.logger
}

type nsqConsumer struct {
	NsqConfig        *nsq.Config
	NsqLookupdURLs   []string
	Topic            string
	Channel          string
	MessageHandler   func(*NsqMessageDeserialize) error
	MaxInFlight      int
	PostponeProducer nsqproducer.Producer
	count            uint64
	logger           logrus.FieldLogger
}

type ConsumerOpts struct {
	NsqConfig      *nsq.Config
	NsqLookupdURLs []string
	Topic          string
	Channel        string
	MaxInFlight    int
	// PostponeProducer is an NSQ producer user to send postponed messages
	PostponeProducer nsqproducer.Producer
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
		consumer.Channel = defaultChannel
	}
	return consumer, nil
}

func (c *nsqConsumer) Start() func() {
	c.logger = logrus.New().WithFields(logrus.Fields{
		"source":  "nsq-consumer",
		"topic":   c.Topic,
		"channel": c.Channel,
	})
	c.logger.Println("starting consumer")

	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, c.NsqConfig)
	if err != nil {
		rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
		c.logger.WithField("error", err).Fatalf("fail to create new NSQ consumer")
	}

	consumer.SetLogger(log.New(os.Stderr, fmt.Sprintf("[nsq-consumer(%s)]", c.Topic), log.Flags()), nsq.LogLevelWarning)

	consumer.AddConcurrentHandlers(nsq.HandlerFunc(func(message *nsq.Message) (err error) {
		defer func() {
			if errRecovered := recover(); errRecovered != nil {
				err = errgo.Newf("recover panic from nsq consumer: %+v", errRecovered)
				c.logger.WithFields(logrus.Fields{"error": errRecovered.(error), "stacktrace": errgo.Details(errRecovered.(error))}).Error("recover panic")
				rollbar.Error(rollbar.ERR, errRecovered.(error), &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			}
		}()

		if len(message.Body) == 0 {
			err := errgo.New("body is blank, re-enqueued message")
			c.logger.Error(err)
			rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			return err
		}
		var msg NsqMessageDeserialize
		err = json.Unmarshal(message.Body, &msg)
		if err != nil {
			rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			c.logger.WithField("error", err).Error("failed to unmarshal message")
			return errgo.Mask(err, errgo.Any)
		}
		msg.NsqMsg = message

		msg.logger = c.logger.WithFields(logrus.Fields{"message-id": fmt.Sprintf("%s", message.ID), "message-type": msg.Type})

		if msg.At != 0 {
			now := time.Now().Unix()
			delay := msg.At - now
			if delay > 0 {
				return c.postponeMessage(msg, delay)
			}
		}

		before := time.Now()
		msg.Logger().Printf("BEGIN Message")
		err = c.MessageHandler(&msg)
		if err != nil {
			rollbar.ErrorWithStack(rollbar.ERR, err, errgorollbar.BuildStack(err), &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
			msg.Logger().WithField("stacktrace", errgo.Details(err)).Error(err)
			return errgo.Mask(err, errgo.Any)
		}
		c.logger.WithField("duration", time.Since(before)).Printf("END Message")
		return nil
	}), c.MaxInFlight)

	if err = consumer.ConnectToNSQLookupds(c.NsqLookupdURLs); err != nil {
		rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "authenticator"})
		c.logger.WithField("error", err).Fatalf("Fail to connect to NSQ lookupd")
	}

	return func() {
		consumer.Stop()
		// block until stop process is complete
		<-consumer.StopChan
	}
}

func (c *nsqConsumer) postponeMessage(msg NsqMessageDeserialize, delay int64) error {
	if delay > maxPostponeDelay {
		delay = maxPostponeDelay
	}

	publishedMsg := nsqproducer.NsqMessageSerialize{
		At:      msg.At,
		Type:    msg.Type,
		Payload: msg.Payload,
	}

	msg.Logger().Printf("POSTPONE Messaage")

	if c.PostponeProducer == nil {
		return errors.New("no postpone producer configured in this consumer")
	}

	return c.PostponeProducer.DeferredPublish(c.Topic, delay, publishedMsg)
}
