package mack

import (
	"strings"
)

type LanguageSpec struct {
	Name string
	Exts []string
}

// borrowed from https://github.com/ggreer/the_silver_searcher/blob/a61f1780b64266587e7bc30f0f5f71c6cca97c0f/src/lang.c
var lang_specs = []LanguageSpec{
	{Name: "actionscript", Exts: []string{"as", "mxml"}},
	{Name: "ada", Exts: []string{"ada", "adb", "ads"}},
	{Name: "asciidoc", Exts: []string{"adoc", "ad", "asc", "asciidoc"}},
	{Name: "apl", Exts: []string{"apl"}},
	{Name: "asm", Exts: []string{"asm", "s"}},
	{Name: "asp", Exts: []string{"asp", "asa", "aspx", "asax", "ashx", "ascx", "asmx"}},
	{Name: "aspx", Exts: []string{"asp", "asa", "aspx", "asax", "ashx", "ascx", "asmx"}},
	{Name: "batch", Exts: []string{"bat", "cmd"}},
	{Name: "bazel", Exts: []string{"bazel"}},
	{Name: "bitbake", Exts: []string{"bb", "bbappend", "bbclass", "inc"}},
	{Name: "cc", Exts: []string{"c", "h", "xs"}},
	{Name: "cfmx", Exts: []string{"cfc", "cfm", "cfml"}},
	{Name: "chpl", Exts: []string{"chpl"}},
	{Name: "clojure", Exts: []string{"clj", "cljs", "cljc", "cljx", "edn"}},
	{Name: "coffee", Exts: []string{"coffee", "cjsx"}},
	{Name: "config", Exts: []string{"config"}},
	{Name: "coq", Exts: []string{"coq", "g", "v"}},
	{Name: "cpp", Exts: []string{"cpp", "cc", "C", "cxx", "m", "hpp", "hh", "h", "H", "hxx", "tpp"}},
	{Name: "crystal", Exts: []string{"cr", "ecr"}},
	{Name: "csharp", Exts: []string{"cs"}},
	{Name: "cshtml", Exts: []string{"cshtml"}},
	{Name: "css", Exts: []string{"css"}},
	{Name: "cython", Exts: []string{"pyx", "pxd", "pxi"}},
	{Name: "delphi", Exts: []string{"pas", "int", "dfm", "nfm", "dof", "dpk", "dpr", "dproj", "groupproj", "bdsgroup", "bdsproj"}},
	{Name: "dlang", Exts: []string{"d", "di"}},
	{Name: "dot", Exts: []string{"dot", "gv"}},
	{Name: "dts", Exts: []string{"dts", "dtsi"}},
	{Name: "ebuild", Exts: []string{"ebuild", "eclass"}},
	{Name: "elisp", Exts: []string{"el"}},
	{Name: "elixir", Exts: []string{"ex", "eex", "exs"}},
	{Name: "elm", Exts: []string{"elm"}},
	{Name: "erlang", Exts: []string{"erl", "hrl"}},
	{Name: "factor", Exts: []string{"factor"}},
	{Name: "fortran", Exts: []string{"f", "F", "f77", "f90", "F90", "f95", "f03", "for", "ftn", "fpp", "FPP"}},
	{Name: "fsharp", Exts: []string{"fs", "fsi", "fsx"}},
	{Name: "gettext", Exts: []string{"po", "pot", "mo"}},
	{Name: "glsl", Exts: []string{"vert", "tesc", "tese", "geom", "frag", "comp"}},
	{Name: "go", Exts: []string{"go"}},
	{Name: "gradle", Exts: []string{"gradle"}},
	{Name: "groovy", Exts: []string{"groovy", "gtmpl", "gpp", "grunit", "gradle"}},
	{Name: "haml", Exts: []string{"haml"}},
	{Name: "handlebars", Exts: []string{"hbs"}},
	{Name: "haskell", Exts: []string{"hs", "hsig", "lhs"}},
	{Name: "haxe", Exts: []string{"hx"}},
	{Name: "hh", Exts: []string{"h"}},
	{Name: "html", Exts: []string{"htm", "html", "shtml", "xhtml"}},
	{Name: "idris", Exts: []string{"idr", "ipkg", "lidr"}},
	{Name: "ini", Exts: []string{"ini"}},
	{Name: "ipython", Exts: []string{"ipynb"}},
	{Name: "isabelle", Exts: []string{"thy"}},
	{Name: "j", Exts: []string{"ijs"}},
	{Name: "jade", Exts: []string{"jade"}},
	{Name: "java", Exts: []string{"java", "properties"}},
	{Name: "jinja2", Exts: []string{"j2"}},
	{Name: "js", Exts: []string{"es6", "js", "jsx", "vue"}},
	{Name: "json", Exts: []string{"json"}},
	{Name: "jsp", Exts: []string{"jsp", "jspx", "jhtm", "jhtml", "jspf", "tag", "tagf"}},
	{Name: "julia", Exts: []string{"jl"}},
	{Name: "kotlin", Exts: []string{"kt"}},
	{Name: "less", Exts: []string{"less"}},
	{Name: "liquid", Exts: []string{"liquid"}},
	{Name: "lisp", Exts: []string{"lisp", "lsp"}},
	{Name: "log", Exts: []string{"log"}},
	{Name: "lua", Exts: []string{"lua"}},
	{Name: "m4", Exts: []string{"m4"}},
	{Name: "make", Exts: []string{"Makefiles", "mk", "mak"}},
	{Name: "mako", Exts: []string{"mako"}},
	{Name: "markdown", Exts: []string{"markdown", "mdown", "mdwn", "mkdn", "mkd", "md"}},
	{Name: "mason", Exts: []string{"mas", "mhtml", "mpl", "mtxt"}},
	{Name: "matlab", Exts: []string{"m"}},
	{Name: "mathematica", Exts: []string{"m", "wl"}},
	{Name: "md", Exts: []string{"markdown", "mdown", "mdwn", "mkdn", "mkd", "md"}},
	{Name: "mercury", Exts: []string{"m", "moo"}},
	{Name: "naccess", Exts: []string{"asa", "rsa"}},
	{Name: "nim", Exts: []string{"nim"}},
	{Name: "nix", Exts: []string{"nix"}},
	{Name: "objc", Exts: []string{"m", "h"}},
	{Name: "objcpp", Exts: []string{"mm", "h"}},
	{Name: "ocaml", Exts: []string{"ml", "mli", "mll", "mly"}},
	{Name: "octave", Exts: []string{"m"}},
	{Name: "org", Exts: []string{"org"}},
	{Name: "parrot", Exts: []string{"pir", "pasm", "pmc", "ops", "pod", "pg", "tg"}},
	{Name: "pdb", Exts: []string{"pdb"}},
	{Name: "perl", Exts: []string{"pl", "pm", "pm6", "pod", "t"}},
	{Name: "php", Exts: []string{"php", "phpt", "php3", "php4", "php5", "phtml"}},
	{Name: "pike", Exts: []string{"pike", "pmod"}},
	{Name: "plist", Exts: []string{"plist"}},
	{Name: "plone", Exts: []string{"pt", "cpt", "metadata", "cpy", "py", "xml", "zcml"}},
	{Name: "powershell", Exts: []string{"ps1"}},
	{Name: "proto", Exts: []string{"proto"}},
	{Name: "ps1", Exts: []string{"ps1"}},
	{Name: "pug", Exts: []string{"pug"}},
	{Name: "puppet", Exts: []string{"pp"}},
	{Name: "python", Exts: []string{"py"}},
	{Name: "qml", Exts: []string{"qml"}},
	{Name: "racket", Exts: []string{"rkt", "ss", "scm"}},
	{Name: "rake", Exts: []string{"Rakefile"}},
	{Name: "razor", Exts: []string{"cshtml"}},
	{Name: "restructuredtext", Exts: []string{"rst"}},
	{Name: "rs", Exts: []string{"rs"}},
	{Name: "r", Exts: []string{"r", "R", "Rmd", "Rnw", "Rtex", "Rrst"}},
	{Name: "rdoc", Exts: []string{"rdoc"}},
	{Name: "ruby", Exts: []string{"rb", "rhtml", "rjs", "rxml", "erb", "rake", "spec"}},
	{Name: "rust", Exts: []string{"rs"}},
	{Name: "salt", Exts: []string{"sls"}},
	{Name: "sass", Exts: []string{"sass", "scss"}},
	{Name: "scala", Exts: []string{"scala"}},
	{Name: "scheme", Exts: []string{"scm", "ss"}},
	{Name: "shell", Exts: []string{"sh", "bash", "csh", "tcsh", "ksh", "zsh", "fish"}},
	{Name: "smalltalk", Exts: []string{"st"}},
	{Name: "sml", Exts: []string{"sml", "fun", "mlb", "sig"}},
	{Name: "sql", Exts: []string{"sql", "ctl"}},
	{Name: "stata", Exts: []string{"do", "ado"}},
	{Name: "stylus", Exts: []string{"styl"}},
	{Name: "swift", Exts: []string{"swift"}},
	{Name: "tcl", Exts: []string{"tcl", "itcl", "itk"}},
	{Name: "terraform", Exts: []string{"tf", "tfvars"}},
	{Name: "tex", Exts: []string{"tex", "cls", "sty"}},
	{Name: "thrift", Exts: []string{"thrift"}},
	{Name: "tla", Exts: []string{"tla"}},
	{Name: "tt", Exts: []string{"tt", "tt2", "ttml"}},
	{Name: "toml", Exts: []string{"toml"}},
	{Name: "ts", Exts: []string{"ts", "tsx"}},
	{Name: "twig", Exts: []string{"twig"}},
	{Name: "vala", Exts: []string{"vala", "vapi"}},
	{Name: "vb", Exts: []string{"bas", "cls", "frm", "ctl", "vb", "resx"}},
	{Name: "velocity", Exts: []string{"vm", "vtl", "vsl"}},
	{Name: "verilog", Exts: []string{"v", "vh", "sv", "svh"}},
	{Name: "vhdl", Exts: []string{"vhd", "vhdl"}},
	{Name: "vim", Exts: []string{"vim"}},
	{Name: "vue", Exts: []string{"vue"}},
	{Name: "wix", Exts: []string{"wxi", "wxs"}},
	{Name: "wsdl", Exts: []string{"wsdl"}},
	{Name: "wadl", Exts: []string{"wadl"}},
	{Name: "xml", Exts: []string{"xml", "dtd", "xsl", "xslt", "xsd", "ent", "tld", "plist", "wsdl"}},
	{Name: "yaml", Exts: []string{"yaml", "yml"}},
	{Name: "zeek", Exts: []string{"zeek", "bro", "bif"}},
	{Name: "zephir", Exts: []string{"zep"}},
}

