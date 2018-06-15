package gohome

type Tweenset struct {
	Tweens []Tween
	Loop bool

	currentTweens []Tween
	addedTweensNotAdd []uint32
	tweenBeforeNextTween Tween
	currentStoppedTween uint32
	paused bool
	allTweensAdded bool
	parent interface{}
}

func (this *Tweenset) SetParent(twobj interface{}) {
	this.parent = twobj
}

func (this *Tweenset) Pause() {
	this.paused = true
}

func (this *Tweenset) Resume() {
	this.paused = false
}

func (this *Tweenset) Stop() {
	this.Reset()
}

func (this *Tweenset) Start() {
	this.Reset()
	this.paused = false
}

func (this *Tweenset) Update(delta_time float32) {
	if this.paused {
		return
	}

	if this.shouldAddNewTweens() {
		this.checkTweensToStart()
	}

	for i:=0;i<len(this.currentTweens);i++ {
		if this.currentTweens[i].Update(delta_time) {
			this.currentTweens[i].End()
			if i+1 == len(this.currentTweens) {
				this.currentTweens = this.currentTweens[:i]
			} else {
				this.currentTweens = append(this.currentTweens[:i],this.currentTweens[i+1:]...)
				i--
			}
		}
	}
}

func (this *Tweenset) shouldAddNewTweens() bool {
	if len(this.currentTweens) == 0 && !this.allTweensAdded {
		return true
	} else if len(this.currentTweens) == 0 && this.allTweensAdded {
		if this.Loop {
			this.Reset()
			this.paused = false
			return true
		}
		return false
	}

	var allAlways bool
	var tweenBeforeNextEnded bool = this.tweenBeforeNextTween != nil

	for i:=0;i<len(this.currentTweens);i++ {
		if this.currentTweens[i].GetType() != TWEEN_TYPE_ALWAYS {
			allAlways = false
			if this.tweenBeforeNextTween != nil && this.tweenBeforeNextTween == this.currentTweens[i] {
				tweenBeforeNextEnded = false
			}
		}
	}

	return allAlways || tweenBeforeNextEnded
}

func (this *Tweenset) Reset() {
	if len(this.currentTweens) != 0 {
		this.currentTweens = this.currentTweens[:0]
	}
	if len(this.addedTweensNotAdd) != 0 {
		this.addedTweensNotAdd = this.addedTweensNotAdd[:0]
	}
	this.tweenBeforeNextTween = nil
	this.currentStoppedTween = 0
	this.paused = true
	this.allTweensAdded = false

	for i:=len(this.Tweens)-1;i>=0;i-- {
		this.Tweens[i].Reset()
	}
}

func (this *Tweenset) checkTweensToStart() {
	afterPreviousFound := false
	for i:=this.currentStoppedTween;i<uint32(len(this.Tweens));i++ {
		if i == this.currentStoppedTween && this.Tweens[i].GetType() == TWEEN_TYPE_AFTER_PREVIOUS {
			this.addTweenToCurrent(i)
			this.tweenBeforeNextTween = nil
		} else {
			if !afterPreviousFound {
				if this.Tweens[i].GetType() != TWEEN_TYPE_AFTER_PREVIOUS {
					this.addTweenToCurrent(i)
				} else {
					startThisTween := false
					notAlwaysFound := int(i)-1
					if this.Tweens[i-1].GetType() == TWEEN_TYPE_ALWAYS {
						notAlwaysFound = this.searchNextTweenNotAlways(int(i)-2,this.currentStoppedTween)
						if notAlwaysFound > 0 {
							startThisTween = false
						} else {
							startThisTween = true
						}
					}

					if !startThisTween {
						afterPreviousFound = true
						this.tweenBeforeNextTween = this.Tweens[notAlwaysFound]
						this.currentStoppedTween = i
					} else {
						this.addTweenToCurrent(i)
					}
				}
			} else {
				if this.Tweens[i].GetType() == TWEEN_TYPE_ALWAYS {
					this.addTweenToCurrent(i)
				}
			}
		}
	}

	if !afterPreviousFound {
		this.allTweensAdded = true
		this.currentStoppedTween = uint32(len(this.Tweens))
	}
}

func (this *Tweenset) searchNextTweenNotAlways(start int,end uint32) int {
	for j:=start;j>=int(end);j-- {
		if this.Tweens[j].GetType() != TWEEN_TYPE_ALWAYS {
			return j
		}
	}

	return -1
}

func (this *Tweenset) addTweenToCurrent(i uint32) {
	if !this.hasAlreadyBeenAdded(i) {
		this.currentTweens = append(this.currentTweens,this.Tweens[i])
		this.addedTweensNotAdd = append(this.addedTweensNotAdd,i)
		this.Tweens[i].Start(this.parent)
	}
}

func (this *Tweenset) hasAlreadyBeenAdded(i uint32) bool {
	for j:=0;j<len(this.addedTweensNotAdd);j++ {
		if this.addedTweensNotAdd[j] == i {
			return true
		}
	}
	return false
}
