# schemaverify
JSON schema linter for Go

# Install
  
```
go get -u github.com/tdakkota/schemaverify/cmd/schemaverify
```

# Usage
schemaverify uses [`singlechecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/singlechecker) package to run:

```
schemaverify: checks that generated Go structs matches JSON Schema

Usage: schemaverify [-flag] [package]


Flags:
  -V    print version and exit
  -all
        no effect (deprecated)
  -c int
        display offending line with this many lines of context (default -1)
  -cpuprofile string
        write CPU profile to this file
  -debug string
        debug flags, any subset of "fpstv"
  -fix
        apply all suggested fixes
  -flags
        print analyzer flags in JSON
  -json
        emit JSON output
  -memprofile string
        write memory profile to this file
  -report-skipped
        report definition which not found
  -schema-dir string
        path to JSON Schema directory (default "vk-api-schema")
  -source
        no effect (deprecated)
  -tags string
        no effect (deprecated)
  -trace string
        write trace log to this file
  -v    no effect (deprecated)

```
