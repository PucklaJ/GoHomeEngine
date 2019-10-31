package gohome

import (
	"strconv"
)

// An object that gets updated in every frame
type UpdateObject interface {
	// Gets called every frame with the time elapsed from the last frame to the current in seconds
	Update(delta_time float32)
}

// This manager handles the updating of all objects
type UpdateManager struct {
	updateObjects []UpdateObject
	breakLoop     bool
}

func (upmgr *UpdateManager) Init() {
}

// Adds an object to the loop
func (upmgr *UpdateManager) AddObject(upobj UpdateObject) {
	upmgr.updateObjects = append(upmgr.updateObjects, upobj)
}

// Removes an object from the loop
func (upmgr *UpdateManager) RemoveObject(upobj UpdateObject) {
	for i := 0; i < len(upmgr.updateObjects); i++ {
		if upmgr.updateObjects[i] == upobj {
			upmgr.updateObjects = append(upmgr.updateObjects[:i], upmgr.updateObjects[i+1:]...)
			return
		}
	}
}

// Tells the manager to break out of the update loop
func (upmgr *UpdateManager) BreakUpdateLoop() {
	upmgr.breakLoop = true
}

// Gets called every frame from the framework
func (upmgr *UpdateManager) Update(delta_time float32) {
	upmgr.breakLoop = false
	var obj UpdateObject
	objlen := len(upmgr.updateObjects)
	for i := 0; i < objlen && i < len(upmgr.updateObjects); i++ {
		if upmgr.breakLoop {
			upmgr.breakLoop = false
			break
		}
		obj = upmgr.updateObjects[i]
		if obj == nil {
			ErrorMgr.Error("UpdateManager", strconv.Itoa(i), "UpdateObject is nil")
		} else {
			obj.Update(delta_time)

			if i >= len(upmgr.updateObjects) || obj != upmgr.updateObjects[i] {
				i--
			}
		}
	}
}

// Removes all update objects
func (upmgr *UpdateManager) Terminate() {
	if len(upmgr.updateObjects) == 0 {
		return
	}

	upmgr.updateObjects = append(upmgr.updateObjects[:0], upmgr.updateObjects[len(upmgr.updateObjects):]...)
}

// Returns the number of currently attached update objects
func (upmgr *UpdateManager) NumUpdateObjects() int {
	return len(upmgr.updateObjects)
}

// The UpdateManager that should be used for everything
var UpdateMgr UpdateManager
