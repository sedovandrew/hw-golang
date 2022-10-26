Before optimizations:
```bash
$ go test -v -count=1 -timeout=30s -tags bench .
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:45: time used: 616.923983ms / 300ms
    stats_optimization_test.go:46: memory used: 294Mb / 30Mb
    assertion_compare.go:332:
                Error Trace:    stats_optimization_test.go:48
                Error:          "616923983" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too slow]
--- FAIL: TestGetDomainStat_Time_And_Memory (20.15s)
FAIL
FAIL    github.com/sedovandrew/hw10_program_optimization        20.154s
FAIL
```

Refused regular expressions:
```bash
$ go test -v -count=1 -timeout=30s -tags bench .
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:45: time used: 367.001795ms / 300ms
    stats_optimization_test.go:46: memory used: 162Mb / 30Mb
    assertion_compare.go:332:
                Error Trace:    stats_optimization_test.go:48
                Error:          "367001795" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too slow]
--- FAIL: TestGetDomainStat_Time_And_Memory (6.87s)
FAIL
FAIL    github.com/sedovandrew/hw10_program_optimization        6.872s
FAIL
```

Replaced "encoding/json" with "github.com/goccy/go-json":
```bash
$ go test -v -count=1 -timeout=30s -tags bench .
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:45: time used: 224.886218ms / 300ms
    stats_optimization_test.go:46: memory used: 148Mb / 30Mb
    assertion_compare.go:332:
                Error Trace:    stats_optimization_test.go:49
                Error:          "156032824" is not less than "31457280"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too greedy]
--- FAIL: TestGetDomainStat_Time_And_Memory (2.63s)
FAIL
FAIL    github.com/sedovandrew/hw10_program_optimization        2.630s
FAIL
```

Removed unnecassary fields from the structure:
```bash
$ go test -v -count=1 -timeout=30s -tags bench .
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:45: time used: 153.455524ms / 300ms
    stats_optimization_test.go:46: memory used: 139Mb / 30Mb
    assertion_compare.go:332:
                Error Trace:    stats_optimization_test.go:49
                Error:          "146429904" is not less than "31457280"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too greedy]
--- FAIL: TestGetDomainStat_Time_And_Memory (1.85s)
FAIL
FAIL    github.com/sedovandrew/hw10_program_optimization        1.853s
FAIL
```

Reduced memory consumption by line-by-line reading:
```bash
$ go test -v -count=1 -timeout=30s -tags bench .
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:45: time used: 122.720115ms / 300ms
    stats_optimization_test.go:46: memory used: 18Mb / 30Mb
--- PASS: TestGetDomainStat_Time_And_Memory (1.39s)
PASS
ok      github.com/sedovandrew/hw10_program_optimization        1.390s
```
