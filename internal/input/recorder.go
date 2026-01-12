package input

import (
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type RecorderService struct {
	isRecording bool
	stopChan    chan bool
}

func NewRecorderService() *RecorderService {
	return &RecorderService{
		isRecording: false,
		stopChan:    make(chan bool),
	}
}

type Coordinates struct {
	X int
	Y int
}

func (r *RecorderService) StartRecording() (int, int) {
	done := make(chan Coordinates)
	evChan := hook.Start()

	go func() {
		for ev := range evChan {
			if ev.Kind == hook.MouseDown && ev.Button == 1 {
				x, y := robotgo.Location()
				done <- Coordinates{X: x, Y: y}
				hook.End()
				break
			}
		}
	}()

	result := <-done
	close(done)

	return result.X, result.Y
}
