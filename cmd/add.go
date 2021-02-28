package cmd

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
	Desc string
}

var (
	ctx AddContext
	cmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Add(&ctx); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVar(&ctx.Desc, "desc", "", "")
	cmd.MarkFlagRequired("desc")
}

func Add(ctx *AddContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()

	task := biz.NewTask(context.Background(), dao.NewTaskManager(db))

	priorities := []core.Priority{core.Priority0, core.Priority1, core.Priority2}
	prompt := promptui.Select{
		Label: "Select Day",
		Items: priorities,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		log.Info().Msgf("Prompt failed %v\n", err)
		return err
	}

	_, err = task.Create(&core.Task{
		Description: ctx.Desc,
		Priority:    priorities[idx],
		Name:        "test",
		State:       core.StateTodo,
	})

	return err
}
