package main

import (
	"fmt"
	"math"
	"sync"

	"sort"

	// "strconv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"

	"github.com/piquette/finance-go/datetime"

	"github.com/lmpizarro/go-finance/pkg/cedear"
)

type tprice struct {
	name  string
	value float64
	index int
}

func main() {

	var wg sync.WaitGroup

	file := "../simulator/pep.csv"

	cedears := cedear.ReadCedear(file)
	cedear_index := make(map[string]int)
	thisMap := make(map[int][]tprice)

	start := datetime.Datetime{Month: 6, Day: 16, Year: 2021}
	end := datetime.Datetime{Month: 7, Day: 7, Year: 2021}

	wg.Add(len(cedears))
	for index, cear := range cedears {
		cedear_index[cear.Ticket] = index

		go func(cear_ cedear.Cedear, index_ int) {
			defer wg.Done()

			values := cedear.Historical(cear_, start, end)
			
			for timestamp, value := range values {
				val_t := tprice{name: cear_.Ticket, value: value, index: index_}
				thisMap[timestamp] = append(thisMap[timestamp], val_t)
			}
		}(cear, index)
	}
	wg.Wait()

	timestamps__ := make([]int, 0, len(thisMap))
	for timestamp := range thisMap {
		timestamps__ = append(timestamps__, timestamp)
	}

	sort.Ints(timestamps__)

	

	matr_val := make([][]float64, 0)
	for _, ts := range timestamps__ {
		array_val := make([]float64, len(cedears)+1)
		array_val[0] = float64(ts)
		for _, val := range thisMap[ts] {
			array_val[val.index + 1] = val.value
		}
		matr_val = append(matr_val, array_val)
	}

	fmt.Println(cedear_index)

	var buffer_vals_matrix = make([]float64, 0)
	for _, v := range matr_val {
		for j, vv := range v {
			if j != 0 {
				buffer_vals_matrix = append(buffer_vals_matrix, vv)
			}
		}
	}

	matrix_num := mat.NewDense(len(matr_val), len(matr_val[0])-1, buffer_vals_matrix)

	nr, nc := matrix_num.Dims()   
	delayed := matrix_num.Slice(0, nr-1, 0, nc).(*mat.Dense)
	matrix_num = matrix_num.Slice(1, nr, 0, nc).(*mat.Dense)


	ratio_price := mat.NewDense(nr-1, nc, nil)
	ratio_price.DivElem(matrix_num, delayed)

	log_return := mat.NewDense(nr-1, nc, nil)
	log_return.Apply(log_elem, ratio_price)
   
	print_matrix(log_return)

	col := mat.Col(nil, 0, log_return)
	vec_col := mat.NewVecDense(len(col), col)
	mean := stat.Mean(col, nil)
	varc := stat.Variance(col, nil)
	min := mat.Min(vec_col)
	max := mat.Max(vec_col)
	sigma := math.Sqrt(varc)

	fmt.Printf("%e %e %e %e %e %e %d \n", 
	       mean, sigma, min, max, mean - 2 * sigma, mean + 2 * sigma, count_pos(*vec_col))
}

func count_pos(vec_col mat.VecDense) int {

	count_pos := 0

	nr, _ := vec_col.Dims()

	
    for i := 0; i < nr; i++ {
	    if vec_col.At(i, 0) > 0 {
			count_pos++
		}	
	}
	return count_pos
}

func log_elem(i, j int, v float64) float64 {
	return math.Log(v)
}

func print_matrix(a mat.Matrix) {
	r, c := a.Dims()

	fmt.Println()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%f ", a.At(i, j))
		}

		fmt.Println()
	}
	fmt.Println()
}