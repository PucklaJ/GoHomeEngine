package gohome

import (
	"strconv"
)

type UpdateObject interface {
	Update(delta_time float32)
}

type UpdateManager struct {
	updateObjects []UpdateObject
	breakLoop     bool
}

func (upmgr *UpdateManager) Init() {
}

func (upmgr *UpdateManager) AddObject(upobj UpdateObject) {
	upmgr.updateObjects = append(upmgr.updateObjects, upobj)
}

func (upmgr *UpdateManager) RemoveObject(upobj UpdateObject) {
	for i := 0; i < len(upmgr.updateObjects); i++ {
		if upmgr.updateObjects[i] == upobj {
			upmgr.updateObjects = append(upmgr.updateObjects[:i], upmgr.updateObjects[i+1:]...)
			return
		}
	}
}

func (upmgr *UpdateManager) BreakUpdateLoop() {
	upmgr.breakLoop = true
}

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

func (upmgr *UpdateManager) Terminate() {
	if len(upmgr.updateObjects) == 0 {
		return
	}

	upmgr.updateObjects = append(upmgr.updateObjects[:0], upmgr.updateObjects[len(upmgr.updateObjects):]...)
}

func (upmgr *UpdateManager) NumUpdateObjects() uint32 {
	return uint32(len(upmgr.updateObjects))
}

var UpdateMgr UpdateManager
