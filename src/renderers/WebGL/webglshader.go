package renderer

import (
	"bytes"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"github.com/gopherjs/gopherjs/js"
	"strconv"
)

type WebGLShader struct {
	program             *js.Object
	name                string
	shaders             [6]*js.Object
	uniform_locations   map[string]*js.Object
	attribute_locations map[string]uint32
	attribute_sizes     map[string]int
	validated           bool
}

func CreateWebGLShader(name string) (*WebGLShader, error) {
	shader := &WebGLShader{
		program:             nil,
		name:                name,
		shaders:             [6]*js.Object{nil, nil, nil, nil, nil, nil},
		uniform_locations:   make(map[string]*js.Object),
		attribute_locations: make(map[string]uint32),
		validated:           false,
	}
	program := gl.CreateProgram()
	if program == js.Undefined {
		return shader, &WebGLError{errorString: "Couldn't create program"}
	} else {
		shader.program = program
		return shader, nil
	}
}

func getShaderTypeName(shader_type uint8) string {
	var shader_type_name string
	switch shader_type {
	case gohome.VERTEX:
		shader_type_name = "Vertex Shader"
	case gohome.FRAGMENT:
		shader_type_name = "Fragment Shader"
	case gohome.GEOMETRY:
		shader_type_name = "Geometry Shader"
	case gohome.TESSELLETION:
		shader_type_name = "Tesselletion Shader"
	case gohome.EVELUATION:
		shader_type_name = "Eveluation Shader"
	case gohome.COMPUTE:
		shader_type_name = "Compute Shader"
	}

	return shader_type_name
}

func toGohomeShaderType(shader_type int) uint8 {
	switch shader_type {
	case gl.VERTEX_SHADER:
		return gohome.VERTEX
	case gl.FRAGMENT_SHADER:
		return gohome.FRAGMENT
	}

	return 255
}

func toShaderTypeName(shader_type int) string {
	return getShaderTypeName(toGohomeShaderType(shader_type))
}

func getAttributeSizeForType(atype string) int {
	switch atype {
	case "mat3":
		return 3
	case "mat4":
		return 4
	}

	return 1
}

func (s *WebGLShader) getAttributeNames(program *js.Object, src string) []string {
	var line bytes.Buffer
	var lineString string
	var attributeNames []string
	var curChar byte = ' '
	var curIndex uint32 = 0
	var curWordIndex uint32 = 0
	var curWord uint32 = 0
	var wordBuffer bytes.Buffer
	var wordsString []string
	var readWord bool = false
	var version uint32 = 0

	s.attribute_sizes = make(map[string]int)

	for curIndex < uint32(len(src)) {
		for curChar = ' '; curChar != '\n' && curChar != 13; curChar = src[curIndex] {
			line.WriteByte(curChar)
			curIndex++
			if curIndex == uint32(len(src)) {
				break
			}
		}

		lineString = line.String()
		readWord = false
		curWord = 0
		for curWordIndex = 0; curWordIndex < uint32(len(lineString)); curWordIndex++ {
			curChar = lineString[curWordIndex]
			if curChar == ' ' || curChar == '\t' {
				if readWord {
					wordsString[curWord] = wordBuffer.String()
					wordBuffer.Reset()
					curWord++
					readWord = false
				}
			} else {
				if !readWord {
					readWord = true
					wordsString = append(wordsString, string(' '))
				}
				wordBuffer.WriteByte(curChar)
			}
		}
		if readWord {
			wordsString[curWord] = wordBuffer.String()
		}
		wordBuffer.Reset()
		line.Reset()
		if len(wordsString) >= 2 && wordsString[0] == "#version" {
			versionInt, _ := strconv.Atoi(wordsString[1])
			version = uint32(versionInt)
		}
		if len(wordsString) >= 2 && wordsString[0] == "void" && wordsString[1] == "main()" {
			break
		} else if len(wordsString) >= 3 {
			if (version >= 130 && wordsString[0] == "in") || (version <= 120 && wordsString[0] == "attribute") {
				if wordsString[1] == "highp" || wordsString[1] == "mediump" || wordsString[1] == "lowp" {
					wordsString = append(wordsString[:1], wordsString[2:]...)
				}
				if wordsString[2][len(wordsString[2])-1] == ';' {
					wordsString[2] = wordsString[2][:len(wordsString[2])-1]
				}
				attributeNames = append(attributeNames, wordsString[2])
				s.attribute_sizes[wordsString[2]] = getAttributeSizeForType(wordsString[1])
			}
		}
		wordsString = append(wordsString[len(wordsString):], wordsString[:0]...)
	}

	return attributeNames
}

