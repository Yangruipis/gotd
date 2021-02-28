package core

import (
	"context"
)

type TaskFilterParam struct {
}

type TaskManager interface {
	CreateOrUpdate(ctx context.Context, task *Task) (*Task, error)
	Get(ctx context.Context, id string) (*Task, error)
	Delete(ctx context.Context, id string) error

	List(ctx context.Context, filter TaskFilterParam) ([]*Task, error)
}

type EventFilterParam struct {
}
type EventManager interface {
	CreateOrUpdate(ctx context.Context, task *Event) (*Event, error)
	Get(ctx context.Context, id string) (*Event, error)
	Delete(ctx context.Context, id string) error

	List(ctx context.Context, filter EventFilterParam) ([]*Event, error)
}
