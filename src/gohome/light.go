package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/color"
	"math"
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
	SHADOW_DISTANCE_UNIFORM_NAME        string = "shadowDistance"
	FAR_PLANE_UNIFORM_NAME              string = "farPlane"
	SHADOWMAP_SIZE_UNIFORM_NAME         string = "shadowMapSize"

	SHADOWMAP_SHADER_NAME           string = "ShadowMap"
	SHADOWMAP_INSTANCED_SHADER_NAME string = "ShadowMapInstanced"

	DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE uint32 = 1024 * 4
	DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE        uint32 = 1024
)

type Attentuation struct {
	Constant  float32
	Linear    float32
	Quadratic float32
}

func (a Attentuation) SetUniforms(s Shader, variableName string, arrayIndex uint32) {
	s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_CONSTANT_UNIFORM_NAME, a.Constant)
	s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_LINEAR_UNIFORM_NAME, a.Linear)
	s.SetUniformF(variableName+"["+strconv.Itoa(int(arrayIndex))+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_QUADRATIC_UNIFORM_NAME, a.Quadratic)
}

type PointLight struct {
	Position mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	Attentuation
}

func (pl PointLight) SetUniforms(s Shader, arrayIndex uint32) {
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position)
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
	pl.Attentuation.SetUniforms(s, POINT_LIGHTS_UNIFORM_NAME, arrayIndex)
}

type DirectionalLight struct {
	Direction mgl32.Vec3

	DiffuseColor  color.Color
	SpecularColor color.Color

	ShadowMap        RenderTexture
	CastsShadows     uint8
	LightSpaceMatrix mgl32.Mat4

	ShadowDistance float32

	lightCam Camera3D
}

func (pl *DirectionalLight) SetUniforms(s Shader, arrayIndex uint32) {
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction)
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
	s.SetUniformM4(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix)
	if pl.ShadowMap != nil {
		size := make([]int32, 2)
		size[0] = int32(pl.ShadowMap.GetWidth())
		size[1] = int32(pl.ShadowMap.GetHeight())
		s.SetUniformIV2(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_SIZE_UNIFORM_NAME, size)
	}
	s.SetUniformB(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows)
	if pl.CastsShadows == 1 {
		maxtextures := Render.GetMaxTextures()
		currentTextureUnit := Render.NextTextureUnit()
		if currentTextureUnit >= uint32(maxtextures)-1 {
			s.SetUniformI(DIRECTIONAL_LIGHTS_UNIFORM_NAME+SHADOWMAP_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]", 0)
			s.SetUniformB(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, 0)
		} else {
			pl.ShadowMap.Bind(currentTextureUnit)
			s.SetUniformI(DIRECTIONAL_LIGHTS_UNIFORM_NAME+SHADOWMAP_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]", int32(currentTextureUnit))
		}

		s.SetUniformF(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOW_DISTANCE_UNIFORM_NAME, pl.ShadowDistance)
	}
}

func (this *DirectionalLight) InitShadowmap(width, height uint32) {
	if this.CastsShadows == 0 {
		return
	}
	loadShadowMapShader()
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("DirectionallightShadowmap", width, height, 1, true, false, true, false)
		this.ShadowMap.SetBorderDepth(1.0)
		this.ShadowMap.SetWrapping(WRAPPING_CLAMP_TO_BORDER)
		this.ShadowMap.SetFiltering(FILTERING_LINEAR)
	}
}

