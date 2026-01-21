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
    StrExpr,
    Var,
)


class Go(Language):
    """Go language Code Generator."""

    def __init__(
        self,
        *,
        relib: str = "regexp",
        debug: bool = False,
        with_report: bool = True,
        with_path: bool = True,
        with_package: bool = True,
        with_predef: bool = True,
    ):
        super().__init__(
            "Go",
            debug=debug,
            relib=relib,
            with_report=with_report,
            with_path=with_path,
            with_package=with_package,
            with_predef=with_predef,
            not_op="!",
            and_op="&&",
            or_op="||",
            lcom="//",
            true="true",
            false="false",
            null="nil",
            check_t="bool",
            json_t="interface{}",
            path_t="*Path",
            float_t="float64",
            str_t="string",
            bool_t="bool",
            int_t="int",
            match_t="bool",
            eoi="",
            set_caps=(type(None), bool, int, float, str),  # type: ignore
        )

        assert relib in ("regexp"), f"only standard regexp supported, not {relib}"

    # --- Expressions ---

    def val(self, variable: Var) -> JsonExpr:
        """Dereference interface{} to actual type is tricky in Go without type assertion.
        Here we assume the variable is already cast or we use a helper.
        For the generated code, we usually expect 'val' to be interface{}.
        """
        return variable

    def is_null(self, var: Var) -> BoolExpr:
        return f"{var} == nil"

    def is_bool(self, var: Var) -> BoolExpr:
        return f"IsBool({var})"

    def is_int(self, var: Var) -> BoolExpr:
        return f"IsInteger({var})"

    def is_float(self, var: Var) -> BoolExpr:
        return f"IsNumber({var})"

    def is_str(self, var: Var) -> BoolExpr:
        return f"IsString({var})"

    def is_list(self, var: Var) -> BoolExpr:
        return f"IsArray({var})"

    def is_obj(self, var: Var) -> BoolExpr:
        return f"IsObject({var})"

    # --- Type Conversions (Helpers expected in runtime) ---

    def to_bool(self, var: Var) -> BoolExpr:
        return f"AsBool({var})"

    def to_int(self, var: Var) -> IntExpr:
        return f"AsInt({var})"

    def to_float(self, var: Var) -> Expr:
        return f"AsFloat({var})"

    def to_str(self, var: Var) -> StrExpr:
        return f"AsString({var})"

    def to_list(self, var: Var) -> Expr:
        return f"AsArray({var})"

    def to_obj(self, var: Var) -> Expr:
        return f"AsObject({var})"

    # --- Structure ---

    def header(self, name: str, imports: list[str]) -> Block:
        code = [f"package {name or 'main'}", ""]

        # Standard imports
        std_imports = ['"fmt"']
        if self._relib == "regexp":
            std_imports.append('"regexp"')

        # Add provided imports if any
        all_imports = std_imports + [f'"{i}"' for i in imports]

        code.append("import (")
        for i in sorted(set(all_imports)):
            code.append(f"    {i}")
        code.append(")")
        return code

    def footer(self) -> Block:
        return []

    # --- Function Definition ---

    def sub_fun(self, name: str, body: Block, inline: bool = False) -> Block:
        # func Name(val interface{}, path *Path, rep *Report) bool
        return (
            [f"func {name}(val interface{{}}, path *Path, rep *Report) bool {{"]
            + self.indent(body)
            + ["}"]
        )

    def call_fun(self, name: str, val: Var, path: PathExpr, report: Var) -> BoolExpr:
        return f"{name}({val}, {path}, {report})"

    # --- Regular Expressions ---

    def def_re(self, name: str, regex: str, opts: str) -> Block:
        # Go regexp doesn't support many flags like Perl/Python (e.g. standard library is simple)
        # We usually compile them at package level (var) or init
        return [f"var {name}_re = regexp.MustCompile(`{regex}`)"]

    def match_re(self, name: str, var: Var, regex: str, opts: str) -> BoolExpr:
        # var must be string here
        return f"{name}_re.MatchString({var})"

    # --- Collections (Maps/Sets for constants) ---

    def def_cset(self, name: str, constants: ConstList) -> Block:
        # Define a map[type]bool for O(1) lookup.
        return (
            [f"var {name}_set = map[interface{{}}]bool{{"]
            + [f"    {self.json_cst(c)}: true," for c in constants]
            + ["}"]
        )

    def in_cset(self, name: str, var: Var, constants: ConstList) -> BoolExpr:
        return f"{name}_set[{var}]"

    # --- Constant Maps (Used for property optimization) ---

    def def_cmap(self, name: str, mapping: dict[JsonScalar, str]) -> Block:
        # A map from constant values (usually property names) to check functions.
        # In Go: map[interface{}]func(any, *Path, *Report) bool
        return (
            [f"var {name} = map[interface{{}}]func(interface{{}}, *Path, *Report) bool{{"]
            + [f"    {self.json_cst(k)}: {v}," for k, v in mapping.items()]
            + ["}"]
        )

    def get_cmap(self, name: str, tag: Var, ttag: type) -> Expr:
        return f"{name}[{tag}]"

    def ini_cmap(self, name: str, mapping: dict[JsonScalar, str]) -> Block:
        return []

    def del_cmap(self, name: str, mapping: dict[JsonScalar, str]) -> Block:
        return []

    # --- Utilities ---

    def length(self, var: Var) -> IntExpr:
        return f"len({var})"

    # Go specific formatting for constants
    def json_cst(self, c: JsonScalar) -> str:
        if c is None:
            return "nil"
        if isinstance(c, bool):
            return "true" if c else "false"
        if isinstance(c, str):
            return json.dumps(c)  # Handles escaping quotes, newlines, etc.
        return str(c)
