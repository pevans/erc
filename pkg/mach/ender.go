package mach

// An Ender is a type which allows us to shut down an emulator properly.
// Examples of things that may need to be done here are writing changes
// to files, wishing the user goodbye in some form, etc.
type Ender interface {
	End() error
}
