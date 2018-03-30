package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/assimp"
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

func (this *LevelObject) SetTransform(mat assimp.Matrix4x4) {
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			this.Transform.TransformMatrix[this.Transform.TransformMatrix.Index(r, c)] = mat.Values()[c][r]
		}
	}
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
				this.LevelObjects[i].Entity3D.InitModel(m, nil)
			}
			RenderMgr.AddObject(&this.LevelObjects[i].Entity3D, &this.LevelObjects[i].Transform)
		}
	}
}

func (this *Level) RemoveFromScene() {
	for i := 0; i < len(this.LevelObjects); i++ {
		m := ResourceMgr.GetModel(this.LevelObjects[i].Name)
		if m != nil {
			RenderMgr.RemoveObject(&this.LevelObjects[i].Entity3D, &this.LevelObjects[i].Transform)
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
