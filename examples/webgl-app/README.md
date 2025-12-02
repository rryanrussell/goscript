# WebGL Application Example

This example demonstrates goscript transpiling a complete WebGL application from Go to JavaScript.

## What It Does

- Renders objects using WebGL
- Visualization of agent-based simulation (using fixed data)

## Files

- `main.go` - Main application logic
- `bg.go` - Background rendering
- `fish.go` - Object rendering (likely agents/entities)
- `dom.extern.go` - Browser API bindings (WebGL, DOM)

## Transpile

```bash
# From goscript root directory:
go run main.go examples/webgl-app > examples/webgl-app/output.js
```

## Run

1. Include `output.js` in HTML:

```html
<!DOCTYPE html>
<html>
<head>
    <title>WebGL Simulation</title>
</head>
<body>
    <canvas width="800" height="600"></canvas>
    <script src="output.js"></script>
    <script>Main();</script>
</body>
</html>
```

## External Bindings

The `.extern.go` files define TypeScript-like interfaces for browser APIs:

```go
type WebGLRenderingContext struct{}
func (gl WebGLRenderingContext) ClearColor(r, g, b, a float64)
func (gl WebGLRenderingContext) Clear(mask int)
```

These map directly to browser globals like `WebGLRenderingContext`.

## Notes

- This was originally part of a larger project for neural network visualization
- WebGL provides real-time rendering

A working example of Go driving browser graphics with minimal overhead!
