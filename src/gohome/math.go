package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
	"math"
	"strconv"
)

const (
	MESH2DVERTEXSIZE  = 2 * 2 * 4 // 2*2*sizeof(float32)
	INDEXSIZE         = 4         // sizeof(unsigned int 32)
	SHAPE3DVERTEXSIZE = 3*4 + 4*4
	SHAPE2DVERTEXSIZE = 2*4 + 4*4
)

// A vertex of a 2D mesh
type Mesh2DVertex [4]float32

// Sets the position of the vertex
func (m *Mesh2DVertex) Vertex(x, y float32) {
	m[0] = x
	m[1] = y
}

// Sets the uv of the mesh
func (m *Mesh2DVertex) TexCoord(u, v float32) {
	m[2] = u
	m[3] = v
}

// A vertex of a 3D mesh
type Mesh3DVertex [3 + 3 + 2 + 3]float32

// Returns wether one vertex is the same as another
func (this *Mesh3DVertex) Equals(other *Mesh3DVertex) bool {
	for i := 0; i < len(*this); i++ {
		if (*this)[i] != (*other)[i] {
			return false
		}
	}

	return true
}

// Returns the index in the float array
func VertexPosIndex(which int) int {
	return which
}

// Returns the index in the float array
func VertexNormalIndex(which int) int {
	return 3 + which
}

// Returns the index in the float array
func VertexTexCoordIndex(which int) int {
	return 2*3 + which
}

// A bounding box stretching from min to max
type AxisAlignedBoundingBox struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

// Return the values of the bounding box as a string
func (this *AxisAlignedBoundingBox) String() string {
	maxX := strconv.FormatFloat(float64(this.Max.X()), 'f', 3, 32)
	maxY := strconv.FormatFloat(float64(this.Max.Y()), 'f', 3, 32)
	maxZ := strconv.FormatFloat(float64(this.Max.Z()), 'f', 3, 32)

	minX := strconv.FormatFloat(float64(this.Min.X()), 'f', 3, 32)
	minY := strconv.FormatFloat(float64(this.Min.Y()), 'f', 3, 32)
	minZ := strconv.FormatFloat(float64(this.Min.Z()), 'f', 3, 32)

	return "(Max: " + maxX + "; " + maxY + "; " + maxZ + " | Min: " + minX + "; " + minY + "; " + minZ + ")"
}

// Returns wether a bounding box intersects with this bounding box
func (this AxisAlignedBoundingBox) Intersects(thisPos mgl32.Vec3, other AxisAlignedBoundingBox, otherPos mgl32.Vec3) bool {
	newThisMax := this.Max.Add(thisPos)
	newThisMin := this.Min.Add(thisPos)

	newOtherMax := other.Max.Add(otherPos)
	newOtherMin := other.Min.Add(otherPos)

	return newThisMax.X() > newOtherMin.X() && newThisMax.Y() > newOtherMin.Y() && newThisMax.Z() > newOtherMin.Z() &&
		newThisMin.X() < newOtherMax.X() && newThisMin.Y() < newOtherMax.Y() && newThisMin.Z() < newOtherMax.Z()
}

// The vertex of a vertex of a shape 3D
type Shape3DVertex [3 + 4]float32 // Position + Color
// The vertices of a 3D line
type Line3D [2]Shape3DVertex

// Sets the color of the vertex
func (this *Shape3DVertex) SetColor(col color.Color) {
	vec4Col := ColorToVec4(col)
	for i := 0; i < 4; i++ {
		(*this)[3+i] = vec4Col[i]
	}
}

// Sets the color of the line
func (this *Line3D) SetColor(col color.Color) {
	for j := 0; j < 2; j++ {
		(*this)[j].SetColor(col)
	}
}

// Returns the color of the line
func (this *Line3D) Color() color.Color {
	return Color{
		R: uint8((*this)[0][3] * 255.0),
		G: uint8((*this)[0][4] * 255.0),
		B: uint8((*this)[0][5] * 255.0),
		A: uint8((*this)[0][6] * 255.0),
	}
}

