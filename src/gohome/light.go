package gohome

import (
	"image/color"
	"math"
	"strconv"

	"github.com/PucklaJ/mathgl/mgl32"
)

const (
	AMBIENT_LIGHT_UNIFORM_NAME          = "ambientLight"
	POINT_LIGHTS_UNIFORM_NAME           = "pointLights"
	DIRECTIONAL_LIGHTS_UNIFORM_NAME     = "directionalLights"
	SPOT_LIGHTS_UNIFORM_NAME            = "spotLights"
	POSITION_UNIFORM_NAME               = "position"
	DIRECTION_UNIFORM_NAME              = "direction"
	DIFFUSE_COLOR_UNIFORM_NAME          = "diffuseColor"
	SPECULAR_COLOR_UNIFORM_NAME         = "specularColor"
	INNERCUTOFF_UNIFORM_NAME            = "innerCutOff"
	OUTERCUTOFF_UNIFORM_NAME            = "outerCutOff"
	ATTENTUATION_UNIFORM_NAME           = "attentuation"
	ATTENTUATION_CONSTANT_UNIFORM_NAME  = "constant"
	ATTENTUATION_LINEAR_UNIFORM_NAME    = "linear"
	ATTENTUATION_QUADRATIC_UNIFORM_NAME = "quadratic"
	NUM_POINT_LIGHTS_UNIFORM_NAME       = "numPointLights"
	NUM_DIRECTIONAL_LIGHTS_UNIFORM_NAME = "numDirectionalLights"
	NUM_SPOT_LIGHTS_UNIFORM_NAME        = "numSpotLights"
	LIGHT_SPACE_MATRIX_UNIFORM_NAME     = "lightSpaceMatrix"
	SHADOWMAP_UNIFORM_NAME              = "shadowmap"
	CASTSSHADOWS_UNIFORM_NAME           = "castsShadows"
	SHADOW_DISTANCE_UNIFORM_NAME        = "shadowDistance"
	FAR_PLANE_UNIFORM_NAME              = "farPlane"
	SHADOWMAP_SIZE_UNIFORM_NAME         = "shadowMapSize"

	SHADOWMAP_SHADER_NAME           = "ShadowMap"
	SHADOWMAP_INSTANCED_SHADER_NAME = "ShadowMap Instanced"

	DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE = 1024 * 4
	DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE        = 1024
)

// A strcut holding the values for attenuation of lights
type Attentuation struct {
	// The constant factor
	Constant float32
	// The linear factor
	Linear float32
	// The quadratic factor
	Quadratic float32
}

// Sets the uniform values of s
func (a Attentuation) SetUniforms(s Shader, variableName string, arrayIndex int) {
	s.SetUniformF(variableName+"["+strconv.Itoa(arrayIndex)+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_CONSTANT_UNIFORM_NAME, a.Constant)
	s.SetUniformF(variableName+"["+strconv.Itoa(arrayIndex)+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_LINEAR_UNIFORM_NAME, a.Linear)
	s.SetUniformF(variableName+"["+strconv.Itoa(arrayIndex)+"]."+ATTENTUATION_UNIFORM_NAME+"."+ATTENTUATION_QUADRATIC_UNIFORM_NAME, a.Quadratic)
}

// A Light with a position emitting in all directions
type PointLight struct {
	// The position of the light in world coordinates
	Position mgl32.Vec3

	// The diffuse color of the light
	DiffuseColor color.Color
	// The specular color of the light
	SpecularColor color.Color

	// The attenuation values of the light
	Attentuation
}

// Sets the uniform values of s
func (pl PointLight) SetUniforms(s Shader, arrayIndex int) {
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position)
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
	pl.Attentuation.SetUniforms(s, POINT_LIGHTS_UNIFORM_NAME, arrayIndex)
}

