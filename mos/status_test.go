package mos_test

import (
	"testing"

	"github.com/pevans/erc/mos"
	"github.com/stretchr/testify/assert"
)

func TestApplyStatus(t *testing.T) {
	c := mos.CPU{}

	cases := []struct {
		b        bool
		in, want uint8
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
	c := mos.CPU{}

	cases := []struct {
		in, want uint8
	}{
		{64, 0},
		{128, mos.NEGATIVE},
	}

	for _, cas := range cases {
		c.ApplyN(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestApplyZ(t *testing.T) {
	c := mos.CPU{}

	cases := []struct {
		in, want uint8
	}{
		{0, mos.ZERO},
		{1, 0},
	}

	for _, cas := range cases {
		c.ApplyZ(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestApplyNZ(t *testing.T) {
	c := mos.CPU{}

	cases := []struct {
		in, want uint8
	}{
		{0, mos.ZERO},
		{1, 0},
		{128, mos.NEGATIVE},
	}

	for _, cas := range cases {
		c.ApplyNZ(cas.in)
		assert.Equal(t, cas.want, c.P)
	}
}

func TestCompare(t *testing.T) {
	c := mos.CPU{}

	cases := []struct {
		a, oper, pWant uint8
	}{
		{1, 3, mos.NEGATIVE},
		{3, 1, mos.CARRY},
		{3, 3, mos.ZERO | mos.CARRY},
	}

	for _, cas := range cases {
		c.A = cas.a
		c.P = 0
		c.EffVal = cas.oper

		mos.Compare(&c, c.A)
		assert.Equal(t, cas.pWant, c.P)
	}
}
