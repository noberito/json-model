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
        """Déclare et/ou assigne une variable en Go."""
        if tname:
            # Déclaration explicite : var nom type = valeur
            assign = f" = {val}" if val is not None else ""
            return [f"var {var} {tname}{assign}"]
        else:
            # Assignation simple : nom = valeur
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
                # En Go, on vérifie si c'est un float sans partie fractionnaire
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
    # Predefs
    #
    def predef(self, var: Var, name: str, path: Var, is_str: bool = False) -> BoolExpr:
        # Si les prédefs sont désactivées, on valide juste le type string
        if not self._with_predef and self.str_content_predef(name):
            return self.const(True) if is_str else self.is_a(var, str)

        # Préparation du préfixe de type et de la conversion en string
        prefix = "" if is_str else f"jm.IsString({var}) && "
        val = var if is_str else f"jm.AsString({var})"

        # Mapping des prédefs vers le runtime Go
        rt_funcs = {
            "$UUID": "IsValidUUID",
            "$DATE": "IsValidDate",  # Appelera jm.IsValidDate(val)
            "$TIME": "IsValidTime",
            "$DATETIME": "IsValidDateTime",
            "$REGEX": "IsValidRegex",
            "$URL": "IsValidURL",
            "$URI": "IsValidURI",
            "$EMAIL": "IsValidEmail",
            "$JSON": "IsValidJSON",
        }

        if name in rt_funcs:
            # FIX: On retire {path} et rep pour n'avoir qu'un seul argument
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
        # Go ne supporte pas l'assignation dans une expression de comparaison de la même manière que Python (:=)
        # Sauf dans le cadre d'un 'if'. Cette fonction est souvent utilisée par l'optimiseur.
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
        sseg = pseg if is_var else self.esc(pseg) if is_prop else f"jm.IntToSeg({pseg})"
        return f"jm.ExtendPath({pvar}, {sseg})"

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
        # Importation de votre runtime local
        code.append('    jm "jsonmodel/json_model/runtime/go/jsonmodel"')
        code.append(")")

        code += ["", "type Checker func(interface{}, *jm.Path, *jm.Report) bool", ""]
        return code

    # --- Gestion de l'initialisation et de la libération ---

    def gen_init(self, init: Block) -> Block:
        """
        Génère la fonction d'initialisation globale.
        Elle utilise un modèle de fichier 'go_init.go' pour structurer le code.
        """
        return self.file_subs("go_init.go", init)

    def gen_free(self, free: Block) -> Block:
        """
        Génère la fonction de nettoyage (free/cleanup).
        En Go, cela est moins critique qu'en C, mais utile pour réinitialiser des singletons.
        """
        return self.file_subs("go_free.go", free)

    def gen_code(self, code: Block, entry: str, package: str | None, indent: bool = False) -> Block:
        """Remplace les marqueurs globaux dans le code généré."""
        if indent:
            code = self.indent(code, False)

        # On s'assure que le package par défaut est 'main' s'il n'est pas fourni
        pkg_name = package if package else "main"

        return [
            line.replace("CHECK_FUNCTION_NAME", entry).replace("CHECK_PACKAGE_NAME", pkg_name)
            for line in code
        ]

    # In json_model/go.py

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
        """Génère le code final assemblé pour Go."""

        # 1. Génération de l'entête (Package + Imports)
        # On passe 'package' pour utiliser le nom choisi ou 'main' par défaut
        full_code = self.file_header(exe, package)

        # 2. Ajout des définitions globales (Variables, Regex)
        full_code += defs
        full_code += [""]

        # 3. Ajout des fonctions de vérification (Sub-functions)
        # On traite les marqueurs comme CHECK_FUNCTION_NAME dans les fonctions générées
        processed_subs = self.gen_code(subs, entry, package)
        full_code += processed_subs
        full_code += [""]

        # 4. Injection de l'initialisation et du nettoyage via les templates
        # Ces méthodes utilisent 'file_subs' pour remplacer CODE_BLOCK
        full_code += self.gen_init(inis)
        full_code += [""]
        full_code += self.gen_free(dels)

        # 5. Ajout de la fonction main pour le pipeline CLI
        if exe:
            full_code += self.gen_main_function(entry)

        return full_code

    # Dans json_model/go.py

    def gen_main_function(self, entry: str) -> Block:
        return [
            "",
            "func main() {",
            "    if len(os.Args) < 2 {",
            '        fmt.Printf("Usage: %s <json_file>\\n", os.Args[0])',
            "        os.Exit(1)",
            "    }",
            "",
            "    // Lecture du fichier spécifié en argument",
            "    data, err := os.ReadFile(os.Args[1])",
            "    if err != nil {",
            '        fmt.Printf("Erreur lors de la lecture du fichier: %v\\n", err)',
            "        os.Exit(1)",
            "    }",
            "",
            "    var input interface{}",
            "    if err := json.Unmarshal(data, &input); err != nil {",
            '        fmt.Printf("Erreur de parsing JSON: %v\\n", err)',
            "        os.Exit(1)",
            "    }",
            "",
            "    // Initialisation du validateur (regex, etc.)",
            "    check_model_init()",
            "    report := &jm.Report{}",
            "",
            "    // Appel de la validation via l'entrée par défaut",
            '    isValid := check_model_map[""](input, nil, report)',
            "",
            "    if isValid {",
            '        fmt.Println("✅ Le JSON est valide selon le modèle !")',
            "    } else {",
            '        fmt.Println("❌ Le JSON est invalide :")',
            "        for _, errMsg := range report.Errors {",
            '            fmt.Printf("  - %s\\n", errMsg)',
            "        }",
            "        os.Exit(1)",
            "    }",
            "}",
        ]
