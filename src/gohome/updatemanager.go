package gohome

type UpdateObject interface {
	Update(delta_time float32)
}

type UpdateManager struct {
	updateObjects []UpdateObject
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

func (upmgr *UpdateManager) Update(delta_time float32) {
	for i := 0; i < len(upmgr.updateObjects); i++ {
		upmgr.updateObjects[i].Update(delta_time)
	}
}

func (upmgr *UpdateManager) Terminate() {
	if len(upmgr.updateObjects) == 0 {
		return
	}

	upmgr.updateObjects = append(upmgr.updateObjects[:0], upmgr.updateObjects[len(upmgr.updateObjects):]...)
}

var UpdateMgr UpdateManager
