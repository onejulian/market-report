package main

import (
	"fmt"
	"log"

	"market-report/analyzer"
	"market-report/report"
)

func main() {
	// Paso 1: Análisis causal macro
	client, analysis, err := analyzer.GetGeminiAnalysis()
	if err != nil {
		log.Fatalf("ERROR FATAL: %v\n", err)
	}

	// Paso 2: Estimación de tendencia y catalizadores futuros
	forecast, err := analyzer.GetTrendForecast(client, analysis)
	if err != nil {
		// Si la estimación falla, se sube solo el análisis causal
		fmt.Printf("[ADVERTENCIA] Estimación de tendencia no disponible: %v\n", err)
		fmt.Println(">>> Se procederá solo con el análisis causal.")
		forecast = ""
	}

	// Paso 3: Combinar análisis + estimación y actualizar HTML
	fullReport := analysis
	if forecast != "" {
		fullReport = analysis + "\n" + forecast
	}

	err = report.UpdateHTML(fullReport)
	if err != nil {
		log.Fatalf("ERROR FATAL AL ACTUALIZAR HTML: %v\n", err)
	}

	fmt.Println("SUCCESS: HTML actualizado correctamente con análisis y proyección.")
}
