// This file contains the layouts that are used for testing pkg, including the struct's that shape
// the data for each component of the package (and the package as a whole) and the fullly structured
// data for each of the test packages.
package pkg_test

// testPackage is the main structure for how a package should look after assembly.
type testPackage struct {
	name           string
	importPath     string
	comments       string
	files          []string
	testFiles      []string
	subdirectories []string
	imports        []string
	testImports    []string
	constantBlocks []testConstantBlock
	variableBlocks []testVariableBlock
	functions      []testFunction
	types          []testType
}

type testConstantBlock struct {
	typeName  string
	comments  string
	source    string
	constants []testConstant
}

type testConstant struct {
	name string
}

type testVariableBlock struct {
	typeName  string
	comments  string
	source    string
	variables []testVariable
	errors    []testError
}

type testVariable struct {
	name string
}

type testError struct {
	name string
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
	name     string
	comments string
	receiver testParameter
	inputs   []testParameter
	outputs  []testParameter
}

type testParameter struct {
	s        string
	name     string
	typeName string
	pointer  bool
}

// Structure for package "errors".
var pkgErrors = testPackage{
	name:       "errors",
	importPath: "errors",
	comments: `Package errors implements functions to manipulate errors.

The New function creates errors whose only content is a text message.

The Unwrap, Is and As functions work on errors that may wrap other errors. An error wraps another error if its type has the method

	Unwrap() error

If e.Unwrap() returns a non-nil error w, then we say that e wraps w.

Unwrap unpacks wrapped errors. If its argument's type has an Unwrap method, it calls the method once. Otherwise, it returns nil.

A simple way to create wrapped errors is to call fmt.Errorf and apply the %w verb to the error argument:

	errors.Unwrap(fmt.Errorf("... %w ...", ..., err, ...))

returns err.

Is unwraps its first argument sequentially looking for an error that matches the second. It reports whether it finds a match. It should be used in preference to simple equality checks:

	if errors.Is(err, fs.ErrExist)

is preferable to

	if err == fs.ErrExist

because the former will succeed if err wraps fs.ErrExist.

As unwraps its first argument sequentially looking for an error that can be assigned to its second argument, which must be a pointer. If it succeeds, it performs the assignment and returns true. Otherwise, it returns false. The form

	var perr *fs.PathError
	if errors.As(err, &perr) {
		fmt.Println(perr.Path)
	}

is preferable to

	if perr, ok := err.(*fs.PathError); ok {
		fmt.Println(perr.Path)
	}

because the former will succeed if err wraps an *fs.PathError.
`,
	files:          []string{"errors.go", "wrap.go"},
	testFiles:      []string{"errors_test.go", "example_test.go", "wrap_test.go"},
	subdirectories: []string{}, // no subdirectories in this package
	imports:        []string{"internal/reflectlite"},
	testImports:    []string{"errors", "fmt", "io/fs", "os", "reflect", "testing", "time"},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variableBlocks: []testVariableBlock{}, // no exported variables or errors in this package
	functions: []testFunction{
		{
			name: "As",
			comments: `As finds the first error in err's chain that matches target, and if so, sets target to that error value and returns true. Otherwise, it returns false.

The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.

An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As(interface{}) bool such that As(target) returns true. In the latter case, the As method is responsible for setting target.

An error type might provide an As method so it can be treated as if it were a different error type.

As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type.
`,
			inputs: []testParameter{
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
				{
					s:        "target interface{}",
					name:     "target",
					typeName: "interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "Is",
			comments: `Is reports whether any error in err's chain matches target.

The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.

An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool such that Is(target) returns true.

An error type might provide an Is method so it can be treated as equivalent to an existing error. For example, if MyError defines

	func (m MyError) Is(target error) bool { return target == fs.ErrExist }

then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for an example in the standard library.
`,
			inputs: []testParameter{
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
				{
					s:        "target error",
					name:     "target",
					typeName: "error",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "New",
			comments: `New returns an error that formats as the given text. Each call to New returns a distinct error value even if the text is identical.
`,
			inputs: []testParameter{
				{
					s:        "text string",
					name:     "text",
					typeName: "string",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Unwrap",
			comments: `Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
`,
			inputs: []testParameter{
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
	},
	types: []testType{}, // no types in this package
}

// Structure for package "fmt".
var pkgFmt = testPackage{
	name:       "fmt",
	importPath: "fmt",
	comments: `Package fmt implements formatted I/O with functions analogous to C's printf and scanf. The format 'verbs' are derived from C's but are simpler.


Printing

The verbs:

General:

	%v	the value in a default format
		when printing structs, the plus flag (%+v) adds field names
	%#v	a Go-syntax representation of the value
	%T	a Go-syntax representation of the type of the value
	%%	a literal percent sign; consumes no value

Boolean:

	%t	the word true or false

Integer:

	%b	base 2
	%c	the character represented by the corresponding Unicode code point
	%d	base 10
	%o	base 8
	%O	base 8 with 0o prefix
	%q	a single-quoted character literal safely escaped with Go syntax.
	%x	base 16, with lower-case letters for a-f
	%X	base 16, with upper-case letters for A-F
	%U	Unicode format: U+1234; same as "U+%04X"

Floating-point and complex constituents:

	%b	decimalless scientific notation with exponent a power of two,
		in the manner of strconv.FormatFloat with the 'b' format,
		e.g. -123456p-78
	%e	scientific notation, e.g. -1.234456e+78
	%E	scientific notation, e.g. -1.234456E+78
	%f	decimal point but no exponent, e.g. 123.456
	%F	synonym for %f
	%g	%e for large exponents, %f otherwise. Precision is discussed below.
	%G	%E for large exponents, %F otherwise
	%x	hexadecimal notation (with decimal power of two exponent), e.g. -0x1.23abcp+20
	%X	upper-case hexadecimal notation, e.g. -0X1.23ABCP+20

String and slice of bytes (treated equivalently with these verbs):

	%s	the uninterpreted bytes of the string or slice
	%q	a double-quoted string safely escaped with Go syntax
	%x	base 16, lower-case, two characters per byte
	%X	base 16, upper-case, two characters per byte

Slice:

	%p	address of 0th element in base 16 notation, with leading 0x

Pointer:

	%p	base 16 notation, with leading 0x
	The %b, %d, %o, %x and %X verbs also work with pointers,
	formatting the value exactly as if it were an integer.

The default format for %v is:

	bool:                    %t
	int, int8 etc.:          %d
	uint, uint8 etc.:        %d, %#x if printed with %#v
	float32, complex64, etc: %g
	string:                  %s
	chan:                    %p
	pointer:                 %p

For compound objects, the elements are printed using these rules, recursively, laid out like this:

	struct:             {field0 field1 ...}
	array, slice:       [elem0 elem1 ...]
	maps:               map[key1:value1 key2:value2 ...]
	pointer to above:   &{}, &[], &map[]

Width is specified by an optional decimal number immediately preceding the verb. If absent, the width is whatever is necessary to represent the value. Precision is specified after the (optional) width by a period followed by a decimal number. If no period is present, a default precision is used. A period with no following number specifies a precision of zero. Examples:

	%f     default width, default precision
	%9f    width 9, default precision
	%.2f   default width, precision 2
	%9.2f  width 9, precision 2
	%9.f   width 9, precision 0

Width and precision are measured in units of Unicode code points, that is, runes. (This differs from C's printf where the units are always measured in bytes.) Either or both of the flags may be replaced with the character '*', causing their values to be obtained from the next operand (preceding the one to format), which must be of type int.

For most values, width is the minimum number of runes to output, padding the formatted form with spaces if necessary.

For strings, byte slices and byte arrays, however, precision limits the length of the input to be formatted (not the size of the output), truncating if necessary. Normally it is measured in runes, but for these types when formatted with the %x or %X format it is measured in bytes.

For floating-point values, width sets the minimum width of the field and precision sets the number of places after the decimal, if appropriate, except that for %g/%G precision sets the maximum number of significant digits (trailing zeros are removed). For example, given 12.345 the format %6.3f prints 12.345 while %.3g prints 12.3. The default precision for %e, %f and %#g is 6; for %g it is the smallest number of digits necessary to identify the value uniquely.

For complex numbers, the width and precision apply to the two components independently and the result is parenthesized, so %f applied to 1.2+3.4i produces (1.200000+3.400000i).

Other flags:

	+	always print a sign for numeric values;
		guarantee ASCII-only output for %q (%+q)
	-	pad with spaces on the right rather than the left (left-justify the field)
	#	alternate format: add leading 0b for binary (%#b), 0 for octal (%#o),
		0x or 0X for hex (%#x or %#X); suppress 0x for %p (%#p);
		for %q, print a raw (backquoted) string if strconv.CanBackquote
		returns true;
		always print a decimal point for %e, %E, %f, %F, %g and %G;
		do not remove trailing zeros for %g and %G;
		write e.g. U+0078 'x' if the character is printable for %U (%#U).
	' '	(space) leave a space for elided sign in numbers (% d);
		put spaces between bytes printing strings or slices in hex (% x, % X)
	0	pad with leading zeros rather than spaces;
		for numbers, this moves the padding after the sign

Flags are ignored by verbs that do not expect them. For example there is no alternate decimal format, so %#d and %d behave identically.

For each Printf-like function, there is also a Print function that takes no format and is equivalent to saying %v for every operand. Another variant Println inserts blanks between operands and appends a newline.

Regardless of the verb, if an operand is an interface value, the internal concrete value is used, not the interface itself. Thus:

	var i interface{} = 23
	fmt.Printf("%v\n", i)

will print 23.

Except when printed using the verbs %T and %p, special formatting considerations apply for operands that implement certain interfaces. In order of application:

1. If the operand is a reflect.Value, the operand is replaced by the concrete value that it holds, and printing continues with the next rule.

2. If an operand implements the Formatter interface, it will be invoked. In this case the interpretation of verbs and flags is controlled by that implementation.

3. If the %v verb is used with the # flag (%#v) and the operand implements the GoStringer interface, that will be invoked.

If the format (which is implicitly %v for Println etc.) is valid for a string (%s %q %v %x %X), the following two rules apply:

4. If an operand implements the error interface, the Error method will be invoked to convert the object to a string, which will then be formatted as required by the verb (if any).

5. If an operand implements method String() string, that method will be invoked to convert the object to a string, which will then be formatted as required by the verb (if any).

For compound operands such as slices and structs, the format applies to the elements of each operand, recursively, not to the operand as a whole. Thus %q will quote each element of a slice of strings, and %6.2f will control formatting for each element of a floating-point array.

However, when printing a byte slice with a string-like verb (%s %q %x %X), it is treated identically to a string, as a single item.

To avoid recursion in cases such as

	type X string
	func (x X) String() string { return Sprintf("<%s>", x) }

convert the value before recurring:

	func (x X) String() string { return Sprintf("<%s>", string(x)) }

Infinite recursion can also be triggered by self-referential data structures, such as a slice that contains itself as an element, if that type has a String method. Such pathologies are rare, however, and the package does not protect against them.

When printing a struct, fmt cannot and therefore does not invoke formatting methods such as Error or String on unexported fields.

Explicit argument indexes:

In Printf, Sprintf, and Fprintf, the default behavior is for each formatting verb to format successive arguments passed in the call. However, the notation [n] immediately before the verb indicates that the nth one-indexed argument is to be formatted instead. The same notation before a '*' for a width or precision selects the argument index holding the value. After processing a bracketed expression [n], subsequent verbs will use arguments n+1, n+2, etc. unless otherwise directed.

For example,

	fmt.Sprintf("%[2]d %[1]d\n", 11, 22)

will yield "22 11", while

	fmt.Sprintf("%[3]*.[2]*[1]f", 12.0, 2, 6)

equivalent to

	fmt.Sprintf("%6.2f", 12.0)

will yield " 12.00". Because an explicit index affects subsequent verbs, this notation can be used to print the same values multiple times by resetting the index for the first argument to be repeated:

	fmt.Sprintf("%d %d %#[1]x %#x", 16, 17)

will yield "16 17 0x10 0x11".

Format errors:

If an invalid argument is given for a verb, such as providing a string to %d, the generated string will contain a description of the problem, as in these examples:

	Wrong type or unknown verb: %!verb(type=value)
		Printf("%d", "hi"):        %!d(string=hi)
	Too many arguments: %!(EXTRA type=value)
		Printf("hi", "guys"):      hi%!(EXTRA string=guys)
	Too few arguments: %!verb(MISSING)
		Printf("hi%d"):            hi%!d(MISSING)
	Non-int for width or precision: %!(BADWIDTH) or %!(BADPREC)
		Printf("%*s", 4.5, "hi"):  %!(BADWIDTH)hi
		Printf("%.*s", 4.5, "hi"): %!(BADPREC)hi
	Invalid or invalid use of argument index: %!(BADINDEX)
		Printf("%*[2]d", 7):       %!d(BADINDEX)
		Printf("%.[2]d", 7):       %!d(BADINDEX)

All errors begin with the string "%!" followed sometimes by a single character (the verb) and end with a parenthesized description.

If an Error or String method triggers a panic when called by a print routine, the fmt package reformats the error message from the panic, decorating it with an indication that it came through the fmt package. For example, if a String method calls panic("bad"), the resulting formatted message will look like

	%!s(PANIC=bad)

The %!s just shows the print verb in use when the failure occurred. If the panic is caused by a nil receiver to an Error or String method, however, the output is the undecorated string, "<nil>".


Scanning

An analogous set of functions scans formatted text to yield values. Scan, Scanf and Scanln read from os.Stdin; Fscan, Fscanf and Fscanln read from a specified io.Reader; Sscan, Sscanf and Sscanln read from an argument string.

Scan, Fscan, Sscan treat newlines in the input as spaces.

Scanln, Fscanln and Sscanln stop scanning at a newline and require that the items be followed by a newline or EOF.

Scanf, Fscanf, and Sscanf parse the arguments according to a format string, analogous to that of Printf. In the text that follows, 'space' means any Unicode whitespace character except newline.

In the format string, a verb introduced by the % character consumes and parses input; these verbs are described in more detail below. A character other than %, space, or newline in the format consumes exactly that input character, which must be present. A newline with zero or more spaces before it in the format string consumes zero or more spaces in the input followed by a single newline or the end of the input. A space following a newline in the format string consumes zero or more spaces in the input. Otherwise, any run of one or more spaces in the format string consumes as many spaces as possible in the input. Unless the run of spaces in the format string appears adjacent to a newline, the run must consume at least one space from the input or find the end of the input.

The handling of spaces and newlines differs from that of C's scanf family: in C, newlines are treated as any other space, and it is never an error when a run of spaces in the format string finds no spaces to consume in the input.

The verbs behave analogously to those of Printf. For example, %x will scan an integer as a hexadecimal number, and %v will scan the default representation format for the value. The Printf verbs %p and %T and the flags # and + are not implemented. For floating-point and complex values, all valid formatting verbs (%b %e %E %f %F %g %G %x %X and %v) are equivalent and accept both decimal and hexadecimal notation (for example: "2.3e+7", "0x4.5p-8") and digit-separating underscores (for example: "3.14159_26535_89793").

Input processed by verbs is implicitly space-delimited: the implementation of every verb except %c starts by discarding leading spaces from the remaining input, and the %s verb (and %v reading into a string) stops consuming input at the first space or newline character.

The familiar base-setting prefixes 0b (binary), 0o and 0 (octal), and 0x (hexadecimal) are accepted when scanning integers without a format or with the %v verb, as are digit-separating underscores.

Width is interpreted in the input text but there is no syntax for scanning with a precision (no %5.2f, just %5f). If width is provided, it applies after leading spaces are trimmed and specifies the maximum number of runes to read to satisfy the verb. For example,

	Sscanf(" 1234567 ", "%5s%d", &s, &i)

will set s to "12345" and i to 67 while

	Sscanf(" 12 34 567 ", "%5s%d", &s, &i)

will set s to "12" and i to 34.

In all the scanning functions, a carriage return followed immediately by a newline is treated as a plain newline (\r\n means the same as \n).

In all the scanning functions, if an operand implements method Scan (that is, it implements the Scanner interface) that method will be used to scan the text for that operand. Also, if the number of arguments scanned is less than the number of arguments provided, an error is returned.

All arguments to be scanned must be either pointers to basic types or implementations of the Scanner interface.

Like Scanf and Fscanf, Sscanf need not consume its entire input. There is no way to recover how much of the input string Sscanf used.

Note: Fscan etc. can read one character (rune) past the input they return, which means that a loop calling a scan routine may skip some of the input. This is usually a problem only when there is no space between input values. If the reader provided to Fscan implements ReadRune, that method will be used to read characters. If the reader also implements UnreadRune, that method will be used to save the character and successive calls will not lose data. To attach ReadRune and UnreadRune methods to a reader without that capability, use bufio.NewReader.
`,
	files: []string{
		"doc.go", "errors.go", "format.go", "print.go", "scan.go",
	},
	testFiles: []string{
		"errors_test.go", "example_test.go", "export_test.go", "fmt_test.go", "gostringer_example_test.go",
		"scan_test.go", "stringer_example_test.go", "stringer_test.go",
	},
	subdirectories: []string{}, // no subdirectories in this package
	imports: []string{
		"errors", "internal/fmtsort", "io", "math", "os", "reflect", "strconv", "sync", "unicode/utf8",
	},
	testImports: []string{
		"bufio", "bytes", "errors", "fmt", "internal/race", "io", "math", "os", "reflect", "regexp",
		"runtime", "strings", "testing", "testing/iotest", "time", "unicode", "unicode/utf8",
	},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variableBlocks: []testVariableBlock{}, // no exported variables or errors in this package
	functions: []testFunction{
		{
			name: "Errorf",
			comments: `Errorf formats according to a format specifier and returns the string as a value that satisfies error.

If the format specifier includes a %w verb with an error operand, the returned error will implement an Unwrap method returning the operand. It is invalid to include more than one %w verb or to supply it with an operand that does not implement the error interface. The %w verb is otherwise a synonym for %v.
`,
			inputs: []testParameter{
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fprint",
			comments: `Fprint formats using the default formats for its operands and writes to w. Spaces are added between operands when neither is a string. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "w io.Writer",
					name:     "w",
					typeName: "io.Writer",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fprintf",
			comments: `Fprintf formats according to a format specifier and writes to w. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "w io.Writer",
					name:     "w",
					typeName: "io.Writer",
					pointer:  false,
				},
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fprintln",
			comments: `Fprintln formats using the default formats for its operands and writes to w. Spaces are always added between operands and a newline is appended. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "w io.Writer",
					name:     "w",
					typeName: "io.Writer",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fscan",
			comments: `Fscan scans text read from r, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
			inputs: []testParameter{
				{
					s:        "r io.Reader",
					name:     "r",
					typeName: "io.Reader",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fscanf",
			comments: `Fscanf scans text read from r, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully parsed. Newlines in the input must match newlines in the format.
`,
			inputs: []testParameter{
				{
					s:        "r io.Reader",
					name:     "r",
					typeName: "io.Reader",
					pointer:  false,
				},
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Fscanln",
			comments: `Fscanln is similar to Fscan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
			inputs: []testParameter{
				{
					s:        "r io.Reader",
					name:     "r",
					typeName: "io.Reader",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Print",
			comments: `Print formats using the default formats for its operands and writes to standard output. Spaces are added between operands when neither is a string. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Printf",
			comments: `Printf formats according to a format specifier and writes to standard output. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Println",
			comments: `Println formats using the default formats for its operands and writes to standard output. Spaces are always added between operands and a newline is appended. It returns the number of bytes written and any write error encountered.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Scan",
			comments: `Scan scans text read from standard input, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Scanf",
			comments: `Scanf scans text read from standard input, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why. Newlines in the input must match newlines in the format. The one exception: the verb %c always scans the next rune in the input, even if it is a space (or tab etc.) or newline.
`,
			inputs: []testParameter{
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Scanln",
			comments: `Scanln is similar to Scan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Sprint",
			comments: `Sprint formats using the default formats for its operands and returns the resulting string. Spaces are added between operands when neither is a string.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "string",
					name:     "",
					typeName: "string",
					pointer:  false,
				},
			},
		},
		{
			name: "Sprintf",
			comments: `Sprintf formats according to a format specifier and returns the resulting string.
`,
			inputs: []testParameter{
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "string",
					name:     "",
					typeName: "string",
					pointer:  false,
				},
			},
		},
		{
			name: "Sprintln",
			comments: `Sprintln formats using the default formats for its operands and returns the resulting string. Spaces are always added between operands and a newline is appended.
`,
			inputs: []testParameter{
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "string",
					name:     "",
					typeName: "string",
					pointer:  false,
				},
			},
		},
		{
			name: "Sscan",
			comments: `Sscan scans the argument string, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
			inputs: []testParameter{
				{
					s:        "str string",
					name:     "str",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Sscanf",
			comments: `Sscanf scans the argument string, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully parsed. Newlines in the input must match newlines in the format.
`,
			inputs: []testParameter{
				{
					s:        "str string",
					name:     "str",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "format string",
					name:     "format",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "Sscanln",
			comments: `Sscanln is similar to Sscan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
			inputs: []testParameter{
				{
					s:        "str string",
					name:     "str",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "a ...interface{}",
					name:     "a",
					typeName: "...interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "n int",
					name:     "n",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "err error",
					name:     "err",
					typeName: "error",
					pointer:  false,
				},
			},
		},
	},
	types: []testType{
		{
			name:     "Formatter",
			typeName: "interface",
			source: `type Formatter interface {
	Format(f State, verb rune)
}`,
			comments: `Formatter is implemented by any value that has a Format method. The implementation controls how State and rune are interpreted, and may call Sprint(f) or Fprint(f) etc. to generate its output.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "GoStringer",
			typeName: "interface",
			source: `type GoStringer interface {
	GoString() string
}`,
			comments: `GoStringer is implemented by any value that has a GoString method, which defines the Go syntax for that value. The GoString method is used to print values passed as an operand to a %#v format.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "ScanState",
			typeName: "interface",
			source: `type ScanState interface {
	// ReadRune reads the next rune (Unicode code point) from the input.
	// If invoked during Scanln, Fscanln, or Sscanln, ReadRune() will
	// return EOF after returning the first '\n' or when reading beyond
	// the specified width.
	ReadRune() (r rune, size int, err error)
	// UnreadRune causes the next call to ReadRune to return the same rune.
	UnreadRune() error
	// SkipSpace skips space in the input. Newlines are treated appropriately
	// for the operation being performed; see the package documentation
	// for more information.
	SkipSpace()
	// Token skips space in the input if skipSpace is true, then returns the
	// run of Unicode code points c satisfying f(c).  If f is nil,
	// !unicode.IsSpace(c) is used; that is, the token will hold non-space
	// characters. Newlines are treated appropriately for the operation being
	// performed; see the package documentation for more information.
	// The returned slice points to shared data that may be overwritten
	// by the next call to Token, a call to a Scan function using the ScanState
	// as input, or when the calling Scan method returns.
	Token(skipSpace bool, f func(rune) bool) (token []byte, err error)
	// Width returns the value of the width option and whether it has been set.
	// The unit is Unicode code points.
	Width() (wid int, ok bool)
	// Because ReadRune is implemented by the interface, Read should never be
	// called by the scanning routines and a valid implementation of
	// ScanState may choose always to return an error from Read.
	Read(buf []byte) (n int, err error)
}`,
			comments: `ScanState represents the scanner state passed to custom scanners. Scanners may do rune-at-a-time scanning or ask the ScanState to discover the next space-delimited token.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Scanner",
			typeName: "interface",
			source: `type Scanner interface {
	Scan(state ScanState, verb rune) error
}`,
			comments: `Scanner is implemented by any value that has a Scan method, which scans the input for the representation of a value and stores the result in the receiver, which must be a pointer to be useful. The Scan method is called for any argument to Scan, Scanf, or Scanln that implements it.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "State",
			typeName: "interface",
			source: `type State interface {
	// Write is the function to call to emit formatted output to be printed.
	Write(b []byte) (n int, err error)
	// Width returns the value of the width option and whether it has been set.
	Width() (wid int, ok bool)
	// Precision returns the value of the precision option and whether it has been set.
	Precision() (prec int, ok bool)

	// Flag reports whether the flag c, a character, has been set.
	Flag(c int) bool
}`,
			comments: `State represents the printer state passed to custom formatters. It provides access to the io.Writer interface plus information about the flags and options for the operand's format specifier.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Stringer",
			typeName: "interface",
			source: `type Stringer interface {
	String() string
}`,
			comments: `Stringer is implemented by any value that has a String method, which defines the “native” format for that value. The String method is used to print values passed as an operand to any format that accepts a string or to an unformatted printer such as Print.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
	},
}

// Structure for package "hash".
var pkgHash = testPackage{
	name:       "hash",
	importPath: "hash",
	comments: `Package hash provides interfaces for hash functions.
`,
	files: []string{
		"hash.go",
	},
	testFiles: []string{
		"example_test.go", "marshal_test.go",
	},
	subdirectories: []string{"adler32", "crc32", "crc64", "fnv", "maphash"}, // no subdirectories in this package
	imports: []string{
		"io",
	},
	testImports: []string{
		"bytes", "crypto/md5", "crypto/sha1", "crypto/sha256", "crypto/sha512", "encoding", "encoding/hex",
		"fmt", "hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "log", "testing",
	},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variableBlocks: []testVariableBlock{}, // no exported variables or errors in this package
	functions:      []testFunction{},      // no functions in this package
	types: []testType{
		{
			name:     "Hash",
			typeName: "interface",
			source: `type Hash interface {
	// Write (via the embedded io.Writer interface) adds more data to the running hash.
	// It never returns an error.
	io.Writer

	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	Sum(b []byte) []byte

	// Reset resets the Hash to its initial state.
	Reset()

	// Size returns the number of bytes Sum will return.
	Size() int

	// BlockSize returns the hash's underlying block size.
	// The Write method must be able to accept any amount
	// of data, but it may operate more efficiently if all writes
	// are a multiple of the block size.
	BlockSize() int
}`,
			comments: `Hash is the common interface implemented by all hash functions.

Hash implementations in the standard library (e.g. hash/crc32 and crypto/sha256) implement the encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interfaces. Marshaling a hash implementation allows its internal state to be saved and used for additional processing later, without having to re-write the data previously written to the hash. The hash state may contain portions of the input in its original form, which users are expected to handle for any possible security implications.

Compatibility: Any future changes to hash or crypto packages will endeavor to maintain compatibility with state encoded using previous versions. That is, any released versions of the packages should be able to decode data written with any previously released version, subject to issues such as security fixes. See the Go compatibility document for background: https://golang.org/doc/go1compat
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Hash32",
			typeName: "interface",
			source: `type Hash32 interface {
	Hash
	Sum32() uint32
}`,
			comments: `Hash32 is the common interface implemented by all 32-bit hash functions.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Hash64",
			typeName: "interface",
			source: `type Hash64 interface {
	Hash
	Sum64() uint64
}`,
			comments: `Hash64 is the common interface implemented by all 64-bit hash functions.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
	},
}

// Structure for package "archive/tar".
var pkgArchiveTar = testPackage{
	name:       "tar",
	importPath: "archive/tar",
	comments: `Package tar implements access to tar archives.

Tape archives (tar) are a file format for storing a sequence of files that can be read and written in a streaming manner. This package aims to cover most variations of the format, including those produced by GNU and BSD tar tools.
`,
	files: []string{
		"common.go", "format.go", "reader.go", "stat_actime1.go", "stat_actime2.go", "stat_unix.go",
		"strconv.go", "writer.go",
	},
	testFiles: []string{
		"example_test.go", "reader_test.go", "strconv_test.go", "tar_test.go", "writer_test.go",
	},
	subdirectories: []string{"testdata"}, // no subdirectories in this package
	imports: []string{
		"bytes", "errors", "fmt", "io", "io/fs", "math", "os/user", "path", "reflect", "runtime",
		"sort", "strconv", "strings", "sync", "syscall", "time",
	},
	testImports: []string{
		"archive/tar", "bytes", "crypto/md5", "encoding/hex", "errors", "fmt", "internal/testenv",
		"io", "io/fs", "log", "math", "os", "path", "path/filepath", "reflect", "sort", "strconv",
		"strings", "testing", "testing/iotest", "time",
	},
	constantBlocks: []testConstantBlock{
		{
			typeName: "", // no general type for this block of constants
			comments: `Type flags for Header.Typeflag.
`,
			source: `const (
	// Type '0' indicates a regular file.
	TypeReg  = '0'
	TypeRegA = '\x00' // Deprecated: Use TypeReg instead.

	// Type '1' to '6' are header-only flags and may not have a data body.
	TypeLink    = '1' // Hard link
	TypeSymlink = '2' // Symbolic link
	TypeChar    = '3' // Character device node
	TypeBlock   = '4' // Block device node
	TypeDir     = '5' // Directory
	TypeFifo    = '6' // FIFO node

	// Type '7' is reserved.
	TypeCont = '7'

	// Type 'x' is used by the PAX format to store key-value records that
	// are only relevant to the next file.
	// This package transparently handles these types.
	TypeXHeader = 'x'

	// Type 'g' is used by the PAX format to store key-value records that
	// are relevant to all subsequent files.
	// This package only supports parsing and composing such headers,
	// but does not currently support persisting the global state across files.
	TypeXGlobalHeader = 'g'

	// Type 'S' indicates a sparse file in the GNU format.
	TypeGNUSparse = 'S'

	// Types 'L' and 'K' are used by the GNU format for a meta file
	// used to store the path or link name for the next file.
	// This package transparently handles these types.
	TypeGNULongName = 'L'
	TypeGNULongLink = 'K'
)`,
			constants: []testConstant{
				{
					name: "TypeReg",
				},
				{
					name: "TypeRegA",
				},
				{
					name: "TypeLink",
				},
				{
					name: "TypeSymlink",
				},
				{
					name: "TypeChar",
				},
				{
					name: "TypeBlock",
				},
				{
					name: "TypeDir",
				},
				{
					name: "TypeFifo",
				},
				{
					name: "TypeCont",
				},
				{
					name: "TypeXHeader",
				},
				{
					name: "TypeXGlobalHeader",
				},
				{
					name: "TypeGNUSparse",
				},
				{
					name: "TypeGNULongName",
				},
				{
					name: "TypeGNULongLink",
				},
			},
		},
		{
			typeName: "Format",
			comments: `Constants to identify various tar formats.
`,
			source: `const (
	// FormatUnknown indicates that the format is unknown.
	FormatUnknown Format

	// FormatUSTAR represents the USTAR header format defined in POSIX.1-1988.
	//
	// While this format is compatible with most tar readers,
	// the format has several limitations making it unsuitable for some usages.
	// Most notably, it cannot support sparse files, files larger than 8GiB,
	// filenames larger than 256 characters, and non-ASCII filenames.
	//
	// Reference:
	//	http://pubs.opengroup.org/onlinepubs/9699919799/utilities/pax.html#tag_20_92_13_06
	FormatUSTAR

	// FormatPAX represents the PAX header format defined in POSIX.1-2001.
	//
	// PAX extends USTAR by writing a special file with Typeflag TypeXHeader
	// preceding the original header. This file contains a set of key-value
	// records, which are used to overcome USTAR's shortcomings, in addition to
	// providing the ability to have sub-second resolution for timestamps.
	//
	// Some newer formats add their own extensions to PAX by defining their
	// own keys and assigning certain semantic meaning to the associated values.
	// For example, sparse file support in PAX is implemented using keys
	// defined by the GNU manual (e.g., "GNU.sparse.map").
	//
	// Reference:
	//	http://pubs.opengroup.org/onlinepubs/009695399/utilities/pax.html
	FormatPAX

	// FormatGNU represents the GNU header format.
	//
	// The GNU header format is older than the USTAR and PAX standards and
	// is not compatible with them. The GNU format supports
	// arbitrary file sizes, filenames of arbitrary encoding and length,
	// sparse files, and other features.
	//
	// It is recommended that PAX be chosen over GNU unless the target
	// application can only parse GNU formatted archives.
	//
	// Reference:
	//	https://www.gnu.org/software/tar/manual/html_node/Standard.html
	FormatGNU
)`,
			constants: []testConstant{
				{
					name: "FormatUnknown",
				},
				{
					name: "FormatUSTAR",
				},
				{
					name: "FormatPAX",
				},
				{
					name: "FormatGNU",
				},
			},
		},
	},
	variableBlocks: []testVariableBlock{
		{
			typeName: "",
			comments: ``, // no comments for this block of variables
			source: `var (
	ErrHeader          = errors.New("archive/tar: invalid tar header")
	ErrWriteTooLong    = errors.New("archive/tar: write too long")
	ErrFieldTooLong    = errors.New("archive/tar: header field too long")
	ErrWriteAfterClose = errors.New("archive/tar: write after close")
)`,
			variables: []testVariable{
				{
					name: "ErrHeader",
				},
				{
					name: "ErrWriteTooLong",
				},
				{
					name: "ErrFieldTooLong",
				},
				{
					name: "ErrWriteAfterClose",
				},
			},
			errors: []testError{
				{
					name: "ErrHeader",
				},
				{
					name: "ErrWriteTooLong",
				},
				{
					name: "ErrFieldTooLong",
				},
				{
					name: "ErrWriteAfterClose",
				},
			},
		},
	},
	functions: []testFunction{}, // no functions in this package
	types: []testType{
		{
			name:     "Format",
			typeName: "int",
			source:   `type Format int`,
			comments: `Format represents the tar archive format.

The original tar format was introduced in Unix V7. Since then, there have been multiple competing formats attempting to standardize or extend the V7 format to overcome its limitations. The most common formats are the USTAR, PAX, and GNU formats, each with their own advantages and limitations.

The following table captures the capabilities of each format:

	                  |  USTAR |       PAX |       GNU
	------------------+--------+-----------+----------
	Name              |   256B | unlimited | unlimited
	Linkname          |   100B | unlimited | unlimited
	Size              | uint33 | unlimited |    uint89
	Mode              | uint21 |    uint21 |    uint57
	Uid/Gid           | uint21 | unlimited |    uint57
	Uname/Gname       |    32B | unlimited |       32B
	ModTime           | uint33 | unlimited |     int89
	AccessTime        |    n/a | unlimited |     int89
	ChangeTime        |    n/a | unlimited |     int89
	Devmajor/Devminor | uint21 |    uint21 |    uint57
	------------------+--------+-----------+----------
	string encoding   |  ASCII |     UTF-8 |    binary
	sub-second times  |     no |       yes |        no
	sparse files      |     no |       yes |       yes

The table's upper portion shows the Header fields, where each format reports the maximum number of bytes allowed for each string field and the integer type used to store each numeric field (where timestamps are stored as the number of seconds since the Unix epoch).

The table's lower portion shows specialized features of each format, such as supported string encodings, support for sub-second timestamps, or support for sparse files.

The Writer currently provides no support for sparse files.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:     "String",
					comments: ``, // no comments for this method
					receiver: testParameter{
						s:        "f Format",
						name:     "f",
						typeName: "Format",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Header",
			typeName: "struct",
			source: `type Header struct {
	// Typeflag is the type of header entry.
	// The zero value is automatically promoted to either TypeReg or TypeDir
	// depending on the presence of a trailing slash in Name.
	Typeflag byte

	Name     string // Name of file entry
	Linkname string // Target name of link (valid for TypeLink or TypeSymlink)

	Size  int64  // Logical file size in bytes
	Mode  int64  // Permission and mode bits
	Uid   int    // User ID of owner
	Gid   int    // Group ID of owner
	Uname string // User name of owner
	Gname string // Group name of owner

	// If the Format is unspecified, then Writer.WriteHeader rounds ModTime
	// to the nearest second and ignores the AccessTime and ChangeTime fields.
	//
	// To use AccessTime or ChangeTime, specify the Format as PAX or GNU.
	// To use sub-second resolution, specify the Format as PAX.
	ModTime    time.Time // Modification time
	AccessTime time.Time // Access time (requires either PAX or GNU support)
	ChangeTime time.Time // Change time (requires either PAX or GNU support)

	Devmajor int64 // Major device number (valid for TypeChar or TypeBlock)
	Devminor int64 // Minor device number (valid for TypeChar or TypeBlock)

	// Xattrs stores extended attributes as PAX records under the
	// "SCHILY.xattr." namespace.
	//
	// The following are semantically equivalent:
	//  h.Xattrs[key] = value
	//  h.PAXRecords["SCHILY.xattr."+key] = value
	//
	// When Writer.WriteHeader is called, the contents of Xattrs will take
	// precedence over those in PAXRecords.
	//
	// Deprecated: Use PAXRecords instead.
	Xattrs map[string]string

	// PAXRecords is a map of PAX extended header records.
	//
	// User-defined records should have keys of the following form:
	//	VENDOR.keyword
	// Where VENDOR is some namespace in all uppercase, and keyword may
	// not contain the '=' character (e.g., "GOLANG.pkg.version").
	// The key and value should be non-empty UTF-8 strings.
	//
	// When Writer.WriteHeader is called, PAX records derived from the
	// other fields in Header take precedence over PAXRecords.
	PAXRecords map[string]string

	// Format specifies the format of the tar header.
	//
	// This is set by Reader.Next as a best-effort guess at the format.
	// Since the Reader liberally reads some non-compliant files,
	// it is possible for this to be FormatUnknown.
	//
	// If the format is unspecified when Writer.WriteHeader is called,
	// then it uses the first format (in the order of USTAR, PAX, GNU)
	// capable of encoding this Header (see Format).
	Format Format
}`,
			comments: `A Header represents a single header in a tar archive. Some fields may not be populated.

For forward compatibility, users that retrieve a Header from Reader.Next, mutate it in some ways, and then pass it back to Writer.WriteHeader should do so by creating a new Header and copying the fields that they are interested in preserving.
`,
			functions: []testFunction{
				{
					name: "FileInfoHeader",
					comments: `FileInfoHeader creates a partially-populated Header from fi. If fi describes a symlink, FileInfoHeader records link as the link target. If fi describes a directory, a slash is appended to the name.

Since fs.FileInfo's Name method only returns the base name of the file it describes, it may be necessary to modify Header.Name to provide the full path name of the file.
`,
					inputs: []testParameter{
						{
							s:        "fi fs.FileInfo",
							name:     "fi",
							typeName: "fs.FileInfo",
							pointer:  false,
						},
						{
							s:        "link string",
							name:     "link",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Header",
							name:     "",
							typeName: "*Header",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "FileInfo",
					comments: `FileInfo returns an fs.FileInfo for the Header.
`,
					receiver: testParameter{
						s:        "h *Header",
						name:     "h",
						typeName: "*Header",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "fs.FileInfo",
							name:     "",
							typeName: "fs.FileInfo",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Reader",
			typeName: "struct",
			source: `type Reader struct {
	// contains filtered or unexported fields
}`,
			comments: `Reader provides sequential access to the contents of a tar archive. Reader.Next advances to the next file in the archive (including the first), and then Reader can be treated as an io.Reader to access the file's data.
`,
			functions: []testFunction{
				{
					name: "NewReader",
					comments: `NewReader creates a new Reader reading from r.
`,
					inputs: []testParameter{
						{
							s:        "r io.Reader",
							name:     "r",
							typeName: "io.Reader",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Reader",
							name:     "",
							typeName: "*Reader",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Next",
					comments: `Next advances to the next entry in the tar archive. The Header.Size determines how many bytes can be read for the next file. Any remaining data in the current file is automatically discarded.

io.EOF is returned at the end of the input.
`,
					receiver: testParameter{
						s:        "tr *Reader",
						name:     "tr",
						typeName: "*Reader",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "*Header",
							name:     "",
							typeName: "*Header",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Read",
					comments: `Read reads from the current file in the tar archive. It returns (0, io.EOF) when it reaches the end of that file, until Next is called to advance to the next file.

If the current file is sparse, then the regions marked as a hole are read back as NUL-bytes.

Calling Read on special types like TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, and TypeFifo returns (0, io.EOF) regardless of what the Header.Size claims.
`,
					receiver: testParameter{
						s:        "tr *Reader",
						name:     "tr",
						typeName: "*Reader",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "b []byte",
							name:     "b",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Writer",
			typeName: "struct",
			source: `type Writer struct {
	// contains filtered or unexported fields
}`,
			comments: `Writer provides sequential writing of a tar archive. Write.WriteHeader begins a new file with the provided Header, and then Writer can be treated as an io.Writer to supply that file's data.
`,
			functions: []testFunction{
				{
					name: "NewWriter",
					comments: `NewWriter creates a new Writer writing to w.
`,
					inputs: []testParameter{
						{
							s:        "w io.Writer",
							name:     "w",
							typeName: "io.Writer",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Writer",
							name:     "",
							typeName: "*Writer",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Close",
					comments: `Close closes the tar archive by flushing the padding, and writing the footer. If the current file (from a prior call to WriteHeader) is not fully written, then this returns an error.
`,
					receiver: testParameter{
						s:        "tw *Writer",
						name:     "tw",
						typeName: "*Writer",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Flush",
					comments: `Flush finishes writing the current file's block padding. The current file must be fully written before Flush can be called.

This is unnecessary as the next call to WriteHeader or Close will implicitly flush out the file's padding.
`,
					receiver: testParameter{
						s:        "tw *Writer",
						name:     "tw",
						typeName: "*Writer",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Write",
					comments: `Write writes to the current file in the tar archive. Write returns the error ErrWriteTooLong if more than Header.Size bytes are written after WriteHeader.

Calling Write on special types like TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, and TypeFifo returns (0, ErrWriteTooLong) regardless of what the Header.Size claims.
`,
					receiver: testParameter{
						s:        "tw *Writer",
						name:     "tw",
						typeName: "*Writer",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "b []byte",
							name:     "b",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "WriteHeader",
					comments: `WriteHeader writes hdr and prepares to accept the file's contents. The Header.Size determines how many bytes can be written for the next file. If the current file is not fully written, then this returns an error. This implicitly flushes any padding necessary before writing the header.
`,
					receiver: testParameter{
						s:        "tw *Writer",
						name:     "tw",
						typeName: "*Writer",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "hdr *Header",
							name:     "hdr",
							typeName: "*Header",
							pointer:  true,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
			},
		},
	},
}

// Structure for package "unicode".
var pkgUnicode = testPackage{
	name:       "unicode",
	importPath: "unicode",
	comments: `Package unicode provides data and functions to test some properties of Unicode code points.
`,
	files: []string{
		"casetables.go", "digit.go", "graphic.go", "letter.go", "tables.go",
	},
	testFiles: []string{
		"digit_test.go", "example_test.go", "graphic_test.go", "letter_test.go", "script_test.go",
	},
	subdirectories: []string{"utf16", "utf8"}, // no subdirectories in this package
	imports:        []string{},                // no imports for this package
	testImports: []string{
		"flag", "fmt", "runtime", "sort", "strings", "testing", "unicode",
	},
	constantBlocks: []testConstantBlock{
		{
			typeName: "", // no general type for this block of constants
			comments: ``, // no comments for this block of constants
			source: `const (
	MaxRune         = '\U0010FFFF' // Maximum valid Unicode code point.
	ReplacementChar = '\uFFFD'     // Represents invalid code points.
	MaxASCII        = '\u007F'     // maximum ASCII value.
	MaxLatin1       = '\u00FF'     // maximum Latin-1 value.
)`,
			constants: []testConstant{
				{
					name: "MaxRune",
				},
				{
					name: "ReplacementChar",
				},
				{
					name: "MaxASCII",
				},
				{
					name: "MaxLatin1",
				},
			},
		},
		{
			typeName: "", // no general type for this block of constants
			comments: `Indices into the Delta arrays inside CaseRanges for case mapping.
`,
			source: `const (
	UpperCase = iota
	LowerCase
	TitleCase
	MaxCase
)`,
			constants: []testConstant{
				{
					name: "UpperCase",
				},
				{
					name: "LowerCase",
				},
				{
					name: "TitleCase",
				},
				{
					name: "MaxCase",
				},
			},
		},
		{
			typeName: "", // no general type for this block of constants
			comments: `If the Delta field of a CaseRange is UpperLower, it means this CaseRange represents a sequence of the form (say) Upper Lower Upper Lower.
`,
			source: `const (
	UpperLower = MaxRune + 1 // (Cannot be a valid delta.)
)`,
			constants: []testConstant{
				{
					name: "UpperLower",
				},
			},
		},
		{
			typeName: "", // no general type for this block of constants
			comments: `Version is the Unicode edition from which the tables are derived.
`,
			source: `const Version = "13.0.0"`,
			constants: []testConstant{
				{
					name: "Version",
				},
			},
		},
	},
	variableBlocks: []testVariableBlock{
		{
			typeName: "",
			comments: `These variables have type *RangeTable.
`,
			source: `var (
	Cc     = _Cc // Cc is the set of Unicode characters in category Cc (Other, control).
	Cf     = _Cf // Cf is the set of Unicode characters in category Cf (Other, format).
	Co     = _Co // Co is the set of Unicode characters in category Co (Other, private use).
	Cs     = _Cs // Cs is the set of Unicode characters in category Cs (Other, surrogate).
	Digit  = _Nd // Digit is the set of Unicode characters with the "decimal digit" property.
	Nd     = _Nd // Nd is the set of Unicode characters in category Nd (Number, decimal digit).
	Letter = _L  // Letter/L is the set of Unicode letters, category L.
	L      = _L
	Lm     = _Lm // Lm is the set of Unicode characters in category Lm (Letter, modifier).
	Lo     = _Lo // Lo is the set of Unicode characters in category Lo (Letter, other).
	Lower  = _Ll // Lower is the set of Unicode lower case letters.
	Ll     = _Ll // Ll is the set of Unicode characters in category Ll (Letter, lowercase).
	Mark   = _M  // Mark/M is the set of Unicode mark characters, category M.
	M      = _M
	Mc     = _Mc // Mc is the set of Unicode characters in category Mc (Mark, spacing combining).
	Me     = _Me // Me is the set of Unicode characters in category Me (Mark, enclosing).
	Mn     = _Mn // Mn is the set of Unicode characters in category Mn (Mark, nonspacing).
	Nl     = _Nl // Nl is the set of Unicode characters in category Nl (Number, letter).
	No     = _No // No is the set of Unicode characters in category No (Number, other).
	Number = _N  // Number/N is the set of Unicode number characters, category N.
	N      = _N
	Other  = _C // Other/C is the set of Unicode control and special characters, category C.
	C      = _C
	Pc     = _Pc // Pc is the set of Unicode characters in category Pc (Punctuation, connector).
	Pd     = _Pd // Pd is the set of Unicode characters in category Pd (Punctuation, dash).
	Pe     = _Pe // Pe is the set of Unicode characters in category Pe (Punctuation, close).
	Pf     = _Pf // Pf is the set of Unicode characters in category Pf (Punctuation, final quote).
	Pi     = _Pi // Pi is the set of Unicode characters in category Pi (Punctuation, initial quote).
	Po     = _Po // Po is the set of Unicode characters in category Po (Punctuation, other).
	Ps     = _Ps // Ps is the set of Unicode characters in category Ps (Punctuation, open).
	Punct  = _P  // Punct/P is the set of Unicode punctuation characters, category P.
	P      = _P
	Sc     = _Sc // Sc is the set of Unicode characters in category Sc (Symbol, currency).
	Sk     = _Sk // Sk is the set of Unicode characters in category Sk (Symbol, modifier).
	Sm     = _Sm // Sm is the set of Unicode characters in category Sm (Symbol, math).
	So     = _So // So is the set of Unicode characters in category So (Symbol, other).
	Space  = _Z  // Space/Z is the set of Unicode space characters, category Z.
	Z      = _Z
	Symbol = _S // Symbol/S is the set of Unicode symbol characters, category S.
	S      = _S
	Title  = _Lt // Title is the set of Unicode title case letters.
	Lt     = _Lt // Lt is the set of Unicode characters in category Lt (Letter, titlecase).
	Upper  = _Lu // Upper is the set of Unicode upper case letters.
	Lu     = _Lu // Lu is the set of Unicode characters in category Lu (Letter, uppercase).
	Zl     = _Zl // Zl is the set of Unicode characters in category Zl (Separator, line).
	Zp     = _Zp // Zp is the set of Unicode characters in category Zp (Separator, paragraph).
	Zs     = _Zs // Zs is the set of Unicode characters in category Zs (Separator, space).
)`,
			variables: []testVariable{
				{
					name: "Cc",
				},
				{
					name: "Cf",
				},
				{
					name: "Co",
				},
				{
					name: "Cs",
				},
				{
					name: "Digit",
				},
				{
					name: "Nd",
				},
				{
					name: "Letter",
				},
				{
					name: "L",
				},
				{
					name: "Lm",
				},
				{
					name: "Lo",
				},
				{
					name: "Lower",
				},
				{
					name: "Ll",
				},
				{
					name: "Mark",
				},
				{
					name: "M",
				},
				{
					name: "Mc",
				},
				{
					name: "Me",
				},
				{
					name: "Mn",
				},
				{
					name: "Nl",
				},
				{
					name: "No",
				},
				{
					name: "Number",
				},
				{
					name: "N",
				},
				{
					name: "Other",
				},
				{
					name: "C",
				},
				{
					name: "Pc",
				},
				{
					name: "Pd",
				},
				{
					name: "Pe",
				},
				{
					name: "Pf",
				},
				{
					name: "Pi",
				},
				{
					name: "Po",
				},
				{
					name: "Ps",
				},
				{
					name: "Punct",
				},
				{
					name: "P",
				},
				{
					name: "Sc",
				},
				{
					name: "Sk",
				},
				{
					name: "Sm",
				},
				{
					name: "So",
				},
				{
					name: "Space",
				},
				{
					name: "Z",
				},
				{
					name: "Symbol",
				},
				{
					name: "S",
				},
				{
					name: "Title",
				},
				{
					name: "Lt",
				},
				{
					name: "Upper",
				},
				{
					name: "Lu",
				},
				{
					name: "Zl",
				},
				{
					name: "Zp",
				},
				{
					name: "Zs",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `These variables have type *RangeTable.
`,
			source: `var (
	Adlam                  = _Adlam                  // Adlam is the set of Unicode characters in script Adlam.
	Ahom                   = _Ahom                   // Ahom is the set of Unicode characters in script Ahom.
	Anatolian_Hieroglyphs  = _Anatolian_Hieroglyphs  // Anatolian_Hieroglyphs is the set of Unicode characters in script Anatolian_Hieroglyphs.
	Arabic                 = _Arabic                 // Arabic is the set of Unicode characters in script Arabic.
	Armenian               = _Armenian               // Armenian is the set of Unicode characters in script Armenian.
	Avestan                = _Avestan                // Avestan is the set of Unicode characters in script Avestan.
	Balinese               = _Balinese               // Balinese is the set of Unicode characters in script Balinese.
	Bamum                  = _Bamum                  // Bamum is the set of Unicode characters in script Bamum.
	Bassa_Vah              = _Bassa_Vah              // Bassa_Vah is the set of Unicode characters in script Bassa_Vah.
	Batak                  = _Batak                  // Batak is the set of Unicode characters in script Batak.
	Bengali                = _Bengali                // Bengali is the set of Unicode characters in script Bengali.
	Bhaiksuki              = _Bhaiksuki              // Bhaiksuki is the set of Unicode characters in script Bhaiksuki.
	Bopomofo               = _Bopomofo               // Bopomofo is the set of Unicode characters in script Bopomofo.
	Brahmi                 = _Brahmi                 // Brahmi is the set of Unicode characters in script Brahmi.
	Braille                = _Braille                // Braille is the set of Unicode characters in script Braille.
	Buginese               = _Buginese               // Buginese is the set of Unicode characters in script Buginese.
	Buhid                  = _Buhid                  // Buhid is the set of Unicode characters in script Buhid.
	Canadian_Aboriginal    = _Canadian_Aboriginal    // Canadian_Aboriginal is the set of Unicode characters in script Canadian_Aboriginal.
	Carian                 = _Carian                 // Carian is the set of Unicode characters in script Carian.
	Caucasian_Albanian     = _Caucasian_Albanian     // Caucasian_Albanian is the set of Unicode characters in script Caucasian_Albanian.
	Chakma                 = _Chakma                 // Chakma is the set of Unicode characters in script Chakma.
	Cham                   = _Cham                   // Cham is the set of Unicode characters in script Cham.
	Cherokee               = _Cherokee               // Cherokee is the set of Unicode characters in script Cherokee.
	Chorasmian             = _Chorasmian             // Chorasmian is the set of Unicode characters in script Chorasmian.
	Common                 = _Common                 // Common is the set of Unicode characters in script Common.
	Coptic                 = _Coptic                 // Coptic is the set of Unicode characters in script Coptic.
	Cuneiform              = _Cuneiform              // Cuneiform is the set of Unicode characters in script Cuneiform.
	Cypriot                = _Cypriot                // Cypriot is the set of Unicode characters in script Cypriot.
	Cyrillic               = _Cyrillic               // Cyrillic is the set of Unicode characters in script Cyrillic.
	Deseret                = _Deseret                // Deseret is the set of Unicode characters in script Deseret.
	Devanagari             = _Devanagari             // Devanagari is the set of Unicode characters in script Devanagari.
	Dives_Akuru            = _Dives_Akuru            // Dives_Akuru is the set of Unicode characters in script Dives_Akuru.
	Dogra                  = _Dogra                  // Dogra is the set of Unicode characters in script Dogra.
	Duployan               = _Duployan               // Duployan is the set of Unicode characters in script Duployan.
	Egyptian_Hieroglyphs   = _Egyptian_Hieroglyphs   // Egyptian_Hieroglyphs is the set of Unicode characters in script Egyptian_Hieroglyphs.
	Elbasan                = _Elbasan                // Elbasan is the set of Unicode characters in script Elbasan.
	Elymaic                = _Elymaic                // Elymaic is the set of Unicode characters in script Elymaic.
	Ethiopic               = _Ethiopic               // Ethiopic is the set of Unicode characters in script Ethiopic.
	Georgian               = _Georgian               // Georgian is the set of Unicode characters in script Georgian.
	Glagolitic             = _Glagolitic             // Glagolitic is the set of Unicode characters in script Glagolitic.
	Gothic                 = _Gothic                 // Gothic is the set of Unicode characters in script Gothic.
	Grantha                = _Grantha                // Grantha is the set of Unicode characters in script Grantha.
	Greek                  = _Greek                  // Greek is the set of Unicode characters in script Greek.
	Gujarati               = _Gujarati               // Gujarati is the set of Unicode characters in script Gujarati.
	Gunjala_Gondi          = _Gunjala_Gondi          // Gunjala_Gondi is the set of Unicode characters in script Gunjala_Gondi.
	Gurmukhi               = _Gurmukhi               // Gurmukhi is the set of Unicode characters in script Gurmukhi.
	Han                    = _Han                    // Han is the set of Unicode characters in script Han.
	Hangul                 = _Hangul                 // Hangul is the set of Unicode characters in script Hangul.
	Hanifi_Rohingya        = _Hanifi_Rohingya        // Hanifi_Rohingya is the set of Unicode characters in script Hanifi_Rohingya.
	Hanunoo                = _Hanunoo                // Hanunoo is the set of Unicode characters in script Hanunoo.
	Hatran                 = _Hatran                 // Hatran is the set of Unicode characters in script Hatran.
	Hebrew                 = _Hebrew                 // Hebrew is the set of Unicode characters in script Hebrew.
	Hiragana               = _Hiragana               // Hiragana is the set of Unicode characters in script Hiragana.
	Imperial_Aramaic       = _Imperial_Aramaic       // Imperial_Aramaic is the set of Unicode characters in script Imperial_Aramaic.
	Inherited              = _Inherited              // Inherited is the set of Unicode characters in script Inherited.
	Inscriptional_Pahlavi  = _Inscriptional_Pahlavi  // Inscriptional_Pahlavi is the set of Unicode characters in script Inscriptional_Pahlavi.
	Inscriptional_Parthian = _Inscriptional_Parthian // Inscriptional_Parthian is the set of Unicode characters in script Inscriptional_Parthian.
	Javanese               = _Javanese               // Javanese is the set of Unicode characters in script Javanese.
	Kaithi                 = _Kaithi                 // Kaithi is the set of Unicode characters in script Kaithi.
	Kannada                = _Kannada                // Kannada is the set of Unicode characters in script Kannada.
	Katakana               = _Katakana               // Katakana is the set of Unicode characters in script Katakana.
	Kayah_Li               = _Kayah_Li               // Kayah_Li is the set of Unicode characters in script Kayah_Li.
	Kharoshthi             = _Kharoshthi             // Kharoshthi is the set of Unicode characters in script Kharoshthi.
	Khitan_Small_Script    = _Khitan_Small_Script    // Khitan_Small_Script is the set of Unicode characters in script Khitan_Small_Script.
	Khmer                  = _Khmer                  // Khmer is the set of Unicode characters in script Khmer.
	Khojki                 = _Khojki                 // Khojki is the set of Unicode characters in script Khojki.
	Khudawadi              = _Khudawadi              // Khudawadi is the set of Unicode characters in script Khudawadi.
	Lao                    = _Lao                    // Lao is the set of Unicode characters in script Lao.
	Latin                  = _Latin                  // Latin is the set of Unicode characters in script Latin.
	Lepcha                 = _Lepcha                 // Lepcha is the set of Unicode characters in script Lepcha.
	Limbu                  = _Limbu                  // Limbu is the set of Unicode characters in script Limbu.
	Linear_A               = _Linear_A               // Linear_A is the set of Unicode characters in script Linear_A.
	Linear_B               = _Linear_B               // Linear_B is the set of Unicode characters in script Linear_B.
	Lisu                   = _Lisu                   // Lisu is the set of Unicode characters in script Lisu.
	Lycian                 = _Lycian                 // Lycian is the set of Unicode characters in script Lycian.
	Lydian                 = _Lydian                 // Lydian is the set of Unicode characters in script Lydian.
	Mahajani               = _Mahajani               // Mahajani is the set of Unicode characters in script Mahajani.
	Makasar                = _Makasar                // Makasar is the set of Unicode characters in script Makasar.
	Malayalam              = _Malayalam              // Malayalam is the set of Unicode characters in script Malayalam.
	Mandaic                = _Mandaic                // Mandaic is the set of Unicode characters in script Mandaic.
	Manichaean             = _Manichaean             // Manichaean is the set of Unicode characters in script Manichaean.
	Marchen                = _Marchen                // Marchen is the set of Unicode characters in script Marchen.
	Masaram_Gondi          = _Masaram_Gondi          // Masaram_Gondi is the set of Unicode characters in script Masaram_Gondi.
	Medefaidrin            = _Medefaidrin            // Medefaidrin is the set of Unicode characters in script Medefaidrin.
	Meetei_Mayek           = _Meetei_Mayek           // Meetei_Mayek is the set of Unicode characters in script Meetei_Mayek.
	Mende_Kikakui          = _Mende_Kikakui          // Mende_Kikakui is the set of Unicode characters in script Mende_Kikakui.
	Meroitic_Cursive       = _Meroitic_Cursive       // Meroitic_Cursive is the set of Unicode characters in script Meroitic_Cursive.
	Meroitic_Hieroglyphs   = _Meroitic_Hieroglyphs   // Meroitic_Hieroglyphs is the set of Unicode characters in script Meroitic_Hieroglyphs.
	Miao                   = _Miao                   // Miao is the set of Unicode characters in script Miao.
	Modi                   = _Modi                   // Modi is the set of Unicode characters in script Modi.
	Mongolian              = _Mongolian              // Mongolian is the set of Unicode characters in script Mongolian.
	Mro                    = _Mro                    // Mro is the set of Unicode characters in script Mro.
	Multani                = _Multani                // Multani is the set of Unicode characters in script Multani.
	Myanmar                = _Myanmar                // Myanmar is the set of Unicode characters in script Myanmar.
	Nabataean              = _Nabataean              // Nabataean is the set of Unicode characters in script Nabataean.
	Nandinagari            = _Nandinagari            // Nandinagari is the set of Unicode characters in script Nandinagari.
	New_Tai_Lue            = _New_Tai_Lue            // New_Tai_Lue is the set of Unicode characters in script New_Tai_Lue.
	Newa                   = _Newa                   // Newa is the set of Unicode characters in script Newa.
	Nko                    = _Nko                    // Nko is the set of Unicode characters in script Nko.
	Nushu                  = _Nushu                  // Nushu is the set of Unicode characters in script Nushu.
	Nyiakeng_Puachue_Hmong = _Nyiakeng_Puachue_Hmong // Nyiakeng_Puachue_Hmong is the set of Unicode characters in script Nyiakeng_Puachue_Hmong.
	Ogham                  = _Ogham                  // Ogham is the set of Unicode characters in script Ogham.
	Ol_Chiki               = _Ol_Chiki               // Ol_Chiki is the set of Unicode characters in script Ol_Chiki.
	Old_Hungarian          = _Old_Hungarian          // Old_Hungarian is the set of Unicode characters in script Old_Hungarian.
	Old_Italic             = _Old_Italic             // Old_Italic is the set of Unicode characters in script Old_Italic.
	Old_North_Arabian      = _Old_North_Arabian      // Old_North_Arabian is the set of Unicode characters in script Old_North_Arabian.
	Old_Permic             = _Old_Permic             // Old_Permic is the set of Unicode characters in script Old_Permic.
	Old_Persian            = _Old_Persian            // Old_Persian is the set of Unicode characters in script Old_Persian.
	Old_Sogdian            = _Old_Sogdian            // Old_Sogdian is the set of Unicode characters in script Old_Sogdian.
	Old_South_Arabian      = _Old_South_Arabian      // Old_South_Arabian is the set of Unicode characters in script Old_South_Arabian.
	Old_Turkic             = _Old_Turkic             // Old_Turkic is the set of Unicode characters in script Old_Turkic.
	Oriya                  = _Oriya                  // Oriya is the set of Unicode characters in script Oriya.
	Osage                  = _Osage                  // Osage is the set of Unicode characters in script Osage.
	Osmanya                = _Osmanya                // Osmanya is the set of Unicode characters in script Osmanya.
	Pahawh_Hmong           = _Pahawh_Hmong           // Pahawh_Hmong is the set of Unicode characters in script Pahawh_Hmong.
	Palmyrene              = _Palmyrene              // Palmyrene is the set of Unicode characters in script Palmyrene.
	Pau_Cin_Hau            = _Pau_Cin_Hau            // Pau_Cin_Hau is the set of Unicode characters in script Pau_Cin_Hau.
	Phags_Pa               = _Phags_Pa               // Phags_Pa is the set of Unicode characters in script Phags_Pa.
	Phoenician             = _Phoenician             // Phoenician is the set of Unicode characters in script Phoenician.
	Psalter_Pahlavi        = _Psalter_Pahlavi        // Psalter_Pahlavi is the set of Unicode characters in script Psalter_Pahlavi.
	Rejang                 = _Rejang                 // Rejang is the set of Unicode characters in script Rejang.
	Runic                  = _Runic                  // Runic is the set of Unicode characters in script Runic.
	Samaritan              = _Samaritan              // Samaritan is the set of Unicode characters in script Samaritan.
	Saurashtra             = _Saurashtra             // Saurashtra is the set of Unicode characters in script Saurashtra.
	Sharada                = _Sharada                // Sharada is the set of Unicode characters in script Sharada.
	Shavian                = _Shavian                // Shavian is the set of Unicode characters in script Shavian.
	Siddham                = _Siddham                // Siddham is the set of Unicode characters in script Siddham.
	SignWriting            = _SignWriting            // SignWriting is the set of Unicode characters in script SignWriting.
	Sinhala                = _Sinhala                // Sinhala is the set of Unicode characters in script Sinhala.
	Sogdian                = _Sogdian                // Sogdian is the set of Unicode characters in script Sogdian.
	Sora_Sompeng           = _Sora_Sompeng           // Sora_Sompeng is the set of Unicode characters in script Sora_Sompeng.
	Soyombo                = _Soyombo                // Soyombo is the set of Unicode characters in script Soyombo.
	Sundanese              = _Sundanese              // Sundanese is the set of Unicode characters in script Sundanese.
	Syloti_Nagri           = _Syloti_Nagri           // Syloti_Nagri is the set of Unicode characters in script Syloti_Nagri.
	Syriac                 = _Syriac                 // Syriac is the set of Unicode characters in script Syriac.
	Tagalog                = _Tagalog                // Tagalog is the set of Unicode characters in script Tagalog.
	Tagbanwa               = _Tagbanwa               // Tagbanwa is the set of Unicode characters in script Tagbanwa.
	Tai_Le                 = _Tai_Le                 // Tai_Le is the set of Unicode characters in script Tai_Le.
	Tai_Tham               = _Tai_Tham               // Tai_Tham is the set of Unicode characters in script Tai_Tham.
	Tai_Viet               = _Tai_Viet               // Tai_Viet is the set of Unicode characters in script Tai_Viet.
	Takri                  = _Takri                  // Takri is the set of Unicode characters in script Takri.
	Tamil                  = _Tamil                  // Tamil is the set of Unicode characters in script Tamil.
	Tangut                 = _Tangut                 // Tangut is the set of Unicode characters in script Tangut.
	Telugu                 = _Telugu                 // Telugu is the set of Unicode characters in script Telugu.
	Thaana                 = _Thaana                 // Thaana is the set of Unicode characters in script Thaana.
	Thai                   = _Thai                   // Thai is the set of Unicode characters in script Thai.
	Tibetan                = _Tibetan                // Tibetan is the set of Unicode characters in script Tibetan.
	Tifinagh               = _Tifinagh               // Tifinagh is the set of Unicode characters in script Tifinagh.
	Tirhuta                = _Tirhuta                // Tirhuta is the set of Unicode characters in script Tirhuta.
	Ugaritic               = _Ugaritic               // Ugaritic is the set of Unicode characters in script Ugaritic.
	Vai                    = _Vai                    // Vai is the set of Unicode characters in script Vai.
	Wancho                 = _Wancho                 // Wancho is the set of Unicode characters in script Wancho.
	Warang_Citi            = _Warang_Citi            // Warang_Citi is the set of Unicode characters in script Warang_Citi.
	Yezidi                 = _Yezidi                 // Yezidi is the set of Unicode characters in script Yezidi.
	Yi                     = _Yi                     // Yi is the set of Unicode characters in script Yi.
	Zanabazar_Square       = _Zanabazar_Square       // Zanabazar_Square is the set of Unicode characters in script Zanabazar_Square.
)`,
			variables: []testVariable{
				{
					name: "Adlam",
				},
				{
					name: "Ahom",
				},
				{
					name: "Anatolian_Hieroglyphs",
				},
				{
					name: "Arabic",
				},
				{
					name: "Armenian",
				},
				{
					name: "Avestan",
				},
				{
					name: "Balinese",
				},
				{
					name: "Bamum",
				},
				{
					name: "Bassa_Vah",
				},
				{
					name: "Batak",
				},
				{
					name: "Bengali",
				},
				{
					name: "Bhaiksuki",
				},
				{
					name: "Bopomofo",
				},
				{
					name: "Brahmi",
				},
				{
					name: "Braille",
				},
				{
					name: "Buginese",
				},
				{
					name: "Buhid",
				},
				{
					name: "Canadian_Aboriginal",
				},
				{
					name: "Carian",
				},
				{
					name: "Caucasian_Albanian",
				},
				{
					name: "Chakma",
				},
				{
					name: "Cham",
				},
				{
					name: "Cherokee",
				},
				{
					name: "Chorasmian",
				},
				{
					name: "Common",
				},
				{
					name: "Coptic",
				},
				{
					name: "Cuneiform",
				},
				{
					name: "Cypriot",
				},
				{
					name: "Cyrillic",
				},
				{
					name: "Deseret",
				},
				{
					name: "Devanagari",
				},
				{
					name: "Dives_Akuru",
				},
				{
					name: "Dogra",
				},
				{
					name: "Duployan",
				},
				{
					name: "Egyptian_Hieroglyphs",
				},
				{
					name: "Elbasan",
				},
				{
					name: "Elymaic",
				},
				{
					name: "Ethiopic",
				},
				{
					name: "Georgian",
				},
				{
					name: "Glagolitic",
				},
				{
					name: "Gothic",
				},
				{
					name: "Grantha",
				},
				{
					name: "Greek",
				},
				{
					name: "Gujarati",
				},
				{
					name: "Gunjala_Gondi",
				},
				{
					name: "Gurmukhi",
				},
				{
					name: "Han",
				},
				{
					name: "Hangul",
				},
				{
					name: "Hanifi_Rohingya",
				},
				{
					name: "Hanunoo",
				},
				{
					name: "Hatran",
				},
				{
					name: "Hebrew",
				},
				{
					name: "Hiragana",
				},
				{
					name: "Imperial_Aramaic",
				},
				{
					name: "Inherited",
				},
				{
					name: "Inscriptional_Pahlavi",
				},
				{
					name: "Inscriptional_Parthian",
				},
				{
					name: "Javanese",
				},
				{
					name: "Kaithi",
				},
				{
					name: "Kannada",
				},
				{
					name: "Katakana",
				},
				{
					name: "Kayah_Li",
				},
				{
					name: "Kharoshthi",
				},
				{
					name: "Khitan_Small_Script",
				},
				{
					name: "Khmer",
				},
				{
					name: "Khojki",
				},
				{
					name: "Khudawadi",
				},
				{
					name: "Lao",
				},
				{
					name: "Latin",
				},
				{
					name: "Lepcha",
				},
				{
					name: "Limbu",
				},
				{
					name: "Linear_A",
				},
				{
					name: "Linear_B",
				},
				{
					name: "Lisu",
				},
				{
					name: "Lycian",
				},
				{
					name: "Lydian",
				},
				{
					name: "Mahajani",
				},
				{
					name: "Makasar",
				},
				{
					name: "Malayalam",
				},
				{
					name: "Mandaic",
				},
				{
					name: "Manichaean",
				},
				{
					name: "Marchen",
				},
				{
					name: "Masaram_Gondi",
				},
				{
					name: "Medefaidrin",
				},
				{
					name: "Meetei_Mayek",
				},
				{
					name: "Mende_Kikakui",
				},
				{
					name: "Meroitic_Cursive",
				},
				{
					name: "Meroitic_Hieroglyphs",
				},
				{
					name: "Miao",
				},
				{
					name: "Modi",
				},
				{
					name: "Mongolian",
				},
				{
					name: "Mro",
				},
				{
					name: "Multani",
				},
				{
					name: "Myanmar",
				},
				{
					name: "Nabataean",
				},
				{
					name: "Nandinagari",
				},
				{
					name: "New_Tai_Lue",
				},
				{
					name: "Newa",
				},
				{
					name: "Nko",
				},
				{
					name: "Nushu",
				},
				{
					name: "Nyiakeng_Puachue_Hmong",
				},
				{
					name: "Ogham",
				},
				{
					name: "Ol_Chiki",
				},
				{
					name: "Old_Hungarian",
				},
				{
					name: "Old_Italic",
				},
				{
					name: "Old_North_Arabian",
				},
				{
					name: "Old_Permic",
				},
				{
					name: "Old_Persian",
				},
				{
					name: "Old_Sogdian",
				},
				{
					name: "Old_South_Arabian",
				},
				{
					name: "Old_Turkic",
				},
				{
					name: "Oriya",
				},
				{
					name: "Osage",
				},
				{
					name: "Osmanya",
				},
				{
					name: "Pahawh_Hmong",
				},
				{
					name: "Palmyrene",
				},
				{
					name: "Pau_Cin_Hau",
				},
				{
					name: "Phags_Pa",
				},
				{
					name: "Phoenician",
				},
				{
					name: "Psalter_Pahlavi",
				},
				{
					name: "Rejang",
				},
				{
					name: "Runic",
				},
				{
					name: "Samaritan",
				},
				{
					name: "Saurashtra",
				},
				{
					name: "Sharada",
				},
				{
					name: "Shavian",
				},
				{
					name: "Siddham",
				},
				{
					name: "SignWriting",
				},
				{
					name: "Sinhala",
				},
				{
					name: "Sogdian",
				},
				{
					name: "Sora_Sompeng",
				},
				{
					name: "Soyombo",
				},
				{
					name: "Sundanese",
				},
				{
					name: "Syloti_Nagri",
				},
				{
					name: "Syriac",
				},
				{
					name: "Tagalog",
				},
				{
					name: "Tagbanwa",
				},
				{
					name: "Tai_Le",
				},
				{
					name: "Tai_Tham",
				},
				{
					name: "Tai_Viet",
				},
				{
					name: "Takri",
				},
				{
					name: "Tamil",
				},
				{
					name: "Tangut",
				},
				{
					name: "Telugu",
				},
				{
					name: "Thaana",
				},
				{
					name: "Thai",
				},
				{
					name: "Tibetan",
				},
				{
					name: "Tifinagh",
				},
				{
					name: "Tirhuta",
				},
				{
					name: "Ugaritic",
				},
				{
					name: "Vai",
				},
				{
					name: "Wancho",
				},
				{
					name: "Warang_Citi",
				},
				{
					name: "Yezidi",
				},
				{
					name: "Yi",
				},
				{
					name: "Zanabazar_Square",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `These variables have type *RangeTable.
`,
			source: `var (
	ASCII_Hex_Digit                    = _ASCII_Hex_Digit                    // ASCII_Hex_Digit is the set of Unicode characters with property ASCII_Hex_Digit.
	Bidi_Control                       = _Bidi_Control                       // Bidi_Control is the set of Unicode characters with property Bidi_Control.
	Dash                               = _Dash                               // Dash is the set of Unicode characters with property Dash.
	Deprecated                         = _Deprecated                         // Deprecated is the set of Unicode characters with property Deprecated.
	Diacritic                          = _Diacritic                          // Diacritic is the set of Unicode characters with property Diacritic.
	Extender                           = _Extender                           // Extender is the set of Unicode characters with property Extender.
	Hex_Digit                          = _Hex_Digit                          // Hex_Digit is the set of Unicode characters with property Hex_Digit.
	Hyphen                             = _Hyphen                             // Hyphen is the set of Unicode characters with property Hyphen.
	IDS_Binary_Operator                = _IDS_Binary_Operator                // IDS_Binary_Operator is the set of Unicode characters with property IDS_Binary_Operator.
	IDS_Trinary_Operator               = _IDS_Trinary_Operator               // IDS_Trinary_Operator is the set of Unicode characters with property IDS_Trinary_Operator.
	Ideographic                        = _Ideographic                        // Ideographic is the set of Unicode characters with property Ideographic.
	Join_Control                       = _Join_Control                       // Join_Control is the set of Unicode characters with property Join_Control.
	Logical_Order_Exception            = _Logical_Order_Exception            // Logical_Order_Exception is the set of Unicode characters with property Logical_Order_Exception.
	Noncharacter_Code_Point            = _Noncharacter_Code_Point            // Noncharacter_Code_Point is the set of Unicode characters with property Noncharacter_Code_Point.
	Other_Alphabetic                   = _Other_Alphabetic                   // Other_Alphabetic is the set of Unicode characters with property Other_Alphabetic.
	Other_Default_Ignorable_Code_Point = _Other_Default_Ignorable_Code_Point // Other_Default_Ignorable_Code_Point is the set of Unicode characters with property Other_Default_Ignorable_Code_Point.
	Other_Grapheme_Extend              = _Other_Grapheme_Extend              // Other_Grapheme_Extend is the set of Unicode characters with property Other_Grapheme_Extend.
	Other_ID_Continue                  = _Other_ID_Continue                  // Other_ID_Continue is the set of Unicode characters with property Other_ID_Continue.
	Other_ID_Start                     = _Other_ID_Start                     // Other_ID_Start is the set of Unicode characters with property Other_ID_Start.
	Other_Lowercase                    = _Other_Lowercase                    // Other_Lowercase is the set of Unicode characters with property Other_Lowercase.
	Other_Math                         = _Other_Math                         // Other_Math is the set of Unicode characters with property Other_Math.
	Other_Uppercase                    = _Other_Uppercase                    // Other_Uppercase is the set of Unicode characters with property Other_Uppercase.
	Pattern_Syntax                     = _Pattern_Syntax                     // Pattern_Syntax is the set of Unicode characters with property Pattern_Syntax.
	Pattern_White_Space                = _Pattern_White_Space                // Pattern_White_Space is the set of Unicode characters with property Pattern_White_Space.
	Prepended_Concatenation_Mark       = _Prepended_Concatenation_Mark       // Prepended_Concatenation_Mark is the set of Unicode characters with property Prepended_Concatenation_Mark.
	Quotation_Mark                     = _Quotation_Mark                     // Quotation_Mark is the set of Unicode characters with property Quotation_Mark.
	Radical                            = _Radical                            // Radical is the set of Unicode characters with property Radical.
	Regional_Indicator                 = _Regional_Indicator                 // Regional_Indicator is the set of Unicode characters with property Regional_Indicator.
	STerm                              = _Sentence_Terminal                  // STerm is an alias for Sentence_Terminal.
	Sentence_Terminal                  = _Sentence_Terminal                  // Sentence_Terminal is the set of Unicode characters with property Sentence_Terminal.
	Soft_Dotted                        = _Soft_Dotted                        // Soft_Dotted is the set of Unicode characters with property Soft_Dotted.
	Terminal_Punctuation               = _Terminal_Punctuation               // Terminal_Punctuation is the set of Unicode characters with property Terminal_Punctuation.
	Unified_Ideograph                  = _Unified_Ideograph                  // Unified_Ideograph is the set of Unicode characters with property Unified_Ideograph.
	Variation_Selector                 = _Variation_Selector                 // Variation_Selector is the set of Unicode characters with property Variation_Selector.
	White_Space                        = _White_Space                        // White_Space is the set of Unicode characters with property White_Space.
)`,
			variables: []testVariable{
				{
					name: "ASCII_Hex_Digit",
				},
				{
					name: "Bidi_Control",
				},
				{
					name: "Dash",
				},
				{
					name: "Deprecated",
				},
				{
					name: "Diacritic",
				},
				{
					name: "Extender",
				},
				{
					name: "Hex_Digit",
				},
				{
					name: "Hyphen",
				},
				{
					name: "IDS_Binary_Operator",
				},
				{
					name: "IDS_Trinary_Operator",
				},
				{
					name: "Ideographic",
				},
				{
					name: "Join_Control",
				},
				{
					name: "Logical_Order_Exception",
				},
				{
					name: "Noncharacter_Code_Point",
				},
				{
					name: "Other_Alphabetic",
				},
				{
					name: "Other_Default_Ignorable_Code_Point",
				},
				{
					name: "Other_Grapheme_Extend",
				},
				{
					name: "Other_ID_Continue",
				},
				{
					name: "Other_ID_Start",
				},
				{
					name: "Other_Lowercase",
				},
				{
					name: "Other_Math",
				},
				{
					name: "Other_Uppercase",
				},
				{
					name: "Pattern_Syntax",
				},
				{
					name: "Pattern_White_Space",
				},
				{
					name: "Prepended_Concatenation_Mark",
				},
				{
					name: "Quotation_Mark",
				},
				{
					name: "Radical",
				},
				{
					name: "Regional_Indicator",
				},
				{
					name: "STerm",
				},
				{
					name: "Sentence_Terminal",
				},
				{
					name: "Soft_Dotted",
				},
				{
					name: "Terminal_Punctuation",
				},
				{
					name: "Unified_Ideograph",
				},
				{
					name: "Variation_Selector",
				},
				{
					name: "White_Space",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `CaseRanges is the table describing case mappings for all letters with non-self mappings.
`,
			source: `var CaseRanges = _CaseRanges`,
			variables: []testVariable{
				{
					name: "CaseRanges",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `Categories is the set of Unicode category tables.
`,
			source: `var Categories = map[string]*RangeTable{
	"C":  C,
	"Cc": Cc,
	"Cf": Cf,
	"Co": Co,
	"Cs": Cs,
	"L":  L,
	"Ll": Ll,
	"Lm": Lm,
	"Lo": Lo,
	"Lt": Lt,
	"Lu": Lu,
	"M":  M,
	"Mc": Mc,
	"Me": Me,
	"Mn": Mn,
	"N":  N,
	"Nd": Nd,
	"Nl": Nl,
	"No": No,
	"P":  P,
	"Pc": Pc,
	"Pd": Pd,
	"Pe": Pe,
	"Pf": Pf,
	"Pi": Pi,
	"Po": Po,
	"Ps": Ps,
	"S":  S,
	"Sc": Sc,
	"Sk": Sk,
	"Sm": Sm,
	"So": So,
	"Z":  Z,
	"Zl": Zl,
	"Zp": Zp,
	"Zs": Zs,
}`,
			variables: []testVariable{
				{
					name: "Categories",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `FoldCategory maps a category name to a table of code points outside the category that are equivalent under simple case folding to code points inside the category. If there is no entry for a category name, there are no such points.
`,
			source: `var FoldCategory = map[string]*RangeTable{
	"L":  foldL,
	"Ll": foldLl,
	"Lt": foldLt,
	"Lu": foldLu,
	"M":  foldM,
	"Mn": foldMn,
}`,
			variables: []testVariable{
				{
					name: "FoldCategory",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `FoldScript maps a script name to a table of code points outside the script that are equivalent under simple case folding to code points inside the script. If there is no entry for a script name, there are no such points.
`,
			source: `var FoldScript = map[string]*RangeTable{
	"Common":    foldCommon,
	"Greek":     foldGreek,
	"Inherited": foldInherited,
}`,
			variables: []testVariable{
				{
					name: "FoldScript",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `GraphicRanges defines the set of graphic characters according to Unicode.
`,
			source: `var GraphicRanges = []*RangeTable{
	L, M, N, P, S, Zs,
}`,
			variables: []testVariable{
				{
					name: "GraphicRanges",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `PrintRanges defines the set of printable characters according to Go. ASCII space, U+0020, is handled separately.
`,
			source: `var PrintRanges = []*RangeTable{
	L, M, N, P, S,
}`,
			variables: []testVariable{
				{
					name: "PrintRanges",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `Properties is the set of Unicode property tables.
`,
			source: `var Properties = map[string]*RangeTable{
	"ASCII_Hex_Digit":                    ASCII_Hex_Digit,
	"Bidi_Control":                       Bidi_Control,
	"Dash":                               Dash,
	"Deprecated":                         Deprecated,
	"Diacritic":                          Diacritic,
	"Extender":                           Extender,
	"Hex_Digit":                          Hex_Digit,
	"Hyphen":                             Hyphen,
	"IDS_Binary_Operator":                IDS_Binary_Operator,
	"IDS_Trinary_Operator":               IDS_Trinary_Operator,
	"Ideographic":                        Ideographic,
	"Join_Control":                       Join_Control,
	"Logical_Order_Exception":            Logical_Order_Exception,
	"Noncharacter_Code_Point":            Noncharacter_Code_Point,
	"Other_Alphabetic":                   Other_Alphabetic,
	"Other_Default_Ignorable_Code_Point": Other_Default_Ignorable_Code_Point,
	"Other_Grapheme_Extend":              Other_Grapheme_Extend,
	"Other_ID_Continue":                  Other_ID_Continue,
	"Other_ID_Start":                     Other_ID_Start,
	"Other_Lowercase":                    Other_Lowercase,
	"Other_Math":                         Other_Math,
	"Other_Uppercase":                    Other_Uppercase,
	"Pattern_Syntax":                     Pattern_Syntax,
	"Pattern_White_Space":                Pattern_White_Space,
	"Prepended_Concatenation_Mark":       Prepended_Concatenation_Mark,
	"Quotation_Mark":                     Quotation_Mark,
	"Radical":                            Radical,
	"Regional_Indicator":                 Regional_Indicator,
	"Sentence_Terminal":                  Sentence_Terminal,
	"STerm":                              Sentence_Terminal,
	"Soft_Dotted":                        Soft_Dotted,
	"Terminal_Punctuation":               Terminal_Punctuation,
	"Unified_Ideograph":                  Unified_Ideograph,
	"Variation_Selector":                 Variation_Selector,
	"White_Space":                        White_Space,
}`,
			variables: []testVariable{
				{
					name: "Properties",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: `Scripts is the set of Unicode script tables.
`,
			source: `var Scripts = map[string]*RangeTable{
	"Adlam":                  Adlam,
	"Ahom":                   Ahom,
	"Anatolian_Hieroglyphs":  Anatolian_Hieroglyphs,
	"Arabic":                 Arabic,
	"Armenian":               Armenian,
	"Avestan":                Avestan,
	"Balinese":               Balinese,
	"Bamum":                  Bamum,
	"Bassa_Vah":              Bassa_Vah,
	"Batak":                  Batak,
	"Bengali":                Bengali,
	"Bhaiksuki":              Bhaiksuki,
	"Bopomofo":               Bopomofo,
	"Brahmi":                 Brahmi,
	"Braille":                Braille,
	"Buginese":               Buginese,
	"Buhid":                  Buhid,
	"Canadian_Aboriginal":    Canadian_Aboriginal,
	"Carian":                 Carian,
	"Caucasian_Albanian":     Caucasian_Albanian,
	"Chakma":                 Chakma,
	"Cham":                   Cham,
	"Cherokee":               Cherokee,
	"Chorasmian":             Chorasmian,
	"Common":                 Common,
	"Coptic":                 Coptic,
	"Cuneiform":              Cuneiform,
	"Cypriot":                Cypriot,
	"Cyrillic":               Cyrillic,
	"Deseret":                Deseret,
	"Devanagari":             Devanagari,
	"Dives_Akuru":            Dives_Akuru,
	"Dogra":                  Dogra,
	"Duployan":               Duployan,
	"Egyptian_Hieroglyphs":   Egyptian_Hieroglyphs,
	"Elbasan":                Elbasan,
	"Elymaic":                Elymaic,
	"Ethiopic":               Ethiopic,
	"Georgian":               Georgian,
	"Glagolitic":             Glagolitic,
	"Gothic":                 Gothic,
	"Grantha":                Grantha,
	"Greek":                  Greek,
	"Gujarati":               Gujarati,
	"Gunjala_Gondi":          Gunjala_Gondi,
	"Gurmukhi":               Gurmukhi,
	"Han":                    Han,
	"Hangul":                 Hangul,
	"Hanifi_Rohingya":        Hanifi_Rohingya,
	"Hanunoo":                Hanunoo,
	"Hatran":                 Hatran,
	"Hebrew":                 Hebrew,
	"Hiragana":               Hiragana,
	"Imperial_Aramaic":       Imperial_Aramaic,
	"Inherited":              Inherited,
	"Inscriptional_Pahlavi":  Inscriptional_Pahlavi,
	"Inscriptional_Parthian": Inscriptional_Parthian,
	"Javanese":               Javanese,
	"Kaithi":                 Kaithi,
	"Kannada":                Kannada,
	"Katakana":               Katakana,
	"Kayah_Li":               Kayah_Li,
	"Kharoshthi":             Kharoshthi,
	"Khitan_Small_Script":    Khitan_Small_Script,
	"Khmer":                  Khmer,
	"Khojki":                 Khojki,
	"Khudawadi":              Khudawadi,
	"Lao":                    Lao,
	"Latin":                  Latin,
	"Lepcha":                 Lepcha,
	"Limbu":                  Limbu,
	"Linear_A":               Linear_A,
	"Linear_B":               Linear_B,
	"Lisu":                   Lisu,
	"Lycian":                 Lycian,
	"Lydian":                 Lydian,
	"Mahajani":               Mahajani,
	"Makasar":                Makasar,
	"Malayalam":              Malayalam,
	"Mandaic":                Mandaic,
	"Manichaean":             Manichaean,
	"Marchen":                Marchen,
	"Masaram_Gondi":          Masaram_Gondi,
	"Medefaidrin":            Medefaidrin,
	"Meetei_Mayek":           Meetei_Mayek,
	"Mende_Kikakui":          Mende_Kikakui,
	"Meroitic_Cursive":       Meroitic_Cursive,
	"Meroitic_Hieroglyphs":   Meroitic_Hieroglyphs,
	"Miao":                   Miao,
	"Modi":                   Modi,
	"Mongolian":              Mongolian,
	"Mro":                    Mro,
	"Multani":                Multani,
	"Myanmar":                Myanmar,
	"Nabataean":              Nabataean,
	"Nandinagari":            Nandinagari,
	"New_Tai_Lue":            New_Tai_Lue,
	"Newa":                   Newa,
	"Nko":                    Nko,
	"Nushu":                  Nushu,
	"Nyiakeng_Puachue_Hmong": Nyiakeng_Puachue_Hmong,
	"Ogham":                  Ogham,
	"Ol_Chiki":               Ol_Chiki,
	"Old_Hungarian":          Old_Hungarian,
	"Old_Italic":             Old_Italic,
	"Old_North_Arabian":      Old_North_Arabian,
	"Old_Permic":             Old_Permic,
	"Old_Persian":            Old_Persian,
	"Old_Sogdian":            Old_Sogdian,
	"Old_South_Arabian":      Old_South_Arabian,
	"Old_Turkic":             Old_Turkic,
	"Oriya":                  Oriya,
	"Osage":                  Osage,
	"Osmanya":                Osmanya,
	"Pahawh_Hmong":           Pahawh_Hmong,
	"Palmyrene":              Palmyrene,
	"Pau_Cin_Hau":            Pau_Cin_Hau,
	"Phags_Pa":               Phags_Pa,
	"Phoenician":             Phoenician,
	"Psalter_Pahlavi":        Psalter_Pahlavi,
	"Rejang":                 Rejang,
	"Runic":                  Runic,
	"Samaritan":              Samaritan,
	"Saurashtra":             Saurashtra,
	"Sharada":                Sharada,
	"Shavian":                Shavian,
	"Siddham":                Siddham,
	"SignWriting":            SignWriting,
	"Sinhala":                Sinhala,
	"Sogdian":                Sogdian,
	"Sora_Sompeng":           Sora_Sompeng,
	"Soyombo":                Soyombo,
	"Sundanese":              Sundanese,
	"Syloti_Nagri":           Syloti_Nagri,
	"Syriac":                 Syriac,
	"Tagalog":                Tagalog,
	"Tagbanwa":               Tagbanwa,
	"Tai_Le":                 Tai_Le,
	"Tai_Tham":               Tai_Tham,
	"Tai_Viet":               Tai_Viet,
	"Takri":                  Takri,
	"Tamil":                  Tamil,
	"Tangut":                 Tangut,
	"Telugu":                 Telugu,
	"Thaana":                 Thaana,
	"Thai":                   Thai,
	"Tibetan":                Tibetan,
	"Tifinagh":               Tifinagh,
	"Tirhuta":                Tirhuta,
	"Ugaritic":               Ugaritic,
	"Vai":                    Vai,
	"Wancho":                 Wancho,
	"Warang_Citi":            Warang_Citi,
	"Yezidi":                 Yezidi,
	"Yi":                     Yi,
	"Zanabazar_Square":       Zanabazar_Square,
}`,
			variables: []testVariable{
				{
					name: "Scripts",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "SpecialCase",
			comments: ``, // no comments for this block of variables
			source:   `var AzeriCase SpecialCase = _TurkishCase`,
			variables: []testVariable{
				{
					name: "AzeriCase",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "SpecialCase",
			comments: ``, // no comments for this block of variables
			source:   `var TurkishCase SpecialCase = _TurkishCase`,
			variables: []testVariable{
				{
					name: "TurkishCase",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
	},
	functions: []testFunction{
		{
			name: "In",
			comments: `In reports whether the rune is a member of one of the ranges.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
				{
					s:        "ranges ...*RangeTable",
					name:     "ranges",
					typeName: "...*RangeTable",
					pointer:  true,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "Is",
			comments: `Is reports whether the rune is in the specified table of ranges.
`,
			inputs: []testParameter{
				{
					s:        "rangeTab *RangeTable",
					name:     "rangeTab",
					typeName: "*RangeTable",
					pointer:  true,
				},
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsControl",
			comments: `IsControl reports whether the rune is a control character. The C (Other) Unicode category includes more code points such as surrogates; use Is(C, r) to test for them.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsDigit",
			comments: `IsDigit reports whether the rune is a decimal digit.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsGraphic",
			comments: `IsGraphic reports whether the rune is defined as a Graphic by Unicode. Such characters include letters, marks, numbers, punctuation, symbols, and spaces, from categories L, M, N, P, S, Zs.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsLetter",
			comments: `IsLetter reports whether the rune is a letter (category L).
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsLower",
			comments: `IsLower reports whether the rune is a lower case letter.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsMark",
			comments: `IsMark reports whether the rune is a mark character (category M).
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsNumber",
			comments: `IsNumber reports whether the rune is a number (category N).
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsOneOf",
			comments: `IsOneOf reports whether the rune is a member of one of the ranges. The function "In" provides a nicer signature and should be used in preference to IsOneOf.
`,
			inputs: []testParameter{
				{
					s:        "ranges []*RangeTable",
					name:     "ranges",
					typeName: "[]*RangeTable",
					pointer:  false,
				},
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsPrint",
			comments: `IsPrint reports whether the rune is defined as printable by Go. Such characters include letters, marks, numbers, punctuation, symbols, and the ASCII space character, from categories L, M, N, P, S and the ASCII space character. This categorization is the same as IsGraphic except that the only spacing character is ASCII space, U+0020.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsPunct",
			comments: `IsPunct reports whether the rune is a Unicode punctuation character (category P).
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsSpace",
			comments: `IsSpace reports whether the rune is a space character as defined by Unicode's White Space property; in the Latin-1 space this is

	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).

Other definitions of spacing characters are set by category Z and property Pattern_White_Space.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsSymbol",
			comments: `IsSymbol reports whether the rune is a symbolic character.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsTitle",
			comments: `IsTitle reports whether the rune is a title case letter.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "IsUpper",
			comments: `IsUpper reports whether the rune is an upper case letter.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "bool",
					name:     "",
					typeName: "bool",
					pointer:  false,
				},
			},
		},
		{
			name: "SimpleFold",
			comments: `SimpleFold iterates over Unicode code points equivalent under the Unicode-defined simple case folding. Among the code points equivalent to rune (including rune itself), SimpleFold returns the smallest rune > r if one exists, or else the smallest rune >= 0. If r is not a valid Unicode code point, SimpleFold(r) returns r.

For example:

	SimpleFold('A') = 'a'
	SimpleFold('a') = 'A'

	SimpleFold('K') = 'k'
	SimpleFold('k') = '\u212A' (Kelvin symbol, K)
	SimpleFold('\u212A') = 'K'

	SimpleFold('1') = '1'

	SimpleFold(-2) = -2
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "rune",
					name:     "",
					typeName: "rune",
					pointer:  false,
				},
			},
		},
		{
			name: "To",
			comments: `To maps the rune to the specified case: UpperCase, LowerCase, or TitleCase.
`,
			inputs: []testParameter{
				{
					s:        "_case int",
					name:     "_case",
					typeName: "int",
					pointer:  false,
				},
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "rune",
					name:     "",
					typeName: "rune",
					pointer:  false,
				},
			},
		},
		{
			name: "ToLower",
			comments: `ToLower maps the rune to lower case.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "rune",
					name:     "",
					typeName: "rune",
					pointer:  false,
				},
			},
		},
		{
			name: "ToTitle",
			comments: `ToTitle maps the rune to title case.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "rune",
					name:     "",
					typeName: "rune",
					pointer:  false,
				},
			},
		},
		{
			name: "ToUpper",
			comments: `ToUpper maps the rune to upper case.
`,
			inputs: []testParameter{
				{
					s:        "r rune",
					name:     "r",
					typeName: "rune",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "rune",
					name:     "",
					typeName: "rune",
					pointer:  false,
				},
			},
		},
	},
	types: []testType{
		{
			name:     "CaseRange",
			typeName: "struct",
			source: `type CaseRange struct {
	Lo    uint32
	Hi    uint32
	Delta d
}`,
			comments: `CaseRange represents a range of Unicode code points for simple (one code point to one code point) case conversion. The range runs from Lo to Hi inclusive, with a fixed stride of 1. Deltas are the number to add to the code point to reach the code point for a different case for that character. They may be negative. If zero, it means the character is in the corresponding case. There is a special case representing sequences of alternating corresponding Upper and Lower pairs. It appears with a fixed Delta of

	{UpperLower, UpperLower, UpperLower}

The constant UpperLower has an otherwise impossible delta value.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Range16",
			typeName: "struct",
			source: `type Range16 struct {
	Lo     uint16
	Hi     uint16
	Stride uint16
}`,
			comments: `Range16 represents of a range of 16-bit Unicode code points. The range runs from Lo to Hi inclusive and has the specified stride.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Range32",
			typeName: "struct",
			source: `type Range32 struct {
	Lo     uint32
	Hi     uint32
	Stride uint32
}`,
			comments: `Range32 represents of a range of Unicode code points and is used when one or more of the values will not fit in 16 bits. The range runs from Lo to Hi inclusive and has the specified stride. Lo and Hi must always be >= 1<<16.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "RangeTable",
			typeName: "struct",
			source: `type RangeTable struct {
	R16         []Range16
	R32         []Range32
	LatinOffset int // number of entries in R16 with Hi <= MaxLatin1
}`,
			comments: `RangeTable defines a set of Unicode code points by listing the ranges of code points within the set. The ranges are listed in two slices to save space: a slice of 16-bit ranges and a slice of 32-bit ranges. The two slices must be in sorted order and non-overlapping. Also, R32 should contain only values >= 0x10000 (1<<16).
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "SpecialCase",
			typeName: "[]CaseRange",
			source:   `type SpecialCase []CaseRange`,
			comments: `SpecialCase represents language-specific case mappings such as Turkish. Methods of SpecialCase customize (by overriding) the standard mappings.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name: "ToLower",
					comments: `ToLower maps the rune to lower case giving priority to the special mapping.
`,
					receiver: testParameter{
						s:        "special SpecialCase",
						name:     "special",
						typeName: "SpecialCase",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "r rune",
							name:     "r",
							typeName: "rune",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "rune",
							name:     "",
							typeName: "rune",
							pointer:  false,
						},
					},
				},
				{
					name: "ToTitle",
					comments: `ToTitle maps the rune to title case giving priority to the special mapping.
`,
					receiver: testParameter{
						s:        "special SpecialCase",
						name:     "special",
						typeName: "SpecialCase",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "r rune",
							name:     "r",
							typeName: "rune",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "rune",
							name:     "",
							typeName: "rune",
							pointer:  false,
						},
					},
				},
				{
					name: "ToUpper",
					comments: `ToUpper maps the rune to upper case giving priority to the special mapping.
`,
					receiver: testParameter{
						s:        "special SpecialCase",
						name:     "special",
						typeName: "SpecialCase",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "r rune",
							name:     "r",
							typeName: "rune",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "rune",
							name:     "",
							typeName: "rune",
							pointer:  false,
						},
					},
				},
			},
		},
	},
}

// Structure for package "net/rpc".
var pkgNetRPC = testPackage{
	name:       "rpc",
	importPath: "net/rpc",
	comments: `Package rpc provides access to the exported methods of an object across a network or other I/O connection. A server registers an object, making it visible as a service with the name of the type of the object. After registration, exported methods of the object will be accessible remotely. A server may register multiple objects (services) of different types but it is an error to register multiple objects of the same type.

Only methods that satisfy these criteria will be made available for remote access; other methods will be ignored:

	- the method's type is exported.
	- the method is exported.
	- the method has two arguments, both exported (or builtin) types.
	- the method's second argument is a pointer.
	- the method has return type error.

In effect, the method must look schematically like

	func (t *T) MethodName(argType T1, replyType *T2) error

where T1 and T2 can be marshaled by encoding/gob. These requirements apply even if a different codec is used. (In the future, these requirements may soften for custom codecs.)

The method's first argument represents the arguments provided by the caller; the second argument represents the result parameters to be returned to the caller. The method's return value, if non-nil, is passed back as a string that the client sees as if created by errors.New. If an error is returned, the reply parameter will not be sent back to the client.

The server may handle requests on a single connection by calling ServeConn. More typically it will create a network listener and call Accept or, for an HTTP listener, HandleHTTP and http.Serve.

A client wishing to use the service establishes a connection and then invokes NewClient on the connection. The convenience function Dial (DialHTTP) performs both steps for a raw network connection (an HTTP connection). The resulting Client object has two methods, Call and Go, that specify the service and method to call, a pointer containing the arguments, and a pointer to receive the result parameters.

The Call method waits for the remote call to complete while the Go method launches the call asynchronously and signals completion using the Call structure's Done channel.

Unless an explicit codec is set up, package encoding/gob is used to transport the data.

Here is a simple example. A server wishes to export an object of type Arith:

	package server

	import "errors"

	type Args struct {
		A, B int
	}

	type Quotient struct {
		Quo, Rem int
	}

	type Arith int

	func (t *Arith) Multiply(args *Args, reply *int) error {
		*reply = args.A * args.B
		return nil
	}

	func (t *Arith) Divide(args *Args, quo *Quotient) error {
		if args.B == 0 {
			return errors.New("divide by zero")
		}
		quo.Quo = args.A / args.B
		quo.Rem = args.A % args.B
		return nil
	}

The server calls (for HTTP service):

	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

At this point, clients can see a service "Arith" with methods "Arith.Multiply" and "Arith.Divide". To invoke one, a client first dials the server:

	client, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

Then it can make a remote call:

	// Synchronous call
	args := &server.Args{7,8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

or

	// Asynchronous call
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done	// will be equal to divCall
	// check errors, print, etc.

A server implementation will often provide a simple, type-safe wrapper for the client.

The net/rpc package is frozen and is not accepting new features.
`,
	files: []string{
		"client.go", "debug.go", "server.go",
	},
	testFiles: []string{
		"client_test.go", "server_test.go",
	},
	subdirectories: []string{"jsonrpc"}, // no subdirectories in this package
	imports: []string{
		"bufio", "encoding/gob", "errors", "fmt", "go/token", "html/template", "io", "log", "net",
		"net/http", "reflect", "sort", "strings", "sync",
	},
	testImports: []string{
		"errors", "fmt", "io", "log", "net", "net/http/httptest", "reflect", "runtime", "strings",
		"sync", "sync/atomic", "testing", "time",
	},
	constantBlocks: []testConstantBlock{
		{
			typeName: "", // no general type for this block of constants
			comments: ``, // no comments for this block of constants
			source: `const (
	// Defaults used by HandleHTTP
	DefaultRPCPath   = "/_goRPC_"
	DefaultDebugPath = "/debug/rpc"
)`,
			constants: []testConstant{
				{
					name: "DefaultRPCPath",
				},
				{
					name: "DefaultDebugPath",
				},
			},
		},
	},
	variableBlocks: []testVariableBlock{
		{
			typeName: "",
			comments: `DefaultServer is the default instance of *Server.
`,
			source: `var DefaultServer = NewServer()`,
			variables: []testVariable{
				{
					name: "DefaultServer",
				},
			},
			errors: []testError{}, // no errors in this block of variables
		},
		{
			typeName: "",
			comments: ``, // no comments for this block of variables
			source:   `var ErrShutdown = errors.New("connection is shut down")`,
			variables: []testVariable{
				{
					name: "ErrShutdown",
				},
			},
			errors: []testError{
				{
					name: "ErrShutdown",
				},
			},
		},
	},
	functions: []testFunction{
		{
			name: "Accept",
			comments: `Accept accepts connections on the listener and serves requests to DefaultServer for each incoming connection. Accept blocks; the caller typically invokes it in a go statement.
`,
			inputs: []testParameter{
				{
					s:        "lis net.Listener",
					name:     "lis",
					typeName: "net.Listener",
					pointer:  false,
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name: "HandleHTTP",
			comments: `HandleHTTP registers an HTTP handler for RPC messages to DefaultServer on DefaultRPCPath and a debugging handler on DefaultDebugPath. It is still necessary to invoke http.Serve(), typically in a go statement.
`,
			inputs:  []testParameter{}, // no inputs for this function
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name: "Register",
			comments: `Register publishes the receiver's methods in the DefaultServer.
`,
			inputs: []testParameter{
				{
					s:        "rcvr interface{}",
					name:     "rcvr",
					typeName: "interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "RegisterName",
			comments: `RegisterName is like Register but uses the provided name for the type instead of the receiver's concrete type.
`,
			inputs: []testParameter{
				{
					s:        "name string",
					name:     "name",
					typeName: "string",
					pointer:  false,
				},
				{
					s:        "rcvr interface{}",
					name:     "rcvr",
					typeName: "interface{}",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
		{
			name: "ServeCodec",
			comments: `ServeCodec is like ServeConn but uses the specified codec to decode requests and encode responses.
`,
			inputs: []testParameter{
				{
					s:        "codec ServerCodec",
					name:     "codec",
					typeName: "ServerCodec",
					pointer:  false,
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name: "ServeConn",
			comments: `ServeConn runs the DefaultServer on a single connection. ServeConn blocks, serving the connection until the client hangs up. The caller typically invokes ServeConn in a go statement. ServeConn uses the gob wire format (see package gob) on the connection. To use an alternate codec, use ServeCodec. See NewClient's comment for information about concurrent access.
`,
			inputs: []testParameter{
				{
					s:        "conn io.ReadWriteCloser",
					name:     "conn",
					typeName: "io.ReadWriteCloser",
					pointer:  false,
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name: "ServeRequest",
			comments: `ServeRequest is like ServeCodec but synchronously serves a single request. It does not close the codec upon completion.
`,
			inputs: []testParameter{
				{
					s:        "codec ServerCodec",
					name:     "codec",
					typeName: "ServerCodec",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "error",
					name:     "",
					typeName: "error",
					pointer:  false,
				},
			},
		},
	},
	types: []testType{
		{
			name:     "Call",
			typeName: "struct",
			source: `type Call struct {
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Receives *Call when Go is complete.
}`,
			comments: `Call represents an active RPC.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Client",
			typeName: "struct",
			source: `type Client struct {
	// contains filtered or unexported fields
}`,
			comments: `Client represents an RPC Client. There may be multiple outstanding Calls associated with a single Client, and a Client may be used by multiple goroutines simultaneously.
`,
			functions: []testFunction{
				{
					name: "Dial",
					comments: `Dial connects to an RPC server at the specified network address.
`,
					inputs: []testParameter{
						{
							s:        "network string",
							name:     "network",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "address string",
							name:     "address",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Client",
							name:     "",
							typeName: "*Client",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "DialHTTP",
					comments: `DialHTTP connects to an HTTP RPC server at the specified network address listening on the default HTTP RPC path.
`,
					inputs: []testParameter{
						{
							s:        "network string",
							name:     "network",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "address string",
							name:     "address",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Client",
							name:     "",
							typeName: "*Client",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "DialHTTPPath",
					comments: `DialHTTPPath connects to an HTTP RPC server at the specified network address and path.
`,
					inputs: []testParameter{
						{
							s:        "network string",
							name:     "network",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "address string",
							name:     "address",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "path string",
							name:     "path",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Client",
							name:     "",
							typeName: "*Client",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "NewClient",
					comments: `NewClient returns a new Client to handle requests to the set of services at the other end of the connection. It adds a buffer to the write side of the connection so the header and payload are sent as a unit.

The read and write halves of the connection are serialized independently, so no interlocking is required. However each half may be accessed concurrently so the implementation of conn should protect against concurrent reads or concurrent writes.
`,
					inputs: []testParameter{
						{
							s:        "conn io.ReadWriteCloser",
							name:     "conn",
							typeName: "io.ReadWriteCloser",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Client",
							name:     "",
							typeName: "*Client",
							pointer:  true,
						},
					},
				},
				{
					name: "NewClientWithCodec",
					comments: `NewClientWithCodec is like NewClient but uses the specified codec to encode requests and decode responses.
`,
					inputs: []testParameter{
						{
							s:        "codec ClientCodec",
							name:     "codec",
							typeName: "ClientCodec",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Client",
							name:     "",
							typeName: "*Client",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Call",
					comments: `Call invokes the named function, waits for it to complete, and returns its error status.
`,
					receiver: testParameter{
						s:        "client *Client",
						name:     "client",
						typeName: "*Client",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "serviceMethod string",
							name:     "serviceMethod",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "args interface{}",
							name:     "args",
							typeName: "interface{}",
							pointer:  false,
						},
						{
							s:        "reply interface{}",
							name:     "reply",
							typeName: "interface{}",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Close",
					comments: `Close calls the underlying codec's Close method. If the connection is already shutting down, ErrShutdown is returned.
`,
					receiver: testParameter{
						s:        "client *Client",
						name:     "client",
						typeName: "*Client",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Go",
					comments: `Go invokes the function asynchronously. It returns the Call structure representing the invocation. The done channel will signal when the call is complete by returning the same Call object. If done is nil, Go will allocate a new channel. If non-nil, done must be buffered or Go will deliberately crash.
`,
					receiver: testParameter{
						s:        "client *Client",
						name:     "client",
						typeName: "*Client",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "serviceMethod string",
							name:     "serviceMethod",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "args interface{}",
							name:     "args",
							typeName: "interface{}",
							pointer:  false,
						},
						{
							s:        "reply interface{}",
							name:     "reply",
							typeName: "interface{}",
							pointer:  false,
						},
						{
							s:        "done chan *Call",
							name:     "done",
							typeName: "chan *Call",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Call",
							name:     "",
							typeName: "*Call",
							pointer:  true,
						},
					},
				},
			},
		},
		{
			name:     "ClientCodec",
			typeName: "interface",
			source: `type ClientCodec interface {
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error

	Close() error
}`,
			comments: `A ClientCodec implements writing of RPC requests and reading of RPC responses for the client side of an RPC session. The client calls WriteRequest to write a request to the connection and calls ReadResponseHeader and ReadResponseBody in pairs to read responses. The client calls Close when finished with the connection. ReadResponseBody may be called with a nil argument to force the body of the response to be read and then discarded. See NewClient's comment for information about concurrent access.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Request",
			typeName: "struct",
			source: `type Request struct {
	ServiceMethod string // format: "Service.Method"
	Seq           uint64 // sequence number chosen by client
	// contains filtered or unexported fields
}`,
			comments: `Request is a header written before every RPC call. It is used internally but documented here as an aid to debugging, such as when analyzing network traffic.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Response",
			typeName: "struct",
			source: `type Response struct {
	ServiceMethod string // echoes that of the Request
	Seq           uint64 // echoes that of the request
	Error         string // error, if any.
	// contains filtered or unexported fields
}`,
			comments: `Response is a header written before every RPC return. It is used internally but documented here as an aid to debugging, such as when analyzing network traffic.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Server",
			typeName: "struct",
			source: `type Server struct {
	// contains filtered or unexported fields
}`,
			comments: `Server represents an RPC Server.
`,
			functions: []testFunction{
				{
					name: "NewServer",
					comments: `NewServer returns a new Server.
`,
					inputs: []testParameter{}, // no inputs for this function
					outputs: []testParameter{
						{
							s:        "*Server",
							name:     "",
							typeName: "*Server",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Accept",
					comments: `Accept accepts connections on the listener and serves requests for each incoming connection. Accept blocks until the listener returns a non-nil error. The caller typically invokes Accept in a go statement.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "lis net.Listener",
							name:     "lis",
							typeName: "net.Listener",
							pointer:  false,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "HandleHTTP",
					comments: `HandleHTTP registers an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath. It is still necessary to invoke http.Serve(), typically in a go statement.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "rpcPath string",
							name:     "rpcPath",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "debugPath string",
							name:     "debugPath",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "Register",
					comments: `Register publishes in the server the set of methods of the receiver value that satisfy the following conditions:

	- exported method of exported type
	- two arguments, both of exported type
	- the second argument is a pointer
	- one return value, of type error

It returns an error if the receiver is not an exported type or has no suitable methods. It also logs the error using package log. The client accesses each method using a string of the form "Type.Method", where Type is the receiver's concrete type.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "rcvr interface{}",
							name:     "rcvr",
							typeName: "interface{}",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "RegisterName",
					comments: `RegisterName is like Register but uses the provided name for the type instead of the receiver's concrete type.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "name string",
							name:     "name",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "rcvr interface{}",
							name:     "rcvr",
							typeName: "interface{}",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "ServeCodec",
					comments: `ServeCodec is like ServeConn but uses the specified codec to decode requests and encode responses.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "codec ServerCodec",
							name:     "codec",
							typeName: "ServerCodec",
							pointer:  false,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "ServeConn",
					comments: `ServeConn runs the server on a single connection. ServeConn blocks, serving the connection until the client hangs up. The caller typically invokes ServeConn in a go statement. ServeConn uses the gob wire format (see package gob) on the connection. To use an alternate codec, use ServeCodec. See NewClient's comment for information about concurrent access.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "conn io.ReadWriteCloser",
							name:     "conn",
							typeName: "io.ReadWriteCloser",
							pointer:  false,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "ServeHTTP",
					comments: `ServeHTTP implements an http.Handler that answers RPC requests.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "w http.ResponseWriter",
							name:     "w",
							typeName: "http.ResponseWriter",
							pointer:  false,
						},
						{
							s:        "req *http.Request",
							name:     "req",
							typeName: "*http.Request",
							pointer:  true,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "ServeRequest",
					comments: `ServeRequest is like ServeCodec but synchronously serves a single request. It does not close the codec upon completion.
`,
					receiver: testParameter{
						s:        "server *Server",
						name:     "server",
						typeName: "*Server",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "codec ServerCodec",
							name:     "codec",
							typeName: "ServerCodec",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "ServerCodec",
			typeName: "interface",
			source: `type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(interface{}) error
	WriteResponse(*Response, interface{}) error

	// Close can be called multiple times and must be idempotent.
	Close() error
}`,
			comments: `A ServerCodec implements reading of RPC requests and writing of RPC responses for the server side of an RPC session. The server calls ReadRequestHeader and ReadRequestBody in pairs to read requests from the connection, and it calls WriteResponse to write a response back. The server calls Close when finished with the connection. ReadRequestBody may be called with a nil argument to force the body of the request to be read and discarded. See NewClient's comment for information about concurrent access.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "ServerError",
			typeName: "string",
			source:   `type ServerError string`,
			comments: `ServerError represents an error that has been returned from the remote side of the RPC connection.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:     "Error",
					comments: ``, // no comments for this method
					receiver: testParameter{
						s:        "e ServerError",
						name:     "e",
						typeName: "ServerError",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
	},
}

// Structure for package "time".
var pkgTime = testPackage{
	name:       "time",
	importPath: "time",
	comments: `Package time provides functionality for measuring and displaying time.

The calendrical calculations always assume a Gregorian calendar, with no leap seconds.


Monotonic Clocks

Operating systems provide both a “wall clock,” which is subject to changes for clock synchronization, and a “monotonic clock,” which is not. The general rule is that the wall clock is for telling time and the monotonic clock is for measuring time. Rather than split the API, in this package the Time returned by time.Now contains both a wall clock reading and a monotonic clock reading; later time-telling operations use the wall clock reading, but later time-measuring operations, specifically comparisons and subtractions, use the monotonic clock reading.

For example, this code always computes a positive elapsed time of approximately 20 milliseconds, even if the wall clock is changed during the operation being timed:

	start := time.Now()
	... operation that takes 20 milliseconds ...
	t := time.Now()
	elapsed := t.Sub(start)

Other idioms, such as time.Since(start), time.Until(deadline), and time.Now().Before(deadline), are similarly robust against wall clock resets.

The rest of this section gives the precise details of how operations use monotonic clocks, but understanding those details is not required to use this package.

The Time returned by time.Now contains a monotonic clock reading. If Time t has a monotonic clock reading, t.Add adds the same duration to both the wall clock and monotonic clock readings to compute the result. Because t.AddDate(y, m, d), t.Round(d), and t.Truncate(d) are wall time computations, they always strip any monotonic clock reading from their results. Because t.In, t.Local, and t.UTC are used for their effect on the interpretation of the wall time, they also strip any monotonic clock reading from their results. The canonical way to strip a monotonic clock reading is to use t = t.Round(0).

If Times t and u both contain monotonic clock readings, the operations t.After(u), t.Before(u), t.Equal(u), and t.Sub(u) are carried out using the monotonic clock readings alone, ignoring the wall clock readings. If either t or u contains no monotonic clock reading, these operations fall back to using the wall clock readings.

On some systems the monotonic clock will stop if the computer goes to sleep. On such a system, t.Sub(u) may not accurately reflect the actual time that passed between t and u.

Because the monotonic clock reading has no meaning outside the current process, the serialized forms generated by t.GobEncode, t.MarshalBinary, t.MarshalJSON, and t.MarshalText omit the monotonic clock reading, and t.Format provides no format for it. Similarly, the constructors time.Date, time.Parse, time.ParseInLocation, and time.Unix, as well as the unmarshalers t.GobDecode, t.UnmarshalBinary. t.UnmarshalJSON, and t.UnmarshalText always create times with no monotonic clock reading.

Note that the Go == operator compares not just the time instant but also the Location and the monotonic clock reading. See the documentation for the Time type for a discussion of equality testing for Time values.

For debugging, the result of t.String does include the monotonic clock reading if present. If t != u because of different monotonic clock readings, that difference will be visible when printing t.String() and u.String().
`,
	files: []string{
		"embed.go", "format.go", "genzabbrs.go", "sleep.go", "sys_plan9.go", "sys_unix.go",
		"sys_windows.go", "tick.go", "time.go", "zoneinfo.go", "zoneinfo_abbrs_windows.go",
		"zoneinfo_android.go", "zoneinfo_ios.go", "zoneinfo_js.go", "zoneinfo_plan9.go",
		"zoneinfo_read.go", "zoneinfo_unix.go", "zoneinfo_windows.go",
	},
	testFiles: []string{
		"example_test.go", "export_android_test.go", "export_test.go", "export_windows_test.go",
		"format_test.go", "internal_test.go", "mono_test.go", "sleep_test.go", "tick_test.go",
		"time_test.go", "tzdata_test.go", "zoneinfo_android_test.go", "zoneinfo_test.go",
		"zoneinfo_unix_test.go", "zoneinfo_windows_test.go",
	},
	subdirectories: []string{
		"tzdata",
	},
	imports: []string{
		"errors", "runtime", "sync", "syscall", "unsafe",
	},
	testImports: []string{
		"bytes", "encoding/gob", "encoding/json", "errors", "fmt", "math/big", "math/rand", "os",
		"reflect", "runtime", "strconv", "strings", "sync", "sync/atomic", "testing",
		"testing/quick", "time", "time/tzdata",
	},
	constantBlocks: []testConstantBlock{
		{
			typeName: "",
			comments: `These are predefined layouts for use in Time.Format and time.Parse. The reference time used in the layouts is the specific time:

	Mon Jan 2 15:04:05 MST 2006

which is Unix time 1136239445. Since MST is GMT-0700, the reference time can be thought of as

	01/02 03:04:05PM '06 -0700

To define your own format, write down what the reference time would look like formatted your way; see the values of constants like ANSIC, StampMicro or Kitchen for examples. The model is to demonstrate what the reference time looks like so that the Format and Parse methods can apply the same transformation to a general time value.

Some valid layouts are invalid time values for time.Parse, due to formats such as _ for space padding and Z for zone information.

Within the format string, an underscore _ represents a space that may be replaced by a digit if the following number (a day) has two digits; for compatibility with fixed-width Unix time formats.

A decimal point followed by one or more zeros represents a fractional second, printed to the given number of decimal places. A decimal point followed by one or more nines represents a fractional second, printed to the given number of decimal places, with trailing zeros removed. When parsing (only), the input may contain a fractional second field immediately after the seconds field, even if the layout does not signify its presence. In that case a decimal point followed by a maximal series of digits is parsed as a fractional second.

Numeric time zone offsets format as follows:

	-0700  ±hhmm
	-07:00 ±hh:mm
	-07    ±hh

Replacing the sign in the format with a Z triggers the ISO 8601 behavior of printing Z instead of an offset for the UTC zone. Thus:

	Z0700  Z or ±hhmm
	Z07:00 Z or ±hh:mm
	Z07    Z or ±hh

The recognized day of week formats are "Mon" and "Monday". The recognized month formats are "Jan" and "January".

The formats 2, _2, and 02 are unpadded, space-padded, and zero-padded day of month. The formats __2 and 002 are space-padded and zero-padded three-character day of year; there is no unpadded day of year format.

Text in the format string that is not recognized as part of the reference time is echoed verbatim during Format and expected to appear verbatim in the input to Parse.

The executable example for Time.Format demonstrates the working of the layout string in detail and is a good reference.

Note that the RFC822, RFC850, and RFC1123 formats should be applied only to local times. Applying them to UTC times will use "UTC" as the time zone abbreviation, while strictly speaking those RFCs require the use of "GMT" in that case. In general RFC1123Z should be used instead of RFC1123 for servers that insist on that format, and RFC3339 should be preferred for new protocols. RFC3339, RFC822, RFC822Z, RFC1123, and RFC1123Z are useful for formatting; when used with time.Parse they do not accept all the time formats permitted by the RFCs and they do accept time formats not formally defined. The RFC3339Nano format removes trailing zeros from the seconds field and thus may not sort correctly once formatted.
`,
			source: `const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)`,
			constants: []testConstant{
				{
					name: "ANSIC",
				},
				{
					name: "UnixDate",
				},
				{
					name: "RubyDate",
				},
				{
					name: "RFC822",
				},
				{
					name: "RFC822Z",
				},
				{
					name: "RFC850",
				},
				{
					name: "RFC1123",
				},
				{
					name: "RFC1123Z",
				},
				{
					name: "RFC3339",
				},
				{
					name: "RFC3339Nano",
				},
				{
					name: "Kitchen",
				},
				{
					name: "Stamp",
				},
				{
					name: "StampMilli",
				},
				{
					name: "StampMicro",
				},
				{
					name: "StampNano",
				},
			},
		},
		{
			typeName: "",
			comments: `Common durations. There is no definition for units of Day or larger to avoid confusion across daylight savings time zone transitions.

To count the number of units in a Duration, divide:

	second := time.Second
	fmt.Print(int64(second/time.Millisecond)) // prints 1000

To convert an integer number of units to a Duration, multiply:

	seconds := 10
	fmt.Print(time.Duration(seconds)*time.Second) // prints 10s
`,
			source: `const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)`,
			constants: []testConstant{
				{
					name: "Nanosecond",
				},
				{
					name: "Microsecond",
				},
				{
					name: "Millisecond",
				},
				{
					name: "Second",
				},
				{
					name: "Minute",
				},
				{
					name: "Hour",
				},
			},
		},
		{
			typeName: "Month",
			comments: ``, // no comments for this block of constants
			source: `const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)`,
			constants: []testConstant{
				{
					name: "January",
				},
				{
					name: "February",
				},
				{
					name: "March",
				},
				{
					name: "April",
				},
				{
					name: "May",
				},
				{
					name: "June",
				},
				{
					name: "July",
				},
				{
					name: "August",
				},
				{
					name: "September",
				},
				{
					name: "October",
				},
				{
					name: "November",
				},
				{
					name: "December",
				},
			},
		},
		{
			typeName: "Weekday",
			comments: ``, // no comments for this block of constants
			source: `const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)`,
			constants: []testConstant{
				{
					name: "Sunday",
				},
				{
					name: "Monday",
				},
				{
					name: "Tuesday",
				},
				{
					name: "Wednesday",
				},
				{
					name: "Thursday",
				},
				{
					name: "Friday",
				},
				{
					name: "Saturday",
				},
			},
		},
	},
	variableBlocks: []testVariableBlock{
		{
			typeName: "Location",
			comments: `Local represents the system's local time zone. On Unix systems, Local consults the TZ environment variable to find the time zone to use. No TZ means use the system default /etc/localtime. TZ="" means use UTC. TZ="foo" means use file foo in the system timezone directory.
`,
			source: `var Local *Location = &localLoc`,
			variables: []testVariable{
				{
					name: "Local",
				},
			},
			errors: []testError{}, // no errors in the block of variables
		},
		{
			typeName: "Location",
			comments: `UTC represents Universal Coordinated Time (UTC).
`,
			source: `var UTC *Location = &utcLoc`,
			variables: []testVariable{
				{
					name: "UTC",
				},
			},
			errors: []testError{}, // no errors in the block of variables
		},
	},
	functions: []testFunction{
		{
			name: "After",
			comments: `After waits for the duration to elapse and then sends the current time on the returned channel. It is equivalent to NewTimer(d).C. The underlying Timer is not recovered by the garbage collector until the timer fires. If efficiency is a concern, use NewTimer instead and call Timer.Stop if the timer is no longer needed.
`,
			inputs: []testParameter{
				{
					s:        "d Duration",
					name:     "d",
					typeName: "Duration",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "<-chan Time",
					name:     "",
					typeName: "<-chan Time",
					pointer:  false,
				},
			},
		},
		{
			name: "Sleep",
			comments: `Sleep pauses the current goroutine for at least the duration d. A negative or zero duration causes Sleep to return immediately.
`,
			inputs: []testParameter{
				{
					s:        "d Duration",
					name:     "d",
					typeName: "Duration",
					pointer:  false,
				},
			},
			outputs: []testParameter{}, // no outputs for this function
		},
		{
			name: "Tick",
			comments: `Tick is a convenience wrapper for NewTicker providing access to the ticking channel only. While Tick is useful for clients that have no need to shut down the Ticker, be aware that without a way to shut it down the underlying Ticker cannot be recovered by the garbage collector; it "leaks". Unlike NewTicker, Tick will return nil if d <= 0.
`,
			inputs: []testParameter{
				{
					s:        "d Duration",
					name:     "d",
					typeName: "Duration",
					pointer:  false,
				},
			},
			outputs: []testParameter{
				{
					s:        "<-chan Time",
					name:     "",
					typeName: "<-chan Time",
					pointer:  false,
				},
			},
		},
	},
	types: []testType{
		{
			name:     "Duration",
			typeName: "int64",
			source:   `type Duration int64`,
			comments: `A Duration represents the elapsed time between two instants as an int64 nanosecond count. The representation limits the largest representable duration to approximately 290 years.
`,
			functions: []testFunction{
				{
					name: "ParseDuration",
					comments: `ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
`,
					inputs: []testParameter{
						{
							s:        "s string",
							name:     "s",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Since",
					comments: `Since returns the time elapsed since t. It is shorthand for time.Now().Sub(t).
`,
					inputs: []testParameter{
						{
							s:        "t Time",
							name:     "t",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
					},
				},
				{
					name: "Until",
					comments: `Until returns the duration until t. It is shorthand for t.Sub(time.Now()).
`,
					inputs: []testParameter{
						{
							s:        "t Time",
							name:     "t",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Hours",
					comments: `Hours returns the duration as a floating point number of hours.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "float64",
							name:     "",
							typeName: "float64",
							pointer:  false,
						},
					},
				},
				{
					name: "Microseconds",
					comments: `Microseconds returns the duration as an integer microsecond count.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int64",
							name:     "",
							typeName: "int64",
							pointer:  false,
						},
					},
				},
				{
					name: "Milliseconds",
					comments: `Milliseconds returns the duration as an integer millisecond count.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int64",
							name:     "",
							typeName: "int64",
							pointer:  false,
						},
					},
				},
				{
					name: "Minutes",
					comments: `Minutes returns the duration as a floating point number of minutes.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "float64",
							name:     "",
							typeName: "float64",
							pointer:  false,
						},
					},
				},
				{
					name: "Nanoseconds",
					comments: `Nanoseconds returns the duration as an integer nanosecond count.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int64",
							name:     "",
							typeName: "int64",
							pointer:  false,
						},
					},
				},
				{
					name: "Round",
					comments: `Round returns the result of rounding d to the nearest multiple of m. The rounding behavior for halfway values is to round away from zero. If the result exceeds the maximum (or minimum) value that can be stored in a Duration, Round returns the maximum (or minimum) duration. If m <= 0, Round returns d unchanged.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "m Duration",
							name:     "m",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
					},
				},
				{
					name: "Seconds",
					comments: `Seconds returns the duration as a floating point number of seconds.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "float64",
							name:     "",
							typeName: "float64",
							pointer:  false,
						},
					},
				},
				{
					name: "String",
					comments: `String returns a string representing the duration in the form "72h3m0.5s". Leading zero units are omitted. As a special case, durations less than one second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure that the leading digit is non-zero. The zero duration formats as 0s.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
				{
					name: "Truncate",
					comments: `Truncate returns the result of rounding d toward zero to a multiple of m. If m <= 0, Truncate returns d unchanged.
`,
					receiver: testParameter{
						s:        "d Duration",
						name:     "d",
						typeName: "Duration",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "m Duration",
							name:     "m",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Location",
			typeName: "struct",
			source: `type Location struct {
	// contains filtered or unexported fields
}`,
			comments: `A Location maps time instants to the zone in use at that time. Typically, the Location represents the collection of time offsets in use in a geographical area. For many Locations the time offset varies depending on whether daylight savings time is in use at the time instant.
`,
			functions: []testFunction{
				{
					name: "FixedZone",
					comments: `FixedZone returns a Location that always uses the given zone name and offset (seconds east of UTC).
`,
					inputs: []testParameter{
						{
							s:        "name string",
							name:     "name",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "offset int",
							name:     "offset",
							typeName: "int",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Location",
							name:     "",
							typeName: "*Location",
							pointer:  true,
						},
					},
				},
				{
					name: "LoadLocation",
					comments: `LoadLocation returns the Location with the given name.

If the name is "" or "UTC", LoadLocation returns UTC. If the name is "Local", LoadLocation returns Local.

Otherwise, the name is taken to be a location name corresponding to a file in the IANA Time Zone database, such as "America/New_York".

The time zone database needed by LoadLocation may not be present on all systems, especially non-Unix systems. LoadLocation looks in the directory or uncompressed zip file named by the ZONEINFO environment variable, if any, then looks in known installation locations on Unix systems, and finally looks in $GOROOT/lib/time/zoneinfo.zip.
`,
					inputs: []testParameter{
						{
							s:        "name string",
							name:     "name",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Location",
							name:     "",
							typeName: "*Location",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "LoadLocationFromTZData",
					comments: `LoadLocationFromTZData returns a Location with the given name initialized from the IANA Time Zone database-formatted data. The data should be in the format of a standard IANA time zone file (for example, the content of /etc/localtime on Unix systems).
`,
					inputs: []testParameter{
						{
							s:        "name string",
							name:     "name",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "data []byte",
							name:     "data",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Location",
							name:     "",
							typeName: "*Location",
							pointer:  true,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "String",
					comments: `String returns a descriptive name for the time zone information, corresponding to the name argument to LoadLocation or FixedZone.
`,
					receiver: testParameter{
						s:        "l *Location",
						name:     "l",
						typeName: "*Location",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Month",
			typeName: "int",
			source:   `type Month int`,
			comments: `A Month specifies a month of the year (January = 1, ...).
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name: "String",
					comments: `String returns the English name of the month ("January", "February", ...).
`,
					receiver: testParameter{
						s:        "m Month",
						name:     "m",
						typeName: "Month",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "ParseError",
			typeName: "struct",
			source: `type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
	Message    string
}`,
			comments: `ParseError describes a problem parsing a time string.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name: "Error",
					comments: `Error returns the string representation of a ParseError.
`,
					receiver: testParameter{
						s:        "e *ParseError",
						name:     "e",
						typeName: "*ParseError",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Ticker",
			typeName: "struct",
			source: `type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	// contains filtered or unexported fields
}`,
			comments: `A Ticker holds a channel that delivers “ticks” of a clock at intervals.
`,
			functions: []testFunction{
				{
					name: "NewTicker",
					comments: `NewTicker returns a new Ticker containing a channel that will send the time on the channel after each tick. The period of the ticks is specified by the duration argument. The ticker will adjust the time interval or drop ticks to make up for slow receivers. The duration d must be greater than zero; if not, NewTicker will panic. Stop the ticker to release associated resources.
`,
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Ticker",
							name:     "",
							typeName: "*Ticker",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Reset",
					comments: `Reset stops a ticker and resets its period to the specified duration. The next tick will arrive after the new period elapses.
`,
					receiver: testParameter{
						s:        "t *Ticker",
						name:     "t",
						typeName: "*Ticker",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{}, // no outputs for this method
				},
				{
					name: "Stop",
					comments: `Stop turns off a ticker. After Stop, no more ticks will be sent. Stop does not close the channel, to prevent a concurrent goroutine reading from the channel from seeing an erroneous "tick".
`,
					receiver: testParameter{
						s:        "t *Ticker",
						name:     "t",
						typeName: "*Ticker",
						pointer:  true,
					},
					inputs:  []testParameter{}, // no inputs for this method
					outputs: []testParameter{}, // no outputs for this method
				},
			},
		},
		{
			name:     "Time",
			typeName: "struct",
			source: `type Time struct {
	// contains filtered or unexported fields
}`,
			comments: `A Time represents an instant in time with nanosecond precision.

Programs using times should typically store and pass them as values, not pointers. That is, time variables and struct fields should be of type time.Time, not *time.Time.

A Time value can be used by multiple goroutines simultaneously except that the methods GobDecode, UnmarshalBinary, UnmarshalJSON and UnmarshalText are not concurrency-safe.

Time instants can be compared using the Before, After, and Equal methods. The Sub method subtracts two instants, producing a Duration. The Add method adds a Time and a Duration, producing a Time.

The zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC. As this time is unlikely to come up in practice, the IsZero method gives a simple way of detecting a time that has not been initialized explicitly.

Each Time has associated with it a Location, consulted when computing the presentation form of the time, such as in the Format, Hour, and Year methods. The methods Local, UTC, and In return a Time with a specific location. Changing the location in this way changes only the presentation; it does not change the instant in time being denoted and therefore does not affect the computations described in earlier paragraphs.

Representations of a Time value saved by the GobEncode, MarshalBinary, MarshalJSON, and MarshalText methods store the Time.Location's offset, but not the location name. They therefore lose information about Daylight Saving Time.

In addition to the required “wall clock” reading, a Time may contain an optional reading of the current process's monotonic clock, to provide additional precision for comparison or subtraction. See the “Monotonic Clocks” section in the package documentation for details.

Note that the Go == operator compares not just the time instant but also the Location and the monotonic clock reading. Therefore, Time values should not be used as map or database keys without first guaranteeing that the identical Location has been set for all values, which can be achieved through use of the UTC or Local method, and that the monotonic clock reading has been stripped by setting t = t.Round(0). In general, prefer t.Equal(u) to t == u, since t.Equal uses the most accurate comparison available and correctly handles the case when only one of its arguments has a monotonic clock reading.
`,
			functions: []testFunction{
				{
					name: "Date",
					comments: `Date returns the Time corresponding to

	yyyy-mm-dd hh:mm:ss + nsec nanoseconds

in the appropriate zone for that time in the given location.

The month, day, hour, min, sec, and nsec values may be outside their usual ranges and will be normalized during the conversion. For example, October 32 converts to November 1.

A daylight savings time transition skips or repeats times. For example, in the United States, March 13, 2011 2:15am never occurred, while November 6, 2011 1:15am occurred twice. In such cases, the choice of time zone, and therefore the time, is not well-defined. Date returns a time that is correct in one of the two zones involved in the transition, but it does not guarantee which.

Date panics if loc is nil.
`,
					inputs: []testParameter{
						{
							s:        "year int",
							name:     "year",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "month Month",
							name:     "month",
							typeName: "Month",
							pointer:  false,
						},
						{
							s:        "day int",
							name:     "day",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "hour int",
							name:     "hour",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "min int",
							name:     "min",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "sec int",
							name:     "sec",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "nsec int",
							name:     "nsec",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "loc *Location",
							name:     "loc",
							typeName: "*Location",
							pointer:  true,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "Now",
					comments: `Now returns the current local time.
`,
					inputs: []testParameter{}, // no inputs for this function
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "Parse",
					comments: `Parse parses a formatted string and returns the time value it represents. The layout defines the format by showing how the reference time, defined to be

	Mon Jan 2 15:04:05 -0700 MST 2006

would be interpreted if it were the value; it serves as an example of the input format. The same interpretation will then be made to the input string.

Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and convenient representations of the reference time. For more information about the formats and the definition of the reference time, see the documentation for ANSIC and the other constants defined by this package. Also, the executable example for Time.Format demonstrates the working of the layout string in detail and is a good reference.

Elements omitted from the value are assumed to be zero or, when zero is impossible, one, so parsing "3:04pm" returns the time corresponding to Jan 1, year 0, 15:04:00 UTC (note that because the year is 0, this time is before the zero Time). Years must be in the range 0000..9999. The day of the week is checked for syntax but it is otherwise ignored.

For layouts specifying the two-digit year 06, a value NN >= 69 will be treated as 19NN and a value NN < 69 will be treated as 20NN.

In the absence of a time zone indicator, Parse returns a time in UTC.

When parsing a time with a zone offset like -0700, if the offset corresponds to a time zone used by the current location (Local), then Parse uses that location and zone in the returned time. Otherwise it records the time as being in a fabricated location with time fixed at the given zone offset.

When parsing a time with a zone abbreviation like MST, if the zone abbreviation has a defined offset in the current location, then that offset is used. The zone abbreviation "UTC" is recognized as UTC regardless of location. If the zone abbreviation is unknown, Parse records the time as being in a fabricated location with the given zone abbreviation and a zero offset. This choice means that such a time can be parsed and reformatted with the same layout losslessly, but the exact instant used in the representation will differ by the actual zone offset. To avoid such problems, prefer time layouts that use a numeric zone offset, or use ParseInLocation.
`,
					inputs: []testParameter{
						{
							s:        "layout string",
							name:     "layout",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "value string",
							name:     "value",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "ParseInLocation",
					comments: `ParseInLocation is like Parse but differs in two important ways. First, in the absence of time zone information, Parse interprets a time as UTC; ParseInLocation interprets the time as in the given location. Second, when given a zone offset or abbreviation, Parse tries to match it against the Local location; ParseInLocation uses the given location.
`,
					inputs: []testParameter{
						{
							s:        "layout string",
							name:     "layout",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "value string",
							name:     "value",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "loc *Location",
							name:     "loc",
							typeName: "*Location",
							pointer:  true,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Unix",
					comments: `Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC. It is valid to pass nsec outside the range [0, 999999999]. Not all sec values have a corresponding time value. One such value is 1<<63-1 (the largest int64 value).
`,
					inputs: []testParameter{
						{
							s:        "sec int64",
							name:     "sec",
							typeName: "int64",
							pointer:  false,
						},
						{
							s:        "nsec int64",
							name:     "nsec",
							typeName: "int64",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Add",
					comments: `Add returns the time t+d.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "AddDate",
					comments: `AddDate returns the time corresponding to adding the given number of years, months, and days to t. For example, AddDate(-1, 2, 3) applied to January 1, 2011 returns March 4, 2010.

AddDate normalizes its result in the same way that Date does, so, for example, adding one month to October 31 yields December 1, the normalized form for November 31.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "years int",
							name:     "years",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "months int",
							name:     "months",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "days int",
							name:     "days",
							typeName: "int",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "After",
					comments: `After reports whether the time instant t is after u.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "u Time",
							name:     "u",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
				{
					name: "AppendFormat",
					comments: `AppendFormat is like Format but appends the textual representation to b and returns the extended buffer.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "b []byte",
							name:     "b",
							typeName: "[]byte",
							pointer:  false,
						},
						{
							s:        "layout string",
							name:     "layout",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "[]byte",
							name:     "",
							typeName: "[]byte",
							pointer:  false,
						},
					},
				},
				{
					name: "Before",
					comments: `Before reports whether the time instant t is before u.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "u Time",
							name:     "u",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
				{
					name: "Clock",
					comments: `Clock returns the hour, minute, and second within the day specified by t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "hour int",
							name:     "hour",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "min int",
							name:     "min",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "sec int",
							name:     "sec",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Date",
					comments: `Date returns the year, month, and day in which t occurs.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "year int",
							name:     "year",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "month Month",
							name:     "month",
							typeName: "Month",
							pointer:  false,
						},

						{
							s:        "day int",
							name:     "day",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Day",
					comments: `Day returns the day of the month specified by t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Equal",
					comments: `Equal reports whether t and u represent the same time instant. Two times can be equal even if they are in different locations. For example, 6:00 +0200 and 4:00 UTC are Equal. See the documentation on the Time type for the pitfalls of using == with Time values; most code should use Equal instead.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "u Time",
							name:     "u",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
				{
					name: "Format",
					comments: `Format returns a textual representation of the time value formatted according to layout, which defines the format by showing how the reference time, defined to be

	Mon Jan 2 15:04:05 -0700 MST 2006

would be displayed if it were the value; it serves as an example of the desired output. The same display rules will then be applied to the time value.

A fractional second is represented by adding a period and zeros to the end of the seconds section of layout string, as in "15:04:05.000" to format a time stamp with millisecond precision.

Predefined layouts ANSIC, UnixDate, RFC3339 and others describe standard and convenient representations of the reference time. For more information about the formats and the definition of the reference time, see the documentation for ANSIC and the other constants defined by this package.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "layout string",
							name:     "layout",
							typeName: "string",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
				{
					name: "GobDecode",
					comments: `GobDecode implements the gob.GobDecoder interface.
`,
					receiver: testParameter{
						s:        "t *Time",
						name:     "t",
						typeName: "*Time",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "data []byte",
							name:     "data",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "GobEncode",
					comments: `GobEncode implements the gob.GobEncoder interface.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "[]byte",
							name:     "",
							typeName: "[]byte",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Hour",
					comments: `Hour returns the hour within the day specified by t, in the range [0, 23].
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "ISOWeek",
					comments: `ISOWeek returns the ISO 8601 year and week number in which t occurs. Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "year int",
							name:     "year",
							typeName: "int",
							pointer:  false,
						},
						{
							s:        "week int",
							name:     "week",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "In",
					comments: `In returns a copy of t representing the same time instant, but with the copy's location information set to loc for display purposes.

In panics if loc is nil.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "loc *Location",
							name:     "loc",
							typeName: "*Location",
							pointer:  true,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "IsZero",
					comments: `IsZero reports whether t represents the zero time instant, January 1, year 1, 00:00:00 UTC.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
				{
					name: "Local",
					comments: `Local returns t with the location set to local time.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "Location",
					comments: `Location returns the time zone information associated with t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "*Location",
							name:     "",
							typeName: "*Location",
							pointer:  true,
						},
					},
				},
				{
					name: "MarshalBinary",
					comments: `MarshalBinary implements the encoding.BinaryMarshaler interface.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "[]byte",
							name:     "",
							typeName: "[]byte",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "MarshalJSON",
					comments: `MarshalJSON implements the json.Marshaler interface. The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "[]byte",
							name:     "",
							typeName: "[]byte",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "MarshalText",
					comments: `MarshalText implements the encoding.TextMarshaler interface. The time is formatted in RFC 3339 format, with sub-second precision added if present.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "[]byte",
							name:     "",
							typeName: "[]byte",
							pointer:  false,
						},
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Minute",
					comments: `Minute returns the minute offset within the hour specified by t, in the range [0, 59].
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Month",
					comments: `Month returns the month of the year specified by t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "Month",
							name:     "",
							typeName: "Month",
							pointer:  false,
						},
					},
				},
				{
					name: "Nanosecond",
					comments: `Nanosecond returns the nanosecond offset within the second specified by t, in the range [0, 999999999].
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Round",
					comments: `Round returns the result of rounding t to the nearest multiple of d (since the zero time). The rounding behavior for halfway values is to round up. If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.

Round operates on the time as an absolute duration since the zero time; it does not operate on the presentation form of the time. Thus, Round(Hour) may return a time with a non-zero minute, depending on the time's Location.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "Second",
					comments: `Second returns the second offset within the minute specified by t, in the range [0, 59].
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "String",
					comments: `String returns the time formatted using the format string

	"2006-01-02 15:04:05.999999999 -0700 MST"

If the time has a monotonic clock reading, the returned string includes a final field "m=±<value>", where value is the monotonic clock reading formatted as a decimal number of seconds.

The returned string is meant for debugging; for a stable serialized representation, use t.MarshalText, t.MarshalBinary, or t.Format with an explicit format string.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no input parameters for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
				{
					name: "Sub",
					comments: `Sub returns the duration t-u. If the result exceeds the maximum (or minimum) value that can be stored in a Duration, the maximum (or minimum) duration will be returned. To compute t-d for a duration d, use t.Add(-d).
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "u Time",
							name:     "u",
							typeName: "Time",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Duration",
							name:     "",
							typeName: "Duration",
							pointer:  false,
						},
					},
				},
				{
					name: "Truncate",
					comments: `Truncate returns the result of rounding t down to a multiple of d (since the zero time). If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.

Truncate operates on the time as an absolute duration since the zero time; it does not operate on the presentation form of the time. Thus, Truncate(Hour) may return a time with a non-zero minute, depending on the time's Location.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "UTC",
					comments: `UTC returns t with the location set to UTC.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "Time",
							name:     "",
							typeName: "Time",
							pointer:  false,
						},
					},
				},
				{
					name: "Unix",
					comments: `Unix returns t as a Unix time, the number of seconds elapsed since January 1, 1970 UTC. The result does not depend on the location associated with t. Unix-like operating systems often record time as a 32-bit count of seconds, but since the method here returns a 64-bit value it is valid for billions of years into the past or future.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int64",
							name:     "",
							typeName: "int64",
							pointer:  false,
						},
					},
				},
				{
					name: "UnixNano",
					comments: `UnixNano returns t as a Unix time, the number of nanoseconds elapsed since January 1, 1970 UTC. The result is undefined if the Unix time in nanoseconds cannot be represented by an int64 (a date before the year 1678 or after 2262). Note that this means the result of calling UnixNano on the zero Time is undefined. The result does not depend on the location associated with t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int64",
							name:     "",
							typeName: "int64",
							pointer:  false,
						},
					},
				},
				{
					name: "UnmarshalBinary",
					comments: `UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
`,
					receiver: testParameter{
						s:        "t *Time",
						name:     "t",
						typeName: "*Time",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "data []byte",
							name:     "data",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "UnmarshalJSON",
					comments: `UnmarshalJSON implements the json.Unmarshaler interface. The time is expected to be a quoted string in RFC 3339 format.
`,
					receiver: testParameter{
						s:        "t *Time",
						name:     "t",
						typeName: "*Time",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "data []byte",
							name:     "data",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "UnmarshalText",
					comments: `UnmarshalText implements the encoding.TextUnmarshaler interface. The time is expected to be in RFC 3339 format.
`,
					receiver: testParameter{
						s:        "t *Time",
						name:     "t",
						typeName: "*Time",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "data []byte",
							name:     "data",
							typeName: "[]byte",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "error",
							name:     "",
							typeName: "error",
							pointer:  false,
						},
					},
				},
				{
					name: "Weekday",
					comments: `Weekday returns the day of the week specified by t.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "Weekday",
							name:     "",
							typeName: "Weekday",
							pointer:  false,
						},
					},
				},
				{
					name: "Year",
					comments: `Year returns the year in which t occurs.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "YearDay",
					comments: `YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years, and [1,366] in leap years.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "int",
							name:     "",
							typeName: "int",
							pointer:  false,
						},
					},
				},
				{
					name: "Zone",
					comments: `Zone computes the time zone in effect at time t, returning the abbreviated name of the zone (such as "CET") and its offset in seconds east of UTC.
`,
					receiver: testParameter{
						s:        "t Time",
						name:     "t",
						typeName: "Time",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "name string",
							name:     "name",
							typeName: "string",
							pointer:  false,
						},
						{
							s:        "offset int",
							name:     "offset",
							typeName: "int",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Timer",
			typeName: "struct",
			source: `type Timer struct {
	C <-chan Time
	// contains filtered or unexported fields
}`,
			comments: `The Timer type represents a single event. When the Timer expires, the current time will be sent on C, unless the Timer was created by AfterFunc. A Timer must be created with NewTimer or AfterFunc.
`,
			functions: []testFunction{
				{
					name: "AfterFunc",
					comments: `AfterFunc waits for the duration to elapse and then calls f in its own goroutine. It returns a Timer that can be used to cancel the call using its Stop method.
`,
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
						{
							s:        "f func()",
							name:     "f",
							typeName: "func()",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Timer",
							name:     "",
							typeName: "*Timer",
							pointer:  true,
						},
					},
				},
				{
					name: "NewTimer",
					comments: `NewTimer creates a new Timer that will send the current time on its channel after at least duration d.
`,
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "*Timer",
							name:     "",
							typeName: "*Timer",
							pointer:  true,
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Reset",
					comments: `Reset changes the timer to expire after duration d. It returns true if the timer had been active, false if the timer had expired or been stopped.

For a Timer created with NewTimer, Reset should be invoked only on stopped or expired timers with drained channels.

If a program has already received a value from t.C, the timer is known to have expired and the channel drained, so t.Reset can be used directly. If a program has not yet received a value from t.C, however, the timer must be stopped and—if Stop reports that the timer expired before being stopped—the channel explicitly drained:

	if !t.Stop() {
		<-t.C
	}
	t.Reset(d)

This should not be done concurrent to other receives from the Timer's channel.

Note that it is not possible to use Reset's return value correctly, as there is a race condition between draining the channel and the new timer expiring. Reset should always be invoked on stopped or expired channels, as described above. The return value exists to preserve compatibility with existing programs.

For a Timer created with AfterFunc(d, f), Reset either reschedules when f will run, in which case Reset returns true, or schedules f to run again, in which case it returns false. When Reset returns false, Reset neither waits for the prior f to complete before returning nor does it guarantee that the subsequent goroutine running f does not run concurrently with the prior one. If the caller needs to know whether the prior execution of f is completed, it must coordinate with f explicitly.
`,
					receiver: testParameter{
						s:        "t *Timer",
						name:     "t",
						typeName: "*Timer",
						pointer:  true,
					},
					inputs: []testParameter{
						{
							s:        "d Duration",
							name:     "d",
							typeName: "Duration",
							pointer:  false,
						},
					},
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
				{
					name: "Stop",
					comments: `Stop prevents the Timer from firing. It returns true if the call stops the timer, false if the timer has already expired or been stopped. Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.

To ensure the channel is empty after a call to Stop, check the return value and drain the channel. For example, assuming the program has not received from t.C already:

	if !t.Stop() {
		<-t.C
	}

This cannot be done concurrent to other receives from the Timer's channel or other calls to the Timer's Stop method.

For a timer created with AfterFunc(d, f), if t.Stop returns false, then the timer has already expired and the function f has been started in its own goroutine; Stop does not wait for f to complete before returning. If the caller needs to know whether f is completed, it must coordinate with f explicitly.
`,
					receiver: testParameter{
						s:        "t *Timer",
						name:     "t",
						typeName: "*Timer",
						pointer:  true,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "bool",
							name:     "",
							typeName: "bool",
							pointer:  false,
						},
					},
				},
			},
		},
		{
			name:     "Weekday",
			typeName: "int",
			source:   `type Weekday int`,
			comments: `A Weekday specifies a day of the week (Sunday = 0, ...).
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name: "String",
					comments: `String returns the English name of the day ("Sunday", "Monday", ...).
`,
					receiver: testParameter{
						s:        "d Weekday",
						name:     "d",
						typeName: "Weekday",
						pointer:  false,
					},
					inputs: []testParameter{}, // no inputs for this method
					outputs: []testParameter{
						{
							s:        "string",
							name:     "",
							typeName: "string",
							pointer:  false,
						},
					},
				},
			},
		},
	},
}
