// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package bootcommand

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type event struct {
	u    uint32
	down bool
}

type sender struct {
	e []event
}

func (s *sender) KeyEvent(u uint32, down bool) error {
	s.e = append(s.e, event{u, down})
	return nil
}

func Test_vncSpecialLookup(t *testing.T) {
	in := "<rightShift><rightshiftoff><RIGHTSHIFTON><menu>"
	expected := []event{
		{0xFFE2, true},
		{0xFFE2, false},
		{0xFFE2, false},
		{0xFFE2, true},
		{0xFF67, true},
		{0xFF67, false},
	}
	s := &sender{}
	d := NewVNCDriver(s, time.Duration(0))
	seq, err := GenerateExpressionSequence(in)
	assert.NoError(t, err)
	err = seq.Do(context.Background(), d)
	assert.NoError(t, err)
	assert.Equal(t, expected, s.e)
}

func Test_vncShiftSequence(t *testing.T) {
	in := "AbC"
	expected := []event{
		{0xFFE1, true},
		{0x041, true},
		{0x041, false},
		{0xFFE1, false},
		{0x062, true},
		{0x062, false},
		{0xFFE1, true},
		{0x043, true},
		{0x043, false},
		{0xFFE1, false},
	}
	s := &sender{}
	d := NewVNCDriver(s, time.Duration(0))
	seq, err := GenerateExpressionSequence(in)
	assert.NoError(t, err)
	err = seq.Do(context.Background(), d)
	assert.NoError(t, err)
	assert.Equal(t, expected, s.e)
}

func Test_vncShiftSingle(t *testing.T) {
	in := "<Aon><Aoff>"
	expected := []event{
		{0xFFE1, true},
		{0x041, true},
		{0x041, false},
		{0xFFE1, false},
	}
	s := &sender{}
	d := NewVNCDriver(s, time.Duration(0))
	seq, err := GenerateExpressionSequence(in)
	assert.NoError(t, err)
	err = seq.Do(context.Background(), d)
	assert.NoError(t, err)
	assert.Equal(t, expected, s.e)
}

func Test_vncIntervalNotGiven(t *testing.T) {
	s := &sender{}
	d := NewVNCDriver(s, time.Duration(0))
	assert.Equal(t, d.interval, time.Duration(100)*time.Millisecond)
}

func Test_vncIntervalGiven(t *testing.T) {
	s := &sender{}
	d := NewVNCDriver(s, time.Duration(5000)*time.Millisecond)
	assert.Equal(t, d.interval, time.Duration(5000)*time.Millisecond)
}
