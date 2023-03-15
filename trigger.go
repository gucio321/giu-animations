package animations

// TriggerFunc is a function determining whether the animation should
// be auto-triggered. See also TriggerType.
type TriggerFunc func() bool

// TriggerType represents a strategy of automated triggering of an animation.
type TriggerType byte

const (
	// TriggerNever is animation default value.
	// Animation will not be started automatically.
	TriggerNever TriggerType = iota
	// TriggerOnTrue will start animation whenever trigger becomes true.
	TriggerOnTrue
	// TriggerOnChange will trigger animation, when value of trigger's function
	// will change.
	TriggerOnChange
)
