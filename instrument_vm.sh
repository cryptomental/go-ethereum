GO_INSTRUMENT=/media/jhg/5c980a66-cf42-464d-878a-85016937237c/gopkg/src/github.com/guidovranken/go-coverage-instrumentation/instrument/main
cd core
rm -rf vm_instrumented
mkdir vm_instrumented
$GO_INSTRUMENT vm vm_instrumented
