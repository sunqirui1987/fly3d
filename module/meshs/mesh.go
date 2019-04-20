package meshs

import (
	"math"
	"reflect"

	log "github.com/suiqirui1987/fly3d/tools/logrus"
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/cullings"
	"github.com/suiqirui1987/fly3d/tools"
)

type MeshCache struct {
	position *math32.Vector3
	scaling  *math32.Vector3
	rotation *math32.Vector3
}

type Mesh struct {
	Name   string
	Id     string
	_scene *engines.Scene

	_totalVertices int
	_worldMatrix   *math32.Matrix4

	Position     *math32.Vector3
	Rotation     *math32.Vector3
	Scaling      *math32.Vector3
	_scaleFactor float32

	_vertexStrideSize int

	_indices      []uint16
	SubMeshes     []*SubMesh
	_childrenFlag bool

	Material      IMaterial
	MutilMaterial IMultiMaterial

	Parent          *Mesh
	_isReady        bool
	_isEnabled      bool
	Isvisible       bool
	Ispickable      bool
	Visibility      float32
	BillboardMode   int
	Checkcollisions bool
	OnDispose       func()
	_isDisposed     bool
	ReceiveShadows  bool

	_boundingInfo *cullings.BoundingInfo

	_animationStarted bool
	_vertexBuffers    map[string]*VertexBuffer
	_indexBuffer      *gl.GLIndexBuffer

	_collisionsScalingMatrix   *math32.Matrix4
	_collisionsTransformMatrix *math32.Matrix4

	//cache
	_cache           *MeshCache
	_cache_positions []*math32.Vector3
}

func NewMesh(name string, scene *engines.Scene) *Mesh {
	this := &Mesh{}
	this.Name = name
	this.Id = name
	this._scene = scene

	this.Init()

	this._scene.Meshes = append(this._scene.Meshes, this)

	return this

}

func (this *Mesh) Init() {

	this._totalVertices = 0
	this._worldMatrix = (math32.NewMatrix4()).Identity()

	this.Position = math32.NewVector3(0, 0, 0)
	this.Rotation = math32.NewVector3(0, 0, 0)
	this.Scaling = math32.NewVector3(1, 1, 1)

	this._indices = make([]uint16, 0)
	this.SubMeshes = make([]*SubMesh, 0)
	this._childrenFlag = false

	this.Material = nil
	this.MutilMaterial = nil

	this.Parent = nil
	this._isReady = true
	this._isEnabled = true
	this.Isvisible = true
	this.Ispickable = true
	this.Visibility = 1.0
	this.BillboardMode = BILLBOARDMODE_NONE
	this.Checkcollisions = false
	this._isDisposed = false

	this._boundingInfo = nil

	this._animationStarted = false

	//cache
	this._cache_positions = nil
	this._cache = &MeshCache{}
}

func (this *Mesh) GetBoundingInfo() *cullings.BoundingInfo {
	return this._boundingInfo
}
func (this *Mesh) GetScene() *engines.Scene {
	return this._scene
}

func (this *Mesh) GetVerticesData(kind string) []float32 {
	return this._vertexBuffers[kind].GetData()
}

func (this *Mesh) GetTotalIndices() int {
	return len(this._indices)
}
func (this *Mesh) GetIndices() []uint16 {
	return this._indices
}
func (this *Mesh) GetVertexStrideSize() int {
	return this._vertexStrideSize
}

func (this *Mesh) _needToSynchonizeChildren() bool {
	return this._childrenFlag
}

func (this *Mesh) IsSynchronized() bool {
	if this.BillboardMode != BILLBOARDMODE_NONE {
		return false
	}

	if this._cache.position == nil || this._cache.rotation == nil || this._cache.scaling == nil {
		return false
	}

	if !this._cache.position.Equals(this.Position) {
		return false
	}
	if !this._cache.rotation.Equals(this.Rotation) {
		return false
	}
	if !this._cache.scaling.Equals(this.Scaling) {
		return false
	}
	if this.Parent != nil {
		return !this.Parent._needToSynchonizeChildren()
	}
	return true

}

