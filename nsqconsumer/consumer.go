package nsqconsumer

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/stvp/rollbar"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/nsqproducer"
)

const (
	// defaultChannel is the name of the channel we're using when we want the
	// message to be receive only by 1 consumer, but no matter which one
	defaultChannel = "default"
)

var (
	maxPostponeDelay int64 = 3600
)

// LogLevel is a wrapper around nsq.LogLevel to ensure that the default log level is set to Warning and not Debug
type LogLevel int

const (
	// DefaultLogLevel is the default log level for NSQ when no log level is provided
	DefaultLogLevel LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

func ParseLogLevel(logLevel string) LogLevel {
	switch logLevel {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warning":
		return LogLevelWarning
	case "error":
		return LogLevelError
	default:
		return DefaultLogLevel
	}
}

func (l LogLevel) toNSQLogLevel() nsq.LogLevel {
	switch l {
	case LogLevelDebug:
		return nsq.LogLevelDebug
	case LogLevelInfo:
		return nsq.LogLevelInfo
	case LogLevelWarning:
		return nsq.LogLevelWarning
	case LogLevelError:
		return nsq.LogLevelError
	case DefaultLogLevel:
		return nsq.LogLevelWarning
	default:
		return nsq.LogLevelWarning
	}
}

type Error struct {
	error   error
	noRetry bool
}

func (nsqerr Error) Error() string {
	return nsqerr.error.Error()
}

// Unwrap returns the cause of the error to be compatible with errors.As/Is()
func (nsqerr Error) Unwrap() error {
	return nsqerr.error
}

// NoRetry returns true if the message should not be retried to be handled
func (nsqerr Error) NoRetry() bool {
	return nsqerr.noRetry
}

type ErrorOpts struct {
	NoRetry bool
}

func NewError(err error, opts ErrorOpts) error {
	return Error{error: err, noRetry: opts.NoRetry}
}

type NsqMessageDeserialize struct {
	RequestID string          `json:"request_id"`
	Type      string          `json:"type"`
	At        int64           `json:"at"`
	Payload   json.RawMessage `json:"payload"`
	NsqMsg    *nsq.Message
}

// FromMessageSerialize let you transform a Serialized message to a DeserializeMessage for a consumer
// Its use is mostly for testing as writing manual `json.RawMessage` is boring
func FromMessageSerialize(msg *nsqproducer.NsqMessageSerialize) *NsqMessageDeserialize {
	res := &NsqMessageDeserialize{
		At:   msg.At,
		Type: msg.Type,
	}
	buffer, _ := json.Marshal(msg.Payload)
	res.Payload = json.RawMessage(buffer)
	return res
}

// TouchUntilClosed returns a channel which has to be closed by the called
// Until the channel is closed, the NSQ message will be touched every 40 secs
// to ensure NSQ does not consider the message as failed because of time out.
func (msg *NsqMessageDeserialize) TouchUntilClosed() chan<- struct{} {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(40 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				msg.NsqMsg.Touch()
			}
		}
	}()
	return done
}

type nsqConsumer struct {
	NsqConfig        *nsq.Config
	NsqLookupdURLs   []string
	Topic            string
	Channel          string
	MessageHandler   func(context.Context, *NsqMessageDeserialize) error
	MaxInFlight      int
	SkipLogSet       map[string]bool
	PostponeProducer nsqproducer.Producer
	disableBackoff   bool
	count            uint64
	logger           logrus.FieldLogger
	logLevel         LogLevel
	telemetry        *telemetry
	withTelemetry    bool
}

type ConsumerOpts struct {
	NsqConfig      *nsq.Config
	NsqLookupdURLs []string
	Topic          string
	Channel        string
	MaxInFlight    int
	SkipLogSet     map[string]bool
	LogLevel       LogLevel
	// PostponeProducer is an NSQ producer user to send postponed messages
	PostponeProducer nsqproducer.Producer
	// How long can the consumer keep the message before the message is considered as 'Timed Out'
	MsgTimeout     time.Duration
	MessageHandler func(context.Context, *NsqMessageDeserialize) error
	DisableBackoff bool
	// WithoutTelemetry indicates whether OpenTelemetry instrumentation should be disabled
	WithoutTelemetry bool
}

type Consumer interface {
	Start(ctx context.Context) func()
}

func New(opts ConsumerOpts) (Consumer, error) {
	if opts.MsgTimeout != 0 {
		opts.NsqConfig.MsgTimeout = opts.MsgTimeout
	}

	if opts.SkipLogSet == nil {
		opts.SkipLogSet = map[string]bool{}
	}

	consumer := &nsqConsumer{
		NsqConfig:      opts.NsqConfig,
		NsqLookupdURLs: opts.NsqLookupdURLs,
		Topic:          opts.Topic,
		Channel:        opts.Channel,
		MessageHandler: opts.MessageHandler,
		MaxInFlight:    opts.MaxInFlight,
		SkipLogSet:     opts.SkipLogSet,
		logLevel:       opts.LogLevel,
		disableBackoff: opts.DisableBackoff,
		withTelemetry:  !opts.WithoutTelemetry,
	}
	if consumer.MaxInFlight == 0 {
		consumer.MaxInFlight = opts.NsqConfig.MaxInFlight
	} else {
		opts.NsqConfig.MaxInFlight = consumer.MaxInFlight
	}
	if opts.Topic == "" {
		return nil, stderrors.New("topic can't be blank")
	}
	if opts.MessageHandler == nil {
		return nil, stderrors.New("message handler can't be blank")
	}
	if opts.Channel == "" {
		consumer.Channel = defaultChannel
	}
	return consumer, nil
}

