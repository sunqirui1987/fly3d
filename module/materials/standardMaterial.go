package materials

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/effects"
	"github.com/suiqirui1987/fly3d/module/lights"
	"strconv"
	"strings"
)

type StandardMaterial struct {
	Material

	DiffuseTexture    ITexture
	AmbientTexture    ITexture
	OpacityTexture    ITexture
	ReflectionTexture ITexture
	EmissiveTexture   ITexture
	SpecularTexture   ITexture
	BumpTexture       ITexture

	AmbientColor  *math32.Color3
	DiffuseColor  *math32.Color3
	SpecularColor *math32.Color3
	SpecularPower float32
	EmissiveColor *math32.Color3

	_cachedDefines string
	_renderTargets []interface{}

	//Internals
	_worldViewProjectionMatrix *math32.Matrix4
	_lightMatrix               *math32.Matrix4
	_globalAmbientColor        *math32.Color3
	_baseColor                 *math32.Color3
	_scaledDiffuse             *math32.Color3
	_scaledSpecular            *math32.Color3
}

func NewStandardMaterial(name string, scene *engines.Scene) *StandardMaterial {
	this := &StandardMaterial{}
	this.Name = name
	this.Id = name
	this._scene = scene
	this._scene.Materials = append(this._scene.Materials, this)

	this.Init()
	return this
}

func (this *StandardMaterial) Init() {

	this.Material.Init()

	this.AmbientColor = math32.NewColor3(0, 0, 0)
	this.DiffuseColor = math32.NewColor3(1, 1, 1)
	this.SpecularColor = math32.NewColor3(1, 1, 1)
	this.SpecularPower = 64
	this.EmissiveColor = math32.NewColor3(0, 0, 0)

	this._cachedDefines = ""
	this._renderTargets = make([]interface{}, 0)

	this._worldViewProjectionMatrix = math32.NewMatrix4().Zero()
	this._lightMatrix = math32.NewMatrix4().Zero()
	this._globalAmbientColor = math32.NewColor3(0, 0, 0)
	this._baseColor = math32.NewColor3(0, 0, 0)
	this._scaledDiffuse = math32.NewColor3(0, 0, 0)
	this._scaledSpecular = math32.NewColor3(0, 0, 0)
}

func (this *StandardMaterial) NeedAlphaBlending() bool {
	return (this.Alpha < 1.0) || (this.OpacityTexture != nil)
}

