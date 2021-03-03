package todo

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/biz"
	"github.com/Yangruipis/gotd/pkg/core"
	dao "github.com/Yangruipis/gotd/pkg/db/gorm"
	"github.com/jinzhu/gorm"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type AddContext struct {
	Name string
}

var (
	addCtx AddContext
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Add(&addCtx); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func init() {
	TodoCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addCtx.Name, "name", "n", "", "")
	addCmd.MarkFlagRequired("name")
}

func Add(ctx *AddContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()

	task := biz.NewBiz(context.Background(), dao.NewTaskManager(db), dao.NewEventManager(db))

	priorities := []core.Priority{core.Priority0, core.Priority1, core.Priority2}
	prompt := promptui.Select{
		Label: "Priority",
		Items: priorities,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}

	promptStr := promptui.Prompt{
		Label: "Description",
	}
	result, err := promptStr.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}

	_, err = task.CreateTask(&core.Task{
		Name:        ctx.Name,
		Description: result,
		Priority:    priorities[idx],
		State:       core.StateTodo,
	})

	return err
}
