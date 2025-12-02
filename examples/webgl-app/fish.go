package app

var fish WebGLProgram

func InitFish(gl WebGLRenderingContext) {
	vSrc := `
    attribute vec2 position;
	uniform vec2 center;
	uniform vec2 scale;
    void main() {
      gl_Position = vec4(scale * position + center, 0.0, 1.0);
    }
  `
	fSrc := `
	precision mediump float;
	uniform float gray;
	uniform vec2 screen;
    void main() {
		float light = gl_FragCoord.y / screen.y;
		gl_FragColor = vec4(.5, .5, .5, 1.0);
    }
    `
	fish = CompileShader(gl, vSrc, fSrc)
}

func RenderFish(gl WebGLRenderingContext, fishes []Obj) {
	if len(fishes) == 0 {
		return
	}

	gl.UseProgram(fish)

	posAttrib := gl.GetAttribLocation(fish, "position")
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 0, 0)

	gray := gl.GetUniformLocation(fish, "gray")
	gl.Uniform1f(gray, .5)

	screen := gl.GetUniformLocation(fish, "screen")
	gl.Uniform2f(screen, 1280, 720)

	scale := gl.GetUniformLocation(fish, "scale")

	center := gl.GetUniformLocation(fish, "center")

	for _, f := range fishes {
		size := 20.0

		if f.Type == 2 {
			size = 5.0
		}

		gl.Uniform2f(scale, size/1280.0, size/720.0)
		x := f.Pos.X/1280.0*2.0 - 1.0
		y := 1.0 - f.Pos.Y/720.0*2.0
		gl.Uniform2f(center, x, y)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_SHORT, 0)
	}
}
