package materials

import (
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
)

type Material struct {
	Name            string
	Id              string
	_scene          *engines.Scene
	Alpha           float32
	Wireframe       bool
	BackFaceCulling bool
	_effect         IEffect
	OnDispose       func()
}

func NewMaterial(name string, scene *engines.Scene) *Material {
	this := &Material{}
	this.Name = name
	this.Id = name
	this._scene = scene

	this._scene.Materials = append(this._scene.Materials, this)

	this.Init()
	return this
}

func (this *Material) Init() {
	this.Wireframe = false
	this.Alpha = 1.0
	this.BackFaceCulling = true
}

func (this *Material) BaseDispose() {
	index := tools.IndexOf(this._scene.Materials, this)
	if index == -1 {
		return
	}
	this._scene.Materials = append(this._scene.Materials[:index], this._scene.Materials[index+1:]...)
	// Callback
	if this.OnDispose != nil {
		this.OnDispose()
	}
}

/** interface IMaterial*/
func (this *Material) IsReady(IMesh) bool {
	return true
}
func (this *Material) GetId() string {
	return this.Id
}
func (this *Material) GetEffect() IEffect {
	return this._effect
}
func (this *Material) GetRenderTargetTextures() []ITexture {
	return nil
}
func (this *Material) NeedAlphaTesting() bool {
	return false
}
func (this *Material) NeedAlphaBlending() bool {
	return this.Alpha < 1.0
}

func (this *Material) GetAlpha() float32 {
	return this.Alpha
}
func (this *Material) HasWireframe() bool {
	return this.Wireframe
}

func (this *Material) PreBind() {
	engine := this._scene.GetEngine()

	engine.EnableEffect(this._effect)
	engine.SetState(this.BackFaceCulling)
}
func (this *Material) Bind(world *math32.Matrix4, mesh IMesh) {
	return
}
func (this *Material) UnBind() {
	return
}
func (this *Material) Dispose() {
	this.BaseDispose()
}

/***/
