package event

type Intent struct {
	EventId  byte
	FileName string
}

var handler *Handler = nil

func StartListening() {
	handler = &Handler{
		Chan: make(chan *Intent, 128),
	}
	go handler.StartProcessing()
}

func Notify(id byte, filename string) {
	handler.Chan <- &Intent{
		EventId:  id,
		FileName: filename,
	}
}
