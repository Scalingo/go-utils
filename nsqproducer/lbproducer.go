package nsqproducer

import (
	"context"
	"fmt"
	"math/rand"

	nsq "github.com/nsqio/go-nsq"
	errgo "gopkg.in/errgo.v1"
)

type NsqLBProducer struct {
	producers []Producer
}

type Host struct {
	Host string
	Port string
}

type LBProducerOpts struct {
	Hosts      []Host
	NsqConfig  *nsq.Config
	SkipLogSet map[string]bool
}

func NewLB(opts LBProducerOpts) (*NsqLBProducer, error) {
	if len(opts.Hosts) == 0 {
		return nil, fmt.Errorf("A producer must have at least one host")
	}
	producer := &NsqLBProducer{
		producers: make([]Producer, len(opts.Hosts)),
	}

	for i, h := range opts.Hosts {
		p, err := New(ProducerOpts{
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

	return producer, nil
}

func (p *NsqLBProducer) Publish(ctx context.Context, topic string, message NsqMessageSerialize) error {
	firstProducer := rand.Int() % len(p.producers)

	var err error
	for i := 0; i < len(p.producers); i++ {
		err = p.producers[(i+firstProducer)%len(p.producers)].Publish(ctx, topic, message)
		if err == nil {
			return nil
		}
	}

	return errgo.Notef(err, "fail to send message on %v hosts", len(p.producers))
}

func (p *NsqLBProducer) DeferredPublish(ctx context.Context, topic string, delay int64, message NsqMessageSerialize) error {
	firstProducer := rand.Int() % len(p.producers)

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
