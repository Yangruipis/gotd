package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
)

const (
	buildVersion string = "0.0.1"
	buildInfo    string = ""
)

var rootCmd = &cobra.Command{
	Use:     "gotd",
	Short:   "",
	Long:    ``,
	Version: fmt.Sprintf("%s (%s)", buildVersion, buildInfo),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var loglevel uint8

func init() {
	cobra.OnInitialize(initLog)

	rootCmd.PersistentFlags().Uint8Var(&loglevel, "loglevel", 1, "Logging Level")
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))

	var out = zerolog.NewConsoleWriter()
	out.TimeFormat = time.RFC3339Nano
	out.NoColor = false
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano
	log.Logger = log.Output(out).With().Caller().Stack().Logger()
}
