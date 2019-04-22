package lights

import (
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
)

type Light struct {
	Name   string
	Id     string
	_scene *engines.Scene

	ShadowGenerator IShadowGenerator

	Position  *math32.Vector3
	Direction *math32.Vector3

	Intensity float32
	Isenable  bool
	Diffuse   *math32.Color3
	Specular  *math32.Color3
}

func NewLight(name string, scene *engines.Scene) *Light {

	this := &Light{}
	this.Name = name
	this.Id = name

	this._scene = scene
	this._scene.Lights = append(this._scene.Lights, this)

	return this
}

func (this *Light) Init() {
	this.Intensity = 1.0
	this.Isenable = true
}

func (this *Light) GetScene() *engines.Scene {
	return this._scene
}

func (this *Light) IsEnabled() bool {
	return this.Isenable
}
func (this *Light) GetIntensity() float32 {
	return this.Intensity
}

func (this *Light) GetDiffuse() *math32.Color3 {
	return this.Diffuse
}
func (this *Light) GetSpecular() *math32.Color3 {
	return this.Specular
}

func (this *Light) GetPosition() *math32.Vector3 {
	return this.Position
}

func (this *Light) GetDirection() *math32.Vector3 {
	return this.Direction
}

func (this *Light) GetShadowGenerator() IShadowGenerator {
	return this.ShadowGenerator
}
func (this *Light) SetShadowGenerator(val IShadowGenerator) {

	this.ShadowGenerator = val
}

func (this *Light) IsSupportShadow() bool {
	return false
}

// Methods
func (this *Light) Dispose() {
	if this.ShadowGenerator != nil {
		this.ShadowGenerator.Dispose()
		this.ShadowGenerator = nil
	}

	// Remove from scene
	index := tools.IndexOf(this, this._scene.Lights)
	if index > -1 {
		this._scene.Lights = append(this._scene.Lights[:index], this._scene.Lights[index+1:]...)
	}

}
