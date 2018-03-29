package gohome

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderType uint8

const (
	TYPE_2D RenderType = iota
	TYPE_3D RenderType = iota
)

type RenderObject interface {
	Render()
	SetShader(s Shader)
	GetShader() Shader
	GetType() RenderType
	IsVisible() bool
	NotRelativeCamera() int
}

type TransformableObject interface {
	CalculateTransformMatrix(rmgr *RenderManager, notRelativeToCamera int)
	SetTransformMatrix(rmgr *RenderManager)
}

type RenderPair struct {
	RenderObject
	TransformableObject
}

type Viewport struct {
	CameraIndex         uint32
	X, Y, Width, Height int
}

type RenderManager struct {
	renderObjects []RenderPair
	currentShader Shader
	camera2Ds     []*Camera2D
	camera3Ds     []*Camera3D
	viewport2Ds   []*Viewport
	viewport3Ds   []*Viewport
	Projection2D  Projection
	Projection3D  Projection
	ForceShader3D Shader
	ForceShader2D Shader

	backBufferMS     RenderTexture
	backBuffer       RenderTexture
	backBuffer2D     RenderTexture
	backBuffer3D     RenderTexture
	backBufferShader Shader

	PostProcessingShader Shader
	renderScreenShader   Shader

	WireFrameMode bool

	currentCamera2D *Camera2D
	currentCamera3D *Camera3D
	currentViewport *Viewport
}

func (rmgr *RenderManager) Init() {
	rmgr.currentShader = nil
	rmgr.backBufferMS = Render.CreateRenderTexture("BackBufferMS", uint32(Framew.WindowGetSize()[0]), uint32(Framew.WindowGetSize()[1]), 1, true, true, false)
	rmgr.backBuffer = Render.CreateRenderTexture("BackBuffer", uint32(Framew.WindowGetSize()[0]), uint32(Framew.WindowGetSize()[1]), 1, true, false, false)
	rmgr.backBuffer2D = Render.CreateRenderTexture("BackBuffer2D", uint32(Framew.WindowGetSize()[0]), uint32(Framew.WindowGetSize()[1]), 1, true, true, false)
	rmgr.backBuffer3D = Render.CreateRenderTexture("BackBuffer3D", uint32(Framew.WindowGetSize()[0]), uint32(Framew.WindowGetSize()[1]), 1, true, true, false)
	ResourceMgr.LoadShader("BackBufferShader", "backBufferShaderVert.glsl", "backBufferShaderFrag.glsl", "", "", "", "")
	ResourceMgr.LoadShader("PostProcessingShader", "postProcessingShaderVert.glsl", "postProcessingShaderFrag.glsl", "", "", "", "")
	ResourceMgr.LoadShader("RenderScreenShader", "postProcessingShaderVert.glsl", "renderScreenFrag.glsl", "", "", "", "")
	rmgr.backBufferShader = ResourceMgr.GetShader("BackBufferShader")
	rmgr.PostProcessingShader = ResourceMgr.GetShader("PostProcessingShader")
	rmgr.renderScreenShader = ResourceMgr.GetShader("RenderScreenShader")
	rmgr.AddViewport2D(&Viewport{
		0,
		0, 0,
		int(Framew.WindowGetSize()[0]),
		int(Framew.WindowGetSize()[1]),
	})
	rmgr.AddViewport3D(&Viewport{
		0,
		0, 0,
		int(Framew.WindowGetSize()[0]),
		int(Framew.WindowGetSize()[1]),
	})
}

func (rmgr *RenderManager) AddObject(robj RenderObject, tobj TransformableObject) {
	rmgr.renderObjects = append(rmgr.renderObjects, RenderPair{RenderObject: robj, TransformableObject: tobj})
}

func (rmgr *RenderManager) RemoveObject(robj RenderObject, tobj TransformableObject) {
	for i := 0; i < len(rmgr.renderObjects); i++ {
		if robj == rmgr.renderObjects[i].RenderObject && tobj == rmgr.renderObjects[i].TransformableObject {
			rmgr.renderObjects = append(rmgr.renderObjects[:i], rmgr.renderObjects[i+1:]...)
			return
		}
	}
}

