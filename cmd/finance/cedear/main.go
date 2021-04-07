package main

import (
	"errors"
	"fmt"

	"github.com/montanaflynn/stats"
	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/shopspring/decimal"

	"github.com/piquette/finance-go/equity"

	"log"

	"github.com/lmpizarro/go-finance/pkg/utils"

	"math"
)

// Cedear ...
type Cedear struct {
	ticket   string
	ratio    float64
	price    float64
	quantity float64
	value    float64
}

func main() {

	var cedears []Cedear
	var errormessages []string

	lines, err := utils.ReadCsv("cedear.csv")
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

		quo, err := equity.Get(line[0])
		if err != nil {
			log.Printf("Error quote.Get %s", line[0])
			// panic(err)
			continue
		}

		data := Cedear{
			ticket:   line[0],
			ratio:    ratio,
			quantity: quantity,
			price:    quo.RegularMarketPrice,
		}

		data.value = calcValue(data)

		// printCedear(data)

		pcfutearn := 100 * quo.EpsForward / quo.RegularMarketPrice

		fmt.Printf("%8s %8.2f %8.2f %8.2f", quo.Symbol, math.Logb(float64(quo.MarketCap)), quo.TrailingPE, quo.ForwardPE)
		fmt.Printf(" %8.2f %8.2f %8.2f %8.2f\n", quo.EpsForward, pcfutearn, quo.PriceToBook, quo.TrailingAnnualDividendYield)
		cedears = append(cedears, data)
	}

	fmt.Printf("\n")

	fmt.Println("retrieved data: ")
	for _, c := range cedears {
		printCedear(c)
	}

	pfv := calcPortFolioValue(cedears)

	fmt.Printf("\n portfolioValue %f \n", pfv)

	start := datetime.Datetime{Month: 10, Day: 10, Year: 2020}
	end := datetime.Datetime{Month: 3, Day: 31, Year: 2021}

	printHeader()

	for _, c := range cedears {
		err = historical(c, start, end)
		if err != nil {
			message := fmt.Sprintf("historical error %s\n", c.ticket)
			errormessages = append(errormessages, message)
			continue
		}
	}

	fmt.Println("\n Error Messages")
	for _, mess := range errormessages {
		fmt.Println(mess)
	}
}

func historical(quote Cedear, start, end datetime.Datetime) error {

	p := &chart.Params{
		Symbol:   quote.ticket,
		Start:    &start,
		End:      &end,
		Interval: datetime.OneDay,
	}

	iter := chart.Get(p)

	var values []float64

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		b := iter.Bar()
		avg := decimal.Avg(b.Low, b.Close, b.Open, b.High)
		val, _ := avg.Float64()
		values = append(values, val)
		// Meta-data for the iterator - (*finance.ChartMeta).
		// fmt.Println(iter.Meta())
	}

	std, _ := stats.StandardDeviationPopulation(values)
	min, _ := stats.Min(values)
	max, _ := stats.Max(values)
	mean, _ := stats.Mean(values)
	desviacion := 100 * (quote.price - mean) / mean
	std = 100 * std / mean

	fmt.Printf("%10s %10.2f %10.2f %10.2f %10.2f %10.2f %10.2f\n",
		quote.ticket, std, min, max, desviacion, mean, quote.price)

	// Catch an error, if there was one.
	if iter.Err() != nil {
		// Uh-oh!
		return errors.New("Statistics Error")
	}

	return nil
}

func printHeader() {
	header := [7]string{"ticket", "stdpu", "min", "max", "devi", "mean", "price"}

	fmt.Printf("\n %10s", header[0])

	for _, e := range header[1:] {
		fmt.Printf("%10s", e)
	}

	fmt.Printf("\n")
	// %10.2f %10.2f %10.2f %10.2f %10.2f \n", )
}

func printCedear(c Cedear) {
	fmt.Printf("%10s %10.2f %10.2f %10.2f %10.2f\n", c.ticket, c.ratio, c.quantity, c.price, c.value)
}

func calcPortFolioValue(c []Cedear) float64 {

	portFolioValue := 0.0

	for _, ced := range c {
		portFolioValue = portFolioValue + ced.value
	}

	return portFolioValue
}

func calcValue(c Cedear) float64 {
	return c.price * c.quantity / c.ratio
}
