package renderer

import (
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"strconv"
	"sync"
)

type valueTypeIndexOffset struct {
	valueType uint32
	index     int
	offset    int
}

type indexValueType struct {
	index     int
	valueType uint32
}

type WebGLInstancedMesh3D struct {
	vertices         []gohome.Mesh3DVertex
	indices          []uint16
	numVertices      int
	numIndices       int
	numInstances     int
	numUsedInstances int

	vao             *js.Object
	buffer          *js.Object
	canUseVAOs      bool
	canUseInstanced bool
	hasUV           bool
	loaded          bool
	aabb            gohome.AxisAlignedBoundingBox

	Name     string
	Material *gohome.Material

	tangentsCalculated    bool
	customValues          []uint32
	valueTypeIndexOffsets []valueTypeIndexOffset
	instancedSize         int
	sizePerInstance       int

	floats        [][]float32
	vec2s         [][]mgl32.Vec2
	vec3s         [][]mgl32.Vec3
	vec4s         [][]mgl32.Vec4
	mat2s         [][]mgl32.Mat2
	mat3s         [][]mgl32.Mat3
	mat4s         [][]mgl32.Mat4
	namesForIndex map[indexValueType]string
}

func CreateWebGLInstancedMesh3D(name string) *WebGLInstancedMesh3D {
	mesh := &WebGLInstancedMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}
	render, _ := gohome.Render.(*WebGLRenderer)
	mesh.canUseVAOs = render.HasFunctionAvailable("VERTEX_ARRAY")
	mesh.canUseInstanced = render.HasFunctionAvailable("INSTANCED")
	if !mesh.canUseInstanced {
		mesh.namesForIndex = make(map[indexValueType]string)
	}

	return mesh
}

func (this *WebGLInstancedMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	this.vertices = append(this.vertices, vertices...)
	index := len(this.indices)
	this.indices = append(this.indices, make([]uint16, len(indices))...)
	for id, i := range indices {
		this.indices[index+id] = uint16(i)
	}
	this.checkAABB()
}

func (this *WebGLInstancedMesh3D) CalculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
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

func (this *WebGLInstancedMesh3D) CalculateTangents() {
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

func getSize(valueType uint32) int {
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

func (this *WebGLInstancedMesh3D) getInstancedSize() int {
	var sumSize = 0
	for i := 0; i < len(this.customValues); i++ {
		sumSize += getSize(this.customValues[i])
	}

	return sumSize
}

func vertexAttribPointerForValueType(valueType uint32, offset *int, index *int, sizeOfOneInstance int) {
	switch valueType {
	case gohome.VALUE_FLOAT:
		gl.VertexAttribPointer(*index, 1, gl.FLOAT, false, sizeOfOneInstance, *offset)
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4
		break
	case gohome.VALUE_VEC2:
		gl.VertexAttribPointer(*index, 2, gl.FLOAT, false, sizeOfOneInstance, *offset)
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2
		break
	case gohome.VALUE_VEC3:
		gl.VertexAttribPointer(*index, 3, gl.FLOAT, false, sizeOfOneInstance, *offset)
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 3
		break
	case gohome.VALUE_VEC4:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, sizeOfOneInstance, *offset)
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 4
		break
	case gohome.VALUE_MAT2:
		gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, sizeOfOneInstance, *offset)
		gl.VertexAttribDivisor(*index, 1)
		(*index)++
		(*offset) += 4 * 2 * 2
		break
	case gohome.VALUE_MAT3:
		for i := 0; i < 3; i++ {
			gl.VertexAttribPointer(*index, 3, gl.FLOAT, false, sizeOfOneInstance, *offset)
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 3
		}
		break
	case gohome.VALUE_MAT4:
		for i := 0; i < 4; i++ {
			gl.VertexAttribPointer(*index, 4, gl.FLOAT, false, sizeOfOneInstance, *offset)
			gl.VertexAttribDivisor(*index, 1)
			(*index)++
			(*offset) += 4 * 4
		}
		break
	}
}

func (this *WebGLInstancedMesh3D) instancedVertexAttribPointer(verticesSize int, indicesSize int, sizeOfOneInstance int) {
	offset := verticesSize + indicesSize
	var index = 4

	for i := 0; i < len(this.customValues); i++ {
		vertexAttribPointerForValueType(this.customValues[i], &offset, &index, sizeOfOneInstance)
	}
}

