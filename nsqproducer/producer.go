package nsqproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-utils/logger"
)

type Producer interface {
	Publish(ctx context.Context, topic string, message NsqMessageSerialize) error
	DeferredPublish(ctx context.Context, topic string, delay int64, message NsqMessageSerialize) error
	Stop()
}

type NsqProducer struct {
	producer   *nsq.Producer
	config     *nsq.Config
	skipLogSet map[string]bool
	telemetry  *telemetry
}

type ProducerOpts struct {
	Host       string
	Port       string
	NsqConfig  *nsq.Config
	SkipLogSet map[string]bool
	// WithoutTelemetry indicates whether OpenTelemetry instrumentation should be disabled
	WithoutTelemetry bool
}

type WithLoggableFields interface {
	LoggableFields() logrus.Fields
}

type NsqMessageSerialize struct {
	At      int64       `json:"at"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`

	// Automatically set by context if existing, generated otherwise
	RequestID string `json:"request_id"`
}

var _ Producer = &NsqProducer{} // Ensure that NsqProducer implements the Producer interface

func New(opts ProducerOpts) (*NsqProducer, error) {
	client, err := nsq.NewProducer(opts.Host+":"+opts.Port, opts.NsqConfig)
	if err != nil {
		return nil, fmt.Errorf("init-nsq: cannot initialize nsq producer: %v:%v", opts.Host, opts.Port)
	}

	if opts.SkipLogSet == nil {
		opts.SkipLogSet = map[string]bool{}
	}

	var telemetry *telemetry
	if !opts.WithoutTelemetry {
		telemetry, err = newTelemetry()
		if err != nil {
			return nil, fmt.Errorf("init-nsq: cannot initialize telemetry: %v", err)
		}
	}

	return &NsqProducer{
		producer:   client,
		config:     opts.NsqConfig,
		skipLogSet: opts.SkipLogSet,
		telemetry:  telemetry,
	}, nil
}

func (p *NsqProducer) Stop() {
	p.producer.Stop()
}

func (p *NsqProducer) Ping() error {
	return p.producer.Ping()
}

func (p *NsqProducer) Publish(ctx context.Context, topic string, message NsqMessageSerialize) error {
	startedAt := time.Now()
	messageType := message.Type
	if messageType == "" {
		messageType = unknownMessageType
	}
	publishType := publishTypeImmediate
	var telemetryErr error
	defer func() {
		if p.telemetry != nil {
			p.telemetry.record(ctx, startedAt, topic, messageType, publishType, telemetryErr)
		}
	}()

	var err error
	message.RequestID, err = p.requestID(ctx)
	if err != nil {
		err = errgo.Notef(err, "fail to get requestID")
		telemetryErr = err
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		telemetryErr = err
		return err
	}

	err = p.producer.Publish(topic, body)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		telemetryErr = err
		return err
	}

	p.log(ctx, message, logrus.Fields{})

	return nil
}

func (p *NsqProducer) DeferredPublish(ctx context.Context, topic string, delay int64, message NsqMessageSerialize) error {
	startedAt := time.Now()
	messageType := message.Type
	if messageType == "" {
		messageType = unknownMessageType
	}
	publishType := publishTypeDeferred
	var telemetryErr error
	defer func() {
		if p.telemetry != nil {
			p.telemetry.record(ctx, startedAt, topic, messageType, publishType, telemetryErr)
		}
	}()

	var err error
	message.RequestID, err = p.requestID(ctx)
	if err != nil {
		err = errgo.Notef(err, "fail to get requestID")
		telemetryErr = err
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		telemetryErr = err
		return err
	}

	err = p.producer.DeferredPublish(topic, time.Duration(delay)*time.Second, body)
	if err != nil {
		err = errgo.Mask(err, errgo.Any)
		telemetryErr = err
		return err
	}

	p.log(ctx, message, logrus.Fields{"message_delay": delay})

	return nil
}

func (p *NsqProducer) requestID(ctx context.Context) (string, error) {
	reqid, ok := ctx.Value("request_id").(string)
	if !ok {
		uuid, err := uuid.NewV4()
		if err != nil {
			return "", errgo.Notef(err, "fail to generate UUID v4")
		}
		return uuid.String(), nil
	}
	return reqid, nil
}

func (p *NsqProducer) logger(ctx context.Context) logrus.FieldLogger {
	return logger.Get(ctx)
}

func (p *NsqProducer) log(ctx context.Context, message NsqMessageSerialize, fields logrus.Fields) {
	if p.skipLogSet[message.Type] {
		return
	}

	logger := p.logger(ctx).WithFields(fields)

	if logger.Level == logrus.DebugLevel {
		logger.WithFields(logrus.Fields{"message_type": message.Type, "message_payload": message.Payload}).Debug("publish message")
	} else {
		// We don't want the complete payload to be dump in the logs With this
		// interface we can, for each type of payload, add fields in the logs.
		if payload, ok := message.Payload.(WithLoggableFields); ok {
			logger = logger.WithFields(payload.LoggableFields())
		}
		logger.WithFields(logrus.Fields{"message_type": message.Type}).Info("publish message")
	}
}
