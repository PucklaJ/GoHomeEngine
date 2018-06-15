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
}

type Level struct {
	Name         string
	LevelObjects []LevelObject
}

func (this *Level) AddToScene() {
	for i := 0; i < len(this.LevelObjects); i++ {
		m := ResourceMgr.GetModel(this.LevelObjects[i].Name)
		if m != nil {
			if this.LevelObjects[i].Entity3D.Model3D == nil {
				this.LevelObjects[i].Entity3D.InitModel(m)
				this.LevelObjects[i].Entity3D.SetTransformableObject(&this.LevelObjects[i].Transform)
			}
			RenderMgr.AddObject(&this.LevelObjects[i].Entity3D)
		}
	}
}

func (this *Level) RemoveFromScene() {
	for i := 0; i < len(this.LevelObjects); i++ {
		m := ResourceMgr.GetModel(this.LevelObjects[i].Name)
		if m != nil {
			RenderMgr.RemoveObject(&this.LevelObjects[i].Entity3D)
		}
	}
}

func (this *Level) GetModel(name string) *Model3D {
	for i := 0; i < len(this.LevelObjects); i++ {
		if this.LevelObjects[i].Name == name {
			return ResourceMgr.GetModel(name)
		}
	}
	return nil
}
