package etf_analyzer

import (
	"encoding/json"
	"fmt"
	"github.com/KwangwonChoi/etf-backtester/pkg/etf"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func NewEtfBackTester(bytes []byte) EtfBackTester {

	metadata := getBackTesterMetadataList(bytes)

	start := GetIntegerListFromStringList(strings.Split(metadata.Start, "-"))
	end := GetIntegerListFromStringList(strings.Split(metadata.End, "-"))

	metadata.startDate = time.Date(start[0], time.Month(start[1]), start[2], 0, 0, 0, 0, time.UTC)
	metadata.endDate = time.Date(end[0], time.Month(end[1]), end[2], 0, 0, 0, 0, time.UTC)

	dataMap := make(map[string]etf.Info)

	for _, v := range metadata.Etf {
		dataMap[v.Name] = etf.Info{
			ETF: etf.ETF{
				Name:      v.Name,
				DailyData: etf.LoadData(v.DataFile),
			},
			MetaData: etf.MetaData{
				Leverage:      v.Leverage,
				InvestPercent: v.InvestPercent,
			},
			InvestAmount: make(map[time.Time]float64),
		}
	}

	return EtfBackTester{
		BackTestMetadata: metadata,
		DataMap:          dataMap,
		Logger:           logrus.New(),
	}
}

type EtfBackTester struct {
	BackTestMetadata
	DataMap map[string]etf.Info
	*logrus.Logger
}

type BackTestMetadata struct {
	InitialAmount    float64 `json:"initialAmount,omitempty"`
	PurchaseAmount   float64 `json:"purchaseAmount,omitempty"`
	Start            string  `json:"start,omitempty"`
	End              string  `json:"end,omitempty"`
	startDate        time.Time
	endDate          time.Time
	RebalancePeriod  int                  `json:"rebalancePeriod,omitempty"`
	Accumulative     bool                 `json:"accumulative,omitempty"`
	AccumulateAmount float64              `json:"accumulateAmount,omitempty"`
	AccumulatePeriod int                  `json:"accumulatePeriod,omitempty"`
	Etf              []IndivisualMetadata `json:"etf,omitempty"`
}

type IndivisualMetadata struct {
	Name          string  `json:"name,omitempty"`
	Alias         string  `json:"alias,omitempty"`
	DataFile      string  `json:"dataFile,omitempty"`
	Leverage      float64 `json:"leverage,omitempty"`
	InvestPercent float64 `json:"investPercent,omitempty"`
}

func (a *EtfBackTester) BackTest(printLog bool) {

	var nextAccumulateDate time.Time

	accumulateDateSet := make(map[time.Time]bool)
	rebalanceDateSet := make(map[time.Time]bool)

	initialAmount := a.InitialAmount
	purchaseAmount := a.InitialAmount

	startDate := a.startDate
	endDate := a.endDate

	rebalancePeriod := a.RebalancePeriod

	accumulative := a.Accumulative
	accumulateAmount := a.AccumulateAmount
	accumulatePeriod := a.AccumulatePeriod

	today := startDate
	nextRebalancedDate := today.AddDate(0, 0, rebalancePeriod)

	investData := a.DataMap

	if accumulative {
		nextAccumulateDate = today.AddDate(0, 0, accumulatePeriod)
	}

	// initialize by invest percent
	for _, v := range investData {
		v.InvestAmount[today] = getAmountByPercent(initialAmount, v.InvestPercent)
	}

	yesterday := today
	today = today.AddDate(0, 0, 1)

	for {
		if today.After(endDate) {
			break
		}

		// 변동 적용
		for _, v := range investData {
			v.InvestAmount[today] = applyDailyVariation(v.InvestAmount[yesterday], v.DailyData[today].ChangePercent, v.Leverage)
		}

		// 적립식 투자
		if accumulative && today == nextAccumulateDate {
			for _, v := range investData {
				v.InvestAmount[today] += getAmountByPercent(accumulateAmount, v.InvestPercent)
			}

			accumulateDateSet[today] = true
			purchaseAmount += accumulateAmount
			nextAccumulateDate = today.AddDate(0, 0, accumulatePeriod)
		}

		// rebalance
		if today == nextRebalancedDate {
			var totalAmount float64
			for _, v := range investData {
				totalAmount += v.InvestAmount[today]
			}

			for _, v := range investData {
				newAmount := getAmountByPercent(totalAmount, v.InvestPercent)
				v.InvestAmount[today] = newAmount
			}

			rebalanceDateSet[today] = true
			nextRebalancedDate = today.AddDate(0, 0, rebalancePeriod)
		}

		yesterday = today
		today = today.AddDate(0, 0, 1)
	}

	a.PurchaseAmount = purchaseAmount
	a.printResult(printLog, accumulateDateSet, rebalanceDateSet)
}

func (a *EtfBackTester) printResult(printLog bool, accumulateDateSet, rebalanceDateSet map[time.Time]bool) {

	today := a.startDate
	endDate := a.endDate
	keys := make([]string, 0)

	for k, _ := range a.DataMap {
		keys = append(keys, k)
	}

	if printLog {

		fmt.Print("Date,")

		for _, k := range keys {
			fmt.Print(k + ",")
		}

		fmt.Println()

		for {
			if today == endDate {
				break
			}

			fmt.Print(today.String() + ",")

			appendix := ""

			for _, k := range keys {
				individualData := a.DataMap[k]
				changePercent := strconv.FormatFloat(individualData.DailyData[today].ChangePercent * individualData.MetaData.Leverage, 'f', 2, 64)

				fmt.Print(strconv.FormatFloat(individualData.InvestAmount[today], 'f', 2, 64) + "("+ changePercent +"),")
			}

			if accumulateDateSet[today] {
				appendix += "Accumulate"
			}

			if rebalanceDateSet[today] {
				if appendix != "" {
					appendix += "|"
				}

				appendix += "Rebalance"
			}

			fmt.Print(appendix)

			today = today.AddDate(0, 0, 1)
			fmt.Println()
		}
	}

	totalAmount := float64(0)
	for _, v := range a.DataMap {
		totalAmount += v.InvestAmount[endDate]
	}

	totalReturnPercent := (totalAmount/a.PurchaseAmount)*100 - 100

	fmt.Println("totalPurchaseAmount, " + strconv.FormatFloat(a.PurchaseAmount, 'f', 2, 64))
	fmt.Println("totalAmount, " + strconv.FormatFloat(totalAmount, 'f', 2, 64))
	fmt.Println("totalReturnPercent, " + strconv.FormatFloat(totalReturnPercent, 'f', 2, 64) + "%")
}

func applyDailyVariation(amount, percent, leverage float64) float64 {
	amount += getAmountByPercent(amount, percent*leverage)
	return amount
}

func getAmountByPercent(amount, percent float64) float64 {
	return amount * percent / 100
}

func getBackTesterMetadataList(b []byte) BackTestMetadata {
	var metadataList BackTestMetadata

	err := json.Unmarshal(b, &metadataList)

	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to unmarshal")))
	}

	return metadataList
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