// The vertices of a 3D triangle
type Triangle3D [3]Shape3DVertex

// A struct representing a part of a texture
type TextureRegion struct {
	Min [2]float32
	Max [2]float32
}

// Returns the whole struct as a Vec4
func (this TextureRegion) Vec4() mgl32.Vec4 {
	return [4]float32{this.Min[0], this.Min[1], this.Max[0], this.Max[1]}
}

// Gets the struct values from a Vec4
func (this *TextureRegion) FromVec4(v mgl32.Vec4) {
	this.Min = [2]float32{v[0], v[1]}
	this.Max = [2]float32{v[2], v[3]}
}

// Calculates the uv values of the region
func (this TextureRegion) Normalize(tex Texture) TextureRegion {
	width := float32(tex.GetWidth())
	height := float32(tex.GetHeight())

	this.Min[0] = this.Min[0]/width + 0.5/width
	this.Min[1] = this.Min[1]/height + 0.5/height
	this.Max[0] = this.Max[0]/width - 0.5/width
	this.Max[1] = this.Max[1]/height - 0.5/height

	return this
}

// Returns the values of the struct as a string
func (this TextureRegion) String() string {
	return "(" +
		strconv.FormatFloat(float64(this.Min[0]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Min[1]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Max[0]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Max[1]), 'f', 2, 32) +

		")"
}

// Returns the width of the texture region
func (this TextureRegion) Width() float32 {
	return this.Max[0] - this.Min[0]
}

// Returns the height of the texture region
func (this TextureRegion) Height() float32 {
	return this.Max[1] - this.Min[1]
}

// A vertex of a Shape2D
type Shape2DVertex [2 + 4]float32 // Position + Color

// Returns the position of the vertex
func (this *Shape2DVertex) Vec2() mgl32.Vec2 {
	return [2]float32{this[0], this[1]}
}

// Creates a vertex from pos and col
func (this *Shape2DVertex) Make(pos mgl32.Vec2, col color.Color) {
	vecCol := ColorToVec4(col)
	this[0] = pos[0]
	this[1] = pos[1]
	this[2] = vecCol[0]
	this[3] = vecCol[1]
	this[4] = vecCol[2]
	this[5] = vecCol[3]
}

// The vertices of a 2D line
type Line2D [2]Shape2DVertex

// Sets the color of the line
func (this *Line2D) SetColor(col color.Color) {
	vec4Col := ColorToVec4(col)
	for j := 0; j < 2; j++ {
		for i := 0; i < 4; i++ {
			(*this)[j][i+2] = vec4Col[i]
		}
	}
}

// Returns the color of the line
func (this *Line2D) Color() color.Color {
	return Color{
		R: uint8((*this)[0][2] * 255.0),
		G: uint8((*this)[0][3] * 255.0),
		B: uint8((*this)[0][4] * 255.0),
		A: uint8((*this)[0][5] * 255.0),
	}
}

// The vertices of a 3D triangle
type Triangle2D [3]Shape2DVertex

// Converts the triangle to lines going around the triangle
func (this *Triangle2D) ToLines() (lines []Line2D) {
	var j = 1

	for i := 0; i < 3; i++ {
		if j == 3 {
			j = 0
		}

		lines = append(lines, Line2D{
			this[i],
			this[j],
		})

		j++
	}

	return
}

// The vertices of a 2D rectangle
type Rectangle2D [4]Shape2DVertex

// Converts the rectangle into 2 triangles
func (this *Rectangle2D) ToTriangles() (tris [2]Triangle2D) {
	for i := 0; i < 3; i++ {
		tris[0][i] = this[i]
	}
	tris[1][0] = this[2]
	tris[1][1] = this[3]
	tris[1][2] = this[0]

	return
}

// Converts the triangle into lines going around the rectangle
func (this *Rectangle2D) ToLines() (lines []Line2D) {
	var j = 1

	for i := 0; i < 4; i++ {
		if j == 4 {
			j = 0
		}

		lines = append(lines, Line2D{
			this[i],
			this[j],
		})

		j++
	}

	return
}

