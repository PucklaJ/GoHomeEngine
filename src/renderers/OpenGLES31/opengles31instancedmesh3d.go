package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles31"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"strconv"
	"sync"
	"unsafe"
)

type valueTypeIndexOffset struct {
	valueType uint32
	index     uint32
	offset    uint32
}

type indexValueType struct {
	index     uint32
	valueType uint32
}

type OpenGLES31InstancedMesh3D struct {
	vertices         []gohome.Mesh3DVertex
	indices          []uint32
	numVertices      uint32
	numIndices       uint32
	numInstances     uint32
	numUsedInstances uint32

	vao    uint32
	buffer uint32
	hasUV  bool
	loaded bool
	aabb   gohome.AxisAlignedBoundingBox

	Name     string
	Material *gohome.Material

	tangentsCalculated    bool
	customValues          []uint32
	valueTypeIndexOffsets []valueTypeIndexOffset
	instancedSize         uint32
	sizePerInstance       uint32
}

func CreateOpenGLES31InstancedMesh3D(name string) *OpenGLES31InstancedMesh3D {
	mesh := &OpenGLES31InstancedMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}

	return mesh
}

func (this *OpenGLES31InstancedMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	this.vertices = append(this.vertices, vertices...)
	this.indices = append(this.indices, indices...)
	this.checkAABB()
}

func (this *OpenGLES31InstancedMesh3D) CalculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
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

		if t0.X() == 0.0 && t0.Y() == 0.0 && t1.X() == 0.0 && t1.Y() == 0.0 && t2.X() == 0.0 && t2.Y() == 0.0 {
			this.hasUV = false
			continue
		}
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

func (this *OpenGLES31InstancedMesh3D) CalculateTangents() {
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
		return getSize(gohome.VALUE_FLOAT) * 2
	case gohome.VALUE_VEC3:
		return getSize(gohome.VALUE_FLOAT) * 3
	case gohome.VALUE_VEC4:
		return getSize(gohome.VALUE_FLOAT) * 4
	case gohome.VALUE_MAT2:
		return getSize(gohome.VALUE_VEC2) * 2
	case gohome.VALUE_MAT3:
		return getSize(gohome.VALUE_VEC3) * 3
	case gohome.VALUE_MAT4:
		return getSize(gohome.VALUE_VEC4) * 4
	}

	return 0
}

func (this *OpenGLES31InstancedMesh3D) getInstancedSize() uint32 {
	var sumSize uint32 = 0
	for i := 0; i < len(this.customValues); i++ {
		sumSize += getSize(this.customValues[i])
	}

	return sumSize
}