func (s *WebGLShader) bindAttributesFromFile(program *js.Object, src string) {
	attributeNames := s.getAttributeNames(program, src)

	var index int = 0
	for i := 0; i < len(attributeNames); i++ {
		gl.BindAttribLocation(program, index, attributeNames[i]+"\x00")
		index += s.attribute_sizes[attributeNames[i]]
	}
}

func (s *WebGLShader) compileWebGLShader(shader_name string, shader_type int, src string, program *js.Object) (*js.Object, error) {
	shader := gl.CreateShader(shader_type)
	gl.ShaderSource(shader, src)
	gl.CompileShader(shader)

	status := gl.GetShaderiv(shader, gl.COMPILE_STATUS)
	if int(status[0]) == gl.FALSE {
		logText := gl.GetShaderInfoLog(shader)

		return nil, &WebGLError{errorString: logText}
	}
	gl.GetError()
	gl.AttachShader(program, shader)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return nil, &WebGLError{errorString: "Couldn't attach " + toShaderTypeName(shader_type) + " of " + shader_name + ": ErrorCode: " + strconv.Itoa(int(err))}
	}
	if shader_type == gl.VERTEX_SHADER {
		s.bindAttributesFromFile(program, src)
	}

	return shader, nil
}

func (s *WebGLShader) AddShader(shader_type uint8, src string) error {
	if shader_type == gohome.GEOMETRY {
		render, _ := gohome.Render.(*WebGLRenderer)
		if !render.HasFunctionAvailable("GEOMETRY_SHADER") {
			return &WebGLError{errorString: "Geometry shaders are not supported by this implementation"}
		}
	}

	var err error
	var shaderName *js.Object
	switch shader_type {
	case gohome.VERTEX:
		shaderName, err = s.compileWebGLShader(s.name, gl.VERTEX_SHADER, src, s.program)
	case gohome.FRAGMENT:
		shaderName, err = s.compileWebGLShader(s.name, gl.FRAGMENT_SHADER, src, s.program)
	}

	if err != nil {
		return &WebGLError{errorString: "Couldn't compile " + getShaderTypeName(shader_type) + ": " + err.Error()}
	}

	s.shaders[shader_type] = shaderName

	return nil
}

func (s *WebGLShader) deleteAllShaders() {
	for i := 0; i < 6; i++ {
		if s.shaders[i] != js.Undefined {
			gl.DetachShader(s.program, s.shaders[i])
			gl.DeleteShader(s.shaders[i])
		}
	}
}

func (s *WebGLShader) Link() error {
	defer s.deleteAllShaders()

	gl.GetError()
	gl.LinkProgram(s.program)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return &WebGLError{errorString: "Couldn't link: ErrorCode: " + strconv.Itoa(int(err))}
	}

	gl.GetError()
	status := gl.GetProgrami(s.program, gl.LINK_STATUS)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return &WebGLError{errorString: "Couldn't link: Couldn't get link status: ErrorCode: " + strconv.Itoa(int(err))}
	}
	if status == gl.FALSE {
		gl.GetError()
		logtext := gl.GetProgramInfoLog(s.program)
		if err := gl.GetError(); err != gl.NO_ERROR {
			return &WebGLError{errorString: "Couldn't link: Couldn't get info log: ErrorCode: " + strconv.Itoa(int(err))}
		}

		return &WebGLError{errorString: "Couldn't link: " + logtext}
	}

	return nil
}

