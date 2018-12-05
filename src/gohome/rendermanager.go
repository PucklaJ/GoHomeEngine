package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
)

type Viewport struct {
	CameraIndex         int32
	X, Y, Width, Height int
	StrapToWindow       bool
}

type RenderManager struct {
	renderObjects      []RenderObject
	afterRenderObjects []RenderObject
	CurrentShader      Shader
	camera2Ds          []*Camera2D
	camera3Ds          []*Camera3D
	viewport2Ds        []*Viewport
	viewport3Ds        []*Viewport
	Projection2D       Projection
	Projection3D       Projection
	ForceShader3D      Shader
	ForceShader2D      Shader

	BackBufferMS RenderTexture
	BackBuffer   RenderTexture
	BackBuffer2D RenderTexture
	BackBuffer3D RenderTexture

	BackBufferShader     Shader
	PostProcessingShader Shader
	renderScreenShader   Shader

	currentCamera2D *Camera2D
	currentCamera3D *Camera3D
	currentViewport *Viewport

	EnableBackBuffer             bool
	WireFrameMode                bool
	UpdateProjectionWithViewport bool
	RenderToScreenFirst          bool
	AutoRender                   bool
	ReRender                     bool
}

func (rmgr *RenderManager) Init() {
	windowSize := Framew.WindowGetSize()

	rmgr.CurrentShader = nil
	rmgr.BackBufferMS = Render.CreateRenderTexture("BackBufferMS", uint32(windowSize[0]), uint32(windowSize[1]), 1, true, true, false, false)
	rmgr.BackBuffer = Render.CreateRenderTexture("BackBuffer", uint32(windowSize[0]), uint32(windowSize[1]), 1, true, false, false, false)
	rmgr.BackBuffer2D = Render.CreateRenderTexture("BackBuffer2D", uint32(windowSize[0]), uint32(windowSize[1]), 1, true, true, false, false)
	rmgr.BackBuffer3D = Render.CreateRenderTexture("BackBuffer3D", uint32(windowSize[0]), uint32(windowSize[1]), 1, true, true, false, false)
	ResourceMgr.LoadShaderSource("BackBufferShader", BACKBUFFER_SHADER_VERTEX_SOURCE_OPENGL, BACKBUFFER_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	ResourceMgr.LoadShaderSource("PostProcessingShader", POST_PROCESSING_SHADER_VERTEX_SOURCE_OPENGL, POST_PROCESSING_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	ResourceMgr.LoadShaderSource("RenderScreenShader", POST_PROCESSING_SHADER_VERTEX_SOURCE_OPENGL, RENDER_SCREEN_SHADER_FRAGMENT_SOURCE_OPENGL, "", "", "", "")
	rmgr.BackBufferShader = ResourceMgr.GetShader("BackBufferShader")
	rmgr.PostProcessingShader = ResourceMgr.GetShader("PostProcessingShader")
	rmgr.renderScreenShader = ResourceMgr.GetShader("RenderScreenShader")

	rmgr.AddViewport2D(&Viewport{
		0,
		0, 0,
		int(windowSize[0]),
		int(windowSize[1]),
		true,
	})
	rmgr.AddViewport3D(&Viewport{
		0,
		0, 0,
		int(windowSize[0]),
		int(windowSize[1]),
		true,
	})
	rmgr.SetProjection2D(&Ortho2DProjection{
		Left:   0.0,
		Right:  windowSize[0],
		Top:    0.0,
		Bottom: windowSize[1],
	})
	rmgr.SetProjection3D(&PerspectiveProjection{
		Width:     windowSize[0],
		Height:    windowSize[1],
		FOV:       70.0,
		NearPlane: 0.1,
		FarPlane:  1000.0,
	})

	rmgr.EnableBackBuffer = true
	rmgr.UpdateProjectionWithViewport = false
	rmgr.RenderToScreenFirst = false
	rmgr.AutoRender = true
	rmgr.ReRender = true
}

func (rmgr *RenderManager) AddObject(robj RenderObject) {
	if robj.RendersLast() {
		rmgr.afterRenderObjects = append(rmgr.afterRenderObjects, robj)
	} else {
		rmgr.renderObjects = append(rmgr.renderObjects, robj)
	}
}

func (rmgr *RenderManager) RemoveObject(robj RenderObject) {
	if robj.RendersLast() {
		for i := 0; i < len(rmgr.afterRenderObjects); i++ {
			if robj == rmgr.afterRenderObjects[i] {
				rmgr.afterRenderObjects = append(rmgr.afterRenderObjects[:i], rmgr.afterRenderObjects[i+1:]...)
				return
			}
		}
	} else {
		for i := 0; i < len(rmgr.renderObjects); i++ {
			if robj == rmgr.renderObjects[i] {
				rmgr.renderObjects = append(rmgr.renderObjects[:i], rmgr.renderObjects[i+1:]...)
				return
			}
		}
	}
}

func (rmgr *RenderManager) handleShader(robj RenderObject) {
	shader := robj.GetShader()
	if rmgr.ForceShader2D != nil && TYPE_2D.Compatible(robj.GetType()) {
		if rmgr.CurrentShader != rmgr.ForceShader2D {
			rmgr.CurrentShader = rmgr.ForceShader2D
			if rmgr.CurrentShader != nil {
				rmgr.CurrentShader.Use()
			}
		}
	} else if rmgr.ForceShader3D != nil && TYPE_3D.Compatible(robj.GetType()) {
		if rmgr.CurrentShader != rmgr.ForceShader3D {
			rmgr.CurrentShader = rmgr.ForceShader3D
			if rmgr.CurrentShader != nil {
				rmgr.CurrentShader.Use()
			}
		}
	} else {
		if rmgr.CurrentShader == nil {
			rmgr.CurrentShader = shader
			if rmgr.CurrentShader != nil {
				rmgr.CurrentShader.Use()
			}
		} else if rmgr.CurrentShader != shader {
			rmgr.CurrentShader.Unuse()
			rmgr.CurrentShader = shader
			if rmgr.CurrentShader != nil {
				rmgr.CurrentShader.Use()
			}
		}
	}
}

func (rmgr *RenderManager) updateCamera(robj RenderObject) {
	if TYPE_2D.Compatible(robj.GetType()) {
		if rmgr.currentCamera2D != nil && rmgr.CurrentShader != nil {
			rmgr.currentCamera2D.CalculateViewMatrix()
			rmgr.CurrentShader.SetUniformM3("viewMatrix2D", rmgr.currentCamera2D.GetViewMatrix())
		} else if rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformM3("viewMatrix2D", mgl32.Ident3())
		}
	} else {
		if rmgr.currentCamera3D != nil && rmgr.CurrentShader != nil {
			rmgr.currentCamera3D.CalculateViewMatrix()
			rmgr.CurrentShader.SetUniformM4("viewMatrix3D", rmgr.currentCamera3D.GetViewMatrix())
			rmgr.CurrentShader.SetUniformM4("inverseViewMatrix3D", rmgr.currentCamera3D.GetInverseViewMatrix())
		} else if rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformM4("viewMatrix3D", mgl32.Ident4())
			rmgr.CurrentShader.SetUniformM4("inverseViewMatrix3D", mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateProjection(t RenderType) {
	if TYPE_2D.Compatible(t) {
		if rmgr.Projection2D != nil && rmgr.CurrentShader != nil {
			rmgr.Projection2D.CalculateProjectionMatrix()
			rmgr.CurrentShader.SetUniformM4("projectionMatrix2D", rmgr.Projection2D.GetProjectionMatrix())
		} else if rmgr.Projection2D == nil && rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformM4("projectionMatrix2D", mgl32.Ident4())
		}
	} else {
		if rmgr.Projection3D != nil && rmgr.CurrentShader != nil {
			rmgr.Projection3D.CalculateProjectionMatrix()
			rmgr.CurrentShader.SetUniformM4("projectionMatrix3D", rmgr.Projection3D.GetProjectionMatrix())
		} else if rmgr.Projection3D == nil && rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformM4("projectionMatrix3D", mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateTransformMatrix(robj RenderObject) {
	if robj != nil && robj.GetTransformableObject() != nil {
		robj.GetTransformableObject().CalculateTransformMatrix(rmgr, robj.NotRelativeCamera())
		robj.GetTransformableObject().SetTransformMatrix(rmgr)
	} else {
		if robj.GetType() == TYPE_2D {
			rmgr.setTransformMatrix2D(mgl32.Ident3())
		} else {
			rmgr.setTransformMatrix3D(mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateLights(lightCollectionIndex int32, rtype RenderType) {
	if rtype.Compatible(TYPE_3D) {
		if rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformLights(lightCollectionIndex)
		}
	}
}

func (rmgr *RenderManager) GetBackBuffer() RenderTexture {
	return rmgr.BackBuffer
}

func (rmgr *RenderManager) render3D() {
	if rmgr.BackBuffer3D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer3D.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 0})
	}
	for i := 0; i < len(rmgr.viewport3Ds); i++ {
		rmgr.Render(TYPE_3D, rmgr.viewport3Ds[i].CameraIndex, int32(i), LightMgr.CurrentLightCollection)
	}
	if rmgr.BackBuffer3D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer3D.UnsetAsTarget()
	}
}

func (rmgr *RenderManager) render2D() {
	if rmgr.BackBuffer2D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer2D.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 0})
	}
	for i := 0; i < len(rmgr.viewport2Ds); i++ {
		rmgr.Render(TYPE_2D, rmgr.viewport2Ds[i].CameraIndex, int32(i), LightMgr.CurrentLightCollection)
	}
	if rmgr.BackBuffer2D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer2D.UnsetAsTarget()
	}
}

func (rmgr *RenderManager) renderBackBuffers() {
	if !rmgr.EnableBackBuffer {
		return
	}

	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.SetAsTarget()
		rmgr.clearToBackgroundColor()
	}

	if rmgr.BackBufferShader != nil {
		rmgr.BackBufferShader.Use()
		rmgr.BackBufferShader.SetUniformI("BackBuffer", 0)
		rmgr.BackBufferShader.SetUniformF("depth", 0.5)
	}
	if rmgr.BackBuffer3D != nil {
		rmgr.BackBuffer3D.Bind(0)
	}
	Render.RenderBackBuffer()

	if rmgr.BackBufferShader != nil {
		rmgr.BackBufferShader.SetUniformF("depth", 0.0)
	}
	if rmgr.BackBuffer2D != nil {
		rmgr.BackBuffer2D.Bind(0)
	}
	Render.RenderBackBuffer()
	if rmgr.BackBuffer2D != nil {
		rmgr.BackBuffer2D.Unbind(0)
	}
	if rmgr.BackBufferShader != nil {
		rmgr.BackBufferShader.Unuse()
	}

	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.UnsetAsTarget()
	}
}

