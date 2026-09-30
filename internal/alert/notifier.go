package alert

// Notifier interface, delivers an alert about an Event to somewhere

type Notifier interface {
	Notify(e Event, count int) error
}
