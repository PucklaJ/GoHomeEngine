package gohome

import (
	"math"
	"runtime"
	"sync"

	"github.com/PucklaJ/mathgl/mgl32"
)

// A viewport that is on a certain part of the screen.
// Showing a certain part of the world
type Viewport struct {
	// The index of the camera belonging to this viewport
	CameraIndex int
	// The position and dimensions of the viewport
	X, Y, Width, Height int
	// Wether the viewport should adjust base don the window size
	StrapToWindow bool
}

// The manager that handles the rendering of all objects
type RenderManager struct {
	renderObjects      []RenderObject
	afterRenderObjects []RenderObject
	CurrentShader      Shader
	camera2Ds          []*Camera2D
	camera3Ds          []*Camera3D
	viewport2Ds        []*Viewport
	viewport3Ds        []*Viewport
	// The projetion used for 2D objects
	Projection2D Projection
	// The projection used for 3D objects
	Projection3D Projection
	// If set this shader is forced onto every 3D object
	ForceShader3D Shader
	// If set this shader is forced onto every 2D object
	ForceShader2D Shader

	// The back buffer that will be rendered to the screen
	// and onto which BackBuffer2D and 3D will be rendered
	BackBufferMS RenderTexture
	// The BackBuffer to which all 2D objects will be rendered
	BackBuffer2D RenderTexture
	// The BackBuffer to which all 3D objects will be rendered
	BackBuffer3D RenderTexture

	// The shader used for rendering the back buffers
	BackBufferShader Shader

	currentCamera2D *Camera2D
	currentCamera3D *Camera3D
	currentViewport *Viewport

	// Wether the objects should be rendered to the back buffers or directly to the screen
	EnableBackBuffer bool
	// Wether the objects should be rendered in wire frame mode
	WireFrameMode bool
	// Wether the projection should be updated every frame based on the viewport
	UpdateProjectionWithViewport bool
	// Wether the back buffers of the last frame should be rendered before the objects
	RenderToScreenFirst bool
	// If false ReRender must be set to true everytime you want to re-render the scene
	AutoRender bool
	// If true the scene will be re-rendered
	ReRender bool

	calculatingTransformMatricesParallel bool
}

