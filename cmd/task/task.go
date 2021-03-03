package task

import (
	"context"

	"github.com/Yangruipis/gotd/pkg/biz"
	dao "github.com/Yangruipis/gotd/pkg/db/gorm"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "",
	Long:  ``,
}

func NewBiz(db *gorm.DB) *biz.Biz {
	return biz.NewBiz(
		context.Background(),
		dao.NewTaskDao(db),
		dao.NewEventDao(db),
		dao.NewTagDao(db),
	)
}
