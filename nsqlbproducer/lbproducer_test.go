package nsqlbproducer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Scalingo/go-utils/nsqlbproducer/nsqlbproducermock"
	"github.com/Scalingo/go-utils/nsqproducer"
)

type mockedRandSource struct {
	current int
	values  []int
}

func (m *mockedRandSource) Int() int {
	val := m.values[m.current]
	m.current++
	return val
}

type example struct {
	LBProducer   func([]producer) *NsqLBProducer
	ExpectP1Call bool
	ExpectP2Call bool
	P1Error      error
	P2Error      error
	P1Delay      time.Duration
	P2Delay      time.Duration
	ExpectError  bool
	RandInt      func() int
}

func randLBProducer(order []int) func(producers []producer) *NsqLBProducer {
	return func(producers []producer) *NsqLBProducer {
		return &NsqLBProducer{
			producers:      producers,
			randInt:        (&mockedRandSource{current: 0, values: order}).Int,
			publishTimeout: 1 * time.Second,
		}
	}
}

func TestLBPublish(t *testing.T) {
	examples := map[string]example{
		"when all host works": {
			LBProducer:   randLBProducer([]int{1}),
			ExpectP1Call: false,
			ExpectP2Call: true,
			P1Error:      nil,
			P2Error:      nil,
			ExpectError:  false,
		},
		"when a single host is down": {
			LBProducer:   randLBProducer([]int{1, 0}),
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      nil,
			P2Error:      errors.New("NOP"),
			ExpectError:  false,
		},
		"when all hosts are down": {
			LBProducer:   randLBProducer([]int{1, 0}),
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      errors.New("NOP"),
			P2Error:      errors.New("NOP"),
			ExpectError:  true,
		},
		"when using the fallback mode, the first node ": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers:      producers,
					randInt:        alwaysZero,
					publishTimeout: 1 * time.Second,
				}
			},
			ExpectP1Call: true,
			ExpectP2Call: false,
			P1Error:      nil,
		},
		"when using the fallback mode and the firs node is failing, it should call the second one": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers:      producers,
					randInt:        alwaysZero,
					publishTimeout: 1 * time.Second,
				}
			},
			ExpectP1Call: true,
			ExpectP2Call: true,
			P1Error:      errors.New("FAIL"),
		},
		"when there is a timeout on the first load producer": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers:      producers,
					randInt:        alwaysZero,
					publishTimeout: 100 * time.Millisecond,
				}
			},
			P1Delay:      200 * time.Millisecond,
			ExpectP1Call: true,
			ExpectP2Call: true,
			ExpectError:  false,
		},
		"when there is a timeout on the both producer": {
			LBProducer: func(producers []producer) *NsqLBProducer {
				return &NsqLBProducer{
					producers:      producers,
					randInt:        alwaysZero,
					publishTimeout: 10 * time.Millisecond,
				}
			},
			P1Delay:      20 * time.Millisecond,
			P2Delay:      20 * time.Millisecond,
			ExpectP1Call: true,
			ExpectP2Call: true,
			ExpectError:  true,
		},
	}

	for name, example := range examples {
		t.Run(name, func(t *testing.T) {
			t.Run("Publish", func(t *testing.T) {
				runPublishExample(t, example, false)
			})
			t.Run("DeferredPublish", func(t *testing.T) {
				runPublishExample(t, example, true)
			})
		})
	}
}

func runPublishExample(t *testing.T, example example, deferred bool) {
	ctx := context.Background()
	message := nsqproducer.NsqMessageSerialize{}
	topic := "topic"
	delay := int64(0)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p1 := nsqlbproducermock.NewMockPublishPinger(ctrl)
	p2 := nsqlbproducermock.NewMockPublishPinger(ctrl)

	if example.ExpectP1Call {
		if deferred {
			p1.EXPECT().DeferredPublish(gomock.Any(), topic, delay, message).DoAndReturn(func(ctx context.Context, _, _, _ interface{}) error {
				return timeoutOrError(ctx, example.P1Delay, example.P1Error)
			})
		} else {
			p1.EXPECT().Publish(gomock.Any(), topic, message).DoAndReturn(func(ctx context.Context, _, _ interface{}) error {
				return timeoutOrError(ctx, example.P1Delay, example.P1Error)
			})
		}
	}

	if example.ExpectP2Call {
		if deferred {
			p2.EXPECT().DeferredPublish(gomock.Any(), topic, delay, message).DoAndReturn(func(ctx context.Context, _, _, _ interface{}) error {
				return timeoutOrError(ctx, example.P2Delay, example.P2Error)
			})
		} else {
			p2.EXPECT().Publish(gomock.Any(), topic, message).DoAndReturn(func(ctx context.Context, _, _ interface{}) error {
				return timeoutOrError(ctx, example.P2Delay, example.P2Error)
			})
		}
	}

	producer := example.LBProducer([]producer{{producer: p1, host: Host{}}, {producer: p2, host: Host{}}})

	var err error
	if deferred {
		err = producer.DeferredPublish(ctx, topic, delay, message)
	} else {
		err = producer.Publish(ctx, topic, message)
	}

	if example.ExpectError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}

func timeoutOrError(ctx context.Context, delay time.Duration, err error) error {
	timer := time.NewTimer(delay)
	select {
	case <-timer.C:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