// Initialises all values of the manager
func (rmgr *RenderManager) Init() {
	if Render.HasFunctionAvailable("MULTISAMPLE") {
		bn, bv, bf := GenerateShaderBackBuffer(0)
		ls(bn, bv, bf)
	} else {
		bn, bv, bf := GenerateShaderBackBuffer(SHADER_FLAG_NO_MS)
		ls(bn, bv, bf)
	}

	windowSize := Framew.WindowGetSize()

	rmgr.CurrentShader = nil
	rmgr.BackBufferMS = Render.CreateRenderTexture("BackBufferMS", int(windowSize[0]), int(windowSize[1]), 1, true, true, false, false)
	rmgr.BackBuffer2D = Render.CreateRenderTexture("BackBuffer2D", int(windowSize[0]), int(windowSize[1]), 1, true, true, false, false)
	rmgr.BackBuffer3D = Render.CreateRenderTexture("BackBuffer3D", int(windowSize[0]), int(windowSize[1]), 1, true, true, false, false)
	rmgr.BackBufferShader = ResourceMgr.GetShader("BackBufferShader")

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

// Adds a RenderObject to the scene so that it will be rendered
func (rmgr *RenderManager) AddObject(robj RenderObject) {
	if robj.RendersLast() {
		rmgr.afterRenderObjects = append(rmgr.afterRenderObjects, robj)
	} else {
		rmgr.renderObjects = append(rmgr.renderObjects, robj)
	}
}

// Removes a RenderObject from the scene so that it won't be rendered
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

func (rmgr *RenderManager) calculateTransformMatrices(rtype RenderType) {
	rmgr.calculatingTransformMatricesParallel = true
	calcFunc := func(_robj RenderObject, wg *sync.WaitGroup) {
		tobj := _robj.GetTransformableObject()
		if tobj != nil {
			tobj.CalculateTransformMatrix(rmgr, _robj.NotRelativeCamera())
		}
		wg.Done()
	}
	var wg sync.WaitGroup
	for _, robj := range rmgr.renderObjects {
		if rtype.Compatible(robj.GetType()) {
			wg.Add(1)
			go calcFunc(robj, &wg)
		}
	}
	for _, arobj := range rmgr.afterRenderObjects {
		if rtype.Compatible(arobj.GetType()) {
			wg.Add(1)
			go calcFunc(arobj, &wg)
		}
	}
	wg.Wait()
	rmgr.calculatingTransformMatricesParallel = false
}

func (rmgr *RenderManager) updateTransformMatrix(robj RenderObject) {
	if robj != nil && robj.GetTransformableObject() != nil {
		if runtime.GOOS == "android" {
			robj.GetTransformableObject().CalculateTransformMatrix(rmgr, robj.NotRelativeCamera())
		}
		robj.GetTransformableObject().SetTransformMatrix(rmgr)
	} else {
		if TYPE_2D.Compatible(robj.GetType()) {
			rmgr.setTransformMatrix2D(mgl32.Ident3())
		} else {
			rmgr.setTransformMatrix3D(mgl32.Ident4())
		}
	}
}

func (rmgr *RenderManager) updateLights(lightCollectionIndex int, rtype RenderType) {
	if TYPE_3D.Compatible(rtype) {
		if rmgr.CurrentShader != nil {
			rmgr.CurrentShader.SetUniformLights(lightCollectionIndex)
		}
	}
}

// Returns the back buffer as a RenderTexture
func (rmgr *RenderManager) GetBackBuffer() RenderTexture {
	return rmgr.BackBufferMS
}

func (rmgr *RenderManager) render3D() {
	if rmgr.BackBuffer3D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer3D.SetAsTarget()
		Render.ClearScreen(nil)
	}
	for i := 0; i < len(rmgr.viewport3Ds); i++ {
		rmgr.Render(TYPE_3D, rmgr.viewport3Ds[i].CameraIndex, i, LightMgr.CurrentLightCollection)
	}
	if rmgr.BackBuffer3D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer3D.UnsetAsTarget()
	}
}

