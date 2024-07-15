package mos

import "errors"

// Decimal is a representation of a binary-coded decimal. It is used to
// perform mathematics in a base-10 fashion and to later return the
// correct binary form of the result.
type Decimal struct {
	Result   int
	Carry    bool
	Negative bool
	Error    error
}

var (
	ErrInvalid error = errors.New("invalid bcd format")
)

// NewDecimal returns a new Decimal type with integer as the starting
// result.
func NewDecimal(integer int) Decimal {
	lsd := integer & 0xF
	msd := (integer >> 4) & 0xF

	if lsd > 9 || msd > 9 {
		return Decimal{
			Error: ErrInvalid,
		}
	}

	return Decimal{
		Result: (msd * 10) + lsd,
	}
}

// Add will add any number of terms to a given decimal.
func (d *Decimal) Add(terms ...Decimal) {
	for _, term := range terms {
		d.Result += term.Result
	}

	if d.Result > 99 {
		d.Result -= 100
		d.Carry = true
	}
}

// Subtract will subtract any number of terms from a given decimal.
func (d *Decimal) Subtract(terms ...Decimal) {
	for _, term := range terms {
		d.Result -= term.Result
	}

	if d.Result < 0 {
		d.Result += 100
		d.Negative = true
	}
}

// Binary returns the binary coded decimal that would be expected within
// any value in a register or memory.
func (d *Decimal) Binary() int {
	result := d.Result % 100

	low4 := result % 10
	high4 := result / 10

	return (high4 << 4) | low4
}
