package command

import (
	"fmt"
	"os"

	"github.com/antelman107/parser/app"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "app",
}

func Execute() {
	csvAppBuilder = func(k uint32, cfg app.CsvParserAppConfig) app.CsvApp {
		return app.NewCsvParserApp(k, cfg)
	}
	rootCmd.AddCommand(topCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