func (s *WebGLShader) Use() {
	gl.GetError()
	gl.UseProgram(s.program)
	handleWebGLError("Shader", s.name, "glUseProgram")
}

func (s *WebGLShader) Unuse() {
	gl.GetError()
	gl.UseProgram(nil)
	handleWebGLError("Shader", s.name, "glUseProgram with 0")
}

func (s *WebGLShader) Setup() error {
	return s.validate()
}

func (s *WebGLShader) Terminate() {
	gl.DeleteProgram(s.program)
}

func (s *WebGLShader) getUniformLocation(name string) *js.Object {
	var loc *js.Object
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		gl.GetError()
		loc = gl.GetUniformLocation(s.program, name+"\x00")
		handleWebGLError("Shader", s.name, "glGetUniformLocation")
		s.uniform_locations[name] = loc
	}
	if loc == js.Undefined {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_WARNING, "Shader", s.name, "Couldn't find uniform "+name)
	}
	return loc
}

func (s *WebGLShader) SetUniformV2(name string, value mgl32.Vec2) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform2f(loc, value[0], value[1])
		handleWebGLError("Shader", s.name, "glUniform2f")
	}
}
func (s *WebGLShader) SetUniformV3(name string, value mgl32.Vec3) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform3f(loc, value[0], value[1], value[2])
		handleWebGLError("Shader", s.name, "glUniform3f")
	}
}
func (s *WebGLShader) SetUniformV4(name string, value mgl32.Vec4) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform4f(loc, value[0], value[1], value[2], value[3])
		handleWebGLError("Shader", s.name, "glUniform4f")
	}
}
func (s *WebGLShader) SetUniformIV2(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform2i(loc, int(value[0]), int(value[1]))
		handleWebGLError("Shader", s.name, "glUniform2i")
	}
}
func (s *WebGLShader) SetUniformIV3(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform3i(loc, int(value[0]), int(value[1]), int(value[2]))
		handleWebGLError("Shader", s.name, "glUniform3i")
	}
}
func (s *WebGLShader) SetUniformIV4(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform4i(loc, int(value[0]), int(value[1]), int(value[2]), int(value[3]))
		handleWebGLError("Shader", s.name, "glUniform4i")
	}
}
func (s *WebGLShader) SetUniformF(name string, value float32) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform1f(loc, value)
		handleWebGLError("Shader", s.name, "glUniform1f")
	}
}
func (s *WebGLShader) SetUniformI(name string, value int32) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.Uniform1i(loc, int(value))
		handleWebGLError("Shader", s.name, "glUniform1i")
	}
}
func (s *WebGLShader) SetUniformUI(name string, value uint32) {
	s.SetUniformI(name, int32(value))
}
func (s *WebGLShader) SetUniformB(name string, value uint8) {
	s.SetUniformI(name, int32(value))
}
func (s *WebGLShader) SetUniformM2(name string, value mgl32.Mat2) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.UniformMatrix2fv(loc, false, value[:])
		handleWebGLError("Shader", s.name, "glUniformMatrix2fv")
	}
}
func (s *WebGLShader) SetUniformM3(name string, value mgl32.Mat3) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.UniformMatrix3fv(loc, false, value[:])
		handleWebGLError("Shader", s.name, "glUniformMatrix3fv")
	}
}
func (s *WebGLShader) SetUniformM4(name string, value mgl32.Mat4) {
	loc := s.getUniformLocation(name)
	if loc != js.Undefined {
		gl.GetError()
		gl.UniformMatrix4fv(loc, false, value[:])
		handleWebGLError("Shader", s.name, "glUniformMatrix4fv")
	}
}