func enableValueType(valueType uint32, index *int) {
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

func (this *WebGLInstancedMesh3D) instancedEnableVertexAttribArray() {
	var i uint32
	var index = 4
	for i = 0; i < uint32(len(this.customValues)); i++ {
		enableValueType(this.customValues[i], &index)
	}
}

func (this *WebGLInstancedMesh3D) deleteElements() {
	this.vertices = append(this.vertices[:0], this.vertices[len(this.vertices):]...)
	this.indices = append(this.indices[:0], this.indices[len(this.indices):]...)
}

func (this *WebGLInstancedMesh3D) calculateOffsets() {
	var i uint32
	var offset = this.numVertices*gohome.MESH3DVERTEXSIZE + this.numIndices*2
	for i = 0; i < uint32(len(this.valueTypeIndexOffsets)); i++ {
		this.valueTypeIndexOffsets[i].offset = offset
		offset += getSize(this.valueTypeIndexOffsets[i].valueType)
	}
}

func (this *WebGLInstancedMesh3D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4))
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4+3*4))
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4+3*4+2*4))
	this.instancedVertexAttribPointer(this.numVertices*gohome.MESH3DVERTEXSIZE, this.numIndices*gohome.INDEX_SIZE, this.sizePerInstance)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)
	this.instancedEnableVertexAttribArray()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
}

