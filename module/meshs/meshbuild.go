package meshs

import (
	"math"

	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
)

func CreateBox(name string, size float32, scene *engines.Scene, updatable bool) *Mesh {
	box := NewMesh(name, scene)

	normalsSource := []*math32.Vector3{
		math32.NewVector3(0, 0, 1),
		math32.NewVector3(0, 0, -1),
		math32.NewVector3(1, 0, 0),
		math32.NewVector3(-1, 0, 0),
		math32.NewVector3(0, 1, 0),
		math32.NewVector3(0, -1, 0),
	}

	indices := make([]uint16, 0)
	positions := make([]float32, 0)
	normals := make([]float32, 0)
	uvs := make([]float32, 0)

	// Create each face in turn.
	for index := 0; index < len(normalsSource); index++ {
		normal := normalsSource[index].Clone()

		// Get two vectors perpendicular to the face normal and to each other.
		side1 := math32.NewVector3(normal.Y, normal.Z, normal.X)
		side2 := normal.Cross(side1)

		// Six indices (two triangles) per face.
		verticesLength := uint16(len(positions) / 3)
		indices = append(indices, verticesLength)
		indices = append(indices, verticesLength+1)
		indices = append(indices, verticesLength+2)

		indices = append(indices, verticesLength)
		indices = append(indices, verticesLength+2)
		indices = append(indices, verticesLength+3)

		// Four vertices per face.
		vertex := normal.Sub(side1).Sub(side2).Scale(size / 2)
		positions = append(positions, vertex.X, vertex.Y, vertex.Z)
		normals = append(normals, normal.X, normal.Y, normal.Z)
		uvs = append(uvs, 1.0, 1.0)

		vertex = normal.Sub(side1).Add(side2).Scale(size / 2)
		positions = append(positions, vertex.X, vertex.Y, vertex.Z)
		normals = append(normals, normal.X, normal.Y, normal.Z)
		uvs = append(uvs, 0.0, 1.0)

		vertex = normal.Add(side1).Add(side2).Scale(size / 2)
		positions = append(positions, vertex.X, vertex.Y, vertex.Z)
		normals = append(normals, normal.X, normal.Y, normal.Z)
		uvs = append(uvs, 0.0, 0.0)

		vertex = normal.Add(side1).Sub(side2).Scale(size / 2)
		positions = append(positions, vertex.X, vertex.Y, vertex.Z)
		normals = append(normals, normal.X, normal.Y, normal.Z)
		uvs = append(uvs, 1.0, 0.0)
	}

	box.SetVerticesData(positions, IMesh_VB_PositionKind, updatable)
	box.SetVerticesData(normals, IMesh_VB_NormalKind, updatable)
	box.SetVerticesData(uvs, IMesh_VB_UVKind, updatable)
	box.SetIndices(indices)

	return box
}

func CreateSphere(name string, segments int, diameter float32, scene *engines.Scene, updatable bool) *Mesh {
	sphere := NewMesh(name, scene)

	var totalZRotationSteps, totalYRotationSteps, zRotationStep, yRotationStep, verticesCount uint16

	radius := float32(diameter) / 2.0

	totalZRotationSteps = uint16(2 + segments)
	totalYRotationSteps = uint16(2 * totalZRotationSteps)

	indices := make([]uint16, 0)
	positions := make([]float32, 0)
	normals := make([]float32, 0)
	uvs := make([]float32, 0)

	for zRotationStep = 0; zRotationStep <= totalZRotationSteps; zRotationStep++ {
		normalizedZ := float32(zRotationStep) / float32(totalZRotationSteps)
		angleZ := (normalizedZ * math.Pi)

		for yRotationStep = 0; yRotationStep <= totalYRotationSteps; yRotationStep++ {
			normalizedY := float32(yRotationStep) / float32(totalYRotationSteps)

			angleY := normalizedY * math.Pi * 2.0

			rotationZ := math32.NewMatrix4().RotationZ(-angleZ)
			rotationY := math32.NewMatrix4().RotationY(angleY)
			afterRotZ := math32.NewVector3Up().TransformCoordinates(rotationZ)
			complete := afterRotZ.TransformCoordinates(rotationY)

			vertex := complete.Scale(radius)
			normal := vertex.NormalizeTo()

			positions = append(positions, vertex.X, vertex.Y, vertex.Z)
			normals = append(normals, normal.X, normal.Y, normal.Z)
			uvs = append(uvs, normalizedZ, normalizedY)
		}

		if zRotationStep > 0 {

			verticesCount = uint16(len(positions) / 3)
			var firstIndex uint16
			for firstIndex = verticesCount - 2*(totalYRotationSteps+1); (firstIndex + totalYRotationSteps + 2) < verticesCount; firstIndex++ {
				indices = append(indices, (firstIndex))
				indices = append(indices, (firstIndex + 1))
				indices = append(indices, firstIndex+totalYRotationSteps+1)

				indices = append(indices, (firstIndex + totalYRotationSteps + 1))
				indices = append(indices, (firstIndex + 1))
				indices = append(indices, (firstIndex + totalYRotationSteps + 2))
			}
		}
	}

	sphere.SetVerticesData(positions, IMesh_VB_PositionKind, updatable)
	sphere.SetVerticesData(normals, IMesh_VB_NormalKind, updatable)
	sphere.SetVerticesData(uvs, IMesh_VB_UVKind, updatable)
	sphere.SetIndices(indices)

	return sphere
}

