package renderer

import (
	"strconv"
	"sync"
	"unsafe"

	"github.com/PucklaJ/GoHomeEngine/src/gohome"
	gl "github.com/PucklaJ/android-go/gles2"
	"github.com/PucklaJ/mathgl/mgl32"
)

type valueTypeIndexOffset struct {
	valueType int
	index     int
	offset    int
}

type indexValueType struct {
	index     int
	valueType int
}

type OpenGLES2InstancedMesh3D struct {
	vertices         []gohome.Mesh3DVertex
	indices          []uint32
	numVertices      int
	numIndices       int
	numInstances     int
	numUsedInstances int

	buffer uint32
	hasUV  bool
	loaded bool
	aabb   gohome.AxisAlignedBoundingBox

	Name     string
	Material *gohome.Material

	tangentsCalculated bool
	customValues       []int

	floats        [][]float32
	vec2s         [][]mgl32.Vec2
	vec3s         [][]mgl32.Vec3
	vec4s         [][]mgl32.Vec4
	mat2s         [][]mgl32.Mat2
	mat3s         [][]mgl32.Mat3
	mat4s         [][]mgl32.Mat4
	namesForIndex map[indexValueType]string
}

func CreateOpenGLES2InstancedMesh3D(name string) *OpenGLES2InstancedMesh3D {
	mesh := &OpenGLES2InstancedMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}
	mesh.namesForIndex = make(map[indexValueType]string)

	return mesh
}

func (this *OpenGLES2InstancedMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	this.vertices = append(this.vertices, vertices...)
	this.indices = append(this.indices, indices...)
	this.checkAABB()
}

func (this *OpenGLES2InstancedMesh3D) CalculateTangentsRoutine(startIndex, maxIndex int, wg *sync.WaitGroup) {
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
	for i := startIndex; i < maxIndex && i < len(indices); i += 3 {
		if i > len(indices)-3 {
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
		for j := 0; j < 3; j++ {
			(*vertices)[indices[i+j]][8] = tangent[0]
			(*vertices)[indices[i+j]][9] = tangent[1]
			(*vertices)[indices[i+j]][10] = tangent[2]
		}
	}
}

func (this *OpenGLES2InstancedMesh3D) CalculateTangents() {
	if this.tangentsCalculated {
		return
	}
	var wg sync.WaitGroup

	deltaIndex := len(this.indices) / NUM_GO_ROUTINES_TANGENTS_CALCULATING
	if deltaIndex == 0 {
		deltaIndex = len(this.indices) / 3
	}
	if deltaIndex > 3 {
		deltaIndex -= deltaIndex % 3
	} else {
		deltaIndex = 3
	}

	for i := 0; i < NUM_GO_ROUTINES_TANGENTS_CALCULATING*2; i++ {
		wg.Add(1)
		go this.CalculateTangentsRoutine(i*deltaIndex, i*deltaIndex+deltaIndex, &wg)
		if i*deltaIndex+deltaIndex >= len(this.indices) {
			break
		}
	}

	wg.Wait()

	this.tangentsCalculated = true
}

func getSize(valueType int) int {
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

func (this *OpenGLES2InstancedMesh3D) deleteElements() {
	this.vertices = append(this.vertices[:0], this.vertices[len(this.vertices):]...)
	this.indices = append(this.indices[:0], this.indices[len(this.indices):]...)
}

func (this *OpenGLES2InstancedMesh3D) attributePointer() {
	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4))
	gl.VertexAttribPointer(2, 2, gl.FLOAT, gl.FALSE, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4+3*4))
	gl.VertexAttribPointer(3, 3, gl.FLOAT, gl.FALSE, gohome.MESH3DVERTEXSIZE, gl.PtrOffset(3*4+3*4+2*4))

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
}

