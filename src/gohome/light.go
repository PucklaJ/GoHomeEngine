package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
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

	SHADOWMAP_SHADER_NAME                       string = "ShadowMap"
	SHADOWMAP_INSTANCED_SHADER_NAME             string = "ShadowMapInstanced"
	POINT_LIGHT_SHADOWMAP_SHADER_NAME           string = "PointlightShadowMap"
	POINT_LIGHT_SHADOWMAP_INSTANCED_SHADER_NAME string = "PointlightShadowMapInstanced"

	DEFAULT_DIRECTIONAL_LIGHTS_SHADOWMAP_SIZE uint32 = 1024 * 4
	DEFAULT_SPOT_LIGHTS_SHADOWMAP_SIZE        uint32 = 1024
	DEFAULT_POINT_LIGHTS_SHADOWMAP_SIZE       uint32 = 1024
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

	ShadowMap          RenderTexture
	CastsShadows       uint8
	LightSpaceMatrices [6]mgl32.Mat4
	FarPlane           float32
}

func (pl PointLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+POSITION_UNIFORM_NAME, pl.Position); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = pl.Attentuation.SetUniforms(s, POINT_LIGHTS_UNIFORM_NAME, arrayIndex); err != nil {
		return err
	}
	if err = s.SetUniformB(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows); err != nil {
		return err
	}
	if pl.CastsShadows == 1 {
		maxtextures := Render.GetMaxTextures()
		currentTextureUnit := Render.NextTextureUnit()
		if currentTextureUnit > uint32(maxtextures)-1 {
			s.SetUniformI(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, maxtextures-1)
			s.SetUniformB(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, 0)
		} else {
			pl.ShadowMap.Bind(currentTextureUnit)
			// fmt.Println("Binding PointLight to ", rnd.CurrentTextureUnit)
			s.SetUniformI(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, int32(currentTextureUnit))
		}
		s.SetUniformF(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+FAR_PLANE_UNIFORM_NAME, pl.FarPlane)
		for i := 0; i < 6; i++ {
			s.SetUniformM4(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME+"["+strconv.Itoa(i)+"]", pl.LightSpaceMatrices[i])
		}
	} else {
		maxtextures := Render.GetMaxTextures()
		s.SetUniformI(POINT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, maxtextures-1)
	}

	return nil
}

func (this *PointLight) InitShadowmap(width, height uint32) {
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("PointlightShadowmap", width, height, 1, false, false, true, true)
	}
}