func (this *Mesh) SetEnabled(val bool) {
	this._isEnabled = val
}

func (this *Mesh) IsAnimated() bool {
	return this._animationStarted
}
func (this *Mesh) IsDisposed() bool {
	return this._isDisposed
}

func (this *Mesh) _computeWorldMatrix() *math32.Matrix4 {
	if this.IsSynchronized() {
		this._childrenFlag = false
		return this._worldMatrix
	}

	var localWorld *math32.Matrix4
	localWorld = math32.NewMatrix4().Zero()

	this._childrenFlag = true
	this._cache.position = this.Position.Clone()
	this._cache.rotation = this.Rotation.Clone()
	this._cache.scaling = this.Scaling.Clone()

	localScaling := math32.NewMatrix4().RotationYawPitchRoll(this.Rotation.Y, this.Rotation.X, this.Rotation.Z)
	localRotation := math32.NewMatrix4().Scaling(this.Scaling.X, this.Scaling.Y, this.Scaling.Z)

	localScalingRotation := localScaling.Multiply(localRotation)

	// Billboarding
	localTranslation := math32.NewMatrix4().Translation(this.Position.X, this.Position.Y, this.Position.Z)
	if this.BillboardMode != BILLBOARDMODE_NONE {
		localPosition := this.Position.Clone()
		zero := this._scene.ActiveCamera.GetPosition().Clone()

		if this.Parent != nil {
			localPosition = localPosition.Add(this.Parent.Position)
			localTranslation = math32.NewMatrix4().Translation(localPosition.X, localPosition.Y, localPosition.Z)
		}

		if this.BillboardMode&BILLBOARDMODE_ALL == BILLBOARDMODE_ALL {
			zero = this._scene.ActiveCamera.GetPosition()
		} else {
			if this.BillboardMode&BILLBOARDMODE_X > 0 {
				zero.X = localPosition.X + core.Epsilon
			}
			if this.BillboardMode&BILLBOARDMODE_Y > 0 {
				zero.Y = localPosition.Y + core.Epsilon
			}
			if this.BillboardMode&BILLBOARDMODE_Z > 0 {
				zero.Z = localPosition.Z + core.Epsilon
			}
		}

		localBillboard := math32.NewMatrix4().LookAtLH(localPosition, zero, math32.NewVector3Up())
		localBillboard[12] = 0.0
		localBillboard[13] = 0.0
		localBillboard[14] = 0.0

		localBillboard.Invert()

		localWorld = localScalingRotation.Multiply(localBillboard)

		localScalingRotation = math32.NewMatrix4().RotationY(math.Pi).Multiply(localWorld)
	}

	// Parent
	if this.Parent != nil && this.BillboardMode == BILLBOARDMODE_NONE {
		localWorld = localScalingRotation.Multiply(localTranslation)
		tmp_parentWorld := this.Parent.GetWorldMatrix()
		this._worldMatrix = localWorld.Multiply(tmp_parentWorld)
	} else {
		this._worldMatrix = localScalingRotation.Multiply(localTranslation)
	}

	// Bounding info
	if this._boundingInfo != nil {
		this._scaleFactor = math32.Max(this.Scaling.X, this.Scaling.Y)
		this._scaleFactor = math32.Max(this._scaleFactor, this.Scaling.Z)

		if this.Parent != nil {
			this._scaleFactor = this._scaleFactor * this.Parent._scaleFactor

		}

		this._boundingInfo.Update(localWorld, this._scaleFactor)

		for subIndex := 0; subIndex < len(this.SubMeshes); subIndex++ {
			subMesh := this.SubMeshes[subIndex]

			subMesh.UpdateBoundingInfo(localWorld, this._scaleFactor)
		}
	}

	return localWorld
}
func (this *Mesh) CreateGlobalSubMesh() *SubMesh {

	if this._totalVertices == 0 || this._indices == nil {
		return nil
	}

	this.SubMeshes = make([]*SubMesh, 0)
	return NewSubMesh(0, 0, this._totalVertices, 0, len(this._indices), this)
}

