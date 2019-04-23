package materialicons

import (
	"github.com/suiqirui1987/fly3d/gui/canvas"
)

func LoadFont(ctx *canvas.Context) {
	ctx.CreateFontFromMemory("materialicons", MustAsset("font/MaterialIcons-Regular.ttf"), 0)
}

func LoadFontAs(ctx *canvas.Context, name string) {
	ctx.CreateFontFromMemory(name, MustAsset("font/MaterialIcons-Regular.ttf"), 0)
}