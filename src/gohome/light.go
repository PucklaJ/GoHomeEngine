package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image/color"
	"strconv"
)

const (
	AMBIENT_LIGHT_UNIFORM_NAME          string = "ambientLight"
	POINT_LIGHTS_UNIFORM_NAME           string = "pointLights"
	DIRECTIONAL_LIGHTS_UNIFORM_NAME     string = "directionalLights"
	SPOT_LIGHTS_UNIFORM_NAME            string = "spotLights"
	POSITION_UNIFORM_NAME               string = "position"
	DIRECTION_UNIFORM_NAME              string = "direction"
	DIFFUSE_COLOR_UNIFORM_NAME          string = "diffuseColor"
	SPECULAR_COLOR_UNIFORM_NAME         string = "specularColor"
	INNERCUTOFF_UNIFORM_NAME            string = "innerCutOff"
	OUTERCUTOFF_UNIFORM_NAME            string = "outerCutOff"
	ATTENTUATION_UNIFORM_NAME           string = "attentuation"
	ATTENTUATION_CONSTANT_UNIFORM_NAME  string = "constant"
	ATTENTUATION_LINEAR_UNIFORM_NAME    string = "linear"
	ATTENTUATION_QUADRATIC_UNIFORM_NAME string = "quadratic"
	NUM_POINT_LIGHTS_UNIFORM_NAME       string = "numPointLights"
	NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME string = "numDirectionalLights"
	NUM_SPOT_LIGHTS_UNIFORM_NAME        string = "numSpotLights"
	LIGHT_SPACE_MATRIX_UNIFORM_NAME     string = "lightSpaceMatrix"
	SHADOWMAP_UNIFORM_NAME              string = "shadowmap"
	CASTSSHADOWS_UNIFORM_NAME           string = "castsShadows"

	SHADOWMAP_SHADER_NAME string = "ShadowMap"
)

type Attentuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

func (a Attentuation) SetUniforms(s Shader, variableName string, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_CONSTANT_UNIFORM_NAME, a.Constant); err != nil {
		return err
	}
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_LINEAR_UNIFORM_NAME, a.Linear); err != nil {
		return err
	}
	if err = s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_QUADRATIC_UNIFORM_NAME, a.Quadratic); err != nil {
		return err
	}

	return nil
}

type PointLight struct {
	Position mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	Attentuation
}

func (pl PointLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = pl.Attentuation.SetUniforms(s, POINT_LIGHTS_UNIFORM_NAME, arrayIndex); err != nil {
		return err
	}
	return nil
}

type DirectionalLight struct {
	Direction mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	ShadowMap        RenderTexture
	CastsShadows     uint8
	LightSpaceMatrix mgl32.Mat4

	lightCam Camera3D
}

func (pl DirectionalLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = s.SetUniformM4(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix); err != nil {
		return err
	}
	if err = s.SetUniformB(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows); err != nil {
		return err
	}
	if pl.CastsShadows == 1 {
		rnd, ok := Render.(*OpenGLRenderer)
		if ok {
			s.SetUniformI(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, int32(rnd.CurrentTextureUnit))
			pl.ShadowMap.Bind(rnd.CurrentTextureUnit)
			rnd.CurrentTextureUnit++
		}
	}

	return nil
}

func (this *DirectionalLight) InitShadowmap(width, height uint32) {
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("DirectionallightShadowmap", width, height, 1, true, false, true)
		this.ShadowMap.SetBorderDepth(1.0)
		this.ShadowMap.SetWrapping(WRAPPING_CLAMP_TO_BORDER)
	}
}

func Mat4MulVec3(matrix mgl32.Mat4, vec mgl32.Vec3) mgl32.Vec3 {
	vec4 := vec.Vec4(1.0)
	temp := vec4

	temp[0] = matrix.At(0, 0)*vec4[0] + matrix.At(0, 1)*vec4[1] + matrix.At(0, 2)*vec4[2] + matrix.At(0, 3)*vec4[3]
	temp[1] = matrix.At(1, 0)*vec4[0] + matrix.At(1, 1)*vec4[1] + matrix.At(1, 2)*vec4[2] + matrix.At(1, 3)*vec4[3]
	temp[2] = matrix.At(2, 0)*vec4[0] + matrix.At(2, 1)*vec4[1] + matrix.At(2, 2)*vec4[2] + matrix.At(2, 3)*vec4[3]
	temp[3] = matrix.At(3, 0)*vec4[0] + matrix.At(3, 1)*vec4[1] + matrix.At(3, 2)*vec4[2] + matrix.At(3, 3)*vec4[3]

	vec4 = temp

	return vec4.Vec3()
}

