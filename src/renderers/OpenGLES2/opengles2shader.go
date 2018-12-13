package renderer

import (
	"bytes"
	"github.com/PucklaMotzer09/GoHomeEngine/src/gohome"
	gl "github.com/PucklaMotzer09/android-go/gles2"
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"strconv"
)

type OpenGLES2Shader struct {
	program             uint32
	name                string
	shaders             [6]uint32
	uniform_locations   map[string]int32
	attribute_locations map[string]uint32
	attribute_sizes     map[string]uint32
	validated           bool
}

func CreateOpenGLES2Shader(name string) (*OpenGLES2Shader, error) {
	shader := &OpenGLES2Shader{
		program:             0,
		name:                name,
		shaders:             [6]uint32{0, 0, 0, 0, 0, 0},
		uniform_locations:   make(map[string]int32),
		attribute_locations: make(map[string]uint32),
		validated:           false,
	}
	program := gl.CreateProgram()
	if program == 0 {
		return shader, &OpenGLError{errorString: "Couldn't create program"}
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

func toGohomeShaderType(shader_type uint32) uint8 {
	switch shader_type {
	case gl.VERTEX_SHADER:
		return gohome.VERTEX
	case gl.FRAGMENT_SHADER:
		return gohome.FRAGMENT
	}

	return 255
}

func toShaderTypeName(shader_type uint32) string {
	return getShaderTypeName(toGohomeShaderType(shader_type))
}

func getAttributeSizeForType(atype string) uint32 {
	switch atype {
	case "mat3":
		return 3
	case "mat4":
		return 4
	}

	return 1
}

func (s *OpenGLES2Shader) getAttributeNames(program uint32, src string) []string {
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

	s.attribute_sizes = make(map[string]uint32)

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

func (s *OpenGLES2Shader) bindAttributesFromFile(program uint32, src string) {
	attributeNames := s.getAttributeNames(program, src)

	var index uint32 = 0
	for i := 0; i < len(attributeNames); i++ {
		gl.BindAttribLocation(program, index, attributeNames[i]+"\x00")
		index += s.attribute_sizes[attributeNames[i]]
	}
}

func (s *OpenGLES2Shader) compileOpenGLES2Shader(shader_name string, shader_type uint32, src string, program uint32) (uint32, error) {
	shader := gl.CreateShader(shader_type)
	var srcs [1]string
	srcs[0] = src
	var lengts [1]int32
	lengts[0] = int32(len(src))
	gl.ShaderSource(shader, 1, srcs[:], lengts[:])
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		buf := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, buf)
		logText := string(buf)

		return 0, &OpenGLError{errorString: logText}
	}
	gl.GetError()
	gl.AttachShader(program, shader)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return 0, &OpenGLError{errorString: "Couldn't attach " + toShaderTypeName(shader_type) + " of " + shader_name + ": ErrorCode: " + strconv.Itoa(int(err))}
	}
	if shader_type == gl.VERTEX_SHADER {
		s.bindAttributesFromFile(program, src)
	}

	return shader, nil
}

func (s *OpenGLES2Shader) AddShader(shader_type uint8, src string) error {
	if shader_type == gohome.GEOMETRY {
		render, _ := gohome.Render.(*OpenGLES2Renderer)
		if !render.HasFunctionAvailable("GEOMETRY_SHADER") {
			return &OpenGLError{errorString: "Geometry shaders are not supported by this implementation"}
		}
	}

	var err error
	var shaderName uint32
	switch shader_type {
	case gohome.VERTEX:
		shaderName, err = s.compileOpenGLES2Shader(s.name, gl.VERTEX_SHADER, src, s.program)
	case gohome.FRAGMENT:
		shaderName, err = s.compileOpenGLES2Shader(s.name, gl.FRAGMENT_SHADER, src, s.program)
	}

	if err != nil {
		return &OpenGLError{errorString: "Couldn't compile " + getShaderTypeName(shader_type) + ": " + err.Error()}
	}

	s.shaders[shader_type] = shaderName

	return nil
}

func (s *OpenGLES2Shader) deleteAllShaders() {
	for i := 0; i < 6; i++ {
		if s.shaders[i] != 0 {
			gl.DetachShader(s.program, s.shaders[i])
			gl.DeleteShader(s.shaders[i])
		}
	}
}

func (s *OpenGLES2Shader) Link() error {
	defer s.deleteAllShaders()

	gl.GetError()
	gl.LinkProgram(s.program)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return &OpenGLError{errorString: "Couldn't link: ErrorCode: " + strconv.Itoa(int(err))}
	}

	var status int32
	gl.GetError()
	gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
	if err := gl.GetError(); err != gl.NO_ERROR {
		return &OpenGLError{errorString: "Couldn't link: Couldn't get link status: ErrorCode: " + strconv.Itoa(int(err))}
	}
	if status == gl.FALSE {
		var logLength int32
		gl.GetError()
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)
		if err := gl.GetError(); err != gl.NO_ERROR {
			return &OpenGLError{errorString: "Couldn't link: Couldn't get info log length: ErrorCode: " + strconv.Itoa(int(err))}
		}

		buf := make([]byte, logLength)
		gl.GetError()
		gl.GetProgramInfoLog(s.program, logLength, nil, buf)
		logtext := string(buf)
		if err := gl.GetError(); err != gl.NO_ERROR {
			return &OpenGLError{errorString: "Couldn't link: Couldn't get info log: ErrorCode: " + strconv.Itoa(int(err))}
		}

		return &OpenGLError{errorString: "Couldn't link: " + logtext}
	}

	return nil
}

