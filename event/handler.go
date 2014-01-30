package event

import (
	"log"
)

const (
	Unknown = iota

	Open
	Read
	Write
	Fsync
	Close

	Create
	Rename

	Unlink
	Link
	Symlink
	Readlink

	Chmod
	Chown
	Trunc

	OpenDir
	Mkdir
	Rmdir

	Mknod
	Fallocate
	Access
)

type HandlerFunction func(*Intent, Config)

type EventMap map[string]HandlerFunction

type Handler struct {
	Chan chan *Intent
	Config
	EventMap
}

var EventName [21]string = [21]string{
	"unknown",
	"open", "read", "write", "fsync", "close",
	"create", "rename",
	"unlink", "link", "symlink", "Readlink",
	"chmod", "chown", "trunc",
	"opendir", "mkdir", "rmdir",
	"mknod", "fallocate", "access",
}

var handler *Handler

func StartListening(config Config) *Handler {
	handler = &Handler{
		Chan:     make(chan *Intent, 128),
		EventMap: make(EventMap),
		Config:   config,
	}
	go handler.StartProcessing()
	return handler
}

func (h *Handler) RegisterHandler(events []byte, fn HandlerFunction) {
	h.EventMap[string(events)] = fn
}

func (h *Handler) StartProcessing() {
	for intent := range h.Chan {
		if handler, ok := h.EventMap[string([]byte{intent.EventId})]; ok {
			handler(intent, h.Config)
		}
		log.Printf("> %s\t%s\n", EventName[intent.EventId], intent.FileName)
	}
}
