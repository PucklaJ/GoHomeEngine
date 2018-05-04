package renderer

import (
	// "fmt"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"sync"
)

type valueTypeIndexOffset struct {
	valueType uint32
	index     uint32
	offset    uint32
}

type OpenGLInstancedMesh3D struct {
	vertices     []gohome.Mesh3DVertex
	indices      []uint32
	numVertices  uint32
	numIndices   uint32
	numInstances uint32

	vao                        uint32
	vertexIndexInstancedBuffer uint32

	Name     string
	Material *gohome.Material

	tangentsCalculated    bool
	customValues          []uint32
	valueTypeIndexOffsets []valueTypeIndexOffset
	instancedSize         uint32
}

func CreateOpenGLInstancedMesh3D(name string) *OpenGLInstancedMesh3D {
	return &OpenGLInstancedMesh3D{
		Name: name,
	}
}

func (this *OpenGLInstancedMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	this.vertices = append(this.vertices, vertices...)
	this.indices = append(this.indices, indices...)
}

func (this *OpenGLInstancedMesh3D) CalculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	vertices := &this.vertices
	indices := this.indices

	var p0, p1, p2 mgl32.Vec3
	var t0, t1, t2 mgl32.Vec2
	var r float32
	var deltaPos1, deltaPos2 mgl32.Vec3
	var deltaUv1, deltaUv2 mgl32.Vec2
	var tangent mgl32.Vec3
	var normal mgl32.Vec3
	var bitangent mgl32.Vec3
	var i uint32
	for i = startIndex; i < maxIndex && i < uint32(len(indices)); i += 3 {
		if i > uint32(len(indices)-3) {
			break
		}

		p0 = mgl32.Vec3{(*vertices)[indices[i]][0], (*vertices)[indices[i]][1], (*vertices)[indices[i]][2]}
		p1 = mgl32.Vec3{(*vertices)[indices[i+1]][0], (*vertices)[indices[i+1]][1], (*vertices)[indices[i+1]][2]}
		p2 = mgl32.Vec3{(*vertices)[indices[i+2]][0], (*vertices)[indices[i+2]][1], (*vertices)[indices[i+2]][2]}

		t0 = mgl32.Vec2{(*vertices)[indices[i]][6], (*vertices)[indices[i]][7]}
		t1 = mgl32.Vec2{(*vertices)[indices[i+1]][6], (*vertices)[indices[i+1]][7]}
		t2 = mgl32.Vec2{(*vertices)[indices[i+2]][6], (*vertices)[indices[i+2]][7]}

		normal = mgl32.Vec3{(*vertices)[indices[i]][3], (*vertices)[indices[i]][4], (*vertices)[indices[i]][5]}

		deltaPos1 = p1.Sub(p0)
		deltaPos2 = p2.Sub(p0)

		deltaUv1 = t1.Sub(t0)
		deltaUv2 = t2.Sub(t0)

		r = 1.0 / (deltaUv1[0]*deltaUv2[1] - deltaUv1[1]*deltaUv2[0])

		tangent = (deltaPos1.Mul(deltaUv2[1]).Sub(deltaPos2.Mul(deltaUv1[1]))).Mul(r).Normalize()
		tangent = tangent.Sub(normal.Mul(normal.Dot(tangent))).Normalize()
		bitangent = (deltaPos2.Mul(deltaUv1[0]).Sub(deltaPos1.Mul(deltaUv2[0]))).Mul(r).Normalize()
		if normal.Cross(tangent).Dot(bitangent) < 0.0 {
			tangent = tangent.Mul(-1.0)
		}
		var j uint32
		for j = 0; j < 3; j++ {
			(*vertices)[indices[i+j]][8] = tangent[0]
			(*vertices)[indices[i+j]][9] = tangent[1]
			(*vertices)[indices[i+j]][10] = tangent[2]
		}
	}
}

func (this *OpenGLInstancedMesh3D) CalculateTangents() {
	if this.tangentsCalculated {
		return
	}
	var wg sync.WaitGroup

	deltaIndex := uint32(len(this.indices)) / NUM_GO_ROUTINES_TANGENTS_CALCULATING
	if deltaIndex == 0 {
		deltaIndex = uint32(len(this.indices)) / 3
	}
	if deltaIndex > 3 {
		deltaIndex -= deltaIndex % 3
	} else {
		deltaIndex = 3
	}

	var i uint32
	for i = 0; i < NUM_GO_ROUTINES_TANGENTS_CALCULATING*2; i++ {
		wg.Add(1)
		go this.CalculateTangentsRoutine(i*deltaIndex, i*deltaIndex+deltaIndex, &wg)
		if i*deltaIndex+deltaIndex >= uint32(len(this.indices)) {
			break
		}
	}

	wg.Wait()

	this.tangentsCalculated = true
}

