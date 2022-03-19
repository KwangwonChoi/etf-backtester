package common

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func ReadFileMatrix(filename string) [][]string {
	file, err := os.Open(filename)

	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to load file %s", filename)))
	}

	defer file.Close()

	rdr := csv.NewReader(bufio.NewReader(file))
	rdr.LazyQuotes = true

	rows, err := rdr.ReadAll()

	if err != nil {
		panic(errors.Wrap(err, "failed to create new reader"))
	}

	return rows
}

func ReadFileByte(filename string) []byte {

	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to load file %s", filename)))
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to read file %s", filename)))
	}

	return byteValue
}