func (s *WebGLShader) SetUniformMaterial(mat gohome.Material) {
	var diffBind int32 = 0
	var specBind int32 = 0
	var normBind int32 = 0
	var boundTextures uint32

	maxtextures := gohome.Render.GetMaxTextures()

	if mat.DiffuseTexture != nil {
		diffBind = int32(gohome.Render.NextTextureUnit())
		if diffBind >= maxtextures-1 {
			diffBind = 0
			mat.DiffuseTextureLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.DiffuseTexture.Bind(uint32(diffBind))
			mat.DiffuseTextureLoaded = 1
			boundTextures++
		}
	} else {
		mat.DiffuseTextureLoaded = 0
	}

	if mat.SpecularTexture != nil {
		specBind = int32(gohome.Render.NextTextureUnit())
		if specBind >= maxtextures-1 {
			specBind = 0
			mat.SpecularTextureLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.SpecularTexture.Bind(uint32(specBind))
			mat.SpecularTextureLoaded = 1
			boundTextures++
		}
	} else {
		mat.SpecularTextureLoaded = 0
	}

	if mat.NormalMap != nil {
		normBind = int32(gohome.Render.NextTextureUnit())
		if normBind >= maxtextures-1 {
			normBind = 0
			mat.NormalMapLoaded = 0
			gohome.Render.DecrementTextureUnit(1)
		} else {
			mat.NormalMap.Bind(uint32(normBind))
			mat.NormalMapLoaded = 1
			boundTextures++
		}
	} else {
		mat.NormalMapLoaded = 0
	}

	gohome.Render.DecrementTextureUnit(boundTextures)

	s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.DiffuseColor))
	s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.SpecularColor))
	s.SetUniformF(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SHINYNESS_UNIFORM_NAME, mat.Shinyness)

	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME, mat.DiffuseTextureLoaded)
	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME, mat.SpecularTextureLoaded)
	s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME, mat.NormalMapLoaded)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME, diffBind)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME, specBind)
	s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+gohome.MATERIAL_NORMALMAP_UNIFORM_NAME, normBind)

	s.SetUniformF(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_TRANSPARENCY_UNIFORM_NAME, mat.Transparency)
}

func (s *WebGLShader) SetUniformLights(lightCollectionIndex int32) {
	if lightCollectionIndex == -1 || lightCollectionIndex > int32(len(gohome.LightMgr.LightCollections)-1) {
		s.SetUniformI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, 0)
		s.SetUniformI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, 0)
		s.SetUniformI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, 0)

		s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, mgl32.Vec3{1.0, 1.0, 1.0})
		return
	}

	lightColl := gohome.LightMgr.LightCollections[lightCollectionIndex]

	s.SetUniformI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, int32(len(lightColl.PointLights)))
	s.SetUniformI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, int32(len(lightColl.DirectionalLights)))
	s.SetUniformI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, int32(len(lightColl.SpotLights)))

	s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, gohome.ColorToVec3(lightColl.AmbientLight))

	var i uint32
	for i = 0; i < uint32(len(lightColl.PointLights)); i++ {
		lightColl.PointLights[i].SetUniforms(s, i)
	}
	for i = 0; i < uint32(len(lightColl.DirectionalLights)); i++ {
		lightColl.DirectionalLights[i].SetUniforms(s, i)
	}
	for i = 0; i < uint32(len(lightColl.SpotLights)); i++ {
		lightColl.SpotLights[i].SetUniforms(s, i)
	}
}

func (s *WebGLShader) GetName() string {
	return s.name
}

func (s *WebGLShader) validate() error {
	if s.validated {
		return nil
	}
	s.validated = true
	gl.ValidateProgram(s.program)
	status := gl.GetProgrami(s.program, gl.VALIDATE_STATUS)
	if status == gl.FALSE {
		logtext := gl.GetProgramInfoLog(s.program)
		s.validated = false
		return &WebGLError{"Couldn't validate: " + logtext}
	}

	return nil
}

func (s *WebGLShader) AddAttribute(name string, location uint32) {
	gl.BindAttribLocation(s.program, int(location), name)
}