func calculateDirectionalLightShadowMapProjection(cam *Camera3D, lightCam *Camera3D, proj Projection, dl *DirectionalLight) Ortho3DProjection {
	var pointsViewSpace /*pointsWorldSpace,*/, pointsLightViewSpace [8]mgl32.Vec3
	var inverseViewMatrix, lightViewMatrix mgl32.Mat4
	var projection Ortho3DProjection
	var minX, minY, minZ float32
	var maxX, maxY, maxZ float32
	var center mgl32.Vec3
	const OFFSET float32 = 10.0

	cam.CalculateViewMatrix()
	lightCam.CalculateViewMatrix()
	inverseViewMatrix = cam.GetInverseViewMatrix()
	lightViewMatrix = lightCam.GetViewMatrix()

	pointsViewSpace = proj.GetFrustum()

	for i := 0; i < 8; i++ {
		pointsLightViewSpace[i] = Mat4MulVec3(lightViewMatrix, Mat4MulVec3(inverseViewMatrix, pointsViewSpace[i]))
		if i == 0 {
			minX = pointsLightViewSpace[i][0]
			minY = pointsLightViewSpace[i][1]
			minZ = pointsLightViewSpace[i][2]

			maxX = pointsLightViewSpace[i][0]
			maxY = pointsLightViewSpace[i][1]
			maxZ = pointsLightViewSpace[i][2]
		} else {
			mgl32.SetMin(&minX, &pointsLightViewSpace[i][0])
			mgl32.SetMin(&minY, &pointsLightViewSpace[i][1])
			mgl32.SetMin(&minZ, &pointsLightViewSpace[i][2])

			mgl32.SetMax(&maxX, &pointsLightViewSpace[i][0])
			mgl32.SetMax(&maxY, &pointsLightViewSpace[i][1])
			mgl32.SetMax(&maxZ, &pointsLightViewSpace[i][2])
		}
	}
	// maxZ += OFFSET

	center[0] = (minX + maxX) / 2.0
	center[1] = (minY + maxY) / 2.0
	center[2] = (minZ + maxZ) / 2.0

	lightCam.Position = Mat4MulVec3(lightViewMatrix.Inv(), center)
	lightCam.LookDirection = dl.Direction.Add(mgl32.Vec3{1e-19, 1e-19, 1e-19})
	lightCam.CalculateViewMatrix()
	lightViewMatrix = lightCam.GetViewMatrix()

	// for i := 0; i < 8; i++ {
	// 	pointsLightViewSpace[i] = Mat4MulVec3(lightViewMatrix, pointsWorldSpace[i])
	// 	if i == 0 {
	// 		minX = pointsLightViewSpace[i][0]
	// 		minY = pointsLightViewSpace[i][1]
	// 		minZ = pointsLightViewSpace[i][2]

	// 		maxX = pointsLightViewSpace[i][0]
	// 		maxY = pointsLightViewSpace[i][1]
	// 		maxZ = pointsLightViewSpace[i][2]
	// 	} else {
	// 		mgl32.SetMin(&minX, &pointsLightViewSpace[i][0])
	// 		mgl32.SetMin(&minY, &pointsLightViewSpace[i][1])
	// 		mgl32.SetMin(&minZ, &pointsLightViewSpace[i][2])

	// 		mgl32.SetMax(&maxX, &pointsLightViewSpace[i][0])
	// 		mgl32.SetMax(&maxY, &pointsLightViewSpace[i][1])
	// 		mgl32.SetMax(&maxZ, &pointsLightViewSpace[i][2])
	// 	}
	// }
	// // maxZ += OFFSET

	projection.Left = minX
	projection.Right = maxX
	projection.Bottom = minY
	projection.Top = maxY
	projection.Near = minY
	projection.Far = maxY

	return projection
}

func (this *DirectionalLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(1024, 1024)
	}

	prevCamera := RenderMgr.camera3Ds[0]
	this.Direction = this.Direction.Normalize()
	if this.lightCam.LookDirection[0] == 0.0 && this.lightCam.LookDirection[1] == 0.0 && this.lightCam.LookDirection[2] == 0.0 {
		this.lightCam.Init()
	}
	RenderMgr.SetCamera3D(&this.lightCam, 6)

	prevProjection := RenderMgr.Projection3D

	projection := calculateDirectionalLightShadowMapProjection(prevCamera, &this.lightCam, prevProjection, this)

	RenderMgr.SetProjection3D(&projection)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME)
	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 0}, 1.0)
	RenderMgr.Render(TYPE_3D, 6, -1, -1)
	this.ShadowMap.UnsetAsTarget()

	RenderMgr.ForceShader3D = nil

	RenderMgr.SetProjection3D(prevProjection)

	this.LightSpaceMatrix = projection.GetProjectionMatrix().Mul4(this.lightCam.GetViewMatrix())
}

type SpotLight struct {
	Position  mgl32.Vec3
	Direction mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	InnerCutOff float32
	OuterCutOff float32

	Attentuation

	ShadowMap        RenderTexture
	CastsShadows     uint8
	LightSpaceMatrix mgl32.Mat4
}

