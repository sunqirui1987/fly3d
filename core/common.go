package core

import (
	"image"
	"image/jpeg"
	"image/png"

	"github.com/suiqirui1987/fly3d/math32"
	log "github.com/suiqirui1987/fly3d/tools/logrus"
)

const (
	ALPHA_DISABLE = 0
	ALPHA_ADD     = 1
	ALPHA_COMBINE = 2

	Epsilon           = 0.001
	CollisionsEpsilon = 0.001

	// Statics
	FOGMODE_NONE   = 0
	FOGMODE_EXP    = 1
	FOGMODE_EXP2   = 2
	FOGMODE_LINEAR = 3
)

var ()

type Fly3D struct {
	ClipPlane         *math32.Plane
	IsIE              bool //IE
	ResRepository     string
	ShadersRepository string
}

//全局
var GlobalFly3D = &Fly3D{
	ClipPlane:         nil,
	IsIE:              false,
	ShadersRepository: "github.com/suiqirui1987/fly3d/shaders/",
}

func (this *Fly3D) SetIsDebug(val bool) {
	log.IsDebug = val
}

func init() {

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

}
