package gohome

// import "fmt"

const (
	TEST_CAMERA_MOVEMENT_MOVE_SPEED           float32 = 30.0
	TEST_CAMERA_MOVEMENT_ROTATE_SPEED         float32 = 0.5
	TEST_CAMERA_MOVEMENT_MOVE_SPEED_MAGNIFIER float32 = 5.0
)

type TestCameraMovement3D struct {
	cam *Camera3D
}

func (this *TestCameraMovement3D) Init(cam *Camera3D) {
	this.cam = cam
	this.cam.Init()
}

func (this *TestCameraMovement3D) updateLookDirection() {
	pitch := float32(InputMgr.Mouse.DPos[1]) * TEST_CAMERA_MOVEMENT_ROTATE_SPEED
	yaw := float32(InputMgr.Mouse.DPos[0]) * TEST_CAMERA_MOVEMENT_ROTATE_SPEED

	this.cam.AddRotation([2]float32{-pitch, -yaw})
}

func (this *TestCameraMovement3D) updatePosition(delta_time float32) {
	var pos [3]float32
	var speed float32 = TEST_CAMERA_MOVEMENT_MOVE_SPEED

	if InputMgr.IsPressed(KeyLeftControl) {
		speed *= 1.0 / TEST_CAMERA_MOVEMENT_MOVE_SPEED_MAGNIFIER
	} else if InputMgr.IsPressed(KeyLeftShift) {
		speed *= TEST_CAMERA_MOVEMENT_MOVE_SPEED_MAGNIFIER
	}

	if InputMgr.IsPressed(KeyW) {
		pos[2] -= speed * delta_time
	}

	if InputMgr.IsPressed(KeyS) {
		pos[2] += speed * delta_time
	}

	if InputMgr.IsPressed(KeyA) {
		pos[0] -= speed * delta_time
	}

	if InputMgr.IsPressed(KeyD) {
		pos[0] += speed * delta_time
	}

	this.cam.AddPositionRelative(pos)

	if InputMgr.IsPressed(KeyI) {
		pos[1] += speed * delta_time
	}

	if InputMgr.IsPressed(KeyK) {
		pos[1] -= speed * delta_time
	}

	this.cam.Position = this.cam.Position.Add([3]float32{0.0, pos[1], 0.0})
}

func (this *TestCameraMovement3D) Update(delta_time float32) {
	if InputMgr.JustPressed(KeyM) {
		Framew.CursorDisable()
	} else if InputMgr.JustPressed(KeyEscape) {
		Framew.CurserShow()
	}

	if Framew.CursorDisabled() {
		this.updateLookDirection()
		this.updatePosition(delta_time)
	}
}
