package analyzer

import (
	"context"
	"fmt"
	"market-report/config"
	"strings"
	"time"

	"google.golang.org/genai"
)

// GetTrendForecast toma el análisis causal previo como contexto y genera
// una estimación de duración de la tendencia actual + próximos eventos clave.
func GetTrendForecast(client *genai.Client, previousAnalysis string) (string, error) {
	ctx := context.Background()
	nowStr := time.Now().UTC().Format("2006-01-02 15:04 UTC")

	query := fmt.Sprintf(`
    CURRENT TIME: %s

    **CONTEXT — ANÁLISIS CAUSAL PREVIO (ya publicado en el reporte):**
    %s

    ---

    **TASK: Estimación de Persistencia de Tendencia y Próximos Catalizadores para EUR/USD.**

    Basándote en el análisis causal previo que acabas de recibir como contexto, y realizando búsquedas
    adicionales en Google para validar y complementar tu estimación, genera el siguiente reporte:

    SECCIÓN A — ESTIMACIÓN DE PERSISTENCIA DE LA TENDENCIA ACTUAL:
    1. Identifica la tendencia o fenómeno dominante que está moviendo el precio AHORA (según el análisis previo).
    2. Estima con la mayor precisión posible el rango de tiempo en el que esta tendencia podría mantenerse.
       - Usa lógica causal: ¿el driver es un dato puntual (efecto 2-6 horas), un cambio de política monetaria (efecto semanas/meses), o un flujo de posicionamiento (efecto días)?
       - Establece un escenario BASE (más probable) y un escenario ALTERNATIVO (si algo cambia antes).
    3. Indica qué condiciones específicas invalidarían la tendencia actual (triggers de reversión).

    SECCIÓN B — PRÓXIMOS EVENTOS CATALIZADORES:
    1. Busca "forex economic calendar this week" y "forex economic calendar next week" en Google.
    2. Lista los próximos 5-8 eventos de ALTO IMPACTO que podrían POTENCIAR o REVERTIR la tendencia actual.
    3. Para CADA evento indica:
       - Nombre exacto del evento
       - Fecha y hora exacta (UTC)
       - Impacto esperado: ¿potencia la tendencia actual o la revierte? ¿Por qué?
       - Nivel de importancia: CRÍTICO / ALTO / MODERADO
    4. Ordena los eventos cronológicamente.

    SECCIÓN C — VEREDICTO DE PROYECCIÓN:
    - Resume en 2-3 oraciones tu proyección integral: duración estimada de la tendencia,
      primer evento que podría alterar el curso, y nivel de confianza de tu estimación.

    OUTPUT FORMAT (HTML puro, sin markdown, sin tags html/body):
    <div class="forecast-section">
      <h3>🔮 Persistencia de la Tendencia Actual</h3>
      <p><strong>Tendencia dominante:</strong> [descripción]</p>
      <p><strong>Duración estimada:</strong> [rango temporal]</p>
      <p><strong>Escenario base:</strong> [descripción]</p>
      <p><strong>Escenario alternativo:</strong> [descripción]</p>
      <p><strong>Triggers de invalidación:</strong> [condiciones que romperían la tendencia]</p>

      <h3>📅 Próximos Catalizadores Clave</h3>
      <table class="events-table">
        <thead>
          <tr>
            <th>Fecha/Hora (UTC)</th>
            <th>Evento</th>
            <th>Impacto en Tendencia</th>
            <th>Importancia</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>[fecha]</td>
            <td>[nombre]</td>
            <td>[potencia/revierte + razón]</td>
            <td><span class="badge-[critico|alto|moderado]">[CRÍTICO|ALTO|MODERADO]</span></td>
          </tr>
          <!-- más filas -->
        </tbody>
      </table>

      <h3>⚡ Veredicto de Proyección</h3>
      <p><strong>PROYECCIÓN:</strong> [tu síntesis]</p>
    </div>
    `, nowStr, previousAnalysis)

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
				{Text: config.ForecastSysInstruct},
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

	fmt.Println("\n>>> Iniciando estimación de tendencia y catalizadores futuros...")

	var lastErr error
	for attemptIdx, attempt := range attempts {
		callCtx, cancel := context.WithTimeout(ctx, config.CallTimeout)
		result, err := client.Models.GenerateContent(callCtx, attempt.model, contents, reqConfig)
		cancel()

		if err == nil {
			fmt.Println("\n>>> Estimación de tendencia completada.")
			return result.Text(), nil
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
			fmt.Println(">>> Terminando estimación de tendencia.")
		}
	}

	return "", fmt.Errorf("estimación de tendencia no disponible después de %d intentos, ultimo error: %w", len(attempts), lastErr)
}
