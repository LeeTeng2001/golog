package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/logviewer/v2/src/pkg/parser"
	"github.com/logviewer/v2/src/pkg/source"
	"github.com/spf13/cobra"
)

var (
	// source cfg
	cliSourceType string
	cliSourcePath string
	// parser cfg
	cliParserType string
	// app cfg
)

func execute(cmd *cobra.Command, args []string) error {
	// init source
	rd := source.AllReaders[cliSourceType]
	if rd == nil {
		return errors.New("source does not exist for: " + cliSourceType)
	}
	if err := rd.Init(cliSourcePath); err != nil {
		return fmt.Errorf("failed to init source: %w", err)
	}

	// init parser
	ps := parser.AllParser[cliParserType]
	if ps == nil {
		return errors.New("parser does not exist for: " + cliParserType)
	}
	if err := ps.Init(rd); err != nil {
		return fmt.Errorf("failed to init parser: %w", err)
	}

	ps.GetLogs(0, 10)

	// run application
	// return app.Main(ps)
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "logviewer",
	Short: "Logviewer is a fast & interactive structured log viewer",
	RunE:  execute,
}

func main() {
	rootCmd.Flags().StringVarP(&cliSourceType, "source", "s", "file", "source type [file]")
	rootCmd.Flags().StringVarP(&cliSourcePath, "from", "f", "", "source input")
	rootCmd.MarkFlagRequired("from")
	rootCmd.Flags().StringVarP(&cliParserType, "parser", "p", "zap", "parser type [zap]")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