func getSize(valueType uint32) uint32 {
	switch valueType {
	case gohome.VALUE_FLOAT:
		return 4
	case gohome.VALUE_VEC2:
		return 4 * 2
	case gohome.VALUE_VEC3:
		return 4 * 3
	case gohome.VALUE_VEC4:
		return 4 * 4
	case gohome.VALUE_MAT2:
		return 4 * 2 * 2
	case gohome.VALUE_MAT3:
		return 4 * 3 * 3
	case gohome.VALUE_MAT4:
		return 4 * 4 * 4
	}

	return 0
}

func (this *OpenGLInstancedMesh3D) getInstancedSize() uint32 {
	var sumSize uint32 = 0
	for i := 0; i < len(this.customValues); i++ {
		sumSize += getSize(this.customValues[i])
	}

	return sumSize
}

func vertexAttribPointerForValueType(valueType uint32, offset *uint32, index *uint32, sizeOfOneInstance uint32) {
	switch valueType {
	case gohome.VALUE_FLOAT:
		gl.VertexAttribPointer(*index, 1, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4
		break
	case gohome.VALUE_VEC2:
		gl.VertexAttribPointer(*index, 2, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2
		break
	case gohome.VALUE_VEC3:
		gl.VertexAttribPointer(*index, 3, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 3
		break
	case gohome.VALUE_VEC4:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 4
		break
	case gohome.VALUE_MAT2:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2 * 2
		break
	case gohome.VALUE_MAT3:
		for i := 0; i < 3; i++ {
			gl.VertexAttribPointer(*index, 3, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 3
		}
		break
	case gohome.VALUE_MAT4:
		for i := 0; i < 4; i++ {
			gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 4
		}
		break
	}
}

func (this *OpenGLInstancedMesh3D) instancedVertexAttribPointer(verticesSize uint32, indicesSize uint32, sizeOfOneInstance uint32) {
	offset := verticesSize + indicesSize
	var index uint32 = 4

	for i := 0; i < len(this.customValues); i++ {
		vertexAttribPointerForValueType(this.customValues[i], &offset, &index, sizeOfOneInstance)
	}
}

func enableValueType(valueType uint32, index *uint32) {
	switch valueType {
	case gohome.VALUE_FLOAT:
		gl.EnableVertexAttribArray(*index)
		(*index)++
		break
	case gohome.VALUE_VEC2:
		gl.EnableVertexAttribArray(*index)
		(*index)++
		break
	case gohome.VALUE_VEC3:
		gl.EnableVertexAttribArray(*index)
		(*index)++
		break
	case gohome.VALUE_VEC4:
		gl.EnableVertexAttribArray(*index)
		(*index)++
		break
	case gohome.VALUE_MAT2:
		gl.EnableVertexAttribArray(*index)
		(*index)++
		break
	case gohome.VALUE_MAT3:
		for i := 0; i < 3; i++ {
			gl.EnableVertexAttribArray(*index)
			(*index)++
		}
		break
	case gohome.VALUE_MAT4:
		for i := 0; i < 4; i++ {
			gl.EnableVertexAttribArray(*index)
			(*index)++
		}
		break
	}
}

func (this *OpenGLInstancedMesh3D) instancedEnableVertexAttribArray() {
	var i uint32
	var index uint32 = 4
	for i = 0; i < uint32(len(this.customValues)); i++ {
		enableValueType(this.customValues[i], &index)
	}
}

func (this *OpenGLInstancedMesh3D) deleteElements() {
	this.vertices = append(this.vertices[:0], this.vertices[len(this.vertices):]...)
	this.indices = append(this.indices[:0], this.indices[len(this.indices):]...)
}

func (this *OpenGLInstancedMesh3D) calculateOffsets() {
	var i uint32
	var offset uint32 = this.numVertices*MESH3DVERTEX_SIZE + this.numIndices*gohome.INDEX_SIZE
	for i = 0; i < uint32(len(this.valueTypeIndexOffsets)); i++ {
		this.valueTypeIndexOffsets[i].offset = offset
		offset += getSize(this.valueTypeIndexOffsets[i].valueType)
	}
}

func (this *OpenGLInstancedMesh3D) Load() {
	this.numVertices = uint32(len(this.vertices))
	this.numIndices = uint32(len(this.indices))

	if this.numVertices == 0 || this.numIndices == 0 {
		log.Println("No vertices or indices have been added for instanced mesh", this.Name, "!")
		return
	}
	if this.numInstances == 0 {
		log.Println("Num Instances hasn't been set for instanced mesh", this.Name, "! Will be set to 1")
		this.numInstances = 1
	}

	var verticesSize uint32 = this.numVertices * MESH3DVERTEX_SIZE
	var indicesSize uint32 = this.numIndices * gohome.INDEX_SIZE
	this.instancedSize = this.getInstancedSize() * this.numInstances

	this.CalculateTangents()

	gl.GenVertexArrays(1, &this.vao)
	gl.GenBuffers(1, &this.vertexIndexInstancedBuffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, int(verticesSize+indicesSize+this.instancedSize), nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(verticesSize), gl.Ptr(&this.vertices[0][0]))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.vertexIndexInstancedBuffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), int(indicesSize), gl.Ptr(&this.indices[0]))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(this.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4))
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4+3*4))
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, int32(MESH3DVERTEX_SIZE), gl.PtrOffset(3*4+3*4+2*4))
	this.instancedVertexAttribPointer(verticesSize, indicesSize, this.instancedSize/this.numInstances)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	this.instancedEnableVertexAttribArray()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	gl.BindVertexArray(0)

	this.instancedSize /= this.numInstances

	this.deleteElements()
	this.calculateOffsets()
}
func (this *OpenGLInstancedMesh3D) Render() {
	if this.numVertices == 0 || this.numIndices == 0 {
		return
	}
	if gohome.RenderMgr.CurrentShader != nil {
		if err := gohome.RenderMgr.CurrentShader.SetUniformMaterial(*this.Material); err != nil {
			// fmt.Println("Error instanced:", err)
		}
	}
	gl.BindVertexArray(this.vao)
	gl.DrawElementsInstanced(gl.TRIANGLES, int32(this.numIndices), gl.UNSIGNED_INT, gl.PtrOffset(int(this.numVertices*MESH3DVERTEX_SIZE)), int32(this.numInstances))
	gl.BindVertexArray(0)

}
func (this *OpenGLInstancedMesh3D) Terminate() {
	defer gl.DeleteVertexArrays(1, &this.vao)
	defer gl.DeleteBuffers(1, &this.vertexIndexInstancedBuffer)
}