// A circle in 2D space
type Circle2D struct {
	// The mid point of the circle
	Position mgl32.Vec2
	// The radius of the circle
	Radius   float32
	// The color of the circle
	Col      color.Color
}

// Creates a Vec2 from polar coordinates
func FromPolar(radius float32, angle float32) mgl32.Vec2 {
	return mgl32.Vec2{
		radius * float32(math.Cos(float64(mgl32.DegToRad(angle)))),
		radius * float32(math.Sin(float64(mgl32.DegToRad(angle)))),
	}
}

// Converts the circle into triangles
func (this *Circle2D) ToTriangles(numTriangles int) (tris []Triangle2D) {
	tris = append(tris, make([]Triangle2D, numTriangles)...)

	var pos1, pos2, pos3 mgl32.Vec2
	var vertex1, vertex2, vertex3 Shape2DVertex
	pos3 = this.Position
	vertex3.Make(pos3, this.Col)
	for i := 0; i < numTriangles; i++ {
		pos1 = FromPolar(this.Radius, -(float32(i) * 360.0 / float32(numTriangles))).Add(this.Position)
		pos2 = FromPolar(this.Radius, -(float32(i+1) * 360.0 / float32(numTriangles))).Add(this.Position)

		vertex1.Make(pos1, this.Col)
		vertex2.Make(pos2, this.Col)

		tris[i][0] = vertex1
		tris[i][1] = vertex2
		tris[i][2] = vertex3
	}

	return
}

// Converts the circle into lines going around the circle
func (this *Circle2D) ToLines(numLines int) (lines []Line2D) {
	lines = append(lines, make([]Line2D, numLines)...)

	var pos1, pos2 mgl32.Vec2
	var vertex1, vertex2 Shape2DVertex
	for i := 0; i < numLines; i++ {
		pos1 = FromPolar(this.Radius, -(float32(i) * 360.0 / float32(numLines))).Add(this.Position)
		pos2 = FromPolar(this.Radius, -(float32(i+1) * 360.0 / float32(numLines))).Add(this.Position)

		vertex1.Make(pos1, this.Col)
		vertex2.Make(pos2, this.Col)

		lines[i][0] = vertex1
		lines[i][1] = vertex2
	}

	return
}

// The vertices of a 2D polygon
type Polygon2D struct {
	Points []Shape2DVertex
}

// Converts the polygon into lines going around the polygon
func (this *Polygon2D) ToLines() (lines []Line2D) {
	lines = append(lines, make([]Line2D, len(this.Points))...)

	var j int = 1
	for i := 0; i < len(this.Points); i++ {
		if j == len(this.Points) {
			j = 0
		}

		lines[i] = Line2D{
			this.Points[i],
			this.Points[j],
		}
		j++
	}

	return
}

// Converts the polygon into triangles
func (this *Polygon2D) ToTriangles() (tris []Triangle2D) {
	var vertices, ears, reflex []int

	vertices = append(vertices, make([]int, len(this.Points))...)
	for i := 0; i < len(this.Points); i++ {
		vertices[i] = i
	}

	for i := 0; i < len(vertices); i++ {
		if isReflex(i, vertices, reflex, this.Points) {
			reflex = append(reflex, i)
		} else if isEar(i, vertices, reflex, this.Points) {
			ears = append(ears, i)
		} else {
		}
	}

	for len(ears) != 0 && len(vertices) > 2 {
		for i := len(ears) - 1; i >= 0 && i < len(ears); i-- {
			tris = append(tris, makeTriangle(ears[i], vertices, this.Points))
			prev, next := getIndices(ears[i], vertices)
			prevReflex := isReflex(prev, vertices, reflex, this.Points)
			prevEar := !prevReflex && isEar(prev, vertices, reflex, this.Points)
			nextReflex := isReflex(next, vertices, reflex, this.Points)
			nextEar := !nextReflex && isEar(next, vertices, reflex, this.Points)
			takingFront := ears[i] == 0
			takingBack := !takingFront && ears[i] == len(vertices)-1
			takingMid := !(takingFront || takingBack)
			vertices = append(vertices[:ears[i]], vertices[ears[i]+1:]...)

			if takingFront {
				next = next - 1
				prev = prev - 1
			} else if takingMid {
				next = next - 1
			}

			ears = append(ears[:i], ears[i+1:]...)

			if prevReflex || prevEar {
				if isEar(prev, vertices, reflex, this.Points) {
					if prevReflex {
						ears = append(ears, prev)
						remove(prev, &reflex)
					}
				} else if prevEar {
					remove(prev, &ears)
				} else if prevReflex && isConvex(prev, vertices, reflex, this.Points) {
					remove(prev, &reflex)
				}
			}

			if nextReflex || nextEar {
				if isEar(next, vertices, reflex, this.Points) {
					if nextReflex {
						ears = append(ears, next)
						remove(next, &reflex)
					}
				} else if nextEar {
					remove(next, &ears)
				} else if nextReflex && isConvex(next, vertices, reflex, this.Points) {
					remove(next, &reflex)
				}
			}
		}
	}

	return
}

