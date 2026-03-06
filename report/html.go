package report

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

func UpdateHTML(newReportContent string) error {
	filename := "index.html"
	nowDisplay := time.Now().UTC().Format("02-Jan-2006 15:04 UTC")

	reportsStart := "<!-- MR_REPORTS_START -->"
	reportsEnd := "<!-- MR_REPORTS_END -->"
	cardStart := "<!-- MR_REPORT_CARD_START -->"
	cardEnd := "<!-- MR_REPORT_CARD_END -->"
	maxReports := 10

	baseHTML := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#0f172a">
    <title>EUR/USD Macro Monitor</title>
    <style>
        :root {
            --bg-body: #0f172a;       /* Slate 900 */
            --bg-card: #1e293b;       /* Slate 800 */
            --text-main: #f1f5f9;     /* Slate 100 */
            --text-muted: #94a3b8;    /* Slate 400 */
            --border: #334155;        /* Slate 700 */
            --accent: #38bdf8;        /* Sky 400 */
            --accent-glow: rgba(56, 189, 248, 0.15);
            --success: #34d399;       /* Emerald 400 */
            --font-main: 'Inter', 'Segoe UI', system-ui, -apple-system, sans-serif;
        }

        body {
            font-family: var(--font-main);
            background-color: var(--bg-body);
            color: var(--text-main);
            margin: 0;
            padding: 20px;
            line-height: 1.6;
            -webkit-font-smoothing: antialiased;
        }

        .container {
            max-width: 850px;
            margin: 0 auto;
        }

        .header {
            text-align: center;
            margin-bottom: 40px;
            padding-bottom: 25px;
            border-bottom: 1px solid var(--border);
            animation: fadeIn 0.8s ease-out;
        }

        .header h1 {
            margin: 0 0 10px 0;
            font-size: 2rem;
            font-weight: 700;
            letter-spacing: -0.02em;
            background: linear-gradient(to right, #fff, #cbd5e1);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        .header p {
            margin: 0;
            color: var(--text-muted);
            font-size: 0.95rem;
        }

        /* Tarjetas de Reporte */
        .report-card {
            background-color: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 24px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.2);
            position: relative;
            overflow: hidden;
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .report-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
            border-color: var(--accent);
        }

        /* Borde lateral de acento */
        .report-card::before {
            content: "";
            position: absolute;
            left: 0; top: 0; bottom: 0;
            width: 4px;
            background: var(--accent);
            opacity: 0.8;
        }

        .timestamp {
            display: inline-block;
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            color: var(--accent);
            background: var(--accent-glow);
            padding: 4px 12px;
            border-radius: 99px;
            margin-bottom: 20px;
            border: 1px solid rgba(56, 189, 248, 0.2);
        }

        /* Estilos del contenido generado */
        .content h3 {
            color: var(--text-main);
            font-size: 1.1rem;
            margin-top: 24px;
            margin-bottom: 12px;
            padding-bottom: 8px;
            border-bottom: 1px solid var(--border);
            display: flex;
            align-items: center;
        }
        
        .content h3::before {
            content: "▹";
            margin-right: 8px;
            color: var(--accent);
        }

        .content p {
            color: var(--text-muted);
            margin-bottom: 16px;
        }

        .content strong {
            color: var(--success); /* Resalta datos clave en verde */
            font-weight: 600;
        }
        
        /* Diagnóstico final destacado */
        .content p:last-child strong {
            color: #fbbf24; /* Amber para la conclusión */
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Mobile */
        @media (max-width: 600px) {
            body { padding: 12px; }
            .header h1 { font-size: 1.5rem; }
            .report-card { padding: 20px; border-radius: 12px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>EUR/USD Algorithmic Causality</h1>
            <p>Monitor de Estructura de Mercado & Flujos en Tiempo Real</p>
        </div>
        
        <div id="archive">
            <!-- MR_REPORTS_START -->
            <!-- MR_REPORTS_END -->
        </div>
    </div>
    <script>
        function updateTimeAgo() {
            const months = {
                'JAN': '01', 'FEB': '02', 'MAR': '03', 'APR': '04', 'MAY': '05', 'JUN': '06',
                'JUL': '07', 'AUG': '08', 'SEP': '09', 'OCT': '10', 'NOV': '11', 'DEC': '12'
            };

            function parseDateString(dateStr) {
                let parts = dateStr.toUpperCase().match(/(\d{2})-([A-Z]{3})-(\d{4})\s+(\d{2}):(\d{2})\s+UTC/);
                if (parts) {
                    const isoStr = parts[3] + "-" + months[parts[2]] + "-" + parts[1] + "T" + parts[4] + ":" + parts[5] + ":00Z";
                    return new Date(isoStr);
                }
                return null;
            }

            function timeAgo(date) {
                if (!date || isNaN(date.getTime())) return "";

                const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
                if (seconds < 60) return "hace menos de 1 minuto";

                const intervals = [
                    { label: 'año', seconds: 31536000, plural: 's' },
                    { label: 'mes', seconds: 2592000, plural: 'es' },
                    { label: 'día', seconds: 86400, plural: 's' },
                    { label: 'hora', seconds: 3600, plural: 's' },
                    { label: 'minuto', seconds: 60, plural: 's' }
                ];

                for (let i = 0; i < intervals.length; i++) {
                    const interval = intervals[i];
                    const count = Math.floor(seconds / interval.seconds);
                    if (count >= 1) {
                        return "hace " + count + " " + interval.label + (count === 1 ? "" : interval.plural);
                    }
                }
                return "hace menos de 1 minuto";
            }

            document.querySelectorAll('.timestamp').forEach(function (ts) {
                if (!ts.dataset.original) {
                    ts.dataset.original = ts.textContent.trim();
                }

                const text = ts.dataset.original;
                const dateMatch = text.match(/REPORTE GENERADO:\s+(.*)/i);

                if (dateMatch && dateMatch[1]) {
                    const date = parseDateString(dateMatch[1]);
                    if (date && !isNaN(date.getTime())) {
                        const ago = timeAgo(date);
                        ts.textContent = text + " (" + ago + ")";
                    }
                }
            });
        }

        document.addEventListener('DOMContentLoaded', function () {
            updateTimeAgo();
            setInterval(updateTimeAgo, 60000);
        });
    </script>
</body>
</html>`

	var existingReportsHTML string

	// Intentar recuperar historial existente
	if _, err := os.Stat(filename); err == nil {
		contentBytes, err := os.ReadFile(filename)
		if err == nil {
			oldContent := string(contentBytes)
			s := strings.Index(oldContent, reportsStart)
			e := strings.Index(oldContent, reportsEnd)
			if s != -1 && e != -1 && s+len(reportsStart) < e {
				existingReportsHTML = oldContent[s+len(reportsStart) : e]
			}
		}
	}

	newBlock := fmt.Sprintf(`
        %s
        <div class="report-card">
            <div class="timestamp">REPORTE GENERADO: %s</div>
            <div class="content">
                %s
            </div>
        </div>
        %s
    `, cardStart, nowDisplay, newReportContent, cardEnd)

	// Extraer tarjetas antiguas
	pattern := "(?s)" + regexp.QuoteMeta(cardStart) + ".*?" + regexp.QuoteMeta(cardEnd)
	cardRe := regexp.MustCompile(pattern)
	oldCards := cardRe.FindAllString(existingReportsHTML, -1)

	// Limitar a maxReports - 1
	limit := maxReports - 1
	if len(oldCards) > limit {
		oldCards = oldCards[:limit]
	}

	// Combinar tarjetas
	var allCards []string
	allCards = append(allCards, newBlock)
	allCards = append(allCards, oldCards...)

	// Limpiar espacios extra y unir
	for i := range allCards {
		allCards[i] = strings.TrimSpace(allCards[i])
	}
	finalReportsSection := strings.Join(allCards, "\n")

	// Inyectar en la NUEVA base
	targetPattern := fmt.Sprintf("            %s\n            %s", reportsStart, reportsEnd)
	replacement := fmt.Sprintf("            %s\n%s\n            %s", reportsStart, finalReportsSection, reportsEnd)
	finalHTML := strings.Replace(baseHTML, targetPattern, replacement, 1)

	err := os.WriteFile(filename, []byte(finalHTML), 0644)
	if err != nil {
		return fmt.Errorf("error al escribir el archivo html: %v", err)
	}

	return nil
}
