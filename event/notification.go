package event

type Intent struct {
	Name string
	File string
}

var handler *Handler = nil

func StartListening() {
	handler = &Handler{
		Chan: make(chan *Intent, 128),
	}
    go handler.StartProcessing()
}

func Notify(name string, filename string) {
	handler.Chan <- &Intent{
		Name: name,
		File: filename,
	}
}
