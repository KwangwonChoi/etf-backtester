/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	etf_analyzer "github.com/KwangwonChoi/etf-backtester/pkg/analyzer/etf-analyzer"
	"github.com/KwangwonChoi/etf-backtester/pkg/common"
	"github.com/spf13/cobra"
)

func NewAnalyzeCmd() *cobra.Command {

	printLog := false
	files := ""

	// analyzeCmd represents the analyze command
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {

			bytes := common.ReadFileByte(files)

			backTester := etf_analyzer.NewEtfBackTester(bytes)

			backTester.BackTest(printLog)
		},
	}

	analyzeCmd.PersistentFlags().StringVar(&files, "fileName", "/Users/user/Desktop/LIGHTSRC/finance/etf-backtester/data/nasdaq100.csv", "{\"vt\" : \"/Users/user/Desktop/LIGHTSRC/finance/etf-backtester/data/nasdaq100.csv\"}")
	analyzeCmd.PersistentFlags().BoolVar(&printLog, "printLog", false, "true/false")

	return analyzeCmd
}