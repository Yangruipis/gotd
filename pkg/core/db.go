package core

import (
	"context"
)

type TaskFilterParam struct {
	NameKeyword string
	DescKeyword string
	State       uint8
	Priority    uint8
	MinTime     int64
	MaxTime     int64
}

type TaskDao interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	Update(ctx context.Context, task *Task) (*Task, error)
	Get(ctx context.Context, id uint32) (*Task, error)
	Delete(ctx context.Context, id uint32) error

	List(ctx context.Context, filter TaskFilterParam) ([]*Task, error)
}

type EventFilterParam struct {
}
type EventDao interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	Get(ctx context.Context, id uint32) (*Event, error)
	GetByTaskID(ctx context.Context, taskId uint32) ([]*Event, error)

	List(ctx context.Context, filter EventFilterParam) ([]*Event, error)
}

type TagFilterParam struct {
}
type TagDao interface {
	Create(ctx context.Context, tag *Tag) (*Tag, error)
	Delete(ctx context.Context, id uint32) error
	Get(ctx context.Context, id uint32) (*Tag, error)
	List(ctx context.Context, filter TagFilterParam) ([]*Tag, error)

	CreateTaskTag(ctx context.Context, tag *Tag, taskId uint32) (*Tag, error)
	DeleteTaskTag(ctx context.Context, tagId uint32, taskId uint32) error
}
