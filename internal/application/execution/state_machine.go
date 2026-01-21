package execution

import (
	"errors"
	"quant-trading/internal/domain/execution"
)

func (o *Order) Apply(evt *execution.Event) error {
	switch evt.Type {
	case execution.OrderSubmitted:
		return o.onSubmitted(evt)

	case execution.OrderAccepted:
		return o.onAccepted(evt)

	case execution.OrderPartiallyFilled:
		return o.onPartiallyFilled(evt)

	case execution.OrderFilled:
		return o.onFilled(evt)

	case execution.OrderCanceled:
		return o.onCanceled(evt)

	case execution.OrderRejected:
		return o.onRejected(evt)

	default:
		return errors.New("unknow event")
	}
}

func (o *Order) onSubmitted(evt *execution.Event) error {
	if o.State != StateNew {
		return errors.New("invalid state transition")
	}
	o.State = StateSubmitted
	return nil
}

func (o *Order) onAccepted(evt *execution.Event) error {
	if o.State != StateSubmitted {
		return errors.New("invalid state transition")
	}
	o.State = StateAccepted
	return nil
}

func (o *Order) onPartiallyFilled(evt *execution.Event) error {
	if o.State != StateAccepted && o.State != StatePartiallyFilled {
		return errors.New("invalid state transition")
	}

	o.FilledQty += evt.FilledQty

	if o.FilledQty >= o.Qty {
		o.State = StateFilled
	} else {
		o.State = StatePartiallyFilled
	}
	return nil
}

func (o *Order) onFilled(evt *execution.Event) error {
	if o.State != StateAccepted && o.State != StatePartiallyFilled {
		return errors.New("invalid state transition")
	}
	o.FilledQty = o.Qty
	o.State = StateFilled
	return nil
}

func (o *Order) onCanceled(evt *execution.Event) error {
	if o.State == StateFilled || o.State == StateRejected {
		return errors.New("cannot cancel terminal order")
	}
	o.State = StateCanceled
	return nil
}

func (o *Order) onRejected(evt *execution.Event) error {
	if o.State != StateSubmitted {
		return errors.New("invalid reject state")
	}
	o.State = StateRejected
	return nil
}
