package sse

type client struct {
	uid                 uint64
	notificationChannel NotifierChan
}

func NewClient(uid uint64) *client {
	return &client{
		uid:                 uid,
		notificationChannel: make(NotifierChan),
	}
}