func (this *PointLight) RenderShadowMap() {
	if this.CastsShadows == 0 {
		return
	}
	if this.ShadowMap == nil {
		this.InitShadowmap(DEFAULT_POINT_LIGHTS_SHADOWMAP_SIZE, DEFAULT_POINT_LIGHTS_SHADOWMAP_SIZE)
	}

	prevProjection := RenderMgr.Projection3D
	projection := &PerspectiveProjection{
		Width:     float32(this.ShadowMap.GetWidth()),
		Height:    float32(this.ShadowMap.GetHeight()),
		FOV:       90.0,
		NearPlane: 0.01,
		FarPlane:  this.FarPlane,
	}

	RenderMgr.Projection3D = projection

	var cameras [6]Camera3D
	for i := 0; i < 6; i++ {
		cameras[i].Init()
		cameras[i].Position = this.Position
	}

	cameras[0].LookDirection = mgl32.Vec3{1.0, 0.0, 0.0}
	cameras[1].LookDirection = mgl32.Vec3{-1.0, 0.0, 0.0}
	cameras[2].LookDirection = mgl32.Vec3{0.0, 1.0, 0.0}
	cameras[3].LookDirection = mgl32.Vec3{0.0, -1.0, 0.0}
	cameras[4].LookDirection = mgl32.Vec3{0.0, 0.0, 1.0}
	cameras[5].LookDirection = mgl32.Vec3{0.0, 0.0, -1.0}
	cameras[0].Up = mgl32.Vec3{0.0, -1.0, 0.0}
	cameras[1].Up = mgl32.Vec3{0.0, -1.0, 0.0}
	cameras[2].Up = mgl32.Vec3{0.0, 0.0, 1.0}
	cameras[3].Up = mgl32.Vec3{0.0, 0.0, -1.0}
	cameras[4].Up = mgl32.Vec3{0.0, -1.0, 0.0}
	cameras[5].Up = mgl32.Vec3{0.0, -1.0, 0.0}

	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 0}, 1.0)
	Render.SetBacckFaceCulling(false)

	shader := ResourceMgr.GetShader(POINT_LIGHT_SHADOWMAP_SHADER_NAME)
	shader.Use()
	for i := 0; i < 6; i++ {
		cameras[i].CalculateViewMatrix()
		viewMatrix := cameras[i].GetViewMatrix()
		shader.SetUniformM4("lightSpaceMatrices["+strconv.Itoa(i)+"]", viewMatrix)
	}
	shader.SetUniformV3("lightPos", this.Position)
	shader.SetUniformF(FAR_PLANE_UNIFORM_NAME, this.FarPlane)
	RenderMgr.ForceShader3D = shader
	RenderMgr.Render(TYPE_3D_NORMAL, -1, -1, -1)

	shader = ResourceMgr.GetShader(POINT_LIGHT_SHADOWMAP_INSTANCED_SHADER_NAME)
	shader.Use()
	for i := 0; i < 6; i++ {
		cameras[i].CalculateViewMatrix()
		viewMatrix := cameras[i].GetViewMatrix()
		shader.SetUniformM4("lightSpaceMatrices["+strconv.Itoa(i)+"]", viewMatrix)
	}
	shader.SetUniformV3("lightPos", this.Position)
	shader.SetUniformF(FAR_PLANE_UNIFORM_NAME, this.FarPlane)
	RenderMgr.ForceShader3D = shader
	RenderMgr.Render(TYPE_3D_INSTANCED, -1, -1, -1)

	shader.Unuse()

	Render.SetBacckFaceCulling(true)
	this.ShadowMap.UnsetAsTarget()

	RenderMgr.ForceShader3D = nil

	RenderMgr.Projection3D = prevProjection
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

func (pl *DirectionalLight) SetUniforms(s Shader, arrayIndex uint32) error {
	var err error
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIRECTION_UNIFORM_NAME, pl.Direction); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor)); err != nil {
		return err
	}
	if err = s.SetUniformM4(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+LIGHT_SPACE_MATRIX_UNIFORM_NAME, pl.LightSpaceMatrix); err != nil {
		return err
	}
	if err = s.SetUniformB(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, pl.CastsShadows); err != nil {
		return err
	}
	if pl.CastsShadows == 1 {
		maxtextures := Render.GetMaxTextures()
		currentTextureUnit := Render.NextTextureUnit()
		if currentTextureUnit >= uint32(maxtextures)-1 {
			s.SetUniformI(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, 0)
			s.SetUniformB(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, 0)
		} else {
			pl.ShadowMap.Bind(currentTextureUnit)
			// fmt.Println("Binding Directional to ", rnd.CurrentTextureUnit)
			s.SetUniformI(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, int32(currentTextureUnit))
		}

		s.SetUniformF(DIRECTIONAL_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOW_DISTANCE_UNIFORM_NAME, pl.ShadowDistance)
	}

	return nil
}

