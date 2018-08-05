package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
	"math"
	"strconv"
)

const (
	MESH2DVERTEX_SIZE  uint32 = 2 * 2 * 4 // 2*2*sizeof(float32)
	INDEX_SIZE         uint32 = 4         // sizeof(uint32)
	LINE3D_VERTEX_SIZE uint32 = 3*4 + 4*4
)

type Mesh2DVertex [4]float32

func (m *Mesh2DVertex) Vertex(x, y float32) {
	m[0] = x
	m[1] = y
}

func (m *Mesh2DVertex) TexCoord(u, v float32) {
	m[2] = u
	m[3] = v
}

type Mesh3DVertex [3 + 3 + 2 + 3]float32

func VertexPosIndex(which int) int {
	return which
}

func VertexNormalIndex(which int) int {
	return 3 + which
}

func VertexTexCoordIndex(which int) int {
	return 2*3 + which
}

type AxisAlignedBoundingBox struct {
	Min mgl32.Vec3
	Max mgl32.Vec3
}

func (this *AxisAlignedBoundingBox) String() string {
	maxX := strconv.FormatFloat(float64(this.Max.X()), 'f', 3, 32)
	maxY := strconv.FormatFloat(float64(this.Max.Y()), 'f', 3, 32)
	maxZ := strconv.FormatFloat(float64(this.Max.Z()), 'f', 3, 32)

	minX := strconv.FormatFloat(float64(this.Min.X()), 'f', 3, 32)
	minY := strconv.FormatFloat(float64(this.Min.Y()), 'f', 3, 32)
	minZ := strconv.FormatFloat(float64(this.Min.Z()), 'f', 3, 32)

	return "(Max: " + maxX + "; " + maxY + "; " + maxZ + " | Min: " + minX + "; " + minY + "; " + minZ + ")"
}

func (this AxisAlignedBoundingBox) Intersects(thisPos mgl32.Vec3, other AxisAlignedBoundingBox, otherPos mgl32.Vec3) bool {
	newThisMax := this.Max.Add(thisPos)
	newThisMin := this.Min.Add(thisPos)

	newOtherMax := other.Max.Add(otherPos)
	newOtherMin := other.Min.Add(otherPos)

	return newThisMax.X() > newOtherMin.X() && newThisMax.Y() > newOtherMin.Y() && newThisMax.Z() > newOtherMin.Z() &&
		newThisMin.X() < newOtherMax.X() && newThisMin.Y() < newOtherMax.Y() && newThisMin.Z() < newOtherMax.Z()
}

type Shape3DVertex [3 + 4]float32 // Position + Color
type Line3D [2]Shape3DVertex

func (this *Line3D) SetColor(col color.Color) {
	vec4Col := ColorToVec4(col)
	for j := 0; j < 2; j++ {
		for i := 0; i < 4; i++ {
			(*this)[j][i+3] = vec4Col[i]
		}
	}
}

func (this *Line3D) Color() color.Color {
	return Color{
		R: uint8((*this)[0][3] * 255.0),
		G: uint8((*this)[0][4] * 255.0),
		B: uint8((*this)[0][5] * 255.0),
		A: uint8((*this)[0][6] * 255.0),
	}
}

type TextureRegion struct {
	Min [2]float32
	Max [2]float32
}

func (this TextureRegion) Vec4() mgl32.Vec4 {
	return [4]float32{this.Min[0], this.Min[1], this.Max[0], this.Max[1]}
}

func (this *TextureRegion) FromVec4(v mgl32.Vec4) {
	this.Min = [2]float32{v[0], v[1]}
	this.Max = [2]float32{v[2], v[3]}
}

func (this TextureRegion) Normalize(tex Texture) TextureRegion {
	width := float32(tex.GetWidth())
	height := float32(tex.GetHeight())

	this.Min[0] = this.Min[0]/width + 0.5/width
	this.Min[1] = this.Min[1]/height + 0.5/height
	this.Max[0] = this.Max[0]/width - 0.5/width
	this.Max[1] = this.Max[1]/height - 0.5/height

	return this
}

func (this TextureRegion) String() string {
	return "(" +
		strconv.FormatFloat(float64(this.Min[0]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Min[1]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Max[0]), 'f', 2, 32) + ";" +
		strconv.FormatFloat(float64(this.Max[1]), 'f', 2, 32) +

		")"
}

func (this TextureRegion) Width() float32 {
	return this.Max[0] - this.Min[0]
}

