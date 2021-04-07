package main

import (
	"fmt"
    "github.com/lmpizarro/go-finance/pkg/utils"
	"github.com/gocolly/colly"


	eqt "github.com/piquette/finance-go/equity"
)

func main() {

	// Instantiate default collector
	c := colly.NewCollector(
		 
		colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

		// On every a element which has href attribute call callback
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			// Print link
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			// Visit link found on page
			// Only those links are visited which are in AllowedDomains
			c.Visit(e.Request.AbsoluteURL(link))
		})
	
		// Before making a request print "Visiting ..."
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
	
		// Start scraping on https://hackerspaces.org
		c.Visit("https://hackerspaces.org/")

}


func IterateEquity() {
	var tickets []string

	lines, err := utils.ReadCsv("nasdaq.csv")
	if err != nil {
		panic(err)
	}

	for _, line := range lines {
		tickets = append(tickets, line[0])
	}

	fmt.Println(tickets)

	iter := eqt.List(tickets)

	// Iterate over results. Will exit upon any error.
	for iter.Next() {
		q := iter.Equity()
		fmt.Printf("%T \n", q)
	}

}

func GetEquityAndPrint() {

	lines, err := utils.ReadCsv("nasdaq.csv")
	if err != nil {
		panic(err)
	}

	format0 := "%6s %6s %8s %8s %8s %6s %6s %10s \n"
	fmt.Printf(format0, "tckt", "epsFw", "epsTr12", "fwPE", "trPE", "BkV", "PrtBk", "mktPr")
	for _, line := range lines {
		// fmt.Println(line)

		eqo, err := eqt.Get(line[0])

		if err != nil {
			panic(err)
		}

		format1 := "%6s %6.2f %8.2f %8.2f %8.2f "
		format2 := "%6.2f %6.2f %10.2e\n"
		fmt.Printf(format1, line[0], eqo.EpsForward, eqo.EpsTrailingTwelveMonths, eqo.ForwardPE, eqo.TrailingPE)
		fmt.Printf(format2, eqo.BookValue, eqo.PriceToBook, float64(eqo.MarketCap))
	}
}

// https://golang.org/pkg/sort/
