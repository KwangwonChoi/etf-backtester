package etf_analyzer

import (
	"github.com/KwangwonChoi/etf-backtester/pkg/analyzer"
	"testing"
	"time"
)

func TestGetData(t *testing.T) {
	backTester := NewEtfBackTester("/Users/user/Desktop/LIGHTSRC/finance/etf-backtester/nasdaq100.csv")

	startDate := time.Date(1999, 04, 27, 0, 0, 0, 0, time.UTC)
	//endDate := time.Date(1999, 04, 28, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2022, 02, 01, 0, 0, 0, 0, time.UTC)
	investPercent := float64(50)
	rebalancePeriod := int(365)
	leverate := float64(3)
	initialAmount := int64(100)

	backTester.BackTest(analyzer.BackTesterConfig{
		StartDate:           startDate,
		EndDate:             endDate,
		InvestPercent:       investPercent,
		RebalancePeriodDay:  rebalancePeriod,
		LeverageMultiple:    leverate,
		InitialInvestAmount: initialAmount,
	})
}
