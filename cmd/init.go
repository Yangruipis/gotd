package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := Init(); err != nil {
				log.Fatal().Msgf("error executing command: %+v", err)
			}
		},
	}
)

func Init() error {
	return nil
}
