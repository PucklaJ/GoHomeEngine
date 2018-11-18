package gohome

type RenderType uint8

const (
	TYPE_2D           RenderType = iota
	TYPE_3D           RenderType = iota
	TYPE_3D_NORMAL    RenderType = iota
	TYPE_2D_NORMAL    RenderType = iota
	TYPE_3D_INSTANCED RenderType = iota
	TYPE_2D_INSTANCED RenderType = iota
	TYPE_EVERYTHING   RenderType = iota
)

func (this RenderType) Compatible(rtype RenderType) bool {
	if this == TYPE_EVERYTHING || rtype == TYPE_EVERYTHING {
		return true
	}
	switch this {
	case TYPE_2D:
		switch rtype {
		case TYPE_2D:
			return true
		case TYPE_2D_NORMAL:
			return true
		case TYPE_2D_INSTANCED:
			return true
		}
		break
	case TYPE_3D:
		switch rtype {
		case TYPE_3D:
			return true
		case TYPE_3D_NORMAL:
			return true
		case TYPE_3D_INSTANCED:
			return true
		}
		break
	case TYPE_3D_NORMAL:
		switch rtype {
		case TYPE_3D:
			return true
		case TYPE_3D_NORMAL:
			return true
		}
		break
	case TYPE_2D_NORMAL:
		switch rtype {
		case TYPE_2D:
			return true
		case TYPE_2D_NORMAL:
			return true
		}
		break
	case TYPE_3D_INSTANCED:
		switch rtype {
		case TYPE_3D:
			return true
		case TYPE_3D_INSTANCED:
			return true
		}
		break
	case TYPE_2D_INSTANCED:
		switch rtype {
		case TYPE_2D:
			return true
		case TYPE_2D_INSTANCED:
			return true
		}
		break
	}

	return false
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
