
package effects

var ShadersStore = map[string]string{}

func init() { 

ShadersStore["default_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

#define MAP_PROJECTION	4.

// Constants
uniform vec3 vEyePosition;
uniform vec3 vAmbientColor;
uniform vec4 vDiffuseColor;
uniform vec4 vSpecularColor;
uniform vec3 vEmissiveColor;

// Input
varying vec3 vPositionW;
varying vec3 vNormalW;

#ifdef VERTEXCOLOR
varying vec3 vColor;
#endif

// Lights
#ifdef LIGHT0
uniform vec4 vLightData0;
uniform vec3 vLightDiffuse0;
uniform vec3 vLightSpecular0;
#ifdef SHADOW0
varying vec4 vPositionFromLight0;
uniform sampler2D shadowSampler0;
#endif
#ifdef SPOTLIGHT0
uniform vec4 vLightDirection0;
#endif
#ifdef HEMILIGHT0
uniform vec3 vLightGround0;
#endif
#endif

#ifdef LIGHT1
uniform vec4 vLightData1;
uniform vec3 vLightDiffuse1;
uniform vec3 vLightSpecular1;
#ifdef SHADOW1
varying vec4 vPositionFromLight1;
uniform sampler2D shadowSampler1;
#endif
#ifdef SPOTLIGHT1
uniform vec4 vLightDirection1;
#endif
#ifdef HEMILIGHT1
uniform vec3 vLightGround1;
#endif
#endif

#ifdef LIGHT2
uniform vec4 vLightData2;
uniform vec3 vLightDiffuse2;
uniform vec3 vLightSpecular2;
#ifdef SHADOW2
varying vec4 vPositionFromLight2;
uniform sampler2D shadowSampler2;
#endif
#ifdef SPOTLIGHT2
uniform vec4 vLightDirection2;
#endif
#ifdef HEMILIGHT2
uniform vec3 vLightGround2;
#endif
#endif

#ifdef LIGHT3
uniform vec4 vLightData3;
uniform vec3 vLightDiffuse3;
uniform vec3 vLightSpecular3;
#ifdef SHADOW3
varying vec4 vPositionFromLight3;
uniform sampler2D shadowSampler3;
#endif
#ifdef SPOTLIGHT3
uniform vec4 vLightDirection3;
#endif
#ifdef HEMILIGHT3
uniform vec3 vLightGround3;
#endif
#endif

// Samplers
#ifdef DIFFUSE
varying vec2 vDiffuseUV;
uniform sampler2D diffuseSampler;
uniform vec2 vDiffuseInfos;
#endif

#ifdef AMBIENT
varying vec2 vAmbientUV;
uniform sampler2D ambientSampler;
uniform vec2 vAmbientInfos;
#endif

#ifdef OPACITY	
varying vec2 vOpacityUV;
uniform sampler2D opacitySampler;
uniform vec2 vOpacityInfos;
#endif

#ifdef REFLECTION
varying vec3 vReflectionUVW;
uniform samplerCube reflectionCubeSampler;
uniform sampler2D reflection2DSampler;
uniform vec3 vReflectionInfos;
#endif

#ifdef EMISSIVE
varying vec2 vEmissiveUV;
uniform vec2 vEmissiveInfos;
uniform sampler2D emissiveSampler;
#endif

#ifdef SPECULAR
varying vec2 vSpecularUV;
uniform vec2 vSpecularInfos;
uniform sampler2D specularSampler;
#endif

// Shadows
#ifdef SHADOWS

float unpack(vec4 color)
{
	const vec4 bitShift = vec4(1. / (255. * 255. * 255.), 1. / (255. * 255.), 1. / 255., 1.);
	return dot(color, bitShift);
}

float unpackHalf(vec2 color) 
{ 
	return color.x + (color.y / 255.0);
}

float computeShadow(vec4 vPositionFromLight, sampler2D shadowSampler)
{
	vec3 depth = vPositionFromLight.xyz / vPositionFromLight.w;
	vec2 uv = 0.5 * depth.xy + vec2(0.5, 0.5);

	if (uv.x < 0. || uv.x > 1.0 || uv.y < 0. || uv.y > 1.0)
	{
		return 1.0;
	}

	float shadow = unpack(texture2D(shadowSampler, uv));

	if (depth.z > shadow)
	{
		return 0.;
	}
	return 1.;
}

// Thanks to http://devmaster.net/
float ChebychevInequality(vec2 moments, float t) 
{
	if (t <= moments.x)
	{
		return 1.0;
	}
	
	float variance = moments.y - (moments.x * moments.x); 
	variance = max(variance, 0.);

	float d = t - moments.x; 	
	return variance / (variance + d * d); 
}

float computeShadowWithVSM(vec4 vPositionFromLight, sampler2D shadowSampler)
{
	vec3 depth = vPositionFromLight.xyz / vPositionFromLight.w;
	vec2 uv = 0.5 * depth.xy + vec2(0.5, 0.5);

	if (uv.x < 0. || uv.x > 1.0 || uv.y < 0. || uv.y > 1.0)
	{
		return 1.0;
	}

	vec4 texel = texture2D(shadowSampler, uv);

	vec2 moments = vec2(unpackHalf(texel.xy), unpackHalf(texel.zw));
	return clamp(1.3 - ChebychevInequality(moments, depth.z), 0., 1.0);
}
#endif

// Bump
#ifdef BUMP
#extension GL_OES_standard_derivatives : enable
varying vec2 vBumpUV;
uniform vec2 vBumpInfos;
uniform sampler2D bumpSampler;

// Thanks to http://www.thetenthplanet.de/archives/1180
mat3 cotangent_frame(vec3 normal, vec3 p, vec2 uv)
{
	// get edge vectors of the pixel triangle
	vec3 dp1 = dFdx(p);
	vec3 dp2 = dFdy(p);
	vec2 duv1 = dFdx(uv);
	vec2 duv2 = dFdy(uv);

	// solve the linear system
	vec3 dp2perp = cross(dp2, normal);
	vec3 dp1perp = cross(normal, dp1);
	vec3 tangent = dp2perp * duv1.x + dp1perp * duv2.x;
	vec3 binormal = dp2perp * duv1.y + dp1perp * duv2.y;

	// construct a scale-invariant frame 
	float invmax = inversesqrt(max(dot(tangent, tangent), dot(binormal, binormal)));
	return mat3(tangent * invmax, binormal * invmax, normal);
}

vec3 perturbNormal(vec3 viewDir)
{
	vec3 map = texture2D(bumpSampler, vBumpUV).xyz * vBumpInfos.y;
	map = map * 255. / 127. - 128. / 127.;
	mat3 TBN = cotangent_frame(vNormalW, -viewDir, vBumpUV);
	return normalize(TBN * map);
}
#endif

#ifdef CLIPPLANE
varying float fClipDistance;
#endif

// Fog
#ifdef FOG

#define FOGMODE_NONE    0.
#define FOGMODE_EXP     1.
#define FOGMODE_EXP2    2.
#define FOGMODE_LINEAR  3.
#define E 2.71828

uniform vec4 vFogInfos;
uniform vec3 vFogColor;
varying float fFogDistance;

float CalcFogFactor()
{
	float fogCoeff = 1.0;
	float fogStart = vFogInfos.y;
	float fogEnd = vFogInfos.z;
	float fogDensity = vFogInfos.w;

	if (FOGMODE_LINEAR == vFogInfos.x)
	{
		fogCoeff = (fogEnd - fFogDistance) / (fogEnd - fogStart);
	}
	else if (FOGMODE_EXP == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fogDensity);
	}
	else if (FOGMODE_EXP2 == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fFogDistance * fogDensity * fogDensity);
	}

	return min(1., max(0., fogCoeff));
}
#endif