func vertexAttribPointerForValueType(valueType uint32, offset *uint32, index *uint32, sizeOfOneInstance uint32) {
	switch valueType {
	case gohome.VALUE_FLOAT:
		gl.VertexAttribPointer(*index, 1, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4
		break
	case gohome.VALUE_VEC2:
		gl.VertexAttribPointer(*index, 2, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2
		break
	case gohome.VALUE_VEC3:
		gl.VertexAttribPointer(*index, 3, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 3
		break
	case gohome.VALUE_VEC4:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 4
		break
	case gohome.VALUE_MAT2:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2 * 2
		break
	case gohome.VALUE_MAT3:
		for i := 0; i < 3; i++ {
			gl.VertexAttribPointer(*index, 3, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 3
		}
		break
	case gohome.VALUE_MAT4:
		for i := 0; i < 4; i++ {
			gl.VertexAttribPointer(*index, 4, gl.FLOAT, gl.FALSE, int32(sizeOfOneInstance), gl.PtrOffset(int(*offset)))
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 4
		}
		break
	}
}

func (this *OpenGLES31InstancedMesh3D) instancedVertexAttribPointer(verticesSize uint32, indicesSize uint32, sizeOfOneInstance uint32) {
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

func (this *OpenGLES31InstancedMesh3D) instancedEnableVertexAttribArray() {
	var i uint32
	var index uint32 = 4
	for i = 0; i < uint32(len(this.customValues)); i++ {
		enableValueType(this.customValues[i], &index)
	}
}

func (this *OpenGLES31InstancedMesh3D) deleteElements() {
	this.vertices = append(this.vertices[:0], this.vertices[len(this.vertices):]...)
	this.indices = append(this.indices[:0], this.indices[len(this.indices):]...)
}

func (this *OpenGLES31InstancedMesh3D) calculateOffsets() {
	var i uint32
	var offset uint32 = this.numVertices*gohome.MESH3DVERTEXSIZE + this.numIndices*gohome.INDEXSIZE
	for i = 0; i < uint32(len(this.valueTypeIndexOffsets)); i++ {
		this.valueTypeIndexOffsets[i].offset = offset
		offset += getSize(this.valueTypeIndexOffsets[i].valueType)
	}
}

func (this *OpenGLES31InstancedMesh3D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, int32(gohome.MESH3DVERTEXSIZE), gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, int32(gohome.MESH3DVERTEXSIZE), gl.PtrOffset(3*4))
	gl.VertexAttribPointer(2, 2, gl.FLOAT, gl.FALSE, int32(gohome.MESH3DVERTEXSIZE), gl.PtrOffset(3*4+3*4))
	gl.VertexAttribPointer(3, 3, gl.FLOAT, gl.FALSE, int32(gohome.MESH3DVERTEXSIZE), gl.PtrOffset(3*4+3*4+2*4))
	this.instancedVertexAttribPointer(this.numVertices*gohome.MESH3DVERTEXSIZE, this.numIndices*gohome.INDEXSIZE, this.sizePerInstance)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	this.instancedEnableVertexAttribArray()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
}

func (this *OpenGLES31InstancedMesh3D) Load() {
	if this.loaded {
		return
	}
	defer func() {
		this.loaded = true
	}()
	this.numVertices = uint32(len(this.vertices))
	this.numIndices = uint32(len(this.indices))

	if this.numVertices == 0 || this.numIndices == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "No vertices or indices have been added!")
		return
	}
	if this.numInstances == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_WARNING, "InstancedMesh3D", this.Name, "Num Instances hasn't been set! Will be set to 1")
		this.numInstances = 1
	}

	var verticesSize uint32 = this.numVertices * gohome.MESH3DVERTEXSIZE
	var indicesSize uint32 = this.numIndices * gohome.INDEXSIZE
	this.instancedSize = this.getInstancedSize() * this.numInstances
	var bufferSize int
	bufferSize = int(verticesSize) + int(indicesSize) + int(this.instancedSize)

	var usage uint32
	usage = gl.DYNAMIC_DRAW

	this.hasUV = true
	this.CalculateTangents()

	var buf [1]uint32
	gl.GenVertexArrays(1, buf[:])
	this.vao = buf[0]

	gl.GenBuffers(1, buf[:])
	this.buffer = buf[0]
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "GenBuffer: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, usage)
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "BufferData: ")

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(verticesSize), unsafe.Pointer(&this.vertices[0]))
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "BufferSubData Vertices: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), int(indicesSize), unsafe.Pointer(&this.indices[0]))
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "BufferSubData Indices: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	this.sizePerInstance = this.instancedSize / this.numInstances
	this.numUsedInstances = this.numInstances

	gl.BindVertexArray(this.vao)
	this.attributePointer()
	gl.BindVertexArray(0)

	this.deleteElements()
	this.calculateOffsets()
}
func (this *OpenGLES31InstancedMesh3D) Render() {
	if this.numUsedInstances == 0 {
		gohome.ErrorMgr.Warning("InstancedMesh3D", this.Name, "numUsedInstances == 0")
		return
	}

	if this.numVertices == 0 || this.numIndices == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "No Vertices or Indices have been loaded!")
		return
	}
	if gohome.RenderMgr.CurrentShader != nil {
		if this.Material == nil {
			this.Material = &gohome.Material{}
			this.Material.InitDefault()
			gohome.ErrorMgr.Log("", this.Name, "InitDefault")
		}
		gohome.RenderMgr.CurrentShader.SetUniformMaterial(*this.Material)
	}

	gl.BindVertexArray(this.vao)

	gl.GetError()
	gl.DrawElementsInstanced(gl.TRIANGLES, int32(this.numIndices), gl.UNSIGNED_INT, gl.PtrOffset(int(this.numVertices*gohome.MESH3DVERTEXSIZE)), int32(this.numUsedInstances))
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "RenderError: ")

	gl.BindVertexArray(0)
}

func (this *OpenGLES31InstancedMesh3D) Terminate() {
	var vbuf [1]uint32
	vbuf[0] = this.vao
	defer gl.DeleteVertexArrays(1, vbuf[:])
	var buf [1]uint32
	buf[0] = this.buffer
	defer gl.DeleteBuffers(1, buf[:])
}

func (this *OpenGLES31InstancedMesh3D) SetMaterial(mat *gohome.Material) {
	this.Material = mat
}
func (this *OpenGLES31InstancedMesh3D) GetMaterial() *gohome.Material {
	if this.Material == nil {
		this.Material = &gohome.Material{}
		this.Material.InitDefault()
	}
	return this.Material
}
func (this *OpenGLES31InstancedMesh3D) GetName() string {
	return this.Name
}

