# go-nef
Golang New Error Framework

After years of looking at other teams code which adopted the
standard golang error framework which implements the "Error() string" interface signature, this does little to allow callers to properly respond to the error except to try to parse the error string returned for intelligent processing (Not just PANIC!).  So we allow for custom control codes and let the user implement their own named constants to be included in a package.  Call stack tracing is also provided for easy location of the code that triggered the error.

The github/pkg/errors is another common package that tries to add some improvements on the base concepts, but the Wrap() does not have easy accessibility for devs and it suffers from the same string parsing challenges.  They also have 4 new interfaces: causer, withStack, withMessage, and fundamental interfaces making it difficult to process them 

We do not try to merge Logging characteristics into the error capture library.  Most teams will have to write an interface library to map this NEF to whatever Logging interface they are using, whether sumologic or corelogix or AWS logstash.
