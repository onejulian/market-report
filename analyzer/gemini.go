package analyzer

import (
	"context"
	"fmt"
	"market-report/config"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func GetGeminiAnalysis() (*genai.Client, string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		godotenv.Load()
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		return nil, "", fmt.Errorf("falta la API key, define GEMINI_API_KEY en las variables de entorno")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
		APIKey:  apiKey,
	})
	if err != nil {
		return nil, "", fmt.Errorf("error al crear el cliente de genai: %v", err)
	}

	nowStr := time.Now().UTC().Format("2006-01-02 15:04 UTC")

	// --- 2. USER TASK ---
	query := fmt.Sprintf(`
    CURRENT TIME: %s
    
    **TASK: Ejecuta el reporte de causalidad para EUR/USD.**
    
    Sigue estrictamente estos pasos de investigación usando Google Search:

    PASO 1: Validación de Datos Macro (Recencia)
    - Busca noticias "High Impact Forex News" de las últimas 6 horas.
    - Confirma fecha y hora. Si no hay, declara "CONTEXTO VACÍO".

    PASO 2: Rastreo de Flujo de Dinero (Bonos y Acciones)
    - Busca el rendimiento actual (yield) del "US 10 Year Treasury Note" y "Germany 10 Year Bund".
    - Busca "S&P 500 futures price" y "VIX index now".

    PASO 3: Análisis de Narrativa
    - Busca titulares recientes en Bloomberg/Reuters sobre EUR/USD.
    
    OUTPUT FORMAT (HTML puro):
    <div class="report-section">
      <h3>1. Estado del Conductor (Bonos)</h3>
      <p>[Tu análisis del spread aquí]</p>
      
      <h3>2. Datos Recientes y Narrativa</h3>
      <p>[Tu validación de noticias]</p>
      
      <h3>3. Conclusión Causal Rigurosa</h3>
      <p><strong>DIAGNÓSTICO:</strong> [Tu síntesis final]</p>
    </div>
    `, nowStr)

	tools := []*genai.Tool{
		{
			GoogleSearch: &genai.GoogleSearch{},
		},
	}

	reqConfig := &genai.GenerateContentConfig{
		Tools:       tools,
		Temperature: genai.Ptr[float32](0.3),
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: config.SysInstruct},
			},
		},
	}

	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				{Text: query},
			},
		},
	}

	attempts := []struct {
		model string
		label string
	}{
		{config.ModelID, "principal"},
		{config.ModelID, "principal"},
		{config.ModelIDFallback, "fallback"},
		{config.ModelIDFallback, "fallback"},
	}

	fmt.Println(">>> Iniciando contacto con Gemini (Estratega Macro)...")

	var lastErr error
	for attemptIdx, attempt := range attempts {
		callCtx, cancel := context.WithTimeout(ctx, config.CallTimeout)
		result, err := client.Models.GenerateContent(callCtx, attempt.model, contents, reqConfig)
		cancel()

		if err == nil {
			fmt.Println("\n>>> Análisis completado.")
			return client, result.Text(), nil
		}

		lastErr = err
		infoErr := err.Error()

		isTimeout := strings.Contains(infoErr, "504") || strings.Contains(infoErr, "context deadline exceeded")
		isInternalError := strings.Contains(infoErr, "500")

		reason := "Error 503 o de servicio"
		if isTimeout {
			reason = fmt.Sprintf("Timeout (%ds)", int(config.CallTimeout.Seconds()))
		}
		if isInternalError {
			reason = "Error interno"
		}

		if attemptIdx < len(attempts)-1 {
			nextModelLabel := attempts[attemptIdx+1].label
			fmt.Printf("\n[ADVERTENCIA] %s en modelo %s (Intento %d/%d)\n", reason, attempt.label, attemptIdx+1, len(attempts))
			fmt.Printf(">>> Esperando %ds antes de reintentar con modelo %s...\n", int(config.RetryDelay.Seconds()), nextModelLabel)
			time.Sleep(config.RetryDelay)
		} else {
			fmt.Printf("\n[ERROR] %s persistente en todos los modelos (%d intentos).\n", strings.ToLower(reason), len(attempts))
			fmt.Println(">>> Terminando ejecución para evitar costos o bucles infinitos.")
		}
	}

	return nil, "", fmt.Errorf("API de Gemini no disponible después de %d intentos, ultimo error: %w", len(attempts), lastErr)
}
