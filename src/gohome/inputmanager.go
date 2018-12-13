package gohome

import (
	"github.com/PucklaMotzer09/mathgl/mgl32"
	"image/png"
	"os"
)

type Mouse struct {
	Pos   [2]int16
	DPos  [2]int16
	Wheel [2]int8
}

func (this *Mouse) ToWorldPosition2D() mgl32.Vec2 {
	return this.ToWorldPosition2DAdv(0, 0)
}

func (this *Mouse) ToWorldPosition2DAdv(cameraIndex int32, viewportIndex uint32) mgl32.Vec2 {
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

func (this *Mouse) ToScreenPosition() (vec mgl32.Vec2) {
	vec[0], vec[1] = float32(this.Pos[0]), float32(this.Pos[1])

	ns := Render.GetNativeResolution()
	ws := Framew.WindowGetSize()
	rel := ns.DivVec(ws)
	vec = vec.MulVec(rel)

	return
}

func (this *Mouse) ToRay() mgl32.Vec3 {
	return this.ToRayAdv(0, 0)
}

func (this *Mouse) ToRayAdv(viewportIndex, cameraIndex int32) mgl32.Vec3 {
	return ScreenPositionToRayAdv(this.ToScreenPosition(), viewportIndex, cameraIndex)
}

type Touch struct {
	ID   uint8
	Pos  [2]int16
	DPos [2]int16
	PPos [2]int16
}

type InputManager struct {
	prevKeys       map[Key]bool
	currentKeys    map[Key]bool
	holdTimes      map[Key]float32
	prevTouches    map[uint8]bool
	currentTouches map[uint8]bool

	Mouse   Mouse
	Touches map[uint8]Touch
}

func (inmgr *InputManager) Init() {
	inmgr.prevKeys = make(map[Key]bool)
	inmgr.currentKeys = make(map[Key]bool)
	inmgr.holdTimes = make(map[Key]float32)
	inmgr.currentTouches = make(map[uint8]bool)
	inmgr.prevTouches = make(map[uint8]bool)
	inmgr.Touches = make(map[uint8]Touch)
}

func (inmgr *InputManager) PressKey(key Key) {
	inmgr.currentKeys[key] = true
	if _, ok := inmgr.holdTimes[key]; !ok {
		inmgr.holdTimes[key] = 0.0
	}
}

func (inmgr *InputManager) ReleaseKey(key Key) {
	inmgr.currentKeys[key] = false
	inmgr.holdTimes[key] = 0.0
}

func (inmgr *InputManager) Touch(id uint8) {
	inmgr.currentTouches[id] = true
}

func (inmgr *InputManager) ReleaseTouch(id uint8) {
	inmgr.currentTouches[id] = false
}

func (inmgr *InputManager) IsPressed(key Key) bool {
	if v, ok := inmgr.currentKeys[key]; ok && v {
		return true
	} else {
		return false
	}
}

func (inmgr *InputManager) WasPressed(key Key) bool {
	if v, ok := inmgr.prevKeys[key]; ok && v {
		return true
	} else {
		return false
	}
}

func (inmgr *InputManager) JustPressed(key Key) bool {
	return inmgr.IsPressed(key) && !inmgr.WasPressed(key)
}

func (inmgr *InputManager) IsTouched(id uint8) bool {
	touched, ok := inmgr.currentTouches[id]
	return ok && touched
}

func (inmgr *InputManager) WasTouched(id uint8) bool {
	touched, ok := inmgr.prevTouches[id]
	return ok && touched
}

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

func (inmgr *InputManager) Update(delta_time float32) {
	if inmgr.JustPressed(KeyF12) {
		img := TextureToImage(RenderMgr.GetBackBuffer(), false, true)
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

var InputMgr InputManager
