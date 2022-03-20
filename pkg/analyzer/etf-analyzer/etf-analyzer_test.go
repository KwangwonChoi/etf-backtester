package etf_analyzer

import (
	"testing"
)

func TestGetData(t *testing.T) {

	jsonData := `
{
  "initialAmount": 10000,
  "start": "2001-01-01",
  "end": "2012-01-01",
  "rebalancePeriod": 365,
  "accumulative": true,
  "accumulateAmount": 100,
  "accumulatePeriod": 30,
  "etf": [{
    "name": "vt",
    "alias": "vt",
    "dataFile": "../../../data/vt.csv",
    "leverage": 1,
    "investPercent": 50
  }, {
    "name": "tqqq",
    "alias": "vt",
    "dataFile": "../../../data/nasdaq100.csv",
    "leverage": 3,
    "investPercent": 50
  }]
}
`

	backTester := NewEtfBackTester([]byte(jsonData))
	backTester.BackTest(true)

}
