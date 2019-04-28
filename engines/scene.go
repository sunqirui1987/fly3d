package engines

import (
	"math"

	"github.com/suiqirui1987/fly3d/core"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/tools"
	log "github.com/suiqirui1987/fly3d/tools/logrus"
)

type ISceneDisposed interface {
	Dispose()
}

type Scene struct {

	//
	_engine      *Engine
	AutoClear    bool
	ClearColor   *math32.Color3
	AmbientColor *math32.Color3

	//maxtrix
	_viewMatrix       *math32.Matrix4
	_projectionMatrix *math32.Matrix4
	_transformMatrix  *math32.Matrix4
	_frustumPlanes    []*math32.Plane

	_totalVertices   int
	_activeVertices  int
	_activeParticles int
	_renderId        int

	_lastFrameDuration            int
	_evaluateActiveMeshesDuration int
	_renderTargetsDuration        int
	_renderDuration               int
	_particlesDuration            int
	_spritesDuration              int

	_animationRatio float32
	_pendingData    []string

	//callback
	BeforeRender func()
	AfterRender  func()

	_onReadyCallbacks        []func()
	_onBeforeRenderCallbacks []func()

	//disposed
	ToBeDisposed []ISceneDisposed

	// Fog
	FogMode    int
	FogColor   *math32.Color3
	FogDensity float32
	FogStart   float32
	FogEnd     float32

	//lights
	Lights []ILight

	//camera
	Cameras      []ICamera
	ActiveCamera ICamera

	//meshes
	Meshes                []IMesh
	_activeMeshes         []IMesh
	_opaqueSubMeshes      []ISubMesh
	_transparentSubMeshes []ISubMesh
	_alphaTestSubMeshes   []ISubMesh

	//target
	_renderTargets []ITexture

	//Materials
	Materials           []IMaterial
	MultiMaterials      []IMultiMaterial
	_processedMaterials []IMaterial

	// Textures
	Textures []ITexture

	// Particles
	ParticlesEnabled       bool
	ParticleSystems        []IParticleSystem
	_activeParticleSystems []IParticleSystem

	//Sprites
	SpriteManagers []ISpriteManager

	// Layers
	Layers []ILayer

	// Collisions
	CollisionsEnabled bool
	Gravity           *math32.Vector3

	// Animations
	ActiveAnimatables []IAnimatable
}

func NewScene(engine *Engine) *Scene {
	this := &Scene{}

	this._engine = engine
	this.AutoClear = true
	this.ClearColor = math32.NewColor3(0.2, 0.2, 0.3)
	this.AmbientColor = math32.NewColor3(0.0, 0.0, 0.0)

	this._engine.Scenes = append(this._engine.Scenes, this)

	this._totalVertices = 0
	this._activeVertices = 0
	this._activeParticles = 0
	this._lastFrameDuration = 0
	this._evaluateActiveMeshesDuration = 0
	this._renderTargetsDuration = 0
	this._renderDuration = 0

	this._renderId = 0

	this.ToBeDisposed = make([]ISceneDisposed, 0)

	this._onReadyCallbacks = make([]func(), 0)
	this._pendingData = make([]string, 0)

	this._onBeforeRenderCallbacks = make([]func(), 0)

	//fog
	this.FogMode = core.FOGMODE_NONE
	this.FogColor = math32.NewColor3(0.2, 0.2, 0.3)
	this.FogDensity = 0.1
	this.FogStart = 0
	this.FogEnd = 1000.0

	//Lights
	this.Lights = make([]ILight, 0)

	//Camera
	this.Cameras = make([]ICamera, 0)
	this.ActiveCamera = nil

	//Mesh
	this.Meshes = make([]IMesh, 0)
	this._activeMeshes = make([]IMesh, 0)

	// Materials
	this.Materials = make([]IMaterial, 0)
	this.MultiMaterials = make([]IMultiMaterial, 0)

	// Textures
	this.Textures = make([]ITexture, 0)

	// Particles
	this.ParticlesEnabled = true
	this.ParticleSystems = make([]IParticleSystem, 0)

	// Sprites
	this.SpriteManagers = make([]ISpriteManager, 0)

	// Layers
	this.Layers = make([]ILayer, 0)

	// Collisions
	this.CollisionsEnabled = true
	this.Gravity = math32.NewVector3(0, 0, -9)

	// Animations
	this.ActiveAnimatables = make([]IAnimatable, 0)

	// Matrices
	this._transformMatrix = math32.NewMatrix4().Zero()

	return this
}

