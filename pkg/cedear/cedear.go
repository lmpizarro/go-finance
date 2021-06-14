package cedear

import (

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/lmpizarro/go-finance/pkg/utils"
)

// Cedear ...
type Cedear struct {
	Ticket   string
	Ratio    float64
	Price    float64
	Quantity float64
	Value    float64
	Percent  float64
}

func ReadCedear(file string) []Cedear {

	var cedears []Cedear

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

		data := Cedear{
			Ticket:   line[0],
			Ratio:    ratio,
			Quantity: quantity,
			Price:    0.0,
		}

		cedears = append(cedears, data)
	}

	return cedears
}

func Historical(quote Cedear, start, end datetime.Datetime) map[int]float64 {

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
