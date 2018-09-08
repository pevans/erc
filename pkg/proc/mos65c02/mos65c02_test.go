package mos65c02

import (
	"testing"

	"github.com/pevans/erc/pkg/mach"
	"github.com/stretchr/testify/assert"
)

func TestApplyStatus(t *testing.T) {
	c := CPU{}

	cases := []struct {
		b        bool
		in, want mach.Byte
	}{
		{true, 1, 1},
		{true, 2, 2},
		{false, 4, 0},
	}

	for _, cas := range cases {
		c.P = 0
		c.ApplyStatus(cas.b, cas.in)

		assert.Equal(t, cas.want, c.P)
	}
}

func TestApplyN(t *testing.T) {
	c := CPU{}

	cases := []struct {
		in, want mach.Byte
	}{
		{64, 0},
		{128, NEGATIVE},
	}

	for _, cas := range cases {
		c.ApplyN(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestApplyZ(t *testing.T) {
	c := CPU{}

	cases := []struct {
		in, want mach.Byte
	}{
		{0, ZERO},
		{1, 0},
	}

	for _, cas := range cases {
		c.ApplyZ(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestApplyNZ(t *testing.T) {
	c := CPU{}

	cases := []struct {
		in, want mach.Byte
	}{
		{0, ZERO},
		{1, 0},
		{128, NEGATIVE},
	}

	for _, cas := range cases {
		c.ApplyNZ(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestCompare(t *testing.T) {
	c := CPU{}

	cases := []struct {
		a, oper, pWant mach.Byte
	}{
		{1, 3, NEGATIVE | CARRY},
		{3, 1, CARRY},
		{3, 3, ZERO},
	}

	for _, cas := range cases {
		c.A = cas.a
		c.P = 0
		c.EffVal = cas.oper

		Compare(&c, c.A)
		assert.Equal(t, cas.pWant, c.P)
	}
}
