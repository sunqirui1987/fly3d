package collisions

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

// Collisions
func GetNewPosition(scene *engines.Scene, position *math32.Vector3, velocity *math32.Vector3, collider *Collider, maximumRetry int) *math32.Vector3 {
	scaledPosition := position.Divide(collider.Radius)
	scaledVelocity := velocity.Divide(collider.Radius)

	collider.Retry = 0
	collider._initialVelocity = scaledVelocity
	collider._initialPosition = scaledPosition
	finalPosition := _collideWithWorld(scene, scaledPosition, scaledVelocity, collider, maximumRetry)

	finalPosition = finalPosition.Multiply(collider.Radius)

	return finalPosition
}

func _collideWithWorld(scene *engines.Scene, position *math32.Vector3, velocity *math32.Vector3, collider *Collider, maximumRetry int) *math32.Vector3 {
	closeDistance := float32(core.CollisionsEpsilon * 10.0)

	if collider.Retry >= maximumRetry {
		return position
	}

	collider.Initialize(position, velocity, closeDistance)

	// Check all meshes
	for index := 0; index < len(scene.Meshes); index++ {
		mesh := scene.Meshes[index]
		if mesh.IsEnabled() {
			mesh.CheckCollision(collider)
		}
	}

	if !collider.CollisionFound {
		return position.Add(velocity)
	}

	if velocity.X != 0 || velocity.Y != 0 || velocity.Z != 0 {
		response := collider.GetResponse(position, velocity)
		position = response.Position
		velocity = response.Velocity
	}

	if velocity.Length() <= closeDistance {
		return position
	}

	collider.Retry = collider.Retry + 1
	return _collideWithWorld(scene, position, velocity, collider, maximumRetry)
}
