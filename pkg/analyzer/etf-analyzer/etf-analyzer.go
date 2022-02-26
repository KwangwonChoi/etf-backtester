package etf_analyzer

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/KwangwonChoi/etf-backtester/pkg/analyzer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewEtfBackTester(files string) EtfBackTester {

	dataMap := make(map[string]map[time.Time]EtfData)
	logger := logrus.New()
	aliases := make([]string, 0)

	fileMap := make(map[string]string)
	json.Unmarshal([]byte(files), &fileMap)

	backTester := EtfBackTester{
		DataMap: dataMap,
		Logger: logger,
	}

	for key, value := range fileMap {
		backTester.DataMap[key] = make(map[time.Time]EtfData)
		backTester.loadData(key, value)
		aliases = append(aliases, key)
	}

	backTester.aliases = aliases

	return backTester
}

type EtfBackTester struct {
	DataMap map[string]map[time.Time]EtfData
	aliases []string
	*logrus.Logger
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

func (a *EtfBackTester) loadData(alias, filename string) {
	file, err := os.Open(filename)

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

		a.DataMap[alias][date] = EtfData{
			date:          date,
			endPrice:      endPrice,
			startPrice:    startPrice,
			highPrice:     highPrice,
			lowPrice:      lowPrice,
			changePercent: changePercent,
		}
	}
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
	investPercentMap := getInvestPercentMap(config.InvestPercent)
	leverageMap := getLeverageMap(config.LeverageMultiple, a.aliases)
	nextRebalanceDate := config.StartDate.AddDate(0, 0, config.RebalancePeriodDay)
	today := config.StartDate
	tomorrow := today.AddDate(0,0,1)
	resultMap := make(map[string]map[time.Time]float64)

	if config.Accumulative {
		nextAccumulateDate = config.StartDate.AddDate(0, 0, config.AccumulativePeriodDay)
	}

	for _, v := range a.aliases {
		resultMap[v] = make(map[time.Time]float64)
		resultMap[v][today] = initialAmount * (investPercentMap[v]/100)
	}

	if printLog {
		fmt.Printf("%s || ", today.String())
		for _, v := range a.aliases {
			fmt.Printf("%s : %f |", v, resultMap[v][today])
		}
		fmt.Println()
	}

	for {
		for _, v := range a.aliases {
			resultMap[v][tomorrow] = resultMap[v][today] * (1 + (a.DataMap[v][today].changePercent/100) * leverageMap[v])
		}

		if printLog {
			fmt.Printf("%s || ", tomorrow.String())
			for _, v := range a.aliases {
				fmt.Printf("%s : %f |", v, resultMap[v][tomorrow])
			}
			fmt.Println()
		}

		if config.Accumulative && tomorrow == nextAccumulateDate {

			ownAmount += float64(config.AccumulativeAmount)

			if printLog {
				fmt.Println("accumulate cash :", config.AccumulativeAmount)
			}

			for _, v := range a.aliases {
				resultMap[v][tomorrow] += float64(config.AccumulativeAmount) * (investPercentMap[v]/100)
			}

			if printLog {
				fmt.Printf("%s :: ", tomorrow.String())
				for _, v := range a.aliases {
					fmt.Printf("%s : %f |", v, resultMap[v][tomorrow])
				}
				fmt.Println()
			}

			nextAccumulateDate = tomorrow.AddDate(0, 0, config.AccumulativePeriodDay)
		}

		if tomorrow == nextRebalanceDate {

			totalAmount := float64(0)

			for _, v := range a.aliases {
				totalAmount += resultMap[v][tomorrow]
			}

			if printLog {
				fmt.Println("rebalance (inv, cash)")
			}

			for _, v := range a.aliases {
				resultMap[v][tomorrow] = totalAmount * (investPercentMap[v]/100)
			}

			if printLog {
				fmt.Printf("%s :: ", tomorrow.String())
				for _, v := range a.aliases {
					fmt.Printf("%s : %f |", v, resultMap[v][tomorrow])
				}
				fmt.Println()
			}

			nextRebalanceDate = today.AddDate(0, 0, config.RebalancePeriodDay)
		}

		today = today.AddDate(0, 0, 1)
		tomorrow = tomorrow.AddDate(0, 0, 1)

		if tomorrow == config.EndDate {
			break
		}
	}

	totalAmount := float64(0)

	for _, v := range a.aliases {
		totalAmount += resultMap[v][today]
	}

	incomeRate := (totalAmount - ownAmount) / ownAmount * 100

	fmt.Println("totalAmount:", int(totalAmount))
	fmt.Println("originAmount:", int(ownAmount))
	fmt.Println("incomeRate:", incomeRate, "%")
}

func getInvestPercentMap( jsonString string ) map[string]float64 {
	res := getFloatMap(jsonString)

	total := float64(0)

	for _, v := range res {
		total += v
	}

	if total != 100 {
		panic("percent total should be 100")
	}

	return res
}

func getLeverageMap( jsonString string, aliases []string ) map[string]float64 {
	res := getFloatMap(jsonString)

	for _, v := range aliases {
		if res[v] == 0 {
			res[v] = 1
		}
	}

	return res
}

func getFloatMap( jsonString string ) map[string]float64 {
	investPercentMap := make(map[string]float64)

	err := json.Unmarshal([]byte(jsonString), &investPercentMap)

	if err != nil {
		panic(err)
	}

	return investPercentMap
}