//engine
func (this *Scene) GetEngine() *Engine {
	return this._engine
}

//stats
func (this *Scene) GetTotalVertices() int {
	return this._totalVertices
}
func (this *Scene) GetActiveVertices() int {
	return this._activeVertices
}

func (this *Scene) GetActiveParticles() int {
	return this._activeParticles
}
func (this *Scene) GetLastFrameDuration() int {
	return this._lastFrameDuration
}

func (this *Scene) GetEvaluateActiveMeshesDuration() int {
	return this._evaluateActiveMeshesDuration
}
func (this *Scene) GetRenderTargetsDuration() int {
	return this._renderTargetsDuration
}
func (this *Scene) GetRenderDuration() int {
	return this._renderDuration
}
func (this *Scene) GetParticlesDuration() int {
	return this._particlesDuration
}
func (this *Scene) GetSpritesDuration() int {
	return this._spritesDuration
}
func (this *Scene) GetAnimationRatio() float32 {
	return this._animationRatio
}

//ready

func (this *Scene) IsReady() bool {

	for index := 0; index < len(this.Materials); index++ {
		if !this.Materials[index].IsReady(nil) {
			return false
		}
	}

	return true
}
func (this *Scene) GetWaitingItemsCount() int {
	return len(this._pendingData)
}

func (this *Scene) ExecuteWhenReady(f func()) {
	if this.IsReady() {
		f()
		return
	}

	if len(this._pendingData) == 0 {
		f()
		return
	}
	this._onReadyCallbacks = append(this._onReadyCallbacks, f)

}

func (this *Scene) RegisterBeforeRender(f func()) {
	this._onBeforeRenderCallbacks = append(this._onBeforeRenderCallbacks, f)
}
func (this *Scene) UnregisterBeforeRender(f func()) {
	index := tools.IndexOf(f, this._onBeforeRenderCallbacks)
	if index > -1 {
		this._onBeforeRenderCallbacks = append(this._onBeforeRenderCallbacks[:index], this._onBeforeRenderCallbacks[index+1:]...)
	}

}
func (this *Scene) AddPendingData(url string) {
	this._pendingData = append(this._pendingData, url)
}

func (this *Scene) RemovePendingData(url string) {
	index := tools.IndexOf(url, this._pendingData)
	if index > -1 {
		this._pendingData = append(this._pendingData[:index], this._pendingData[index+1:]...)

		if len(this._pendingData) == 0 {
			//callback
			for _, f := range this._onReadyCallbacks {
				f()
			}

			this._onReadyCallbacks = make([]func(), 0)
		}
	}
}

func (this *Scene) _animate() {
	for index := 0; index < len(this.ActiveAnimatables); index++ {
		if !this.ActiveAnimatables[index].Animate() {

			this.ActiveAnimatables = append(this.ActiveAnimatables[:index], this.ActiveAnimatables[index+1:]...)
			index--
		}
	}
}

// Matrix
func (this *Scene) GetViewMatrix() *math32.Matrix4 {
	return this._viewMatrix
}
func (this *Scene) GetProjectionMatrix() *math32.Matrix4 {
	return this._projectionMatrix
}
func (this *Scene) GetTransformMatrix() *math32.Matrix4 {
	return this._transformMatrix
}
func (this *Scene) SetTransformMatrix(view *math32.Matrix4, projection *math32.Matrix4) {
	this._viewMatrix = view
	this._projectionMatrix = projection
	this._transformMatrix = this._viewMatrix.Multiply(this._projectionMatrix)
}

// Methods
func (this *Scene) ActiveCameraByID(id string) {
	for index := 0; index < len(this.Cameras); index++ {
		if this.Cameras[index].GetId() == id {
			this.ActiveCamera = this.Cameras[index]
			return
		}
	}
}

func (this *Scene) GetMaterialByID(id string) IMaterial {
	for index := 0; index < len(this.Materials); index++ {
		if this.Materials[index].GetId() == id {
			return this.Materials[index]
		}
	}
	return nil
}

func (this *Scene) GetMeshByID(id string) IMesh {
	for index := 0; index < len(this.Meshes); index++ {
		if this.Meshes[index].GetId() == id {
			return this.Meshes[index]
		}
	}
	return nil
}

func (this *Scene) GetLastMeshByID(id string) IMesh {
	var result IMesh

	for index := 0; index < len(this.Meshes); index++ {
		if this.Meshes[index].GetId() == id {
			result = this.Meshes[index]
		}
	}
	return result
}

