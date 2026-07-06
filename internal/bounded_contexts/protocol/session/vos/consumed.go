package vos

type Consumed bool

func NewConsumed(v bool) Consumed    { return Consumed(v) }
func (c *Consumed) IsConsumed() bool { return bool(*c) }
