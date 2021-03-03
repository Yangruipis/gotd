package dao

import (
	"context"
	"time"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/Yangruipis/gotd/pkg/utils"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	db *gorm.DB
}

func NewTagDao(db *gorm.DB) *Tag {
	return &Tag{
		db: db,
	}
}

var _ core.TagDao = (*Tag)(nil)

func (m *Tag) withoutDelete() *gorm.DB {
	return m.db.Where("deleted_ts = 0")
}

func (m *Tag) Create(ctx context.Context, tag *core.Tag) (*core.Tag, error) {
	tag.CreateTime = time.Now()
	if err := m.db.AutoMigrate(tag).Error; err != nil {
		return tag, err
	}
	err := m.db.Create(tag).Error
	if utils.UniqueConstraintFailed(err) {
		return m.GetByName(ctx, tag.TagName)
	}
	return tag, err
}

func (m *Tag) createTaskTag(ctx context.Context, taskTag *core.TaskTag) (*core.TaskTag, error) {
	if err := m.db.AutoMigrate(taskTag).Error; err != nil {
		return taskTag, err
	}
	err := m.db.Create(taskTag).Error
	return taskTag, err
}

func (m *Tag) CreateTaskTag(ctx context.Context, tag *core.Tag, taskId uint32) (*core.Tag, error) {
	// FIXME(ryang) with transaction
	tag, err := m.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	taskTag := &core.TaskTag{
		TaskID: taskId,
		TagID:  tag.ID,
	}
	if _, err := m.createTaskTag(ctx, taskTag); err != nil {
		return tag, err
	}
	return tag, nil
}

func (m *Tag) Get(ctx context.Context, id uint32) (*core.Tag, error) {
	tag := new(core.Tag)
	if err := m.withoutDelete().First(&tag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (m *Tag) GetByName(ctx context.Context, name string) (*core.Tag, error) {
	tag := new(core.Tag)
	if err := m.withoutDelete().First(&tag, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (m *Tag) Delete(ctx context.Context, id uint32) error {
	query := m.db.Table("tag").Where("id = ?", id).
		UpdateColumn(map[string]interface{}{
			"deleted_ts": 1,
		})
	if query.Error != nil {
		return query.Error
	}

	query = m.db.Table("task_tag").
		Where("tag_id = ?", id).
		UpdateColumn(map[string]interface{}{
			"deleted_ts": 1,
		})
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (m *Tag) DeleteTaskTag(ctx context.Context, tagId uint32, taskId uint32) error {
	query := m.db.Table("task_tag").Where("task_id = ?", taskId).
		Where("tag_id = ?", tagId).
		UpdateColumn(map[string]interface{}{
			"deleted_ts": 1,
		})
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (m *Tag) List(ctx context.Context, filter core.TagFilterParam) ([]*core.Tag, error) {
	query := m.withoutDelete()
	tags := make([]*core.Tag, 0)
	if err := query.Find(&tags).Error; err != nil {
		return []*core.Tag{}, nil
	}
	return tags, nil
}
