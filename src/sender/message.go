package sender

import "github.com/google/uuid"

func NewMessage[V string | []byte](text V) Message {
	return Message{
		id:   uuid.New(),
		text: []byte(text),
	}
}

type Message struct {
	id   uuid.UUID
	text []byte
}

func (m Message) GetID() uuid.UUID {
	return m.id
}

func (m Message) GetText() []byte {
	return m.text
}
