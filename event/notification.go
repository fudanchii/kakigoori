package event

type Intent struct {
	EventId  byte
	FileName string
}

func Notify(id byte, filename string) {
	handler.Chan <- &Intent{
		EventId:  id,
		FileName: filename,
	}
}
