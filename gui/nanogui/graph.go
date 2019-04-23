package nanogui

import (
	"github.com/suiqirui1987/fly3d/gui/canvas"
)

type Graph struct {
	WidgetImplement

	caption, header, footer                     string
	backgroundColor, foregroundColor, textColor canvas.Color
	values                                      []float32
}

func NewGraph(parent Widget, captions ...string) *Graph {
	var caption string
	switch len(captions) {
	case 0:
		caption = "Untitled"
	case 1:
		caption = captions[0]
	default:
		panic("NewGraph can accept only one extra parameter (label)")
	}
	graph := &Graph{
		caption:         caption,
		backgroundColor: canvas.MONO(20, 128),
		foregroundColor: canvas.RGBA(255, 192, 0, 128),
		textColor:       canvas.MONO(240, 192),
	}
	InitWidget(graph, parent)
	return graph
}

func (g *Graph) Caption() string {
	return g.caption
}

func (g *Graph) SetCaption(caption string) {
	g.caption = caption
}

func (g *Graph) Header() string {
	return g.header
}

func (g *Graph) SetHeader(header string) {
	g.header = header
}

func (g *Graph) Footer() string {
	return g.footer
}

func (g *Graph) SetFooter(footer string) {
	g.footer = footer
}

func (g *Graph) BackgroundColor() canvas.Color {
	return g.backgroundColor
}

func (g *Graph) SetBackgroundColor(color canvas.Color) {
	g.backgroundColor = color
}

func (g *Graph) ForegroundColor() canvas.Color {
	return g.foregroundColor
}

func (g *Graph) SetForegroundColor(color canvas.Color) {
	g.foregroundColor = color
}

func (g *Graph) TextColor() canvas.Color {
	return g.textColor
}

func (g *Graph) SetTextColor(color canvas.Color) {
	g.textColor = color
}

func (g *Graph) Values() []float32 {
	return g.values
}

func (g *Graph) SetValues(values []float32) {
	g.values = values
}

func (g *Graph) PreferredSize(self Widget, ctx *canvas.Context) (int, int) {
	return 180, 45

}

func (g *Graph) Draw(self Widget, ctx *canvas.Context) {
	g.WidgetImplement.Draw(self, ctx)

	x := float32(g.x)
	y := float32(g.y)
	w := float32(g.w)
	h := float32(g.h)

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetFillColor(g.backgroundColor)
	ctx.Fill()

	if len(g.values) < 2 {
		return
	}

	ctx.BeginPath()
	ctx.MoveTo(x, y+h)
	dx := float32(len(g.values) - 1)
	for i, v := range g.values {
		vx := x + float32(i)*w/dx
		vy := y + (1.0-v)*h
		ctx.LineTo(vx, vy)
	}

	ctx.LineTo(x+w, y+h)
	ctx.SetStrokeColor(canvas.MONO(100, 255))
	ctx.Stroke()
	ctx.SetFillColor(g.foregroundColor)
	ctx.Fill()

	ctx.SetFontFace(g.theme.FontNormal)
	ctx.SetFillColor(g.textColor)
	if g.caption != "" {
		ctx.SetFontSize(14)
		ctx.SetTextAlign(canvas.AlignLeft | canvas.AlignTop)
		ctx.Text(x+3, y+1, g.caption)
	}

	if g.header != "" {
		ctx.SetFontSize(18)
		ctx.SetTextAlign(canvas.AlignRight | canvas.AlignTop)
		ctx.Text(x+w-3, y+1, g.header)
	}

	if g.footer != "" {
		ctx.SetFontSize(15)
		ctx.SetTextAlign(canvas.AlignRight | canvas.AlignBottom)
		ctx.Text(x+w-3, y+h-1, g.footer)
	}

	ctx.BeginPath()
	ctx.Rect(x, y, w, h)
	ctx.SetStrokeColor(canvas.MONO(100, 255))
	ctx.Stroke()
}

func (g *Graph) String() string {
	return g.StringHelper("Graph", g.caption)
}
