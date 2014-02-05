package event

import (
	"log"
)

const (
	Unknown = iota

	Open  = 1
	Read  = 1 << 1
	Write = 1 << 2
	Fsync = 1 << 3
	Close = 1 << 4

	Create = 1 << 5
	Rename = 1 << 6

	Unlink   = 1 << 7
	Link     = 1 << 8
	Symlink  = 1 << 9
	Readlink = 1 << 10

	Chmod = 1 << 11
	Chown = 1 << 12
	Trunc = 1 << 13

	OpenDir = 1 << 14
	Mkdir   = 1 << 15
	Rmdir   = 1 << 16

	Mknod     = 1 << 17
	Fallocate = 1 << 18
	Access    = 1 << 19
)

type HandlerFunction func(*Intent, Config)

type EventMap map[uint32]HandlerFunction

type Handler struct {
	Chan chan *Intent
	Config
	EventMap
	TrackedEvents map[string]uint32
}

var EventName = map[uint32]string{
	Unknown: "unknown",

	Open:  "open",
	Read:  "read",
	Write: "write",
	Fsync: "fsync",
	Close: "close",

	Create: "create",
	Rename: "rename",

	Unlink:   "unlink",
	Link:     "link",
	Symlink:  "symlink",
	Readlink: "Readlink",

	Chmod: "chmod",
	Chown: "chown",
	Trunc: "trunc",

	OpenDir: "opendir",
	Mkdir:   "mkdir",
	Rmdir:   "rmdir",

	Mknod:     "mknod",
	Fallocate: "fallocate",
	Access:    "access",
}

var handler *Handler

func StartListening(config Config) *Handler {
	handler = &Handler{
		Chan:          make(chan *Intent, 512),
		Config:        config,
		EventMap:      make(EventMap),
		TrackedEvents: make(map[string]uint32),
	}
	go handler.StartProcessing()
	return handler
}

func (h *Handler) RegisterHandler(events uint32, fn HandlerFunction) {
	if events == 0 {
		return // Nothing
	}
	h.EventMap[events] = fn
}

func (h *Handler) StartProcessing() {
	var events uint32
	for intent := range h.Chan {
		// Check if this file already recorded
		if evseq, ok := h.TrackedEvents[intent.FileName]; ok {
			events = evseq | intent.EventId
		} else {
			events = intent.EventId
		}

		if handler, ok := h.EventMap[events]; ok {
			go handler(intent, h.Config)
			delete(h.TrackedEvents, intent.FileName)
		} else if events & (Close|Unlink|Rename|Rmdir) > 0 {
			delete(h.TrackedEvents, intent.FileName)
		} else {
			for key, _ := range h.EventMap {
				if events & key > 0 {
					h.TrackedEvents[intent.FileName] = events
					log.Printf("%d registered for %s\n", events, intent.FileName)
					break
				}
			}
		}
		log.Printf("> %s\t%s\n", EventName[intent.EventId], intent.FileName)
	}
}
