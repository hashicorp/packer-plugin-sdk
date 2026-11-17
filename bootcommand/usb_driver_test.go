// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package bootcommand

import (
	"context"
	"testing"
	"time"

	"golang.org/x/mobile/event/key"
)

func TestUSBDriver(t *testing.T) {
	tc := []struct {
		command string
		code    key.Code
		shift   bool
	}{
		{
			"<leftShift>",
			key.CodeLeftShift,
			false,
		},
		{
			"<pageUp>",
			key.CodePageUp,
			false,
		},
		{
			"<pageDown>",
			key.CodePageDown,
			false,
		},
		{
			"<leftShiftOff>",
			key.CodeLeftShift,
			false,
		},
		{
			"<leftShiftOn>",
			key.CodeLeftShift,
			true,
		},
		{
			"<leftsuper>",
			key.CodeLeftGUI,
			false,
		},
		{
			"<leftCommand>",
			key.CodeLeftGUI,
			false,
		},
		{
			"<RIGHTCOMMAND>",
			key.CodeRightGUI,
			false,
		},
		{
			"<leftoption>",
			key.CodeLeftAlt,
			false,
		},
		{
			"<rightOption>",
			key.CodeRightAlt,
			false,
		},
		{
			"<spacebar>",
			key.CodeSpacebar,
			false,
		},
		{
			"<menu>",
			key.Code(0x65),
			false,
		},
		{
			"<return>",
			key.CodeReturnEnter,
			false,
		},
		{
			"a",
			key.CodeA,
			false,
		},
		{
			"A",
			key.CodeA,
			true,
		},
		{
			"_",
			key.CodeHyphenMinus,
			true,
		},
	}
	for _, tt := range tc {
		t.Run(tt.command, func(t *testing.T) {
			var code key.Code
			var shift bool
			sendCodes := func(c key.Code, d bool) error {
				code = c
				shift = d
				return nil
			}
			d := NewUSBDriver(sendCodes, time.Duration(0))
			seq, err := GenerateExpressionSequence(tt.command)
			if err != nil {
				t.Fatalf("bad: not expected error: %s", err.Error())
			}
			err = seq.Do(context.Background(), d)
			if err != nil {
				t.Fatalf("bad: not expected error: %s", err.Error())
			}
			if code != tt.code {
				t.Fatalf("bad: wrong scan code: \n expected: %s \n actual: %s", tt.code, code)
			}
			if shift != tt.shift {
				t.Fatalf("bad: wrong shift: \n expected: %t \n actual: %t", tt.shift, shift)
			}
		})
	}
}

// Every special key the PC-XT driver knows must also be known by the USB
// driver, otherwise a boot command that validates fine fails at run time.
func TestUSBDriver_SpecialKeysMatchPCXTDriver(t *testing.T) {
	usb := NewUSBDriver(nil, time.Duration(0))
	pcxt := NewPCXTDriver(nil, 0, time.Duration(0))
	for special := range pcxt.specialMap {
		if _, ok := usb.specialMap[special]; !ok {
			t.Errorf("special key <%s> is supported by the PC-XT driver but not by the USB driver", special)
		}
	}
}

func TestUSBDriver_KeyIntervalNotGiven(t *testing.T) {
	d := NewUSBDriver(nil, time.Duration(0))
	if d.interval != time.Duration(100)*time.Millisecond {
		t.Fatal("not expected key interval")
	}
}

func TestUSBDriver_KeyIntervalGiven(t *testing.T) {
	d := NewUSBDriver(nil, time.Duration(5000)*time.Millisecond)
	if d.interval != time.Duration(5000)*time.Millisecond {
		t.Fatal("not expected key interval")
	}
}
