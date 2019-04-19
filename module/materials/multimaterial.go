package materials

import (
	"github.com/suiqirui1987/fly3d/engines"
	. "github.com/suiqirui1987/fly3d/interfaces"
)

type MultiMaterial struct {
	Name         string
	Id           string
	_scene       *engines.Scene
	SubMaterials []IMaterial
}

func NewMultiMaterial(name string, scene *engines.Scene) *MultiMaterial {
	this := &MultiMaterial{}
	this.Name = name
	this.Id = name
	this._scene = scene
	this._scene.MultiMaterials = append(this._scene.MultiMaterials, this)

	return this
}

/**
interface IMultiMaterial start
*/
func (this *MultiMaterial) GetId() string {
	return this.Id
}

func (this *MultiMaterial) GetSubMaterial(index int) IMaterial {
	if index < 0 || index >= len(this.SubMaterials) {
		return NewStandardMaterial("default material", this._scene)
	}

	return this.SubMaterials[index]
}

/**
interface IMultiMaterial end
*/