func (rmgr *RenderManager) renderPostProcessing() {
	if !rmgr.EnableBackBuffer {
		return
	}

	if rmgr.BackBuffer != nil {
		rmgr.BackBuffer.SetAsTarget()
		Render.ClearScreen(Color{0, 0, 0, 0})
	}

	if rmgr.PostProcessingShader != nil {
		rmgr.PostProcessingShader.Use()
		rmgr.PostProcessingShader.SetUniformI("BackBuffer", 0)
	}
	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.Bind(0)
	}
	Render.RenderBackBuffer()
	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.Unbind(0)
	}
	if rmgr.PostProcessingShader != nil {
		rmgr.PostProcessingShader.Unuse()
	}

	if rmgr.BackBuffer != nil {
		rmgr.BackBuffer.UnsetAsTarget()
	}
}

func (rmgr *RenderManager) renderToScreen() {
	if !rmgr.EnableBackBuffer {
		return
	}

	if rmgr.renderScreenShader != nil {
		rmgr.renderScreenShader.Use()
		rmgr.renderScreenShader.SetUniformI("BackBuffer", 0)
	}
	if rmgr.BackBuffer != nil {
		rmgr.BackBuffer.Bind(0)
	}
	Render.RenderBackBuffer()
	if rmgr.BackBuffer != nil {
		rmgr.BackBuffer.Unbind(0)
	}
	if rmgr.renderScreenShader != nil {
		rmgr.renderScreenShader.Unuse()
	}
}