func (this *Scene) GetMeshByName(name string) IMesh {
	for index := 0; index < len(this.Meshes); index++ {
		if this.Meshes[index].GetName() == name {
			return this.Meshes[index]
		}
	}
	return nil
}

func (this *Scene) IsActiveMesh(mesh IMesh) bool {

	index := tools.IndexOf(mesh, this.Meshes)
	return index != -1
}

func (this *Scene) _evaluateSubMesh(subMesh ISubMesh, mesh IMesh) {
	if len(mesh.GetSubMeshes()) == 1 || subMesh.IsInFrustrum(this._frustumPlanes) {
		material := subMesh.GetMaterial()

		if material != nil {
			// Render targets
			rendertargets := material.GetRenderTargetTextures()
			if rendertargets != nil {
				if tools.IndexOf(material, this._processedMaterials) == -1 {
					this._processedMaterials = append(this._processedMaterials, material)
					this._renderTargets = append(this._renderTargets, rendertargets...)
				}

			}

			// Dispatch
			if material.NeedAlphaBlending() || mesh.GetVisibility() < 1.0 { // Transparent
				if material.GetAlpha() > 0 || mesh.GetVisibility() < 1.0 {
					this._transparentSubMeshes = append(this._transparentSubMeshes, subMesh) // Opaque
				}
			} else if material.NeedAlphaBlending() { // Alpha test
				this._alphaTestSubMeshes = append(this._alphaTestSubMeshes, subMesh)
			} else {
				this._opaqueSubMeshes = append(this._opaqueSubMeshes, subMesh)
			}
		}
	}
}

func (this *Scene) _evaluateActiveMeshes() {

	this._activeMeshes = make([]IMesh, 0)
	this._opaqueSubMeshes = make([]ISubMesh, 0)
	this._transparentSubMeshes = make([]ISubMesh, 0)
	this._alphaTestSubMeshes = make([]ISubMesh, 0)
	this._processedMaterials = make([]IMaterial, 0)
	this._renderTargets = make([]ITexture, 0)
	this._activeParticleSystems = make([]IParticleSystem, 0)

	if this._frustumPlanes == nil || len(this._frustumPlanes) == 0 {
		this._frustumPlanes = math32.NewFrustum().GetPlanes(this._transformMatrix)
	} else {
		math32.NewFrustum().GetPlanesToRef(this._transformMatrix, this._frustumPlanes)
	}

	this._totalVertices = 0
	this._activeVertices = 0

	for meshIndex := 0; meshIndex < len(this.Meshes); meshIndex++ {
		mesh := this.Meshes[meshIndex]
		this._totalVertices += mesh.GetTotalVertices()

		if !mesh.IsReady() {
			continue
		}

		mesh.ComputeWorldMatrix()

		if mesh.IsEnabled() && mesh.IsVisible() && mesh.GetVisibility() > 0.0 && mesh.IsInFrustrum(this._frustumPlanes) {
			this._activeMeshes = append(this._activeMeshes, mesh)

			for _, subMesh := range mesh.GetSubMeshes() {
				this._evaluateSubMesh(subMesh, mesh)
			}

		}
	}

	// Particle systems
	beforeParticlesDate := tools.GetCurrentTimeMs()
	if this.ParticlesEnabled {
		for particleIndex := 0; particleIndex < len(this.ParticleSystems); particleIndex++ {
			particleSystem := this.ParticleSystems[particleIndex]
			emitter := particleSystem.GetEmitter()
			if (emitter != nil && emitter.GetPosition() == nil) ||
				(emitter != nil && emitter.IsEnabled() == true) {
				this._activeParticleSystems = append(this._activeParticleSystems, particleSystem)
				particleSystem.Animate()
			}
		}
	}
	this._particlesDuration = tools.GetCurrentTimeMs() - beforeParticlesDate

}

