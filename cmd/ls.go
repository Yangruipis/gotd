package cmd

import (
	"context"
	"time"

	"github.com/Yangruipis/gotd/pkg/biz"
	"github.com/Yangruipis/gotd/pkg/core"
	dao "github.com/Yangruipis/gotd/pkg/db/gorm"
	"github.com/Yangruipis/gotd/pkg/render"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ListContext struct {
	id     uint32
	detail bool

	week uint32
	day  uint32

	state    bool
	priority bool

	keyword bool
}

var (
	listCtx ListContext
	listCmd = &cobra.Command{
		Use:   "ls",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := List(&listCtx); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().Uint32VarP(&listCtx.id, "id", "i", 0, "")
	listCmd.Flags().BoolVarP(&listCtx.detail, "detail", "d", false, "")

	listCmd.Flags().Uint32Var(&listCtx.week, "week", 0, "")
	listCmd.Flags().Uint32Var(&listCtx.day, "day", 0, "")

	listCmd.Flags().BoolVarP(&listCtx.state, "state", "s", false, "")
	listCmd.Flags().BoolVarP(&listCtx.priority, "prior", "p", false, "")

	listCmd.Flags().BoolVarP(&listCtx.keyword, "keyword", "k", false, "")

}

func List(ctx *ListContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()
	gotdBiz := biz.NewBiz(context.Background(), dao.NewTaskManager(db), dao.NewEventManager(db))

	if ctx.id > 0 {
		taskGot, err := gotdBiz.GetTask(ctx.id)
		if err != nil {
			return err
		}
		err = renderTasks(ctx, []*core.Task{taskGot}, gotdBiz)
	} else {
		filter := &core.TaskFilterParam{}
		if err := parseTime(ctx, filter); err != nil {
			return err
		}

		if err := parseState(ctx, filter); err != nil {
			return err
		}

		if err := parsePrior(ctx, filter); err != nil {
			return err
		}

		if err := parseKeyword(ctx, filter); err != nil {
			return err
		}

		tasksGot, err := gotdBiz.List(*filter)
		if err != nil {
			return err
		}
		err = renderTasks(ctx, tasksGot, gotdBiz)
	}
	return err
}

func renderTasks(ctx *ListContext, tasksGot []*core.Task, gotdBiz *biz.Biz) error {
	renderer := render.NewTaskRenderer()
	for _, taskGot := range tasksGot {
		if ctx.detail {
			events, err := gotdBiz.GetEventByTaskID(taskGot.ID)
			if err != nil {
				return err
			}
			if _, err := renderer.RenderAll(taskGot, events); err != nil {
				return err
			}
		} else {
			if _, err := renderer.RenderTask(taskGot); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseTime(ctx *ListContext, filter *core.TaskFilterParam) error {
	if ctx.day != 0 && ctx.week != 0 {
		log.Warn().Msgf("argument 'week' and 'day' are conflict, only use 'week'")
	}

	if ctx.week != 0 {
		now.WeekStartDay = time.Monday
		tBegin := now.BeginningOfWeek().Unix()
		tEnd := now.EndOfWeek().Unix()
		gap := tEnd - tBegin
		filter.MinTime = tBegin - int64(ctx.week-1)*(gap)
		filter.MaxTime = tEnd - int64(ctx.week-1)*(gap)
	} else if ctx.day != 0 {
		tBegin := now.BeginningOfDay().Unix()
		tEnd := now.EndOfDay().Unix()
		gap := tEnd - tBegin
		filter.MinTime = tBegin - int64(ctx.day-1)*(gap)
		filter.MaxTime = tEnd - int64(ctx.day-1)*(gap)
	}
	return nil
}

func parseState(ctx *ListContext, filter *core.TaskFilterParam) error {
	if !ctx.state {
		return nil
	}
	tmpList := append([]string{"ALL"}, core.StateList...)
	prompt := promptui.Select{
		Label:     "State",
		Items:     tmpList,
		CursorPos: 0,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}
	if idx == 0 {
		return nil
	} else {
		filter.State = uint8(idx)
	}
	return nil
}

func parsePrior(ctx *ListContext, filter *core.TaskFilterParam) error {
	if !ctx.priority {
		return nil
	}
	priorities := []string{"ALL", "P0", "P1", "P2"}
	prompt := promptui.Select{
		Label:     "Priority",
		Items:     priorities,
		CursorPos: 0,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}
	if idx == 0 {
		return nil
	} else {
		filter.Priority = uint8(idx)
	}
	return nil
}

func parseKeyword(ctx *ListContext, filter *core.TaskFilterParam) error {
	if !ctx.keyword {
		return nil
	}
	promptStr := promptui.Prompt{
		Label: "Search Keyword",
	}
	result, err := promptStr.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}
	if result != "" {
		filter.NameKeyword = result
		filter.DescKeyword = result
	}
	return nil
}