func CreatePlane(name string, size int, scene *engines.Scene, updatable bool) *Mesh {
	plane := NewMesh(name, scene)

	indices := make([]uint16, 0)
	positions := make([]float32, 0)
	normals := make([]float32, 0)
	uvs := make([]float32, 0)
	// Vertices
	halfSize := float32(size / 2.0)
	positions = append(positions, -halfSize, -halfSize, 0)
	normals = append(normals, 0, 0, -1.0)
	uvs = append(uvs, 0.0, 0.0)

	positions = append(positions, halfSize, -halfSize, 0)
	normals = append(normals, 0, 0, -1.0)
	uvs = append(uvs, 1.0, 0.0)

	positions = append(positions, halfSize, halfSize, 0)
	normals = append(normals, 0, 0, -1.0)
	uvs = append(uvs, 1.0, 1.0)

	positions = append(positions, -halfSize, halfSize, 0)
	normals = append(normals, 0, 0, -1.0)
	uvs = append(uvs, 0.0, 1.0)

	// Indices
	indices = append(indices, uint16(0))
	indices = append(indices, uint16(1))
	indices = append(indices, uint16(2))

	indices = append(indices, uint16(0))
	indices = append(indices, uint16(2))
	indices = append(indices, uint16(3))

	plane.SetVerticesData(positions, IMesh_VB_PositionKind, updatable)
	plane.SetVerticesData(normals, IMesh_VB_NormalKind, updatable)
	plane.SetVerticesData(uvs, IMesh_VB_UVKind, updatable)
	plane.SetIndices(indices)

	return plane
}

func CreateGround(name string, width int, height int, subdivisions int, scene *engines.Scene, updatable bool) *Mesh {
	ground := NewMesh(name, scene)

	indices := make([]uint16, 0)
	positions := make([]float32, 0)
	normals := make([]float32, 0)
	uvs := make([]float32, 0)
	var row, col int

	var x, y, z float32

	for row = 0; row <= subdivisions; row++ {
		for col = 0; col <= subdivisions; col++ {
			x = float32(col*width)/float32(subdivisions) - (float32(width) / 2.0)
			y = 0.0
			z = float32((subdivisions-row)*height)/float32(subdivisions) - (float32(height) / 2.0)
			position := math32.NewVector3(x, y, z)
			normal := math32.NewVector3(0, 1.0, 0)

			positions = append(positions, position.X, position.Y, position.Z)
			normals = append(normals, normal.X, normal.Y, normal.Z)
			uvs = append(uvs, float32(col)/float32(subdivisions), 1.0-float32(row)/float32(subdivisions))
		}
	}

	for row = 0; row < subdivisions; row++ {
		for col = 0; col < subdivisions; col++ {
			indices = append(indices, uint16(col+1+(row+1)*(subdivisions+1)))
			indices = append(indices, uint16(col+1+row*(subdivisions+1)))
			indices = append(indices, uint16(col+row*(subdivisions+1)))

			indices = append(indices, uint16(col+(row+1)*(subdivisions+1)))
			indices = append(indices, uint16(col+1+(row+1)*(subdivisions+1)))
			indices = append(indices, uint16(col+row*(subdivisions+1)))
		}
	}

	ground.SetVerticesData(positions, IMesh_VB_PositionKind, updatable)
	ground.SetVerticesData(normals, IMesh_VB_NormalKind, updatable)
	ground.SetVerticesData(uvs, IMesh_VB_UVKind, updatable)
	ground.SetIndices(indices)

	return ground
}

//"torus", 8, 2, 32, scene, false
func CreateTorus(name string, diameter float32, thickness float32, tessellation int, scene *engines.Scene, updatable bool) *Mesh {

	torus := NewMesh(name, scene)

	indices := make([]uint16, 0)
	positions := make([]float32, 0)
	normals := make([]float32, 0)
	uvs := make([]float32, 0)

	stride := tessellation + 1

	for i := 0; i <= tessellation; i++ {
		var u, v float32
		var outerAngle, innerAngle float32

		u = float32(i) / float32(tessellation)

		outerAngle = float32(i)*math.Pi*2.0/float32(tessellation) - math.Pi/2.0

		transform := math32.NewMatrix4().Translation(diameter/2.0, 0, 0).Multiply(math32.NewMatrix4().RotationY(outerAngle))

		for j := 0; j <= tessellation; j++ {
			v = float32(1 - float32(j)/float32(tessellation))

			innerAngle = float32(j)*math.Pi*2.0/float32(tessellation) + math.Pi
			dx := math32.Cos(innerAngle)
			dy := math32.Sin(innerAngle)

			// Create a vertex.
			normal := math32.NewVector3(dx, dy, 0)
			position := normal.Scale(thickness / 2.0)
			textureCoordinate := math32.NewVector2(u, v)

			position = position.TransformCoordinates(transform)
			normal = normal.TransformNormal(transform)

			positions = append(positions, position.X, position.Y, position.Z)
			normals = append(normals, normal.X, normal.Y, normal.Z)
			uvs = append(uvs, textureCoordinate.X, textureCoordinate.Y)

			// And create indices for two triangles.
			nextI := (i + 1) % stride
			nextJ := (j + 1) % stride

			indices = append(indices, uint16(i*stride+j))
			indices = append(indices, uint16(i*stride+nextJ))
			indices = append(indices, uint16(nextI*stride+j))

			indices = append(indices, uint16(i*stride+nextJ))
			indices = append(indices, uint16(nextI*stride+nextJ))
			indices = append(indices, uint16(nextI*stride+j))
		}
	}

	torus.SetVerticesData(positions, IMesh_VB_PositionKind, updatable)
	torus.SetVerticesData(normals, IMesh_VB_NormalKind, updatable)
	torus.SetVerticesData(uvs, IMesh_VB_UVKind, updatable)
	torus.SetIndices(indices)

	return torus

}