func (this *Scene) LocalRender(opaqueSubMeshes []ISubMesh, alphaTestSubMeshes []ISubMesh, transparentSubMeshes []ISubMesh, activeMeshes []IMesh) {
	engine := this._engine
	// Opaque
	var subIndex int
	var submesh ISubMesh
	for subIndex = 0; subIndex < len(opaqueSubMeshes); subIndex++ {
		submesh = opaqueSubMeshes[subIndex]
		this._activeVertices += submesh.GetVerticesCount()

		submesh.Render()
	}

	// Alpha test
	engine.SetAlphaTesting(true)
	for subIndex := 0; subIndex < len(alphaTestSubMeshes); subIndex++ {
		submesh = alphaTestSubMeshes[subIndex]
		this._activeVertices += submesh.GetVerticesCount()

		submesh.Render()
	}
	engine.SetAlphaTesting(false)

	if activeMeshes == nil {
		// Sprites
		beforeSpritessDate := tools.GetCurrentTimeMs()
		for index := 0; index < len(this.SpriteManagers); index++ {
			spriteManager := this.SpriteManagers[index]

			spriteManager.Render()
		}
		this._spritesDuration = tools.GetCurrentTimeMs() - beforeSpritessDate
	}

	// Transparent
	engine.SetAlphaMode(core.ALPHA_COMBINE)
	for subIndex = 0; subIndex < len(transparentSubMeshes); subIndex++ {
		submesh = transparentSubMeshes[subIndex]
		this._activeVertices += submesh.GetVerticesCount()

		submesh.Render()
	}
	engine.SetAlphaMode(core.ALPHA_DISABLE)

	// Particle systems
	beforeParticlesDate := tools.GetCurrentTimeMs()
	for particleIndex := 0; particleIndex < len(this._activeParticleSystems); particleIndex++ {
		particleSystem := this._activeParticleSystems[particleIndex]
		emitter := particleSystem.GetEmitter()
		if (emitter != nil && emitter.GetPosition() == nil) ||
			(activeMeshes == nil) ||
			(tools.IndexOf(emitter, activeMeshes) != -1) {
			this._activeParticles += particleSystem.Render()
		}
	}
	this._particlesDuration = tools.GetCurrentTimeMs() - beforeParticlesDate

}

func (this *Scene) Render() {

	// Camera
	if this.ActiveCamera == nil {
		log.Println("Active camera not set")
		return
	}

	startDate := tools.GetCurrentTimeMs()
	this._particlesDuration = 0
	this._activeParticles = 0
	engine := this._engine

	// Before render
	if this.BeforeRender != nil {
		this.BeforeRender()
	}

	for callbackIndex := 0; callbackIndex < len(this._onBeforeRenderCallbacks); callbackIndex++ {
		this._onBeforeRenderCallbacks[callbackIndex]()
	}

	this.SetTransformMatrix(this.ActiveCamera.GetViewMatrix(), this.ActiveCamera.GetProjectionMatrix())

	// Animations
	this._animationRatio = tools.GetDeltaTime() * (60.0 / 1000.0)
	this._animate()

	// Meshes
	beforeEvaluateActiveMeshesDate := tools.GetCurrentTimeMs()
	this._evaluateActiveMeshes()
	this._evaluateActiveMeshesDuration = tools.GetCurrentTimeMs() - beforeEvaluateActiveMeshesDate

	// Shadows
	for lightIndex := 0; lightIndex < len(this.Lights); lightIndex++ {
		light := this.Lights[lightIndex]
		shadowGenerator := light.GetShadowGenerator()

		if light.IsEnabled() && shadowGenerator != nil && shadowGenerator.IsReady() {
			this._renderTargets = append(this._renderTargets, shadowGenerator.GetShadowMap())
		}
	}

	// Render targets
	beforeRenderTargetDate := tools.GetCurrentTimeMs()
	for renderIndex := 0; renderIndex < len(this._renderTargets); renderIndex++ {
		renderTarget := this._renderTargets[renderIndex]

		renderTarget.Render()
	}

	if len(this._renderTargets) > 0 { // Restore back buffer
		engine.RestoreDefaultFramebuffer()
	}
	this._renderTargetsDuration = tools.GetCurrentTimeMs() - beforeRenderTargetDate

	var layerIndex int
	var layer ILayer

	// Clear
	beforeRenderDate := tools.GetCurrentTimeMs()
	engine.Clear(this.ClearColor, this.AutoClear, true)

	// Backgrounds
	if len(this.Layers) > 0 {
		engine.SetDepthBuffer(false)
		for layerIndex = 0; layerIndex < len(this.Layers); layerIndex++ {
			layer = this.Layers[layerIndex]
			if layer.IsBackground() {
				layer.Render()
			}
		}
		engine.SetDepthBuffer(true)
	}

	// Render
	this.LocalRender(this._opaqueSubMeshes, this._alphaTestSubMeshes, this._transparentSubMeshes, nil)

	// Foregrounds
	if len(this.Layers) > 0 {
		engine.SetDepthBuffer(false)
		for layerIndex = 0; layerIndex < len(this.Layers); layerIndex++ {
			layer = this.Layers[layerIndex]
			if !layer.IsBackground() {
				layer.Render()
			}
		}
		engine.SetDepthBuffer(true)
	}

	this._renderDuration = tools.GetCurrentTimeMs() - beforeRenderDate

	// Update camera
	this.ActiveCamera.Update()

	// After render
	if this.AfterRender != nil {
		this.AfterRender()
	}

	// Cleaning
	for index := 0; index < len(this.ToBeDisposed); index++ {
		this.ToBeDisposed[index].Dispose()
	}

	this.ToBeDisposed = make([]ISceneDisposed, 0)

	this._lastFrameDuration = tools.GetCurrentTimeMs() - startDate
}

