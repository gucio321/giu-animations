package animations

type TriggerFunc func() bool

type TriggerType byte

const (
	// TriggerNever is a devault value.
	// Animation will not be started automatically.
	TriggerNever TriggerType = iota
	// TriggerOnTrue will start animation whenever trigger becomes true.
	TriggerOnTrue
	// TriggerOnChange will trigger animation, when value of trigger's function
	// will change.
	TriggerOnChange
)