func (this TextureRegion) Height() float32 {
	return this.Max[1] - this.Min[1]
}

type Shape2DVertex [2 + 4]float32 // Position + Color

func (this *Shape2DVertex) Vec2() mgl32.Vec2 {
	return [2]float32{this[0], this[1]}
}

func (this *Shape2DVertex) Make(pos mgl32.Vec2, col color.Color) {
	vecCol := ColorToVec4(col)
	this[0] = pos[0]
	this[1] = pos[1]
	this[2] = vecCol[0]
	this[3] = vecCol[1]
	this[4] = vecCol[2]
	this[5] = vecCol[3]
}

type Line2D [2]Shape2DVertex

func (this *Line2D) SetColor(col color.Color) {
	vec4Col := ColorToVec4(col)
	for j := 0; j < 2; j++ {
		for i := 0; i < 4; i++ {
			(*this)[j][i+2] = vec4Col[i]
		}
	}
}

func (this *Line2D) Color() color.Color {
	return Color{
		R: uint8((*this)[0][2] * 255.0),
		G: uint8((*this)[0][3] * 255.0),
		B: uint8((*this)[0][4] * 255.0),
		A: uint8((*this)[0][5] * 255.0),
	}
}

type Triangle2D [3]Shape2DVertex

func (this *Triangle2D) ToLines() (lines []Line2D) {
	var j uint32 = 1
	var i uint32 = 0

	for i = 0; i < 3; i++ {
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

type Rectangle2D [4]Shape2DVertex

func (this *Rectangle2D) ToTriangles() (tris [2]Triangle2D) {
	for i := 0; i < 3; i++ {
		tris[0][i] = this[i]
	}
	tris[1][0] = this[2]
	tris[1][1] = this[3]
	tris[1][2] = this[0]

	return
}

func (this *Rectangle2D) ToLines() (lines []Line2D) {
	var j uint32 = 1
	var i uint32 = 0

	for i = 0; i < 4; i++ {
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

type Circle2D struct {
	Position mgl32.Vec2
	Radius   float32
	Col      color.Color
}

func FromPolar(radius float32, angle float32) mgl32.Vec2 {
	return mgl32.Vec2{
		radius * float32(math.Cos(float64(mgl32.DegToRad(angle)))),
		radius * float32(math.Sin(float64(mgl32.DegToRad(angle)))),
	}
}

func (this *Circle2D) ToTriangles(numTriangles uint32) (tris []Triangle2D) {
	tris = append(tris, make([]Triangle2D, numTriangles)...)

	var pos1, pos2, pos3 mgl32.Vec2
	var vertex1, vertex2, vertex3 Shape2DVertex
	pos3 = this.Position
	vertex3.Make(pos3, this.Col)
	for i := uint32(0); i < numTriangles; i++ {
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

func (this *Circle2D) ToLines(numLines uint32) (lines []Line2D) {
	lines = append(lines, make([]Line2D, numLines)...)

	var pos1, pos2 mgl32.Vec2
	var vertex1, vertex2 Shape2DVertex
	for i := uint32(0); i < numLines; i++ {
		pos1 = FromPolar(this.Radius, -(float32(i) * 360.0 / float32(numLines))).Add(this.Position)
		pos2 = FromPolar(this.Radius, -(float32(i+1) * 360.0 / float32(numLines))).Add(this.Position)

		vertex1.Make(pos1, this.Col)
		vertex2.Make(pos2, this.Col)

		lines[i][0] = vertex1
		lines[i][1] = vertex2
	}

	return
}

type Polygon2D struct {
	Points []Shape2DVertex
}

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

func (this *Polygon2D) ToTriangles() (tris []Triangle2D) {
	var vertices, ears, reflex []uint32
	var mid mgl32.Vec2
	var max mgl32.Vec2
	var min mgl32.Vec2

	vertices = append(vertices, make([]uint32, len(this.Points))...)
	for i := 0; i < len(this.Points); i++ {
		vertices[i] = uint32(i)
		v := this.Points[i].Vec2()
		mgl32.SetMax(&max[0], &v[0])
		mgl32.SetMax(&max[1], &v[1])

		mgl32.SetMin(&min[0], &v[0])
		mgl32.SetMin(&min[1], &v[1])
	}

	mid = max.Sub(min).Mul(0.5).Add(min)

	for i := 0; i < len(vertices); i++ {
		if isReflex(uint32(i), vertices, reflex, this.Points, mid) {
			reflex = append(reflex, uint32(i))
		} else if isEar(uint32(i), vertices, reflex, this.Points, mid) {
			ears = append(ears, uint32(i))
		} else {
		}
	}

	for len(ears) != 0 && len(vertices) > 2 {
		for i := len(ears) - 1; i >= 0 && i < len(ears); i-- {
			tris = append(tris, makeTriangle(ears[i], vertices, this.Points))
			prev, next := getIndices(ears[i], vertices)
			prevReflex := isReflex(prev, vertices, reflex, this.Points, mid)
			prevEar := !prevReflex && isEar(prev, vertices, reflex, this.Points, mid)
			nextReflex := isReflex(next, vertices, reflex, this.Points, mid)
			nextEar := !nextReflex && isEar(next, vertices, reflex, this.Points, mid)
			takingFront := ears[i] == 0
			takingBack := !takingFront && ears[i] == uint32(len(vertices))-1
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
				if isEar(prev, vertices, reflex, this.Points, mid) {
					if prevReflex {
						ears = append(ears, prev)
						remove(prev, &reflex)
					}
				} else if prevEar {
					remove(prev, &ears)
				} else if prevReflex && isConvex(prev, vertices, reflex, this.Points, mid) {
					remove(prev, &reflex)
				}
			}

			if nextReflex || nextEar {
				if isEar(next, vertices, reflex, this.Points, mid) {
					if nextReflex {
						ears = append(ears, next)
						remove(next, &reflex)
					}
				} else if nextEar {
					remove(next, &ears)
				} else if nextReflex && isConvex(next, vertices, reflex, this.Points, mid) {
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

func getIndices(index uint32, vertices []uint32) (prev, next uint32) {
	next = index + 1
	if next == uint32(len(vertices)) {
		next = 0
	}
	if int32(index)-1 < 0 {
		prev = uint32(len(vertices)) - 1
	} else {
		prev = index - 1
	}

	return
}

func getPoints(index uint32, vertices []uint32, points []Shape2DVertex) (vim1, vi, vip1 mgl32.Vec2) {
	vi = points[vertices[index]].Vec2()
	prev, next := getIndices(index, vertices)
	vim1 = points[vertices[prev]].Vec2()
	vip1 = points[vertices[next]].Vec2()
	return
}

func isConvex(index uint32, vertices []uint32, reflex []uint32, points []Shape2DVertex, mid mgl32.Vec2) bool {
	for i := 0; i < len(reflex); i++ {
		if index == reflex[i] {
			return false
		}
	}

	vim1, vi, vip1 := getPoints(index, vertices, points)

	vprev := vim1.Sub(vi).Normalize()
	vnext := vip1.Sub(vi).Normalize()
	toMid := mid.Sub(vi).Normalize()

	angleprev := mgl32.RadToDeg(float32(math.Acos(float64(vprev.Dot(toMid)))))
	anglenext := mgl32.RadToDeg(float32(math.Acos(float64(vnext.Dot(toMid)))))

	return (anglenext + angleprev) < 180.0
}

func isReflex(index uint32, vertices []uint32, reflex []uint32, points []Shape2DVertex, mid mgl32.Vec2) bool {
	return !isConvex(index, vertices, reflex, points, mid)
}

func isEar(index uint32, vertices []uint32, reflex []uint32, points []Shape2DVertex, mid mgl32.Vec2) bool {
	if !isConvex(index, vertices, reflex, points, mid) {
		return false
	}

	prev, next := getIndices(index, vertices)
	vim1, vi, vip1 := getPoints(index, vertices, points)

	for i := 0; i < len(vertices); i++ {
		if uint32(i) == index || uint32(i) == prev || uint32(i) == next || !isReflex(uint32(i), vertices, reflex, points, mid) {
			continue
		}
		v := points[vertices[i]].Vec2()
		if triangleContains(vim1, vi, vip1, v) {
			return false
		}
	}

	return true
}

func makeTriangle(index uint32, vertices []uint32, points []Shape2DVertex) (tri Triangle2D) {
	prev, next := getIndices(index, vertices)

	tri[0] = points[vertices[prev]]
	tri[1] = points[vertices[index]]
	tri[2] = points[vertices[next]]

	return
}

func remove(index uint32, reflex *[]uint32) {
	for i := 0; i < len(*reflex); i++ {
		if index == (*reflex)[i] {
			*reflex = append((*reflex)[:i], (*reflex)[i+1:]...)
		}
	}
}