func (rmgr *RenderManager) Update() {
	defer func() {
		rmgr.ReRender = rmgr.AutoRender
	}()
	if !rmgr.ReRender {
		return
	}
	Render.ClearScreen(Render.GetBackgroundColor())
	if rmgr.RenderToScreenFirst {
		rmgr.renderToScreen()
	}
	rmgr.render3D()
	rmgr.render2D()
	rmgr.renderBackBuffers()
	rmgr.renderPostProcessing()
	if !rmgr.RenderToScreenFirst {
		rmgr.renderToScreen()
	}
	Framew.WindowSwap()
}

func (rmgr *RenderManager) handleCurrentCameraAndViewport(rtype RenderType, cameraIndex int32, viewportIndex int32) {
	if TYPE_2D.Compatible(rtype) {
		if cameraIndex == -1 || len(rmgr.camera2Ds) == 0 || int32(len(rmgr.camera2Ds)-1) < cameraIndex {
			rmgr.currentCamera2D = nil
		} else {
			rmgr.currentCamera2D = rmgr.camera2Ds[cameraIndex]
		}
		if viewportIndex == -1 || len(rmgr.viewport2Ds) == 0 || int32(len(rmgr.viewport2Ds)-1) < viewportIndex {
			rmgr.currentViewport = nil
		} else {
			rmgr.currentViewport = rmgr.viewport2Ds[viewportIndex]
		}

	} else if TYPE_3D.Compatible(rtype) {
		if cameraIndex == -1 || len(rmgr.camera3Ds) == 0 || int32(len(rmgr.camera3Ds)-1) < cameraIndex {
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

	if rmgr.currentViewport == nil {
		vp := Render.GetViewport()
		rmgr.currentViewport = &vp
		rmgr.currentViewport.StrapToWindow = false
	}

	if rmgr.currentViewport.StrapToWindow {
		var wSize mgl32.Vec2
		if !rmgr.EnableBackBuffer {
			wSize = Framew.WindowGetSize()
		} else {
			wSize = Render.GetNativeResolution()
		}
		rmgr.currentViewport.Width = int(wSize.X())
		rmgr.currentViewport.Height = int(wSize.Y())
		rmgr.currentViewport.X = 0
		rmgr.currentViewport.Y = 0
	}
	Render.SetViewport(*rmgr.currentViewport)
	if rmgr.UpdateProjectionWithViewport {
		if TYPE_2D.Compatible(rtype) {
			if rmgr.Projection2D != nil {
				rmgr.Projection2D.Update(*rmgr.currentViewport)
			}
		}
		if TYPE_3D.Compatible(rtype) {
			if rmgr.Projection3D != nil {
				rmgr.Projection3D.Update(*rmgr.currentViewport)
			}
		}
	}
}

func (rmgr *RenderManager) prepareRenderRenderObject(robj RenderObject, lightCollectionIndex int32) {
	rmgr.handleShader(robj)
	rmgr.updateCamera(robj)
	rmgr.updateProjection(robj.GetType())
	rmgr.updateTransformMatrix(robj)
	rmgr.updateLights(lightCollectionIndex, robj.GetType())
}

func (rmgr *RenderManager) renderRenderObject(robj RenderObject, lightCollectionIndex int32) {
	rmgr.prepareRenderRenderObject(robj, lightCollectionIndex)
	robj.Render()
}

func (rmgr *RenderManager) renderInnerLoop(rtype RenderType, robj RenderObject, lightCollectionIndex int32) {
	if !robj.IsVisible() || !rtype.Compatible(robj.GetType()) {
		return
	}
	Render.PreRender()
	if !robj.HasDepthTesting() {
		Render.SetDepthTesting(false)
	}
	rmgr.renderRenderObject(robj, lightCollectionIndex)
	if !robj.HasDepthTesting() {
		Render.SetDepthTesting(true)
	}
	Render.AfterRender()
}

func (rmgr *RenderManager) Render(rtype RenderType, cameraIndex int32, viewportIndex int32, lightCollectionIndex int32) {
	if len(rmgr.renderObjects) == 0 && len(rmgr.afterRenderObjects) == 0 {
		return
	}

	rmgr.handleCurrentCameraAndViewport(rtype, cameraIndex, viewportIndex)

	Render.SetWireFrame(rmgr.WireFrameMode)

	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.Use()
	}

	for i := 0; i < len(rmgr.renderObjects); i++ {
		rmgr.renderInnerLoop(rtype, rmgr.renderObjects[i], lightCollectionIndex)
	}

	for i := 0; i < len(rmgr.afterRenderObjects); i++ {
		rmgr.renderInnerLoop(rtype, rmgr.afterRenderObjects[i], lightCollectionIndex)
	}

	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.Unuse()
	}

	Render.SetWireFrame(false)
}

