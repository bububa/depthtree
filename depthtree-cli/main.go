package main

import (
	"flag"
	"fmt"
	"github.com/bububa/depthtree/client"
	//"github.com/mkideal/log"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		host string
		csv  string
	)
	flag.StringVar(&host, "host", "", "host")
	flag.StringVar(&csv, "csv", "", "import csv file")
	flag.Parse()
	c := client.NewClient(fmt.Sprintf("http://%s", host))
	if csv != "" {
		ImportCsv(c, csv)
	}
}
