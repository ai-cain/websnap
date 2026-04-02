# websnap

> Propuesta de CLI en Go para capturas reproducibles de interfaces web desde terminal.

---

## Estado del proyecto

| Campo | Estado |
| --- | --- |
| Fase actual | Proposal / pre-alpha |
| Release publicado | Ninguno |
| Próximo objetivo | `v0.1.0` |
| Stack base elegido | **Go + chromedp** |
| Estado de GIF | Diferido a versiones posteriores |

> **Importante:** este repositorio todavía **no contiene una implementación ejecutable**.  
> Este README describe la propuesta, el alcance de la primera versión y la dirección arquitectónica.

---

## Qué problema resuelve

`websnap` busca resolver un problema concreto: **capturar interfaces web de forma reproducible, rápida y scriptable**, sin depender de abrir herramientas manuales ni repetir el mismo flujo visual una y otra vez.

Casos típicos:

- documentar componentes o pantallas
- generar evidencia visual en revisiones
- capturar vistas locales (`localhost`) o remotas
- preparar material para demos, PRs y portafolio
- automatizar capturas desde scripts o CI

---

## La pregunta importante: ¿cómo una CLI toma una captura si vive en terminal?

La terminal **no renderiza la web**. Lo que hace la CLI es **orquestar un navegador headless**.

Flujo propuesto:

1. El usuario ejecuta `websnap shot <url>`.
2. La CLI parsea argumentos y valida la entrada.
3. La CLI abre un navegador Chromium en modo headless.
4. El navegador carga la URL y renderiza la interfaz fuera de pantalla.
5. La CLI le pide al navegador capturar:
   - viewport actual
   - página completa
   - o un elemento por selector
6. La CLI guarda el PNG en disco y devuelve la ruta generada.

En resumen: **la terminal no “toma la foto”**; la CLI **dirige** a un navegador headless para que la tome.

---

## Objetivo del primer release

El primer release útil apunta a una V1 pequeña pero seria: **captura de screenshots estable**.

### Alcance comprometido para `v0.1.0`

- comando `shot`
- captura desde URL
- configuración de viewport (`--width`, `--height`)
- ruta de salida explícita con `--out`
- creación automática de `media/img`
- mensajes de error claros
- contrato CLI simple y defendible

### Fuera de alcance para `v0.1.0`

- GIF
- video
- watch mode
- uploads a servicios externos
- archivo de configuración
- autenticación compleja
- automatizaciones avanzadas tipo test runner

---

## CLI objetivo de la propuesta

> Sintaxis objetivo. **Todavía no implementada**.

```bash
websnap shot https://example.com
websnap shot https://example.com --width 1440 --height 900
websnap shot https://example.com --out ./captures/home.png
```

Capacidades previstas inmediatamente después del bootstrap inicial:

```bash
websnap shot https://example.com --selector ".hero"
websnap shot https://example.com --full-page
```

---

## Por qué Go

Se elige **Go** por razones de producto y distribución, no por moda:

- binario único y fácil de distribuir
- buena experiencia para CLI y automatización
- menor fricción en CI
- arranque rápido y mantenimiento simple
- base sólida para crecer sin cargar el proyecto con tooling innecesario

Para la primera etapa, la decisión propuesta es:

- **Lenguaje:** Go
- **Motor de navegador:** `chromedp`
- **Procesamiento de GIF a futuro:** FFmpeg, pero fuera de `v0.1.0`

---

## Estructura de salida propuesta

La herramienta creará la salida en la **ruta desde la que se ejecute el comando**, no dentro del repositorio:

```text
media/
  └── img/
```

Versiones futuras podrían sumar:

```text
media/
  ├── img/
  └── gif/
```

---

## Estado real del repositorio hoy

Hoy este repositorio contiene la **documentación base de la propuesta**:

- definición del problema
- alcance inicial
- roadmap por versiones
- arquitectura propuesta

Pendiente por construir:

- módulo Go
- comando `websnap`
- caso de uso `shot`
- adaptador `chromedp`
- manejo de salida a disco
- validación de flags y errores

---

## Documentación del proyecto

- [`docs/README.md`](docs/README.md) — índice documental
- [`docs/FEATURES.md`](docs/FEATURES.md) — roadmap por versiones
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — arquitectura propuesta en Go

---

## Roadmap resumido

- `v0.1.0` — bootstrap del CLI + captura básica
- `v0.2.0` — selector y full-page
- `v0.3.0` — mejoras de reproducibilidad y DX
- `v0.4.0` — clip y refinamientos de captura
- `v0.5.0` — GIF experimental
- `v1.0.0` — release estable

El detalle fino vive en [`docs/FEATURES.md`](docs/FEATURES.md).

---

## Principios de diseño

- **reproducibilidad sobre magia**
- **una V1 pequeña, pero sólida**
- **contrato CLI simple**
- **separar screenshots de GIF para no contaminar el diseño**
- **arquitectura extensible sin sobreingeniería**

---

## Licencia

Pendiente de definir formalmente en archivo `LICENSE`.
