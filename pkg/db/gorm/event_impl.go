package dao

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Event struct {
	db *gorm.DB
}

func NewEventDao(db *gorm.DB) *Event {
	return &Event{
		db: db,
	}
}

var _ core.EventDao = (*Event)(nil)

func (m *Event) withoutDelete() *gorm.DB {
	return m.db
}

func (m *Event) Create(ctx context.Context, event *core.Event) (*core.Event, error) {
	log.Debug().Msgf("event: %+v", event)
	if err := m.db.AutoMigrate(event).Error; err != nil {
		return event, err
	}
	err := m.db.Create(event).Error
	return event, err
}

func (m *Event) Update(ctx context.Context, event *core.Event) (*core.Event, error) {
	log.Debug().Msgf("event: %+v", event)
	err := m.db.Save(event).Error
	return event, err
}

func (m *Event) Get(ctx context.Context, id uint32) (*core.Event, error) {
	event := new(core.Event)
	if err := m.withoutDelete().First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return event, nil
}

func (m *Event) GetByTaskID(ctx context.Context, taskId uint32) ([]*core.Event, error) {
	events := make([]*core.Event, 0)
	if err := m.withoutDelete().Find(&events, "task_id = ?", taskId).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (m *Event) List(ctx context.Context, filter core.EventFilterParam) ([]*core.Event, error) {
	return nil, nil
}
