package analyzer

type BackTester interface {
	BackTest()
	loadData(filename string)
}