func (this *Scene) Dispose() {

	this.BeforeRender = nil
	this.AfterRender = nil

	log.Println("Detach cameras")
	// Detach cameras
	canvas := this._engine.GetRenderingCanvas()
	var index int
	for index = 0; index < len(this.Cameras); index++ {
		this.Cameras[index].DetachControl(canvas)
	}

	log.Println("Release meshes")
	// Release meshes
	for _, m := range this.Meshes {
		m.Dispose()
	}
	this.Meshes = make([]IMesh, 0)

	log.Println("Release materials")
	// Release materials
	for _, m := range this.Materials {
		m.Dispose()
	}
	this.Materials = make([]IMaterial, 0)

	log.Println("Release particles")
	// Release particles
	for _, m := range this.ParticleSystems {
		m.Dispose()
	}
	this.ParticleSystems = make([]IParticleSystem, 0)

	log.Println("Release sprites")
	// Release sprites
	for _, m := range this.SpriteManagers {
		m.Dispose()
	}
	this.SpriteManagers = make([]ISpriteManager, 0)

	log.Println("Release layers")
	// Release layers
	for _, m := range this.Layers {
		m.Dispose()
	}
	this.Layers = make([]ILayer, 0)

	log.Println("Release textures")
	// Release textures
	for _, m := range this.Textures {
		m.Dispose()
	}
	this.Textures = make([]ITexture, 0)

	log.Println("Remove from engine")
	// Remove from engine
	index = tools.IndexOf(this, this._engine.Scenes)
	if index > -1 {
		this._engine.Scenes = append(this._engine.Scenes[:index], this._engine.Scenes[index+1:]...)
	}

	log.Println("engine.WipeCaches")
	this._engine.WipeCaches()
}

func (this *Scene) CreatePickingRay(x, y float32, world *math32.Matrix4) *math32.Ray {
	engine := this._engine

	if this._viewMatrix != nil {
		if this.ActiveCamera != nil {
			log.Println("Active camera not set")
			return nil
		}

		this.SetTransformMatrix(this.ActiveCamera.GetViewMatrix(), this.ActiveCamera.GetProjectionMatrix())
	}

	rw := float32(engine.GetRenderWidth()) * engine.GetHardwareScalingLevel()
	rh := float32(engine.GetRenderHeight()) * engine.GetHardwareScalingLevel()
	if world == nil {
		world = (math32.NewMatrix4()).Identity()
	}
	return math32.CreateNew(x, y, rw, rh, world, this._viewMatrix, this._projectionMatrix)
}

func (this *Scene) Pick(x, y float32) *ColliderPickingInfo {

	var distance float32
	var pickedPoint *math32.Vector3
	var pickedMesh IMesh

	distance = math.MaxFloat32

	for meshIndex := 0; meshIndex < len(this.Meshes); meshIndex++ {
		mesh := this.Meshes[meshIndex]
		if !mesh.IsEnabled() || !mesh.IsVisible() || !mesh.IsPickable() {
			continue
		}

		world := mesh.GetWorldMatrix()
		ray := this.CreatePickingRay(x, y, world)

		result := mesh.Intersects(ray)
		if !result.Hit {
			continue
		}
		if result.Distance >= distance {
			continue
		}

		distance = result.Distance
		pickedMesh = mesh

		// Get picked point
		worldOrigin := ray.Origin.TransformCoordinates(world)
		direction := ray.Direction.Clone()
		direction.Normalize()
		direction = direction.Scale(result.Distance)

		worldDirection := direction.TransformNormal(world)
		pickedPoint = worldOrigin.Add(worldDirection)
	}

	pickinfo := &ColliderPickingInfo{
		Hit:         (distance != math.MaxFloat32),
		Distance:    distance,
		PickedMesh:  pickedMesh,
		PickedPoint: pickedPoint,
	}

	return pickinfo

}
