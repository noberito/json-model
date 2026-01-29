import json

from .language import (
    Block,
    BoolExpr,
    ConstList,
    Expr,
    IntExpr,
    JsonExpr,
    JsonScalar,
    Language,
    PathExpr,
    PropMap,
    StrExpr,
    Var,
)
from .mtypes import Conditionals, Number, TestHint


class Go(Language):
    """Go language Code Generator."""

    def __init__(
        self,
        *,
        relib: str = "regexp",
        debug: bool = False,
        with_report: bool = True,
        with_path: bool = True,
        with_predef: bool = True,
        with_package: bool = True,
    ):
        super().__init__(
            "Go",
            relib=relib,
            debug=debug,
            with_report=with_report,
            with_path=with_path,
            with_predef=with_predef,
            with_package=with_package,
            not_op="!",
            and_op="&&",
            or_op="||",
            lcom="//",
            true="true",
            false="false",
            null="nil",
            check_t="bool",
            json_t="interface{}",
            path_t="*jm.Path",
            float_t="float64",
            str_t="string",
            bool_t="bool",
            int_t="int64",
        )

        assert relib in ("regexp"), f"support for regexp (standard library), not {relib}"
        self.reindent = True
        self._report_t = "*jm.Report"

    def _var(self, var: Var, val: Expr | None, tname: str | None) -> Block:
        """Declares and/or assigns a variable in Go."""
        if tname:
            # Explicit declaration: var name type = value
            assign = f" = {val}" if val is not None else ""
            return [f"var {var} {tname}{assign}"]
        else:
            # Simple assignment: name = value
            return [f"{var} = {val}"]

    #
    # Type tests
    #
    def is_num(self, var: Var) -> BoolExpr:
        return f"jm.IsNumber({var})"

    def is_def(self, var: Var) -> BoolExpr:
        return f"{var} != nil"

    def is_scalar(self, var: Var) -> BoolExpr:
        return f"jm.IsScalar({var})"

    def is_a(self, var: Var, tval: type | None, loose: bool | None = None) -> BoolExpr:
        if tval is None or tval is type(None):
            return f"{var} == nil"
        elif tval is bool:
            return f"jm.IsBool({var})"
        elif tval is int:
            check = f"jm.IsInteger({var})"
            if loose:
                # In Go, we check if it is a float with no fractional part
                check = f"({check} || (jm.IsFloat({var}) && jm.AsFloat({var}) == float64(int64(jm.AsFloat({var})))))"
            return check
        elif tval is float:
            return self.is_num(var) if loose else f"jm.IsFloat({var})"
        elif tval is Number:
            return self.is_num(var)
        elif tval is str:
            return f"jm.IsString({var})"
        elif tval is list:
            return f"jm.IsArray({var})"
        elif tval is dict:
            return f"jm.IsObject({var})"
        return "false"

    #
    # Predefined checks (Predefs)
    #
    def predef(self, var: Var, name: str, path: Var, is_str: bool = False) -> BoolExpr:
        # If predefs are disabled, just validate the string type
        if not self._with_predef and self.str_content_predef(name):
            return self.const(True) if is_str else self.is_a(var, str)

        # Prepare type prefix and string conversion
        prefix = "" if is_str else f"jm.IsString({var}) && "
        val = var if is_str else f"jm.AsString({var})"

        # Mapping predefs to the Go runtime functions
        rt_funcs = {
            "$UUID": "IsValidUUID",
            "$DATE": "IsValidDate",  # Calls jm.IsValidDate(val)
            "$TIME": "IsValidTime",
            "$DATETIME": "IsValidDateTime",
            "$REGEX": "IsValidRegex",
            "$URL": "IsValidURL",
            "$URI": "IsValidURI",
            "$EMAIL": "IsValidEmail",
            "$JSON": "IsValidJSON",
        }

        if name in rt_funcs:
            # FIX: Remove {path} and rep to use only a single argument
            return f"{prefix}jm.{rt_funcs[name]}({val})"

        return super().predef(var, name, path, is_str)

    #
    # Expressions & Extraction
    #
    def value(self, var: Var, tvar: type) -> Expr:
        if tvar is bool:
            return f"jm.AsBool({var})"
        if tvar is int:
            return f"jm.AsInt({var})"
        if tvar is float:
            return f"jm.AsFloat({var})"
        if tvar is str:
            return f"jm.AsString({var})"
        return var

    def obj_prop_val(self, obj: Var, prop: str | StrExpr, is_var: bool = False) -> JsonExpr:
        p = prop if is_var else self.esc(prop)
        return f"jm.ObjectGet({obj}, {p})"

    def has_prop(self, obj: Var, prop: str) -> BoolExpr:
        return f"jm.ObjectHasProp({obj}, {self.esc(prop)})"

    def obj_has_prop_val(
        self, dst: Var, obj: Var, prop: str | StrExpr, is_var: bool = False
    ) -> BoolExpr:
        # Go does not support assignment within comparison expressions like Python (:=),
        # except within an 'if' statement. This function is frequently used by the optimizer.
        return f"jm.ObjectHasPropVal({obj}, {self.esc(prop) if not is_var else prop}, &{dst})"

    def any_len(self, var: Var) -> IntExpr:
        return f"jm.Len({var})"

    #
    # Loops & Blocks
    #
    def arr_loop(self, arr: Var, idx: Var, val: Var, body: Block) -> Block:
        return [f"for {idx}, {val} := range jm.AsArray({arr}) {{"] + self.indent(body) + ["}"]

    def obj_loop(self, obj: Var, key: Var, val: Var, body: Block) -> Block:
        return [f"for {key}, {val} := range jm.AsObject({obj}) {{"] + self.indent(body) + ["}"]

    def int_loop(self, idx: Var, start: IntExpr, end: IntExpr, body: Block) -> Block:
        return [f"for {idx} := {start}; {idx} < {end}; {idx}++ {{"] + self.indent(body) + ["}"]

    def if_stmt(
        self, cond: BoolExpr, true: Block, false: Block = [], likely: TestHint = None
    ) -> Block:
        code = [f"if {cond} {{"] + self.indent(true)
        if false:
            code += ["} else {"] + self.indent(false)
        code += ["}"]
        return code

    def mif_stmt(self, cond_true: Conditionals, false: Block = []) -> Block:
        code, op = [], "if"
        for cond, likely, true in cond_true:
            code += [f"{op} {cond} {{"]
            code += self.indent(true)
            op = "} else if"
        if false:
            code += ["} else {"] + self.indent(false)
        code += ["}"]
        return code

    #
    # Reporting & Path
    #
    def is_reporting(self) -> BoolExpr:
        return "rep != nil"

    def report(self, msg: str, path: Var) -> Block:
        return (
            [f"if rep != nil {{ rep.Add({self.esc(msg)}, {path}) }}"] if self._with_report else []
        )

    def clean_report(self) -> Block:
        return ["if rep != nil { rep.Clear() }"]

    def path_val(self, pvar: Var, pseg: str | int, is_prop: bool, is_var: bool) -> PathExpr:
        if not self._with_path:
            return "nil"

        if is_prop:
            # For object properties: use ExtendPath(parent, "name")
            name = pseg if is_var else self.esc(pseg)
            return f"jm.ExtendPath({pvar}, {name})"
        else:
            # For array indices: use ExtendPathIndex(parent, index)
            # We cast to int() to be safe with different Go int types
            return f"jm.ExtendPathIndex({pvar}, int({pseg}))"

    def path_lvar(self, lvar: Var, rvar: Var) -> PathExpr:
        return f"jm.SelectPath({lvar}, {rvar} != nil)" if self._with_path else "nil"

    #
    # Function definitions
    #

    def sub_re(self, name: str, regex: str, opts: str) -> Block:
        return [
            f"func {name}(val string, path *jm.Path, rep *jm.Report) bool {{",
            f"    return {name}_re.MatchString(val)",
            "}",
        ]

    def sub_fun(self, name: str, body: Block, inline: bool = False) -> Block:
        return (
            [f"func {name}(val interface{{}}, path *jm.Path, rep *jm.Report) bool {{"]
            + self.indent(body)
            + ["}"]
        )

    def def_pmap(self, name: str, pmap: PropMap, public: bool) -> Block:
        return [f"var {name} map[string]Checker"]

    def ini_pmap(self, name: str, pmap: PropMap, public: bool) -> Block:
        lines = [f"{name} = map[string]Checker{{"]
        for p, f in pmap.items():
            lines.append(f"    {self.esc(p)}: {f},")
        lines.append("}")
        return lines

    #
    # Regexp
    #
    def def_re(self, name: str, regex: str, opts: str) -> Block:
        return [f"var {name}_re *regexp.Regexp"]

    def ini_re(self, name: str, regex: str, opts: str) -> Block:
        sregex = self.esc((f"(?{opts})" if opts else "") + regex)
        return [f"{name}_re = regexp.MustCompile({sregex})"]

    def match_re(self, name: str, var: str, regex: str, opts: str) -> BoolExpr:
        return f"{name}_re.MatchString({var})"

    #
    # File handling
    #

    def file_header(self, exe: bool, package: str | None) -> Block:
        pkg = package if package else "main"
        code = [f"package {pkg}", ""]

        code.append("import (")
        code.append('    "fmt"')
        code.append('    "regexp"')
        if exe:
            code.append('    "os"')
            code.append('    "encoding/json"')
            code.append('    "flag"')  # Added for argument parsing
        # Import your local runtime
        code.append('    jm "jsonmodel/json_model/runtime/go/jsonmodel"')
        code.append(")")

        code += ["", "type Checker func(interface{}, *jm.Path, *jm.Report) bool", ""]
        return code

    # --- Initialization and Liberation handling ---

    def gen_init(self, init: Block) -> Block:
        """
        Generates the global initialization function.
        Uses the 'go_init.go' file template to structure the code.
        """
        return self.file_subs("go_init.go", init)

    def gen_free(self, free: Block) -> Block:
        """
        Generates the cleanup function (free/cleanup).
        In Go, this is less critical than in C, but useful for resetting singletons.
        """
        return self.file_subs("go_free.go", free)

    def gen_code(self, code: Block, entry: str, package: str | None, indent: bool = False) -> Block:
        """Replaces global markers in the generated code."""
        if indent:
            code = self.indent(code, False)

        # Ensure the default package is 'main' if not provided
        pkg_name = package if package else "main"

        return [
            line.replace("CHECK_FUNCTION_NAME", entry).replace("CHECK_PACKAGE_NAME", pkg_name)
            for line in code
        ]

    def gen_full_code(
        self,
        defs: Block,
        inis: Block,
        dels: Block,
        subs: Block,
        entry: str,
        package: str | None,
        exe: bool,
    ) -> Block:
        """Generates the final assembled code for Go."""

        # 1. Header generation (Package + Imports)
        # Use provided 'package' name or default to 'main'
        full_code = self.file_header(exe, package)

        # 2. Add global definitions (Variables, Regex)
        full_code += defs
        full_code += [""]

        # 3. Add verification functions (Sub-functions)
        # Process markers like CHECK_FUNCTION_NAME in generated functions
        processed_subs = self.gen_code(subs, entry, package)
        full_code += processed_subs
        full_code += [""]

        # 4. Inject initialization and cleanup via templates
        # These methods use 'file_subs' to replace CODE_BLOCK
        full_code += self.gen_init(inis)
        full_code += [""]
        full_code += self.gen_free(dels)

        # 5. Add main function for the CLI pipeline
        if exe:
            full_code += self.gen_main_function(entry)

        return full_code

    def gen_main_function(self, entry: str) -> Block:
        return [
            "",
            "func main() {",
            '    testMode := flag.Bool("t", false, "run in test mode with a list of [[expected, val], ...]")',
            "    flag.Parse()",
            "    args := flag.Args()",
            "",
            "    if len(args) < 1 {",
            '        fmt.Printf("Usage: %s [-t] <json_file>\\n", os.Args[0])',
            "        os.Exit(1)",
            "    }",
            "",
            "    data, err := os.ReadFile(args[0])",
            "    if err != nil {",
            '        fmt.Printf("Error reading file: %v\\n", err)',
            "        os.Exit(1)",
            "    }",
            "",
            "    check_model_init()",
            "    defer check_model_free()",
            "",
            "    if *testMode {",
            "        var cases [][]interface{}",
            "        if err := json.Unmarshal(data, &cases); err != nil {",
            '            fmt.Printf("Test suite JSON error: %v\\n", err)',
            "            os.Exit(1)",
            "        }",
            "",
            "        failed := 0",
            "        for i, c := range cases {",
            "            if len(c) < 2 { continue }",
            "            expected, _ := c[0].(bool)",
            "            input := c[1]",
            "",
            "            report := &jm.Report{}",
            '            isValid := check_model_map[""](input, nil, report)',
            "",
            "            if isValid == expected {",
            '                fmt.Printf("Test #%d: PASS\\n", i)',
            "            } else {",
            '                fmt.Printf("Test #%d: FAIL (expected %v, got %v)\\n", i, expected, isValid)',
            "                for _, errMsg := range report.Errors {",
            '                    fmt.Printf("  - %s\\n", errMsg)',
            "                }",
            "                failed++",
            "            }",
            "        }",
            "        if failed > 0 {",
            '            fmt.Printf("\\nDone: %d tests failed\\n", failed)',
            "            os.Exit(1)",
            "        }",
            '        fmt.Println("\\nDone: All tests passed")',
            "    } else {",
            "        // Standard single-file validation logic",
            "        var input interface{}",
            "        if err := json.Unmarshal(data, &input); err != nil {",
            '            fmt.Printf("JSON parsing error: %v\\n", err)',
            "            os.Exit(1)",
            "        }",
            "        report := &jm.Report{}",
            '        isValid := check_model_map[""](input, nil, report)',
            "        if isValid {",
            '            fmt.Println("PASS")',
            "        } else {",
            '            fmt.Println("FAIL")',
            "            for _, errMsg := range report.Errors {",
            '                fmt.Printf("  - %s\\n", errMsg)',
            "            }",
            "            os.Exit(1)",
            "        }",
            "    }",
            "}",
        ]
