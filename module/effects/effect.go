package effects

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
	"log"
	"reflect"
)

type Effect struct {
	_engine *engines.Engine
	Name    string
	Defines string

	_attributesNames []string
	_uniformsNames   []string
	_samplers        []string
	_valueCache      map[string]interface{}
	_isReady         bool

	_program    gl.Program
	_attributes []gl.Attrib
	_uniforms   []gl.Uniform
}

func CreateEffect(engine *engines.Engine, baseName string, attributesNames []string, uniformsNames []string, samplers []string, defines string) *Effect {
	name := baseName + "@" + defines
	var e *Effect

	effect, ok := engine.CompiledEffects[name]
	if e, ok = effect.(*Effect); ok == true {
		return e
	}

	e = NewEffect(baseName, attributesNames, uniformsNames, samplers, engine, defines)
	engine.CompiledEffects[name] = e

	return e
}

func NewEffect(baseName string, attributesNames []string, uniformsNames []string, samplers []string, engine *engines.Engine, defines string) *Effect {

	this := &Effect{}

	this._engine = engine
	this.Name = baseName
	this.Defines = defines
	this._attributesNames = attributesNames
	this._uniformsNames = append(uniformsNames, samplers...)
	this._samplers = samplers
	this._isReady = false
	this._valueCache = map[string]interface{}{}

	// Is in local store
	if _, ok := ShadersStore[baseName+"_vertex"]; ok {
		vs_str, _ := ShadersStore[baseName+"_vertex"]
		ps_str, _ := ShadersStore[baseName+"_fragment"]

		this._prepareEffect(vs_str, ps_str, attributesNames, defines)

	} else {
		shaderUrl := core.GlobalFly3D.ResRepository + core.GlobalFly3D.ShadersRepository + baseName

		that := this
		// Vertex shader
		tools.LoadFile(shaderUrl+".vertex.fx",
			func(vertexSourceCode string) {
				// Fragment shader
				tools.LoadFile(shaderUrl+".fragment.fx",
					func(fragmentSourceCode string) {

						that._prepareEffect(vertexSourceCode, fragmentSourceCode, attributesNames, defines)
					}, nil)

			}, nil)
	}
	return this
}

// Properties
func (this *Effect) IsReady() bool {
	return this._isReady
}

func (this *Effect) GetProgram() gl.Program {
	return this._program
}

func (this *Effect) GetAttribute(index int) gl.Attrib {
	return this._attributes[index]
}
func (this *Effect) GetAttributesNames() []string {
	return this._attributesNames
}

func (this *Effect) GetAttributesCount() int {
	return len(this._attributes)
}

func (this *Effect) GetUniformIndex(uniformName string) int {
	index := -1
	for i, name := range this._uniformsNames {
		if name == uniformName {
			index = i
		}
	}
	return index
}

func (this *Effect) GetUniform(uniformName string) gl.Uniform {
	index := this.GetUniformIndex(uniformName)
	if index > -1 {
		return this._uniforms[index]
	}
	return gl.Uniform{}
}

func (this *Effect) GetSamplers() []string {
	return this._samplers
}

func (this *Effect) _prepareEffect(vertexSourceCode string, fragmentSourceCode string, attributesNames []string, defines string) {
	log.Printf("start _prepareEffect \r\n")
	engine := this._engine
	this._program = engine.CreateShaderProgram(vertexSourceCode, fragmentSourceCode, defines)

	this._uniforms = engine.GetUniforms(this._program, this._uniformsNames)
	this._attributes = engine.GetAttributes(this._program, attributesNames)

	for index := 0; index < len(this._samplers); index++ {
		sampler := this.GetUniform(this._samplers[index])

		if !sampler.Valid() {
			this._samplers = append(this._samplers[:index], this._samplers[index+1:]...)
			index--
		}
	}

	engine.BindSamplers(this)

	this._isReady = true
}

func (this *Effect) SetTexture(channel string, texture *gl.GLTextureBuffer) {
	index := tools.IndexOf(channel, this._samplers)
	if index > -1 {
		this._engine.SetTexture(index, texture)
	}

}

func (this *Effect) SetMatrix(uniformName string, val *math32.Matrix4) {
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetMatrix(this.GetUniform(uniformName), val)
}

func (this *Effect) SetBool(uniformName string, val bool) {
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetBool(this.GetUniform(uniformName), val)
}

func (this *Effect) SetVector2(uniformName string, x, y float32) {
	val := math32.NewVector2(x, y)
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetVector2(this.GetUniform(uniformName), val)
}
func (this *Effect) SetVector2i(uniformName string, x, y int) {
	val := math32.NewVector2(float32(x), float32(y))
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetVector2(this.GetUniform(uniformName), val)
}

func (this *Effect) SetVector3(uniformName string, val *math32.Vector3) {

	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetVector3(this.GetUniform(uniformName), val)
}

func (this *Effect) SetFloat2(uniformName string, x, y float32) {
	val := math32.NewVector2(x, y)
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetFloat2(this.GetUniform(uniformName), x, y)
}

func (this *Effect) SetFloat3(uniformName string, x, y, z float32) {
	val := math32.NewVector3(x, y, z)
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetFloat3(this.GetUniform(uniformName), x, y, z)
}

func (this *Effect) SetFloat4(uniformName string, x, y, z, w float32) {
	val := math32.NewVector4(x, y, z, w)
	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetFloat4(this.GetUniform(uniformName), x, y, z, w)
}

func (this *Effect) SetColor3(uniformName string, val *math32.Color3) {

	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetColor3(this.GetUniform(uniformName), val)
}

func (this *Effect) SetColor42(uniformName string, val *math32.Color4) {

	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetColor4(this.GetUniform(uniformName), val)
}
func (this *Effect) SetColor4(uniformName string, c3 *math32.Color3, a float32) {
	val := math32.NewColor4(c3.R, c3.G, c3.B, a)

	uniform_val, ok := this._valueCache[uniformName]
	if ok && reflect.DeepEqual(uniform_val, val) {
		return
	}

	this._valueCache[uniformName] = val
	this._engine.SetColor4(this.GetUniform(uniformName), val)
}
