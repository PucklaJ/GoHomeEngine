package gohome

import (
	"github.com/PucklaJ/mathgl/mgl32"
)

// A wrapper for the TransformableObject
type LevelTransformableObject struct {
	TransformMatrix mgl32.Mat4
}

// Does nothing
func (this *LevelTransformableObject) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {

}

// Sets the transform matrix to the current transform matrix in rendering
func (this *LevelTransformableObject) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix3D(this.TransformMatrix)
}

// A struct holding data of an object from a level
type LevelObject struct {
	// The name of the object
	Name string
	// The transform matrix of the object
	Transform LevelTransformableObject
	// The Entity that is the object
	Entity3D
	// The model that is the object
	*Model3D
}

// A level containing level objects
type Level struct {
	// The name of the level
	Name string
	// All objects of the level
	LevelObjects []LevelObject
	// All entities of the level
	entities []*Entity3D
}

// Adds all objects of the level to the scene
func (this *Level) AddToScene() {
	for i := 0; i < len(this.LevelObjects); i++ {
		var entity Entity3D
		m := this.LevelObjects[i].Model3D
		if m != nil {
			entity.InitModel(m)
			entity.SetTransformableObject(&this.LevelObjects[i].Transform)
			RenderMgr.AddObject(&entity)
			this.entities = append(this.entities, &entity)
		}
	}
}

// Removes all objects of the level from the scene
func (this *Level) RemoveFromScene() {
	for i := 0; i < len(this.entities); i++ {
		RenderMgr.RemoveObject(this.entities[i])
	}
}

// Returns the model with name
func (this *Level) GetModel(name string) *Model3D {
	for i := 0; i < len(this.LevelObjects); i++ {
		if this.LevelObjects[i].Name == name {
			return this.LevelObjects[i].Model3D
		}
	}
	return nil
}
