package utils

import (
	"strconv"
	"encoding/csv"
	"os"
)

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}


// FloatToString ...
func FloatToString(num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(num, 'f', 6, 64)
}

func StringToFloat64(nums string) (float64, error) {
		num, err := strconv.ParseFloat(nums, 64)
		if err != nil {
			 //panic(err)
			 return 0.0, err
		}

	return num, err
}
