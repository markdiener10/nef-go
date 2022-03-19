# nef-go

Golang New Error Framework

## Usage

First import the library from the open source repo:

import (  
    nef "github.com/markdiener10/nef-go"  
)  
(Preface the the library reference with "nef" for easier package access)  

There is a single call in the interface for error construction and diagnostic data capture:

nef.New(STACKTRACESIZE, VARIABLE PARAMETERS)
	
For convenience, a form of Panic() has been provided to allow for the Nef* object to be properly allocated and then the system panic() function is called with Nef passed in.

The variable parameters can be described:

Reference Code: Supply a code useable for error processing further up the stack  
Previous Error: Pass an error,[]error, or *Nef  
Developer Note: A formatted string for specific error information capture  

All variable parameters are optional. The order of the parameters affect the behavior of the output.Generally pass in a reference code as the first variable parameters.  Then pass in the previous error.  And finally pass a formatted string with additional parameters as described in the standard library fmt.Sprintf().  

While possible to pass in multiple error values, only the first error passed will be captured as the previous error.  Additional errors passed will be ignored.

Some examples are helpful:

nef.New(0) -> Generate *Nef with no stack trace  
nef.New(2) -> Generate *Nef with a stack trace size of 2  
nef.New(0,errors.New("Error String")) -> Generate *Nef with no stack trace and previous error  
nef.New(2,MYREFRENCECODE) -> Generate *Nef with stacksize of 2 and reference code  
nef.New(3,"Developer Note:%d:%d:%s",87,92,"String Parm") -> Generate *Nef with stacksize of 3 and formatted string with parameters  
nef.New(2,MYREFRENCECODE,errors.New("Error String"),"Developer Note:%d:%d:%s",87,92,"String Parm")  

Upon generation of a *Nef value, interface functions are available:

*Nef.Stack() - return an array of NefStackFrame structures that allow for easy formatting of call stack information (Only returns the caller of nef.New() and upwards.)  
*Nef.Code() - package defined error code  
*Nef.Note() - return DevNote with all parameters formatted  
*Nef.Error() - return DevNote but also conform to system error interface  

And a block of access functions for previous error information:  

*Nef.IsPrev() - returns true if there was a previous error passed  
*Nef.PrevErr() - return any previous error value or nil if not passed as a parameter  
*Nef.PrevErrs() - return any previous []error value or nil if not passed as a parameter  
*Nef.PrevNef() - return any previous *Nef value or nil if not passed as a parameter  

## Developer Notes

After years of looking at other teams code which adopted the
standard golang error framework which implements the "Error() string" interface signature, this does little to allow callers to properly respond to the error except to try to parse the error string returned for intelligent processing (Not just PANIC!).  So we allow for custom control codes and let the user implement their own named constants to be included in a package.  Call stack tracing is also provided for easy location of the code that triggered the error, but this information would need to be delivered to the logging subsystem.

The github/pkg/errors is another common package that tries to add some improvements on the base concepts, but the Wrap() does not have easy accessibility for devs and it suffers from the same string parsing challenges.  They also have 4 new interfaces: causer, withStack, withMessage, and fundamental interfaces making it difficult to process them 

We do not try to merge Logging characteristics into the error capture library.  Most teams will have to write an interface library to map this NEF to whatever Logging interface they are using, whether sumologic or corelogix or AWS logstash or etc.

Nothing prevents a developer from construction an Nef value and passing an error value for the PreviousError parameter.  But one should only do this when they are constructing the first Nef and wrapping
a standard error value from code lower down the stack.  Repeated calls to nef.New() should pass in Nef* values up the chain. This way the CastNef() function can repeatedly be called until it runs out of previous errors.

An example loop to extract the Nef/error chain is given in the unit test code nef_test.go 