func calculateDirectionalLightShadowMapProjection(cam *Camera3D, lightCam *Camera3D, proj Projection, dl *DirectionalLight) Ortho3DProjection {
	var pointsViewSpace, pointsLightViewSpace [8]mgl32.Vec3
	var inverseViewMatrix, lightViewMatrix mgl32.Mat4
	var projection Ortho3DProjection
	var minX, minY, minZ float32
	var maxX, maxY, maxZ float32
	var center mgl32.Vec3
	var farPlaneHalfHeightShadowDistance, farPlaneHalfWidthShadowDistance float32
	var shadowDistanceVector mgl32.Vec3
	var persProj *PerspectiveProjection
	var ok bool
	const OFFSET float32 = 15.0

	if persProj, ok = proj.(*PerspectiveProjection); ok {
		farPlaneHalfWidthShadowDistance = float32(math.Tan(float64(persProj.FOV)/180.0*math.Pi) * float64(dl.ShadowDistance))
		farPlaneHalfHeightShadowDistance = farPlaneHalfWidthShadowDistance / (persProj.Width / persProj.Height)
	}
	cam.CalculateViewMatrix()
	lightCam.CalculateViewMatrix()
	inverseViewMatrix = cam.GetInverseViewMatrix()
	lightViewMatrix = lightCam.GetViewMatrix()

	pointsViewSpace = proj.GetFrustum()

	var i uint32
	for i = 0; i < 8; i++ {
		if ok {
			if i == FAR_LEFT_DOWN {
				shadowDistanceVector = mgl32.Vec3{-farPlaneHalfWidthShadowDistance, -farPlaneHalfHeightShadowDistance, -dl.ShadowDistance}
				length := shadowDistanceVector.Len()
				pointsViewSpace[i] = pointsViewSpace[i].Sub(pointsViewSpace[i].Normalize().Mul(pointsViewSpace[i].Len() - length))
			} else if i == FAR_RIGHT_DOWN {
				shadowDistanceVector = mgl32.Vec3{farPlaneHalfWidthShadowDistance, -farPlaneHalfHeightShadowDistance, -dl.ShadowDistance}
				length := shadowDistanceVector.Len()
				pointsViewSpace[i] = pointsViewSpace[i].Sub(pointsViewSpace[i].Normalize().Mul(pointsViewSpace[i].Len() - length))
			} else if i == FAR_RIGHT_UP {
				shadowDistanceVector = mgl32.Vec3{farPlaneHalfWidthShadowDistance, farPlaneHalfHeightShadowDistance, -dl.ShadowDistance}
				length := shadowDistanceVector.Len()
				pointsViewSpace[i] = pointsViewSpace[i].Sub(pointsViewSpace[i].Normalize().Mul(pointsViewSpace[i].Len() - length))
			} else if i == FAR_LEFT_UP {
				shadowDistanceVector = mgl32.Vec3{-farPlaneHalfWidthShadowDistance, farPlaneHalfHeightShadowDistance, -dl.ShadowDistance}
				length := shadowDistanceVector.Len()
				pointsViewSpace[i] = pointsViewSpace[i].Sub(pointsViewSpace[i].Normalize().Mul(pointsViewSpace[i].Len() - length))
			}

		}
		pointsLightViewSpace[i] = lightViewMatrix.Mul4(inverseViewMatrix).Mul4x1(pointsViewSpace[i].Vec4(1)).Vec3()
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
	maxZ += OFFSET

	center[0] = (minX + maxX) / 2.0
	center[1] = (minY + maxY) / 2.0
	center[2] = (minZ + maxZ) / 2.0

	lightCam.Position = lightViewMatrix.Inv().Mul4x1(center.Vec4(1)).Vec3()
	lightCam.LookDirection = dl.Direction.Add(mgl32.Vec3{1e-19, 1e-19, 1e-19})
	lightCam.CalculateViewMatrix()
	lightViewMatrix = lightCam.GetViewMatrix()

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
		this.InitShadowmap(DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE, DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE)
	}
	if ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME) == nil {
		this.ShadowMap.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 255})
		this.ShadowMap.UnsetAsTarget()
		return
	}

	prevCamera := RenderMgr.camera3Ds[0]
	this.Direction = this.Direction.Normalize()
	if this.lightCam.LookDirection[0] == 0.0 && this.lightCam.LookDirection[1] == 0.0 && this.lightCam.LookDirection[2] == 0.0 {
		this.lightCam.Init()
	}
	RenderMgr.SetCamera3D(&this.lightCam, 0)

	prevProjection := RenderMgr.Projection3D

	projection := calculateDirectionalLightShadowMapProjection(prevCamera, &this.lightCam, prevProjection, this)

	RenderMgr.SetProjection3D(&projection)

	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 255})
	Render.SetBacckFaceCulling(false)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_NORMAL, 0, -1, -1)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_INSTANCED, 0, -1, -1)

	Render.SetBacckFaceCulling(true)
	this.ShadowMap.UnsetAsTarget()
	RenderMgr.SetCamera3D(prevCamera, 0)

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

