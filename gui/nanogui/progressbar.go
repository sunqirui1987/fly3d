package nanogui

import (
	"fmt"
	"github.com/suiqirui1987/fly3d/gui/canvas"
)

type ProgressBar struct {
	WidgetImplement

	value float32
}

func NewProgressBar(parent Widget) *ProgressBar {
	progressBar := &ProgressBar{}
	InitWidget(progressBar, parent)
	return progressBar
}

func (p *ProgressBar) Value() float32 {
	return p.value
}

func (p *ProgressBar) SetValue(value float32) {
	p.value = value
}

func (p *ProgressBar) PreferredSize(self Widget, ctx *canvas.Context) (int, int) {
	return 70, 12
}

func (p *ProgressBar) Draw(self Widget, ctx *canvas.Context) {
	px := float32(p.x)
	py := float32(p.y)
	pw := float32(p.w)
	ph := float32(p.h)
	p.WidgetImplement.Draw(self, ctx)
	paint := canvas.BoxGradient(px+1, py+1, pw-2, ph, 3, 4, canvas.MONO(0, 32), canvas.MONO(0, 92))
	ctx.BeginPath()
	ctx.RoundedRect(px, py, pw, ph, 3)
	ctx.SetFillPaint(paint)
	ctx.Fill()

	value := clampF(p.value, 0.0, 1.0)
	barPos := (pw - 2) * value
	barPaint := canvas.BoxGradient(px, py, barPos+1.5, ph-1, 3, 4, canvas.MONO(220, 100), canvas.MONO(128, 100))
	ctx.BeginPath()
	ctx.RoundedRect(px+1, py+1, barPos, ph-2, 3)
	ctx.SetFillPaint(barPaint)
	ctx.Fill()
}

func (p *ProgressBar) String() string {
	return p.StringHelper("ProgressBar", fmt.Sprintf("%f", p.value))
}
