package nef

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNoContent(t *testing.T) {
	nef := New(0, nil, 0, "")
	stack := nef.Stack()
	if stack != nil { //"No Trace"
		t.Errorf("Stack should be nil:%d", len(*stack))
	}
}

func TestDevNote(t *testing.T) {

	nef := New(0, nil, 0, "Dev Note Added:%d", 25)
	if nef.Note() == nil {
		t.Errorf("Note should not be nil")
	}
	if !strings.Contains(*nef.Note(), "Dev Note Added") {
		t.Errorf("Note was not set correctly")
	}

}

func TestStackSizeGreaterThanZero(t *testing.T) {
	nef := New(10, nil, 0, "")
	stack := nef.Stack()
	if stack == nil {
		t.Errorf("Stack should be greater than zero (A)")
	}
	if len(*stack) != 2 {
		t.Errorf("Stack depth should be 2:%d", len(*stack))
	}
	//for index, frame := range *stack {
	//	fmt.Println(index, frame)
	//}
}

func TestPreviousSystemError(t *testing.T) {

	err := errors.New("Previous System Error")
	nef := New(0, err, 0, "")
	if nef.CastErr() == nil {
		t.Errorf("System Error should be retrievable")
	}

}

func TestPreviousNef(t *testing.T) {

	nefPrevious := New(0, nil, 0, "Previous Error")
	nef := New(0, nefPrevious, 0, "Current Error")

	if nef.CastNef() == nil {
		t.Errorf("Previous Nef should be retrievable")
	}

	count := 0
	prev := nef
	for {
		prev = nef.CastNef()
		if prev == nil {
			break
		}
		count++
		if prev.CastErr() != nil {
			t.Errorf("Previous System Error should not exist")
		}
		nef = prev
	}
	if count != 1 {
		t.Errorf("Layered error count != 1")
	}

}

func TestPreviousNefAndSystemError(t *testing.T) {

	err := errors.New("First System Error")
	nefPrevious := New(0, err, 0, "Second Nef Error")
	nef := New(0, nefPrevious, 0, "Third Nef Error")

	count := 0
	prev := nef
	for {

		prev = nef.CastNef()
		if prev == nil {
			break
		}

		if nef.CastErr() == nil {
			t.Errorf("Previous System Error should not exist")
		}
		count++
		nef = prev
	}

	if count != 1 {
		t.Errorf("Layered error count != 1:%d", count)
	}

	err = nef.CastErr()
	if err == nil {
		t.Errorf("Previous System Error should exist")
	}

	if err != nil {
		fmt.Println(err.Error())
	}

}
