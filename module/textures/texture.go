package textures

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	"github.com/suiqirui1987/fly3d/math32"
)

type Texture struct {
	BaseTexture
	UOffset float32
	VOffset float32
	UScale  float32
	VScale  float32

	UAng float32
	VAng float32
	WAng float32

	_cachedUOffset float32
	_cachedVOffset float32
	_cachedUScale  float32
	_cachedVScale  float32

	_cachedUAng float32
	_cachedVAng float32
	_cachedWAng float32

	_cachedTextureMatrix   *math32.Matrix4
	_projectionModeMatrix  *math32.Matrix4
	_rowGenerationMatrix   *math32.Matrix4
	_cachedCoordinatesMode int
	_t0                    *math32.Vector3
	_t1                    *math32.Vector3
	_t2                    *math32.Vector3
}

func NewTexture(url string, scene *engines.Scene, noMipmap bool, invertY int) *Texture {
	if url == "" {
		return nil
	}
	this := &Texture{}

	this._scene = scene
	this._texture = this._getFromCache(url, noMipmap)

	if this._texture == nil {
		this._texture = this._scene.GetEngine().CreateTexture(url, noMipmap, invertY, scene)
	}

	this._scene.Textures = append(this._scene.Textures, this)

	this.Init()
	return this
}

func (this *Texture) Init() {
	this.BaseTexture.Init()

	this.UOffset = 0.0
	this.VOffset = 0.0

	this.UScale = 1.0
	this.VScale = 1.0

	this.UAng = 0.0
	this.VAng = 0.0
	this.WAng = 0.0

	this._texture.WrapU = gl.WRAP_ADDRESSMODE
	this._texture.WrapV = gl.WRAP_ADDRESSMODE

	this.CoordinatesIndex = 0
	this.CoordinatesMode = EXPLICIT_MODE

}
func (this *Texture) _prepareRowForTextureGeneration(x, y, z float32, t *math32.Vector3) {
	x -= this.UOffset + 0.5
	y -= this.VOffset + 0.5
	z -= 0.5

	math32.NewVector3Zero().TransformCoordinatesFromFloatsToRef(x, y, z, this._rowGenerationMatrix, t)

	t.X *= this.UScale
	t.Y *= this.VScale

	t.X += 0.5
	t.Y += 0.5
	t.Z += 0.5

}

func (this *Texture) ComputeTextureMatrix() *math32.Matrix4 {
	if this.UOffset == this._cachedUOffset &&
		this.VOffset == this._cachedVOffset &&
		this.UScale == this._cachedUScale &&
		this.VScale == this._cachedVScale &&
		this.UAng == this._cachedUAng &&
		this.VAng == this._cachedVAng &&
		this.WAng == this._cachedWAng {
		return this._cachedTextureMatrix
	}

	this._cachedUOffset = this.UOffset
	this._cachedVOffset = this.VOffset
	this._cachedUScale = this.UScale
	this._cachedVScale = this.VScale
	this._cachedUAng = this.UAng
	this._cachedVAng = this.VAng
	this._cachedWAng = this.WAng

	if this._cachedTextureMatrix == nil {
		this._cachedTextureMatrix = math32.NewMatrix4().Zero()
		this._rowGenerationMatrix = math32.NewMatrix4()
		this._t0 = math32.NewVector3Zero()
		this._t1 = math32.NewVector3Zero()
		this._t2 = math32.NewVector3Zero()
	}

	math32.NewMatrix4().RotationYawPitchRollToRef(this.VAng, this.UAng, this.WAng, this._rowGenerationMatrix)

	this._prepareRowForTextureGeneration(0, 0, 0, this._t0)
	this._prepareRowForTextureGeneration(1.0, 0, 0, this._t1)
	this._prepareRowForTextureGeneration(0, 1.0, 0, this._t2)

	this._t1 = this._t1.Sub(this._t0)
	this._t2 = this._t2.Sub(this._t0)

	this._cachedTextureMatrix = this._cachedTextureMatrix.Identity()
	this._cachedTextureMatrix[0] = this._t1.X
	this._cachedTextureMatrix[1] = this._t1.Y
	this._cachedTextureMatrix[2] = this._t1.Z
	this._cachedTextureMatrix[4] = this._t2.X
	this._cachedTextureMatrix[5] = this._t2.Y
	this._cachedTextureMatrix[6] = this._t2.Z
	this._cachedTextureMatrix[8] = this._t0.X
	this._cachedTextureMatrix[9] = this._t0.Y
	this._cachedTextureMatrix[10] = this._t0.Z

	return this._cachedTextureMatrix
}
func (this *Texture) ComputeReflectionTextureMatrix() *math32.Matrix4 {
	if this.UOffset == this._cachedUOffset &&
		this.VOffset == this._cachedVOffset &&
		this.UScale == this._cachedUScale &&
		this.VScale == this._cachedVScale &&
		this.CoordinatesMode == this._cachedCoordinatesMode {
		return this._cachedTextureMatrix
	}

	if this._cachedTextureMatrix == nil {
		this._cachedTextureMatrix = math32.NewMatrix4().Zero()
		this._projectionModeMatrix = math32.NewMatrix4().Zero()
	}

	switch this.CoordinatesMode {
	case SPHERICAL_MODE:
		this._cachedTextureMatrix = this._cachedTextureMatrix.Identity()
		this._cachedTextureMatrix[0] = -0.5 * this.UScale
		this._cachedTextureMatrix[5] = -0.5 * this.VScale
		this._cachedTextureMatrix[12] = 0.5 + this.UOffset
		this._cachedTextureMatrix[13] = 0.5 + this.VOffset
		break
	case PLANAR_MODE:
		this._cachedTextureMatrix = this._cachedTextureMatrix.Identity()
		this._cachedTextureMatrix[0] = this.UScale
		this._cachedTextureMatrix[5] = this.VScale
		this._cachedTextureMatrix[12] = this.UOffset
		this._cachedTextureMatrix[13] = this.VOffset
		break
	case PROJECTION_MODE:
		this._projectionModeMatrix = this._projectionModeMatrix.Identity()
		this._projectionModeMatrix[0] = 0.5
		this._projectionModeMatrix[5] = -0.5
		this._projectionModeMatrix[10] = 0.0
		this._projectionModeMatrix[12] = 0.5
		this._projectionModeMatrix[13] = 0.5
		this._projectionModeMatrix[14] = 1.0
		this._projectionModeMatrix[15] = 1.0

		this._cachedTextureMatrix = this._scene.GetProjectionMatrix().Multiply(this._projectionModeMatrix)
		break
	}

	return this._cachedTextureMatrix
}
