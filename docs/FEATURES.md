# Roadmap de features por versión

Este documento separa dos cosas que MUCHAS veces se mezclan mal:

1. **lo comprometido**
2. **lo deseable**

La idea es simple: una propuesta seria no promete todo al mismo tiempo. Versiona el crecimiento.

---

## Estado actual

- Estado del repo: **proposal / pre-alpha**
- Versión publicada: **ninguna**
- Próximo objetivo: **`v0.1.0`**

---

## Convenciones del roadmap

- **Committed**: entra en esa versión salvo cambio fuerte de alcance
- **Candidate**: idea válida, pero no comprometida todavía
- **Deferred**: intencionalmente postergado

---

## `v0.1.0` — Bootstrap útil del CLI

**Objetivo:** tener una primera versión ejecutable, pequeña y defendible.

**Estado:** Committed

### Alcance

- comando `websnap shot <url>`
- validación de URL de entrada
- viewport configurable:
  - `--width`
  - `--height`
- salida configurable con `--out`
- generación automática de `media/img`
- mensajes de error claros
- contrato mínimo de salida usable desde terminal

### Lo que NO entra

- `--selector`
- `--full-page`
- `--clip`
- GIF
- video
- watch mode
- configuración por archivo

---

## `v0.2.0` — Targets de captura

**Objetivo:** pasar de “captura básica” a “captura útil para UI work”.

**Estado:** Committed

### Alcance

- `--selector`
- `--full-page`
- estrategia consistente de nombre de archivo
- protección básica ante rutas inválidas
- validación de selector y errores más descriptivos

### Riesgo principal

Capturar por selector introduce dependencia fuerte del render y del timing del DOM. Debe resolverse sin volver opaca la CLI.

---

## `v0.3.0` — Reproducibilidad y experiencia de uso

**Objetivo:** hacer la herramienta más confiable para demos, documentación y automatización.

**Estado:** Committed

### Alcance

- `--delay`
- `--timeout`
- salida más explícita en consola
- códigos de salida previsibles
- mejor ergonomía para `localhost`

### Valor

Esta versión vuelve la herramienta mucho más defendible en escenarios reales, no solo en demos felices.

---

## `v0.4.0` — Captura avanzada

**Objetivo:** agregar control más fino del área capturada.

**Estado:** Candidate

### Alcance candidato

- `--clip x,y,width,height`
- validación geométrica del área
- presets simples de viewport

### Motivo de no comprometerlo todavía

El recorte fino parece simple, pero introduce validaciones visuales y reglas de consistencia que conviene agregar después de estabilizar el flujo básico.

---

## `v0.5.0` — GIF experimental

**Objetivo:** abrir la puerta a capturas animadas sin prometer estabilidad temprana.

**Estado:** Candidate

### Alcance candidato

- comando `websnap gif <url>`
- `--duration`
- `--fps`
- captura secuencial de frames
- integración con FFmpeg

### Razón para dejarlo aquí

GIF **no es una extensión menor** del screenshot. Es otro pipeline:

- múltiple captura
- sincronización temporal
- encoding
- costo de rendimiento
- dependencia adicional

Por eso debe entrar como track separado y explícitamente experimental.

---

## `v1.0.0` — Release estable

**Objetivo:** consolidar `websnap` como herramienta confiable para screenshots desde terminal.

**Estado:** Committed como dirección, alcance exacto por refinar

### Resultado esperado

- contrato CLI estable
- flujo de screenshot sólido
- documentación madura
- errores previsibles
- distribución clara del binario

### Nota importante

`v1.0.0` no obliga a incluir GIF. La estabilidad del producto puede centrarse primero en screenshots.

---

## Backlog posterior a `v1.0.0`

**Estado:** Deferred

Ideas válidas, pero fuera del compromiso actual:

- video (`mp4` / `webm`)
- archivo `.websnap.yaml`
- modo watch
- uploads a S3 / Cloudinary
- presets por proyecto
- autenticación avanzada

---

## Regla de evolución

Si una feature agrega complejidad de runtime, dependencias extra o un pipeline distinto, **no entra por entusiasmo**.  
Primero se protege la ruta principal: `shot`.
