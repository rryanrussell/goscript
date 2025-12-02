package app

type Context struct {
	gl    WebGLRenderingContext
	state []Obj
}

type Vec2 struct {
	X float64
	Y float64
}

type Obj struct {
	Pos  Vec2
	Type int
}

func setup(c *Context) {
	c.gl = Document.GetElementsByTagName("canvas").([]HTMLCanvasElement)[0].GetContext1(Webgl, WebGLContextAttributes{})

	InitBg(c.gl)
	InitFish(c.gl)
}

func render(c *Context) {
	c.gl.ClearColor(0, 0, 0, 1)
	c.gl.Clear(c.gl.COLOR_BUFFER_BIT)

	RenderBg(c.gl)
	RenderFish(c.gl, c.state)
}

func Main() {
	var ctx Context
	setup(&ctx)

	// Fixed slice of data instead of WebSocket
	ctx.state = []Obj{
		{Pos: Vec2{X: 200, Y: 200}, Type: 1},
		{Pos: Vec2{X: 300, Y: 300}, Type: 2},
		{Pos: Vec2{X: 400, Y: 250}, Type: 1},
		{Pos: Vec2{X: 500, Y: 100}, Type: 2},
		{Pos: Vec2{X: 600, Y: 400}, Type: 1},
	}

	render(&ctx)
}