// Light Computing
struct lightingInfo
{
	vec3 diffuse;
	vec3 specular;
};

lightingInfo computeLighting(vec3 viewDirectionW, vec3 vNormal, vec4 lightData, vec3 diffuseColor, vec3 specularColor) {
	lightingInfo result;

	vec3 lightVectorW;
	if (lightData.w == 0.)
	{
		lightVectorW = normalize(lightData.xyz - vPositionW);
	}
	else
	{
		lightVectorW = normalize(-lightData.xyz);
	}

	// diffuse
	float ndl = max(0., dot(vNormal, lightVectorW));

	// Specular
	vec3 angleW = normalize(viewDirectionW + lightVectorW);
	float specComp = max(0., dot(vNormal, angleW));
	specComp = pow(specComp, vSpecularColor.a);

	result.diffuse = ndl * diffuseColor;
	result.specular = specComp * specularColor;

	return result;
}

lightingInfo computeSpotLighting(vec3 viewDirectionW, vec3 vNormal, vec4 lightData, vec4 lightDirection, vec3 diffuseColor, vec3 specularColor) {
	lightingInfo result;

	vec3 lightVectorW = normalize(lightData.xyz - vPositionW);

	// diffuse
	float cosAngle = max(0., dot(-lightDirection.xyz, lightVectorW));
	float spotAtten = 0.0;

	if (cosAngle >= lightDirection.w)
	{
		cosAngle = max(0., pow(cosAngle, lightData.w));
		spotAtten = max(0., (cosAngle - lightDirection.w) / (1. - cosAngle));

		// Diffuse
		float ndl = max(0., dot(vNormal, -lightDirection.xyz));

		// Specular
		vec3 angleW = normalize(viewDirectionW - lightDirection.xyz);
		float specComp = max(0., dot(vNormal, angleW));
		specComp = pow(specComp, vSpecularColor.a);

		result.diffuse = ndl * spotAtten * diffuseColor;
		result.specular = specComp * specularColor * spotAtten;

		return result;
	}

	result.diffuse = vec3(0.);
	result.specular = vec3(0.);

	return result;
}

lightingInfo computeHemisphericLighting(vec3 viewDirectionW, vec3 vNormal, vec4 lightData, vec3 diffuseColor, vec3 specularColor, vec3 groundColor) {
	lightingInfo result;

	// Diffuse
	float ndl = dot(vNormal, lightData.xyz) * 0.5 + 0.5;

	// Specular
	vec3 angleW = normalize(viewDirectionW + lightData.xyz);
	float specComp = max(0., dot(vNormal, angleW));
	specComp = pow(specComp, vSpecularColor.a);

	result.diffuse = mix(groundColor, diffuseColor, ndl);
	result.specular = specComp * specularColor;

	return result;
}