func (s *OpenGLES2Shader) Use() {
	gl.GetError()
	gl.UseProgram(s.program)
	handleOpenGLError("Shader", s.name, "glUseProgram")
}

func (s *OpenGLES2Shader) Unuse() {
	gl.GetError()
	gl.UseProgram(0)
	handleOpenGLError("Shader", s.name, "glUseProgram with 0")
}

func (s *OpenGLES2Shader) Setup() error {
	return s.validate()
}

func (s *OpenGLES2Shader) Terminate() {
	gl.DeleteProgram(s.program)
}

func (s *OpenGLES2Shader) getUniformLocation(name string) int32 {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		gl.GetError()
		loc = gl.GetUniformLocation(s.program, name+"\x00")
		handleOpenGLError("Shader", s.name, "glGetUniformLocation")
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		gohome.ErrorMgr.Message(gohome.ERROR_LEVEL_WARNING, "Shader", s.name, "Couldn't find uniform "+name)
	}
	return loc
}

func (s *OpenGLES2Shader) SetUniformV2(name string, value mgl32.Vec2) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform2f(loc, value[0], value[1])
		handleOpenGLError("Shader", s.name, "glUniform2f")
	}
}
func (s *OpenGLES2Shader) SetUniformV3(name string, value mgl32.Vec3) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform3f(loc, value[0], value[1], value[2])
		handleOpenGLError("Shader", s.name, "glUniform3f")
	}
}
func (s *OpenGLES2Shader) SetUniformV4(name string, value mgl32.Vec4) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform4f(loc, value[0], value[1], value[2], value[3])
		handleOpenGLError("Shader", s.name, "glUniform4f")
	}
}
func (s *OpenGLES2Shader) SetUniformIV2(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform2i(loc, value[0], value[1])
		handleOpenGLError("Shader", s.name, "glUniform2i")
	}
}
func (s *OpenGLES2Shader) SetUniformIV3(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform3i(loc, value[0], value[1], value[2])
		handleOpenGLError("Shader", s.name, "glUniform3i")
	}
}
func (s *OpenGLES2Shader) SetUniformIV4(name string, value []int32) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform4i(loc, value[0], value[1], value[2], value[3])
		handleOpenGLError("Shader", s.name, "glUniform4i")
	}
}
func (s *OpenGLES2Shader) SetUniformF(name string, value float32) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform1f(loc, value)
		handleOpenGLError("Shader", s.name, "glUniform1f")
	}
}
func (s *OpenGLES2Shader) SetUniformI(name string, value int32) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.Uniform1i(loc, value)
		handleOpenGLError("Shader", s.name, "glUniform1i")
	}
}
func (s *OpenGLES2Shader) SetUniformUI(name string, value uint32) {
	gohome.ErrorMgr.Error("Shader", s.name, "SetUniformUI does not work in OpenGLES2 2.0")
}

func (s *OpenGLES2Shader) SetUniformB(name string, value uint8) {
	s.SetUniformI(name, int32(value))
}
func (s *OpenGLES2Shader) SetUniformM2(name string, value mgl32.Mat2) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.UniformMatrix2fv(loc, 1, gl.FALSE, &value[0])
		handleOpenGLError("Shader", s.name, "glUniformMatrix2fv")
	}
}
func (s *OpenGLES2Shader) SetUniformM3(name string, value mgl32.Mat3) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.UniformMatrix3fv(loc, 1, gl.FALSE, &value[0])
		handleOpenGLError("Shader", s.name, "glUniformMatrix3fv")
	}
}
func (s *OpenGLES2Shader) SetUniformM4(name string, value mgl32.Mat4) {
	loc := s.getUniformLocation(name)
	if loc != -1 {
		gl.GetError()
		gl.UniformMatrix4fv(loc, 1, gl.FALSE, &value[0])
		handleOpenGLError("Shader", s.name, "glUniformMatrix4fv")
	}
}

func (s *OpenGLES2Shader) SetUniformMaterial(mat gohome.Material) {
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

func (s *OpenGLES2Shader) SetUniformLights(lightCollectionIndex int32) {
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

func (s *OpenGLES2Shader) GetName() string {
	return s.name
}

func (s *OpenGLES2Shader) validate() error {
	if s.validated {
		return nil
	}
	s.Use()
	maxtextures := gohome.Render.GetMaxTextures()
	for i := 0; i < 31; i++ {
		s.SetUniformI(gohome.POINT_LIGHTS_UNIFORM_NAME+gohome.SHADOWMAP_UNIFORM_NAME+"["+strconv.Itoa(i)+"]", maxtextures-1)
	}
	s.Unuse()
	s.validated = true
	gl.ValidateProgram(s.program)
	var status int32
	gl.GetProgramiv(s.program, gl.VALIDATE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		buf := make([]byte, logLength)
		gl.GetProgramInfoLog(s.program, logLength, nil, buf)
		logtext := string(buf)
		s.validated = false
		return &OpenGLError{"Couldn't validate: " + logtext}
	}

	return nil
}

func (s *OpenGLES2Shader) AddAttribute(name string, location uint32) {
	gl.BindAttribLocation(s.program, location, name)
}
