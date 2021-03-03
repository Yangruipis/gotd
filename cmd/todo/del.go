package todo

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/biz"
	dao "github.com/Yangruipis/gotd/pkg/db/gorm"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DelContext struct {
	id uint32
}

var (
	delCtx DelContext
	delCmd = &cobra.Command{
		Use:   "del",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Del(&delCtx); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func init() {
	TodoCmd.AddCommand(delCmd)

	delCmd.Flags().Uint32VarP(&delCtx.id, "id", "i", 0, "")
	delCmd.MarkFlagRequired("id")
}

func Del(ctx *DelContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()

	task := biz.NewBiz(context.Background(), dao.NewTaskManager(db), dao.NewEventManager(db))
	return task.DeleteTask(ctx.id)
}
