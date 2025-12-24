# 💶 EUR/USD Algorithmic Causality Monitor

> **AI-Powered Macro Strategist**: Un sistema autónomo que analiza la estructura del mercado forex en tiempo real, buscando causalidad rigurosa entre flujos de bonos, datos macroeconómicos y sentimiento de riesgo.

![Python](https://img.shields.io/badge/Python-3.10%2B-blue)
![AI](https://img.shields.io/badge/AI-Gemini%203%20Flash-orange)

## 🧠 ¿Qué es esto?

Este no es un bot de trading convencional. No utiliza análisis técnico tradicional (RSI, MACD, etc.). En su lugar, actúa como un **Estratega Macro Senior**, utilizando IA Generativa con acceso a herramientas en tiempo real para explicar el **"POR QUÉ"** del movimiento del precio.

### Jerarquía de Análisis
El sistema sigue una estricta jerarquía de validación de datos:
1.  **La Verdad (Bonos):** Analiza el spread entre el *US 10Y Treasury* y el *German Bund*. El flujo de dinero real manda.
2.  **Las Expectativas (Bancos Centrales):** Revisa el precio de los futuros de fondos federales (Fed Funds Futures).
3.  **El Sentimiento (Riesgo):** Correlación con S&P 500 y VIX.
4.  **La Recencia (Noticias):** Filtra el ruido validando la fecha/hora de los datos macro (evita "stale data").

---

## 🚀 Características

-   **Análisis Autónomo:** Se ejecuta automáticamente cada 4 horas.
-   **Live Data Access:** Utiliza Google Search Tooling para obtener rendimientos de bonos y noticias financieras al instante.
-   **Dashboard Web Moderno:** Genera un reporte en `index.html` con diseño profesional, modo oscuro y responsive.
-   **Cero Alucinaciones:** Implementa controles estrictos para verificar fechas y fuentes antes de emitir un juicio.
-   **Historial Persistente:** Mantiene un archivo de los últimos 10 reportes para análisis de tendencias.

---

## 🛠️ Instalación y Uso Local

### Prerrequisitos
-   Python 3.10 o superior.
-   Una API Key de Google Gemini (Google AI Studio).

### Pasos

1.  **Clonar el repositorio:**
    ```bash
    git clone https://github.com/tu-usuario/market-report.git
    cd market-report
    ```

2.  **Instalar dependencias:**
    ```bash
    pip install -r requirements.txt
    ```

3.  **Configurar API Key:**
    *   **Mac/Linux:** `export GEMINI_API_KEY="tu_api_key_aqui"`
    *   **Windows (PowerShell):** `$env:GEMINI_API_KEY="tu_api_key_aqui"`

4.  **Ejecutar el análisis:**
    ```bash
    python main.py
    ```
    *Esto generará o actualizará el archivo `index.html` en la raíz.*

---

## 🤖 Automatización (GitHub Actions)

Este repositorio incluye un workflow configurado (`.github/workflows/monitor.yml`) para funcionar 100% en la nube de forma gratuita.

### Configuración
1.  Ve a `Settings` > `Secrets and variables` > `Actions` en tu repositorio de GitHub.
2.  Crea un **New repository secret** llamado `GEMINI_API_KEY` y pega tu clave.
3.  Habilita los permisos de escritura para el workflow en `Settings` > `Actions` > `General` > `Workflow permissions` (Seleccionar "Read and write permissions").

### Funcionamiento
-   **Frecuencia:** Se ejecuta automáticamente cada **4 horas**.
-   **Auto-Commit:** El bot analiza el mercado, actualiza el HTML y hace un `git push` automático con los nuevos cambios.
-   **GitHub Pages:** Puedes activar GitHub Pages (Source: `main` branch) para ver tu dashboard en vivo en `https://tu-usuario.github.io/market-report/`.

---

## 📂 Estructura del Proyecto

```text
market-report/
├── .github/workflows/
│   └── monitor.yml    # Configuración del cron job (cada 4h)
├── main.py            # Cerebro del sistema (Prompt Engineering + Lógica)
├── index.html         # Dashboard (Frontend generado automáticamente)
├── requirements.txt   # Dependencias (google-genai)
└── README.md          # Documentación
```

---

## ⚠️ Disclaimer

Esta herramienta es un experimento de **investigación algorítmica**. Los reportes generados son diagnósticos automatizados basados en modelos de lenguaje y datos públicos. **No constituyen asesoramiento financiero ni recomendación de inversión.** El trading de divisas conlleva un alto nivel de riesgo.