func triangleContains(pos0, pos1, pos2, pos mgl32.Vec2) bool {
	oneoverdoublearea := 1.0 / (-pos1.Y()*pos2.X() + pos0.Y()*(-pos1.X()+pos2.X()) + pos0.X()*(pos1.Y()-pos2.Y()) + pos1.X()*pos2.Y())
	s := oneoverdoublearea * (pos0.Y()*pos2.X() - pos0.X()*pos2.Y() + (pos2.Y()-pos0.Y())*pos.X() + (pos0.X()-pos2.X())*pos.Y())
	t := oneoverdoublearea * (pos0.X()*pos1.Y() - pos0.Y()*pos1.X() + (pos0.Y()-pos1.Y())*pos.X() + (pos1.X()-pos0.X())*pos.Y())

	return s > 0.0 && t > 0.0 && (1.0-s-t) > 0.0
}

func getIndices(index int, vertices []int) (prev, next int) {
	next = index + 1
	if next == len(vertices) {
		next = 0
	}
	if int32(index)-1 < 0 {
		prev = len(vertices) - 1
	} else {
		prev = index - 1
	}

	return
}

func getPoints(index int, vertices []int, points []Shape2DVertex) (vim1, vi, vip1 mgl32.Vec2) {
	vi = points[vertices[index]].Vec2()
	prev, next := getIndices(index, vertices)
	vim1 = points[vertices[prev]].Vec2()
	vip1 = points[vertices[next]].Vec2()
	return
}

func isConvex(index int, vertices []int, reflex []int, points []Shape2DVertex) bool {
	for i := 0; i < len(reflex); i++ {
		if index == reflex[i] {
			return false
		}
	}

	vim1, vi, vip1 := getPoints(index, vertices, points)

	toPrev := vim1.Sub(vi)
	toNext := vip1.Sub(vi)

	theta1 := mgl32.RadToDeg(toPrev.Angle())
	theta2 := mgl32.RadToDeg(toNext.Angle())
	angle := 180.0 + theta1 - theta2 + 360.0
	for angle > 360.0 {
		angle -= 360.0
	}
	return angle < 180.0
}

func isReflex(index int, vertices []int, reflex []int, points []Shape2DVertex) bool {
	return !isConvex(index, vertices, reflex, points)
}

func isEar(index int, vertices []int, reflex []int, points []Shape2DVertex) bool {
	if !isConvex(index, vertices, reflex, points) {
		return false
	}

	prev, next := getIndices(index, vertices)
	vim1, vi, vip1 := getPoints(index, vertices, points)

	for i := 0; i < len(vertices); i++ {
		if i == index || i == prev || i == next || !isReflex(i, vertices, reflex, points) {
			continue
		}
		v := points[vertices[i]].Vec2()
		if triangleContains(vim1, vi, vip1, v) {
			return false
		}
	}

	return true
}

func makeTriangle(index int, vertices []int, points []Shape2DVertex) (tri Triangle2D) {
	prev, next := getIndices(index, vertices)

	tri[0] = points[vertices[prev]]
	tri[1] = points[vertices[index]]
	tri[2] = points[vertices[next]]

	return
}

