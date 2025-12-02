package app

import "fmt"

func CompileShader(gl WebGLRenderingContext, vertexSrc, fragmentSrc string) WebGLProgram {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, vertexSrc)
	gl.CompileShader(vertexShader)
	ok := gl.GetShaderParameter(vertexShader, gl.COMPILE_STATUS)

	if ok == false {
		panic(fmt.Errorf("failed to compile vertex shader: %s", gl.GetShaderInfoLog(vertexShader)))
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, fragmentSrc)
	gl.CompileShader(fragmentShader)

	ok = gl.GetShaderParameter(fragmentShader, gl.COMPILE_STATUS)

	if ok == false {
		panic(fmt.Errorf("failed to compile fragment shader: %s", gl.GetShaderInfoLog(fragmentShader)))
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	gl.LinkProgram(program)

	return program
}