func (rmgr *RenderManager) handleShader(robj RenderObject) {
	shader := robj.GetShader()
	if rmgr.ForceShader2D != nil && robj.GetType() == TYPE_2D {
		if rmgr.currentShader != rmgr.ForceShader2D {
			rmgr.currentShader = rmgr.ForceShader2D
			if rmgr.currentShader != nil {
				rmgr.currentShader.Use()
			}
		}
	} else if rmgr.ForceShader3D != nil && robj.GetType() == TYPE_3D {
		if rmgr.currentShader != rmgr.ForceShader3D {
			rmgr.currentShader = rmgr.ForceShader3D
			if rmgr.currentShader != nil {
				rmgr.currentShader.Use()
			}
		}
	} else {
		if rmgr.currentShader == nil {
			rmgr.currentShader = shader
			if rmgr.currentShader != nil {
				rmgr.currentShader.Use()
			}
		} else if rmgr.currentShader != shader {
			rmgr.currentShader.Unuse()
			rmgr.currentShader = shader
			if rmgr.currentShader != nil {
				rmgr.currentShader.Use()
			}
		}
	}
}

func (rmgr *RenderManager) updateCamera(robj RenderObject) {
	if robj.GetType() == TYPE_2D {
		if rmgr.currentCamera2D != nil && rmgr.currentShader != nil {
			rmgr.currentCamera2D.CalculateViewMatrix()
			rmgr.currentShader.SetUniformM3("viewMatrix2D", rmgr.currentCamera2D.GetViewMatrix())
		} else if rmgr.currentShader != nil {
			rmgr.currentShader.SetUniformM3("viewMatrix2D", mgl32.Ident3())
		}
	} else {
		if rmgr.currentCamera3D != nil && rmgr.currentShader != nil {
			rmgr.currentCamera3D.CalculateViewMatrix()
			rmgr.currentShader.SetUniformM4("viewMatrix3D", rmgr.currentCamera3D.GetViewMatrix())
		} else if rmgr.currentShader != nil {
			rmgr.currentShader.SetUniformM4("viewMatrix3D", mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateProjection(t RenderType) {
	if t == TYPE_2D {
		if rmgr.Projection2D != nil && rmgr.currentShader != nil {
			rmgr.Projection2D.CalculateProjectionMatrix()
			rmgr.currentShader.SetUniformM4("projectionMatrix2D", rmgr.Projection2D.GetProjectionMatrix())
		} else if rmgr.Projection2D == nil && rmgr.currentShader != nil {
			rmgr.currentShader.SetUniformM4("projectionMatrix2D", mgl32.Ident4())
		}
	} else {
		if rmgr.Projection3D != nil && rmgr.currentShader != nil {
			rmgr.Projection3D.CalculateProjectionMatrix()
			rmgr.currentShader.SetUniformM4("projectionMatrix3D", rmgr.Projection3D.GetProjectionMatrix())
		} else if rmgr.Projection3D == nil && rmgr.currentShader != nil {
			rmgr.currentShader.SetUniformM4("projectionMatrix3D", mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateTransformMatrix(robj *RenderPair, t RenderType) {
	if robj != nil && robj.TransformableObject != nil {
		robj.TransformableObject.CalculateTransformMatrix(rmgr, robj.RenderObject.NotRelativeCamera())
		robj.TransformableObject.SetTransformMatrix(rmgr)
	} else {
		if t == TYPE_2D {
			rmgr.setTransformMatrix2D(mgl32.Ident3())
		} else {
			rmgr.setTransformMatrix3D(mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateLights(lightCollectionIndex int32) {
	if rmgr.currentShader != nil {
		if err := rmgr.currentShader.SetUniformLights(lightCollectionIndex); err != nil {
			// fmt.Println("Error:", err)
		}
	}
}

func (rmgr *RenderManager) renderBackBuffers() {
	rmgr.backBufferMS.SetAsTarget()

	rmgr.backBufferShader.Use()
	rmgr.backBufferShader.SetUniformI("backBuffer", 0)
	rmgr.backBufferShader.SetUniformF("depth", 0.5)
	rmgr.backBuffer3D.Bind(0)
	Render.RenderBackBuffer()

	rmgr.backBufferShader.SetUniformF("depth", 0.0)
	rmgr.backBuffer2D.Bind(0)
	Render.RenderBackBuffer()
	rmgr.backBuffer2D.Unbind(0)
	rmgr.backBufferShader.Unuse()

	rmgr.backBufferMS.UnsetAsTarget()
}

func (rmgr *RenderManager) render3D() {
	rmgr.backBuffer3D.SetAsTarget()
	for i := 0; i < len(rmgr.viewport3Ds); i++ {
		rmgr.Render(TYPE_3D, rmgr.viewport3Ds[i].CameraIndex, int32(i), LightMgr.CurrentLightCollection)
	}
	rmgr.backBuffer3D.UnsetAsTarget()
}

func (rmgr *RenderManager) render2D() {
	rmgr.backBuffer2D.SetAsTarget()
	for i := 0; i < len(rmgr.viewport2Ds); i++ {
		rmgr.Render(TYPE_2D, rmgr.viewport2Ds[i].CameraIndex, int32(i), LightMgr.CurrentLightCollection)
	}
	rmgr.backBuffer2D.UnsetAsTarget()
}

func (rmgr *RenderManager) GetBackBuffer() RenderTexture {
	return rmgr.backBuffer
}

func (rmgr *RenderManager) renderPostProcessing() {
	rmgr.backBuffer.SetAsTarget()

	rmgr.PostProcessingShader.Use()
	rmgr.PostProcessingShader.SetUniformI("backBuffer", 0)
	rmgr.backBufferMS.Bind(0)
	Render.RenderBackBuffer()
	rmgr.backBufferMS.Unbind(0)
	rmgr.PostProcessingShader.Unuse()

	rmgr.backBuffer.UnsetAsTarget()
}

func (rmgr *RenderManager) renderToScreen() {
	rmgr.renderScreenShader.Use()
	rmgr.renderScreenShader.SetUniformI("backBuffer", 0)
	rmgr.backBuffer.Bind(0)
	Render.RenderBackBuffer()
	rmgr.backBuffer.Unbind(0)
	rmgr.renderScreenShader.Unuse()
}

func (rmgr *RenderManager) Update() {
	rmgr.render3D()
	rmgr.render2D()
	rmgr.renderBackBuffers()
	rmgr.renderPostProcessing()
	rmgr.renderToScreen()
}

func (rmgr *RenderManager) handleCurrentCameraAndViewport(rtype RenderType, cameraIndex uint32, viewportIndex int32) {
	if rtype == TYPE_2D {
		if len(rmgr.camera2Ds) == 0 || uint32(len(rmgr.camera2Ds)-1) < cameraIndex {
			rmgr.currentCamera2D = nil
		} else {
			rmgr.currentCamera2D = rmgr.camera2Ds[cameraIndex]
		}
		if viewportIndex == -1 || len(rmgr.viewport2Ds) == 0 || int32(len(rmgr.viewport2Ds)-1) < viewportIndex {
			rmgr.currentViewport = nil
		} else {
			rmgr.currentViewport = rmgr.viewport2Ds[viewportIndex]
		}

	} else if rtype == TYPE_3D {
		if len(rmgr.camera3Ds) == 0 || uint32(len(rmgr.camera3Ds)-1) < cameraIndex {
			rmgr.currentCamera3D = nil
		} else {
			rmgr.currentCamera3D = rmgr.camera3Ds[cameraIndex]
		}
		if viewportIndex == -1 || len(rmgr.viewport3Ds) == 0 || int32(len(rmgr.viewport3Ds)-1) < viewportIndex {
			rmgr.currentViewport = nil
		} else {
			rmgr.currentViewport = rmgr.viewport3Ds[viewportIndex]
		}
	}

	if rmgr.currentViewport != nil {
		Render.SetViewport(*rmgr.currentViewport)
		if rtype == TYPE_2D {
			if rmgr.Projection2D != nil {
				rmgr.Projection2D.Update(*rmgr.currentViewport)
			}
		}
		if rtype == TYPE_3D {
			if rmgr.Projection3D != nil {
				rmgr.Projection3D.Update(*rmgr.currentViewport)
			}
		}
	}
}

func (rmgr *RenderManager) prepareRenderRenderObject(robj *RenderPair, lightCollectionIndex int32) {
	rmgr.handleShader(robj.RenderObject)
	rmgr.updateCamera(robj.RenderObject)
	rmgr.updateProjection(robj.RenderObject.GetType())
	rmgr.updateTransformMatrix(robj, robj.RenderObject.GetType())
	rmgr.updateLights(lightCollectionIndex)
}

func (rmgr *RenderManager) renderRenderObject(robj *RenderPair, lightCollectionIndex int32) {
	rmgr.prepareRenderRenderObject(robj, lightCollectionIndex)
	robj.Render()
}

func (rmgr *RenderManager) Render(rtype RenderType, cameraIndex uint32, viewportIndex int32, lightCollectionIndex int32) {
	if len(rmgr.renderObjects) == 0 {
		return
	}

	rmgr.handleCurrentCameraAndViewport(rtype, cameraIndex, viewportIndex)

	Render.SetWireFrame(rmgr.WireFrameMode)

	if rmgr.currentShader != nil {
		rmgr.currentShader.Use()
	}

	for i := 0; i < len(rmgr.renderObjects); i++ {
		if rtype != rmgr.renderObjects[i].GetType() || !rmgr.renderObjects[i].IsVisible() {
			continue
		}
		Render.PreRender()
		rmgr.renderRenderObject(&rmgr.renderObjects[i], lightCollectionIndex)
		Render.AfterRender()
	}

	if rmgr.currentShader != nil {
		rmgr.currentShader.Unuse()
	}

	Render.SetWireFrame(false)
}

func (rmgr *RenderManager) SetCamera2D(cam *Camera2D, index uint32) {
	if len(rmgr.camera2Ds) == 0 {
		rmgr.camera2Ds = make([]*Camera2D, 1)
	}
	if uint32(len(rmgr.camera2Ds)-1) < index {
		rmgr.camera2Ds = append(rmgr.camera2Ds, make([]*Camera2D, index-uint32(len(rmgr.camera2Ds)-1))...)
	}
	rmgr.camera2Ds[index] = cam
}

func (rmgr *RenderManager) SetCamera3D(cam *Camera3D, index uint32) {
	if len(rmgr.camera3Ds) == 0 {
		rmgr.camera3Ds = make([]*Camera3D, 1)
	}
	if uint32(len(rmgr.camera3Ds)-1) < index {
		rmgr.camera3Ds = append(rmgr.camera3Ds, make([]*Camera3D, index-uint32(len(rmgr.camera3Ds)-1))...)
	}
	rmgr.camera3Ds[index] = cam
}

func (rmgr *RenderManager) AddViewport2D(viewport *Viewport) {
	rmgr.viewport2Ds = append(rmgr.viewport2Ds, viewport)
}

func (rmgr *RenderManager) AddViewport3D(viewport *Viewport) {
	rmgr.viewport3Ds = append(rmgr.viewport3Ds, viewport)
}

func (rmgr *RenderManager) SetViewport2D(viewport *Viewport, index uint32) {
	if len(rmgr.viewport2Ds) == 0 {
		rmgr.viewport2Ds = make([]*Viewport, 1)
	}
	if uint32(len(rmgr.viewport2Ds)-1) < index {
		rmgr.viewport2Ds = append(rmgr.viewport2Ds, make([]*Viewport, index-uint32(len(rmgr.viewport2Ds)-1))...)
	}
	rmgr.viewport2Ds[index] = viewport
}

func (rmgr *RenderManager) SetViewport3D(viewport *Viewport, index uint32) {
	if len(rmgr.viewport3Ds) == 0 {
		rmgr.viewport3Ds = make([]*Viewport, 1)
	}
	if uint32(len(rmgr.viewport3Ds)-1) < index {
		rmgr.viewport3Ds = append(rmgr.viewport3Ds, make([]*Viewport, index-uint32(len(rmgr.viewport3Ds)-1))...)
	}
	rmgr.viewport3Ds[index] = viewport
}

func (rmgr *RenderManager) setTransformMatrix2D(mat mgl32.Mat3) {
	if rmgr.currentShader != nil {
		rmgr.currentShader.SetUniformM3("transformMatrix2D", mat)
	}
}

func (rmgr *RenderManager) setTransformMatrix3D(mat mgl32.Mat4) {
	if rmgr.currentShader != nil {
		rmgr.currentShader.SetUniformM4("transformMatrix3D", mat)
	}
}

func (rmgr *RenderManager) SetProjection2D(proj Projection) {
	rmgr.Projection2D = proj
}

func (rmgr *RenderManager) SetProjection3D(proj Projection) {
	rmgr.Projection3D = proj
}

func (rmgr *RenderManager) Terminate() {
	rmgr.backBufferMS.Terminate()
	rmgr.backBuffer.Terminate()
	rmgr.backBuffer2D.Terminate()
	rmgr.backBuffer3D.Terminate()
	if len(rmgr.renderObjects) == 0 {
		return
	}
	rmgr.renderObjects = append(rmgr.renderObjects[:0], rmgr.renderObjects[len(rmgr.renderObjects):]...)
}

var RenderMgr RenderManager
