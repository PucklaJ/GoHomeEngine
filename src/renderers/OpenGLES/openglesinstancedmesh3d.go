package renderer

import (
	// "fmt"
	"encoding/binary"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/gl"
	"log"
	"sync"
)

type indexValueTypePair struct {
	Index     uint32
	ValueType uint32
}

type OpenGLESInstancedMesh3D struct {
	vertices    []gohome.Mesh3DVertex
	indices     []uint32
	numVertices uint32
	numIndices  uint32

	vao gl.VertexArray
	vbo gl.Buffer
	ibo gl.Buffer

	Name     string
	Material *gohome.Material

	tangentsCalculated bool
	gles               *gl.Context
	isgles3            bool

	numInstances uint32
	floats       [][]float32
	vec2s        [][]mgl32.Vec2
	vec3s        [][]mgl32.Vec3
	vec4s        [][]mgl32.Vec4
	mat2s        [][]mgl32.Mat2
	mat3s        [][]mgl32.Mat3
	mat4s        [][]mgl32.Mat4
	names        map[indexValueTypePair]string
}

func (this *OpenGLESInstancedMesh3D) CalculateTangentsRoutine(startIndex, maxIndex uint32, wg *sync.WaitGroup) {
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

func (this *OpenGLESInstancedMesh3D) CalculateTangents() {
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

func (this *OpenGLESInstancedMesh3D) AddVertices(vertices []gohome.Mesh3DVertex, indices []uint32) {
	this.vertices = append(this.vertices, vertices...)
	this.indices = append(this.indices, indices...)
}

func CreateOpenGLESInstancedMesh3D(name string) *OpenGLESInstancedMesh3D {
	mesh := OpenGLESInstancedMesh3D{
		Name:               name,
		tangentsCalculated: false,
	}

	render, _ := gohome.Render.(*OpenGLESRenderer)
	mesh.gles = &render.gles
	_, mesh.isgles3 = (*mesh.gles).(gl.Context3)

	mesh.names = make(map[indexValueTypePair]string)

	return &mesh
}

func (this *OpenGLESInstancedMesh3D) deleteElements() {
	this.vertices = append(this.vertices[:0], this.vertices[len(this.vertices):]...)
	this.indices = append(this.indices[:0], this.indices[len(this.indices):]...)
}

func (this *OpenGLESInstancedMesh3D) toByteArrays() ([]byte, []byte) {
	var verticesBytes []byte
	var indicesBytes []byte

	verticesFloats := make([]float32, this.numVertices*MESH3DVERTEX_SIZE/4)
	var index uint32
	for i := 0; uint32(i) < this.numVertices; i++ {
		for j := 0; uint32(j) < MESH3DVERTEX_SIZE/4; j++ {
			verticesFloats[index+uint32(j)] = this.vertices[i][j]
		}
		index += MESH3DVERTEX_SIZE / 4
	}

	verticesBytes = f32.Bytes(binary.LittleEndian, verticesFloats...)

	indicesBytes = make([]byte, this.numIndices*gohome.INDEX_SIZE)
	for i := 0; uint32(i) < this.numIndices; i++ {
		binary.LittleEndian.PutUint32(indicesBytes[uint32(i)*gohome.INDEX_SIZE:uint32(i)*gohome.INDEX_SIZE+3+1], this.indices[i])
	}

	return verticesBytes, indicesBytes
}

func (this *OpenGLESInstancedMesh3D) attributePointer() {
	(*this.gles).BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	(*this.gles).VertexAttribPointer(gl.Attrib{0}, 3, gl.FLOAT, false, int(MESH3DVERTEX_SIZE), 0)
	(*this.gles).VertexAttribPointer(gl.Attrib{1}, 3, gl.FLOAT, false, int(MESH3DVERTEX_SIZE), 3*4)
	(*this.gles).VertexAttribPointer(gl.Attrib{2}, 2, gl.FLOAT, false, int(MESH3DVERTEX_SIZE), 3*4+3*4)
	(*this.gles).VertexAttribPointer(gl.Attrib{3}, 3, gl.FLOAT, false, int(MESH3DVERTEX_SIZE), 3*4+3*4+2*4)

	(*this.gles).EnableVertexAttribArray(gl.Attrib{0})
	(*this.gles).EnableVertexAttribArray(gl.Attrib{1})
	(*this.gles).EnableVertexAttribArray(gl.Attrib{2})
	(*this.gles).EnableVertexAttribArray(gl.Attrib{3})

	(*this.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.ibo)
}

func (this *OpenGLESInstancedMesh3D) Load() {

	this.numVertices = uint32(len(this.vertices))
	this.numIndices = uint32(len(this.indices))

	if this.numVertices == 0 || this.numIndices == 0 {
		log.Println("No vertices or indices have been added for mesh", this.Name, "!")
		return
	}

	this.CalculateTangents()
	verticesBytes, indicesBytes := this.toByteArrays()

	if this.isgles3 {
		this.vao = (*this.gles).CreateVertexArray()
	}
	this.vbo = (*this.gles).CreateBuffer()
	this.ibo = (*this.gles).CreateBuffer()

	(*this.gles).BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	(*this.gles).BufferData(gl.ARRAY_BUFFER, verticesBytes, gl.STATIC_DRAW)
	(*this.gles).BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{0})

	(*this.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.ibo)
	(*this.gles).BufferData(gl.ELEMENT_ARRAY_BUFFER, indicesBytes, gl.STATIC_DRAW)
	(*this.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{0})

	if this.isgles3 {
		(*this.gles).BindVertexArray(this.vao)
		this.attributePointer()
		(*this.gles).BindVertexArray(gl.VertexArray{0})
	}

	this.deleteElements()
}

func (this *OpenGLESInstancedMesh3D) setInstancedUniforms(shader gohome.Shader, instanceIndex uint32) {
	if shader == nil {
		return
	}

	for i := 0; i < len(this.floats); i++ {
		shader.SetUniformF(this.names[indexValueTypePair{uint32(i), gohome.VALUE_FLOAT}], this.floats[i][instanceIndex])
	}
	for i := 0; i < len(this.vec2s); i++ {
		shader.SetUniformV2(this.names[indexValueTypePair{uint32(i), gohome.VALUE_VEC2}], this.vec2s[i][instanceIndex])
	}
	for i := 0; i < len(this.vec3s); i++ {
		shader.SetUniformV3(this.names[indexValueTypePair{uint32(i), gohome.VALUE_VEC3}], this.vec3s[i][instanceIndex])
	}
	for i := 0; i < len(this.vec4s); i++ {
		shader.SetUniformV4(this.names[indexValueTypePair{uint32(i), gohome.VALUE_VEC4}], this.vec4s[i][instanceIndex])
	}
	for i := 0; i < len(this.mat2s); i++ {
		shader.SetUniformM2(this.names[indexValueTypePair{uint32(i), gohome.VALUE_MAT2}], this.mat2s[i][instanceIndex])
	}
	for i := 0; i < len(this.mat3s); i++ {
		shader.SetUniformM3(this.names[indexValueTypePair{uint32(i), gohome.VALUE_MAT3}], this.mat3s[i][instanceIndex])
	}
	for i := 0; i < len(this.mat4s); i++ {
		shader.SetUniformM4(this.names[indexValueTypePair{uint32(i), gohome.VALUE_MAT4}], this.mat4s[i][instanceIndex])
	}
}

func (this *OpenGLESInstancedMesh3D) Render() {
	if this.numVertices == 0 || this.numIndices == 0 {
		return
	}
	if gohome.RenderMgr.CurrentShader != nil && this.Material != nil {
		if err := gohome.RenderMgr.CurrentShader.SetUniformMaterial(*this.Material); err != nil {
			// fmt.Println("Error:", err)
		}
	}
	if this.isgles3 {
		(*this.gles).BindVertexArray(this.vao)
	} else {
		this.attributePointer()
	}
	for i := 0; uint32(i) < this.numInstances; i++ {
		this.setInstancedUniforms(gohome.RenderMgr.CurrentShader, uint32(i))
		(*this.gles).DrawElements(gl.TRIANGLES, int(this.numIndices), gl.UNSIGNED_INT, 0)
	}
	if this.isgles3 {
		(*this.gles).BindVertexArray(gl.VertexArray{0})
	} else {
		(*this.gles).BindBuffer(gl.ARRAY_BUFFER, gl.Buffer{0})
		(*this.gles).BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.Buffer{0})
	}
}

