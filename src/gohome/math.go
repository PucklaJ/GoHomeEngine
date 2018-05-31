package gohome

import (
	"github.com/go-gl/mathgl/mgl32"
	"strconv"
)

const (
	MESH2DVERTEX_SIZE uint32 = 2 * 2 * 4 // 2*2*sizeof(float32)
	INDEX_SIZE        uint32 = 4         // sizeof(uint32)
	LINE3D_VERTEX_SIZE uint32 = 3 * 4 + 4 * 4
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
	maxX := strconv.FormatFloat(float64(this.Max.X()),'f',3,32)
	maxY := strconv.FormatFloat(float64(this.Max.Y()),'f',3,32)
	maxZ := strconv.FormatFloat(float64(this.Max.Z()),'f',3,32)

	minX := strconv.FormatFloat(float64(this.Min.X()),'f',3,32)
	minY := strconv.FormatFloat(float64(this.Min.Y()),'f',3,32)
	minZ := strconv.FormatFloat(float64(this.Min.Z()),'f',3,32)

	return "(Max: " + maxX + "; " + maxY + "; " + maxZ + " | Min: " + minX + "; " + minY + "; " + minZ + ")"
}

func (this *AxisAlignedBoundingBox) Intersects(thisPos mgl32.Vec3,other AxisAlignedBoundingBox,otherPos mgl32.Vec3) bool {
	newThisMax := this.Max.Add(thisPos)
	newThisMin := this.Min.Add(thisPos)

	newOtherMax := other.Max.Add(otherPos)
	newOtherMin := other.Min.Add(otherPos)

	return newThisMax.X() > newOtherMin.X() && newThisMax.Y() > newOtherMin.Y() && newThisMax.Z() > newOtherMin.Z() &&
		   newThisMin.X() < newOtherMax.X() && newThisMin.Y() < newOtherMax.Y() && newThisMin.Z() < newOtherMax.Z()
}

type Line3DVertex [3 + 4]float32 // Position + Color