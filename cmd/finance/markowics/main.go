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

func GetPrices(fileAssetNames string, start, end datetime.Datetime) (*mat.Dense, []string) {

	cedears := ced.ReadCedear(fileAssetNames)
	mapAssetNametoTimeValues := make(map[string]map[int]float64)
	for _, cer := range cedears {
		mapTStoValue  := make(map[int]float64) 
		
		values := ced.Historical(cer, start, end)

		for timeStamp, val := range values{
			mapTStoValue[timeStamp] = val
		}

		mapAssetNametoTimeValues[cer.Ticket] = mapTStoValue 
	}

	var assetNames []string
	assetNames = make([]string, len(mapAssetNametoTimeValues))

	i := 0
	mapTStoVals := make(map[int][]float64)
	for assetName, map__ := range mapAssetNametoTimeValues {
		assetNames[i] = assetName
		i++
		for ts, val := range map__{
			mapTStoVals[ts] = append(mapTStoVals[ts], val)
		}
	}

	timestamps := make([]int, 0, len(mapTStoVals))
	for timestamp := range mapTStoVals {
		timestamps = append(timestamps, timestamp)
	}
	sort.Ints(timestamps)

	var b []float64	
	for _, tss := range timestamps {
		vals := mapTStoVals[tss]
		if len(vals) == len(mapAssetNametoTimeValues) {
			for _, val := range vals {
            	b = append(b, val)
			}
		}
	}

	bb := mat.NewDense(len(mapTStoVals), len(mapAssetNametoTimeValues), b)
	return bb, assetNames
}

func EqualComposition(numberOfAssets int) *mat.VecDense {

	var data []float64
	for i := 0; i < numberOfAssets; i++ {
		data = append(data, 1.0 / float64(numberOfAssets) )				
	} 
	vv := mat.NewVecDense(numberOfAssets, data )

	return vv
}

func CalcVarianceByComposition(cov mat.SymDense, comp mat.VecDense) float64{
	cN, _ := cov.Dims()

	cc := mat.NewVecDense(cN, nil)
	cc.MulVec(&cov, &comp)

	dd := mat.NewVecDense(1, nil)

	dd.MulVec(comp.T(), cc)
 
	return dd.At(0,0)
}

func RandomComposition(comp mat.VecDense) {

	rN, _ := comp.Dims()

	for i := 0; i < rN; i++ {
		comp.SetVec(i, rand.Float64())
	}

	s := mat.Sum(&comp)
	
	for i := 0; i < rN; i++ {
		comp.SetVec(i, comp.AtVec(i) / s)
	}

}

func main() {
	rand.Seed(1)
    start := datetime.Datetime{Month: 5, Day: 14, Year: 2020}
	end := datetime.Datetime{Month: 6, Day: 14, Year: 2021}

	file := "../csvs/pep.csv"

	prices, assetNames := GetPrices(file, start, end)

	_, numberOfAssets := prices.Dims()

	quant.MatPrint("final", prices)

	logReturn := quant.LogReturn(prices)
	cov := quant.Covariance(&logReturn)

	quant.MatPrint("cov", &cov)

	composition := EqualComposition(numberOfAssets)

	variance := CalcVarianceByComposition(cov, *composition)

	new_composition := EqualComposition(numberOfAssets)
	for j := 0; j < 1000000; j++ {
		RandomComposition(*composition)
		new_variance := CalcVarianceByComposition(cov, *composition)
		if new_variance < variance {
			variance = new_variance
			new_composition.CopyVec(composition)
		}
	}
	fmt.Printf("\n variance %.4v \n",variance)

	for i, asNa := range assetNames {
		fmt.Printf("%8s %8.2f\n", asNa, 100*new_composition.AtVec(i))		
	}

	nP := mat.NewDense(8, 2, nil)
	
	for i := 0; i < 8; i++ {
		for j := 0;  j < 2; j++  {
			nP.Set(i, j, float64(i*(j+1)))
		}
	}
		

	normal := quant.Normalize(nP)

	quant.MatPrint("normal", normal)
}