func (pl *SpotLight) SetUniforms(s Shader, arrayIndex uint32) {
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position)
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction)
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
	s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+INNERCUTOFF_UNIFORM_NAME, pl.InnerCutOff)
	s.SetUniformF(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+OUTERCUTOFF_UNIFORM_NAME, pl.OuterCutOff)
	pl.Attentuation.SetUniforms(s, SPOT_LIGHTS_UNIFORM_NAME, arrayIndex)
	s.SetUniformM4(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix)
	if pl.ShadowMap != nil {
		size := make([]int32, 2)
		size[0] = int32(pl.ShadowMap.GetWidth())
		size[1] = int32(pl.ShadowMap.GetHeight())
		s.SetUniformIV2(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_SIZE_UNIFORM_NAME, size)
	}
	s.SetUniformB(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows)
	if pl.CastsShadows == 1 {
		maxtextures := Render.GetMaxTextures()
		currentTextureUnit := Render.NextTextureUnit()
		if currentTextureUnit >= uint32(maxtextures)-1 {
			s.SetUniformI(SPOT_LIGHTS_UNIFORM_NAME+SHADOWMAP_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]", 0)
			s.SetUniformB(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, 0)
		} else {
			pl.ShadowMap.Bind(currentTextureUnit)
			s.SetUniformI(SPOT_LIGHTS_UNIFORM_NAME+SHADOWMAP_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]", int32(currentTextureUnit))
		}

	}
}

func loadShadowMapShader() {
	if ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME) == nil {
		ResourceMgr.LoadShaderSource(SHADOWMAP_SHADER_NAME, SHADOWMAP_SHADER_VERTEX_SOURCE_OPENGL, SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	}
	if ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME) == nil {
		ResourceMgr.LoadShaderSource(SHADOWMAP_INSTANCED_SHADER_NAME, SHADOWMAP_INSTANCED_SHADER_VERTEX_SOURCE_OPENGL, SHADOWMAP_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	}
}

func (this *SpotLight) InitShadowmap(width, height uint32) {
	if this.CastsShadows == 0 {
		return
	}
	loadShadowMapShader()
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("SpotlightShadowmap", width, height, 1, true, false, true, false)
		this.ShadowMap.SetBorderDepth(1.0)
		this.ShadowMap.SetWrapping(WRAPPING_CLAMP_TO_BORDER)
	}
}

func (this *SpotLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE, DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE)
	}
	if ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME) == nil {
		this.ShadowMap.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 255})
		this.ShadowMap.UnsetAsTarget()
		return
	}

	prevProjection := RenderMgr.Projection3D
	prevPerspective, ok := prevProjection.(*PerspectiveProjection)
	var prevFOV float32
	if ok {
		prevFOV = prevPerspective.FOV
	} else {
		prevFOV = 70.0
	}
	ns := Render.GetNativeResolution()
	projection := &PerspectiveProjection{
		Width:     ns[0],
		Height:    ns[1],
		FOV:       prevFOV,
		NearPlane: 0.1,
		FarPlane:  1000.0,
	}
	RenderMgr.SetProjection3D(projection)
	var camera Camera3D
	camera.Init()
	camera.Position = this.Position
	camera.LookDirection = this.Direction.Add(mgl32.Vec3{1e-19, 1e-19, 1e-19})
	RenderMgr.SetCamera3D(&camera, 6)

	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 255})
	Render.SetBacckFaceCulling(false)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_NORMAL, 6, -1, -1)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_INSTANCED, 6, -1, -1)

	Render.SetBacckFaceCulling(true)
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
	LightCollections       []LightCollection
	CurrentLightCollection int32
}

func (this *LightManager) Init() {
	this.LightCollections = make([]LightCollection, 1)
	this.CurrentLightCollection = 0
}

func (this *LightManager) Update() {
	for i := 0; i < len(this.LightCollections); i++ {
		this.LightCollections[i].RenderShadowMaps()
	}
}

func (this *LightManager) SetAmbientLight(color color.Color, lightCollectionIndex uint32) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.LightCollections)-1) < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.LightCollections)-1))...)
	}
	this.LightCollections[lightCollectionIndex].AmbientLight = color
}

func (this *LightManager) AddPointLight(pl *PointLight, lightCollectionIndex uint32) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.LightCollections)-1) < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.LightCollections)-1))...)
	}
	this.LightCollections[lightCollectionIndex].AddPointLight(pl)
}

func (this *LightManager) AddDirectionalLight(pl *DirectionalLight, lightCollectionIndex uint32) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.LightCollections)-1) < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.LightCollections)-1))...)
	}
	this.LightCollections[lightCollectionIndex].AddDirectionalLight(pl)
}

func (this *LightManager) AddSpotLight(pl *SpotLight, lightCollectionIndex uint32) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if uint32(len(this.LightCollections)-1) < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-uint32(len(this.LightCollections)-1))...)
	}
	this.LightCollections[lightCollectionIndex].AddSpotLight(pl)
}

func (this *LightManager) DisableLighting() {
	this.CurrentLightCollection = -1
}

func (this *LightManager) EnableLighting() {
	this.CurrentLightCollection = 0
}

var LightMgr LightManager
