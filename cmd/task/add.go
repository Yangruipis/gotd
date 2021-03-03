package task

import (
	"github.com/Yangruipis/gotd/pkg/core"
	"github.com/jinzhu/gorm"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type AddContext struct {
	Name string

	tag bool
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
	TaskCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addCtx.Name, "name", "n", "", "")
	addCmd.Flags().BoolVarP(&addCtx.tag, "tag", "t", false, "")
	addCmd.MarkFlagRequired("name")
}

func Add(ctx *AddContext) error {

	db, err := gorm.Open("sqlite3", ".sqlite3.db")
	if err != nil {
		return err
	}
	defer db.Close()

	biz := NewBiz(db)

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

	taskGot, err := biz.CreateTask(&core.Task{
		Name:        ctx.Name,
		Description: result,
		Priority:    priorities[idx],
		State:       core.StateTodo,
	})
	if err != nil {
		return err
	}

	if ctx.tag {
		tags, err := biz.ListAllTag()
		if err != nil {
			return err
		}
		index := -1

		tagNames := make([]string, 0, len(tags))
		for _, tag := range tags {
			tagNames = append(tagNames, tag.TagName)
		}
		for index < 0 {
			prompt := promptui.SelectWithAdd{
				Label:    "Tag Name",
				Items:    tagNames,
				AddLabel: "New Tag",
			}

			index, result, err = prompt.Run()
			if index == -1 {
				tagNames = append(tagNames, result)
			}
		}
		if err != nil {
			log.Error().Err(err).Msgf("Prompt failed %v\n", err)
			return err
		}
		if _, err := biz.CreateTaskTag(tagNames[index], taskGot.ID); err != nil {
			return err
		}

	}
	return nil
}
