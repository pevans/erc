package mach

// A Loader is a type which can load data from some source and represent
// that as a kind of input to an emulated machine.
type Loader interface {
	Load(string) error
}
