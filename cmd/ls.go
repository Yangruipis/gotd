package cmd

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/biz"
	"github.com/Yangruipis/gotd/pkg/core"
	dao "github.com/Yangruipis/gotd/pkg/db/gorm"
	"github.com/Yangruipis/gotd/pkg/render"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type ListContext struct {
	id     uint32
	detail bool
	month  uint32
	week   uint32
	day    uint32
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
	listCmd.Flags().BoolVar(&listCtx.detail, "detail", false, "")

	listCmd.Flags().Uint32Var(&listCtx.month, "month", 0, "")
	listCmd.Flags().Uint32Var(&listCtx.week, "week", 0, "")
	listCmd.Flags().Uint32Var(&listCtx.day, "day", 0, "")
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

		tasksGot, err := gotdBiz.ListAll()
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

func parseTime(ctx *ListContext) (uint64, uint64, error) {

	if ctx.week != 0 {
		now.BeginningOfWeek()
	}

	return 0, 0, nil
}
