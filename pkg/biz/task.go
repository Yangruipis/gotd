package biz

import (
	"context"
	"time"

	"github.com/Yangruipis/gotd/pkg/core"
)

type Biz struct {
	ctx      context.Context
	taskDao  core.TaskDao
	eventDao core.EventDao
	tagDao   core.TagDao
}

func NewBiz(ctx context.Context,
	taskDao core.TaskDao,
	eventDao core.EventDao,
	tagDao core.TagDao) *Biz {
	return &Biz{
		ctx:      ctx,
		taskDao:  taskDao,
		eventDao: eventDao,
		tagDao:   tagDao,
	}
}

func (b *Biz) CreateTask(task *core.Task) (*core.Task, error) {
	if _, err := b.eventDao.Create(b.ctx, &core.Event{
		TaskID:    task.ID,
		OccurTime: time.Now(),
		PrevState: 0,
		CurState:  task.State,
	}); err != nil {
		return nil, err
	}
	return b.taskDao.Create(b.ctx, task)
}

func (b *Biz) CreateTaskWithTags(task *core.Task, tags ...*core.Tag) (*core.Task, error) {
	task, err := b.CreateTask(task)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		_, err := b.tagDao.CreateTaskTag(b.ctx, tag, task.ID)
		if err != nil {
			return task, err
		}
	}
	return task, nil
}

func (b *Biz) GetTask(id uint32) (*core.Task, error) {
	return b.taskDao.Get(b.ctx, id)
}

func (b *Biz) UpdateTaskState(task *core.Task, prevState core.State) (*core.Task, error) {
	if prevState != task.State {
		if _, err := b.eventDao.Create(b.ctx, &core.Event{
			TaskID:    task.ID,
			OccurTime: time.Now(),
			PrevState: prevState,
			CurState:  task.State,
		}); err != nil {
			return nil, err
		}
	}
	return b.taskDao.Update(b.ctx, task)
}

func (b *Biz) DeleteTask(taskId uint32) error {
	return b.taskDao.Delete(b.ctx, taskId)
}

func (b *Biz) List(filter core.TaskFilterParam) ([]*core.Task, error) {
	return b.taskDao.List(b.ctx, filter)
}

func (b *Biz) GetEventByTaskID(id uint32) ([]*core.Event, error) {
	return b.eventDao.GetByTaskID(b.ctx, id)
}
