package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/logviewer/v2/src/pkg/app"
	"github.com/logviewer/v2/src/pkg/common"
	"github.com/logviewer/v2/src/pkg/parser"
	"github.com/logviewer/v2/src/pkg/source"
	"github.com/spf13/cobra"
)

var (
	// source cfg
	cliDebugOutput string
	cliSourceType  string
	cliSourcePath  string
	// parser cfg
	cliParserType string
	// app cfg
)

var rootCmd = &cobra.Command{
	Use:   "logviewer",
	Short: "Logviewer is a fast & interactive structured log viewer",
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) error {
	// init log
	common.ConfigureGlobalLog(cliDebugOutput)

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

	// run application
	return app.Main(ps)
}

func main() {
	rootCmd.Flags().StringVarP(&cliDebugOutput, "debug", "", "", "debug output")
	rootCmd.Flags().StringVarP(&cliSourceType, "source", "s", "file", "source type [file]")
	rootCmd.Flags().StringVarP(&cliSourcePath, "from", "f", "", "source input")
	rootCmd.MarkFlagRequired("from")
	rootCmd.Flags().StringVarP(&cliParserType, "parser", "p", "zap", "parser type [zap]")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
