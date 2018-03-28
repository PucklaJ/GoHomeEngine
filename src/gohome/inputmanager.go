package gohome

type Mouse struct {
	Pos   [2]int16
	DPos  [2]int16
	Wheel [2]int8
}

type InputManager struct {
	prevKeys    map[Key]bool
	currentKeys map[Key]bool
	holdTimes   map[Key]float32
	Mouse       Mouse
}

func (inmgr *InputManager) Init() {
	inmgr.prevKeys = make(map[Key]bool)
	inmgr.currentKeys = make(map[Key]bool)
	inmgr.holdTimes = make(map[Key]float32)
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
}

var InputMgr InputManager
