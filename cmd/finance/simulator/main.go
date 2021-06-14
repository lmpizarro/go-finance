package main

import (
	"fmt"
	"sort"

	"strconv"

	"github.com/montanaflynn/stats"
	"github.com/piquette/finance-go/datetime"

	// "github.com/shopspring/decimal"

	"github.com/lmpizarro/go-finance/pkg/cedear"
)


func main() {

	file := "pep.csv"

	thisMap := make(map[string]float64)

	cedears := cedear.ReadCedear(file)

	start := datetime.Datetime{Month: 6, Day: 1, Year: 2021}
	end := datetime.Datetime{Month: 6, Day: 4, Year: 2021}

	for _, cear := range cedears {
		fmt.Println(cear)
		values := cedear.Historical(cear, start, end)

		for timestamp, v := range values {
			kj := strconv.FormatInt(int64(timestamp), 10)
			if val, ok := thisMap[kj]; ok {
				fmt.Println(kj, ".....", val, v)
				thisMap[kj] = val + v
			} else {
				thisMap[kj] = v
			}
		}
	}

	timestamps := make([]string, 0, len(thisMap))
	for timestamp := range thisMap {
		timestamps = append(timestamps, timestamp)
	}

	sort.Strings(timestamps)

	values := make([]float64, 0, len(thisMap))
	for _, v := range timestamps {
		values = append(values, thisMap[v])
		fmt.Printf("%10.2f \n", thisMap[v])
	}

	std, _ := stats.StandardDeviationPopulation(values)
	mean, _ := stats.Mean(values)
	std = 100 * std / mean
	initValue := thisMap[timestamps[0]]
	finalValue := thisMap[timestamps[len(timestamps)-1]]

	percent := 100 * (finalValue - initValue) / initValue

	fmt.Printf("\n      %10.2f %10.2f %10.2f %10.2f\n", initValue, finalValue, percent, std)

}

