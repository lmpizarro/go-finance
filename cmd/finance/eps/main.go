package main

/*

Earnings per share (EPS) is calculated as a company's
profit divided by the outstanding shares of its common stock.

Profit / # shares

The price-to-earnings ratio (P/E ratio) is the ratio for valuing a
company that measures its current share price relative to its
per-share earnings (EPS).

(Price/share) / EPS = # shares * (Price/share) / Profit  =  MktCap / Profit

Price/share = PriceToBook/BookValue

pcEPS = epsFw / Price / share

*/

import (
	"fmt"

	"sort"

	"github.com/lmpizarro/go-finance/pkg/utils"
	"github.com/piquette/finance-go"
	eqt "github.com/piquette/finance-go/equity"
)

type ByMktCap []finance.Equity

func (a ByMktCap) Len() int           { return len(a) }
func (a ByMktCap) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMktCap) Less(i, j int) bool { return a[i].MarketCap < a[j].MarketCap }

func main() {

	var equities []finance.Equity
	var symbols []string

	lines, err := utils.ReadCsv("tickets.csv")
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		symbols = append(symbols, line[0])
	}

	symbols = removeDuplicateValues(symbols)

	for _, symbol := range symbols {
		q, err := eqt.Get(symbol)
		if err != nil {
			panic(err)
		}
		equities = append(equities, *q)
	}

	// ordeno para
	sort.Sort(ByMktCap(equities))
	reverse(equities)

	format0 := "%6s, %6s, %8s, %8s, %8s, %6s, %6s, %8s, %12s\n"
	fmt.Printf(format0, "tckt", "epsTr12", "epsFw", "trPE", "fwPE", "BkV", "PrtBk", "mktPr", "pcEps")

	for i, equity := range equities {

		printEquity(i, equity)
	}
}

func printEquity(i int, eqo finance.Equity) {

	// Price/share = PriceToBook/BookValue
	// pcEPS = epsFw / Price / share

	pcEPS := 100 * eqo.EpsForward / (eqo.PriceToBook * eqo.BookValue)

	format1 := "%2d %6s, %6.2f, %8.2f, %8.2f, %8.2f, "
	format2 := "%6.2f, %6.2f, %10.2e, %8.2f\n"

	if eqo.EpsTrailingTwelveMonths > 999 {
		eqo.EpsTrailingTwelveMonths = 999
	}

	if eqo.BookValue > 999 {
		eqo.BookValue = 999
	}

	fmt.Printf(format1, i, eqo.Symbol,
		eqo.EpsTrailingTwelveMonths,
		eqo.EpsForward,
		eqo.TrailingPE,
		eqo.ForwardPE)
	fmt.Printf(format2, eqo.BookValue,
		eqo.PriceToBook,
		float64(eqo.MarketCap), pcEPS)
}

func removeDuplicateValues(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func reverse(equities []finance.Equity) {
	for i, j := 0, len(equities)-1; i < j; i, j = i+1, j-1 {
		equities[i], equities[j] = equities[j], equities[i]
	}
}
