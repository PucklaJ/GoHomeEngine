package gohome

type Mouse struct {
	Pos   [2]int16
	DPos  [2]int16
	Wheel [2]int8
}

type Touch struct {
	ID 		uint8
	Pos 	[2]int16
	DPos 	[2]int16
	PPos 	[2]int16
}

type InputManager struct {
	prevKeys    map[Key]bool
	currentKeys map[Key]bool
	holdTimes   map[Key]float32
	prevTouches map[uint8]bool
	currentTouches map[uint8]bool

	Mouse       Mouse
	Touches		map[uint8]Touch
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
	touched,ok := inmgr.currentTouches[id]
	return touched && ok
}

func (inmgr *InputManager) WasTouched(id uint8) bool {
	touched,ok := inmgr.prevTouches[id]
	return touched && ok
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
	for k, v := range inmgr.currentKeys {
		inmgr.prevKeys[k] = v
	}
	for k := range inmgr.holdTimes {
		if inmgr.IsPressed(k) {
			inmgr.holdTimes[k] += delta_time
		}
	}
	for k,v := range inmgr.currentTouches {
		inmgr.prevTouches[k] = v
	}
}

var InputMgr InputManager
