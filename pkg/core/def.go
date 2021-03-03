package core

import "time"

type State uint8

const (
	StateTodo State = iota + 1
	StateDone
	StateSomeday
	StateCanceled
	StateDoing
)

var (
	StateList = []string{"TODO", "DONE", "SOMEDAY", "CANCELED", "DOING"}
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
	Priority    Priority `gorm:"column:priority;type:int unsigned;not null;default:2"`

	State      State     `gorm:"column:state;type:int unsigned;not null"`
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;not null"`
	UpdateTime time.Time `gorm:"column:update_time;type:TIMESTAMP;not null"`
	DeletedTS  uint8     `gorm:"column:deleted_ts;type:bigint unsigned;not null"`
}

func (*Task) TableName() string {
	return "task"
}

type Event struct {
	ID        uint32    `gorm:"column:id;primary_key;auto_increment"`
	TaskID    uint32    `gorm:"column:task_id;type:int unsigned"`
	OccurTime time.Time `gorm:"column:occur_time;type:TIMESTAMP;not null"`
	PrevState State     `gorm:"column:prev_state;type:int unsigned;not null"`
	CurState  State     `gorm:"column:cur_state;type:int unsigned;not null"`
}

func (*Event) TableName() string {
	return "event"
}

func (*Event) ForeignKeys() map[string]string {
	return map[string]string{
		"task_id": "task(id)",
	}
}

type Tag struct {
	ID         uint32    `gorm:"column:id;primary_key;auto_increment"`
	TagName    string    `gorm:"column:name;type:varchar(64);not null;UNIQUE"`
	CreateTime time.Time `gorm:"column:create_time;type:TIMESTAMP;not null"`
	DeletedTS  uint8     `gorm:"column:deleted_ts;type:bigint unsigned;not null"`
}

func (*Tag) TableName() string {
	return "tag"
}

type TaskTag struct {
	ID        uint32 `gorm:"column:id;primary_key;auto_increment"`
	TaskID    uint32 `gorm:"column:task_id;type:int unsigned"`
	TagID     uint32 `gorm:"column:tag_id;type:int unsigned"`
	DeletedTS uint8  `gorm:"column:deleted_ts;type:bigint unsigned;not null"`
}

func (*TaskTag) TableName() string {
	return "task_tag"
}

func (*TaskTag) ForeignKeys() map[string]string {
	return map[string]string{
		"task_id": "task(id)",
		"tag_id":  "tag(id)",
	}
}

type AgendaType uint8

const (
	AgendaTypeDaily AgendaType = iota + 1
	AgendaTypeWeekly
	AgendaTypeMonthly
)

type SkipType uint8

const (
	SkipTypeOnWeekend SkipType = iota + 1
	SkipTypeOnHoliday
	SkipTypeOnWeekendAndHoliday
)

type Agenda struct {
	*Task

	ScheduleTime time.Time
	DeadlineTIme time.Time
	AgendaType   AgendaType
	SkipType     SkipType
}

type Report struct {
}
