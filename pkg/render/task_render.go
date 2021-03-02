package render

import (
	"fmt"
	"time"

	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/gookit/color"
)

type RenderOption struct {
	defaultTheme *color.Theme
	StateOption  map[core.State]*color.Theme
	PriOption    map[core.Priority]*color.Theme
}

type ColorRenderer struct {
	option RenderOption

	stateRender func(theme *color.Theme, state core.State)
	priRender   func(theme *color.Theme, pri core.Priority)
	timeRender  func(time time.Time)
	nameRender  func(name string)
	idRender    func(id uint32)
	descRender  func(desc string)
}

func NewTaskRenderer() *ColorRenderer {
	return &ColorRenderer{
		option: RenderOption{
			defaultTheme: color.Info,
			StateOption: map[core.State]*color.Theme{
				core.StateCanceled: color.Secondary,
				core.StateTodo:     color.Danger,
				core.StateSomeday:  color.Warn,
				core.StateDone:     color.Success,
				core.StateDoing:    color.Notice,
			},
			PriOption: map[core.Priority]*color.Theme{
				core.Priority0: color.Danger,
				core.Priority1: color.Notice,
				core.Priority2: color.Success,
			},
		},
		stateRender: func(theme *color.Theme, state core.State) {
			theme.Printf("%-8s ", core.StateList[int(state)-1])
		},
		priRender: func(theme *color.Theme, pri core.Priority) {
			theme.Printf("[%s]", fmt.Sprintf("P%d", pri-1))
		},
		timeRender: func(time time.Time) {
			color.FgBlue.Printf(" [%10s] ", time.Format("2006-01-02 Mon 15:04"))
		},
		nameRender: func(name string) {
			color.FgLightWhite.Printf("%s", name)
		},
		idRender: func(id uint32) {
			color.White.Printf("  %4s  ", fmt.Sprint(id)+".")
		},
		descRender: func(desc string) {
			color.White.Printf("\t%s", desc)
		},
	}
}

var _ core.IRender = (*ColorRenderer)(nil)

func (t *ColorRenderer) RenderTask(task *core.Task) (string, error) {
	t.idRender(task.ID)
	theme, ok := t.option.StateOption[task.State]
	if !ok {
		theme = t.option.defaultTheme
	}
	t.stateRender(theme, task.State)

	theme, ok = t.option.PriOption[task.Priority]
	if !ok {
		theme = t.option.defaultTheme
	}
	t.priRender(theme, task.Priority)
	t.timeRender(task.CreateTime)
	t.nameRender(task.Name)
	fmt.Print("\n")
	return "", nil
}

func (t *ColorRenderer) RenderEvents(events []*core.Event) (string, error) {
	for _, event := range events {
		color.FgGray.Print("\t- State ")
		theme, ok := t.option.StateOption[event.CurState]
		if !ok {
			theme = t.option.defaultTheme
		}
		t.stateRender(theme, event.CurState)

		color.FgGray.Print("     from ")

		theme, ok = t.option.StateOption[event.PrevState]
		if !ok {
			theme = t.option.defaultTheme
		}
		t.stateRender(theme, event.PrevState)

		t.timeRender(event.OccurTime)

		fmt.Print("\n")
	}
	return "", nil
}

func (t *ColorRenderer) RenderDesc(desc string) (string, error) {
	t.descRender(desc)
	fmt.Print("\n")
	return "", nil
}

func (t *ColorRenderer) RenderAll(task *core.Task, events []*core.Event) (string, error) {
	t.RenderTask(task)
	t.RenderEvents(events)
	t.RenderDesc(task.Description)
	return "", nil
}
