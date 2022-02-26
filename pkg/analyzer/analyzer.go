package analyzer

import "time"

type BackTester interface {
	BackTest()
	loadData(filename string)
}

type BackTesterConfig struct {
	StartDate             time.Time
	EndDate               time.Time
	InvestPercent         string
	RebalancePeriodDay    int
	LeverageMultiple      string
	InitialInvestAmount   int64
	Accumulative          bool
	AccumulativePeriodDay int
	AccumulativeAmount    int64
}
