package core

import (
	"context"
)

type TaskFilterParam struct {
	NameKeyword string
	DescKeyword string
	State       uint8
	Priority    uint8
	MinTime     uint64
	MaxTime     uint64
}

type TaskManager interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Get(ctx context.Context, id uint32) (*Task, error)
	Delete(ctx context.Context, id uint32) error

	List(ctx context.Context, filter TaskFilterParam) ([]*Task, error)
}

type EventFilterParam struct {
}
type EventManager interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	Get(ctx context.Context, id uint32) (*Event, error)
	GetByTaskID(ctx context.Context, taskId uint32) ([]*Event, error)

	List(ctx context.Context, filter EventFilterParam) ([]*Event, error)
}