func (this *OpenGLES31InstancedMesh3D) GetNumVertices() uint32 {
	return this.numVertices
}

func (this *OpenGLES31InstancedMesh3D) GetNumIndices() uint32 {
	return this.numIndices
}

func (this *OpenGLES31InstancedMesh3D) recreateBuffer(numInstances uint32) {
	verticesSize := this.numVertices * gohome.MESH3DVERTEXSIZE
	indicesSize := this.numIndices * gohome.INDEXSIZE
	this.instancedSize = this.getInstancedSize() * numInstances
	bufferSize := int(verticesSize) + int(indicesSize) + int(this.instancedSize)
	var tempBuffer uint32

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	tempBuffer = buf[0]
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "SetNumInstances GenBuffer: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, tempBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, gl.DYNAMIC_DRAW)
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "SetNumInstances BufferData: ")

	gl.CopyBufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.ARRAY_BUFFER, 0, 0, int(verticesSize+indicesSize))
	handleOpenGLES31Error("InstancedMesh3D", this.Name, "SetNumInstances CopyBufferSubData Vertices: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	buf[0] = this.buffer
	gl.DeleteBuffers(1, buf[:])
	this.buffer = tempBuffer

	this.sizePerInstance = this.instancedSize / numInstances

	gl.BindVertexArray(this.vao)
	this.attributePointer()
	gl.BindVertexArray(0)

	this.numInstances = numInstances
	this.calculateOffsets()
}

func (this *OpenGLES31InstancedMesh3D) SetNumInstances(n uint32) {
	if this.numInstances != n {
		if this.loaded {
			this.recreateBuffer(n)
		}
		this.numInstances = n
		this.numUsedInstances = n
	}
}
func (this *OpenGLES31InstancedMesh3D) GetNumInstances() uint32 {
	return this.numInstances
}

