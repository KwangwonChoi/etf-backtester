/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/KwangwonChoi/etf-backtester/pkg/analyzer"
	etf_analyzer "github.com/KwangwonChoi/etf-backtester/pkg/analyzer/etf-analyzer"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func NewAnalyzeCmd() *cobra.Command {

	info := analyzer.BackTesterConfig{}
	printLog := false
	files := ""
	startDate := ""
	endDate := ""

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
			start := GetIntegerListFromStringList(strings.Split(startDate, "-"))
			end := GetIntegerListFromStringList(strings.Split(endDate, "-"))

			info.StartDate = time.Date(start[0], time.Month(start[1]), start[2], 0, 0, 0, 0, time.UTC)
			info.EndDate = time.Date(end[0], time.Month(end[1]), end[2], 0, 0, 0, 0, time.UTC)

			backTester := etf_analyzer.NewEtfBackTester(files)

			backTester.BackTest(info, printLog)
		},
	}

	analyzeCmd.PersistentFlags().StringVar(&files, "fileName", "/Users/user/Desktop/LIGHTSRC/finance/etf-backtester/data/nasdaq100.csv", "{\"vt\" : \"/Users/user/Desktop/LIGHTSRC/finance/etf-backtester/data/nasdaq100.csv\"}")
	analyzeCmd.PersistentFlags().StringVar(&startDate, "startDate", "1985-09-26", "1985-09-26")
	analyzeCmd.PersistentFlags().StringVar(&endDate, "endDate", "2022-02-18", "2022-02-18")
	analyzeCmd.PersistentFlags().StringVar(&info.InvestPercent, "investPercent", "", "{\"vt\" : 50, \"qqq\": 30, \"tlt\": 20}")
	analyzeCmd.PersistentFlags().IntVar(&info.RebalancePeriodDay, "rebalancePeriod", 365, "365")
	analyzeCmd.PersistentFlags().StringVar(&info.LeverageMultiple, "leverage", "", "{\"vt\" : 3, \"qqq\": 3, \"tlt\": 2}")
	analyzeCmd.PersistentFlags().BoolVar(&printLog, "printLog", false, "true/false")
	analyzeCmd.PersistentFlags().BoolVar(&info.Accumulative, "accumulative", false, "true/false")
	analyzeCmd.PersistentFlags().Int64Var(&info.InitialInvestAmount, "initAmount", 100000000, "100000000")
	analyzeCmd.PersistentFlags().Int64Var(&info.AccumulativeAmount, "accumulativeAmount", 2000000, "2000000")
	analyzeCmd.PersistentFlags().IntVar(&info.AccumulativePeriodDay, "accumulativePeriod", 30, "30")

	return analyzeCmd
}

func GetIntegerListFromStringList(num []string) []int {

	returnValue := make([]int, 3)

	for i, v := range num {
		convInt, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			panic(err)
		}

		returnValue[i] = int(convInt)
	}

	return returnValue
}