func (rmgr *RenderManager) RenderRenderObject(robj RenderObject) {
	rmgr.RenderRenderObjectAdv(robj, 0, -1)
}

func (rmgr *RenderManager) RenderRenderObjectAdv(robj RenderObject, cameraIndex, viewportIndex int32) {
	rmgr.handleCurrentCameraAndViewport(robj.GetType(), cameraIndex, viewportIndex)

	Render.SetWireFrame(rmgr.WireFrameMode)

	rmgr.CurrentShader = nil

	Render.PreRender()
	rmgr.renderRenderObject(robj, LightMgr.CurrentLightCollection)
	Render.AfterRender()

	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.Unuse()
		rmgr.CurrentShader = nil
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
	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.SetUniformM3("transformMatrix2D", mat)
	}
}

func (rmgr *RenderManager) setTransformMatrix3D(mat mgl32.Mat4) {
	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.SetUniformM4("transformMatrix3D", mat)
	}
}

func (rmgr *RenderManager) SetProjection2D(proj Projection) {
	rmgr.Projection2D = proj
}

func (rmgr *RenderManager) SetProjection3D(proj Projection) {
	rmgr.Projection3D = proj
}

func (rmgr *RenderManager) Terminate() {
	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.Terminate()
	}
	if rmgr.BackBuffer != nil {
		rmgr.BackBuffer.Terminate()
	}
	if rmgr.BackBuffer2D != nil {
		rmgr.BackBuffer2D.Terminate()
	}
	if rmgr.BackBuffer3D != nil {
		rmgr.BackBuffer3D.Terminate()
	}
	if len(rmgr.renderObjects) == 0 {
		return
	}
	rmgr.renderObjects = append(rmgr.renderObjects[:0], rmgr.renderObjects[len(rmgr.renderObjects):]...)
	rmgr.viewport2Ds = append(rmgr.viewport2Ds[:0], rmgr.viewport2Ds[len(rmgr.viewport2Ds):]...)
	rmgr.viewport3Ds = append(rmgr.viewport3Ds[:0], rmgr.viewport3Ds[len(rmgr.viewport3Ds):]...)
}

