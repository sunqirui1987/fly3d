package textures

import (
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
)

type RenderTargetTexture struct {
	Texture
	Name string

	OnBeforeRender func()
	OnAfterRender  func()

	CustomRenderFunction func([]ISubMesh, []ISubMesh, []ISubMesh, []IMesh)

	_waitingRenderList []string
	_renderList        []IMesh

	_opaqueSubMeshes      []ISubMesh
	_transparentSubMeshes []ISubMesh
	_alphaTestSubMeshes   []ISubMesh
}

func NewRenderTargetTexture(name string, size int, scene *engines.Scene, generateMipMaps bool) *RenderTargetTexture {
	this := &RenderTargetTexture{}
	this.Name = name
	this._scene = scene
	this._scene.Textures = append(this._scene.Textures, this)

	this._texture = scene.GetEngine().CreateRenderTargetTexture(size, generateMipMaps)
	return this
}

func (this *RenderTargetTexture) Init() {
	this.Texture.Init()

	this._waitingRenderList = make([]string, 0)
	this._renderList = make([]IMesh, 0)

	this._opaqueSubMeshes = make([]ISubMesh, 0)
	this._transparentSubMeshes = make([]ISubMesh, 0)
	this._alphaTestSubMeshes = make([]ISubMesh, 0)

}

func (this *RenderTargetTexture) IsRenderTarget() bool {
	return true
}

func (this *RenderTargetTexture) GetRenderList() []IMesh {
	return this._renderList
}

func (this *RenderTargetTexture) Resize(size int, generateMipMaps bool) {
	this.ReleaseGLTexture()
	this._texture = this._scene.GetEngine().CreateRenderTargetTexture(size, generateMipMaps)
}

func (this *RenderTargetTexture) Render() {

	if this.OnBeforeRender != nil {
		this.OnBeforeRender()
	}

	scene := this._scene
	engine := scene.GetEngine()

	for index := 0; index < len(this._waitingRenderList); index++ {
		id := this._waitingRenderList[index]
		this._renderList = append(this._renderList, this._scene.GetMeshByID(id))

	}
	this._waitingRenderList = make([]string, 0)

	if len(this._renderList) == 0 {
		return
	}

	// Bind
	engine.BindFramebuffer(this._texture)

	// Clear
	engine.Clear(scene.ClearColor, true, true)

	// Dispatch subMeshes
	this._opaqueSubMeshes = make([]ISubMesh, 0)
	this._transparentSubMeshes = make([]ISubMesh, 0)
	this._alphaTestSubMeshes = make([]ISubMesh, 0)

	for meshIndex := 0; meshIndex < len(this._renderList); meshIndex++ {
		mesh := this._renderList[meshIndex]

		if mesh.IsEnabled() && mesh.IsVisible() {

			for _, subMesh := range mesh.GetSubMeshes() {
				material := subMesh.GetMaterial()

				if material.NeedAlphaTesting() { // Alpha test
					this._alphaTestSubMeshes = append(this._alphaTestSubMeshes, subMesh)
				} else if material.NeedAlphaBlending() { // Transparent
					if material.GetAlpha() > 0 {
						//Opaque
						this._transparentSubMeshes = append(this._transparentSubMeshes, subMesh)

					}
				} else {
					this._opaqueSubMeshes = append(this._opaqueSubMeshes, subMesh)
				}
			}
		}
	}

	// Render
	if this.CustomRenderFunction != nil {
		this.CustomRenderFunction(this._opaqueSubMeshes, this._alphaTestSubMeshes, this._transparentSubMeshes, this._renderList)
	} else {
		scene.LocalRender(this._opaqueSubMeshes, this._alphaTestSubMeshes, this._transparentSubMeshes, this._renderList)
	}
	// Unbind
	engine.UnBindFramebuffer(this._texture)

	if this.OnAfterRender != nil {
		this.OnAfterRender()
	}

}
