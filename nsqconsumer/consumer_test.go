package nsqconsumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-utils/nsqproducer"
)

type fakeProducer struct {
	called  bool
	topic   string
	delay   int64
	message nsqproducer.NsqMessageSerialize
	err     error
}

func (f *fakeProducer) Publish(context.Context, string, nsqproducer.NsqMessageSerialize) error {
	return nil
}

func (f *fakeProducer) DeferredPublish(_ context.Context, topic string, delay int64, message nsqproducer.NsqMessageSerialize) error {
	f.called = true
	f.topic = topic
	f.delay = delay
	f.message = message
	return f.err
}

func (f *fakeProducer) Stop() {}

type testMessageDelegate struct {
	finished bool
	touched  bool
	requeued bool
	backoff  bool
	delay    time.Duration
}

func (d *testMessageDelegate) OnFinish(*nsq.Message) {
	d.finished = true
}

func (d *testMessageDelegate) OnRequeue(_ *nsq.Message, delay time.Duration, backoff bool) {
	d.requeued = true
	d.delay = delay
	d.backoff = backoff
}

func (d *testMessageDelegate) OnTouch(*nsq.Message) {
	d.touched = true
}

func newTestConsumer(t *testing.T, handler func(context.Context, *NsqMessageDeserialize) error) *nsqConsumer {
	t.Helper()
	return &nsqConsumer{
		Topic:          "events",
		Channel:        "default",
		MessageHandler: handler,
		SkipLogSet:     map[string]bool{},
		logger:         logrus.New(),
	}
}

func newTestMessage(t *testing.T, body []byte) *nsq.Message {
	t.Helper()
	msg := nsq.NewMessage(nsq.MessageID{}, body)
	delegate := &testMessageDelegate{}
	msg.Delegate = delegate
	return msg
}

func testGetDelegateFromMessage(t *testing.T, msg *nsq.Message) *testMessageDelegate {
	t.Helper()
	delegate, ok := msg.Delegate.(*testMessageDelegate)
	require.True(t, ok)
	return delegate
}

func TestNsqConsumer_nsqHandler(t *testing.T) {
	t.Parallel()

	t.Run("rejects invalid bodies", func(t *testing.T) {
		c := newTestConsumer(t, func(context.Context, *NsqMessageDeserialize) error {
			t.Fatal("handler should not be called")
			return nil
		})

		blankMsg := newTestMessage(t, nil)
		err := c.nsqHandler(blankMsg)
		require.EqualError(t, err, "body is blank, re-enqueued message")

		invalidMsg := newTestMessage(t, []byte("{"))
		err = c.nsqHandler(invalidMsg)
		require.EqualError(t, err, "unexpected end of JSON input")
	})

	t.Run("success", func(t *testing.T) {
		called := false
		c := newTestConsumer(t, func(_ context.Context, msg *NsqMessageDeserialize) error {
			called = true
			require.Equal(t, "user.created", msg.Type)
			require.Equal(t, "req-123", msg.RequestID)
			return nil
		})

		body := []byte(`{"request_id":"req-123","type":"user.created","payload":{"id":42}}`)
		msg := newTestMessage(t, body)
		delegate := testGetDelegateFromMessage(t, msg)
		err := c.nsqHandler(msg)

		require.NoError(t, err)
		require.True(t, called)
		require.False(t, delegate.finished)
		require.False(t, delegate.requeued)
	})

	t.Run("no retry error finishes", func(t *testing.T) {
		c := newTestConsumer(t, func(context.Context, *NsqMessageDeserialize) error {
			return NewError(errors.New("boom"), ErrorOpts{NoRetry: true})
		})

		msg := newTestMessage(t, []byte(`{"type":"job","payload":{}}`))
		delegate := testGetDelegateFromMessage(t, msg)
		err := c.nsqHandler(msg)

		require.NoError(t, err)
		require.True(t, delegate.finished)
		require.False(t, delegate.requeued)
	})

	t.Run("retry errors requeue", func(t *testing.T) {
		tests := []struct {
			name            string
			disableBackoff  bool
			expectedBackoff bool
		}{
			{name: "with backoff", disableBackoff: false, expectedBackoff: true},
			{name: "without backoff", disableBackoff: true, expectedBackoff: false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				c := newTestConsumer(t, func(context.Context, *NsqMessageDeserialize) error {
					return errors.New("boom")
				})
				c.disableBackoff = test.disableBackoff

				msg := newTestMessage(t, []byte(`{"type":"job","payload":{}}`))
				delegate := testGetDelegateFromMessage(t, msg)
				err := c.nsqHandler(msg)

				require.NoError(t, err)
				require.False(t, delegate.finished)
				require.True(t, delegate.requeued)
				require.Equal(t, time.Duration(-1), delegate.delay)
				require.Equal(t, test.expectedBackoff, delegate.backoff)
			})
		}
	})

	t.Run("postpones future messages", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fake := &fakeProducer{}
			called := false
			c := newTestConsumer(t, func(context.Context, *NsqMessageDeserialize) error {
				called = true
				return nil
			})
			c.PostponeProducer = fake

			at := time.Now().Add(2 * time.Second).Unix()
			body, err := json.Marshal(map[string]any{
				"request_id": "req-1",
				"type":       "deferred",
				"at":         at,
				"payload":    map[string]string{"k": "v"},
			})
			require.NoError(t, err)

			firstMsg := newTestMessage(t, body)
			err = c.nsqHandler(firstMsg)
			require.NoError(t, err)
			require.False(t, called)
			require.True(t, fake.called)
			require.Equal(t, c.Topic, fake.topic)
			require.Positive(t, fake.delay)
			require.Equal(t, "deferred", fake.message.Type)
			require.Equal(t, at, fake.message.At)

			time.Sleep(3 * time.Second)
			synctest.Wait()

			republishedBody, err := json.Marshal(fake.message)
			require.NoError(t, err)
			secondMsg := newTestMessage(t, republishedBody)
			err = c.nsqHandler(secondMsg)
			require.NoError(t, err)
			require.True(t, called)
		})
	})
}