// A light with only a direction and no position
type DirectionalLight struct {
	// The direction in which the light is shining
	Direction mgl32.Vec3

	// The diffuse color of the light
	DiffuseColor color.Color
	// The specular color of the light
	SpecularColor color.Color

	// The shadow map texture of the light
	ShadowMap RenderTexture
	// Wether the light casts shadows
	CastsShadows uint8
	// A view matrix using the direction as the look direction
	LightSpaceMatrix mgl32.Mat4

	ShadowDistance float32

	lightCam Camera3D
}

// Sets the uniforms of s
func (pl *DirectionalLight) SetUniforms(s Shader, arrayIndex int) {
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+DIRECTION_UNIFORM_NAME, pl.Direction)
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
	s.SetUniformM4(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix)
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

// Initialises the shadow map of this light
func (this *DirectionalLight) InitShadowmap(width, height int) {
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
	if cam != nil {
		cam.CalculateViewMatrix()
		inverseViewMatrix = cam.GetInverseViewMatrix()
	} else {
		inverseViewMatrix = mgl32.Ident4()
	}
	lightCam.CalculateViewMatrix()
	lightViewMatrix = lightCam.GetViewMatrix()

	pointsViewSpace = proj.GetFrustum()

	for i := 0; i < 8; i++ {
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

// Renders all objects that cast shadows on the shadow map
func (this *DirectionalLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE, DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE)
	}
	if this.ShadowMap == nil {
		this.CastsShadows = 0
		return
	}
	if ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME) == nil {
		this.ShadowMap.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 255})
		this.ShadowMap.UnsetAsTarget()
		return
	}

	var prevCamera *Camera3D
	if len(RenderMgr.camera3Ds) != 0 {
		prevCamera = RenderMgr.camera3Ds[0]
	}
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
	RenderMgr.Render(TYPE_3D_NORMAL|TYPE_CASTS_SHADOWS, 0, -1, -1)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_INSTANCED|TYPE_CASTS_SHADOWS, 0, -1, -1)

	Render.SetBacckFaceCulling(true)
	this.ShadowMap.UnsetAsTarget()
	RenderMgr.SetCamera3D(prevCamera, 0)

	RenderMgr.ForceShader3D = nil

	RenderMgr.SetProjection3D(prevProjection)

	this.LightSpaceMatrix = projection.GetProjectionMatrix().Mul4(this.lightCam.GetViewMatrix())
}

// A light with a position and a direction
type SpotLight struct {
	// The position of the light
	Position mgl32.Vec3
	// The direction of the light
	Direction mgl32.Vec3

	// The diffuse color of the light
	DiffuseColor color.Color
	// The specular color of the light
	SpecularColor color.Color

	// The angle at which the light starts to fade away in degrees
	InnerCutOff float32
	// The angle at which the light is completely faded away in degrees
	OuterCutOff float32

	// The attenuation values of this light
	Attentuation

	// The shadow map of this light
	ShadowMap RenderTexture
	// Wether this light should cast shadows
	CastsShadows uint8
	// A view matrix that uses the position and direction of the light
	LightSpaceMatrix mgl32.Mat4
	// The near plane used for the rendering of the shadow map
	NearPlane float32
	// The far plane used for the rendering of the shadow map
	FarPlane float32
}

// Sets the uniforms of s
func (pl *SpotLight) SetUniforms(s Shader, arrayIndex int) {
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+POSITION_UNIFORM_NAME, pl.Position)
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+DIRECTION_UNIFORM_NAME, pl.Direction)
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor))
	s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(arrayIndex)+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor))
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
		LoadGeneratedShaderShadowMap(0)
	}
	if ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME) == nil {
		LoadGeneratedShaderShadowMap(SHADER_FLAG_INSTANCED)
	}
}

// Initialises the shadow map of this light
func (this *SpotLight) InitShadowmap(width, height int) {
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

// Renders all objects that cast shadows onto this shadow map
func (this *SpotLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE, DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE)
	}
	if this.ShadowMap == nil {
		this.CastsShadows = 0
		return
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
		NearPlane: this.NearPlane,
		FarPlane:  this.FarPlane,
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
	RenderMgr.Render(TYPE_3D_NORMAL|TYPE_CASTS_SHADOWS, 6, -1, -1)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_INSTANCED|TYPE_CASTS_SHADOWS, 6, -1, -1)

	Render.SetBacckFaceCulling(true)
	this.ShadowMap.UnsetAsTarget()

	RenderMgr.ForceShader3D = nil

	RenderMgr.SetProjection3D(prevProjection)

	this.LightSpaceMatrix = projection.GetProjectionMatrix().Mul4(camera.GetViewMatrix())
}

