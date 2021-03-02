package biz

import (
	"context"
	"time"

	"github.com/Yangruipis/gotd/pkg/core"
)

type Biz struct {
	ctx      context.Context
	taskMgr  core.TaskManager
	eventMgr core.EventManager
}

func NewBiz(ctx context.Context,
	taskMgr core.TaskManager,
	eventMgr core.EventManager) *Biz {
	return &Biz{
		ctx:      ctx,
		taskMgr:  taskMgr,
		eventMgr: eventMgr,
	}
}

func (b *Biz) CreateTask(task *core.Task) (*core.Task, error) {
	if _, err := b.eventMgr.Create(b.ctx, &core.Event{
		TaskID:    task.ID,
		OccurTime: time.Now(),
		PrevState: 0,
		CurState:  task.State,
	}); err != nil {
		return nil, err
	}
	return b.taskMgr.Create(b.ctx, task)
}

func (b *Biz) GetTask(id uint32) (*core.Task, error) {
	return b.taskMgr.Get(b.ctx, id)
}

func (b *Biz) UpdateTaskState(task *core.Task, prevState core.State) (*core.Task, error) {
	if prevState != task.State {
		if _, err := b.eventMgr.Create(b.ctx, &core.Event{
			TaskID:    task.ID,
			OccurTime: time.Now(),
			PrevState: prevState,
			CurState:  task.State,
		}); err != nil {
			return nil, err
		}
	}
	return b.taskMgr.Update(b.ctx, task)
}

func (b *Biz) DeleteTask(taskId uint32) error {
	return b.taskMgr.Delete(b.ctx, taskId)
}

func (b *Biz) List(filter core.TaskFilterParam) ([]*core.Task, error) {
	return b.taskMgr.List(b.ctx, filter)
}

func (b *Biz) GetEventByTaskID(id uint32) ([]*core.Event, error) {
	return b.eventMgr.GetByTaskID(b.ctx, id)
}
