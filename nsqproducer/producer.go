package nsqproducer

import (
	"encoding/json"

	"log"

	"github.com/nsqio/go-nsq"
	"gopkg.in/errgo.v1"
)

type ProducerOpts struct {
	Host      string
	Port      string
	NsqConfig *nsq.Config
}

type NsqMessageSerialize struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var client *nsq.Producer

func Init(opts ProducerOpts) func() {
	var err error
	client, err = nsq.NewProducer(opts.Host+":"+opts.Port, opts.NsqConfig)
	if err != nil {
		log.Fatalf("init-nsq: cannot initialize nsq producer: %v:%v", opts.Host, opts.Port)
	}
	return func() {
		client.Stop()
	}
}

func Publish(topic string, message NsqMessageSerialize) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	err = client.Publish(topic, body)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	return nil
}