func (this *OpenGLESInstancedMesh3D) Terminate() {
	if this.isgles3 {
		defer (*this.gles).DeleteVertexArray(this.vao)
	}
	defer (*this.gles).DeleteBuffer(this.vbo)
	defer (*this.gles).DeleteBuffer(this.ibo)
}

func (this *OpenGLESInstancedMesh3D) SetMaterial(mat *gohome.Material) {
	this.Material = mat
}

func (this *OpenGLESInstancedMesh3D) GetMaterial() *gohome.Material {
	if this.Material == nil {
		this.Material = &gohome.Material{}
	}
	return this.Material
}

func (this *OpenGLESInstancedMesh3D) GetNumVertices() uint32 {
	return this.numVertices
}
func (this *OpenGLESInstancedMesh3D) GetNumIndices() uint32 {
	return this.numIndices
}

func (this *OpenGLESInstancedMesh3D) GetVertices() []gohome.Mesh3DVertex {
	return this.vertices
}
func (this *OpenGLESInstancedMesh3D) GetIndices() []uint32 {
	return this.indices
}

func (this *OpenGLESInstancedMesh3D) GetName() string {
	return this.Name
}

func (this *OpenGLESInstancedMesh3D) increaseSliceSize(num uint32) {
	for i := 0; i < len(this.floats); i++ {
		this.floats[i] = append(this.floats[i], make([]float32, num)...)
	}
	for i := 0; i < len(this.vec2s); i++ {
		this.vec2s[i] = append(this.vec2s[i], make([]mgl32.Vec2, num)...)
	}
	for i := 0; i < len(this.vec3s); i++ {
		this.vec3s[i] = append(this.vec3s[i], make([]mgl32.Vec3, num)...)
	}
	for i := 0; i < len(this.vec4s); i++ {
		this.vec4s[i] = append(this.vec4s[i], make([]mgl32.Vec4, num)...)
	}
	for i := 0; i < len(this.mat2s); i++ {
		this.mat2s[i] = append(this.mat2s[i], make([]mgl32.Mat2, num)...)
	}
	for i := 0; i < len(this.mat3s); i++ {
		this.mat3s[i] = append(this.mat3s[i], make([]mgl32.Mat3, num)...)
	}
	for i := 0; i < len(this.mat4s); i++ {
		this.mat4s[i] = append(this.mat4s[i], make([]mgl32.Mat4, num)...)
	}
}

