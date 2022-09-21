package main

import (
	"encoding/json"
	"time"
)

var fieldsExcludedFromHeader = []string{
	"GoFiles",
	"CompiledGoFiles",
	"IgnoredGoFiles",
	"IgnoredOtherFiles",
	"EmbedPatterns",
	"EmbedFiles",
	"CFiles",
	"CgoFiles",
	"CXXFiles",
	"MFiles",
	"HFiles",
	"FFiles",
	"SFiles",
	"SwigFiles",
	"SwigCXXFiles",
	"SysoFiles",
	"Imports",
	"ImportMap",
	"Deps",
	"Module",
	"TestGoFiles",
	"TestImports",
	"XTestGoFiles",
	"XTestImports",
	"ForTest",
}

type jsonPackage struct {
	ImportPath        string            `json:",omitempty"`
	Dir               string            `json:",omitempty"`
	Doc               string            `json:",omitempty"`
	Name              string            `json:",omitempty"`
	Export            string            `json:",omitempty"`
	GoFiles           []string          `json:",omitempty"`
	CompiledGoFiles   []string          `json:",omitempty"`
	IgnoredGoFiles    []string          `json:",omitempty"`
	IgnoredOtherFiles []string          `json:",omitempty"`
	EmbedPatterns     []string          `json:",omitempty"`
	EmbedFiles        []string          `json:",omitempty"`
	CFiles            []string          `json:",omitempty"`
	CgoFiles          []string          `json:",omitempty"`
	CXXFiles          []string          `json:",omitempty"`
	MFiles            []string          `json:",omitempty"`
	HFiles            []string          `json:",omitempty"`
	FFiles            []string          `json:",omitempty"`
	SFiles            []string          `json:",omitempty"`
	SwigFiles         []string          `json:",omitempty"`
	SwigCXXFiles      []string          `json:",omitempty"`
	SysoFiles         []string          `json:",omitempty"`
	Target            string            `json:",omitempty"`
	Imports           []string          `json:",omitempty"`
	ImportMap         map[string]string `json:",omitempty"`
	Deps              []string          `json:",omitempty"`
	Module            *Module           `json:",omitempty"`
	TestGoFiles       []string          `json:",omitempty"`
	TestImports       []string          `json:",omitempty"`
	XTestGoFiles      []string          `json:",omitempty"`
	XTestImports      []string          `json:",omitempty"`
	ForTest           string            `json:",omitempty"` // q in a "p [q.test]" package, else ""
	DepOnly           bool              `json:",omitempty"`

	Error      *PackageError   `json:",omitempty"`
	DepsErrors []*PackageError `json:",omitempty"`
}

func (jp *jsonPackage) headerText() string {
	// todo: rewrite in a propper way

	type header struct {
		ImportPath string        `json:",omitempty"`
		Name       string        `json:",omitempty"`
		Doc        string        `json:",omitempty"`
		Dir        string        `json:",omitempty"`
		Export     string        `json:",omitempty"`
		ForTest    string        `json:",omitempty"` // q in a "p [q.test]" package, else ""
		DepOnly    bool          `json:",omitempty"`
		Error      *PackageError `json:",omitempty"`
		Target     string        `json:",omitempty"`
	}

	raw_json, _ := json.MarshalIndent(&header{
		ImportPath: jp.ImportPath,
		Name:       jp.Name,
		Doc:        jp.Doc,
		Dir:        jp.Dir,
		Export:     jp.Export,
		ForTest:    jp.ForTest,
		DepOnly:    jp.DepOnly,
		Error:      jp.Error,
	}, "", "  ")

	// remove first and last closures together with newlines
	return string(raw_json[2 : len(raw_json)-2])
}

type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

type Module struct {
	Path       string       // module path
	Query      string       // version query corresponding to this version
	Version    string       // module version
	Versions   []string     // available module versions
	Replace    *Module      // replaced by this module
	Time       *time.Time   // time version was created
	Update     *Module      // available update (with -u)
	Main       bool         // is this the main module?
	Indirect   bool         // module is only indirectly needed by main module
	Dir        string       // directory holding local copy of files, if any
	GoMod      string       // path to go.mod file describing module, if any
	GoVersion  string       // go version used in module
	Retracted  []string     // retraction information, if any (with -retracted or -u)
	Deprecated string       // deprecation message, if any (with -u)
	Error      *ModuleError // error loading module
	Origin     any          // provenance of module
	Reuse      bool         // reuse of old module info is safe
}

type ModuleError struct {
	Err string // the error itself
}
