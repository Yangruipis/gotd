package dao

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
)

type EventManager struct {
	db *gorm.DB
}

var _ core.EventManager = (*EventManager)(nil)

func (m *EventManager) CreateOrUpdate(ctx context.Context, task *core.Event) (*core.Event, error) {
	return nil, nil
}

func (m *EventManager) Get(ctx context.Context, id string) (*core.Event, error) {
	return nil, nil
}

func (m *EventManager) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *EventManager) List(ctx context.Context, filter core.EventFilterParam) ([]*core.Event, error) {
	return nil, nil
}
