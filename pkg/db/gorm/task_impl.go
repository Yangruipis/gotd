package dao

import (
	"context"
	"time"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Task struct {
	db *gorm.DB
}

func NewTaskManager(db *gorm.DB) *Task {
	return &Task{
		db: db,
	}
}

var _ core.TaskManager = (*Task)(nil)

func (m *Task) withoutDelete() *gorm.DB {
	return m.db.Where("deleted_ts = 0")
}

func (m *Task) Create(ctx context.Context, task *core.Task) (*core.Task, error) {
	task.CreateTime = time.Now()
	task.UpdateTime = time.Now()
	log.Debug().Msgf("task: %+v", task)
	if err := m.db.AutoMigrate(task).Error; err != nil {
		return task, err
	}
	err := m.db.Create(task).Error
	return task, err
}

func (m *Task) Update(ctx context.Context, task *core.Task) (*core.Task, error) {
	task.UpdateTime = time.Now()
	log.Debug().Msgf("task: %+v", task)
	err := m.db.Save(task).Error
	return task, err
}

func (m *Task) Get(ctx context.Context, id uint32) (*core.Task, error) {
	task := new(core.Task)
	if err := m.withoutDelete().First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return task, nil
}

func (m *Task) Delete(ctx context.Context, id uint32) error {
	query := m.db.Table("task").Where("id = ?", id).
		UpdateColumn(map[string]interface{}{
			"update_time": time.Now(),
			"deleted_ts":  1,
		})
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (m *Task) List(ctx context.Context, filter core.TaskFilterParam) ([]*core.Task, error) {
	query := m.withoutDelete().Table("task")

	if filter.Priority != 0 {
		query.Where("priority = ?", filter.Priority)
	}
	if filter.State != 0 {
		query.Where("state = ?", filter.State)
	}
	if filter.NameKeyword != "" {
		query.Where("name LIKE ?", "%"+filter.NameKeyword+"%")
	}
	if filter.DescKeyword != "" {
		query.Where("description LIKE ?", "%"+filter.DescKeyword+"%")
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if filter.MinTime != 0 {
		minTime := time.Unix(int64(filter.MinTime), 0).In(loc)
		query.Where("create_time > ?", minTime.Format(time.RFC3339))
	}
	if filter.MaxTime != 0 && filter.MaxTime > filter.MinTime {
		maxTime := time.Unix(int64(filter.MaxTime), 0).In(loc)
		query.Where("create_time < ?", maxTime.Format(time.RFC3339))
	}

	tasks := make([]*core.Task, 0)
	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
