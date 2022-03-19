package nef

//There are a lot of assumptions in the golang runtime and builtin code that could be surfaced
//to make the formatting of our stack trace more robust and quicker for devs to locate code
//and debug microservices

//TODO: There may be an issue with stripping local directory references in the build utility of go
//go build -trimpath - strip local file name references from cmd and asm

import (
	"fmt"
	"runtime"
)

//New error framework interface (for reflection purposes)
type Nefi interface {
	Error() string
	Code() int
}

//Condensed version of runtime.Frame structure
type NefStackFrame struct {
	File     string
	Line     int
	Function string
}

//New error framework transport
type Nef struct {
	stackTrace    *[]*NefStackFrame
	previousError interface{}
	devNote       string
	codeReference int
}

//Conform to Errors.Error() interface for compatibility with builtin package
func (g *Nef) Error() string {
	return g.devNote
}

func (g *Nef) Code() int {
	return g.codeReference
}

func (g *Nef) Note() string {
	return g.devNote
}

func (g *Nef) Stack() *[]*NefStackFrame {
	return g.stackTrace
}

func (g *Nef) IsPrevErr() bool {
	return g.previousError != nil
}

func (g *Nef) PrevErr() error {
	err, ok := g.previousError.(error)
	if !ok {
		return nil
	}
	return err
}

func (g *Nef) PrevErrs() []error {

	err, ok := g.previousError.([]error)
	if ok {
		return err
	}
	erra, ok := g.previousError.(*[]error)
	if ok {
		return *erra
	}
	return nil
}

func (g *Nef) PrevNef() *Nef {
	nef, ok := g.previousError.(*Nef)
	if !ok {
		return nil
	}
	return nef
}

//Original signature
//func New(stackSize uint, previousError interface{}, codeReference int, devNote string, args ...interface{}) *Nef {

//New signature
func New(stackSize uint, parms ...interface{}) *Nef {

	var pStack *[]*NefStackFrame = nil

	if stackSize > 0 {
		stackTrace := make([]*NefStackFrame, 0)

		programCounters := make([]uintptr, stackSize)
		//Only grab the caller of this New() function, not the New() function itself or the runtime.Callers() function
		//So we skip 2 on the Callers stack to yield the caller of this function.
		numberOfFrames := runtime.Callers(2, programCounters)
		if numberOfFrames > 0 {
			stackFrames := runtime.CallersFrames(programCounters[0:numberOfFrames])
			for {
				frame, more := stackFrames.Next()
				if frame.Function == "runtime.goexit" {
					break
				}
				nefFrame := &NefStackFrame{File: frame.File, Function: frame.Function, Line: frame.Line}
				stackTrace = append(stackTrace, nefFrame)
				pStack = &stackTrace
				if !more {
					break
				}
			}
		}
	}

	var gpreverr interface{} = nil
	var gdevnote string = ""
	var gsubparms = make([]interface{}, 0)
	var gcode int = 0
	var gparm interface{}

	for _, gparm = range parms {

		switch gparm.(type) {
		case error, []error, *error, *[]error, *Nef, Nef:
			//Only take first error
			if gpreverr != nil {
				continue
			}
			gpreverr = gparm
			continue
		case string: //DevNote afterwards by convention
			if gdevnote == "" {
				gdevnote = gparm.(string)
				continue
			}
			gsubparms = append(gsubparms, gparm)
			continue
		case *string:
			if gdevnote == "" {
				gdevnote = *(gparm.(*string))
				continue
			}
			gsubparms = append(gsubparms, gparm)
			continue
		case int:
			//Put the reference code before the devnote in the parms
			if gdevnote == "" {
				if gcode != 0 {
					continue
				}
				gcode = (gparm.(int))
				continue
			}
			gsubparms = append(gsubparms, gparm)
		default:
			gsubparms = append(gsubparms, gparm)
		}
	}

	if gdevnote != "" {
		gdevnote = fmt.Sprintf(gdevnote, gsubparms...)
	}

	//and keep going up until we reach the top
	return &Nef{
		stackTrace:    pStack,
		previousError: gpreverr,
		codeReference: gcode,
		devNote:       gdevnote,
	}
}

func Panic(stackSize uint, args ...interface{}) {
	nef := New(stackSize, args...)
	panic(nef)
}
