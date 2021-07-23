package main

import (
	"github.com/final-project/routers"
	"github.com/final-project/services/reportService"
)

func main() {
	go func() {
		reportService.RunReportService()
	}()
	routers.RunServer()
}
