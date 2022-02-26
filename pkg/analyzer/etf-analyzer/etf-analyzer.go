package etf_analyzer

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/KwangwonChoi/etf-backtester/pkg/analyzer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewEtfBackTester(filename []string) EtfBackTester {

	dataMap := make(map[string]map[time.Time]EtfData)
	logger := logrus.New()

	backTester := EtfBackTester{
		dataMap,
		logger,
	}

	backTester.loadDatas(filename)

	return backTester
}

type EtfBackTester struct {
	DataMap *map[string]map[time.Time]EtfData
	*logrus.Logger
}

type Etf struct {
	name string
	filename string
	data map[time.Time]EtfData
}

type EtfData struct {
	date          time.Time
	endPrice      float64
	startPrice    float64
	highPrice     float64
	lowPrice      float64
	tradingVolume string // 거래량
	changePercent float64
}

const (
	dateColumn          = 0
	endPriceColumn      = 1
	startPriceColumn    = 2
	highPriceColumn     = 3
	lowPriceColumn      = 4
	tradingVolumeColumn = 5
	changePercentColumn = 6
)

func (a *EtfBackTester) loadDatas(filenames []string) {

	datas := make(map[string]map[time.Time]EtfData)

	for _, v := range filenames {
	fileNameList := strings.Split(v, "/")
		datas[fileNameList[len(fileNameList)]] = a.loadData(v)
	}

	a.DataMap = datas
}

func (a *EtfBackTester) loadData(filename string) map[time.Time]EtfData {
	file, err := os.Open(filename)
	returnValue := make(map[time.Time]EtfData)

	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to load file %s", filename)))
	}

	rdr := csv.NewReader(bufio.NewReader(file))
	rdr.LazyQuotes = true

	rows, err := rdr.ReadAll()

	if err != nil {
		panic(errors.Wrap(err, "failed to create new reader"))
	}

	for _, row := range rows {
		date, err := a.getDateFromKRData(row[dateColumn])

		if err != nil {
			panic(errors.Wrap(err, "failed to get data from date format"))
		}

		textReplacer := strings.NewReplacer(",", "", "%", "")

		endPrice, err := strconv.ParseFloat(textReplacer.Replace(row[endPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse endPrice value"))
		}

		startPrice, err := strconv.ParseFloat(textReplacer.Replace(row[startPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse endPrice value"))
		}

		highPrice, err := strconv.ParseFloat(textReplacer.Replace(row[highPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse endPrice value"))
		}

		lowPrice, err := strconv.ParseFloat(textReplacer.Replace(row[lowPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse endPrice value"))
		}

		changePercent, err := strconv.ParseFloat(textReplacer.Replace(row[changePercentColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse endPrice value"))
		}

		returnValue[date] = EtfData{
			date:          date,
			endPrice:      endPrice,
			startPrice:    startPrice,
			highPrice:     highPrice,
			lowPrice:      lowPrice,
			changePercent: changePercent,
		}
	}

	return returnValue
}

func (a *EtfBackTester) getDateFromKRData(data string) (time.Time, error) {

	a.Logger.Debugln("dateData: ", data)
	replacer := strings.NewReplacer("\"", "", " ", "", "\ufeff", "")

	yearStrList := strings.Split(data, "년")
	monthStrList := strings.Split(yearStrList[1], "월")
	dayStrList := strings.Split(monthStrList[1], "일")

	yearStr := replacer.Replace(yearStrList[0])
	monthStr := replacer.Replace(monthStrList[0])
	dayStr := replacer.Replace(dayStrList[0])

	year, _ := strconv.ParseInt(yearStr, 10, 64)
	month, _ := strconv.ParseInt(monthStr, 10, 64)
	day, _ := strconv.ParseInt(dayStr, 10, 64)

	dateStr := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)

	return dateStr, nil
}

func (a *EtfBackTester) BackTest(config analyzer.BackTesterConfig, printLog bool) {

	var nextAccumulateDate time.Time
	initialAmount := float64(config.InitialInvestAmount)
	ownAmount := initialAmount
	nextRebalanceDate := config.StartDate.AddDate(0, 0, config.RebalancePeriodDay)
	today := config.StartDate

	if config.Accumulative {
		nextAccumulateDate = config.StartDate.AddDate(0, 0, config.AccumulativePeriodDay)
	}

	invAmount := float64(initialAmount * (config.InvestPercent / 100))
	cashAmount := float64(initialAmount - invAmount)

	for {
		invAmount = invAmount + (invAmount * (a.DataMap[today].changePercent / 100) * config.LeverageMultiple)

		if config.Accumulative && today == nextAccumulateDate {
			inv := float64(float64(config.AccumulativeAmount) * (config.InvestPercent / 100))
			cash := float64(config.AccumulativeAmount) - inv
			ownAmount += float64(config.AccumulativeAmount)

			invAmount += inv
			cashAmount += cash

			if printLog {
				fmt.Println("accumulate cash :", config.AccumulativeAmount, "(inv, cash)", invAmount, cashAmount)
			}

			nextAccumulateDate = today.AddDate(0, 0, config.AccumulativePeriodDay)
		}

		if today == nextRebalanceDate {
			nextRebalanceDate = today.AddDate(0, 0, config.RebalancePeriodDay)

			totalAmount := invAmount + cashAmount
			invAmount = totalAmount * (config.InvestPercent / 100)
			cashAmount = totalAmount - invAmount

			if printLog {
				fmt.Println("rebalance (inv, cash)", invAmount, cashAmount)
			}
		}

		today = today.AddDate(0, 0, 1)

		if printLog {
			fmt.Println(strings.Split(today.String(), " ")[0], " Amount(inv, cash) ", int64(invAmount), int64(cashAmount))
		}

		if today == config.EndDate {
			break
		}
	}

	totalAmount := invAmount + cashAmount
	incomeRate := (totalAmount - ownAmount) / ownAmount * 100

	fmt.Println("totalAmount:", int(totalAmount))
	fmt.Println("incomeRate:", incomeRate, "%")
}