func (this *WebGLInstancedMesh3D) Load() {
	if this.loaded {
		return
	}
	defer func() {
		this.loaded = true
	}()
	this.numVertices = len(this.vertices)
	this.numIndices = len(this.indices)

	if this.numVertices == 0 || this.numIndices == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "No vertices or indices have been added!")
		return
	}
	if this.numInstances == 0 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_WARNING, "InstancedMesh3D", this.Name, "Num Instances hasn't been set! Will be set to 1")
		this.numInstances = 1
	}

	var verticesSize = this.numVertices * gohome.MESH3DVERTEXSIZE
	var indicesSize = this.numIndices * 2
	if this.canUseInstanced {
		this.instancedSize = this.getInstancedSize() * this.numInstances
	}
	var bufferSize int
	if this.canUseInstanced {
		bufferSize = verticesSize + indicesSize + this.instancedSize
	} else {
		bufferSize = verticesSize + indicesSize
	}
	var usage int
	if this.canUseInstanced {
		usage = gl.DYNAMIC_DRAW
	} else {
		usage = gl.STATIC_DRAW
	}
	this.hasUV = true

	this.CalculateTangents()

	vertexBuffer := gohome.Mesh3DVerticesToFloatArray(this.vertices)

	if this.canUseVAOs {
		this.vao = gl.CreateVertexArray()
	}
	this.buffer = gl.CreateBuffer()
	handleWebGLError("InstancedMesh3D", this.Name, "GenBuffer: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, usage)
	handleWebGLError("InstancedMesh3D", this.Name, "BufferData: ")

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, vertexBuffer)
	handleWebGLError("InstancedMesh3D", this.Name, "BufferSubData Vertices: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), this.indices)
	handleWebGLError("InstancedMesh3D", this.Name, "BufferSubData Indices: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)

	if this.canUseInstanced {
		this.sizePerInstance = this.instancedSize / this.numInstances
	}
	this.numUsedInstances = this.numInstances

	if this.canUseVAOs {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(nil)
	}

	this.deleteElements()
	if this.canUseInstanced {
		this.calculateOffsets()
	}
}
func (this *WebGLInstancedMesh3D) Render() {
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
	if this.canUseVAOs {
		gl.BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}
	if this.canUseInstanced {
		gl.GetError()
		gl.DrawElementsInstanced(gl.TRIANGLES, this.numIndices, gl.UNSIGNED_SHORT, this.numVertices*gohome.MESH3DVERTEXSIZE, this.numUsedInstances)
		handleWebGLError("InstancedMesh3D", this.Name, "RenderError: ")
	} else {
		for i := 0; i < this.numUsedInstances && i < this.numInstances; i++ {
			this.setInstancedValuesUniforms(i)
			gl.DrawElements(gl.TRIANGLES, this.numIndices, gl.UNSIGNED_SHORT, this.numVertices*gohome.MESH3DVERTEXSIZE)
			handleWebGLError("InstancedMesh3D", this.Name, "RenderError: ")
		}
	}

	if this.canUseVAOs {
		gl.BindVertexArray(nil)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
	}
}

func (this *WebGLInstancedMesh3D) setInstancedValuesUniforms(instance int) {
	shader := gohome.RenderMgr.CurrentShader
	var ivt indexValueType
	if shader == nil {
		return
	}
	ivt.valueType = gohome.VALUE_FLOAT
	for i := 0; i < len(this.floats); i++ {
		ivt.index = i
		value := this.floats[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformF(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Float "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_VEC2
	for i := 0; i < len(this.vec2s); i++ {
		ivt.index = i
		value := this.vec2s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformV2(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Vec2 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_VEC3
	for i := 0; i < len(this.vec3s); i++ {
		ivt.index = i
		value := this.vec3s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformV3(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Vec3 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_VEC4
	for i := 0; i < len(this.vec4s); i++ {
		ivt.index = i
		value := this.vec4s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformV4(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Vec4 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_MAT2
	for i := 0; i < len(this.mat2s); i++ {
		ivt.index = i
		value := this.mat2s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformM2(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Mat2 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_MAT3
	for i := 0; i < len(this.mat3s); i++ {
		ivt.index = i
		value := this.mat3s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformM3(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Mat3 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
	ivt.valueType = gohome.VALUE_MAT4
	for i := 0; i < len(this.mat4s); i++ {
		ivt.index = i
		value := this.mat4s[i][instance]
		name, ok := this.namesForIndex[ivt]
		if ok {
			shader.SetUniformM4(name, value)
		} else {
			gohome.ErrorMgr.Error("InstancedMesh3D", this.Name, "No Name has been specified for Mat4 "+strconv.FormatUint(uint64(ivt.index), 10))
		}
	}
}

func (this *WebGLInstancedMesh3D) Terminate() {
	if this.canUseVAOs {
		gl.DeleteVertexArray(this.vao)
	}
	gl.DeleteBuffer(this.buffer)
}

func (this *WebGLInstancedMesh3D) SetMaterial(mat *gohome.Material) {
	this.Material = mat
}
func (this *WebGLInstancedMesh3D) GetMaterial() *gohome.Material {
	if this.Material == nil {
		this.Material = &gohome.Material{}
		this.Material.InitDefault()
	}
	return this.Material
}
func (this *WebGLInstancedMesh3D) GetName() string {
	return this.Name
}

func (this *WebGLInstancedMesh3D) GetNumVertices() uint32 {
	return uint32(this.numVertices)
}

func (this *WebGLInstancedMesh3D) GetNumIndices() uint32 {
	return uint32(this.numIndices)
}

func (this *WebGLInstancedMesh3D) recreateBuffer(numInstances int) {
	verticesSize := this.numVertices * gohome.MESH3DVERTEXSIZE
	indicesSize := this.numIndices * 2
	this.instancedSize = this.getInstancedSize() * numInstances
	bufferSize := verticesSize + indicesSize + this.instancedSize

	tempBuffer := gl.CreateBuffer()
	handleWebGLError("InstancedMesh3D", this.Name, "SetNumInstances GenBuffer: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)

	gl.BindBuffer(gl.ARRAY_BUFFER, tempBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, gl.DYNAMIC_DRAW)
	handleWebGLError("InstancedMesh3D", this.Name, "SetNumInstances BufferData: ")

	gl.CopyBufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.ARRAY_BUFFER, 0, 0, verticesSize+indicesSize)
	handleWebGLError("InstancedMesh3D", this.Name, "SetNumInstances CopyBufferSubData Vertices: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
	gl.DeleteBuffer(this.buffer)
	this.buffer = tempBuffer

	this.sizePerInstance = this.instancedSize / numInstances

	if this.canUseVAOs {
		gl.BindVertexArray(this.vao)
		this.attributePointer()
		gl.BindVertexArray(nil)
	}
	this.numInstances = numInstances
	this.calculateOffsets()
}

func (this *WebGLInstancedMesh3D) changeNumInstancesUniforms(n int) {
	if n > this.numInstances {
		for i := 0; i < len(this.floats); i++ {
			this.floats[i] = append(this.floats[i], make([]float32, n-this.numInstances)...)
		}
		for i := 0; i < len(this.vec2s); i++ {
			this.vec2s[i] = append(this.vec2s[i], make([]mgl32.Vec2, n-this.numInstances)...)
		}
		for i := 0; i < len(this.vec3s); i++ {
			this.vec3s[i] = append(this.vec3s[i], make([]mgl32.Vec3, n-this.numInstances)...)
		}
		for i := 0; i < len(this.vec4s); i++ {
			this.vec4s[i] = append(this.vec4s[i], make([]mgl32.Vec4, n-this.numInstances)...)
		}
		for i := 0; i < len(this.mat2s); i++ {
			this.mat2s[i] = append(this.mat2s[i], make([]mgl32.Mat2, n-this.numInstances)...)
		}
		for i := 0; i < len(this.mat3s); i++ {
			this.mat3s[i] = append(this.mat3s[i], make([]mgl32.Mat3, n-this.numInstances)...)
		}
		for i := 0; i < len(this.mat4s); i++ {
			this.mat4s[i] = append(this.mat4s[i], make([]mgl32.Mat4, n-this.numInstances)...)
		}
	} else {
		for i := 0; i < len(this.floats); i++ {
			this.floats[i] = this.floats[i][:n]
		}
		for i := 0; i < len(this.vec2s); i++ {
			this.vec2s[i] = this.vec2s[i][:n]
		}
		for i := 0; i < len(this.vec3s); i++ {
			this.vec3s[i] = this.vec3s[i][:n]
		}
		for i := 0; i < len(this.vec4s); i++ {
			this.vec4s[i] = this.vec4s[i][:n]
		}
		for i := 0; i < len(this.mat2s); i++ {
			this.mat2s[i] = this.mat2s[i][:n]
		}
		for i := 0; i < len(this.mat3s); i++ {
			this.mat3s[i] = this.mat3s[i][:n]
		}
		for i := 0; i < len(this.mat4s); i++ {
			this.mat4s[i] = this.mat4s[i][:n]
		}
	}
}

func (this *WebGLInstancedMesh3D) SetNumInstances(n uint32) {
	n1 := int(n)
	if this.numInstances != n1 {
		if this.loaded {
			if this.canUseInstanced {
				this.recreateBuffer(n1)
			} else {
				this.changeNumInstancesUniforms(n1)
			}
		}
		this.numInstances = n1
		this.numUsedInstances = n1
	}
}
func (this *WebGLInstancedMesh3D) GetNumInstances() uint32 {
	return uint32(this.numInstances)
}

func (this *WebGLInstancedMesh3D) addValueTypeIndexOffset(valueType uint32) {
	var maxIndex int = 0
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

func (this *WebGLInstancedMesh3D) addValueTypeIndexOffsetFront(valueType uint32) {
	var maxIndex int = 0
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

func (this *WebGLInstancedMesh3D) AddValueFront(valueType uint32) {
	if this.canUseInstanced {
		this.customValues = append(this.customValues, valueType)
		this.addValueTypeIndexOffsetFront(valueType)
	} else {
		this.AddValue(valueType)
	}
}

func (this *WebGLInstancedMesh3D) AddValue(valueType uint32) {
	if this.canUseInstanced {
		this.customValues = append(this.customValues, valueType)
		this.addValueTypeIndexOffset(valueType)
	} else {
		switch valueType {
		case gohome.VALUE_FLOAT:
			this.floats = append(this.floats, nil)
		case gohome.VALUE_VEC2:
			this.vec2s = append(this.vec2s, nil)
		case gohome.VALUE_VEC3:
			this.vec3s = append(this.vec3s, nil)
		case gohome.VALUE_VEC4:
			this.vec4s = append(this.vec4s, nil)
		case gohome.VALUE_MAT2:
			this.mat2s = append(this.mat2s, nil)
		case gohome.VALUE_MAT3:
			this.mat3s = append(this.mat3s, nil)
		case gohome.VALUE_MAT4:
			this.mat4s = append(this.mat4s, nil)
		}
	}
}

func (this *WebGLInstancedMesh3D) getOffset(valueType uint32, index int) int {
	for i := 0; i < len(this.valueTypeIndexOffsets); i++ {
		if this.valueTypeIndexOffsets[i].valueType == valueType && this.valueTypeIndexOffsets[i].index == index {
			return this.valueTypeIndexOffsets[i].offset
		}
	}

	return 0
}

func (this *WebGLInstancedMesh3D) SetF(index uint32, value []float32) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_FLOAT, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Float value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.floats[index] = value
	}

}
func (this *WebGLInstancedMesh3D) SetV2(index uint32, value []mgl32.Vec2) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_VEC2, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec2 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.vec2s[index] = value
	}
}

func (this *WebGLInstancedMesh3D) SetV3(index uint32, value []mgl32.Vec3) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_VEC3, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec3 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.vec3s[index] = value
	}
}
func (this *WebGLInstancedMesh3D) SetV4(index uint32, value []mgl32.Vec4) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_VEC4, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Vec4 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.vec4s[index] = value
	}
}
func (this *WebGLInstancedMesh3D) SetM2(index uint32, value []mgl32.Mat2) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_MAT2, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat2 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.mat2s[index] = value
	}
}
func (this *WebGLInstancedMesh3D) SetM3(index uint32, value []mgl32.Mat3) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_MAT3, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat3 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.mat3s[index] = value
	}

}
func (this *WebGLInstancedMesh3D) SetM4(index uint32, value []mgl32.Mat4) {
	if this.canUseInstanced {
		offset := this.getOffset(gohome.VALUE_MAT4, int(index))
		if offset == 0 {
			return
		}
		if len(value) < this.numInstances {
			gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_ERROR, "InstancedMesh3D", this.Name, "Mat4 value "+strconv.Itoa(int(index))+" is too small!")
			return
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)

		for i := 0; i < this.numInstances; i++ {
			gl.BufferSubData(gl.ARRAY_BUFFER, offset, value[i])
			offset += this.sizePerInstance
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, nil)
	} else {
		this.mat4s[index] = value
	}
}
func (this *WebGLInstancedMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return this.vertices
}
func (this *WebGLInstancedMesh3D) GetIndices() []uint32 {
	inds := make([]uint32, len(this.indices))
	for k, v := range this.indices {
		inds[k] = uint32(v)
	}
	return inds
}

func (this *WebGLInstancedMesh3D) SetName(index uint32, value_type uint32, value string) {
	if this.canUseInstanced {
		return
	}

	var ivt indexValueType
	ivt.index = int(index)
	ivt.valueType = value_type

	this.namesForIndex[ivt] = value
}

func (this *WebGLInstancedMesh3D) HasUV() bool {
	return this.hasUV
}
func (this *WebGLInstancedMesh3D) AABB() gohome.AxisAlignedBoundingBox {
	return this.aabb
}
func (this *WebGLInstancedMesh3D) Copy() gohome.Mesh3D {
	return nil
}

func (this *WebGLInstancedMesh3D) checkAABB() {
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

func (this *WebGLInstancedMesh3D) SetNumUsedInstances(n uint32) {
	this.numUsedInstances = int(n)
}

func (this *WebGLInstancedMesh3D) GetNumUsedInstances() uint32 {
	return uint32(this.numUsedInstances)
}

func (this *WebGLInstancedMesh3D) LoadedToGPU() bool {
	return this.loaded
}

func (this *WebGLRenderer) InstancedMesh3DFromLoadedMesh3D(mesh gohome.Mesh3D) gohome.InstancedMesh3D {
	oglmesh := mesh.(*WebGLMesh3D)
	ioglmesh := CreateWebGLInstancedMesh3D(oglmesh.Name)
	ioglmesh.numVertices = oglmesh.numVertices
	ioglmesh.numIndices = oglmesh.numIndices

	ioglmesh.numInstances = 0

	var verticesSize = ioglmesh.numVertices * gohome.MESH3DVERTEXSIZE
	var indicesSize = ioglmesh.numIndices * 2
	if ioglmesh.canUseInstanced {
		ioglmesh.instancedSize = 0
	}
	bufferSize := verticesSize + indicesSize
	var usage int
	if ioglmesh.canUseInstanced {
		usage = gl.DYNAMIC_DRAW
	} else {
		usage = gl.STATIC_DRAW
	}

	if ioglmesh.canUseVAOs {
		ioglmesh.vao = gl.CreateVertexArray()
	}
	ioglmesh.buffer = gl.CreateBuffer()
	handleWebGLError("InstancedMesh3D", ioglmesh.Name, "GenBuffer: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, ioglmesh.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, usage)
	handleWebGLError("InstancedMesh3D", ioglmesh.Name, "BufferData: ")

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, oglmesh.buffer)
	gl.CopyBufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.ARRAY_BUFFER, 0, 0, bufferSize)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	if ioglmesh.canUseInstanced {
		ioglmesh.sizePerInstance = 0
	}
	ioglmesh.numUsedInstances = 0

	if ioglmesh.canUseVAOs {
		gl.BindVertexArray(ioglmesh.vao)
		ioglmesh.attributePointer()
		gl.BindVertexArray(nil)
	}

	ioglmesh.deleteElements()
	if ioglmesh.canUseInstanced {
		ioglmesh.calculateOffsets()
	}
	ioglmesh.loaded = true

	ioglmesh.aabb = oglmesh.aabb
	ioglmesh.hasUV = oglmesh.hasUV
	ioglmesh.tangentsCalculated = true

	return ioglmesh
}