// A collection of point lights, directional lights and spot lights
type LightCollection struct {
	// The color of the ambient light
	AmbientLight color.Color

	// All point lights of this collection
	PointLights []*PointLight
	// All directional lights of this collection
	DirectionalLights []*DirectionalLight
	// All spot lights of this collection
	SpotLights []*SpotLight
}

// Adds a point light to this collection
func (this *LightCollection) AddPointLight(pl *PointLight) {
	this.PointLights = append(this.PointLights, pl)
}

// Adds a directional light to this collection
func (this *LightCollection) AddDirectionalLight(pl *DirectionalLight) {
	this.DirectionalLights = append(this.DirectionalLights, pl)
}

// Adds a spot light to this collection
func (this *LightCollection) AddSpotLight(pl *SpotLight) {
	this.SpotLights = append(this.SpotLights, pl)
}

// Renders all shadow maps
func (this *LightCollection) RenderShadowMaps() {
	for i := 0; i < len(this.DirectionalLights); i++ {
		this.DirectionalLights[i].RenderShadowMap()
	}
	for i := 0; i < len(this.SpotLights); i++ {
		this.SpotLights[i].RenderShadowMap()
	}
}

// A manager holding multiple light collections
type LightManager struct {
	// All light collections
	LightCollections []LightCollection
	// The index of the currently used light collection
	CurrentLightCollection int
}

// Initialises the values of the light manager
func (this *LightManager) Init() {
	this.LightCollections = make([]LightCollection, 1)
	this.CurrentLightCollection = 0
}

// Updates all light collections / rendering all shadow maps
func (this *LightManager) Update() {
	for i := 0; i < len(this.LightCollections); i++ {
		this.LightCollections[i].RenderShadowMaps()
	}
}

// Sets the ambient light of the given collection
func (this *LightManager) SetAmbientLight(color color.Color, lightCollectionIndex int) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if len(this.LightCollections)-1 < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-len(this.LightCollections)-1)...)
	}
	this.LightCollections[lightCollectionIndex].AmbientLight = color
}

// Adds a point light to a given collection
func (this *LightManager) AddPointLight(pl *PointLight, lightCollectionIndex int) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if len(this.LightCollections)-1 < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-len(this.LightCollections)-1)...)
	}
	this.LightCollections[lightCollectionIndex].AddPointLight(pl)
}

// Adds a directional light to a given collection
func (this *LightManager) AddDirectionalLight(pl *DirectionalLight, lightCollectionIndex int) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if len(this.LightCollections)-1 < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-len(this.LightCollections)-1)...)
	}
	this.LightCollections[lightCollectionIndex].AddDirectionalLight(pl)
}

// Adds a spot light to a given collection
func (this *LightManager) AddSpotLight(pl *SpotLight, lightCollectionIndex int) {
	if len(this.LightCollections) == 0 {
		this.LightCollections = make([]LightCollection, 1)
	} else if len(this.LightCollections)-1 < lightCollectionIndex {
		this.LightCollections = append(this.LightCollections, make([]LightCollection, lightCollectionIndex-len(this.LightCollections)-1)...)
	}
	this.LightCollections[lightCollectionIndex].AddSpotLight(pl)
}

// Sets the current light collection index to -1, not using any lights (uses white as ambient light)
func (this *LightManager) DisableLighting() {
	this.CurrentLightCollection = -1
}

// Sets the current light collection index to 0
func (this *LightManager) EnableLighting() {
	this.CurrentLightCollection = 0
}

// The LightManager that should be used for everything
var LightMgr LightManager
