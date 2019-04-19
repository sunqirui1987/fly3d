package layers

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/effects"
	"github.com/suiqirui1987/fly3d/module/textures"
	"github.com/suiqirui1987/fly3d/tools"
)

type Layer struct {
	Name         string
	Isbackground bool
	Color        *math32.Color4

	_texture ITexture
	_scene   *engines.Scene

	_vertexDeclaration [1]int
	_vertexStrideSize  int
	_vertexBuffer      *gl.GLVertexBuffer
	_indexBuffer       *gl.GLIndexBuffer
	_effect            IEffect

	OnDispose func()
}

func NewLayer(name string, imgUrl string, scene *engines.Scene, isBackgroud bool, color *math32.Color4) *Layer {

	this := &Layer{}
	this.Name = name
	if imgUrl == "" {
		this._texture = nil
	} else {
		this._texture = textures.NewTexture(imgUrl, scene, true, -1)
	}

	this._scene = scene
	this.Isbackground = isBackgroud
	if color != nil {
		this.Color = color
	} else {
		this.Color = math32.NewColor4(1, 1, 1, 1)
	}

	// VBO

	vertices := []float32{
		1, 1,
		-1, 1,
		-1, -1,
		1, -1,
	}

	this._vertexDeclaration[0] = 2
	this._vertexStrideSize = 2 * 4

	this._vertexBuffer = this._scene.GetEngine().CreateVertexBuffer(vertices)

	// Indices
	indices := []uint16{
		0, 1, 2,
		0, 2, 3,
	}

	this._indexBuffer = scene.GetEngine().CreateIndexBuffer(indices, false)

	// Effects
	this._effect = effects.CreateEffect(this._scene.GetEngine(), "layer",
		[]string{"position"},
		[]string{"textureMatrix", "color"},
		[]string{"textureSampler"}, "")

	this._scene.Layers = append(this._scene.Layers, this)
	return this
}

/**** interface start ****/
/*
type ILayer interface {
	IsBackground() bool
	Render()
	Dispose()
}
*/

func (this *Layer) IsBackground() bool {
	return this.Isbackground
}

func (this *Layer) Render() {

	// Check
	if !this._effect.IsReady() || this._texture == nil || !this._texture.IsReady() {
		return
	}
	engine := this._scene.GetEngine()

	// Render
	engine.EnableEffect(this._effect)
	engine.SetState(false)

	// Texture
	this._effect.SetTexture("textureSampler", this._texture.GetGLTexture())
	this._effect.SetMatrix("textureMatrix", this._texture.ComputeTextureMatrix())

	// Color
	this._effect.SetFloat4("color", this.Color.R, this.Color.G, this.Color.B, this.Color.A)

	// VBOs
	engine.BindBuffers(this._vertexBuffer, this._indexBuffer, this._vertexDeclaration[:], this._vertexStrideSize, this._effect)

	// Draw order
	engine.SetAlphaMode(core.ALPHA_COMBINE)
	engine.Draw(true, 0, 6)
	engine.SetAlphaMode(core.ALPHA_DISABLE)
}

func (this *Layer) Dispose() {
	if this._vertexBuffer != nil {
		this._scene.GetEngine().ReleaseVertexBuffer(this._vertexBuffer)
		this._vertexBuffer = nil
	}

	if this._indexBuffer != nil {
		this._scene.GetEngine().ReleaseIndexBuffer(this._indexBuffer)
		this._indexBuffer = nil
	}

	if this._texture != nil {
		this._texture.Dispose()
		this._texture = nil
	}

	// Remove from scene
	index := tools.IndexOf(this, this._scene.Layers)
	if index > -1 {
		this._scene.Layers = append(this._scene.Layers[:index], this._scene.Layers[index+1:]...)
	}

	// Callback
	if this.OnDispose != nil {
		this.OnDispose()
	}
}

/**** interface end ****/
