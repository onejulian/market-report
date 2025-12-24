import os
import datetime
import re
from google import genai
from google.genai import types

# --- CONFIGURACIÓN ---
# Usamos el modelo experimental más capaz para razonamiento. 
# Si te da error de acceso, cambia a "gemini-1.5-pro-latest"
MODEL_ID = "gemini-3-flash-preview" 

# --- 1. SYSTEM PROMPT (La Personalidad y Reglas Rigurosas) ---
SYS_INSTRUCT = """### ROLE
Act as a Senior Macro-Quant Strategist specializing in G10 FX currencies, specifically EUR/USD. Your sole purpose is to explain the RIGOROUS CAUSALITY of current price action based on cross-asset correlation, macro data, and market structure. You do NOT predict future prices; you explain the "current state of truth."

### CORE DIRECTIVE: "RECENCY & RELEVANCE"
You must strictly validate the timestamp of every piece of data you analyze.
- Before citing a macro indicator (CPI, NFP, GDP), you must verify: Is this data from the last 24-48 hours? If it is older, it is "Stale Data" and serves only as context, not as an immediate catalyst.
- You must explicitly contrast the "Hard Data" (numbers) with the "Market Narrative" (headlines).

### ANALYSIS HIERARCHY (Order of Importance)
1. **Bond Yield Spreads (The Truth):** Analyze the real-time spread between US 10Y Treasury and German Bund (10Y). Direction of spread = Direction of flow.
2. **Central Bank Pricing (The Expectations):** What are Fed Fund Futures pricing in today vs. yesterday?
3. **Risk Sentiment (The Mood):** Correlation with S&P 500 and VIX.
4. **Calendar Validation (The Catalyst):** Check the Economic Calendar for High-Impact events released in the last 4 hours.

### RULES FOR OUTPUT
- **No Fluff:** Be concise, professional, and dense with information.
- **Divergence Spotting:** If Price is rising but Bond Spreads are falling, you MUST highlight this as a "Divergence/Anomaly" caused likely by flows/positioning, not fundamentals.
- **Missing Data:** If you cannot find real-time data for a specific metric, state "DATA UNAVAILABLE" rather than hallucinating a number.
- **Language:** Spanish.
- **Format:** Return ONLY raw HTML code for the report content (divs, h3, p, strong, ul). Do NOT use markdown blocks like ```html. Do NOT include <html> or <body> tags.
"""

def get_gemini_analysis():
    # Recomendado: exporta tu key en el entorno para no hardcodearla en el repo.
    # - Windows (PowerShell): $env:GOOGLE_API_KEY="..."
    # - Git Bash: export GOOGLE_API_KEY="..."
    api_key = os.environ.get("GEMINI_API_KEY") or os.environ.get("GOOGLE_API_KEY") or ""
    if not api_key:
        raise ValueError("Falta la API key. Define GEMINI_API_KEY o GOOGLE_API_KEY en las variables de entorno.")

    client = genai.Client(api_key=api_key)
    
    # Fecha exacta para evitar alucinaciones temporales
    now_str = datetime.datetime.now().strftime("%Y-%m-%d %H:%M UTC")

    # --- 2. USER TASK (La Instrucción de Ejecución cada 4h) ---
    query = f"""
    CURRENT TIME: {now_str}
    
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
    """

    # Configuración con Herramienta de Búsqueda
    tools = [types.Tool(google_search=types.GoogleSearch())]
    
    generate_content_config = types.GenerateContentConfig(
        tools=tools,
        system_instruction=[types.Part.from_text(text=SYS_INSTRUCT)],
        # Temperature 0 para máxima precisión, 0.3 para un poco de fluidez
        temperature=0.3 
    )

    full_response = ""
    print(">>> Iniciando contacto con Gemini (Estratega Macro)...")
    
    # Llamada al modelo
    for chunk in client.models.generate_content_stream(
        model=MODEL_ID,
        contents=[types.Content(role="user", parts=[types.Part.from_text(text=query)])],
        config=generate_content_config,
    ):
        if chunk.text:
            full_response += chunk.text
            print(".", end="", flush=True) # Feedback visual de carga
            
    print("\n>>> Análisis completado.")
    return full_response

def update_html(new_report_content):
    filename = "index.html"
    now_display = datetime.datetime.now().strftime("%d-%b-%Y %H:%M UTC")

    REPORTS_START = "<!-- MR_REPORTS_START -->"
    REPORTS_END = "<!-- MR_REPORTS_END -->"
    CARD_START = "<!-- MR_REPORT_CARD_START -->"
    CARD_END = "<!-- MR_REPORT_CARD_END -->"
    MAX_REPORTS = 10

    # --- NUEVO DISEÑO: Dashboard Financiero Dark Mode ---
    base_html = """<!DOCTYPE html>
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
</body>
</html>"""

    existing_reports_html = ""
    # Intentar recuperar historial existente
    if os.path.exists(filename):
        try:
            with open(filename, "r", encoding="utf-8") as f:
                old_content = f.read()
                s = old_content.find(REPORTS_START)
                e = old_content.find(REPORTS_END)
                if s != -1 and e != -1:
                    existing_reports_html = old_content[s + len(REPORTS_START) : e]
        except Exception:
            pass

    # Bloque del nuevo reporte
    new_block = f"""
        {CARD_START}
        <div class="report-card">
            <div class="timestamp">REPORTE GENERADO: {now_display}</div>
            <div class="content">
                {new_report_content}
            </div>
        </div>
        {CARD_END}
    """
    
    # Combinar
    card_re = re.compile(re.escape(CARD_START) + r".*?" + re.escape(CARD_END), flags=re.DOTALL)
    old_cards = card_re.findall(existing_reports_html)
    
    # Nuevo + Viejos (Max 10)
    all_cards = [new_block] + old_cards[:MAX_REPORTS-1]
    
    final_reports_section = "\n".join(c.strip() for c in all_cards)
    
    # Inyectar en la NUEVA base
    # Usamos replace con el formato exacto del string base_html de arriba
    final_html = base_html.replace(
        "            <!-- MR_REPORTS_START -->\n            <!-- MR_REPORTS_END -->",
        f"            <!-- MR_REPORTS_START -->\n{final_reports_section}\n            <!-- MR_REPORTS_END -->"
    )
    
    with open(filename, "w", encoding="utf-8") as f:
        f.write(final_html)

if __name__ == "__main__":
    try:
        report = get_gemini_analysis()
        update_html(report)
        print("SUCCESS: HTML actualizado correctamente.")
    except Exception as e:
        print(f"ERROR FATAL: {e}")