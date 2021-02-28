package core

import "time"

type State uint8

const (
	StateTodo State = iota + 1
	StateDone
	StateSomeday
	StateCanceled
	StateFinished
	StateDoing
)

type Priority uint8

const (
	Priority0 Priority = iota + 1
	Priority1
	Priority2
)

type Task struct {

	ID          uint32   `gorm:"column:id;primary_key;auto_increment"`
	Name        string   `gorm:"column:name;type:varchar(64);not null"`
	Description string   `gorm:"column:description;type:text;not null"`
	Priority    Priority `gorm:"column:priority;type:int unsigned;not null"`

	State      State     `gorm:"column:state;type:int unsigned;not null"`
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;not null"`
	DeletedTS  uint8     `gorm:"column:deleted_ts;type:bigint unsigned;not null"`
}

func (*Task) TableName() string {
	return "task"
}

type Event struct {
	ID uint32

	TaskID    uint32
	Time      time.Time
	PrevState State
	CurState  State
}

type Tag struct {
	ID uint32

	TaskID  uint32
	TagName string
}

type Report struct {
}
