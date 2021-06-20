// This file contains the layouts that are used for testing pkg, including the struct's that shape
// the data for each component of the package (and the package as a whole) and the fullly structured
// data for each of the test packages.
package pkg_test

// testPackage is the main structure for how a package should look after assembly.
type testPackage struct {
	importPath  string
	name        string
	files       []string
	testFiles   []string
	imports     []string
	testImports []string
	functions   []testFunction
	types       []testType
	constants   []string
	variables   []string
	errors      []string
}

type testFunction struct {
	name     string
	comments string
	inputs   []testParameter
	outputs  []testParameter
}

type testType struct {
	name      string
	typeName  string
	source    string
	comments  string
	functions []testFunction
	methods   []testMethod
}

type testMethod struct {
	name        string
	comments    string
	receiver    string
	pointerRcvr bool
	inputs      []testParameter
	outputs     []testParameter
}

type testParameter struct {
	name     string
	typeName string
}

// Structure for package "errors".
var pkgErrors = testPackage{
	importPath:  "errors",
	name:        "errors",
	files:       []string{"errors.go", "wrap.go"},
	testFiles:   []string{"errors_test.go", "example_test.go", "wrap_test.go"},
	imports:     []string{"internal/reflectlite"},
	testImports: []string{"errors", "fmt", "io/fs", "os", "reflect", "testing", "time"},
	functions: []testFunction{
		{
			name:     "As",
			comments: "",
			inputs: []testParameter{
				{
					name:     "err",
					typeName: "error",
				},
				{
					name:     "target",
					typeName: "interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "Is",
			comments: "",
			inputs: []testParameter{
				{
					name:     "err",
					typeName: "error",
				},
				{
					name:     "target",
					typeName: "error",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "New",
			comments: "",
			inputs: []testParameter{
				{
					name:     "text",
					typeName: "string",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
		{
			name:     "Unwrap",
			comments: "",
			inputs: []testParameter{
				{
					name:     "err",
					typeName: "error",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
	},
	types:     []testType{}, // no types in this package
	constants: []string{},   // TODO
	variables: []string{},   // TODO
	errors:    []string{},   // TODO
}

// Structure for package "fmt".
var pkgFmt = testPackage{
	importPath: "fmt",
	name:       "fmt",
	files: []string{
		"doc.go", "errors.go", "format.go", "print.go", "scan.go",
	},
	testFiles: []string{
		"errors_test.go", "example_test.go", "export_test.go", "fmt_test.go", "gostringer_example_test.go",
		"scan_test.go", "stringer_example_test.go", "stringer_test.go",
	},
	imports: []string{
		"errors", "internal/fmtsort", "io", "math", "os", "reflect", "strconv", "sync", "unicode/utf8",
	},
	testImports: []string{
		"bufio", "bytes", "errors", "fmt", "internal/race", "io", "math", "os", "reflect", "regexp",
		"runtime", "strings", "testing", "testing/iotest", "time", "unicode", "unicode/utf8",
	},
	functions: []testFunction{
		{
			name:     "Errorf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fprint",
			comments: "",
			inputs: []testParameter{
				{
					name:     "w",
					typeName: "io.Writer",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fprintf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "w",
					typeName: "io.Writer",
				},
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fprintln",
			comments: "",
			inputs: []testParameter{
				{
					name:     "w",
					typeName: "io.Writer",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fscan",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "io.Reader",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fscanf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "io.Reader",
				},
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Fscanln",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "io.Reader",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Print",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Printf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Println",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Scan",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Scanf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Scanln",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Sprint",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "string",
				},
			},
		},
		{
			name:     "Sprintf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "string",
				},
			},
		},
		{
			name:     "Sprintln",
			comments: "",
			inputs: []testParameter{
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "string",
				},
			},
		},
		{
			name:     "Sscan",
			comments: "",
			inputs: []testParameter{
				{
					name:     "str",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Sscanf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "str",
					typeName: "string",
				},
				{
					name:     "format",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
		{
			name:     "Sscanln",
			comments: "",
			inputs: []testParameter{
				{
					name:     "str",
					typeName: "string",
				},
				{
					name:     "a",
					typeName: "...interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "n",
					typeName: "int",
				},
				{
					name:     "err",
					typeName: "error",
				},
			},
		},
	},
	types: []testType{
		{
			name:      "Formatter",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "GoStringer",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "ScanState",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Scanner",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "State",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Stringer",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
	},
	constants: []string{}, // TODO
	variables: []string{}, // TODO
	errors:    []string{}, // TODO
}

// Structure for package "hash".
var pkgHash = testPackage{
	importPath: "hash",
	name:       "hash",
	files: []string{
		"hash.go",
	},
	testFiles: []string{
		"example_test.go", "marshal_test.go",
	},
	imports: []string{
		"io",
	},
	testImports: []string{
		"bytes", "crypto/md5", "crypto/sha1", "crypto/sha256", "crypto/sha512", "encoding", "encoding/hex",
		"fmt", "hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "log", "testing",
	},
	functions: []testFunction{}, // no functions in this package
	types: []testType{
		{
			name:      "Hash",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Hash32",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Hash64",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
	},
	constants: []string{}, // TODO
	variables: []string{}, // TODO
	errors:    []string{}, // TODO
}

// Structure for package "archive/tar".
var pkgArchiveTar = testPackage{
	importPath: "archive/tar",
	name:       "tar",
	files: []string{
		"common.go", "format.go", "reader.go", "stat_actime1.go", "stat_actime2.go", "stat_unix.go",
		"strconv.go", "writer.go",
	},
	testFiles: []string{
		"example_test.go", "reader_test.go", "strconv_test.go", "tar_test.go", "writer_test.go",
	},
	imports: []string{
		"bytes", "errors", "fmt", "io", "io/fs", "math", "os/user", "path", "reflect", "runtime",
		"sort", "strconv", "strings", "sync", "syscall", "time",
	},
	testImports: []string{
		"archive/tar", "bytes", "crypto/md5", "encoding/hex", "errors", "fmt", "internal/testenv",
		"io", "io/fs", "log", "math", "os", "path", "path/filepath", "reflect", "sort", "strconv",
		"strings", "testing", "testing/iotest", "time",
	},
	functions: []testFunction{}, // no functions in this package
	types: []testType{
		{
			name:      "Format",
			typeName:  "int",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:        "String",
					comments:    "",
					receiver:    "f Format",
					pointerRcvr: false,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "string",
						},
					},
				},
			},
		},
		{
			name:     "Header",
			typeName: "struct",
			source:   "",
			comments: "",
			functions: []testFunction{
				{
					name:     "FileInfoHeader",
					comments: "",
					inputs: []testParameter{
						{
							name:     "fi",
							typeName: "fs.FileInfo",
						},
						{
							name:     "link",
							typeName: "string",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Header",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name:        "FileInfo",
					comments:    "",
					receiver:    "h *Header",
					pointerRcvr: true,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "fs.FileInfo",
						},
					},
				},
			},
		},
		{
			name:     "Reader",
			typeName: "struct",
			source:   "",
			comments: "",
			functions: []testFunction{
				{
					name:     "NewReader",
					comments: "",
					inputs: []testParameter{
						{
							name:     "r",
							typeName: "io.Reader",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Reader",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name:        "Next",
					comments:    "",
					receiver:    "tr *Reader",
					pointerRcvr: true,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Header",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "Read",
					comments:    "",
					receiver:    "tr *Reader",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "b",
							typeName: "[]byte",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "int",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
			},
		},
		{
			name:     "Writer",
			typeName: "struct",
			source:   "",
			comments: "",
			functions: []testFunction{
				{
					name:     "NewWriter",
					comments: "",
					inputs: []testParameter{
						{
							name:     "w",
							typeName: "io.Writer",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Writer",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name:        "Close",
					comments:    "",
					receiver:    "tw *Writer",
					pointerRcvr: true,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "Flush",
					comments:    "",
					receiver:    "tw *Writer",
					pointerRcvr: true,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "Write",
					comments:    "",
					receiver:    "tw *Writer",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "b",
							typeName: "[]byte",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "int",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "WriteHeader",
					comments:    "",
					receiver:    "tw *Writer",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "hdr",
							typeName: "*Header",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
			},
		},
	},
	constants: []string{}, // TODO
	variables: []string{}, // TODO
	errors:    []string{}, // TODO
}

// Structure for package "unicode".
var pkgUnicode = testPackage{
	importPath: "unicode",
	name:       "unicode",
	files: []string{
		"casetables.go", "digit.go", "graphic.go", "letter.go", "tables.go",
	},
	testFiles: []string{
		"digit_test.go", "example_test.go", "graphic_test.go", "letter_test.go", "script_test.go",
	},
	imports: []string{}, // no imports for this package
	testImports: []string{
		"flag", "fmt", "runtime", "sort", "strings", "testing", "unicode",
	},
	functions: []testFunction{
		{
			name:     "In",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
				{
					name:     "ranges",
					typeName: "...*RangeTable",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "Is",
			comments: "",
			inputs: []testParameter{
				{
					name:     "rangeTab",
					typeName: "*RangeTable",
				},
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsControl",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsDigit",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsGraphic",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsLetter",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsLower",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsMark",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsNumber",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsOneOf",
			comments: "",
			inputs: []testParameter{
				{
					name:     "ranges",
					typeName: "[]*RangeTable",
				},
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsPrint",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsPunct",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsSpace",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsSymbol",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsTitle",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "IsUpper",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "bool",
				},
			},
		},
		{
			name:     "SimpleFold",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "rune",
				},
			},
		},
		{
			name:     "To",
			comments: "",
			inputs: []testParameter{
				{
					name:     "_case",
					typeName: "int",
				},
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "rune",
				},
			},
		},
		{
			name:     "ToLower",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "rune",
				},
			},
		},
		{
			name:     "ToTitle",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "rune",
				},
			},
		},
		{
			name:     "ToUpper",
			comments: "",
			inputs: []testParameter{
				{
					name:     "r",
					typeName: "rune",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "rune",
				},
			},
		},
	},
	types: []testType{
		{
			name:      "CaseRange",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Range16",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Range32",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "RangeTable",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "SpecialCase",
			typeName:  "[]CaseRange",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:        "ToLower",
					comments:    "",
					receiver:    "special SpecialCase",
					pointerRcvr: false,
					inputs: []testParameter{
						{
							name:     "r",
							typeName: "rune",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "rune",
						},
					},
				},
				{
					name:        "ToTitle",
					comments:    "",
					receiver:    "special SpecialCase",
					pointerRcvr: false,
					inputs: []testParameter{
						{
							name:     "r",
							typeName: "rune",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "rune",
						},
					},
				},
				{
					name:        "ToUpper",
					comments:    "",
					receiver:    "special SpecialCase",
					pointerRcvr: false,
					inputs: []testParameter{
						{
							name:     "r",
							typeName: "rune",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "rune",
						},
					},
				},
			},
		},
	},
	constants: []string{}, // TODO
	variables: []string{}, // TODO
	errors:    []string{}, // TODO
}

// Structure for package "net/rpc".
var pkgNetRPC = testPackage{
	importPath: "net/rpc",
	name:       "rpc",
	files: []string{
		"client.go", "debug.go", "server.go",
	},
	testFiles: []string{
		"client_test.go", "server_test.go",
	},
	imports: []string{
		"bufio", "encoding/gob", "errors", "fmt", "go/token", "html/template", "io", "log", "net",
		"net/http", "reflect", "sort", "strings", "sync",
	},
	testImports: []string{
		"errors", "fmt", "io", "log", "net", "net/http/httptest", "reflect", "runtime", "strings",
		"sync", "sync/atomic", "testing", "time",
	},
	functions: []testFunction{
		{
			name:     "Accept",
			comments: "",
			inputs: []testParameter{
				{
					name:     "lis",
					typeName: "net.Listener",
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name:     "HandleHTTP",
			comments: "",
			inputs:   []testParameter{}, // no inputs for this function
			outputs:  []testParameter{}, // no outputs for this function
		},
		{
			name:     "Register",
			comments: "",
			inputs: []testParameter{
				{
					name:     "rcvr",
					typeName: "interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
		{
			name:     "RegisterName",
			comments: "",
			inputs: []testParameter{
				{
					name:     "name",
					typeName: "string",
				},
				{
					name:     "rcvr",
					typeName: "interface{}",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
		{
			name:     "ServeCodec",
			comments: "",
			inputs: []testParameter{
				{
					name:     "codec",
					typeName: "ServerCodec",
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name:     "ServeConn",
			comments: "",
			inputs: []testParameter{
				{
					name:     "conn",
					typeName: "io.ReadWriteCloser",
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name:     "ServeRequest",
			comments: "",
			inputs: []testParameter{
				{
					name:     "codec",
					typeName: "ServerCodec",
				},
			},
			outputs: []testParameter{
				{
					name:     "",
					typeName: "error",
				},
			},
		},
	},
	types: []testType{
		{
			name:      "Call",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Client",
			typeName: "struct",
			source:   "",
			comments: "",
			functions: []testFunction{
				{
					name:     "Dial",
					comments: "",
					inputs: []testParameter{
						{
							name:     "network",
							typeName: "string",
						},
						{
							name:     "address",
							typeName: "string",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Client",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:     "DialHTTP",
					comments: "",
					inputs: []testParameter{
						{
							name:     "network",
							typeName: "string",
						},
						{
							name:     "address",
							typeName: "string",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Client",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:     "DialHTTPPath",
					comments: "",
					inputs: []testParameter{
						{
							name:     "network",
							typeName: "string",
						},
						{
							name:     "address",
							typeName: "string",
						},
						{
							name:     "path",
							typeName: "string",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Client",
						},
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:     "NewClient",
					comments: "",
					inputs: []testParameter{
						{
							name:     "conn",
							typeName: "io.ReadWriteCloser",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Client",
						},
					},
				},
				{
					name:     "NewClientWithCodec",
					comments: "",
					inputs: []testParameter{
						{
							name:     "codec",
							typeName: "ClientCodec",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Client",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name:        "Call",
					comments:    "",
					receiver:    "client *Client",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "serviceMethod",
							typeName: "string",
						},
						{
							name:     "args",
							typeName: "interface{}",
						},
						{
							name:     "reply",
							typeName: "interface{}",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "Close",
					comments:    "",
					receiver:    "client *Client",
					pointerRcvr: true,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "Go",
					comments:    "",
					receiver:    "client *Client",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "serviceMethod",
							typeName: "string",
						},
						{
							name:     "args",
							typeName: "interface{}",
						},
						{
							name:     "reply",
							typeName: "interface{}",
						},
						{
							name:     "done",
							typeName: "chan *Call",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Call",
						},
					},
				},
			},
		},
		{
			name:      "ClientCodec",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Request",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "Response",
			typeName:  "struct",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Server",
			typeName: "struct",
			source:   "",
			comments: "",
			functions: []testFunction{
				{
					name:     "NewServer",
					comments: "",
					inputs:   []testParameter{}, // no inputs for this function
					outputs: []testParameter{
						{
							name:     "",
							typeName: "*Server",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name:        "Accept",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "lis",
							typeName: "net.Listener",
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name:        "HandleHTTP",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "rpcPath",
							typeName: "string",
						},
						{
							name:     "debugPath",
							typeName: "string",
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name:        "Register",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "rcvr",
							typeName: "interface{}",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "RegisterName",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "name",
							typeName: "string",
						},
						{
							name:     "rcvr",
							typeName: "interface{}",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
				{
					name:        "ServeCodec",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "codec",
							typeName: "ServerCodec",
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name:        "ServeConn",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "conn",
							typeName: "io.ReadWriteCloser",
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name:        "ServeHTTP",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "w",
							typeName: "http.ResponseWriter",
						},
						{
							name:     "req",
							typeName: "*http.Request",
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name:        "ServeRequest",
					comments:    "",
					receiver:    "server *Server",
					pointerRcvr: true,
					inputs: []testParameter{
						{
							name:     "codec",
							typeName: "ServerCodec",
						},
					},
					outputs: []testParameter{
						{
							name:     "",
							typeName: "error",
						},
					},
				},
			},
		},
		{
			name:      "ServerCodec",
			typeName:  "interface",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:      "ServerError",
			typeName:  "string",
			source:    "",
			comments:  "",
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:        "Error",
					comments:    "",
					receiver:    "e ServerError",
					pointerRcvr: false,
					inputs:      []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							name:     "",
							typeName: "string",
						},
					},
				},
			},
		},
	},
	constants: []string{}, // TODO
	variables: []string{}, // TODO
	errors:    []string{}, // TODO
}
