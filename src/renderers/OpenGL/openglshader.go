package renderer

import (
	// "fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	"strconv"
	// "log"
	"github.com/PucklaMotzer09/gohomeengine/src/gohome"
	"strings"
)

type OpenGLShader struct {
	program             uint32
	name                string
	shaders             [6]uint32
	uniform_locations   map[string]int32
	attribute_locations map[string]uint32
	validated           bool
}

func CreateOpenGLShader(name string) (*OpenGLShader, error) {
	shader := &OpenGLShader{
		program:             0,
		name:                name,
		shaders:             [6]uint32{0, 0, 0, 0, 0, 0},
		uniform_locations:   make(map[string]int32),
		attribute_locations: make(map[string]uint32),
		validated:           false,
	}
	program := gl.CreateProgram()
	if program == 0 {
		return shader, &OpenGLError{errorString: "Couldn't create shader program of " + name}
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

func compileOpenGLShader(shader_type uint32, src **uint8, program uint32) (uint32, error) {
	shader := gl.CreateShader(shader_type)
	gl.ShaderSource(shader, 1, src, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logText := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logText))

		return 0, &OpenGLError{errorString: logText}
	}
	gl.AttachShader(program, shader)

	return shader, nil
}

func (s *OpenGLShader) AddShader(shader_type uint8, src string) error {
	csource, free := gl.Strs(src + "\x00")
	defer free()
	var err error
	var shaderName uint32
	switch shader_type {
	case gohome.VERTEX:
		shaderName, err = compileOpenGLShader(gl.VERTEX_SHADER, csource, s.program)
	case gohome.FRAGMENT:
		shaderName, err = compileOpenGLShader(gl.FRAGMENT_SHADER, csource, s.program)
	case gohome.GEOMETRY:
		shaderName, err = compileOpenGLShader(gl.GEOMETRY_SHADER, csource, s.program)
	case gohome.TESSELLETION:
		shaderName, err = compileOpenGLShader(gl.TESS_CONTROL_SHADER, csource, s.program)
	case gohome.EVELUATION:
		shaderName, err = compileOpenGLShader(gl.TESS_EVALUATION_SHADER, csource, s.program)
	case gohome.COMPUTE:
		shaderName, err = compileOpenGLShader(gl.COMPUTE_SHADER, csource, s.program)
	}

	if err != nil {
		return &OpenGLError{errorString: "Couldn't compile source of " + getShaderTypeName(shader_type) + " of " + s.name + ": " + err.Error()}
	}

	s.shaders[shader_type] = shaderName

	return nil
}

func (s *OpenGLShader) deleteAllShaders() {
	for i := 0; i < 6; i++ {
		if s.shaders[i] != 0 {
			gl.DetachShader(s.program, s.shaders[i])
			gl.DeleteShader(s.shaders[i])
		}
	}
}

func (s *OpenGLShader) Link() error {
	defer s.deleteAllShaders()

	gl.LinkProgram(s.program)

	var status int32
	gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		logtext := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(logtext))

		return &OpenGLError{errorString: "Couldn't link shader " + s.name + ": " + logtext}
	}

	return nil
}

func (s *OpenGLShader) Use() {
	gl.UseProgram(s.program)
}

func (s *OpenGLShader) Unuse() {
	gl.UseProgram(0)
}

func (s *OpenGLShader) Setup() error {
	s.validate()
	if runtime.GOOS != "windows" {
		s.Use()

		var c int32
		var i uint32
		s.uniform_locations = make(map[string]int32)
		s.attribute_locations = make(map[string]uint32)
		gl.GetProgramiv(s.program, gl.ACTIVE_UNIFORMS, &c)
		for i = 0; i < uint32(c); i++ {
			var buf [256]byte
			gl.GetActiveUniform(s.program, i, 256, nil, nil, nil, &buf[0])
			loc := gl.GetUniformLocation(s.program, &buf[0])
			name := gl.GoStr(&buf[0])
			s.uniform_locations[name] = loc
		}
		gl.GetProgramiv(s.program, gl.ACTIVE_ATTRIBUTES, &c)
		for i = 0; i < uint32(c); i++ {
			var buf [256]byte
			gl.GetActiveAttrib(s.program, i, 256, nil, nil, nil, &buf[0])
			loc := gl.GetAttribLocation(s.program, &buf[0])
			name := gl.GoStr(&buf[0])
			s.attribute_locations[name] = uint32(loc)
		}

		s.Unuse()
	}

	return nil
}

func (s *OpenGLShader) Terminate() {
	gl.DeleteProgram(s.program)
}

func (s *OpenGLShader) SetUniformV2(name string, value mgl32.Vec2) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}
	gl.Uniform2f(loc, value[0], value[1])

	return nil
}
func (s *OpenGLShader) SetUniformV3(name string, value mgl32.Vec3) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform3f(loc, value[0], value[1], value[2])

	return nil
}
func (s *OpenGLShader) SetUniformV4(name string, value mgl32.Vec4) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform4f(loc, value[0], value[1], value[2], value[3])

	return nil
}
func (s *OpenGLShader) SetUniformIV2(name string, value []int32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform2i(loc, value[0], value[1])

	return nil
}
func (s *OpenGLShader) SetUniformIV3(name string, value []int32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform3i(loc, value[0], value[1], value[2])

	return nil
}
func (s *OpenGLShader) SetUniformIV4(name string, value []int32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform4i(loc, value[0], value[1], value[2], value[3])

	return nil
}
func (s *OpenGLShader) SetUniformF(name string, value float32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1f(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformI(name string, value int32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1i(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformUI(name string, value uint32) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1ui(loc, value)

	return nil
}
func (s *OpenGLShader) SetUniformB(name string, value uint8) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.Uniform1i(loc, int32(value))

	return nil
}
func (s *OpenGLShader) SetUniformM2(name string, value mgl32.Mat2) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix2fv(loc, 1, false, &value[0])

	return nil
}
func (s *OpenGLShader) SetUniformM3(name string, value mgl32.Mat3) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix3fv(loc, 1, false, &value[0])

	return nil
}
func (s *OpenGLShader) SetUniformM4(name string, value mgl32.Mat4) error {
	var loc int32
	var ok bool
	if loc, ok = s.uniform_locations[name]; !ok {
		loc = gl.GetUniformLocation(s.program, gl.Str(name+"\x00"))
		s.uniform_locations[name] = loc
	}
	if loc == -1 {
		return &OpenGLError{errorString: "Couldn't find uniform " + name + " in shader " + s.name}
	}

	gl.UniformMatrix4fv(loc, 1, false, &value[0])

	return nil
}

