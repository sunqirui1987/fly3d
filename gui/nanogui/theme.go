package nanogui

import (
	"github.com/suiqirui1987/fly3d/gui/canvas"
)

type Theme struct {
	StandardFontSize     int
	ButtonFontSize       int
	TextBoxFontSize      int
	WindowCornerRadius   int
	WindowHeaderHeight   int
	WindowDropShadowSize int
	ButtonCornerRadius   int

	DropShadow        canvas.Color
	Transparent       canvas.Color
	BorderDark        canvas.Color
	BorderLight       canvas.Color
	BorderMedium      canvas.Color
	TextColor         canvas.Color
	DisabledTextColor canvas.Color
	TextColorShadow   canvas.Color
	IconColor         canvas.Color

	ButtonGradientTopFocused   canvas.Color
	ButtonGradientBotFocused   canvas.Color
	ButtonGradientTopUnfocused canvas.Color
	ButtonGradientBotUnfocused canvas.Color
	ButtonGradientTopPushed    canvas.Color
	ButtonGradientBotPushed    canvas.Color

	/* Window-related */
	WindowFillUnfocused  canvas.Color
	WindowFillFocused    canvas.Color
	WindowTitleUnfocused canvas.Color
	WindowTitleFocused   canvas.Color

	WindowHeaderGradientTop canvas.Color
	WindowHeaderGradientBot canvas.Color
	WindowHeaderSepTop      canvas.Color
	WindowHeaderSepBot      canvas.Color

	WindowPopup            canvas.Color
	WindowPopupTransparent canvas.Color

	FontNormal string
	FontBold   string
	FontIcons  string
}

func NewStandardTheme(ctx *canvas.Context) *Theme {
	ctx.CreateFontFromMemory("sans", MustAsset("fonts/Roboto-Regular.ttf"), 0)
	ctx.CreateFontFromMemory("sans-bold", MustAsset("fonts/Roboto-Bold.ttf"), 0)
	ctx.CreateFontFromMemory("icons", MustAsset("fonts/entypo.ttf"), 0)
	return &Theme{
		StandardFontSize:     16,
		ButtonFontSize:       20,
		TextBoxFontSize:      20,
		WindowCornerRadius:   2,
		WindowHeaderHeight:   30,
		WindowDropShadowSize: 10,
		ButtonCornerRadius:   2,

		DropShadow:        canvas.MONO(0, 128),
		Transparent:       canvas.MONO(0, 0),
		BorderDark:        canvas.MONO(29, 255),
		BorderLight:       canvas.MONO(92, 255),
		BorderMedium:      canvas.MONO(35, 255),
		TextColor:         canvas.MONO(255, 160),
		DisabledTextColor: canvas.MONO(255, 80),
		TextColorShadow:   canvas.MONO(0, 160),
		IconColor:         canvas.MONO(255, 160),

		ButtonGradientTopFocused:   canvas.MONO(64, 255),
		ButtonGradientBotFocused:   canvas.MONO(48, 255),
		ButtonGradientTopUnfocused: canvas.MONO(74, 255),
		ButtonGradientBotUnfocused: canvas.MONO(58, 255),
		ButtonGradientTopPushed:    canvas.MONO(41, 255),
		ButtonGradientBotPushed:    canvas.MONO(29, 255),

		WindowFillUnfocused:  canvas.MONO(43, 230),
		WindowFillFocused:    canvas.MONO(45, 230),
		WindowTitleUnfocused: canvas.MONO(220, 160),
		WindowTitleFocused:   canvas.MONO(255, 190),

		WindowHeaderGradientTop: canvas.MONO(74, 255),
		WindowHeaderGradientBot: canvas.MONO(58, 255),
		WindowHeaderSepTop:      canvas.MONO(92, 255),
		WindowHeaderSepBot:      canvas.MONO(29, 255),

		WindowPopup:            canvas.MONO(50, 255),
		WindowPopupTransparent: canvas.MONO(50, 0),

		FontNormal: "sans",
		FontBold:   "sans-bold",
		FontIcons:  "icons",
	}
}