func (rmgr *RenderManager) UpdateViewports(current Viewport, previous Viewport) {
	var xRel, yRel, widthRel, heightRel float32

	for i := 0; i < len(rmgr.viewport2Ds); i++ {
		xRel = float32(rmgr.viewport2Ds[i].X) / float32(previous.X)
		yRel = float32(rmgr.viewport2Ds[i].Y) / float32(previous.Y)
		widthRel = float32(rmgr.viewport2Ds[i].Width) / float32(previous.Width)
		heightRel = float32(rmgr.viewport2Ds[i].Height) / float32(previous.Height)

		rmgr.viewport2Ds[i].X = int(xRel * float32(current.X))
		rmgr.viewport2Ds[i].Y = int(yRel * float32(current.Y))
		rmgr.viewport2Ds[i].Width = int(widthRel * float32(current.Width))
		rmgr.viewport2Ds[i].Height = int(heightRel * float32(current.Height))
	}

	for i := 0; i < len(rmgr.viewport3Ds); i++ {
		xRel = float32(rmgr.viewport3Ds[i].X) / float32(previous.X)
		yRel = float32(rmgr.viewport3Ds[i].Y) / float32(previous.Y)
		widthRel = float32(rmgr.viewport3Ds[i].Width) / float32(previous.Width)
		heightRel = float32(rmgr.viewport3Ds[i].Height) / float32(previous.Height)

		rmgr.viewport3Ds[i].X = int(xRel * float32(current.X))
		rmgr.viewport3Ds[i].Y = int(yRel * float32(current.Y))
		rmgr.viewport3Ds[i].Width = int(widthRel * float32(current.Width))
		rmgr.viewport3Ds[i].Height = int(heightRel * float32(current.Height))
	}
}

func (rmgr *RenderManager) clearToBackgroundColor() {
	backg := Render.GetBackgroundColor()
	var r, g, b uint32
	if backg != nil {
		r, g, b, _ = backg.RGBA()
	} else {
		r, g, b = 0, 0, 0
	}

	newCol := Color{uint8(float32(r) / float32(0xffff) * 255.0),
		uint8(float32(g) / float32(0xffff) * 255.0),
		uint8(float32(b) / float32(0xffff) * 255.0),
		0}

	Render.ClearScreen(newCol)
}

func (rmgr *RenderManager) SetProjection2DToTexture(texture Texture) {
	rmgr.Projection2D = &Ortho2DProjection{
		Left:   0.0,
		Top:    0.0,
		Right:  float32(texture.GetWidth()),
		Bottom: float32(texture.GetHeight()),
	}
}

func (rmgr *RenderManager) UpdateProjection2D(viewportIndex int32) {
	if viewportIndex >= 0 && viewportIndex < int32(len(rmgr.viewport2Ds)) {
		rmgr.Projection2D.Update(*rmgr.viewport2Ds[viewportIndex])
	}
}

func (rmgr *RenderManager) NumRenderObjects() uint32 {
	return uint32(len(rmgr.renderObjects) + len(rmgr.afterRenderObjects))
}

var RenderMgr RenderManager