func (s *OpenGLShader) SetUniformMaterial(mat gohome.Material) error {
	var err error
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
			// fmt.Println("Binding Diffuse Texture to ", diffBind)
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
			// fmt.Println("Binding SpecularTexture to ", specBind)
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
			// fmt.Println("Binding NormalMap to ", normBind)
			mat.NormalMapLoaded = 1
			boundTextures++
		}
	} else {
		mat.NormalMapLoaded = 0
	}

	gohome.Render.DecrementTextureUnit(boundTextures)

	if err = s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.DiffuseColor)); err != nil {
		// return err
	}
	if err = s.SetUniformV3(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_COLOR_UNIFORM_NAME, gohome.ColorToVec3(mat.SpecularColor)); err != nil {
		// return err
	}
	if err = s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_TEXTURE_UNIFORM_NAME, diffBind); err != nil {
		// return err
	}
	if err = s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_TEXTURE_UNIFORM_NAME, specBind); err != nil {
		// return err
	}
	if err = s.SetUniformI(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_NORMALMAP_UNIFORM_NAME, normBind); err != nil {
		// return err
	}
	if err = s.SetUniformF(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SHINYNESS_UNIFORM_NAME, mat.Shinyness); err != nil {
		// return err
	}

	if err = s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_DIFFUSE_TEXTURE_LOADED_UNIFORM_NAME, mat.DiffuseTextureLoaded); err != nil {
		// return err
	}
	if err = s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_SPECULAR_TEXTURE_LOADED_UNIFORM_NAME, mat.SpecularTextureLoaded); err != nil {
		// return err
	}
	if err = s.SetUniformB(gohome.MATERIAL_UNIFORM_NAME+"."+gohome.MATERIAL_NORMALMAP_LOADED_UNIFORM_NAME, mat.NormalMapLoaded); err != nil {
		// return err
	}

	return err
}

func (s *OpenGLShader) SetUniformLights(lightCollectionIndex int32) error {
	if lightCollectionIndex == -1 || lightCollectionIndex > int32(len(gohome.LightMgr.LightCollections)-1) {
		var err error
		if err = s.SetUniformUI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}
		if err = s.SetUniformUI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}
		if err = s.SetUniformUI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, 0); err != nil {
			// return err
		}

		if err = s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, mgl32.Vec3{1.0, 1.0, 1.0}); err != nil {
			// return err
		}
		return nil
	}

	lightColl := gohome.LightMgr.LightCollections[lightCollectionIndex]

	var err error
	if err = s.SetUniformUI(gohome.NUM_POINT_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.PointLights))); err != nil {
		// return err
	}
	if err = s.SetUniformUI(gohome.NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.DirectionalLights))); err != nil {
		// return err
	}
	if err = s.SetUniformUI(gohome.NUM_SPOT_LIGHTS_UNIFORM_NAME, uint32(len(lightColl.SpotLights))); err != nil {
		// return err
	}

	if err = s.SetUniformV3(gohome.AMBIENT_LIGHT_UNIFORM_NAME, gohome.ColorToVec3(lightColl.AmbientLight)); err != nil {
		// return err
	}

	var i uint32
	for i = 0; i < uint32(len(lightColl.PointLights)); i++ {
		if err = lightColl.PointLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}
	for i = 0; i < uint32(len(lightColl.DirectionalLights)); i++ {
		if err = lightColl.DirectionalLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}
	for i = 0; i < uint32(len(lightColl.SpotLights)); i++ {
		if err = lightColl.SpotLights[i].SetUniforms(s, i); err != nil {
			// return err
		}
	}

	return err
}

func (s *OpenGLShader) GetName() string {
	return s.name
}

func (s *OpenGLShader) validate() error {
	if s.validated {
		return nil
	}
	s.Use()
	maxtextures := gohome.Render.GetMaxTextures()
	for i := 0; i < 31; i++ {
		s.SetUniformI("pointLights["+strconv.Itoa(i)+"].shadowmap", maxtextures-1)
	}
	s.Unuse()
	s.validated = true
	var vao uint32
	gl.CreateVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	defer gl.DeleteVertexArrays(1, &vao)
	defer gl.BindVertexArray(0)
	gl.ValidateProgram(s.program)
	var status int32
	gl.GetProgramiv(s.program, gl.VALIDATE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

		logtext := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(logtext))
		s.validated = false
		return &OpenGLError{errorString: "Couldn't validate shader " + s.name + ": " + logtext}
	}

	return nil
}

func (s *OpenGLShader) AddAttribute(name string, location uint32) {
	gl.BindAttribLocation(s.program, location, gl.Str(name))
}
