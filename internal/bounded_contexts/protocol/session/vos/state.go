package vos

type State struct{ value string }

func NewState(v string) (State, error) { return State{value: v}, nil }
func (s State) Value() string          { return s.value }
func (s State) String() string         { return s.value }
