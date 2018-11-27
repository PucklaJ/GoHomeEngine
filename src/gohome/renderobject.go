package gohome

type RenderType uint16

const (
	TYPE_3D_NORMAL    RenderType = (1 << 1)
	TYPE_2D_NORMAL    RenderType = (1 << 2)
	TYPE_3D_INSTANCED RenderType = (1 << 3)
	TYPE_2D_INSTANCED RenderType = (1 << 4)
	TYPE_2D           RenderType = TYPE_2D_NORMAL | TYPE_2D_INSTANCED
	TYPE_3D           RenderType = TYPE_3D_NORMAL | TYPE_3D_INSTANCED
	TYPE_EVERYTHING   RenderType = (1 << 16) - 1
)

func (this RenderType) Compatible(rtype RenderType) bool {
	return (this & rtype) != 0
}

type TransformableObject interface {
	CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int)
	SetTransformMatrix(rmgr *RenderManager)
}

type RenderObject interface {
	Render()
	SetShader(s Shader)
	GetShader() Shader
	SetType(rtype RenderType)
	GetType() RenderType
	IsVisible() bool
	NotRelativeCamera() int
	SetTransformableObject(tobj TransformableObject)
	GetTransformableObject() TransformableObject
	RendersLast() bool
	HasDepthTesting() bool
}

type NilRenderObject struct {
}

func (*NilRenderObject) Render() {

}

func (*NilRenderObject) SetShader(s Shader) {

}
func (*NilRenderObject) GetShader() Shader {
	return nil
}
func (*NilRenderObject) SetType(rtype RenderType) {

}
func (*NilRenderObject) GetType() RenderType {
	return TYPE_EVERYTHING
}
func (*NilRenderObject) IsVisible() bool {
	return true
}
func (*NilRenderObject) NotRelativeCamera() int {
	return -1
}
func (*NilRenderObject) SetTransformableObject(tobj TransformableObject) {

}
func (*NilRenderObject) GetTransformableObject() TransformableObject {
	return nil
}
func (*NilRenderObject) RendersLast() bool {
	return false
}
func (*NilRenderObject) HasDepthTesting() bool {
	return true
}
