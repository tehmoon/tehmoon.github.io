package main

import (
	"log"
)

type Module interface{
	Init() (error)
}

const (
// TYPES OMIT
	MessageTypeReset	MessageType = iota
	MessageTypeData
	MessageTypeClose
// TYPES OMIT
)

type MessageType int

type Message struct{
	Type MessageType
	Payload []byte
}

type GenerateText struct {
	in chan *Message
	out chan *Message
}

func (m *GenerateText) Init() (error) {
	go func(module *GenerateText) {
		words := []string{"Hello", "golang", "boston", "meetup",}

		for message := range module.in {
			switch message.Type {
				case MessageTypeReset:
					module.out <- &Message{Type: MessageTypeReset,}

					for _, word := range words {
						module.out <- &Message{Type: MessageTypeData, Payload: []byte(word),}
					}

					module.out <- &Message{Type: MessageTypeClose,}
				default:
					module.out <- message
			}
		}

		close(module.out)
	}(m)

	return nil
}

func NewGenerateText(in, out chan *Message) (module *GenerateText) {
	return &GenerateText{
		in: in,
		out: out,
	}
}

type Stdout struct {
	in chan *Message
	out chan *Message
}

func (m *Stdout) Init() (error) {
	go func(module *Stdout) {
		for message := range module.in {
			switch message.Type {
				case MessageTypeData:
					log.Println(string(message.Payload[:]))
				default:
					module.out <- message
			}
		}

		close(module.out)
	}(m)

	return nil
}

func NewStdout(in, out chan *Message) (module *Stdout) {
	return &Stdout{
		in: in,
		out: out,
	}
}

func main() {
	in := make(chan *Message, 5)
	out := make(chan *Message, 5)

	modc := make(chan *Message, 5)

	mod1 := NewGenerateText(in, modc)
	err := mod1.Init()
	if err != nil {
		panic(err)
	}

	mod2 := NewStdout(modc, out)
	err = mod2.Init()
	if err != nil {
		panic(err)
	}

	in <- &Message{Type: MessageTypeReset,}

	isclosed := false
	LOOP: for message := range out {
		switch message.Type {
			case MessageTypeReset:
			case MessageTypeClose:
				if isclosed {
					close(in)
					break LOOP
				}

				isclosed = true
				in <- message
			default:
				in <- message
		}
	}

	if !isclosed {
		close(in)
	}
}