func (this *Mesh) Subdivide(count int) {
	if count < 1 {
		return
	}

	subdivisionSize := len(this._indices) / count
	offset := 0

	this.SubMeshes = make([]*SubMesh, 0)
	for index := 0; index < count; index++ {
		num := math32.Min(float32(subdivisionSize), float32(len(this._indices)-offset))
		CreateFromIndices(0, offset, int(num), this)

		offset += subdivisionSize
	}
}

func (this *Mesh) SetVerticesData(data []float32, kind string, updatable bool) {

	if this._vertexBuffers == nil {
		this._vertexBuffers = map[string]*VertexBuffer{}
	}

	if _, ok := this._vertexBuffers[kind]; ok {
		this._vertexBuffers[kind].Dispose()
	}

	this._vertexBuffers[kind] = NewVertexBuffer(this, data, kind, updatable)

	if kind == IMesh_VB_PositionKind {
		stride := this._vertexBuffers[kind].GetStrideSize()
		this._totalVertices = len(data) / stride

		this._boundingInfo = cullings.NewBoundingInfo(data, 0, this._totalVertices)

		this.CreateGlobalSubMesh()
	}
}

func (this *Mesh) UpdateVertices(kind string, data []float32) {
	if _, ok := this._vertexBuffers[kind]; ok {
		this._vertexBuffers[kind].Update(data)
	}
}

func (this *Mesh) SetIndices(indices []uint16) {
	engine := this._scene.GetEngine()

	if this._indexBuffer != nil {
		engine.ReleaseIndexBuffer(this._indexBuffer)
	}

	this._indexBuffer = engine.CreateIndexBuffer(indices, false)
	this._indices = indices

	this.CreateGlobalSubMesh()
}

func (this *Mesh) BindAndDraw(subMesh *SubMesh, effect IEffect, wireframe bool) {
	engine := this._scene.GetEngine()

	// Wireframe
	indexToBind := this._indexBuffer
	var useTriangles bool
	useTriangles = true
	if wireframe {
		indexToBind = subMesh.GetLinesIndexBuffer(this._indices, engine)
		useTriangles = false
	}

	glvertexBuffers := map[string]*gl.GLVertexBuffer{}
	for key, val := range this._vertexBuffers {
		glvertexBuffers[key] = val._buffer
	}

	// VBOs
	engine.BindMultiBuffers(glvertexBuffers, indexToBind, effect)

	// Draw order
	if useTriangles {
		engine.Draw(true, subMesh._indexStart, subMesh._indexCount)
	} else {
		engine.Draw(false, 0, subMesh._linesIndexCount)
	}
}

func (this *Mesh) Render(submesh ISubMesh) {

	subMesh, ok := submesh.(*SubMesh)
	if !ok {
		log.Printf("Mesh Rend Failed : submesh is not SubMesh interface")
		return
	}

	engine := this._scene.GetEngine()

	// World
	world := this.GetWorldMatrix()

	// Material
	effectiveMaterial := subMesh.GetMaterial()

	if effectiveMaterial == nil || !effectiveMaterial.IsReady(this) {
		return
	}

	effectiveMaterial.PreBind()
	effectiveMaterial.Bind(world, this)

	haswireframe := false
	if engine.ForceWireframe || effectiveMaterial.HasWireframe() {
		haswireframe = true
	}
	// Bind and draw
	this.BindAndDraw(subMesh, effectiveMaterial.GetEffect(), haswireframe)

	// UnBind
	effectiveMaterial.UnBind()

}

func (this *Mesh) IsDescendantOf(ancestor *Mesh) bool {
	if this.Parent != nil {
		if reflect.DeepEqual(this.Parent, ancestor) {
			return true
		}

		return this.Parent.IsDescendantOf(ancestor)
	}
	return false
}
func (this *Mesh) GetDescendants() []*Mesh {
	results := make([]*Mesh, 0)

	for _, m := range this._scene.Meshes {
		mesh, ok := m.(*Mesh)
		if !ok {
			continue
		}
		if mesh.IsDescendantOf(this) {
			results = append(results, mesh)
		}
	}
	return results
}

