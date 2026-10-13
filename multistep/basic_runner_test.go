// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package multistep

import (
	"context"
	"reflect"
	"testing"
)

func TestBasicRunner_ImplRunner(t *testing.T) {
	var raw interface{} = &BasicRunner{}
	if _, ok := raw.(Runner); !ok {
		t.Fatalf("BasicRunner must be a Runner")
	}
}

func TestBasicRunner_Run(t *testing.T) {
	data := new(BasicStateBag)
	stepA := &TestStepAcc{Data: "a"}
	stepB := &TestStepAcc{Data: "b"}

	r := &BasicRunner{Steps: []Step{stepA, stepB}}
	r.Run(context.Background(), data)

	// Test run data
	expected := []string{"a", "b"}
	results := data.Get("data").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test cleanup data
	expected = []string{"b", "a"}
	results = data.Get("cleanup").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test no halted or canceled
	if _, ok := data.GetOk(StateCancelled); ok {
		t.Errorf("canceled should not be in state bag")
	}

	if _, ok := data.GetOk(StateHalted); ok {
		t.Errorf("halted should not be in state bag")
	}
}

func TestBasicRunner_Run_Halt(t *testing.T) {
	data := new(BasicStateBag)
	stepA := &TestStepAcc{Data: "a"}
	stepB := &TestStepAcc{Data: "b", Halt: true}
	stepC := &TestStepAcc{Data: "c"}

	r := &BasicRunner{Steps: []Step{stepA, stepB, stepC}}
	r.Run(context.Background(), data)

	// Test run data
	expected := []string{"a", "b"}
	results := data.Get("data").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test cleanup data
	expected = []string{"b", "a"}
	results = data.Get("cleanup").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test that it says it is halted
	halted := data.Get(StateHalted).(bool)
	if !halted {
		t.Errorf("not halted")
	}
}

// confirm that can't run twice
func TestBasicRunner_Run_Run(t *testing.T) {
	defer func() {
		recover()
	}()
	ch := make(chan chan bool)
	stepInt := &TestStepSync{ch}
	stepWait := &TestStepWaitForever{}
	r := &BasicRunner{Steps: []Step{stepInt, stepWait}}

	go r.Run(context.Background(), new(BasicStateBag))
	// wait until really running
	<-ch

	// now try to run aain
	r.Run(context.Background(), new(BasicStateBag))

	// should not get here in nominal codepath
	t.Errorf("Was able to run an already running BasicRunner")
}

func TestBasicRunner_Cancel(t *testing.T) {

	topCtx, topCtxCancel := context.WithCancel(context.Background())

	checkCanceled := func(data StateBag) {
		canceled := data.Get(StateCancelled).(bool)
		if !canceled {
			t.Fatal("state should be canceled")
		}
	}

	data := new(BasicStateBag)
	r := &BasicRunner{}
	r.Steps = []Step{
		&TestStepAcc{Data: "a"},
		&TestStepAcc{Data: "b"},
		TestStepFn{
			run: func(ctx context.Context, sb StateBag) StepAction {
				return ActionContinue
			},
			cleanup: checkCanceled,
		},
		TestStepFn{
			run: func(ctx context.Context, sb StateBag) StepAction {
				topCtxCancel()
				<-ctx.Done()
				return ActionContinue
			},
			cleanup: checkCanceled,
		},
		TestStepFn{
			run: func(context.Context, StateBag) StepAction {
				t.Fatal("I should not be called")
				return ActionContinue
			},
			cleanup: func(StateBag) {
				t.Fatal("I should not be called")
			},
		},
	}

	r.Run(topCtx, data)

	// Test run data
	expected := []string{"a", "b"}
	results := data.Get("data").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test cleanup data
	expected = []string{"b", "a"}
	results = data.Get("cleanup").([]string)
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("unexpected result: %#v", results)
	}

	// Test that it says it is canceled
	checkCanceled(data)

}

func TestBasicRunner_Cancel_Special(t *testing.T) {
	stepOne := &TestStepInjectCancel{}
	stepTwo := &TestStepInjectCancel{}
	r := &BasicRunner{Steps: []Step{stepOne, stepTwo}}

	state := new(BasicStateBag)
	state.Put("runner", r)
	r.Run(context.Background(), state)

	// test that state contains canceled
	if _, ok := state.GetOk(StateCancelled); !ok {
		t.Errorf("canceled should be in state bag")
	}
}