void main(void) {
	// Clip plane
#ifdef CLIPPLANE
	if (fClipDistance > 0.0)
		discard;
#endif

	vec3 viewDirectionW = normalize(vEyePosition - vPositionW);

	// Base color
	vec4 baseColor = vec4(1., 1., 1., 1.);
	vec3 diffuseColor = vDiffuseColor.rgb;

#ifdef VERTEXCOLOR
	diffuseColor *= vColor;
#endif

#ifdef DIFFUSE
	baseColor = texture2D(diffuseSampler, vDiffuseUV);

#ifdef ALPHATEST
	if (baseColor.a < 0.4)
		discard;
#endif

	baseColor.rgb *= vDiffuseInfos.y;
#endif

	// Bump
	vec3 normalW = vNormalW;

#ifdef BUMP
	normalW = perturbNormal(viewDirectionW);
#endif

	// Ambient color
	vec3 baseAmbientColor = vec3(1., 1., 1.);

#ifdef AMBIENT
	baseAmbientColor = texture2D(ambientSampler, vAmbientUV).rgb * vAmbientInfos.y;
#endif

	// Lighting
	vec3 diffuseBase = vec3(0., 0., 0.);
	vec3 specularBase = vec3(0., 0., 0.);
	float shadow = 1.;

#ifdef LIGHT0
#ifdef SPOTLIGHT0
	lightingInfo info = computeSpotLighting(viewDirectionW, normalW, vLightData0, vLightDirection0, vLightDiffuse0, vLightSpecular0);
#endif
#ifdef HEMILIGHT0
	lightingInfo info = computeHemisphericLighting(viewDirectionW, normalW, vLightData0, vLightDiffuse0, vLightSpecular0, vLightGround0);
#endif
#ifdef POINTDIRLIGHT0
	lightingInfo info = computeLighting(viewDirectionW, normalW, vLightData0, vLightDiffuse0, vLightSpecular0);
#endif
#ifdef SHADOW0
	#ifdef SHADOWVSM0
		shadow = computeShadowWithVSM(vPositionFromLight0, shadowSampler0);
	#else
		shadow = computeShadow(vPositionFromLight0, shadowSampler0);
	#endif
#else
	shadow = 1.;
#endif
	diffuseBase += info.diffuse * shadow;
	specularBase += info.specular * shadow;
#endif

#ifdef LIGHT1
#ifdef SPOTLIGHT1
	info = computeSpotLighting(viewDirectionW, normalW, vLightData1, vLightDirection1, vLightDiffuse1, vLightSpecular1);
#endif
#ifdef HEMILIGHT1
	info = computeHemisphericLighting(viewDirectionW, normalW, vLightData1, vLightDiffuse1, vLightSpecular1, vLightGround1);
#endif
#ifdef POINTDIRLIGHT1
	info = computeLighting(viewDirectionW, normalW, vLightData1, vLightDiffuse1, vLightSpecular1);
#endif
#ifdef SHADOW1
	#ifdef SHADOWVSM1
		shadow = computeShadowWithVSM(vPositionFromLight1, shadowSampler1);
	#else
		shadow = computeShadow(vPositionFromLight1, shadowSampler1);
	#endif
#else
	shadow = 1.;
#endif
	diffuseBase += info.diffuse * shadow;
	specularBase += info.specular * shadow;
#endif

#ifdef LIGHT2
#ifdef SPOTLIGHT2
	info = computeSpotLighting(viewDirectionW, normalW, vLightData2, vLightDirection2, vLightDiffuse2, vLightSpecular2);
#endif
#ifdef HEMILIGHT2
	info = computeHemisphericLighting(viewDirectionW, normalW, vLightData2, vLightDiffuse2, vLightSpecular2, vLightGround2);
#endif
#ifdef POINTDIRLIGHT2
	info = computeLighting(viewDirectionW, normalW, vLightData2, vLightDiffuse2, vLightSpecular2);
#endif
#ifdef SHADOW2
	#ifdef SHADOWVSM2
		shadow = computeShadowWithVSM(vPositionFromLight2, shadowSampler2);
	#else
		shadow = computeShadow(vPositionFromLight2, shadowSampler2);
	#endif	
#else
	shadow = 1.;
#endif
	diffuseBase += info.diffuse * shadow;
	specularBase += info.specular * shadow;
#endif

#ifdef LIGHT3
#ifdef SPOTLIGHT3
	info = computeSpotLighting(viewDirectionW, normalW, vLightData3, vLightDirection3, vLightDiffuse3, vLightSpecular3);
#endif
#ifdef HEMILIGHT3
	info = computeHemisphericLighting(viewDirectionW, normalW, vLightData3, vLightDiffuse3, vLightSpecular3, vLightGround3);
#endif
#ifdef POINTDIRLIGHT3
	info = computeLighting(viewDirectionW, normalW, vLightData3, vLightDiffuse3, vLightSpecular3);
#endif
#ifdef SHADOW3
	#ifdef SHADOWVSM3
		shadow = computeShadowWithVSM(vPositionFromLight3, shadowSampler3);
	#else
		shadow = computeShadow(vPositionFromLight3, shadowSampler3);
	#endif	
#else
	shadow = 1.;
#endif
	diffuseBase += info.diffuse * shadow;
	specularBase += info.specular * shadow;
#endif

	// Reflection
	vec3 reflectionColor = vec3(0., 0., 0.);

#ifdef REFLECTION
	if (vReflectionInfos.z != 0.0)
	{
		reflectionColor = textureCube(reflectionCubeSampler, vReflectionUVW).rgb * vReflectionInfos.y;
	}
	else
	{
		vec2 coords = vReflectionUVW.xy;

		if (vReflectionInfos.x == MAP_PROJECTION)
		{
			coords /= vReflectionUVW.z;
		}

		coords.y = 1.0 - coords.y;

		reflectionColor = texture2D(reflection2DSampler, coords).rgb * vReflectionInfos.y;
	}
#endif

	// Alpha
	float alpha = vDiffuseColor.a;

#ifdef OPACITY
	vec3 opacityMap = texture2D(opacitySampler, vOpacityUV).rgb * vec3(0.3, 0.59, 0.11);
	alpha *= (opacityMap.x + opacityMap.y + opacityMap.z)* vOpacityInfos.y;
#endif

	// Emissive
	vec3 emissiveColor = vEmissiveColor;
#ifdef EMISSIVE
	emissiveColor += texture2D(emissiveSampler, vEmissiveUV).rgb * vEmissiveInfos.y;
#endif

	// Specular map
	vec3 specularColor = vSpecularColor.rgb;
#ifdef SPECULAR
	specularColor = texture2D(specularSampler, vSpecularUV).rgb * vSpecularInfos.y;
#endif

	// Composition
	vec3 finalDiffuse = clamp(diffuseBase * diffuseColor + emissiveColor + vAmbientColor, 0.0, 1.0) * baseColor.rgb;
	vec3 finalSpecular = specularBase * specularColor;

	vec4 color = vec4(finalDiffuse * baseAmbientColor + finalSpecular + reflectionColor, alpha);

#ifdef FOG
	float fog = CalcFogFactor();
	color.rgb = fog * color.rgb + (1.0 - fog) * vFogColor;
#endif

	gl_FragColor = color;
}` 

ShadersStore["default_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

#define MAP_EXPLICIT	0.
#define MAP_SPHERICAL	1.
#define MAP_PLANAR		2.
#define MAP_CUBIC		3.
#define MAP_PROJECTION	4.
#define MAP_SKYBOX		5.

// Attributes
attribute vec3 position;
attribute vec3 normal;
#ifdef UV1
attribute vec2 uv;
#endif
#ifdef UV2
attribute vec2 uv2;
#endif
#ifdef VERTEXCOLOR
attribute vec3 color;
#endif

// Uniforms
uniform mat4 world;
uniform mat4 view;
uniform mat4 worldViewProjection;

#ifdef DIFFUSE
varying vec2 vDiffuseUV;
uniform mat4 diffuseMatrix;
uniform vec2 vDiffuseInfos;
#endif

#ifdef AMBIENT
varying vec2 vAmbientUV;
uniform mat4 ambientMatrix;
uniform vec2 vAmbientInfos;
#endif

#ifdef OPACITY
varying vec2 vOpacityUV;
uniform mat4 opacityMatrix;
uniform vec2 vOpacityInfos;
#endif

#ifdef REFLECTION
uniform vec3 vEyePosition;
varying vec3 vReflectionUVW;
uniform vec3 vReflectionInfos;
uniform mat4 reflectionMatrix;
#endif

#ifdef EMISSIVE
varying vec2 vEmissiveUV;
uniform vec2 vEmissiveInfos;
uniform mat4 emissiveMatrix;
#endif

#ifdef SPECULAR
varying vec2 vSpecularUV;
uniform vec2 vSpecularInfos;
uniform mat4 specularMatrix;
#endif

#ifdef BUMP
varying vec2 vBumpUV;
uniform vec2 vBumpInfos;
uniform mat4 bumpMatrix;
#endif

// Output
varying vec3 vPositionW;
varying vec3 vNormalW;

#ifdef VERTEXCOLOR
varying vec3 vColor;
#endif

#ifdef CLIPPLANE
uniform vec4 vClipPlane;
varying float fClipDistance;
#endif

#ifdef FOG
varying float fFogDistance;
#endif

#ifdef SHADOWS
#ifdef LIGHT0
uniform mat4 lightMatrix0;
varying vec4 vPositionFromLight0;
#endif
#ifdef LIGHT1
uniform mat4 lightMatrix1;
varying vec4 vPositionFromLight1;
#endif
#ifdef LIGHT2
uniform mat4 lightMatrix2;
varying vec4 vPositionFromLight2;
#endif
#ifdef LIGHT3
uniform mat4 lightMatrix3;
varying vec4 vPositionFromLight3;
#endif
#endif

#ifdef REFLECTION
vec3 computeReflectionCoords(float mode, vec4 worldPos, vec3 worldNormal)
{
	if (mode == MAP_SPHERICAL)
	{
		vec3 coords = vec3(view * vec4(worldNormal, 0.0));

		return vec3(reflectionMatrix * vec4(coords, 1.0));
	}
	else if (mode == MAP_PLANAR)
	{
		vec3 viewDir = worldPos.xyz - vEyePosition;
		vec3 coords = normalize(reflect(viewDir, worldNormal));

		return vec3(reflectionMatrix * vec4(coords, 1));
	}
	else if (mode == MAP_CUBIC)
	{
		vec3 viewDir = worldPos.xyz - vEyePosition;
		vec3 coords = reflect(viewDir, worldNormal);

		return vec3(reflectionMatrix * vec4(coords, 0));
	}
	else if (mode == MAP_PROJECTION)
	{
		return vec3(reflectionMatrix * (view * worldPos));
	}
	else if (mode == MAP_SKYBOX)
	{
		return position;
	}

	return vec3(0, 0, 0);
}
#endif

void main(void) {
	gl_Position = worldViewProjection * vec4(position, 1.0);

	vec4 worldPos = world * vec4(position, 1.0);
	vPositionW = vec3(worldPos);
	vNormalW = normalize(vec3(world * vec4(normal, 0.0)));

	// Texture coordinates
#ifndef UV1
	vec2 uv = vec2(0., 0.);
#endif
#ifndef UV2
	vec2 uv2 = vec2(0., 0.);
#endif

#ifdef DIFFUSE
	if (vDiffuseInfos.x == 0.)
	{
		vDiffuseUV = vec2(diffuseMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vDiffuseUV = vec2(diffuseMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef AMBIENT
	if (vAmbientInfos.x == 0.)
	{
		vAmbientUV = vec2(ambientMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vAmbientUV = vec2(ambientMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef OPACITY
	if (vOpacityInfos.x == 0.)
	{
		vOpacityUV = vec2(opacityMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vOpacityUV = vec2(opacityMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef REFLECTION
	vReflectionUVW = computeReflectionCoords(vReflectionInfos.x, vec4(vPositionW, 1.0), vNormalW);
#endif

#ifdef EMISSIVE
	if (vEmissiveInfos.x == 0.)
	{
		vEmissiveUV = vec2(emissiveMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vEmissiveUV = vec2(emissiveMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef SPECULAR
	if (vSpecularInfos.x == 0.)
	{
		vSpecularUV = vec2(specularMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vSpecularUV = vec2(specularMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef BUMP
	if (vBumpInfos.x == 0.)
	{
		vBumpUV = vec2(bumpMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vBumpUV = vec2(bumpMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

	// Clip plane
#ifdef CLIPPLANE
	fClipDistance = dot(worldPos, vClipPlane);
#endif

	// Fog
#ifdef FOG
	fFogDistance = (view * worldPos).z;
#endif

	// Shadows
#ifdef SHADOWS
#ifdef LIGHT0
	vPositionFromLight0 = lightMatrix0 * vec4(position, 1.0);
#endif
#ifdef LIGHT1
	vPositionFromLight1 = lightMatrix1 * vec4(position, 1.0);
#endif
#ifdef LIGHT2
	vPositionFromLight2 = lightMatrix2 * vec4(position, 1.0);
#endif
#ifdef LIGHT3
	vPositionFromLight3 = lightMatrix3 * vec4(position, 1.0);
#endif
#endif

	// Vertex color
#ifdef VERTEXCOLOR
	vColor = color;
#endif
}` 

