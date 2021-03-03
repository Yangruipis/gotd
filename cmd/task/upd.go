package task

import (
	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type UpdContext struct {
	id         uint32
	updateDesc bool
}

var (
	updCtx UpdContext
	updCmd = &cobra.Command{
		Use:   "upd",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Upd(&updCtx); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func init() {
	TaskCmd.AddCommand(updCmd)

	updCmd.Flags().Uint32VarP(&updCtx.id, "id", "i", 0, "")
	updCmd.Flags().BoolVarP(&updCtx.updateDesc, "desc", "d", false, "")
	updCmd.MarkFlagRequired("id")
}

func Upd(ctx *UpdContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()
	biz := NewBiz(db)

	taskGot, err := biz.GetTask(ctx.id)
	if err != nil {
		return err
	}
	priorities := []core.Priority{core.Priority0, core.Priority1, core.Priority2}
	prompt := promptui.Select{
		Label:     "Priority",
		Items:     priorities,
		CursorPos: int(taskGot.Priority) - 1,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}
	taskGot.Priority = priorities[idx]

	prompt = promptui.Select{
		Label:     "State",
		Items:     core.StateList,
		CursorPos: int(taskGot.State) - 1,
	}
	idx, _, err = prompt.Run()
	if err != nil {
		log.Error().Err(err).Msgf("Prompt failed %v\n", err)
		return err
	}
	prevState := taskGot.State
	taskGot.State = core.State(idx + 1)

	if updCtx.updateDesc {
		promptStr := promptui.Prompt{
			Label: "Description",
		}
		result, err := promptStr.Run()
		if err != nil {
			log.Error().Err(err).Msgf("Prompt failed %v\n", err)
			return err
		}
		taskGot.Description = result
	}
	_, err = biz.UpdateTaskState(taskGot, prevState)

	return err
}
