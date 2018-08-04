package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
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