func (rmgr *RenderManager) render2D() {
	if rmgr.BackBuffer2D != nil && rmgr.EnableBackBuffer {
		rmgr.BackBuffer2D.SetAsTarget()
		Render.ClearScreen(nil)
	}
	for i := 0; i < len(rmgr.viewport2Ds); i++ {
		rmgr.Render(TYPE_2D, rmgr.viewport2Ds[i].CameraIndex, i, LightMgr.CurrentLightCollection)
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
		rmgr.BackBufferShader.SetUniformI("texture0", 0)
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

func (rmgr *RenderManager) renderToScreen() {
	if !rmgr.EnableBackBuffer {
		return
	}
	if !Render.HasFunctionAvailable("MULTISAMPLE") && Render.HasFunctionAvailable("BLIT_FRAMEBUFFER_SCREEN") {
		rmgr.BackBufferMS.Blit(nil)
	} else {
		if rmgr.BackBufferShader != nil {
			rmgr.BackBufferShader.Use()
			rmgr.BackBufferShader.SetUniformI("texture0", 0)
			rmgr.BackBufferShader.SetUniformF("depth", -1.0)
		}
		if rmgr.BackBufferMS != nil {
			rmgr.BackBufferMS.Bind(0)
		}
		Render.RenderBackBuffer()
		if rmgr.BackBufferMS != nil {
			rmgr.BackBufferMS.Unbind(0)
		}
		if rmgr.BackBufferShader != nil {
			rmgr.BackBufferShader.Unuse()
		}
	}
}

// Updates the manger / renders everything
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
	if !rmgr.RenderToScreenFirst {
		rmgr.renderToScreen()
	}
	Framew.WindowSwap()
}

func (rmgr *RenderManager) handleCurrentCameraAndViewport(rtype RenderType, cameraIndex, viewportIndex int) {
	if TYPE_2D.Compatible(rtype) {
		if cameraIndex == -1 || len(rmgr.camera2Ds) == 0 || len(rmgr.camera2Ds)-1 < cameraIndex {
			rmgr.currentCamera2D = nil
		} else {
			rmgr.currentCamera2D = rmgr.camera2Ds[cameraIndex]
		}
		if viewportIndex == -1 || len(rmgr.viewport2Ds) == 0 || len(rmgr.viewport2Ds)-1 < viewportIndex {
			rmgr.currentViewport = nil
		} else {
			rmgr.currentViewport = rmgr.viewport2Ds[viewportIndex]
		}

	} else if TYPE_3D.Compatible(rtype) {
		if cameraIndex == -1 || len(rmgr.camera3Ds) == 0 || len(rmgr.camera3Ds)-1 < cameraIndex {
			rmgr.currentCamera3D = nil
		} else {
			rmgr.currentCamera3D = rmgr.camera3Ds[cameraIndex]
		}
		if viewportIndex == -1 || len(rmgr.viewport3Ds) == 0 || len(rmgr.viewport3Ds)-1 < viewportIndex {
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

func (rmgr *RenderManager) prepareRenderRenderObject(robj RenderObject, lightCollectionIndex int) {
	rmgr.handleShader(robj)
	rmgr.updateCamera(robj)
	rmgr.updateProjection(robj.GetType())
	rmgr.updateTransformMatrix(robj)
	rmgr.updateLights(lightCollectionIndex, robj.GetType())
}

func (rmgr *RenderManager) renderRenderObject(robj RenderObject, lightCollectionIndex int) {
	rmgr.prepareRenderRenderObject(robj, lightCollectionIndex)
	robj.Render()
}

func (rmgr *RenderManager) renderInnerLoop(rtype RenderType, robj RenderObject, lightCollectionIndex int) {
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

// Renders a certain render type to a certain viewport using a certain light collection
func (rmgr *RenderManager) Render(rtype RenderType, cameraIndex, viewportIndex, lightCollectionIndex int) {
	if len(rmgr.renderObjects) == 0 && len(rmgr.afterRenderObjects) == 0 {
		return
	}

	rmgr.handleCurrentCameraAndViewport(rtype, cameraIndex, viewportIndex)

	Render.SetWireFrame(rmgr.WireFrameMode)

	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.Use()
	}
	if runtime.GOOS != "android" {
		rmgr.calculateTransformMatrices(rtype)
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

// Renders a RenderObject (used for custom rendering)
func (rmgr *RenderManager) RenderRenderObject(robj RenderObject) {
	rmgr.RenderRenderObjectAdv(robj, 0, -1)
}

// Same as RenderRenderObject but with additional arguments for camera and viewport
func (rmgr *RenderManager) RenderRenderObjectAdv(robj RenderObject, cameraIndex, viewportIndex int) {
	rmgr.handleCurrentCameraAndViewport(robj.GetType(), cameraIndex, viewportIndex)

	Render.SetWireFrame(rmgr.WireFrameMode)

	rmgr.CurrentShader = nil

	if runtime.GOOS != "android" {
		if tobj := robj.GetTransformableObject(); tobj != nil {
			tobj.CalculateTransformMatrix(rmgr, robj.NotRelativeCamera())
		}
	}

	Render.PreRender()
	rmgr.renderRenderObject(robj, LightMgr.CurrentLightCollection)
	Render.AfterRender()

	if rmgr.CurrentShader != nil {
		rmgr.CurrentShader.Unuse()
		rmgr.CurrentShader = nil
	}

	Render.SetWireFrame(false)
}

// Attaches a 2D camera to an index
func (rmgr *RenderManager) SetCamera2D(cam *Camera2D, index int) {
	if len(rmgr.camera2Ds) == 0 {
		rmgr.camera2Ds = make([]*Camera2D, 1)
	}
	if len(rmgr.camera2Ds)-1 < index {
		rmgr.camera2Ds = append(rmgr.camera2Ds, make([]*Camera2D, index-len(rmgr.camera2Ds)-1)...)
	}
	rmgr.camera2Ds[index] = cam
}

// Attaches a 3D camera to an index
func (rmgr *RenderManager) SetCamera3D(cam *Camera3D, index int) {
	if len(rmgr.camera3Ds) == 0 {
		rmgr.camera3Ds = make([]*Camera3D, 1)
	}
	if len(rmgr.camera3Ds)-1 < index {
		rmgr.camera3Ds = append(rmgr.camera3Ds, make([]*Camera3D, index-len(rmgr.camera3Ds)-1)...)
	}
	rmgr.camera3Ds[index] = cam
}

// Adds a 2D viewport to the scene
func (rmgr *RenderManager) AddViewport2D(viewport *Viewport) {
	rmgr.viewport2Ds = append(rmgr.viewport2Ds, viewport)
}

// Adds a 3D viewport to the scene
func (rmgr *RenderManager) AddViewport3D(viewport *Viewport) {
	rmgr.viewport3Ds = append(rmgr.viewport3Ds, viewport)
}

// Sets the 2D viewport of a certain index
func (rmgr *RenderManager) SetViewport2D(viewport *Viewport, index int) {
	if len(rmgr.viewport2Ds) == 0 {
		rmgr.viewport2Ds = make([]*Viewport, 1)
	}
	if len(rmgr.viewport2Ds)-1 < index {
		rmgr.viewport2Ds = append(rmgr.viewport2Ds, make([]*Viewport, index-len(rmgr.viewport2Ds)-1)...)
	}
	rmgr.viewport2Ds[index] = viewport
}

// Sets the 3D viewport of a certain index
func (rmgr *RenderManager) SetViewport3D(viewport *Viewport, index int) {
	if len(rmgr.viewport3Ds) == 0 {
		rmgr.viewport3Ds = make([]*Viewport, 1)
	}
	if len(rmgr.viewport3Ds)-1 < index {
		rmgr.viewport3Ds = append(rmgr.viewport3Ds, make([]*Viewport, index-len(rmgr.viewport3Ds)-1)...)
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

// Sets the projection used for 2D rendering
func (rmgr *RenderManager) SetProjection2D(proj Projection) {
	rmgr.Projection2D = proj
}

// Sets the projection used for 3D rendering
func (rmgr *RenderManager) SetProjection3D(proj Projection) {
	rmgr.Projection3D = proj
}

// Cleans everything up
func (rmgr *RenderManager) Terminate() {
	if rmgr.BackBufferMS != nil {
		rmgr.BackBufferMS.Terminate()
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

// Updates the viewport based on a viewport of the GPU
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
	var r, g, b, a uint32
	if backg != nil {
		r, g, b, a = backg.RGBA()
	} else {
		r, g, b, a = 0, 0, 0, math.MaxUint32
	}

	newCol := Color{
		uint8(float32(r) / float32(0xffff) * 255.0),
		uint8(float32(g) / float32(0xffff) * 255.0),
		uint8(float32(b) / float32(0xffff) * 255.0),
		uint8(float32(a) / float32(0xffff) * 255.0),
	}

	Render.ClearScreen(newCol)
}

// Takes the dimensions of a texture and uses it for a projection
func (rmgr *RenderManager) SetProjection2DToTexture(texture Texture) {
	rmgr.Projection2D = &Ortho2DProjection{
		Left:   0.0,
		Top:    0.0,
		Right:  float32(texture.GetWidth()),
		Bottom: float32(texture.GetHeight()),
	}
}

// Updates the 2D projection using the viewport of viewportIndex
func (rmgr *RenderManager) UpdateProjection2D(viewportIndex int32) {
	if viewportIndex >= 0 && viewportIndex < int32(len(rmgr.viewport2Ds)) {
		rmgr.Projection2D.Update(*rmgr.viewport2Ds[viewportIndex])
	}
}

// Returns the number of currently added RenderObjects
func (rmgr *RenderManager) NumRenderObjects() int {
	return len(rmgr.renderObjects) + len(rmgr.afterRenderObjects)
}

// The RenderManager that should be used for everything
var RenderMgr RenderManager
