package nsqlbproducer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Scalingo/go-utils/nsqproducer"
	"github.com/juju/errgo/errors"
	nsq "github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	errgo "gopkg.in/errgo.v1"
)

// NsqLBProducer a producer that distribute nsq messages across a set of node
// if a node send an error when receiving the message it will try with another node of the set
type NsqLBProducer struct {
	producers []producer
	randSrc   randSource
	logger    logrus.FieldLogger
}

type producer struct {
	producer nsqproducer.Producer
	host     Host
}

type Host struct {
	Host string
	Port string
}

func (h Host) String() string {
	return fmt.Sprintf("%s:%s", h.Host, h.Port)
}

type LBProducerOpts struct {
	Hosts      []Host
	NsqConfig  *nsq.Config
	Logger     logrus.FieldLogger
	SkipLogSet map[string]bool
}

type randSource interface {
	Int() int
}

var _ nsqproducer.Producer = &NsqLBProducer{} // Ensure that NsqLBProducer implements the Producer interface

func New(opts LBProducerOpts) (*NsqLBProducer, error) {
	if len(opts.Hosts) == 0 {
		return nil, fmt.Errorf("A producer must have at least one host")
	}
	lbproducer := &NsqLBProducer{
		producers: make([]producer, len(opts.Hosts)),
	}

	for i, h := range opts.Hosts {
		p, err := nsqproducer.New(nsqproducer.ProducerOpts{
			Host:       h.Host,
			Port:       h.Port,
			NsqConfig:  opts.NsqConfig,
			SkipLogSet: opts.SkipLogSet,
		})

		if err != nil {
			return nil, errors.Notef(err, "fail to create producer for host: %s:%s", h.Host, h.Port)
		}

		lbproducer.producers[i] = producer{
			producer: p,
			host:     h,
		}
	}

	lbproducer.randSrc = rand.New(rand.NewSource(time.Now().Unix()))
	lbproducer.logger = opts.Logger

	return lbproducer, nil
}

func (p *NsqLBProducer) Publish(ctx context.Context, topic string, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randSrc.Int() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		producer := p.producers[(i+firstProducer)%len(p.producers)]
		err = producer.producer.Publish(ctx, topic, message)
		if err != nil {
			if p.logger != nil {
				p.logger.WithError(err).WithField("host", producer.host.String()).Error("fail to send nsq message to one nsq node")
			}
		} else {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) DeferredPublish(ctx context.Context, topic string, delay int64, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randSrc.Int() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		producer := p.producers[(i+firstProducer)%len(p.producers)]
		err = producer.producer.DeferredPublish(ctx, topic, delay, message)
		if err != nil {
			if p.logger != nil {
				p.logger.WithError(err).WithField("host", producer.host.String()).Error("fail to send nsq message to one nsq node")
			}
		} else {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) Stop() {
	for _, p := range p.producers {
		p.producer.Stop()
	}
}
