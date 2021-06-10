package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type EtfInCedear struct {
	Etf string
	Weight float64
}

var map_etf_cedear map[string]EtfInCedear


func search_etf(etf, file__ string) map[string]EtfInCedear{

    map_etf_cedear := make(map[string]EtfInCedear)
	map_ticket_weight := make(map[string]float64)

	file, err := os.Open(file__)
	if err != nil {
		log.Fatalf("failed to open")

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	for _, each_ln := range text {
		fields := strings.Fields(each_ln)
		weight, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			panic(err)
		}
		map_ticket_weight[fields[0]] = weight
	}

	file1, err1 := os.Open("todo_cedear.txt")
	if err1 != nil {
		log.Fatalf("failed to open")

	}
	defer file1.Close()

	scanner2 := bufio.NewScanner(file1)
	scanner2.Split(bufio.ScanLines)
	var cedears_ticket []string

	for scanner2.Scan() {
		cedears_ticket = append(cedears_ticket, scanner2.Text())
	}

	for _, each_ln := range cedears_ticket {
		// fields := strings.Fields(each_ln)
		symbol := strings.TrimSpace(each_ln)

		v, ok := map_ticket_weight[symbol]
		if ok {
			fmt.Println(symbol, v, etf)

			_, has := map_etf_cedear[symbol]
			if !has {
				map_etf_cedear[symbol] = EtfInCedear{Etf: etf, Weight: v}
			}
		}
	}

	return map_etf_cedear

}

func main() {
	map_etf_file := make(map[string]string)
    map_etf_cedear := make(map[string]EtfInCedear)

	map_etf_file["arkf"] = "arkf.txt"
	map_etf_file["arkg"] = "arkg.txt"
	map_etf_file["arkk"] = "arkk.txt"
	map_etf_file["arkq"] = "arkq.txt"
	map_etf_file["arkw"] = "arkw.txt"
	map_etf_file["arkx"] = "arkx.txt"

	for key, value := range map_etf_file {
		map_e := search_etf(key, value)

		for k, v := range map_e {
			vv, has := map_etf_cedear[k]
			if !has{
				map_etf_cedear[k] = v
			}else{
				if vv.Weight > v.Weight{
					map_etf_cedear[k] = vv
				}else{
					map_etf_cedear[k] = v
				}
			}
		}
	}


	fmt.Printf("\n\n\n")
	for kk, vv := range map_etf_cedear{
    	fmt.Println(kk, vv.Weight, vv.Etf)

	}
}