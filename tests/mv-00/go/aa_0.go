package main

import (
	"encoding/json"
	"flag"
	"fmt"
	jm "jsonmodel/json_model/runtime/go/jsonmodel"
	"os"
)

type Checker func(interface{}, *jm.Path, *jm.Report) bool

var check_model_map map[string]Checker

// check $ (.)
func json_model_1(val interface{}, path *jm.Path, rep *jm.Report) bool {
    // .
    var res bool = jm.IsArray(val)
    if res {
        for arr_0_idx, arr_0_item := range jm.AsArray(val) {
            var arr_0_lpath *jm.Path = jm.ExtendPathIndex(path, int(arr_0_idx))
            // .0
            res = jm.IsArray(arr_0_item)
            if res {
                for arr_1_idx, arr_1_item := range jm.AsArray(arr_0_item) {
                    var arr_1_lpath *jm.Path = jm.ExtendPathIndex(jm.SelectPath(arr_0_lpath, path != nil), int(arr_1_idx))
                    // .0.0
                    res = jm.IsString(arr_1_item)
                    if ! res {
                        if rep != nil { rep.Add("unexpected string [.0.0]", jm.SelectPath(arr_1_lpath, jm.SelectPath(arr_0_lpath, path != nil) != nil)) }
                        break
                    }
                }
            }
            if ! res {
                if rep != nil { rep.Add("not array or unexpected array [.0]", jm.SelectPath(arr_0_lpath, path != nil)) }
                break
            }
        }
    }
    if ! res {
        if rep != nil { rep.Add("not array or unexpected array [.]", path) }
    }
    return res
}


var initialized bool

func check_model_init() {
	if !initialized {
		defer func() {
			if r := recover(); r != nil {
				panic(fmt.Sprintf("cannot initialize model checker: %v", r))
			}
		}()

        check_model_map = map[string]Checker{
            "": json_model_1,
        }
		initialized = true
	}
}

func check_model_free() {
	if initialized {
		initialized = false
	}
}

func main() {
    testMode := flag.Bool("t", false, "run in test mode with a list of [[expected, val], ...]")
    flag.Parse()
    args := flag.Args()

    if len(args) < 1 {
        fmt.Printf("Usage: %s [-t] <json_file>\n", os.Args[0])
        os.Exit(1)
    }

    data, err := os.ReadFile(args[0])
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    check_model_init()
    defer check_model_free()

    if *testMode {
        var cases [][]interface{}
        if err := json.Unmarshal(data, &cases); err != nil {
            fmt.Printf("Test suite JSON error: %v\n", err)
            os.Exit(1)
        }

        failed := 0
        for i, c := range cases {
            if len(c) < 2 { continue }
            expected, _ := c[0].(bool)
            input := c[1]

            report := &jm.Report{}
            isValid := check_model_map[""](input, nil, report)

            if isValid == expected {
                fmt.Printf("Test #%d: PASS\n", i)
            } else {
                fmt.Printf("Test #%d: FAIL (expected %v, got %v)\n", i, expected, isValid)
                for _, errMsg := range report.Errors {
                    fmt.Printf("  - %s\n", errMsg)
                }
                failed++
            }
        }
        if failed > 0 {
            fmt.Printf("\nDone: %d tests failed\n", failed)
            os.Exit(1)
        }
        fmt.Println("\nDone: All tests passed")
    } else {
        // Standard single-file validation logic
        var input interface{}
        if err := json.Unmarshal(data, &input); err != nil {
            fmt.Printf("JSON parsing error: %v\n", err)
            os.Exit(1)
        }
        report := &jm.Report{}
        isValid := check_model_map[""](input, nil, report)
        if isValid {
            fmt.Println("✅ Valid")
        } else {
            fmt.Println("❌ Invalid:")
            for _, errMsg := range report.Errors {
                fmt.Printf("  - %s\n", errMsg)
            }
            os.Exit(1)
        }
    }
}