func (this *OpenGLInstancedMesh3D) SetMaterial(mat *gohome.Material) {
	this.Material = mat
}
func (this *OpenGLInstancedMesh3D) GetMaterial() *gohome.Material {
	if this.Material == nil {
		this.Material = &gohome.Material{}
		this.Material.InitDefault()
	}
	return this.Material
}
func (this *OpenGLInstancedMesh3D) GetName() string {
	return this.Name
}
func (this *OpenGLInstancedMesh3D) GetNumVertices() uint32 {
	return this.numVertices
}
func (this *OpenGLInstancedMesh3D) GetNumIndices() uint32 {
	return this.numIndices
}
func (this *OpenGLInstancedMesh3D) SetNumInstances(n uint32) {
	this.numInstances = n
}
func (this *OpenGLInstancedMesh3D) GetNumInstances() uint32 {
	return this.numInstances
}

func (this *OpenGLInstancedMesh3D) addValueTypeIndexOffset(valueType uint32) {
	var maxIndex uint32 = 0
	for i := 0; i < len(this.valueTypeIndexOffsets); i++ {
		if this.valueTypeIndexOffsets[i].valueType == valueType {
			if this.valueTypeIndexOffsets[i].index >= maxIndex {
				maxIndex = this.valueTypeIndexOffsets[i].index + 1
			}
		}
	}
	this.valueTypeIndexOffsets = append(this.valueTypeIndexOffsets, valueTypeIndexOffset{
		valueType: valueType,
		index:     maxIndex,
		offset:    0,
	})
}

func (this *OpenGLInstancedMesh3D) AddValue(valueType uint32) {
	this.customValues = append(this.customValues, valueType)
	this.addValueTypeIndexOffset(valueType)
}

func (this *OpenGLInstancedMesh3D) getOffset(valueType, index uint32) uint32 {
	for i := 0; i < len(this.valueTypeIndexOffsets); i++ {
		if this.valueTypeIndexOffsets[i].valueType == valueType && this.valueTypeIndexOffsets[i].index == index {
			return this.valueTypeIndexOffsets[i].offset
		}
	}

	return 0
}

func (this *OpenGLInstancedMesh3D) SetF(index uint32, value []float32) {
	offset := this.getOffset(gohome.VALUE_FLOAT, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Float value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_FLOAT)), gl.Ptr(&value[i]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) SetV2(index uint32, value []mgl32.Vec2) {
	offset := this.getOffset(gohome.VALUE_VEC2, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec2 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC2)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (this *OpenGLInstancedMesh3D) SetV3(index uint32, value []mgl32.Vec3) {
	offset := this.getOffset(gohome.VALUE_VEC3, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec3 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC3)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) SetV4(index uint32, value []mgl32.Vec4) {
	offset := this.getOffset(gohome.VALUE_VEC4, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec4 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC4)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) SetM2(index uint32, value []mgl32.Mat2) {
	offset := this.getOffset(gohome.VALUE_MAT2, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat2 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT2)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) SetM3(index uint32, value []mgl32.Mat3) {
	offset := this.getOffset(gohome.VALUE_MAT3, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat3 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT3)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) SetM4(index uint32, value []mgl32.Mat4) {
	offset := this.getOffset(gohome.VALUE_MAT4, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat4 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.vertexIndexInstancedBuffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT4)), gl.Ptr(&value[i][0]))
		offset += this.instancedSize
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLInstancedMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return this.vertices
}
func (this *OpenGLInstancedMesh3D) GetIndices() []uint32 {
	return this.indices
}

func (this *OpenGLInstancedMesh3D) SetName(index uint32, value_type uint32, value string) {
	// Nothing
}
