package task

import (
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DelContext struct {
	id []uint
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
	TaskCmd.AddCommand(delCmd)

	delCmd.Flags().UintSliceVarP(&delCtx.id, "id", "i", []uint{}, "")
	delCmd.MarkFlagRequired("id")
}

func Del(ctx *DelContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()

	biz := NewBiz(db)
	for _, id := range ctx.id {
		err := biz.DeleteTask(uint32(id))
		if err != nil {
			return err
		}
		log.Info().Msgf("task %d is deleted", id)
	}
	return nil
}
