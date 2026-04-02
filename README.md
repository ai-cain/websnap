# 🚀 websnap

> Fast CLI tool to capture screenshots and GIFs from web pages or specific elements — anywhere, anytime.

---

## ✨ Descripción

**websnap** es una herramienta CLI pensada para desarrolladores que necesitan capturar rápidamente:

* 📸 Screenshots de páginas web o elementos específicos
* 🎞️ GIFs cortos de animaciones o interacciones
* 📂 Salida automática organizada (`media/img`, `media/gif`)

Todo desde cualquier ruta del sistema, sin depender de abrir herramientas externas o grabar manualmente.

---

## ⚡ Características

* CLI rápida y simple
* Funciona desde cualquier carpeta
* Captura por:

  * URL completa
  * Selector CSS (`--selector`)
  * Área específica (`--clip`)
* Generación de GIFs automáticos
* Estructura de salida automática:

  ```
  media/
    ├── img/
    └── gif/
  ```
* Soporte headless (sin UI del navegador)
* Preparado para automatización (scripts, CI, etc.)

---

## 🛠️ Stack (propuesto)

* Lenguaje: **Go** (binario ejecutable)
* Browser automation:

  * `chromedp` (simple y nativo)
  * o `playwright-go` (más avanzado)
* GIF processing: **FFmpeg**

---

## 📦 Instalación

### Opción 1: Binario (recomendado)

Descarga el ejecutable desde releases:

```bash
# Linux / Mac
chmod +x websnap
sudo mv websnap /usr/local/bin/

# Windows
# Agrega websnap.exe a tu PATH
```

---

### Opción 2: Compilar

```bash
git clone https://github.com/tu-user/websnap.git
cd websnap

go build -o websnap ./cmd/websnap
```

---

## 🚀 Uso

### Screenshot básico

```bash
websnap shot https://example.com
```

Salida:

```
./media/img/screenshot-<timestamp>.png
```

---

### Screenshot de un elemento

```bash
websnap shot https://example.com --selector "#app"
```

---

### Screenshot con área específica

```bash
websnap shot https://example.com --clip 0,0,1280,720
```

---

### Generar GIF

```bash
websnap gif https://example.com
```

Salida:

```
./media/gif/animation-<timestamp>.gif
```

---

### GIF de un elemento

```bash
websnap gif https://example.com --selector ".card"
```

---

## ⚙️ Flags disponibles

| Flag          | Descripción                    |
| ------------- | ------------------------------ |
| `--selector`  | Captura un elemento específico |
| `--clip`      | Área: `x,y,width,height`       |
| `--full-page` | Captura toda la página         |
| `--width`     | Ancho viewport                 |
| `--height`    | Alto viewport                  |
| `--delay`     | Espera antes de capturar       |
| `--duration`  | Duración del GIF               |
| `--fps`       | Frames por segundo             |
| `--out`       | Ruta de salida personalizada   |
| `--name`      | Nombre del archivo             |

---

## 📁 Estructura de salida

websnap crea automáticamente:

```
media/
  ├── img/   # screenshots
  └── gif/   # gifs
```

✔ Se crea en la **ruta donde ejecutas el comando**, no en el repo.

---

## 🧠 Cómo funciona

### Screenshot

1. Abre navegador headless
2. Carga la URL
3. Espera render
4. Captura página o elemento
5. Guarda en `media/img/`

---

### GIF

1. Abre navegador
2. Espera (`--delay`)
3. Captura múltiples frames
4. Genera GIF usando FFmpeg
5. Guarda en `media/gif/`

---

## 📌 Ejemplos reales

```bash
# Capturar landing local
websnap shot http://localhost:3000

# Capturar componente específico
websnap shot http://localhost:3000 --selector ".hero"

# Crear gif de animación
websnap gif http://localhost:3000 --duration 3s --fps 12

# Capturar área específica
websnap shot http://localhost:3000 --clip 0,0,1920,1080
```

---

## 🧩 Roadmap

* [ ] Video (mp4/webm)
* [ ] Config file (`.websnap.yaml`)
* [ ] Modo watch
* [ ] Upload automático (S3, Cloudinary)
* [ ] UI minimal opcional

---

## 🏷️ Tags

```
go
cli
screenshot
gif
web-capture
browser-automation
chromedp
playwright
ffmpeg
developer-tools
```

---

## 📄 Licencia

MIT

---

## 💡 Nombre del repo

👉 **websnap** (recomendado)

Alternativas:

* snapgif
* pagecap
* clipshot
* framegrab

