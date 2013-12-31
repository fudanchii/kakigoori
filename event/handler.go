package event

import (
	"log"
)

type Handler struct {
	Chan chan *Intent
}

func (h *Handler) StartProcessing() {
	for intent := range h.Chan {
		//maybe find some handler first before processing
		log.Println("> ", intent.Name, intent.File)
	}
}
