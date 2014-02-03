package event

type Intent struct {
	EventId  uint32
	FileName string
}

func Notify(id uint32, filename string) {
	handler.Chan <- &Intent{
		EventId:  id,
		FileName: filename,
	}
}