func (this *OpenGLES2InstancedMesh3D) Load() {
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

	verticesSize := this.numVertices * gohome.MESH3DVERTEXSIZE
	indicesSize := this.numIndices * 2
	bufferSize := verticesSize + indicesSize
	var usage uint32
	usage = gl.STATIC_DRAW

	this.hasUV = true
	this.CalculateTangents()

	var buf [1]uint32
	gl.GenBuffers(1, buf[:])
	this.buffer = buf[0]
	handleOpenGLError("InstancedMesh3D", this.Name, "GenBuffer: ")

	gl.BindBuffer(gl.ARRAY_BUFFER, this.buffer)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, nil, usage)
	handleOpenGLError("InstancedMesh3D", this.Name, "BufferData: ")

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(verticesSize), unsafe.Pointer(&this.vertices[0]))
	handleOpenGLError("InstancedMesh3D", this.Name, "BufferSubData Vertices: ")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.buffer)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, int(verticesSize), int(indicesSize), unsafe.Pointer(&this.indices[0]))
	handleOpenGLError("InstancedMesh3D", this.Name, "BufferSubData Indices: ")
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	this.numUsedInstances = this.numInstances

	this.deleteElements()
}
func (this *OpenGLES2InstancedMesh3D) Render() {
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
		}
		gohome.RenderMgr.CurrentShader.SetUniformMaterial(*this.Material)
	}
	this.attributePointer()
	for i := 0; i < this.numUsedInstances && i < this.numInstances; i++ {
		this.setInstancedValuesUniforms(i)
		gl.DrawElements(gl.TRIANGLES, int32(this.numIndices), gl.UNSIGNED_SHORT, gl.PtrOffset(this.numVertices*gohome.MESH3DVERTEXSIZE))
		handleOpenGLError("InstancedMesh3D", this.Name, "RenderError: ")
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (this *OpenGLES2InstancedMesh3D) setInstancedValuesUniforms(instance int) {
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

func (this *OpenGLES2InstancedMesh3D) Terminate() {
	var buf [1]uint32
	buf[0] = this.buffer
	defer gl.DeleteBuffers(1, buf[:])
}

func (this *OpenGLES2InstancedMesh3D) SetMaterial(mat *gohome.Material) {
	this.Material = mat
}
func (this *OpenGLES2InstancedMesh3D) GetMaterial() *gohome.Material {
	if this.Material == nil {
		this.Material = &gohome.Material{}
		this.Material.InitDefault()
	}
	return this.Material
}
func (this *OpenGLES2InstancedMesh3D) GetName() string {
	return this.Name
}

func (this *OpenGLES2InstancedMesh3D) GetNumVertices() int {
	return this.numVertices
}

func (this *OpenGLES2InstancedMesh3D) GetNumIndices() int {
	return this.numIndices
}

func (this *OpenGLES2InstancedMesh3D) changeNumInstancesUniforms(n int) {
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

func (this *OpenGLES2InstancedMesh3D) SetNumInstances(n int) {
	if this.numInstances != n {
		if this.loaded {
			this.changeNumInstancesUniforms(n)
		}
		this.numInstances = n
		this.numUsedInstances = n
	}
}
func (this *OpenGLES2InstancedMesh3D) GetNumInstances() int {
	return this.numInstances
}

func (this *OpenGLES2InstancedMesh3D) AddValueFront(valueType int) {
	this.AddValue(valueType)
}

func (this *OpenGLES2InstancedMesh3D) AddValue(valueType int) {
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

func (this *OpenGLES2InstancedMesh3D) SetF(index int, value []float32) {
	this.floats[index] = value
}
func (this *OpenGLES2InstancedMesh3D) SetV2(index int, value []mgl32.Vec2) {
	this.vec2s[index] = value
}

func (this *OpenGLES2InstancedMesh3D) SetV3(index int, value []mgl32.Vec3) {
	this.vec3s[index] = value
}
func (this *OpenGLES2InstancedMesh3D) SetV4(index int, value []mgl32.Vec4) {
	this.vec4s[index] = value
}
func (this *OpenGLES2InstancedMesh3D) SetM2(index int, value []mgl32.Mat2) {
	this.mat2s[index] = value
}
func (this *OpenGLES2InstancedMesh3D) SetM3(index int, value []mgl32.Mat3) {
	this.mat3s[index] = value
}
func (this *OpenGLES2InstancedMesh3D) SetM4(index int, value []mgl32.Mat4) {
	this.mat4s[index] = value
}
func (this *OpenGLES2InstancedMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return this.vertices
}
func (this *OpenGLES2InstancedMesh3D) GetIndices() []uint32 {
	return this.indices
}

func (this *OpenGLES2InstancedMesh3D) SetName(index, value_type int, value string) {
	var ivt indexValueType
	ivt.index = index
	ivt.valueType = value_type

	this.namesForIndex[ivt] = value
}

func (this *OpenGLES2InstancedMesh3D) HasUV() bool {
	return this.hasUV
}
func (this *OpenGLES2InstancedMesh3D) AABB() gohome.AxisAlignedBoundingBox {
	return this.aabb
}
func (this *OpenGLES2InstancedMesh3D) Copy() gohome.Mesh3D {
	return nil
}

func (this *OpenGLES2InstancedMesh3D) checkAABB() {
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

func (this *OpenGLES2InstancedMesh3D) SetNumUsedInstances(n int) {
	this.numUsedInstances = n
}

func (this *OpenGLES2InstancedMesh3D) GetNumUsedInstances() int {
	return this.numUsedInstances
}

func (this *OpenGLES2InstancedMesh3D) LoadedToGPU() bool {
	return this.loaded
}

func (this *OpenGLES2Renderer) InstancedMesh3DFromLoadedMesh3D(mesh gohome.Mesh3D) gohome.InstancedMesh3D {
	oglmesh := mesh.(*OpenGLES2Mesh3D)
	ioglmesh := CreateOpenGLES2InstancedMesh3D(oglmesh.Name)
	ioglmesh.numVertices = oglmesh.numVertices
	ioglmesh.numIndices = oglmesh.numIndices

	ioglmesh.numInstances = 0

	ioglmesh.buffer = oglmesh.buffer

	ioglmesh.numUsedInstances = 0

	ioglmesh.deleteElements()
	ioglmesh.loaded = true

	ioglmesh.aabb = oglmesh.aabb
	ioglmesh.hasUV = oglmesh.hasUV
	ioglmesh.tangentsCalculated = true

	return ioglmesh
}
