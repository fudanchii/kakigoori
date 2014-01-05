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

var EventName [21]string = [21]string{
	"unknown",
	"open", "read", "write", "fsync", "close",
	"create", "rename",
	"unlink", "link", "symlink", "Readlink",
	"chmod", "chown", "trunc",
	"opendir", "mkdir", "rmdir",
	"mknod", "fallocate", "access",
}

type handlerFunction func(*Intent)

type EventMap map[string]handlerFunction

func (evMap *EventMap) RegisterHandler(events []byte, fn handlerFunction) {
	(*evMap)[string(events)] = fn
}

type Handler struct {
	Chan chan *Intent
	EventMap
}

func (h *Handler) StartProcessing() {

	for intent := range h.Chan {
		//maybe find some handler first before processing
		log.Printf("> %s\t%s\n", EventName[intent.EventId], intent.FileName)
	}
}