func remove(index int, reflex *[]int) {
	for i := 0; i < len(*reflex); i++ {
		if index == (*reflex)[i] {
			*reflex = append((*reflex)[:i], (*reflex)[i+1:]...)
		}
	}
}

// Polygon points used for calculations
type PolygonMath2D []mgl32.Vec2

// Wether one polygon intersects the other
func (this *PolygonMath2D) Intersects(other PolygonMath2D) bool {
	// Seperating axis theorem
	return AreIntersecting(this, &other)
}

// Wether the point is inside the polygon
func (this *PolygonMath2D) IntersectsPoint(point mgl32.Vec2) bool {
	return AreIntersectingPoint(this, point)
}

// A quad used for calculations
type QuadMath2D [4]mgl32.Vec2

// Wether this quad intesects with another
func (this *QuadMath2D) Intersects(other QuadMath2D) bool {
	pm2d := this.ToPolygon()
	return pm2d.Intersects(other.ToPolygon())
}

// Wether a point intersects with this quad
func (this *QuadMath2D) IntersectsPoint(point mgl32.Vec2) bool {
	pm2d := this.ToPolygon()
	return pm2d.IntersectsPoint(point)
}

// Converts this quad into a polygon
func (this *QuadMath2D) ToPolygon() PolygonMath2D {
	return PolygonMath2D((*this)[:])
}

// Converts a screen position to a ray pointing from the camera
func ScreenPositionToRay(point mgl32.Vec2) mgl32.Vec3 {
	return ScreenPositionToRayAdv(point, 0, 0)
}

// Same as ScreenPositionToRay with additional viewport and camera arguments
func ScreenPositionToRayAdv(point mgl32.Vec2, viewportIndex, cameraIndex int32) mgl32.Vec3 {
	// Screen position
	var vppos mgl32.Vec2
	if viewportIndex >= int32(len(RenderMgr.viewport3Ds)) {
		vppos = [2]float32{0.0, 0.0}
	} else {
		viewport := RenderMgr.viewport3Ds[viewportIndex]
		if viewport == nil {
			vppos = [2]float32{0.0, 0.0}
		} else {
			vppos[0], vppos[1] = float32(viewport.X), float32(viewport.Y)
		}
	}
	point = point.Sub(vppos)
	var nres mgl32.Vec2
	if RenderMgr.EnableBackBuffer {
		nres = Render.GetNativeResolution()
	} else {
		nres = Framew.WindowGetSize()
	}
	// Normalized device coordinates
	point[0] = (2.0*point[0])/nres[0] - 1.0
	point[1] = ((2.0*point[1])/nres[1] - 1.0) * -1.0
	// Clipspace
	clippos := mgl32.Vec4{point[0], point[1], -1.0, 1.0}
	// Viewspace
	var invprojmat mgl32.Mat4
	if RenderMgr.Projection3D == nil {
		invprojmat = mgl32.Ident4()
	} else {
		RenderMgr.Projection3D.CalculateProjectionMatrix()
		invprojmat = RenderMgr.Projection3D.GetProjectionMatrix().Inv()
	}
	viewpos := invprojmat.Mul4x1(clippos)
	viewpos[2], viewpos[3] = -1.0, 0.0
	// Worldspace
	var invviewmat mgl32.Mat4
	if cameraIndex == -1 {
		invviewmat = mgl32.Ident4()
	} else {
		if cameraIndex >= int32(len(RenderMgr.camera3Ds)) {
			invviewmat = mgl32.Ident4()
		} else {
			cam := RenderMgr.camera3Ds[cameraIndex]
			if cam == nil {
				invviewmat = mgl32.Ident4()
			} else {
				cam.CalculateViewMatrix()
				invviewmat = cam.GetInverseViewMatrix()
			}
		}
	}
	worldpos := invviewmat.Mul4x1(viewpos).Vec3()
	return worldpos.Normalize()
}

// A Plane used for calculations
type PlaneMath3D struct {
	// The normal pointing from the plane
	Normal mgl32.Vec3
	// A random point on the plane
	Point  mgl32.Vec3
}