func (this *DirectionalLight) InitShadowmap(width, height uint32) {
	if this.ShadowMap != nil {
		this.ShadowMap.Terminate()
	} else {
		this.ShadowMap = Render.CreateRenderTexture("DirectionallightShadowmap", width, height, 1, true, false, true, false)
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
	maxZ += OFFSET

	center[0] = (minX + maxX) / 2.0
	center[1] = (minY + maxY) / 2.0
	center[2] = (minZ + maxZ) / 2.0

	lightCam.Position = Mat4MulVec3(lightViewMatrix.Inv(), center)
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

	prevCamera := RenderMgr.camera3Ds[0]
	this.Direction = this.Direction.Normalize()
	if this.lightCam.LookDirection[0] == 0.0 && this.lightCam.LookDirection[1] == 0.0 && this.lightCam.LookDirection[2] == 0.0 {
		this.lightCam.Init()
	}
	RenderMgr.SetCamera3D(&this.lightCam, 6)

	prevProjection := RenderMgr.Projection3D

	projection := calculateDirectionalLightShadowMapProjection(prevCamera, &this.lightCam, prevProjection, this)

	RenderMgr.SetProjection3D(&projection)

	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 0}, 1.0)
	Render.SetBacckFaceCulling(false)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_NORMAL, 6, -1, -1)

	RenderMgr.ForceShader3D = ResourceMgr.GetShader(SHADOWMAP_INSTANCED_SHADER_NAME)
	RenderMgr.Render(TYPE_3D_INSTANCED, 6, -1, -1)

	Render.SetBacckFaceCulling(true)
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
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+DIFFUSE_COLOR_UNIFORM_NAME, ColorToVec3(pl.DiffuseColor)); err != nil {
		return err
	}
	if err = s.SetUniformV3(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SPECULAR_COLOR_UNIFORM_NAME, ColorToVec3(pl.SpecularColor)); err != nil {
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
		maxtextures := Render.GetMaxTextures()
		currentTextureUnit := Render.NextTextureUnit()
		if currentTextureUnit >= uint32(maxtextures)-1 {
			s.SetUniformI(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, 0)
			s.SetUniformB(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+CASTSSHADOWS_UNIFORM_NAME, 0)
		} else {
			pl.ShadowMap.Bind(currentTextureUnit)
			// fmt.Println("Binding SpotLight to ", rnd.CurrentTextureUnit)
			s.SetUniformI(SPOT_LIGHTS_UNIFORM_NAME+"["+strconv.Itoa(int(arrayIndex))+"]."+SHADOWMAP_UNIFORM_NAME, int32(currentTextureUnit))
		}

	}
	return nil
}

func (this *SpotLight) InitShadowmap(width, height uint32) {
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
	camera.Init()
	camera.Position = this.Position
	camera.LookDirection = this.Direction.Add(mgl32.Vec3{1e-19, 1e-19, 1e-19})
	RenderMgr.SetCamera3D(&camera, 6)

	this.ShadowMap.SetAsTarget()
	Render.ClearScreen(&Color{0, 0, 0, 0}, 1.0)
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
	for i := 0; i < len(this.PointLights); i++ {
		this.PointLights[i].RenderShadowMap()
	}
}

type LightManager struct {
	LightCollections       []LightCollection
	CurrentLightCollection int32
}

func (this *LightManager) Init() {
	this.LightCollections = make([]LightCollection, 1)
	this.CurrentLightCollection = 0
	ResourceMgr.LoadShader(SHADOWMAP_SHADER_NAME, "shadowMapVert.glsl", "shadowMapFrag.glsl", "", "", "", "")
	ResourceMgr.LoadShader(SHADOWMAP_INSTANCED_SHADER_NAME, "shadowMapInstancedVert.glsl", "shadowMapFrag.glsl", "", "", "", "")
	ResourceMgr.LoadShader(POINT_LIGHT_SHADOWMAP_SHADER_NAME, "pointLightShadowMapVert.glsl", "pointLightShadowMapFrag.glsl", "pointLightShadowMapGeo.glsl", "", "", "")
	ResourceMgr.LoadShader(POINT_LIGHT_SHADOWMAP_INSTANCED_SHADER_NAME, "pointLightShadowMapInstancedVert.glsl", "pointLightShadowMapFrag.glsl", "pointLightShadowMapGeo.glsl", "", "", "")
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

var LightMgr LightManager