func (this *StandardMaterial) NeedAlphaTesting() bool {
	return this.DiffuseTexture != nil && this.DiffuseTexture.HasAlpha()
}
func (this *StandardMaterial) IsReady(mesh IMesh) bool {
	engine := this._scene.GetEngine()

	// Effect
	defines := make([]string, 0)

	// Textures
	if this.DiffuseTexture != nil {
		if !this.DiffuseTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define DIFFUSE")
		}

	}

	if this.AmbientTexture != nil {
		if !this.AmbientTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define AMBIENT")
		}

	}

	if this.OpacityTexture != nil {
		if !this.OpacityTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define OPACITY")
		}
	}

	if this.ReflectionTexture != nil {
		if !this.ReflectionTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define REFLECTION")
		}

	}

	if this.EmissiveTexture != nil {
		if !this.EmissiveTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define EMISSIVE")
		}
	}

	if this.SpecularTexture != nil {
		if !this.SpecularTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define SPECULAR")
		}
	}

	if engine.GetCaps().StandardDerivatives && this.BumpTexture != nil {
		if !this.BumpTexture.IsReady() {
			return false
		} else {
			defines = append(defines, "#define BUMP")
		}
	}

	if core.GlobalFly3D.ClipPlane != nil {
		defines = append(defines, "#define CLIPPLANE")
	}

	if engine.GetAlphaTesting() {
		defines = append(defines, "#define ALPHATEST")
	}

	// Fog
	if this._scene.FogMode != core.FOGMODE_NONE {
		defines = append(defines, "#define FOG")
	}

	shadowsActivated := false
	var lightIndex int
	lightIndex = 0
	for index := 0; index < len(this._scene.Lights); index++ {
		light := this._scene.Lights[index]

		if !light.IsEnabled() {
			continue
		}

		lightIndex_str := strconv.Itoa(lightIndex)

		defines = append(defines, "#define LIGHT"+lightIndex_str)

		if _, ok := light.(*lights.SpotLight); ok {
			defines = append(defines, "#define SPOTLIGHT"+lightIndex_str)
		} else if _, ok := light.(*lights.HemisphericLight); ok {
			defines = append(defines, "#define HEMILIGHT"+lightIndex_str)
		} else {
			defines = append(defines, "#define POINTDIRLIGHT"+lightIndex_str)
		}

		// Shadows
		shadowGenerator := light.GetShadowGenerator()
		if mesh != nil && mesh.IsReceiveShadows() == true && shadowGenerator != nil && shadowGenerator.IsReady() {
			defines = append(defines, "#define SHADOW"+lightIndex_str)

			if !shadowsActivated {
				defines = append(defines, "#define SHADOWS")
				shadowsActivated = true
			}

			if shadowGenerator.IsUseVarianceShadowMap() {
				defines = append(defines, "#define SHADOWVSM"+lightIndex_str)
			}
		}

		lightIndex++
		if lightIndex == 4 {
			break
		}

	}

	attribs := []string{"position", "normal"}
	if mesh != nil {
		if mesh.IsVerticesDataPresent(IMesh_VB_UVKind) {
			attribs = append(attribs, "uv")
			defines = append(defines, "#define UV1")
		}
		if mesh.IsVerticesDataPresent(IMesh_VB_UV2Kind) {
			attribs = append(attribs, "uv2")
			defines = append(defines, "#define UV2")
		}
		if mesh.IsVerticesDataPresent(IMesh_VB_ColorKind) {
			attribs = append(attribs, "color")
			defines = append(defines, "#define VERTEXCOLOR")
		}
	}

	// Get correct effect
	join := strings.Join(defines, "\n")
	if this._cachedDefines != join {
		this._cachedDefines = join

		// IE patch
		shaderName := "default"
		if core.GlobalFly3D.IsIE == true {
			shaderName = "iedefault"
		}

		this._effect = effects.CreateEffect(
			engine,
			shaderName,
			attribs,
			[]string{"world", "view", "worldViewProjection", "vEyePosition", "vLightsType", "vAmbientColor", "vDiffuseColor", "vSpecularColor", "vEmissiveColor",
				"vLightData0", "vLightDiffuse0", "vLightSpecular0", "vLightDirection0", "vLightGround0", "lightMatrix0",
				"vLightData1", "vLightDiffuse1", "vLightSpecular1", "vLightDirection1", "vLightGround1", "lightMatrix1",
				"vLightData2", "vLightDiffuse2", "vLightSpecular2", "vLightDirection2", "vLightGround2", "lightMatrix2",
				"vLightData3", "vLightDiffuse3", "vLightSpecular3", "vLightDirection3", "vLightGround3", "lightMatrix3",
				"vFogInfos", "vFogColor",
				"vDiffuseInfos", "vAmbientInfos", "vOpacityInfos", "vReflectionInfos", "vEmissiveInfos", "vSpecularInfos", "vBumpInfos",
				"vClipPlane", "diffuseMatrix", "ambientMatrix", "opacityMatrix", "reflectionMatrix", "emissiveMatrix", "specularMatrix", "bumpMatrix",
			},

			[]string{"diffuseSampler", "ambientSampler", "opacitySampler", "reflectionCubeSampler", "reflection2DSampler", "emissiveSampler", "specularSampler", "bumpSampler",
				"shadowSampler0", "shadowSampler1", "shadowSampler2", "shadowSampler3",
			},
			join)
	}
	if !this._effect.IsReady() {
		return false
	}

	return true
}

func (this *StandardMaterial) GetRenderTargetTextures() []ITexture {
	results := make([]ITexture, 0)

	if this.ReflectionTexture != nil && this.ReflectionTexture.IsRenderTarget() {
		results = append(results, this.ReflectionTexture)
	}

	return results
}

func (this *StandardMaterial) Unbind() {
	if this.ReflectionTexture != nil && this.ReflectionTexture.IsRenderTarget() {
		this._effect.SetTexture("reflection2DSampler", nil)
	}
}

