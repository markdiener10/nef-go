# nef-go

Golang New Error Framework

## Usage

First import the library from the open source repo:
import (
	nef "github.com/markdiener10/nef-go"
)

Usually preface the name of the library with "nef" for easier package scope access

There is a single call in the interface for error construction and capture

nef := nef.New(10, err, MYERRORCONSTANT, "Place your devops note here:%d %s",MYTESTVALUE,MYTESTSTRING)
	
For convenience, a form of Panic() has been provided to allow for the Nef* object to be properly allocated
and then the system panic() function is called with Nef passed in.

The parameters and returns are as follows:

func New(StackSize,PreviousErrorOrNef,ConstantValue,DevNote) *Nef

1) StackSize - Set the max depth of the captured call stack from the point of calling nef.New() and upwards. Values less than 1 disables the feature.  Generally 10 is a good number.
2) PreviousErrororNefPtr - Previous err or Nef*. (This is similar to the Wrap() in github/pkg/errors)
3) ConstantValue - A library user defined code to pass upwards to the package user.  Generally use a const NAME = VAL statement to define specific error or control codes.
4) DevNote - The equivalent of "Message" in err.Error() string. Generally to capture specific data at the  point of call.

For later processing, receiver functions are defined to allow for easier processing of the error data.

Nef.Error() - conform to the error interface for compliance reasons (returns DevNote)
Nef.Code() - package defined error code
Nef.Note() - return DevNote as a formatted string
Nef.PrevNef() - return any previous Nef value. Nil if Nef was not constructed with a prior *Nef value
Nef.PrevErr() - return any previous error value. Nil if Nef was not constructed with a prior error value
Nef.Stack() - return an array of NefStackFrame structures that allow for easy formatting of call stack information (Only returns the caller of nef.New() and upwards.

## Developer Notes

After years of looking at other teams code which adopted the
standard golang error framework which implements the "Error() string" interface signature, this does little to allow callers to properly respond to the error except to try to parse the error string returned for intelligent processing (Not just PANIC!).  So we allow for custom control codes and let the user implement their own named constants to be included in a package.  Call stack tracing is also provided for easy location of the code that triggered the error, but this information would need to be delivered to the logging subsystem.

The github/pkg/errors is another common package that tries to add some improvements on the base concepts, but the Wrap() does not have easy accessibility for devs and it suffers from the same string parsing challenges.  They also have 4 new interfaces: causer, withStack, withMessage, and fundamental interfaces making it difficult to process them 

We do not try to merge Logging characteristics into the error capture library.  Most teams will have to write an interface library to map this NEF to whatever Logging interface they are using, whether sumologic or corelogix or AWS logstash or etc.

Nothing prevents a developer from construction an Nef value and passing an error value for the PreviousError parameter.  But one should only do this when they are constructing the first Nef and wrapping
a standard error value from code lower down the stack.  Repeated calls to nef.New() should pass in Nef* values up the chain. This way the CastNef() function can repeatedly be called until it runs out of previous errors.

An example loop to extract the Nef/error chain is given in the unit test code nef_test.go 
