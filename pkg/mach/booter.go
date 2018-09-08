package mach

// A Booter is a type which allows for a "boot" procedure
type Booter interface {
	Boot() error
}
