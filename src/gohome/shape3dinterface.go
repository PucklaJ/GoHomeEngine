package gohome

// An interface that handels low level stuff concerning 3D shapes
type Shape3DInterface interface {
	// Initialises everything
	Init()
	// Adds points to the shape
	AddPoints(points []Shape3DVertex)
	// Returns all points
	GetPoints() []Shape3DVertex
	// Set the draw mode (POINTS,LINES,TRIANGLES)
	SetDrawMode(drawMode uint8)
	// Sets the point size
	SetPointSize(size float32)
	// Sets the line width
	SetLineWidth(size float32)
	// Loads all data to the GPU
	Load()
	// Calls the draw method for the data
	Render()
	// Cleans everything up
	Terminate()
}

// An implementation of Shape3DInterface that does nothing
type NilShape3DInterface struct {
}

func (*NilShape3DInterface) Init() {

}
func (*NilShape3DInterface) AddPoints(points []Shape3DVertex) {

}
func (*NilShape3DInterface) GetPoints() []Shape3DVertex {
	var points []Shape3DVertex
	return points
}

func (*NilShape3DInterface) SetDrawMode(drawMode uint8) {

}

func (*NilShape3DInterface) SetPointSize(size float32) {

}

func (*NilShape3DInterface) SetLineWidth(size float32) {

}

func (*NilShape3DInterface) Load() {

}
func (*NilShape3DInterface) Render() {

}
func (*NilShape3DInterface) Terminate() {

}
