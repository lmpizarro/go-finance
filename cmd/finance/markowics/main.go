package main

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/lmpizarro/go-finance/pkg/quant"
	"gonum.org/v1/gonum/mat"

	ced "github.com/lmpizarro/go-finance/pkg/cedear"
	"github.com/piquette/finance-go/datetime"
)

func test() {
	fmt.Println("hello")

	// Generate a 6×6 matrix of random values.
	a := mat.NewDense(30, 2, nil)

	for i := 0; i < 30; i++ {
		a.Set(i, 0, float64(i)+rand.NormFloat64()*.1)
		a.Set(i, 1, 30.0-float64(i)+rand.NormFloat64()*.1)

	}
	quant.MatPrint("data", a)

	logReturn := quant.LogReturn(a)
	cov := quant.Covariance(&logReturn)

	quant.MatPrint("cov", &cov)

	columnsMean := quant.Mean(a)

	quant.MatPrint("mean", &columnsMean)

}

func main() {
    start := datetime.Datetime{Month: 1, Day: 1, Year: 2020}
	end := datetime.Datetime{Month: 6, Day: 12, Year: 2021}

	file := "pep.csv"

	thisMap := make(map[string]map[int]float64)

	cedears := ced.ReadCedear(file)

	for _, cer := range cedears {
		a  := make(map[int]float64) 
		
		values := ced.Historical(cer, start, end)

		for j, val := range values{

			fmt.Println(j, val)
			a[j] = val
		}

		thisMap[cer.Ticket] = a 

	}

	var tt []string
	tt = make([]string, len(thisMap))

	mapTStoVals := make(map[int][]float64)
	i := 0
	for cername, map__ := range thisMap {
		tt[i] = cername
		i++
		for ts, val := range map__{
			mapTStoVals[ts] = append(mapTStoVals[ts], val)
		}
	}

	timestamps := make([]int, 0, len(mapTStoVals))
	for timestamp, _ := range mapTStoVals {
		timestamps = append(timestamps, timestamp)
	}
	sort.Ints(timestamps)

	var b []float64	
	for _, tss := range timestamps {
		vals := mapTStoVals[tss]
		if len(vals) == len(thisMap) {
			for _, val := range vals {
            	b = append(b, val)
			}
		}
	}

	bb := mat.NewDense(len(mapTStoVals), len(thisMap), b)
	quant.MatPrint("final", bb)

	logReturn := quant.LogReturn(bb)
	cov := quant.Covariance(&logReturn)

	fmt.Println(tt)
	quant.MatPrint("cov", &cov)

	var data []float64

	for i := 0; i < len(thisMap); i++ {
		data = append(data, 1.0 / float64(len(thisMap)) )				
	} 
	vv := mat.NewVecDense(len(thisMap), data )

	cc := mat.NewVecDense(len(thisMap), nil)
	cc.MulVec(&cov, vv)

	dd := mat.NewVecDense(1, nil)

	dd.MulVec(vv.T(), cc)

	quant.MatPrint("var", dd)
}
