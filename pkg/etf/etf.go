package etf

import (
	"github.com/KwangwonChoi/etf-backtester/pkg/common"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type DailyData struct {
	Date               time.Time
	EndPrice           float64
	StartPrice         float64
	HighPrice          float64
	LowPrice           float64
	TradingVolume      string // 거래량
	ChangePercent      float64
}

type ETF struct {
	Name      string
	DailyData map[time.Time]DailyData
}

type MetaData struct {
	Leverage      float64
	InvestPercent float64
}

type Info struct {
	ETF
	MetaData
	InvestAmount map[time.Time]float64
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

func LoadData(filename string) map[time.Time]DailyData {

	rows := common.ReadFileMatrix(filename)

	dailyData := make(map[time.Time]DailyData)

	for _, row := range rows {
		date, err := getDateFromKRData(row[dateColumn])

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
			panic(errors.Wrap(err, "failed parse startPrice value"))
		}

		highPrice, err := strconv.ParseFloat(textReplacer.Replace(row[highPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse highPrice value"))
		}

		lowPrice, err := strconv.ParseFloat(textReplacer.Replace(row[lowPriceColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse lowPrice value"))
		}

		changePercent, err := strconv.ParseFloat(textReplacer.Replace(row[changePercentColumn]), 64)

		if err != nil {
			panic(errors.Wrap(err, "failed parse changePercent value"))
		}

		dailyData[date] = DailyData{
			Date:          date,
			EndPrice:      endPrice,
			StartPrice:    startPrice,
			HighPrice:     highPrice,
			LowPrice:      lowPrice,
			ChangePercent: changePercent,
		}
	}

	return dailyData
}

func getDateFromKRData(data string) (time.Time, error) {

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
