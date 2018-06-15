package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
)

type LevelTransformableObject struct {
	TransformMatrix mgl32.Mat4
}

func (this *LevelTransformableObject) CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int) {

}

func (this *LevelTransformableObject) SetTransformMatrix(rmgr *RenderManager) {
	rmgr.setTransformMatrix3D(this.TransformMatrix)
}

type LevelObject struct {
	Name      string
	Transform LevelTransformableObject
	Entity3D
	*Model3D
}

type Level struct {
	Name         string
	LevelObjects []LevelObject
	entities 	 []*Entity3D
}

func (this *Level) AddToScene() {
	for i := 0; i < len(this.LevelObjects); i++ {
		var entity Entity3D
		m := this.LevelObjects[i].Model3D
		if m != nil {
			entity.InitModel(m)
			entity.SetTransformableObject(&this.LevelObjects[i].Transform)
			RenderMgr.AddObject(&entity)
			this.entities = append(this.entities,&entity)
		}
	}
}

func (this *Level) RemoveFromScene() {
	for i := 0; i < len(this.entities); i++ {
		RenderMgr.RemoveObject(this.entities[i])
	}
}

func (this *Level) GetModel(name string) *Model3D {
	for i := 0; i < len(this.LevelObjects); i++ {
		if this.LevelObjects[i].Name == name {
			return this.LevelObjects[i].Model3D
		}
	}
	return nil
}