func (this *Mesh) GetEmittedParticleSystems() []IParticleSystem {
	results := make([]IParticleSystem, 0)
	for _, particleSystem := range this._scene.ParticleSystems {
		if reflect.DeepEqual(particleSystem.GetEmitter(), this) {
			results = append(results, particleSystem)
		}
	}

	return results
}

func (this *Mesh) GetHierarchyEmittedParticleSystems() []IParticleSystem {
	results := make([]IParticleSystem, 0)
	descendants := this.GetDescendants()
	descendants = append(descendants, this)

	for _, particleSystem := range this._scene.ParticleSystems {
		if tools.IndexOf(particleSystem.GetEmitter(), descendants) != -1 {
			results = append(results, particleSystem)
		}
	}

	return results

}

func (this *Mesh) GetChildren() []*Mesh {
	results := make([]*Mesh, 0)
	for _, m := range this._scene.Meshes {
		mesh, ok := m.(*Mesh)
		if !ok {
			continue
		}
		if reflect.DeepEqual(mesh.Parent, this) {
			results = append(results, mesh)
		}
	}

	return results
}

func (this *Mesh) SetMaterialByID(id string) {
	materials := this._scene.Materials
	for index := 0; index < len(materials); index++ {
		if materials[index].GetId() == id {
			this.Material = materials[index]
			return
		}
	}

	// Multi
	multiMaterials := this._scene.MultiMaterials
	for index := 0; index < len(multiMaterials); index++ {
		if multiMaterials[index].GetId() == id {
			this.MutilMaterial = multiMaterials[index]
			return
		}
	}
}

// Cache
func (this *Mesh) _resetPointsArrayCache() {
	this._cache_positions = nil
}
func (this *Mesh) _generatePointsArray() {
	if this._cache_positions == nil {
		return
	}

	this._cache_positions = make([]*math32.Vector3, 0)

	data := this._vertexBuffers[IMesh_VB_PositionKind].GetData()
	for index := 0; index < len(data); index += 3 {
		this._cache_positions = append(this._cache_positions, math32.NewVector3Zero().FromArray(data, index))

	}
}

//Collisions
func (this *Mesh) _collideForSubMesh(subMesh *SubMesh, transformMatrix *math32.Matrix4, collider ICollider) {
	this._generatePointsArray()
	// Transformation
	if subMesh._lastColliderWorldVertices == nil || !reflect.DeepEqual(subMesh._lastColliderTransformMatrix, transformMatrix) {
		subMesh._lastColliderTransformMatrix = transformMatrix
		subMesh._lastColliderWorldVertices = make([]*math32.Vector3, 0)

		start := subMesh._verticesStart
		end := (subMesh._verticesStart + subMesh._verticesCount)
		for i := start; i < end; i++ {
			pos := this._cache_positions[i].TransformCoordinates(transformMatrix)
			subMesh._lastColliderWorldVertices = append(subMesh._lastColliderWorldVertices, pos)
		}
	}
	// Collide
	collider.Collide(subMesh, subMesh._lastColliderWorldVertices, this._indices, subMesh._indexStart, subMesh._indexStart+subMesh._indexCount, subMesh._verticesStart)
}

func (this *Mesh) _processCollisionsForSubModels(collider ICollider, transformMatrix *math32.Matrix4) {

	for _, subMesh := range this.SubMeshes {

		// Bounding test
		if len(this.SubMeshes) > 1 && !subMesh._checkCollision(collider) {
			continue
		}

		this._collideForSubMesh(subMesh, transformMatrix, collider)
	}
}

func (this *Mesh) _checkCollision(collider ICollider) {
	// Bounding box test
	if !this._boundingInfo.CheckCollision(collider) {
		return
	}

	// Transformation matrix
	this._collisionsScalingMatrix = math32.NewMatrix4().Scaling(1.0/collider.GetRadius().X, 1.0/collider.GetRadius().Y, 1.0/collider.GetRadius().Z)
	this._collisionsTransformMatrix = this._worldMatrix.Multiply(this._collisionsScalingMatrix)

	this._processCollisionsForSubModels(collider, this._collisionsTransformMatrix)
}

