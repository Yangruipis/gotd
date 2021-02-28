package biz

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/core"
)

type Task struct {
	ctx context.Context
	mgr core.TaskManager
}

func NewTask(ctx context.Context, mgr core.TaskManager) *Task {
	return &Task{
		ctx: ctx,
		mgr: mgr,
	}
}

func (t *Task) Create(task *core.Task) (*core.Task, error) {
	return t.mgr.CreateOrUpdate(t.ctx, task)
}