ShadersStore["iedefault_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

#define MAP_PROJECTION	4.

// Constants
uniform vec3 vEyePosition;
uniform vec3 vAmbientColor;
uniform vec4 vDiffuseColor;
uniform vec4 vSpecularColor;
uniform vec3 vEmissiveColor;

// Lights
#ifdef LIGHT0
uniform vec4 vLightData0;
uniform vec3 vLightDiffuse0;
uniform vec3 vLightSpecular0;
#ifdef SHADOW0
varying vec4 vPositionFromLight0;
uniform sampler2D shadowSampler0;
#endif
#endif

//#ifdef LIGHT1
//uniform vec4 vLightData1;
//uniform vec3 vLightDiffuse1;
//uniform vec3 vLightSpecular1;
//#endif

//#ifdef LIGHT2
//uniform vec4 vLightData2;
//uniform vec3 vLightDiffuse2;
//uniform vec3 vLightSpecular2;
//#endif

// Samplers
#ifdef DIFFUSE
varying vec2 vDiffuseUV;
uniform sampler2D diffuseSampler;
uniform vec2 vDiffuseInfos;
#endif

#ifdef AMBIENT
varying vec2 vAmbientUV;
uniform sampler2D ambientSampler;
uniform vec2 vAmbientInfos;
#endif

#ifdef OPACITY	
varying vec2 vOpacityUV;
uniform sampler2D opacitySampler;
uniform vec2 vOpacityInfos;
#endif

#ifdef REFLECTION
varying vec3 vReflectionUVW;
uniform samplerCube reflectionCubeSampler;
uniform sampler2D reflection2DSampler;
uniform vec3 vReflectionInfos;
#endif

#ifdef EMISSIVE
varying vec2 vEmissiveUV;
uniform vec2 vEmissiveInfos;
uniform sampler2D emissiveSampler;
#endif

#ifdef SPECULAR
varying vec2 vSpecularUV;
uniform vec2 vSpecularInfos;
uniform sampler2D specularSampler;
#endif

// Input
varying vec3 vPositionW;
varying vec3 vNormalW;

#ifdef CLIPPLANE
varying float fClipDistance;
#endif

// Shadows
#ifdef SHADOWS

float unpack(vec4 color)
{
	const vec4 bitShift = vec4(1. / (255. * 255. * 255.), 1. / (255. * 255.), 1. / 255., 1.);
	return dot(color, bitShift);
}

float unpackHalf(vec2 color)
{
	return color.x + (color.y / 255.0);
}

// Thanks to http://devmaster.net/
float ChebychevInequality(vec2 moments, float t)
{
	if (t <= moments.x)
	{
		return 1.0;
	}

	float variance = moments.y - (moments.x * moments.x);
	variance = max(variance, 0);

	float d = t - moments.x;
	return variance / (variance + d * d);
}

#endif

// Fog
#ifdef FOG

#define FOGMODE_NONE    0.
#define FOGMODE_EXP     1.
#define FOGMODE_EXP2    2.
#define FOGMODE_LINEAR  3.
#define E 2.71828

uniform vec4 vFogInfos;
uniform vec3 vFogColor;
varying float fFogDistance;

float CalcFogFactor()
{
	float fogCoeff = 1.0;
	float fogStart = vFogInfos.y;
	float fogEnd = vFogInfos.z;
	float fogDensity = vFogInfos.w;

	if (FOGMODE_LINEAR == vFogInfos.x)
	{
		fogCoeff = (fogEnd - fFogDistance) / (fogEnd - fogStart);
	}
	else if (FOGMODE_EXP == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fogDensity);
	}
	else if (FOGMODE_EXP2 == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fFogDistance * fogDensity * fogDensity);
	}

	return min(1., max(0., fogCoeff));
}

