package channelsandfunctions

import (
	"context"
	"time"
)

type Event struct{}

type Service interface {
	Events() <-chan *Event
	Subscribe(func(context.Context, *time.Timer, string) error) (chan<- Event, func() error)
}
