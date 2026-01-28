package main

import (
	"encoding/json"
	"fmt"
	jm "jsonmodel/json_model/runtime/go/jsonmodel"
	"os"
	"regexp"
)

type Checker func(interface{}, *jm.Path, *jm.Report) bool

var _jm_re_0_re *regexp.Regexp
var check_model_map map[string]Checker

func _jm_re_0(val string, path *jm.Path, rep *jm.Report) bool {
    return _jm_re_0_re.MatchString(val)
}

// check $ (.)
func json_model_1(val interface{}, path *jm.Path, rep *jm.Report) bool {
    // A person with a birth date
    // .
    // check close must only props
    if ! jm.IsObject(val) {
        if rep != nil { rep.Add("not an object [.]", path) }
        return false
    }
    if jm.Len(val) != 2 {
        if rep != nil { rep.Add("bad property count [.]", path) }
        return false
    }
    var lpath *jm.Path
    var pval interface{}
    var res bool
    if ! jm.ObjectHasPropVal(val, "name", &pval) {
        if rep != nil { rep.Add("missing mandatory prop <name> [.]", path) }
        return false
    }
    lpath = jm.ExtendPath(path, "name")
    // .name
    // "/^[a-z]+$/i"
    res = jm.IsString(pval) && _jm_re_0(jm.AsString(pval), jm.SelectPath(lpath, path != nil), rep)
    if ! res {
        if rep != nil { rep.Add("unexpected /^[a-z]+$/i [.name]", jm.SelectPath(lpath, path != nil)) }
        if rep != nil { rep.Add("unexpected value for mandatory prop <name> [.]", jm.SelectPath(lpath, path != nil)) }
        return false
    }
    if ! jm.ObjectHasPropVal(val, "born", &pval) {
        if rep != nil { rep.Add("missing mandatory prop <born> [.]", path) }
        return false
    }
    lpath = jm.ExtendPath(path, "born")
    // .born
    res = jm.IsString(pval) && jm.IsValidDate(jm.AsString(pval))
    if ! res {
        if rep != nil { rep.Add("unexpected $DATE [.born]", jm.SelectPath(lpath, path != nil)) }
        if rep != nil { rep.Add("unexpected value for mandatory prop <born> [.]", jm.SelectPath(lpath, path != nil)) }
        return false
    }
    return true
}


var initialized bool

func check_model_init() {
	if !initialized {
		defer func() {
			if r := recover(); r != nil {
				panic(fmt.Sprintf("cannot initialize model checker: %v", r))
			}
		}()

        _jm_re_0_re = regexp.MustCompile("(?i)^[a-z]+$")
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
    if len(os.Args) < 2 {
        fmt.Printf("Usage: %s <json_file>\n", os.Args[0])
        os.Exit(1)
    }

    // Lecture du fichier spécifié en argument
    data, err := os.ReadFile(os.Args[1])
    if err != nil {
        fmt.Printf("Erreur lors de la lecture du fichier: %v\n", err)
        os.Exit(1)
    }

    var input interface{}
    if err := json.Unmarshal(data, &input); err != nil {
        fmt.Printf("Erreur de parsing JSON: %v\n", err)
        os.Exit(1)
    }

    // Initialisation du validateur (regex, etc.)
    check_model_init()
    report := &jm.Report{}

    // Appel de la validation via l'entrée par défaut
    isValid := check_model_map[""](input, nil, report)

    if isValid {
        fmt.Println("✅ Le JSON est valide selon le modèle !")
    } else {
        fmt.Println("❌ Le JSON est invalide :")
        for _, errMsg := range report.Errors {
            fmt.Printf("  - %s\n", errMsg)
        }
        os.Exit(1)
    }
}