#endif

vec3 computeDiffuseLighting(vec3 vNormal, vec4 lightData, vec3 diffuseColor) {
	vec3 lightVectorW;
	if (lightData.w == 0.)
	{
		lightVectorW = normalize(lightData.xyz - vPositionW);
	}
	else
	{
		lightVectorW = normalize(-lightData.xyz);
	}

	// diffuse
	float ndl = max(0., dot(vNormal, lightVectorW));

	return ndl * diffuseColor;
}

vec3 computeSpecularLighting(vec3 viewDirectionW, vec3 vNormal, vec4 lightData, vec3 specularColor) {
	vec3 lightVectorW;
	if (lightData.w == 0.)
	{
		lightVectorW = normalize(lightData.xyz - vPositionW);
	}
	else
	{
		lightVectorW = normalize(-lightData.xyz);
	}

	// Specular
	vec3 angleW = normalize(viewDirectionW + lightVectorW);
	float specComp = max(0., dot(vNormal, angleW));
	specComp = pow(specComp, vSpecularColor.a);

	return specComp * specularColor;
}

void main(void) {
	// Clip plane
#ifdef CLIPPLANE
	if (fClipDistance > 0.0)
		discard;
#endif

	vec3 viewDirectionW = normalize(vEyePosition - vPositionW);

	// Base color
	vec4 baseColor = vec4(1., 1., 1., 1.);
	vec3 diffuseColor = vDiffuseColor.rgb;

#ifdef DIFFUSE
	baseColor = texture2D(diffuseSampler, vDiffuseUV);

#ifdef ALPHATEST
	if (baseColor.a < 0.4)
		discard;
#endif

	baseColor.rgb *= vDiffuseInfos.y;
#endif

	// Bump
	vec3 normalW = vNormalW;

	// Ambient color
	vec3 baseAmbientColor = vec3(1., 1., 1.);

#ifdef AMBIENT
	baseAmbientColor = texture2D(ambientSampler, vAmbientUV).rgb * vAmbientInfos.y;
#endif

	// Lighting
	vec3 diffuseBase = vec3(0., 0., 0.);
	vec3 specularBase = vec3(0., 0., 0.);
	float shadow = 1.0;

#ifdef LIGHT0
	#ifdef SHADOW0
		vec3 depth = vPositionFromLight0.xyz / vPositionFromLight0.w;
		vec2 uv = 0.5 * depth.xy + vec2(0.5, 0.5);
	
		if (uv.x >= 0. && uv.x <= 1.0 && uv.y >= 0. && uv.y <= 1.0)
		{
		#ifdef SHADOWVSM0
			vec4 texel = texture2D(shadowSampler0, uv);

			vec2 moments = vec2(unpackHalf(texel.xy), unpackHalf(texel.zw));
			shadow = clamp(1.3 - ChebychevInequality(moments, depth.z), 0., 1.0);
		#else
			float shadowDepth = unpack(texture2D(shadowSampler0, uv));

			if (depth.z > shadowDepth)
			{
				shadow = 0.;
			}
		#endif
		}
	#endif
	diffuseBase += computeDiffuseLighting(normalW, vLightData0, vLightDiffuse0) * shadow;
	specularBase += computeSpecularLighting(viewDirectionW, normalW, vLightData0, vLightSpecular0) * shadow;
#endif
//#ifdef LIGHT1
//	diffuseBase += computeDiffuseLighting(normalW, vLightData1, vLightDiffuse1);
//	specularBase += computeSpecularLighting(viewDirectionW, normalW, vLightData1, vLightSpecular1);
//#endif
//#ifdef LIGHT2
//	diffuseBase += computeDiffuseLighting(normalW, vLightData2, vLightDiffuse2);
//	specularBase += computeSpecularLighting(viewDirectionW, normalW, vLightData2, vLightSpecular2);
//#endif


	// Reflection
	vec3 reflectionColor = vec3(0., 0., 0.);

#ifdef REFLECTION
	if (vReflectionInfos.z != 0.0)
	{
		reflectionColor = textureCube(reflectionCubeSampler, vReflectionUVW).rgb * vReflectionInfos.y;
	}
	else
	{
		vec2 coords = vReflectionUVW.xy;

		if (vReflectionInfos.x == MAP_PROJECTION)
		{
			coords /= vReflectionUVW.z;
		}

		coords.y = 1.0 - coords.y;

		reflectionColor = texture2D(reflection2DSampler, coords).rgb * vReflectionInfos.y;
	}	
#endif

	// Alpha
	float alpha = vDiffuseColor.a;

#ifdef OPACITY
	vec3 opacityMap = texture2D(opacitySampler, vOpacityUV).rgb * vec3(0.3, 0.59, 0.11);
	alpha *= (opacityMap.x + opacityMap.y + opacityMap.z )* vOpacityInfos.y;
#endif

	// Emissive
	vec3 emissiveColor = vEmissiveColor;
#ifdef EMISSIVE
	emissiveColor += texture2D(emissiveSampler, vEmissiveUV).rgb * vEmissiveInfos.y;
#endif

	// Specular map
	vec3 specularColor = vSpecularColor.rgb;
#ifdef SPECULAR
	specularColor = texture2D(specularSampler, vSpecularUV).rgb * vSpecularInfos.y;	
#endif

	// Composition
	vec3 finalDiffuse = clamp(diffuseBase * diffuseColor + emissiveColor + vAmbientColor, 0.0, 1.0) * baseColor.rgb;
	vec3 finalSpecular = specularBase * specularColor;

	vec4 color = vec4(finalDiffuse * baseAmbientColor + finalSpecular + reflectionColor, alpha);

#ifdef FOG
	float fog = CalcFogFactor();
	color.rgb = fog * color.rgb + (1.0 - fog) * vFogColor;
#endif

	gl_FragColor = color;
}` 

ShadersStore["iedefault_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

#define MAP_EXPLICIT	0.
#define MAP_SPHERICAL	1.
#define MAP_PLANAR		2.
#define MAP_CUBIC		3.
#define MAP_PROJECTION	4.
#define MAP_SKYBOX		5.

// Attributes
attribute vec3 position;
attribute vec3 normal;
#ifdef UV1
attribute vec2 uv;
#endif
#ifdef UV2
attribute vec2 uv2;
#endif

// Uniforms
uniform mat4 world;
uniform mat4 view;
uniform mat4 worldViewProjection;

#ifdef DIFFUSE
varying vec2 vDiffuseUV;
uniform mat4 diffuseMatrix;
uniform vec2 vDiffuseInfos;
#endif

#ifdef AMBIENT
varying vec2 vAmbientUV;
uniform mat4 ambientMatrix;
uniform vec2 vAmbientInfos;
#endif

#ifdef OPACITY
varying vec2 vOpacityUV;
uniform mat4 opacityMatrix;
uniform vec2 vOpacityInfos;
#endif

#ifdef REFLECTION
uniform vec3 vEyePosition;
varying vec3 vReflectionUVW;

uniform vec3 vReflectionInfos;
uniform mat4 reflectionMatrix;
#endif

#ifdef EMISSIVE
varying vec2 vEmissiveUV;
uniform vec2 vEmissiveInfos;
uniform mat4 emissiveMatrix;
#endif

#ifdef SPECULAR
varying vec2 vSpecularUV;
uniform vec2 vSpecularInfos;
uniform mat4 specularMatrix;
#endif

// Output
varying vec3 vPositionW;
varying vec3 vNormalW;

#ifdef CLIPPLANE
uniform vec4 vClipPlane;
varying float fClipDistance;
#endif

#ifdef FOG
varying float fFogDistance;
#endif

#ifdef SHADOWS
#ifdef LIGHT0
uniform mat4 lightMatrix0;
varying vec4 vPositionFromLight0;
#endif
#endif

#ifdef REFLECTION
vec3 computeReflectionCoords(float mode, vec4 worldPos, vec3 worldNormal)
{
	if (mode == MAP_SPHERICAL)
	{
		vec3 coords = vec3(view * vec4(worldNormal, 0.0));

		return vec3(reflectionMatrix * vec4(coords, 1.0));
	}
	else if (mode == MAP_PLANAR)
	{
		vec3 viewDir = worldPos.xyz - vEyePosition;
		vec3 coords = normalize(reflect(viewDir, worldNormal));

		return vec3(reflectionMatrix * vec4(coords, 1));
	}
	else if (mode == MAP_CUBIC)
	{
		vec3 viewDir = worldPos.xyz - vEyePosition;
		vec3 coords = reflect(viewDir, worldNormal);

		return vec3(reflectionMatrix * vec4(coords, 0));
	}
	else if (mode == MAP_PROJECTION)
	{
		return vec3(reflectionMatrix * (view * worldPos));
	}
	else if (mode == MAP_SKYBOX)
	{
		return position;
	}

	return vec3(0, 0, 0);
}
#endif

void main(void) {
	gl_Position = worldViewProjection * vec4(position, 1.0);

	vec4 worldPos = world * vec4(position, 1.0);
	vPositionW = vec3(worldPos);
	vNormalW = normalize(vec3(world * vec4(normal, 0.0)));

	// Texture coordinates
#ifndef UV1
	vec2 uv = vec2(0., 0.);
#endif
#ifndef UV2
	vec2 uv2 = vec2(0., 0.);
#endif

#ifdef DIFFUSE
	if (vDiffuseInfos.x == 0.)
	{
		vDiffuseUV = vec2(diffuseMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vDiffuseUV = vec2(diffuseMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef AMBIENT
	if (vAmbientInfos.x == 0.)
	{
		vAmbientUV = vec2(ambientMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vAmbientUV = vec2(ambientMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef OPACITY
	if (vOpacityInfos.x == 0.)
	{
		vOpacityUV = vec2(opacityMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vOpacityUV = vec2(opacityMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef REFLECTION
	vReflectionUVW = computeReflectionCoords(vReflectionInfos.x, vec4(vPositionW, 1.0), vNormalW);
#endif

#ifdef EMISSIVE
	if (vEmissiveInfos.x == 0.)
	{
		vEmissiveUV = vec2(emissiveMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vEmissiveUV = vec2(emissiveMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

#ifdef SPECULAR
	if (vSpecularInfos.x == 0.)
	{
		vSpecularUV = vec2(specularMatrix * vec4(uv, 1.0, 0.0));
	}
	else
	{
		vSpecularUV = vec2(specularMatrix * vec4(uv2, 1.0, 0.0));
	}
#endif

	// Clip plane
#ifdef CLIPPLANE
	fClipDistance = dot(worldPos, vClipPlane);
#endif

	// Fog
#ifdef FOG
	fFogDistance = (view * worldPos).z;
#endif

	// Shadows
#ifdef SHADOWS
#ifdef LIGHT0
	vPositionFromLight0 = lightMatrix0 * vec4(position, 1.0);
#endif
#endif
}` 

