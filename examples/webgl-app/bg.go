package app

var squareVertices, indexBuffer WebGLBuffer
var bg WebGLProgram

func InitBg(gl WebGLRenderingContext) {
	squareVertices = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, squareVertices)
	gl.BufferData0(
		gl.ARRAY_BUFFER,
		NewFloat32Array([]float32{-1.0, 1.0, 1.0, 1.0, 1.0, -1.0, -1.0, -1.0}),
		gl.STATIC_DRAW,
	)

	indexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData0(
		gl.ELEMENT_ARRAY_BUFFER,
		NewUint16Array([]uint16{0, 1, 3, 1, 2, 3}),
		gl.STATIC_DRAW,
	)

	vSrc := `
    attribute vec2 position;
    void main() {
      gl_Position = vec4(position, 0.0, 1.0);
    }
  `
	fSrc := `
	precision mediump float;
	uniform float gray;
	uniform vec2 screen;
    void main() {
		float light = gl_FragCoord.y / screen.y;
		gl_FragColor = vec4(.3 * light, .3 * light, .3 * light + .2, 1.0);
    }
    `
	bg = CompileShader(gl, vSrc, fSrc)
}

func RenderBg(gl WebGLRenderingContext) {
	gl.UseProgram(bg)

	posAttrib := gl.GetAttribLocation(bg, "position")
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 0, 0)

	gray := gl.GetUniformLocation(bg, "gray")
	gl.Uniform1f(gray, .5)

	screen := gl.GetUniformLocation(bg, "screen")
	gl.Uniform2f(screen, 1280, 720)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_SHORT, 0)
}
