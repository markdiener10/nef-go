package nef

import (
	"errors"
	_ "fmt"
	"strings"
	"testing"
)

func TestNoContent(t *testing.T) {
	nef := New(0)
	stack := nef.Stack()
	if stack != nil { //"No Trace"
		t.Errorf("Stack should be nil:%d", len(*stack))
	}
	if nef.Error() != "" {
		t.Errorf("Stack message should be empty")
	}
	if nef.IsPrevErr() != false {
		t.Errorf("Object should not have previous error")
	}
	if nef.PrevErrs() != nil {
		t.Errorf("Object should not have previous errors array")
	}
}

func TestDevNote(t *testing.T) {

	nef := New(0, "Dev Note Added:%d", 25)
	if nef.Note() == "" {
		t.Errorf("Note should not be nil")
	}
	if !strings.Contains(nef.Note(), "Dev Note Added") {
		t.Errorf("Note was not set correctly")
	}

}

func TestStackSizeGreaterThanZero(t *testing.T) {
	nef := New(2)
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
	nef := New(0, err)
	if nef.PrevErr() == nil {
		t.Errorf("System Error should be retrievable")
	}

}

func TestPreviousNef(t *testing.T) {

	nefPrevious := New(0, "Previous Error")
	nef := New(0, nefPrevious)

	if nef.PrevNef() == nil {
		t.Errorf("Previous Nef should be retrievable")
	}

	count := 0
	prev := nef
	for {
		prev = nef.PrevNef()
		if prev == nil {
			break
		}
		count++
		if prev.PrevErr() != nil {
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
	nefPrevious := New(0, err, "Second Nef Error")
	nef := New(0, nefPrevious, "Third Nef Error")

	count := 0
	prev := nef
	for {

		prev = nef.PrevNef()
		if prev == nil {
			break
		}

		if nef.PrevErr() == nil {
			t.Errorf("Previous System Error should not exist")
		}
		count++
		nef = prev
	}

	if count != 1 {
		t.Errorf("Layered error count != 1:%d", count)
	}

	err = nef.PrevErr()
	if err == nil {
		t.Errorf("Previous System Error should exist")
	}

	//if err != nil {
	//fmt.Println(err.Error())
	//}

}

func TestCodeParameter(t *testing.T) {

	nef := New(0, 35)
	if nef.Code() != 35 {
		t.Errorf("Reference code not recovered:%d", nef.Code())
	}

}

func TestCodeStringPointer(t *testing.T) {

	s := ""

	nef := New(0, 35, &s)
	if nef.Code() != 35 {
		t.Errorf("Reference code not recovered:%d", nef.Code())
	}

	s = "String Pointer"

	nef = New(0, 35, &s)
	if nef.Code() != 35 {
		t.Errorf("Reference code not recovered:%d", nef.Code())
	}

}

func TestAllParameters(t *testing.T) {

	nef := New(0, 45, errors.New("Previous System Error"), "DevNote:%s", "Inserted String")
	if nef.Code() != 45 {
		t.Errorf("Reference code not recovered:%d", nef.Code())
	}
	if nef.Note() != "DevNote:Inserted String" {
		t.Errorf("DevNote not recovered:%s", nef.Note())
	}

}

func TestNefPanic(t *testing.T) {

	defer func() {
		nf := recover()
		if nf == nil {
			t.Errorf("Panic NEF should be non-null")
		}
	}()

	Panic(0, "Dev Note Added:%d", 25)

	t.Errorf("Should not be able to reach post-panic code")

}

func TestCodePrevErrorsArray(t *testing.T) {

	s := "String Pointer"

	errs := []error{errors.New("One"), errors.New("Two")}
	nef := New(0, 35, errs, &s, &s)
	_ = nef.PrevErrs()
	nef = New(0, 35, &errs)
	_ = nef.PrevErrs()

}