func TestPostponeMessage(t *testing.T) {
	t.Parallel()

	msgLogger := logrus.New()
	msg := NsqMessageDeserialize{
		At:      time.Now().Add(2 * time.Hour).Unix(),
		Type:    "deferred",
		Payload: json.RawMessage(`{"id":1}`),
	}

	t.Run("returns error when postpone producer is missing", func(t *testing.T) {
		c := &nsqConsumer{
			Topic: "events",
		}

		err := c.postponeMessage(t.Context(), msgLogger, msg, 10)
		require.EqualError(t, err, "no postpone producer configured in this consumer")
	})

	t.Run("caps delay to max postpone delay and publishes", func(t *testing.T) {
		c := &nsqConsumer{
			Topic: "events",
		}
		fake := &fakeProducer{}
		c.PostponeProducer = fake

		err := c.postponeMessage(t.Context(), msgLogger, msg, maxPostponeDelay+100)
		require.NoError(t, err)
		require.True(t, fake.called)
		require.Equal(t, maxPostponeDelay, fake.delay)
		require.Equal(t, "events", fake.topic)
		require.Equal(t, msg.At, fake.message.At)
		require.Equal(t, msg.Type, fake.message.Type)
	})
}

func TestStopWaitsForMessagesOrTimeout(t *testing.T) {
	t.Parallel()

	newNSQConsumer := func(t *testing.T) *nsq.Consumer {
		t.Helper()
		cfg := nsq.NewConfig()
		cfg.MaxInFlight = 1
		consumer, err := nsq.NewConsumer("topic", "channel", cfg)
		require.NoError(t, err)
		consumer.AddHandler(nsq.HandlerFunc(func(*nsq.Message) error { return nil }))
		return consumer
	}

	t.Run("no in-flight", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := &nsqConsumer{
				logger:               logrus.New(),
				gracefulStopDuration: 2 * time.Second,
			}
			c.wgOngoingMessages.Add(1)
			time.AfterFunc(time.Second, c.wgOngoingMessages.Done)

			consumer := newNSQConsumer(t)
			startedAt := time.Now()
			c.Stop(t.Context(), consumer)
			require.Equal(t, time.Second, time.Since(startedAt))
			require.Less(t, time.Since(startedAt), c.gracefulStopDuration)
			synctest.Wait()
		})
	})

	t.Run("graceful timeout", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			c := &nsqConsumer{
				logger:               logrus.New(),
				gracefulStopDuration: time.Second,
			}
			c.wgOngoingMessages.Add(1)
			time.AfterFunc(2*time.Second, c.wgOngoingMessages.Done)

			consumer := newNSQConsumer(t)
			startedAt := time.Now()
			c.Stop(t.Context(), consumer)
			require.Equal(t, time.Second, time.Since(startedAt))
			// Wait for the end of the AfterFunc call
			time.Sleep(2 * time.Second)
			synctest.Wait()
		})
	})
}
