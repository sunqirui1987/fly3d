package interfaces

import (
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/math32"
)

type IEffect interface {
	IsReady() bool
	GetProgram() gl.Program
	GetAttribute(index int) gl.Attrib
	GetAttributesNames() []string
	GetAttributesCount() int
	GetUniformIndex(uniformName string) int
	GetUniform(uniformName string) gl.Uniform
	GetSamplers() []string
	SetTexture(channel string, texture *gl.GLTextureBuffer)
	SetMatrix(uniformName string, val *math32.Matrix4)
	SetBool(uniformName string, val bool)
	SetVector2(uniformName string, x, y float32)
	SetVector2i(uniformName string, x, y int)
	SetVector3(uniformName string, val *math32.Vector3)
	SetFloat2(uniformName string, x, y float32)
	SetFloat3(uniformName string, x, y, z float32)
	SetFloat4(uniformName string, x, y, z, w float32)
	SetColor3(uniformName string, val *math32.Color3)
	SetColor4(uniformName string, c3 *math32.Color3, a float32)
	SetColor42(uniformName string, val *math32.Color4)
}
