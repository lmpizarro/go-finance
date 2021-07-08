package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/lmpizarro/go-finance/pkg/quant"

	"gonum.org/v1/gonum/mat"

	ced "github.com/lmpizarro/go-finance/pkg/cedear"
	"github.com/piquette/finance-go/datetime"

	// pcca "github.com/sjwhitworth/golearn/pca"

	"gonum.org/v1/gonum/stat"
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
	nP := mat.NewDense(8, 2, nil)

	for i := 0; i < 8; i++ {
		for j := 0; j < 2; j++ {
			nP.Set(i, j, float64(i*(j+1)))
		}
	}

	normal := quant.Normalize(nP)

	quant.MatPrint("normal", normal)
}

func GetPrices(fileAssetNames string, start, end datetime.Datetime) (*mat.Dense, []string) {

	cedears := ced.ReadCedear(fileAssetNames)
	mapAssetNametoTimeValues := make(map[string]map[int]float64)
	for _, cer := range cedears {
		mapTStoValue := make(map[int]float64)

		values := ced.Historical(cer, start, end)

		for timeStamp, val := range values {
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
		for ts, val := range map__ {
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
		data = append(data, 1.0/float64(numberOfAssets))
	}
	vv := mat.NewVecDense(numberOfAssets, data)

	return vv
}

func CalcVarianceByComposition(cov mat.SymDense, comp mat.VecDense) float64 {
	cN, _ := cov.Dims()

	cc := mat.NewVecDense(cN, nil)
	cc.MulVec(&cov, &comp)

	dd := mat.NewVecDense(1, nil)

	dd.MulVec(comp.T(), cc)

	return dd.At(0, 0)
}

func RandomComposition(comp mat.VecDense) {

	rN, _ := comp.Dims()

	for i := 0; i < rN; i++ {
		comp.SetVec(i, getRand())
	}

	s := mat.Sum(&comp)

	for i := 0; i < rN; i++ {
		comp.SetVec(i, comp.AtVec(i)/s)
	}

}

func MC02(numberOfAssets int, meanReturnVec *mat.VecDense, cov *mat.SymDense) *mat.VecDense {
	compositionVec := EqualComposition(numberOfAssets)
	variance := CalcVarianceByComposition(*cov, *compositionVec)
	returns := quant.Returns(meanReturnVec, compositionVec)
	einic := returns / variance
	diff_e := 1000.0

	fmt.Printf("\n return %f variance %f %f \n", returns, variance, einic)

	MM := 0
	NN := 0
	newCompositionVec := EqualComposition(numberOfAssets)
	for j := 0; j < 18000000; j++ {
		RandomComposition(*newCompositionVec)
		new_variance := CalcVarianceByComposition(*cov, *newCompositionVec)
		new_returns := quant.Returns(meanReturnVec, newCompositionVec)
		new_e := new_returns / new_variance

		if new_e > einic {
			variance = new_variance
			compositionVec.CopyVec(newCompositionVec)
			returns = new_returns
			einic = new_e
		} else {
			MM++
			e := math.Exp(-1 * new_variance / new_returns)
			r := rand.Float64()

			if r > e {
				NN++
				variance = new_variance
				compositionVec.CopyVec(newCompositionVec)
				returns = new_returns
				einic = new_e
			}
		}
		if j%1000000 == 0 {
			fmt.Printf("\n return %f variance %f %f \n", returns, variance, einic)
		}

	}

	fmt.Printf("\n return %f variance %f %f %f\n", returns, variance, einic, diff_e)

	fmt.Println(MM, NN)

	return newCompositionVec
}

func MC(numberOfAssets int, meanReturnVec *mat.VecDense, cov *mat.SymDense) *mat.VecDense {
	compositionVec := EqualComposition(numberOfAssets)
	variance := CalcVarianceByComposition(*cov, *compositionVec)
	returns := quant.Returns(meanReturnVec, compositionVec)
	einic := returns / variance
	diff_e := 1000.0

	fmt.Printf("\n return %f variance %f %f \n", returns, variance, einic)

	MM := 0
	NN := 0
	newCompositionVec := EqualComposition(numberOfAssets)
	for j := 0; j < 18000000; j++ {
		RandomComposition(*newCompositionVec)
		new_variance := CalcVarianceByComposition(*cov, *newCompositionVec)
		new_returns := quant.Returns(meanReturnVec, newCompositionVec)
		new_e := new_returns / new_variance

		if new_e > einic {
			variance = new_variance
			compositionVec.CopyVec(newCompositionVec)
			returns = new_returns
			einic = new_e
		}

		if j%1000000 == 0 {
			fmt.Printf("\n return %f variance %f %f \n", returns, variance, einic)
		}

	}

	fmt.Printf("\n return %f variance %f %f %f\n", returns, variance, einic, diff_e)

	fmt.Println(MM, NN)

	return newCompositionVec
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

//Generates random int as function of range
func getRand() float64 {
	return r.Float64()
}

func main() {
	// rand.Seed(1)
	start := datetime.Datetime{Month: 5, Day: 14, Year: 2021}
	end := datetime.Datetime{Month: 6, Day: 15, Year: 2021}

	file := "test.csv"
	prices, assetNames := GetPrices(file, start, end)
	ndata, numberOfAssets := prices.Dims()
	fmt.Printf("days %d tickets %d %s\n", ndata, numberOfAssets, assetNames)

	quant.MatPrint("prices", prices)

	logReturn := quant.LogReturn(prices)
	cov := quant.Covariance(&logReturn)
	meanReturnVec := quant.Mean(&logReturn)
	quant.MatPrint("cov", &cov)
	quant.MatPrint("meanReturns", &meanReturnVec)

	newCompositionVec := MC(numberOfAssets, &meanReturnVec, &cov)
	// newCompositionVec := mat.NewVecDense(numberOfAssets, nil)
	for i, asNa := range assetNames {
		fmt.Printf("%8s %8.2f\n", asNa, 100*newCompositionVec.AtVec(i))
	}

	var ppcc stat.PC

	rN, cN := logReturn.Dims()
	scaled := mat.NewDense(rN, cN, nil)
	scaled.Scale(100.0, &logReturn)
	ok := ppcc.PrincipalComponents(scaled, nil)
	if !ok {
		return
	}
	fmt.Printf("variances = %.4f\n\n", ppcc.VarsTo(nil))
}
