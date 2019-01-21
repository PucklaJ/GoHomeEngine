package gohome

import (
	"time"
)

const (
	NUM_MEASUREMENTS = 11
	NUM_TAKEAWAYS    = 2
)

type FPSLimiter struct {
	MaxFPS    int
	FPS       int
	DeltaTime float32

	timeMeasure    time.Time
	deltaDuration  time.Duration
	additionalTime float32

	realDeltaTimes             [NUM_MEASUREMENTS]float64
	currentRealDeltaTimesIndex int
	numMeasuredRealDeltaTimes  int
	meanDeltaTime              float64
	smoothedDeltaTimes         [NUM_MEASUREMENTS]float64
}

/*
1. Get eleven last delta times
2. Put away highest two and lowest two
3. Calculate mean
*/

func (fps *FPSLimiter) Init() {
	fps.MaxFPS = 60
	fps.FPS = 0
	fps.DeltaTime = 0.0
	fps.currentRealDeltaTimesIndex = 0
	fps.numMeasuredRealDeltaTimes = 0
	fps.additionalTime = 0.0
}

func (fps *FPSLimiter) StartMeasurement() {
	fps.timeMeasure = time.Now()
}

func (fps *FPSLimiter) EndMeasurement() {
	fps.deltaDuration = time.Now().Sub(fps.timeMeasure)
	fps.timeMeasure = time.Now()
	fps.saveCurrentDeltaTime()
}

func (fps *FPSLimiter) AddTime(amount float32) {
	fps.additionalTime += amount
}

func (fps *FPSLimiter) currentFrameIndex() (curIndex int) {
	if fps.currentRealDeltaTimesIndex == 0 {
		curIndex = NUM_MEASUREMENTS - 1
	} else {
		curIndex = fps.currentRealDeltaTimesIndex - 1
	}

	return
}

func (fps *FPSLimiter) LimitFPS() {
	currentDeltaTime := fps.realDeltaTimes[fps.currentFrameIndex()]
	if currentDeltaTime != 0.0 {
		maxDeltaTime := float64(1.0) / float64(fps.MaxFPS)

		if currentDeltaTime < maxDeltaTime {
			time.Sleep(time.Microsecond * time.Duration((maxDeltaTime-currentDeltaTime)*1000.0*1000.0))
		}
	}
	fps.deltaDuration = time.Now().Sub(fps.timeMeasure)
	fps.smoothDeltaTimes()
}

func (fps *FPSLimiter) saveCurrentDeltaTime() {
	fps.realDeltaTimes[fps.currentRealDeltaTimesIndex] = fps.deltaDuration.Seconds() + float64(fps.additionalTime)
	fps.additionalTime = 0.0
	if fps.numMeasuredRealDeltaTimes < NUM_MEASUREMENTS {
		fps.numMeasuredRealDeltaTimes++
	}
	if fps.currentRealDeltaTimesIndex++; fps.currentRealDeltaTimesIndex == NUM_MEASUREMENTS {
		fps.currentRealDeltaTimesIndex = 0
	}
}

func (fps *FPSLimiter) smoothDeltaTimes() {
	currentDeltaTime := fps.realDeltaTimes[fps.currentFrameIndex()]
	fps.smoothedDeltaTimes[fps.currentFrameIndex()] = fps.deltaDuration.Seconds() + currentDeltaTime
	if fps.numMeasuredRealDeltaTimes > 2*NUM_TAKEAWAYS {
		fps.takeAwayHighestAndLowest()
	}

	fps.calculateMeanDeltaTime()
	fps.DeltaTime = float32(fps.meanDeltaTime)
}

func swap(f1 *float64, f2 *float64) {
	temp := *f1
	*f1 = *f2
	*f2 = temp
}

func (fps *FPSLimiter) takeAwayHighestAndLowest() {
	var maxVal float64
	var minVal float64
	var maxIndex int
	var minIndex int

	for i := 0; i < NUM_TAKEAWAYS; i++ {
		maxVal = fps.smoothedDeltaTimes[fps.numMeasuredRealDeltaTimes-1-i]
		minVal = fps.smoothedDeltaTimes[i]
		maxIndex = fps.numMeasuredRealDeltaTimes - 1 - i
		minIndex = i
		for j := i; j < fps.numMeasuredRealDeltaTimes-i; j++ {
			if fps.smoothedDeltaTimes[j] < minVal {
				minVal = fps.smoothedDeltaTimes[j]
				minIndex = j
			}
			if fps.smoothedDeltaTimes[j] > maxVal {
				maxVal = fps.smoothedDeltaTimes[j]
				maxIndex = j
			}
		}

		swap(&fps.smoothedDeltaTimes[minIndex], &fps.smoothedDeltaTimes[i])
		swap(&fps.smoothedDeltaTimes[maxIndex], &fps.smoothedDeltaTimes[fps.numMeasuredRealDeltaTimes-1-i])
	}
}

func (fps *FPSLimiter) calculateMeanDeltaTime() {
	fps.meanDeltaTime = 0.0
	if fps.numMeasuredRealDeltaTimes > 2*NUM_TAKEAWAYS {
		for i := NUM_TAKEAWAYS; i < fps.numMeasuredRealDeltaTimes-NUM_TAKEAWAYS; i++ {
			fps.meanDeltaTime += fps.smoothedDeltaTimes[i]
		}

		fps.meanDeltaTime /= float64(fps.numMeasuredRealDeltaTimes - 2*NUM_TAKEAWAYS)
	} else {
		for i := 0; i < fps.numMeasuredRealDeltaTimes; i++ {
			fps.meanDeltaTime += fps.smoothedDeltaTimes[i]
		}

		fps.meanDeltaTime /= float64(fps.numMeasuredRealDeltaTimes)
	}
}

var FPSLimit FPSLimiter
