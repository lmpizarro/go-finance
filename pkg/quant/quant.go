package quant

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/mat"
)

/*
https://medium.com/wireless-registry-engineering/gonum-tutorial-linear-algebra-in-go-21ef136fc2d7
https://pkg.go.dev/gonum.org/v1/gonum/mat

*/

func Shift(a *mat.Dense) mat.Dense {
	rN, cN := a.Dims()

	var dst []float64
	subm := a.Slice(0, rN-1, 0, cN)
	dst = mat.Row(dst, 0, a)
	oneRow := mat.NewDense(1, cN, dst)
	shifted := mat.NewDense(rN, cN, nil)
	shifted.Stack(oneRow, subm)

	return *shifted
}

func logReturn(a, shifted *mat.Dense) mat.Dense {

	rN, cN := a.Dims()

	for i := 0; i < rN; i++ {
		for j := 0; j < cN; j++ {

			rel := a.At(i, j) / shifted.At(i, j)
			if rel <= 0 {
				rel = 0.0
			} else {
				rel = math.Log10(rel)
			}
			shifted.Set(i, j, rel)

		}
	}
	return *shifted
}

func LogReturn(a *mat.Dense) mat.Dense {
	shifted := Shift(a)

	shifted = logReturn(a, &shifted)

	return shifted
}

func MatPrint(ref string, X mat.Matrix) {

	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("---- %s -----\n", ref)
	fmt.Printf("%0.4v\n", fa)
}

func Covariance(a *mat.Dense) mat.SymDense {
	_, cN := a.Dims()

	b := mat.NewDense(cN, cN, nil)

	var cov mat.SymDense
	cov.SymOuterK(1, b)

	stat.CorrelationMatrix(&cov, a, nil)

	return cov
}

func Mean(a *mat.Dense) mat.VecDense {
	rN, cN := a.Dims()
	columnsMean := mat.NewVecDense(cN, nil)

	var dst []float64
	for i := 0; i < cN; i++ {
		columnsMean.SetVec(i, float64(rN) * stat.Mean(mat.Col(dst, i, a), nil))
	}
	return *columnsMean
}

func Returns(mean *mat.VecDense, comp *mat.VecDense) float64 {

	rN, _ := mean.Dims()

	rnC, _ := comp.Dims()

	if rN != rnC {
		panic("error")
	}

	ret := 0.0

	for i := 0; i < rN; i++ {
		ret += mean.AtVec(i) * comp.AtVec(i)
	}

	return float64(rN) * ret
}

func Variance(a *mat.Dense) mat.VecDense {
	_, cN := a.Dims()
	columnsVar := mat.NewVecDense(cN, nil)

	var dst []float64
	for i := 0; i < cN; i++ {
		columnsVar.SetVec(i, stat.Variance(mat.Col(dst, i, a), nil))
	}
	return *columnsVar
}

func Normalize(a *mat.Dense) *mat.Dense {

	rN, cN := a.Dims()

	normal := mat.NewDense(rN, cN, nil)

	mean := Mean(a)
	variance := Variance(a)

	for i := 0; i < rN; i++ {
		for j := 0; j < cN; j++ {

			normal.Set(i, j, (a.At(i,j) - mean.AtVec(j)))
			if variance.AtVec(j) != 0.0 {
			    normal.Set(i, j, (a.At(i,j) - mean.AtVec(j)) / variance.AtVec(j))
			}
		}
	}

	return normal
}