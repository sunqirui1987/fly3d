// Package event is a simple event system.
package goevent

type IDispatcher interface {
	On(name string, fn interface{}) error
	Emit(name string, params ...interface{}) error
	Has(name string) bool
	List() []string
	Remove(names ...string)
}

// Default event instance
var globalSource = New()

// On set new listener from the default source.
func On(name string, fn interface{}) error {
	return globalSource.On(name, fn)
}

// Emit firing an event from the default source.
func Emit(name string, params ...interface{}) error {
	return globalSource.Emit(name, params...)
}

// Has returns true if a event exists from the default source.
func Has(name string) bool {
	return globalSource.Has(name)
}

// List returns list events from the default source.
func List() []string {
	return globalSource.List()
}

// Remove delete events from the event list from the default source.
func Remove(names ...string) {
	globalSource.Remove(names...)
}