func (this *OpenGLES31InstancedMesh3D) addValueTypeIndexOffset(valueType uint32) {
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

func (this *OpenGLES31InstancedMesh3D) addValueTypeIndexOffsetFront(valueType uint32) {
	var maxIndex uint32 = 0
	for i := 0; i < len(this.valueTypeIndexOffsets); i++ {
		if this.valueTypeIndexOffsets[i].valueType == valueType {
			if this.valueTypeIndexOffsets[i].index >= maxIndex {
				maxIndex = this.valueTypeIndexOffsets[i].index + 1
			}
		}
	}
	this.valueTypeIndexOffsets = append([]valueTypeIndexOffset{
		valueTypeIndexOffset{
			valueType: valueType,
			index:     maxIndex,
			offset:    0,
		},
	}, this.valueTypeIndexOffsets...)
}

func (this *OpenGLES31InstancedMesh3D) AddValueFront(valueType uint32) {
	this.customValues = append(this.customValues, valueType)
	this.addValueTypeIndexOffsetFront(valueType)
}

func (this *OpenGLES31InstancedMesh3D) AddValue(valueType uint32) {
	this.customValues = append(this.customValues, valueType)
	this.addValueTypeIndexOffset(valueType)
}

func (this *OpenGLES31InstancedMesh3D) getOffset(valueType, index uint32) uint32 {
	for i := 0; i < len(this.valueTypeIndexOffsets); i++ {
		if this.valueTypeIndexOffsets[i].valueType == valueType && this.valueTypeIndexOffsets[i].index == index {
			return this.valueTypeIndexOffsets[i].offset
		}
	}

	return 0
}

func (this *OpenGLES31InstancedMesh3D) SetF(index uint32, value []float32) {
	offset := this.getOffset(gohome.VALUE_FLOAT, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Float value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_FLOAT)), unsafe.Pointer(&value[i]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) SetV2(index uint32, value []mgl32.Vec2) {
	offset := this.getOffset(gohome.VALUE_VEC2, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec2 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC2)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (this *OpenGLES31InstancedMesh3D) SetV3(index uint32, value []mgl32.Vec3) {
	offset := this.getOffset(gohome.VALUE_VEC3, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec3 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC3)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) SetV4(index uint32, value []mgl32.Vec4) {
	offset := this.getOffset(gohome.VALUE_VEC4, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec4 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_VEC4)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) SetM2(index uint32, value []mgl32.Mat2) {
	offset := this.getOffset(gohome.VALUE_MAT2, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat2 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT2)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) SetM3(index uint32, value []mgl32.Mat3) {
	offset := this.getOffset(gohome.VALUE_MAT3, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat3 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT3)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) SetM4(index uint32, value []mgl32.Mat4) {
	offset := this.getOffset(gohome.VALUE_MAT4, index)
	if offset == 0 {
		return
	}
	if uint32(len(value)) < this.numInstances {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat4 value "+strconv.Itoa(int(index))+" is too small!")
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

	var i uint32
	for i = 0; i < this.numInstances; i++ {
		gl.BufferSubData(gl.ARRAY_BUFFER, int(offset), int(getSize(gohome.VALUE_MAT4)), unsafe.Pointer(&value[i][0]))
		offset += this.sizePerInstance
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
func (this *OpenGLES31InstancedMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return this.vertices
}
func (this *OpenGLES31InstancedMesh3D) GetIndices() []uint32 {
	return this.indices
}

func (this *OpenGLES31InstancedMesh3D) SetName(index uint32, value_type uint32, value string) {
	return
}

func (this *OpenGLES31InstancedMesh3D) HasUV() bool {
	return this.hasUV
}
func (this *OpenGLES31InstancedMesh3D) AABB() gohome.AxisAlignedBoundingBox {
	return this.aabb
}
func (this *OpenGLES31InstancedMesh3D) Copy() gohome.Mesh3D {
	return nil
}

func (this *OpenGLES31InstancedMesh3D) checkAABB() {
	var max, min mgl32.Vec3 = [3]float32{this.vertices[0][0], this.vertices[0][1], this.vertices[0][2]}, [3]float32{this.vertices[0][0], this.vertices[0][1], this.vertices[0][2]}
	var current gohome.Mesh3DVertex
	for i := 0; i < len(this.vertices); i++ {
		current = this.vertices[i]
		for j := 0; j < 3; j++ {
			if current[j] > max[j] {
				max[j] = current[j]
			} else if current[j] < min[j] {
				min[j] = current[j]
			}
		}
	}

	for i := 0; i < 3; i++ {
		if max[i] > this.aabb.Max[i] {
			this.aabb.Max[i] = max[i]
		}
		if min[i] < this.aabb.Min[i] {
			this.aabb.Min[i] = min[i]
		}
	}
}

func (this *OpenGLES31InstancedMesh3D) SetNumUsedInstances(n uint32) {
	this.numUsedInstances = n
}

func (this *OpenGLES31InstancedMesh3D) GetNumUsedInstances() uint32 {
	return this.numUsedInstances
}

func (this *OpenGLES31InstancedMesh3D) LoadedToGPU() bool {
	return this.loaded
}

func (this *OpenGLES31Renderer) InstancedMesh3DFromLoadedMesh3D(mesh gohome.Mesh3D) gohome.InstancedMesh3D {
	oglmesh := mesh.(*OpenGLES31Mesh3D)
	ioglmesh := CreateOpenGLES31InstancedMesh3D(oglmesh.Name)
	ioglmesh.numVertices = oglmesh.numVertices
	ioglmesh.numIndices = oglmesh.numIndices

	ioglmesh.numInstances = 0

	var verticesSize uint32 = ioglmesh.numVertices * gohome.MESH3DVERTEXSIZE
	var indicesSize uint32 = ioglmesh.numIndices * gohome.INDEXSIZE

	ioglmesh.instancedSize = 0

	bufferSize := int(verticesSize) + int(indicesSize)
	var usage uint32
	usage = gl.DYNAMIC_DRAW

	var buf [1]uint32
	gl.GenVertexArrays(1, buf[:])
	ioglmesh.vao = buf[0]
	gl.GenBuffers(1, buf[:])
	ioglmesh.buffer = buf[0]

	handleOpenGLES31Error("InstancedMesh3D", ioglmesh.Name, "GenBuffer: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, ioglmesh.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, usage)
	handleOpenGLES31Error("InstancedMesh3D", ioglmesh.Name, "BufferData: ")

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglmesh.buffer)
	gl.CopyBufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.ARRAY_BUFFER, 0, 0, bufferSize)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	ioglmesh.sizePerInstance = 0
	ioglmesh.numUsedInstances = 0

	gl.BindVertexArray(ioglmesh.vao)
	ioglmesh.attributePointer()
	gl.BindVertexArray(0)

	ioglmesh.deleteElements()
	ioglmesh.calculateOffsets()
	ioglmesh.loaded = true

	ioglmesh.aabb = oglmesh.aabb
	ioglmesh.hasUV = oglmesh.hasUV
	ioglmesh.tangentsCalculated = true

	return ioglmesh
}
