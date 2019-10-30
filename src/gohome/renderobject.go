package gohome

type RenderType uint16

// The different render types.
// Determines which projection, back buffer, camera etc. will be used for the RenderObject
const (
	TYPE_3D_NORMAL     RenderType = (1 << 1)
	TYPE_2D_NORMAL     RenderType = (1 << 2)
	TYPE_3D_INSTANCED  RenderType = (1 << 3)
	TYPE_2D_INSTANCED  RenderType = (1 << 4)
	TYPE_CASTS_SHADOWS RenderType = (1 << 5)
	TYPE_2D            RenderType = TYPE_2D_NORMAL | TYPE_2D_INSTANCED
	TYPE_3D            RenderType = TYPE_3D_NORMAL | TYPE_3D_INSTANCED
	TYPE_EVERYTHING    RenderType = (1 << 16) - 1
)

// Returns if a RenderType can be drawn using a render type
func (this RenderType) Compatible(rtype RenderType) bool {
	if this == TYPE_2D || this == TYPE_3D || this == TYPE_EVERYTHING {
		return (this & rtype) != 0
	} else {
		return (this & rtype) == this
	}
}

// An interface used for handling the transformation matrices
type TransformableObject interface {
	// Calculates the transformation matrix
	CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int)
	// Sets the current transformation matrix
	SetTransformMatrix(rmgr *RenderManager)
}

// An object that can be rendered
type RenderObject interface {
	// Calls the draw method of the data in the GPU
	Render()
	// Sets the shader used for rendering
	SetShader(s Shader)
	// Returns the shader used for rendering
	GetShader() Shader
	// Set the render type of the object
	SetType(rtype RenderType)
	// Returns the render type of the object
	GetType() RenderType
	// Returns wether this object should be rendered
	IsVisible() bool
	// Returns to which camera this object is not relative to
	NotRelativeCamera() int
	// Sets the transformable object of the RenderObject
	SetTransformableObject(tobj TransformableObject)
	// Returns the transformable object of the RenderObject
	GetTransformableObject() TransformableObject
	// Returns wether this object will be rendered after everything else
	RendersLast() bool
	// Returns wether depth testing is enabled for this object
	HasDepthTesting() bool
}

// An implementation of RenderObject that does nothing
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
