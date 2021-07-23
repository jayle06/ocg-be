package main

import (
	"fmt"
	"github.com/final-project/routers"
	"github.com/final-project/services/reportService"
)

func main() {
	fmt.Println("Hello Nam, cày project đi nhé !!!!")
	go func() {
		reportService.RunReportService()
	}()
	routers.RunServer()
}