var known_languages = make(map[string]LanguageSpec)
var known_extentions = make([]string, len(lang_specs))

func InitFileTypes() {
	// initialize from lang_specs
	ext_set := make(map[string]bool)

	for _, spec := range lang_specs {
		known_languages[spec.Name] = spec
		for _, ext := range spec.Exts {
			ext_set[ext] = true
		}
	}
	for ext := range ext_set {
		known_extentions = append(known_extentions, ext)
	}
}

func uniq_exts_from_file_types(types []string) []string {
	// got file types, gather all extentions
	ext_set := make(map[string]bool)

	for _, ftype := range types {
		if spec, ok := known_languages[ftype]; ok {
			for _, ext := range spec.Exts {
				ext_set[ext] = true
			}
		}
	}

	extentions := make([]string, len(ext_set))
	i := 0
	for key := range opts.FileTypeOption.ExtSet {
		extentions[i] = key
		i++
	}

	return extentions
}

func regex_from_file_exts(exts []string) string {
	builder := strings.Builder{}
	builder.WriteString(".*\\.(")
	for _, ext := range exts {
		builder.WriteString(ext)
		builder.WriteString("|")
	}
	regex_str := builder.String()            // convert to string
	regex_str = regex_str[:len(regex_str)-1] // strip trailing '|'
	regex_str += ")$"                        // add the close
	return regex_str
}
