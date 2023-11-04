package gohome

import (
	"image/png"
	"os"

	"github.com/PucklaJ/mathgl/mgl32"
)

// A struct holding all data of the mouse
type Mouse struct {
	// The current position of the mouse in screen coordinates
	Pos [2]int16
	// The relative mouse movement to the last frame
	DPos [2]int16
	// The wheel movement values [Horizontal,Vertical]
	Wheel [2]int8
}

// Converts the mouse screen coordinates to 2D world coordinates
func (this *Mouse) ToWorldPosition2D() mgl32.Vec2 {
	return this.ToWorldPosition2DAdv(0, 0)
}

// Same as ToWorldPosition2D with additional arguments for the camera and the viewport
func (this *Mouse) ToWorldPosition2DAdv(cameraIndex int, viewportIndex int) mgl32.Vec2 {
	screenPos := mgl32.Vec2{float32(this.Pos[0]), float32(this.Pos[1])}
	if RenderMgr.EnableBackBuffer {
		wsize := Framew.WindowGetSize()
		screenPos = screenPos.DivVec(wsize)
		ns := Render.GetNativeResolution()
		screenPos = screenPos.MulVec(ns)
	}
	viewportPos := screenPos
	viewport := Render.GetViewport()
	if len(RenderMgr.viewport2Ds)-1 >= int(viewportIndex) {
		viewport = *RenderMgr.viewport2Ds[viewportIndex]
	}

	viewportPos = viewportPos.Sub(mgl32.Vec2{float32(viewport.X), float32(viewport.Y)})
	normalizedPos := mgl32.Vec2{viewportPos.X()/float32(viewport.Width)*2.0 - 1.0, (1.0-viewportPos.Y()/float32(viewport.Height))*2.0 - 1.0}
	projection := RenderMgr.Projection2D
	projMatrix := mgl32.Ident4()
	if projection != nil {
		projection.CalculateProjectionMatrix()
		projMatrix = projection.GetProjectionMatrix()
	}
	normalizedPosV3 := normalizedPos.Vec3(-1.0)
	projectedPos := projMatrix.Inv().Mul4x1(normalizedPosV3.Vec4(1)).Vec3()

	if len(RenderMgr.camera2Ds)-1 >= int(cameraIndex) {
		cam := RenderMgr.camera2Ds[cameraIndex]
		cam.CalculateViewMatrix()
		invViewMatrix := cam.GetInverseViewMatrix()
		projectedPos = invViewMatrix.Mat4().Mul4x1(projectedPos.Vec4(1)).Vec3()
	}

	return projectedPos.Vec2()
}

// Converts the raw mouse coordinates to coordinates adapting to the native resolution
func (this *Mouse) ToScreenPosition() (vec mgl32.Vec2) {
	vec[0], vec[1] = float32(this.Pos[0]), float32(this.Pos[1])

	ns := Render.GetNativeResolution()
	ws := Framew.WindowGetSize()
	rel := ns.DivVec(ws)
	vec = vec.MulVec(rel)

	return
}

// Converts the mouse coordinates to a 3D Ray pointing out of the camera
func (this *Mouse) ToRay() mgl32.Vec3 {
	return this.ToRayAdv(0, 0)
}

// Same as ToRay with additional arguments for the camera and the viewport
func (this *Mouse) ToRayAdv(viewportIndex, cameraIndex int32) mgl32.Vec3 {
	return ScreenPositionToRayAdv(this.ToScreenPosition(), viewportIndex, cameraIndex)
}

// A struct holding all the information of a touch on a touchscreen
type Touch struct {
	ID   uint8
	Pos  [2]int16
	DPos [2]int16
	PPos [2]int16
}

// The struct that handles every Input
type InputManager struct {
	prevKeys       map[Key]bool
	currentKeys    map[Key]bool
	holdTimes      map[Key]float32
	prevTouches    map[uint8]bool
	currentTouches map[uint8]bool

	// The data of the mouse
	Mouse Mouse
	// All registered touches
	Touches map[uint8]Touch
}

// Initialises all members of InputManager
func (inmgr *InputManager) Init() {
	inmgr.prevKeys = make(map[Key]bool)
	inmgr.currentKeys = make(map[Key]bool)
	inmgr.holdTimes = make(map[Key]float32)
	inmgr.currentTouches = make(map[uint8]bool)
	inmgr.prevTouches = make(map[uint8]bool)
	inmgr.Touches = make(map[uint8]Touch)
}

// Says that key has been pressed
func (inmgr *InputManager) PressKey(key Key) {
	inmgr.currentKeys[key] = true
	if _, ok := inmgr.holdTimes[key]; !ok {
		inmgr.holdTimes[key] = 0.0
	}
}

// Says that key has been released
func (inmgr *InputManager) ReleaseKey(key Key) {
	inmgr.currentKeys[key] = false
	inmgr.holdTimes[key] = 0.0
}

// Says that the touch with id has been touched
func (inmgr *InputManager) Touch(id uint8) {
	inmgr.currentTouches[id] = true
}

// Says that the touch with id has been released
func (inmgr *InputManager) ReleaseTouch(id uint8) {
	inmgr.currentTouches[id] = false
}

// Returns wether key is currently pressed
func (inmgr *InputManager) IsPressed(key Key) bool {
	if v, ok := inmgr.currentKeys[key]; ok && v {
		return true
	} else {
		return false
	}
}

// Returns wether key was pressed in the last frame
func (inmgr *InputManager) WasPressed(key Key) bool {
	if v, ok := inmgr.prevKeys[key]; ok && v {
		return true
	} else {
		return false
	}
}

// Returns wether key has been just pressed in this frame
func (inmgr *InputManager) JustPressed(key Key) bool {
	return inmgr.IsPressed(key) && !inmgr.WasPressed(key)
}

// Returns wether the touch with id is currently touched
func (inmgr *InputManager) IsTouched(id uint8) bool {
	touched, ok := inmgr.currentTouches[id]
	return ok && touched
}

// Returns wether the touched with id was touched in the last frame
func (inmgr *InputManager) WasTouched(id uint8) bool {
	touched, ok := inmgr.prevTouches[id]
	return ok && touched
}

// Returns wether the touch with id has been touched in this frame
func (inmgr *InputManager) JustTouched(id uint8) bool {
	return inmgr.IsTouched(id) && !inmgr.WasTouched(id)
}

func (inmgr *InputManager) TimeHeld(key Key) float32 {
	if v, ok := inmgr.holdTimes[key]; ok {
		return v
	} else {
		return 0.0
	}
}

// Updates everything and tells the InputManager that the frame is over
func (inmgr *InputManager) Update(delta_time float32) {
	if inmgr.JustPressed(KeyF12) {
		img := TextureToImage(RenderMgr.GetBackBuffer(), false, true)
		if img != nil {
			file, err := os.Create("screenshot.png")
			if err != nil {
				ErrorMgr.Error("Screenshot", "Failed", err.Error())
			} else {
				err = png.Encode(file, img)
				if err != nil {
					ErrorMgr.Error("Screenshot", "Failed", err.Error())
				}
			}
		}
	}

	for k, v := range inmgr.currentKeys {
		inmgr.prevKeys[k] = v
	}
	for k := range inmgr.holdTimes {
		if inmgr.IsPressed(k) {
			inmgr.holdTimes[k] += delta_time
		}
	}
	for k, v := range inmgr.currentTouches {
		inmgr.prevTouches[k] = v
	}
}

// The InputManager that should be used for everything
var InputMgr InputManager