func (pl *SpotLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, colorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, colorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+INNERCUTOFF_UNIFORM_NAME, pl.InnerCutOff); err != nil {
		return err
	}
	if err = s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+OUTERCUTOFF_UNIFORM_NAME, pl.OuterCutOff); err != nil {
		return err
	}
	if err = pl.Attentuation.SetUniforms(s, SPOT_LIGHTS_UNIFORM_NAME, arrayIndex); err != nil {
		return err
	}
	if err = s.SetUniformM4(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix); err != nil {
		return err
	}
	if err = s.SetUniformB(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows); err != nil {
		return err
	}
	if pl.CastsShadows == 1 {
		rnd, ok := Render.(*OpenGLRenderer)
		if ok {
			s.SetUniformI(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, int32(rnd.CurrentTextureUnit))
			pl.ShadowMap.Bind(rnd.CurrentTextureUnit)
			rnd.CurrentTextureUnit++
		}
	}
	return nil
}

func (this *SpotLight) InitShadowmap(width, height uint32) {
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("SpotlightShadowmap", width, height, 1, true, false, true)
		this.ShadowMap.SetBorderDepth(1.0)
		this.ShadowMap.SetWrapping(WRAPPING_CLAMP_TO_BORDER)
	}
}

func (this *SpotLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(1024, 1024)
	}

	prevProjection := RenderMgr.Projection3D
	prevPerspective, ok := prevProjection.(*PerspectiveProjection)
	var prevFOV float32
	if ok {
		prevFOV = prevPerspective.FOV
	} else {
		prevFOV = 70.0
	}
	nW, nH := Render.GetNativeResolution()
	projection := &PerspectiveProjection{
		Width:     float32(nW),
		Height:    float32(nH),
		FOV:       prevFOV,
		NearPlane: 0.1,
		FarPlane:  1000.0,
	}
	RenderMgr.SetProjection3D(projection)
	var camera Camera3D
	camera.Position = this.Position
	camera.LookDirection = this.Direction.Add(mgl32.Vec3{1e-19, 1e-19, 1e-19})
	RenderMgr.SetCamera3D(&camera, 6)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME)
	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 0}, 1.0)
	RenderMgr.Render(TYPE_3D, 6, -1, -1)
	this.ShadowMap.UnsetAsTarget()

	RenderMgr.ForceShader3D = nil

	RenderMgr.SetProjection3D(prevProjection)

	this.LightSpaceMatrix = projection.GetProjectionMatrix().Mul4(camera.GetViewMatrix())
}

type LightCollection struct {
	AmbientLight color.Color

	PointLights       []*PointLight
	DirectionalLights []*DirectionalLight
	SpotLights        []*SpotLight
}

func (this *LightCollection) AddPointLight(pl *PointLight) {
	this.PointLights = append(this.PointLights, pl)
}

func (this *LightCollection) AddDirectionalLight(pl *DirectionalLight) {
	this.DirectionalLights = append(this.DirectionalLights, pl)
}

func (this *LightCollection) AddSpotLight(pl *SpotLight) {
	this.SpotLights = append(this.SpotLights, pl)
}

func (this *LightCollection) RenderShadowMaps() {
	for i := 0; i < len(this.DirectionalLights); i++ {
		this.DirectionalLights[i].RenderShadowMap()
	}
	for i := 0; i < len(this.SpotLights); i++ {
		this.SpotLights[i].RenderShadowMap()
	}
}

type LightManager struct {
	lightCollections       []LightCollection
	CurrentLightCollection int32
}

func (this *LightManager) Init() {
	this.lightCollections = make([]LightCollection, 1)
	this.CurrentLightCollection = 0
	ResourceMgr.LoadShader(SHADOWMAP_SHADER_NAME, "shadowMapVert.glsl", "shadowMapFrag.glsl", "", "", "", "")
}

func (this *LightManager) Update() {
	for i := 0; i < len(this.lightCollections); i++ {
		this.lightCollections[i].RenderShadowMaps()
	}
}

func (this *LightManager) SetAmbientLight(color color.Color, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AmbientLight = color
}

func (this *LightManager) AddPointLight(pl *PointLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddPointLight(pl)
}

func (this *LightManager) AddDirectionalLight(pl *DirectionalLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddDirectionalLight(pl)
}

func (this *LightManager) AddSpotLight(pl *SpotLight, lightCollectionIndex uint32) {
	if len(this.lightCollections) == 0 {
		this.lightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.lightCollections)-1) < lightCollectionIndex {
		this.lightCollections = append(this.lightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.lightCollections)-1))...)
	}
	this.lightCollections[lightCollectionIndex].AddSpotLight(pl)
}

var LightMgr LightManager
