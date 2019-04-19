package lights

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type SpotLight struct {
	Light

	Angle    float32
	Exponent float32

	// Animations
}

func NewSpotLight(name string, position *math32.Vector3, direction *math32.Vector3, angle float32, exponent float32, scene *engines.Scene) *SpotLight {
	this := &SpotLight{}
	this.Init()
	this.Name = name
	this.Id = name
	this._scene = scene
	this.Position = position
	this.Direction = direction
	this.Angle = angle
	this.Exponent = exponent
	this._scene.Lights = append(this._scene.Lights, this)

	return this
}

func (this *SpotLight) Init() {
	this.Light.Init()
	this.Diffuse = math32.NewColor3(1.0, 1.0, 1.0)
	this.Specular = math32.NewColor3(1.0, 1.0, 1.0)
}
func (this *SpotLight) IsSupportShadow() bool {
	return true
}