func (this *StandardMaterial) Bind(world *math32.Matrix4, mesh IMesh) {
	baseColor := this.DiffuseColor

	// Values
	if this.DiffuseTexture != nil {
		this._effect.SetTexture("diffuseSampler", this.DiffuseTexture.GetGLTexture())

		this._effect.SetVector2("vDiffuseInfos", this.DiffuseTexture.GetCoordinatesIndex(), this.DiffuseTexture.GetLevel())
		this._effect.SetMatrix("diffuseMatrix", this.DiffuseTexture.ComputeTextureMatrix())

		baseColor = math32.NewColor3(1, 1, 1)
	}

	if this.AmbientTexture != nil {
		this._effect.SetTexture("ambientSampler", this.AmbientTexture.GetGLTexture())

		this._effect.SetVector2("vAmbientInfos", this.AmbientTexture.GetCoordinatesIndex(), this.AmbientTexture.GetLevel())
		this._effect.SetMatrix("ambientMatrix", this.AmbientTexture.ComputeTextureMatrix())
	}

	if this.OpacityTexture != nil {
		this._effect.SetTexture("opacitySampler", this.OpacityTexture.GetGLTexture())

		this._effect.SetVector2("vOpacityInfos", this.OpacityTexture.GetCoordinatesIndex(), this.OpacityTexture.GetLevel())
		this._effect.SetMatrix("opacityMatrix", this.OpacityTexture.ComputeTextureMatrix())
	}

	if this.ReflectionTexture != nil {
		iscube := 0
		if this.ReflectionTexture.GetGLTexture() != nil && this.ReflectionTexture.GetGLTexture().IsCube {
			iscube = 1
		}
		if iscube == 1 {
			this._effect.SetTexture("reflectionCubeSampler", this.ReflectionTexture.GetGLTexture())
		} else {
			this._effect.SetTexture("reflection2DSampler", this.ReflectionTexture.GetGLTexture())
		}

		this._effect.SetMatrix("reflectionMatrix", this.ReflectionTexture.ComputeReflectionTextureMatrix())
		this._effect.SetFloat3("vReflectionInfos", float32(this.ReflectionTexture.GetCoordinatesMode()), float32(this.ReflectionTexture.GetLevel()), float32(iscube))
	}

	if this.EmissiveTexture != nil {
		this._effect.SetTexture("emissiveSampler", this.EmissiveTexture.GetGLTexture())

		this._effect.SetVector2("vEmissiveInfos", this.EmissiveTexture.GetCoordinatesIndex(), this.EmissiveTexture.GetLevel())
		this._effect.SetMatrix("emissiveMatrix", this.EmissiveTexture.ComputeTextureMatrix())
	}

	if this.SpecularTexture != nil {
		this._effect.SetTexture("specularSampler", this.SpecularTexture.GetGLTexture())

		this._effect.SetVector2("vSpecularInfos", this.SpecularTexture.GetCoordinatesIndex(), this.SpecularTexture.GetLevel())
		this._effect.SetMatrix("specularMatrix", this.SpecularTexture.ComputeTextureMatrix())
	}

	if this.BumpTexture != nil && this._scene.GetEngine().GetCaps().StandardDerivatives {
		this._effect.SetTexture("bumpSampler", this.BumpTexture.GetGLTexture())

		this._effect.SetVector2("vBumpInfos", this.BumpTexture.GetCoordinatesIndex(), this.BumpTexture.GetLevel())
		this._effect.SetMatrix("bumpMatrix", this.BumpTexture.ComputeTextureMatrix())
	}

	this._worldViewProjectionMatrix = world.Multiply(this._scene.GetTransformMatrix())
	this._globalAmbientColor = this._scene.AmbientColor.Multiply(this.AmbientColor)

	this._effect.SetMatrix("world", world)
	this._effect.SetMatrix("worldViewProjection", this._worldViewProjectionMatrix)
	this._effect.SetVector3("vEyePosition", this._scene.ActiveCamera.GetPosition())
	this._effect.SetColor3("vAmbientColor", this._globalAmbientColor)
	this._effect.SetColor4("vDiffuseColor", baseColor, this.Alpha*mesh.GetVisibility())
	this._effect.SetColor4("vSpecularColor", this.SpecularColor, this.SpecularPower)
	this._effect.SetColor3("vEmissiveColor", this.EmissiveColor)

	lightIndex := 0
	for index := 0; index < len(this._scene.Lights); index++ {
		light := this._scene.Lights[index]

		if !light.IsEnabled() {
			continue
		}

		lightIndex_str := strconv.Itoa(lightIndex)

		if polight, ok := light.(*lights.PointLight); ok {
			// Point Light
			this._effect.SetFloat4("vLightData"+lightIndex_str, polight.Position.X, polight.Position.Y, polight.Position.Z, 0)
		} else if dlight, ok := light.(*lights.DirectionalLight); ok {
			// Directional Light
			this._effect.SetFloat4("vLightData"+lightIndex_str, dlight.Direction.X, dlight.Direction.Y, dlight.Direction.Z, 1)
		} else if slight, ok := light.(*lights.SpotLight); ok {
			// Spot Light
			this._effect.SetFloat4("vLightData"+lightIndex_str, slight.Position.X, slight.Position.Y, slight.Position.Z, slight.Exponent)
			normalizeDirection := slight.Direction.NormalizeTo()
			this._effect.SetFloat4("vLightDirection"+lightIndex_str, normalizeDirection.X, normalizeDirection.Y, normalizeDirection.Z, math32.Cos(slight.Angle*0.5))
		} else if hlight, ok := light.(*lights.HemisphericLight); ok {
			// Hemispheric Light
			normalizeDirection := hlight.Direction.NormalizeTo()
			this._effect.SetFloat4("vLightData"+lightIndex_str, normalizeDirection.X, normalizeDirection.Y, normalizeDirection.Z, 0)
			this._effect.SetColor3("vLightGround"+lightIndex_str, hlight.GroundColor.Scale(hlight.Intensity))
		}
		this._scaledDiffuse = light.GetDiffuse().Scale(light.GetIntensity())
		this._scaledSpecular = light.GetSpecular().Scale(light.GetIntensity())

		this._effect.SetColor3("vLightDiffuse"+lightIndex_str, this._scaledDiffuse)
		this._effect.SetColor3("vLightSpecular"+lightIndex_str, this._scaledSpecular)

		lightIndex++

		if lightIndex == 4 {
			break
		}

	}

	if core.GlobalFly3D.ClipPlane != nil {
		this._effect.SetFloat4("vClipPlane", core.GlobalFly3D.ClipPlane.Normal.X, core.GlobalFly3D.ClipPlane.Normal.Y, core.GlobalFly3D.ClipPlane.Normal.Z, core.GlobalFly3D.ClipPlane.D)
	}

	// View
	if this._scene.FogMode != core.FOGMODE_NONE || this.ReflectionTexture != nil {
		this._effect.SetMatrix("view", this._scene.GetViewMatrix())
	}

	// Fog
	if this._scene.FogMode != core.FOGMODE_NONE {
		this._effect.SetFloat4("vFogInfos", float32(this._scene.FogMode), this._scene.FogStart, this._scene.FogEnd, this._scene.FogDensity)
		this._effect.SetColor3("vFogColor", this._scene.FogColor)
	}

}

func (this *StandardMaterial) Dispose() {
	if this.DiffuseTexture != nil {
		this.DiffuseTexture.Dispose()
	}

	if this.AmbientTexture != nil {
		this.AmbientTexture.Dispose()
	}

	if this.OpacityTexture != nil {
		this.OpacityTexture.Dispose()
	}

	if this.ReflectionTexture != nil {
		this.ReflectionTexture.Dispose()
	}

	if this.EmissiveTexture != nil {
		this.EmissiveTexture.Dispose()
	}

	if this.SpecularTexture != nil {
		this.SpecularTexture.Dispose()
	}

	if this.BumpTexture != nil {
		this.BumpTexture.Dispose()
	}

	this.BaseDispose()
}

//IAnimationTarget
func (this *StandardMaterial) GetAnimatables() []IAnimatable {
	return nil
}
func (this *StandardMaterial) GetAnimations() []IAnimation {
	return nil
}
