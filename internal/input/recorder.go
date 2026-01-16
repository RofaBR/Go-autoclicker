package input

type RecorderService struct {
}

func NewRecorderService() *RecorderService {
	return &RecorderService{}
}

type Coordinates struct {
	X int
	Y int
}

func (r *RecorderService) StartRecording() (int, int) {
	return GetHookManager().StartRecording()
}
