package main

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/equity"

	"github.com/lmpizarro/go-finance/pkg/cedear"
	"github.com/lmpizarro/go-finance/pkg/utils"
)

func main() {

	lines, err := utils.ReadCsv("../cedear/cedear.csv")
	if err != nil {
		panic(err)
	}

	vmapTickEquity := get_map_equity(lines)

	map_cedear := get_map_ratios(lines, vmapTickEquity)
	total := get_total(map_cedear)

	trail, forpe := get_trailingPE(vmapTickEquity, map_cedear)

	fmt.Printf("%8s %8s %8s\n", "total", "trail", "fwdpe")
	fmt.Printf("%8.2f %8.2f %8.2f\n", total, trail, forpe)
}


func get_trailingPE(eq_map map[string]*finance.Equity, 
	                cede_map map[string]*cedear.Cedear) (float64, float64) {
	trail := 0.0
    forw := 0.0

	for k, v := range cede_map {
		eqtrail := eq_map[k].TrailingPE
		if eqtrail <=0 {
			eqtrail = eq_map[k].ForwardPE
		}
		trail += v.Percent * eq_map[k].TrailingPE / 100.0

		eqfwdPe := v.Percent * eq_map[k].ForwardPE / 100.0

		forw += eqfwdPe 
	}
	return trail, forw
}



func get_total(map_cedear map[string]*cedear.Cedear) float64 {

	total := 0.0
	for _, v := range map_cedear {
		
		total += v.Value
	}

	for _, v := range map_cedear {
		
		v.Percent = 100 * v.Value / total
	}


	return total
}


func get_map_ratios(lines [][]string, eq_map map[string]*finance.Equity) map[string]*cedear.Cedear {

	var vmap map[string]*cedear.Cedear
	vmap = make(map[string]*cedear.Cedear)

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

		quo := eq_map[line[0]]
		if err != nil {
			fmt.Printf("Error quote.Get %s", line[0])
			// panic(err)
			continue
		}

		data := cedear.Cedear{
			Ticket:   line[0],
			Ratio:    ratio,
			Quantity: quantity,
			Price:    quo.RegularMarketPrice,
		}

		data.Value = data.Price * data.Quantity / data.Ratio

		vmap[line[0]] = &data

	}
    return vmap
}

func get_map_equity(lines [][]string) map[string]*finance.Equity {

	var vmap map[string]*finance.Equity
	vmap = make(map[string]*finance.Equity)

	for _, line := range lines{

		ticket := line[0]
		q, err := equity.Get(ticket)
		if err != nil {
			// Uh-oh!
			panic(err)
		}

		vmap[ticket] = q
    
	}	

	return vmap
}
