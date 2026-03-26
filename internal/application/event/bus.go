package event

// Handler 事件处理器（接口定义的类型）
type Handler func(evt *Envelope)

// Bus 事件总线接口
type Bus interface {
	Publish(evt *Envelope)
	Subscribe(t Type, h Handler)
}
