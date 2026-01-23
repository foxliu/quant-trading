package event

type Handler func(evt *Envelope)

type Bus interface {
	Publish(evt *Envelope)
	Subscribe(t Type, h Handler)
}
