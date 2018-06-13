package nsqlbproducer

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Scalingo/go-utils/nsqproducer"
	nsq "github.com/nsqio/go-nsq"
	errgo "gopkg.in/errgo.v1"
)

type LBStrategy int

const (
	RandomStrategy LBStrategy = iota
	FallbackStrategy
)

type NsqLBProducer struct {
	producers []nsqproducer.Producer
	randInt   func() int
	strategy  LBStrategy
}

type Host struct {
	Host string
	Port string
}

type LBProducerOpts struct {
	Hosts      []Host
	NsqConfig  *nsq.Config
	SkipLogSet map[string]bool
	Strategy   LBStrategy
}

func New(opts LBProducerOpts) (*NsqLBProducer, error) {
	if len(opts.Hosts) == 0 {
		return nil, fmt.Errorf("A producer must have at least one host")
	}
	producer := &NsqLBProducer{
		producers: make([]nsqproducer.Producer, len(opts.Hosts)),
		strategy:  opts.Strategy,
	}

	for i, h := range opts.Hosts {
		p, err := nsqproducer.New(nsqproducer.ProducerOpts{
			Host:       h.Host,
			Port:       h.Port,
			NsqConfig:  opts.NsqConfig,
			SkipLogSet: opts.SkipLogSet,
		})

		if err != nil {
			return nil, errgo.Mask(err)
		}

		producer.producers[i] = p
	}

	switch producer.strategy {
	case FallbackStrategy:
		producer.randInt = alwaysZero
	case RandomStrategy:
		fallthrough
	default:
		producer.randInt = rand.New(rand.NewSource(time.Now().Unix())).Int
	}

	return producer, nil
}

func alwaysZero() int {
	return 0
}

func (p *NsqLBProducer) Publish(ctx context.Context, topic string, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randInt() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		err = p.producers[(i+firstProducer)%len(p.producers)].Publish(ctx, topic, message)
		if err == nil {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) DeferredPublish(ctx context.Context, topic string, delay int64, message nsqproducer.NsqMessageSerialize) error {
	firstProducer := p.randInt() % len(p.producers)

	var err error

	for i := 0; i < len(p.producers); i++ {
		err = p.producers[(i+firstProducer)%len(p.producers)].DeferredPublish(ctx, topic, delay, message)
		if err == nil {
			return nil
		}
	}
	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) Stop() {
	for _, p := range p.producers {
		p.Stop()
	}
}