ShadersStore["layer_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

// Samplers
varying vec2 vUV;
uniform sampler2D textureSampler;

// Color
uniform vec4 color;

void main(void) {
	vec4 baseColor = texture2D(textureSampler, vUV);

	gl_FragColor = baseColor * color;
}` 

ShadersStore["layer_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

// Attributes
attribute vec2 position;

// Uniforms
uniform mat4 textureMatrix;

// Output
varying vec2 vUV;

const vec2 madd = vec2(0.5, 0.5);

void main(void) {	

	vUV = vec2(textureMatrix * vec4(position * madd + madd, 1.0, 0.0));
	gl_Position = vec4(position, 0.0, 1.0);
}` 

ShadersStore["particles_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

// Samplers
varying vec2 vUV;
varying vec4 vColor;
uniform vec4 textureMask;
uniform sampler2D diffuseSampler;

#ifdef CLIPPLANE
varying float fClipDistance;
#endif

void main(void) {
#ifdef CLIPPLANE
	if (fClipDistance > 0.0)
		discard;
#endif
	vec4 baseColor = texture2D(diffuseSampler, vUV);

	gl_FragColor = (baseColor * textureMask + (vec4(1., 1., 1., 1.) - textureMask)) * vColor;
}` 

ShadersStore["particles_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

// Attributes
attribute vec3 position;
attribute vec4 color;
attribute vec4 options;

// Uniforms
uniform mat4 view;
uniform mat4 projection;

// Output
varying vec2 vUV;
varying vec4 vColor;

#ifdef CLIPPLANE
uniform vec4 vClipPlane;
uniform mat4 invView;
varying float fClipDistance;
#endif

void main(void) {	
	vec3 viewPos = (view * vec4(position, 1.0)).xyz; 
	vec3 cornerPos;
	float size = options.y;
	float angle = options.x;
	vec2 offset = options.zw;

	cornerPos = vec3(offset.x - 0.5, offset.y  - 0.5, 0.) * size;

	// Rotate
	vec3 rotatedCorner;
	rotatedCorner.x = cornerPos.x * cos(angle) - cornerPos.y * sin(angle);
	rotatedCorner.y = cornerPos.x * sin(angle) + cornerPos.y * cos(angle);
	rotatedCorner.z = 0.;

	// Position
	viewPos += rotatedCorner;
	gl_Position = projection * vec4(viewPos, 1.0);   
	
	vColor = color;
	vUV = offset;

	// Clip plane
#ifdef CLIPPLANE
	vec4 worldPos = invView * vec4(viewPos, 1.0);
	fClipDistance = dot(worldPos, vClipPlane);
#endif
}` 

ShadersStore["shadowMap_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

vec4 pack(float depth)
{
	const vec4 bitOffset = vec4(255. * 255. * 255., 255. * 255., 255., 1.);
	const vec4 bitMask = vec4(0., 1. / 255., 1. / 255., 1. / 255.);
	
	vec4 comp = fract(depth * bitOffset);
	comp -= comp.xxyz * bitMask;
	
	return comp;
}

// Thanks to http://devmaster.net/
vec2 packHalf(float depth) 
{ 
	const vec2 bitOffset = vec2(1.0 / 255., 0.);
	vec2 color = vec2(depth, fract(depth * 255.));

	return color - (color.yy * bitOffset);
}


void main(void)
{
#ifdef VSM
	float moment1 = gl_FragCoord.z / gl_FragCoord.w; 
	float moment2 = moment1 * moment1;
	gl_FragColor = vec4(packHalf(moment1), packHalf(moment2));
#else
	gl_FragColor = pack(gl_FragCoord.z / gl_FragCoord.w);
#endif
}` 

ShadersStore["shadowMap_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

// Attribute
attribute vec3 position;

// Uniform
uniform mat4 worldViewProjection;

void main(void)
{
	gl_Position = worldViewProjection * vec4(position, 1.0);
}` 

ShadersStore["sprites_fragment"] = `#ifdef GL_ES
precision mediump float;
#endif

uniform bool alphaTest;

varying vec4 vColor;

// Samplers
varying vec2 vUV;
uniform sampler2D diffuseSampler;

// Fog
#ifdef FOG

#define FOGMODE_NONE    0.
#define FOGMODE_EXP     1.
#define FOGMODE_EXP2    2.
#define FOGMODE_LINEAR  3.
#define E 2.71828

uniform vec4 vFogInfos;
uniform vec3 vFogColor;
varying float fFogDistance;

float CalcFogFactor()
{
	float fogCoeff = 1.0;
	float fogStart = vFogInfos.y;
	float fogEnd = vFogInfos.z;
	float fogDensity = vFogInfos.w;

	if (FOGMODE_LINEAR == vFogInfos.x)
	{
		fogCoeff = (fogEnd - fFogDistance) / (fogEnd - fogStart);
	}
	else if (FOGMODE_EXP == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fogDensity);
	}
	else if (FOGMODE_EXP2 == vFogInfos.x)
	{
		fogCoeff = 1.0 / pow(E, fFogDistance * fFogDistance * fogDensity * fogDensity);
	}

	return min(1., max(0., fogCoeff));
}
#endif


void main(void) {
	vec4 baseColor = texture2D(diffuseSampler, vUV);

	if (alphaTest) 
	{
		if (baseColor.a < 0.95)
			discard;
	}

	baseColor *= vColor;

#ifdef FOG
	float fog = CalcFogFactor();
	baseColor.rgb = fog * baseColor.rgb + (1.0 - fog) * vFogColor;
#endif

	gl_FragColor = baseColor;
}` 

ShadersStore["sprites_vertex"] = `#ifdef GL_ES
precision mediump float;
#endif

// Attributes
attribute vec3 position;
attribute vec4 options;
attribute vec4 cellInfo;
attribute vec4 color;

// Uniforms
uniform vec2 textureInfos;
uniform mat4 view;
uniform mat4 projection;

// Output
varying vec2 vUV;
varying vec4 vColor;

#ifdef FOG
varying float fFogDistance;
#endif

void main(void) {	
	vec3 viewPos = (view * vec4(position, 1.0)).xyz; 
	vec3 cornerPos;
	
	float angle = options.x;
	float size = options.y;
	vec2 offset = options.zw;
	vec2 uvScale = textureInfos.xy;

	cornerPos = vec3(offset.x - 0.5, offset.y  - 0.5, 0.) * size;

	// Rotate
	vec3 rotatedCorner;
	rotatedCorner.x = cornerPos.x * cos(angle) - cornerPos.y * sin(angle);
	rotatedCorner.y = cornerPos.x * sin(angle) + cornerPos.y * cos(angle);
	rotatedCorner.z = 0.;

	// Position
	viewPos += rotatedCorner;
	gl_Position = projection * vec4(viewPos, 1.0);   

	// Color
	vColor = color;
	
	// Texture
	vec2 uvOffset = vec2(abs(offset.x - cellInfo.x), 1.0 - abs(offset.y - cellInfo.y));

	vUV = (uvOffset + cellInfo.zw) * uvScale;

	// Fog
#ifdef FOG
	fFogDistance = viewPos.z;
#endif
}` 

} 