func (c *nsqConsumer) Start(ctx context.Context) func() {
	c.logger = logger.Get(ctx).WithFields(logrus.Fields{
		"topic":   c.Topic,
		"channel": c.Channel,
	})
	c.logger.Info("Starting consumer")
	if c.withTelemetry {
		telemetry, err := newTelemetry(ctx)
		if err != nil {
			c.logger.WithError(err).Error("Fail to init telemetry")
		} else {
			c.telemetry = telemetry
		}
	}

	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, c.NsqConfig)
	if err != nil {
		rollbar.Error(rollbar.ERR, err, &rollbar.Field{Name: "worker", Data: "nsq-consumer"})
		c.logger.WithError(err).Fatalf("Fail to create new NSQ consumer")
	}

	consumer.SetLogger(log.New(os.Stderr, fmt.Sprintf("[nsq-consumer(%s)]", c.Topic), log.Flags()), c.logLevel.toNSQLogLevel())

	consumer.AddConcurrentHandlers(nsq.HandlerFunc(c.nsqHandler), c.MaxInFlight)

	err = consumer.ConnectToNSQLookupds(c.NsqLookupdURLs)
	if err != nil {
		c.logger.WithError(err).Error("Fail to connect to NSQ lookupd")
		os.Exit(1)
	}

	return func() {
		consumer.Stop()
		// block until stop process is complete
		<-consumer.StopChan
	}
}

func (c *nsqConsumer) nsqHandler(message *nsq.Message) (err error) {
	ctx := logger.ToCtx(context.Background(), c.logger)
	// We create a new `log` variable because we are so used to call `log.Something` in our codebase
	log := c.logger
	startedAt := time.Now()

	messageType := unknownMessageType
	var telemetryErr error
	defer func() {
		if r := recover(); r != nil {
			var errRecovered error
			switch value := errRecovered.(type) {
			case error:
				errRecovered = value
			default:
				errRecovered = errors.Newf(ctx, "%v", value)
			}
			err = errors.Newf(ctx, "recover panic from nsq consumer: %+v", errRecovered)
			telemetryErr = err
			log.WithError(errRecovered).WithFields(logrus.Fields{
				"stacktrace": string(debug.Stack()),
			}).Error("Recover panic")
		}
		if c.telemetry != nil {
			c.telemetry.record(ctx, startedAt, c.Topic, c.Channel, messageType, telemetryErr)
		}
	}()

	if len(message.Body) == 0 {
		err := errors.New(ctx, "body is blank, re-enqueued message")
		telemetryErr = err
		log.WithError(err).Error("Blank message")
		return err
	}
	var msg NsqMessageDeserialize
	err = json.Unmarshal(message.Body, &msg)
	if err != nil {
		telemetryErr = err
		log.WithError(err).Error("Fail to unmarshal message")
		return err
	}
	msg.NsqMsg = message
	if msg.Type != "" {
		messageType = msg.Type
	}

	// We want to create a new dedicated logger for the NSQ message handling.
	// That way we distinguish between the logger during normal operation, and the error logger (named `errLogger`) found when unwrapping the error raised during message handling.
	msgLogger := logger.Default()
	ctx = logger.ToCtx(context.Background(), msgLogger)
	ctx, msgLogger = logger.WithFieldsToCtx(ctx, logrus.Fields{
		"message_id":   fmt.Sprintf("%s", message.ID),
		"message_type": msg.Type,
		"request_id":   msg.RequestID,
	})

	if msg.At != 0 {
		now := time.Now().Unix()
		delay := msg.At - now
		if delay > 0 {
			err = c.postponeMessage(ctx, msgLogger, msg, delay)
			if err != nil {
				telemetryErr = err
			}
			return err
		}
	}

	before := time.Now()
	if _, ok := c.SkipLogSet[msg.Type]; !ok {
		msgLogger.Info("BEGIN Message")
	}

	err = c.MessageHandler(ctx, &msg)
	if err != nil {
		telemetryErr = err
		var errLogger logrus.FieldLogger
		noRetry := false

		unwrapErr := err
		for unwrapErr != nil {
			switch handlerErr := unwrapErr.(type) {
			case interface{ Ctx() context.Context }:
				errLogger = logger.Get(handlerErr.Ctx())
			case Error:
				noRetry = handlerErr.noRetry
			}
			unwrapErr = errors.UnwrapError(unwrapErr)
		}
		if errLogger == nil {
			errLogger = msgLogger
		}

		// Call this method so the caller won't do anything with the message, we
		// have the control of the flow
		message.DisableAutoResponse()
		if noRetry {
			errLogger.WithError(err).Error("Message handling error - noretry")
			message.Finish()
		} else {
			errLogger.WithError(err).Error("Message handling error")
			if c.disableBackoff {
				// Delay of -1 means the default algorithm will be used to compute the
				// requeue delay (duration to wait before retrying the message)
				message.RequeueWithoutBackoff(-1)
			} else {
				message.Requeue(-1)
			}
		}
		return nil
	}

	if _, ok := c.SkipLogSet[msg.Type]; !ok {
		msgLogger.WithField("duration", time.Since(before)).Info("END Message")
	}
	return nil
}

func (c *nsqConsumer) postponeMessage(ctx context.Context, msgLogger logrus.FieldLogger, msg NsqMessageDeserialize, delay int64) error {
	if delay > maxPostponeDelay {
		delay = maxPostponeDelay
	}

	publishedMsg := nsqproducer.NsqMessageSerialize{
		At:      msg.At,
		Type:    msg.Type,
		Payload: msg.Payload,
	}

	msgLogger.Info("POSTPONE Message")

	if c.PostponeProducer == nil {
		return errors.New(ctx, "no postpone producer configured in this consumer")
	}

	return c.PostponeProducer.DeferredPublish(ctx, c.Topic, delay, publishedMsg)
}
