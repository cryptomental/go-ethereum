rm -rf core_instrumented
mkdir core_instrumented
go run ../../../github.com/guidovranken/go-coverage-instrumentation/instrument/go-coverage-instrumenter.go core core_instrumented
mv core core_uninstrumented
mv core_instrumented core