func (this *Mesh) IntersectsMesh(mesh *Mesh, precise bool) bool {
	if this._boundingInfo == nil || mesh._boundingInfo == nil {
		return false
	}

	return this._boundingInfo.Intersects(mesh._boundingInfo, precise)
}
func (this *Mesh) IntersectsPoint(point *math32.Vector3) bool {
	if this._boundingInfo == nil {
		return false
	}

	return this._boundingInfo.IntersectsPoint(point)
}

//IAnimationTarget
func (this *Mesh) GetAnimations() []IAnimation {
	return nil
}
func (this *Mesh) GetAnimatables() []IAnimatable {
	return nil
}

/***  IMesh interface start ***/

func (this *Mesh) GetId() string {
	return this.Id
}
func (this *Mesh) GetName() string {
	return this.Name
}
func (this *Mesh) GetPosition() *math32.Vector3 {
	return this.Position
}
func (this *Mesh) GetTotalVertices() int {
	return this._totalVertices
}
func (this *Mesh) GetWorldMatrix() *math32.Matrix4 {
	return this._worldMatrix
}

func (this *Mesh) IsReady() bool {
	return this._isReady
}

func (this *Mesh) IsReceiveShadows() bool {
	return this.ReceiveShadows
}
func (this *Mesh) IsVerticesDataPresent(kind string) bool {
	_, ok := this._vertexBuffers[kind]
	return ok
}

func (this *Mesh) ComputeWorldMatrix() {
	this._computeWorldMatrix()
}

func (this *Mesh) IsEnabled() bool {
	if !this.IsReady() || !this._isEnabled {
		return false
	}

	if this.Parent != nil {
		return this.Parent.IsEnabled()
	}

	return true
}

func (this *Mesh) IsVisible() bool {
	return this.Isvisible
}
func (this *Mesh) IsPickable() bool {
	return this.Ispickable
}

func (this *Mesh) Intersects(ray *math32.Ray) *math32.RayIntersectsResult {
	if this._boundingInfo == nil || !ray.IntersectsSphere(this._boundingInfo.Sphere.GetSphere()) {
		return &math32.RayIntersectsResult{Hit: false, Distance: 0}
	}

	this._generatePointsArray()

	var distance float32
	distance = math.MaxFloat32

	for _, subMesh := range this.SubMeshes {
		// Bounding test
		if len(this.SubMeshes) > 1 && !subMesh.CanIntersects(ray) {
			continue
		}

		result := subMesh.Intersects(ray, this._cache_positions, this._indices)

		if result.Hit {
			if result.Distance < distance && result.Distance >= 0 {
				distance = result.Distance
			}
		}
	}

	if distance >= 0 {
		return &math32.RayIntersectsResult{Hit: true, Distance: distance}
	}

	return &math32.RayIntersectsResult{Hit: false, Distance: 0}
}

func (this *Mesh) CheckCollision(collider ICollider) {
	this._checkCollision(collider)
}

func (this *Mesh) GetVisibility() float32 {
	return this.Visibility
}
func (this *Mesh) IsInFrustrum(frustumPlanes []*math32.Plane) bool {
	return this._boundingInfo.IsInFrustrum(frustumPlanes)
}

func (this *Mesh) GetSubMeshes() []ISubMesh {
	submeshs := make([]ISubMesh, 0)
	for _, sub := range this.SubMeshes {
		submeshs = append(submeshs, sub)
	}
	return submeshs
}

func (this *Mesh) Dispose() {
	if this._vertexBuffers != nil {
		for _, vb := range this._vertexBuffers {
			this._scene.GetEngine().ReleaseVertexBuffer(vb._buffer)
		}
		this._vertexBuffers = nil
	}

	if this._indexBuffer != nil {
		this._scene.GetEngine().ReleaseIndexBuffer(this._indexBuffer)
		this._indexBuffer = nil
	}

	// Remove from scene

	index := tools.IndexOf(this, this._scene.Meshes)
	if index > -1 {
		this._scene.Meshes = append(this._scene.Meshes[:index], this._scene.Meshes[index+1:]...)
	}

	this._isDisposed = true

	// Callback
	if this.OnDispose != nil {
		this.OnDispose()
	}
}

/***  IMesh interface end ***/
