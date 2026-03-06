package main

import (
	"fmt"
	"log"

	"market-report/analyzer"
	"market-report/report"
)

func main() {
	analysis, err := analyzer.GetGeminiAnalysis()
	if err != nil {
		log.Fatalf("ERROR FATAL: %v\n", err)
	}

	err = report.UpdateHTML(analysis)
	if err != nil {
		log.Fatalf("ERROR FATAL AL ACTUALIZAR HTML: %v\n", err)
	}

	fmt.Println("SUCCESS: HTML actualizado correctamente.")
}
