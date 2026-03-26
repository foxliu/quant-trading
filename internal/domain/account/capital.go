package account

type Capital struct {
	available float64
}

func newCapital(initial float64) *Capital {
	return &Capital{
		available: initial,
	}
}

func (c *Capital) applyFill(fill Fill) error {
	value := fill.Price * float64(fill.Qty)

	if fill.Side == Buy {
		return c.deduct(value + fill.Fee)
	}
	c.add(value - fill.Fee)
	return nil
}

func (c *Capital) add(amount float64) {
	c.available += amount
}

func (c *Capital) deduct(amount float64) error {
	if c.available < amount {
		return ErrInsufficientFunds
	}
	c.available -= amount
	return nil
}
