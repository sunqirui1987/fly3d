package lights

import (
	"math"

	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/effects"
	"github.com/suiqirui1987/fly3d/module/textures"
	log "github.com/suiqirui1987/fly3d/tools/logrus"
)

type ShadowGenerator struct {
	UseVarianceShadowMap bool
	_light               ILight
	_scene               *engines.Scene
	_shadowMap           *textures.RenderTargetTexture

	_effect    IEffect
	_effectVSM IEffect

	_viewMatrix          *math32.Matrix4
	_projectionMatrix    *math32.Matrix4
	_transformMatrix     *math32.Matrix4
	_worldViewProjection *math32.Matrix4

	_cachedPosition  *math32.Vector3
	_cachedDirection *math32.Vector3
}

func NewShadowGenerator(name string, mapSize int, light ILight, scene *engines.Scene) *ShadowGenerator {
	this := &ShadowGenerator{}

	if light.IsSupportShadow() == false {
		log.Printf(" this light not support shadow")
		return nil
	}

	this._light = light
	this._scene = scene
	this._light.SetShadowGenerator(this)

	engine := this._scene.GetEngine()

	// Render target
	this._shadowMap = textures.NewRenderTargetTexture(name, mapSize, this._scene, false)
	this._shadowMap.GetGLTexture().WrapU = gl.CLAMP_ADDRESSMODE
	this._shadowMap.GetGLTexture().WrapV = gl.CLAMP_ADDRESSMODE

	// Effect
	this._effect = effects.CreateEffect(engine, "shadowMap",
		[]string{"position"},
		[]string{"worldViewProjection"},
		[]string{}, "")

	this._effectVSM = effects.CreateEffect(engine, "shadowMap",
		[]string{"position"},
		[]string{"worldViewProjection"},
		[]string{}, "#define VSM")

	// Custom render function
	that := this

	renderSubMesh := func(subMesh ISubMesh, effect IEffect) {

		mesh := subMesh.GetMesh()

		world := mesh.GetWorldMatrix()

		that._worldViewProjection = world.Multiply(that.GetTransformMatrix())

		effect.SetMatrix("worldViewProjection", that._worldViewProjection)

		subMesh.BindAndDraw(effect, false)

	}

	this._shadowMap.CustomRenderFunction = func(opaqueSubMeshes []ISubMesh, alphaTestSubMeshes []ISubMesh, transparentSubMeshes []ISubMesh, activeMeshes []IMesh) {
		engine := that._scene.GetEngine()

		var effect IEffect
		if that.UseVarianceShadowMap == true {
			effect = that._effectVSM
		} else {
			effect = that._effect
		}

		engine.EnableEffect(effect)

		for index := 0; index < len(opaqueSubMeshes); index++ {
			renderSubMesh(opaqueSubMeshes[index], effect)
		}

		for index := 0; index < len(alphaTestSubMeshes); index++ {
			renderSubMesh(alphaTestSubMeshes[index], effect)
		}
	}

	// Internals
	this._viewMatrix = math32.NewMatrix4().Zero()
	this._projectionMatrix = math32.NewMatrix4().Zero()
	this._transformMatrix = math32.NewMatrix4().Zero()
	this._worldViewProjection = math32.NewMatrix4().Zero()

	return this
}

func (this *ShadowGenerator) GetTransformMatrix() *math32.Matrix4 {

	if this._cachedPosition != nil ||
		this._cachedDirection != nil ||
		!this._light.GetPosition().Equals(this._cachedPosition) ||
		!this._light.GetDirection().Equals(this._cachedDirection) {

		this._cachedPosition = this._light.GetPosition().Clone()
		this._cachedDirection = this._light.GetDirection().Clone()

		activeCamera := this._scene.ActiveCamera

		this._viewMatrix = math32.NewMatrix4().LookAtLH(this._light.GetPosition(), this._light.GetPosition().Add(this._light.GetDirection()), math32.NewVector3Up())
		this._projectionMatrix = math32.NewMatrix4().PerspectiveFovLH(math.Pi/2.0, 1.0, activeCamera.GetMinZ(), activeCamera.GetMaxZ())

		this._viewMatrix.MultiplyToRef(this._projectionMatrix, this._transformMatrix)
	}

	return this._transformMatrix
}

func (this *ShadowGenerator) Dispose() {
	this._shadowMap.Dispose()
}

/**
IShadowGenerator interface start
***/
func (this *ShadowGenerator) IsReady() bool {
	if this == nil {
		return false
	}
	return this._effect.IsReady() && this._effectVSM.IsReady()
}

func (this *ShadowGenerator) IsUseVarianceShadowMap() bool {
	return this.UseVarianceShadowMap
}

func (this *ShadowGenerator) GetShadowMap() IRenderTargetTexture {
	return this._shadowMap
}

/**
IShadowGenerator interface end
***/
