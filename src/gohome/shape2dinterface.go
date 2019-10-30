package gohome

// The draw modes for the shapes
const (
	DRAW_MODE_POINTS    uint8 = iota
	DRAW_MODE_LINES     uint8 = iota
	DRAW_MODE_TRIANGLES uint8 = iota
)

// An interface that handles all low level stuff concerning shape2d
type Shape2DInterface interface {
	// Initialises the values
	Init()
	// Add points to the shape
	AddPoints(points []Shape2DVertex)
	// Add lines to the shape
	AddLines(lines []Line2D)
	// Add triangles to the shape
	AddTriangles(tris []Triangle2D)
	// Returns all vertices
	GetPoints() []Shape2DVertex
	// Sets the draw mode
	SetDrawMode(mode uint8)
	// Sets the point size
	SetPointSize(size float32)
	// Sets the line width
	SetLineWidth(width float32)
	// Loads everything to the GPU
	Load()
	// Calls the draw method on the data
	Render()
	// Cleans everything up
	Terminate()
}

// An implementation of Shape2DInterface that does nothing
type NilShape2DInterface struct {
}

func (*NilShape2DInterface) Init() {

}
func (*NilShape2DInterface) AddPoints(points []Shape2DVertex) {

}
func (*NilShape2DInterface) AddLines(lines []Line2D) {

}
func (*NilShape2DInterface) AddTriangles(tris []Triangle2D) {

}
func (*NilShape2DInterface) GetPoints() []Shape2DVertex {
	var points []Shape2DVertex
	return points
}
func (*NilShape2DInterface) SetDrawMode(mode uint8) {

}
func (*NilShape2DInterface) SetPointSize(size float32) {

}
func (*NilShape2DInterface) SetLineWidth(width float32) {

}
func (*NilShape2DInterface) Load() {

}
func (*NilShape2DInterface) Render() {

}
func (*NilShape2DInterface) Terminate() {

}
