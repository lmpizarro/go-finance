package main

import (
	"fmt"
	"sort"

	"strconv"

	"github.com/montanaflynn/stats"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"

	// "github.com/shopspring/decimal"

	"github.com/lmpizarro/go-finance/pkg/cedear"
	"github.com/lmpizarro/go-finance/pkg/utils"
)

func ReadCedear(file string) []cedear.Cedear {

	var cedears []cedear.Cedear

	lines, err := utils.ReadCsv(file)
	if err != nil {
		panic(err)
	}

	// Loop through lines & turn into object
	for _, line := range lines {
		ratio, err := utils.StringToFloat64(line[1])
		if err != nil {
			// panic(err)
			continue
		}
		quantity, err := utils.StringToFloat64(line[2])
		if err != nil {
			// panic(err)
			continue
		}

		data := cedear.Cedear{
			Ticket:   line[0],
			Ratio:    ratio,
			Quantity: quantity,
			Price:    0.0,
		}

		cedears = append(cedears, data)
	}

	return cedears
}

func main() {

	file := "pep.csv"

	thisMap := make(map[string]float64)

	cedears := ReadCedear(file)

	start := datetime.Datetime{Month: 1, Day: 1, Year: 2020}
	end := datetime.Datetime{Month: 6, Day: 4, Year: 2021}

	for _, cedear := range cedears {
		values := historical(cedear, start, end)

		for k, v := range values {
			kj := strconv.FormatInt(int64(k), 10)
			if val, ok := thisMap[kj]; ok {
				//do something here
				thisMap[kj] = val + v
			} else {
				thisMap[kj] = v
			}
		}
	}

	timestamps := make([]string, 0, len(thisMap))
	for timestamp, _ := range thisMap {
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

func historical(quote cedear.Cedear, start, end datetime.Datetime) map[int]float64 {

	thisMap := make(map[int]float64)

	p := &chart.Params{
		Symbol:   quote.Ticket,
		Start:    &start,
		End:      &end,
		Interval: datetime.OneDay,
	}

	iter := chart.Get(p)

	var values []float64

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		b := iter.Bar()
		da := iter.Bar().Timestamp
		// avg := decimal.Avg(b.Low, b.Close, b.Open, b.High)
		val, _ := b.Close.Float64()
		val *= quote.Quantity / quote.Ratio
		values = append(values, val)
		thisMap[da] = val
		// Meta-data for the iterator - (*finance.ChartMeta).
	}

	return thisMap
}