func (this *OpenGLESInstancedMesh3D) decreaseSliceSize(num uint32) {
	for i := 0; i < len(this.floats); i++ {
		if uint32(len(this.floats[i])) == this.numInstances {
			this.floats[i] = this.floats[i][:uint32(len(this.floats))-num]
		}
	}
	for i := 0; i < len(this.vec2s); i++ {
		if uint32(len(this.vec2s[i])) == this.numInstances {
			this.vec2s[i] = this.vec2s[i][:uint32(len(this.vec2s))-num]
		}
	}
	for i := 0; i < len(this.vec3s); i++ {
		if uint32(len(this.vec3s[i])) == this.numInstances {
			this.vec3s[i] = this.vec3s[i][:uint32(len(this.vec3s))-num]
		}
	}
	for i := 0; i < len(this.vec4s); i++ {
		if uint32(len(this.vec4s[i])) == this.numInstances {
			this.vec4s[i] = this.vec4s[i][:uint32(len(this.vec4s))-num]
		}
	}
	for i := 0; i < len(this.mat2s); i++ {
		if uint32(len(this.mat2s[i])) == this.numInstances {
			this.mat2s[i] = this.mat2s[i][:uint32(len(this.mat2s))-num]
		}
	}
	for i := 0; i < len(this.mat3s); i++ {
		if uint32(len(this.mat3s[i])) == this.numInstances {
			this.mat3s[i] = this.mat3s[i][:uint32(len(this.mat3s))-num]
		}
	}
	for i := 0; i < len(this.mat4s); i++ {
		if uint32(len(this.mat4s[i])) == this.numInstances {
			this.mat4s[i] = this.mat4s[i][:uint32(len(this.mat4s))-num]
		}
	}
}

func (this *OpenGLESInstancedMesh3D) SetNumInstances(n uint32) {
	if n > this.numInstances {
		this.increaseSliceSize(n - this.numInstances)
	} else if n < this.numInstances {
		this.decreaseSliceSize(this.numInstances - n)
	}
	this.numInstances = n
}
func (this *OpenGLESInstancedMesh3D) GetNumInstances() uint32 {
	return this.numInstances
}
func (this *OpenGLESInstancedMesh3D) AddValue(valueType uint32) {
	switch valueType {
	case gohome.VALUE_FLOAT:
		this.floats = append(this.floats, []float32{})
	case gohome.VALUE_VEC2:
		this.vec2s = append(this.vec2s, []mgl32.Vec2{})
	case gohome.VALUE_VEC3:
		this.vec3s = append(this.vec3s, []mgl32.Vec3{})
	case gohome.VALUE_VEC4:
		this.vec4s = append(this.vec4s, []mgl32.Vec4{})
	case gohome.VALUE_MAT2:
		this.mat2s = append(this.mat2s, []mgl32.Mat2{})
	case gohome.VALUE_MAT3:
		this.mat3s = append(this.mat3s, []mgl32.Mat3{})
	case gohome.VALUE_MAT4:
		this.mat4s = append(this.mat4s, []mgl32.Mat4{})
	}
}
func (this *OpenGLESInstancedMesh3D) SetF(index uint32, value []float32) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Float value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Float value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.floats[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetV2(index uint32, value []mgl32.Vec2) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec2 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Vec2 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.vec2s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetV3(index uint32, value []mgl32.Vec3) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec3 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Vec3 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.vec3s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetV4(index uint32, value []mgl32.Vec4) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Vec4 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Vec4 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.vec4s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetM2(index uint32, value []mgl32.Mat2) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat2 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Mat2 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.mat2s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetM3(index uint32, value []mgl32.Mat3) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat3 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Mat3 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.mat3s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetM4(index uint32, value []mgl32.Mat4) {
	if uint32(len(value)) < this.numInstances {
		log.Println("Mat4 value", index, "of instanced mesh", this.Name, "is too small!")
		return
	} else if uint32(len(value)) > this.numInstances {
		log.Println("Mat4 value", index, "of instanced mesh", this.Name, "is too big! Using", this.numInstances, "values")
	}
	this.mat4s[index] = value[:this.numInstances]
}
func (this *OpenGLESInstancedMesh3D) SetName(index uint32, value_type uint32, value string) {
	this.names[indexValueTypePair{index, value_type}] = value
}
