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
	imports        []string
	testImports    []string
	constantBlocks []testConstantBlock
	variables      []string
	errors         []string
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
	name   string
	source string // TODO
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
	imports:        []string{"internal/reflectlite"},
	testImports:    []string{"errors", "fmt", "io/fs", "os", "reflect", "testing", "time"},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variables:      []string{},            // TODO
	errors:         []string{},            // TODO
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
			name: "New",
			comments: `New returns an error that formats as the given text. Each call to New returns a distinct error value even if the text is identical.
`,
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
			name: "Unwrap",
			comments: `Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
`,
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
	imports: []string{
		"errors", "internal/fmtsort", "io", "math", "os", "reflect", "strconv", "sync", "unicode/utf8",
	},
	testImports: []string{
		"bufio", "bytes", "errors", "fmt", "internal/race", "io", "math", "os", "reflect", "regexp",
		"runtime", "strings", "testing", "testing/iotest", "time", "unicode", "unicode/utf8",
	},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variables:      []string{},            // TODO
	errors:         []string{},            // TODO
	functions: []testFunction{
		{
			name: "Errorf",
			comments: `Errorf formats according to a format specifier and returns the string as a value that satisfies error.

If the format specifier includes a %w verb with an error operand, the returned error will implement an Unwrap method returning the operand. It is invalid to include more than one %w verb or to supply it with an operand that does not implement the error interface. The %w verb is otherwise a synonym for %v.
`,
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
			name: "Fprint",
			comments: `Fprint formats using the default formats for its operands and writes to w. Spaces are added between operands when neither is a string. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Fprintf",
			comments: `Fprintf formats according to a format specifier and writes to w. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Fprintln",
			comments: `Fprintln formats using the default formats for its operands and writes to w. Spaces are always added between operands and a newline is appended. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Fscan",
			comments: `Fscan scans text read from r, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
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
			name: "Fscanf",
			comments: `Fscanf scans text read from r, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully parsed. Newlines in the input must match newlines in the format.
`,
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
			name: "Fscanln",
			comments: `Fscanln is similar to Fscan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
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
			name: "Print",
			comments: `Print formats using the default formats for its operands and writes to standard output. Spaces are added between operands when neither is a string. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Printf",
			comments: `Printf formats according to a format specifier and writes to standard output. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Println",
			comments: `Println formats using the default formats for its operands and writes to standard output. Spaces are always added between operands and a newline is appended. It returns the number of bytes written and any write error encountered.
`,
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
			name: "Scan",
			comments: `Scan scans text read from standard input, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
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
			name: "Scanf",
			comments: `Scanf scans text read from standard input, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why. Newlines in the input must match newlines in the format. The one exception: the verb %c always scans the next rune in the input, even if it is a space (or tab etc.) or newline.
`,
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
			name: "Scanln",
			comments: `Scanln is similar to Scan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
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
			name: "Sprint",
			comments: `Sprint formats using the default formats for its operands and returns the resulting string. Spaces are added between operands when neither is a string.
`,
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
			name: "Sprintf",
			comments: `Sprintf formats according to a format specifier and returns the resulting string.
`,
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
			name: "Sprintln",
			comments: `Sprintln formats using the default formats for its operands and returns the resulting string. Spaces are always added between operands and a newline is appended.
`,
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
			name: "Sscan",
			comments: `Sscan scans the argument string, storing successive space-separated values into successive arguments. Newlines count as space. It returns the number of items successfully scanned. If that is less than the number of arguments, err will report why.
`,
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
			name: "Sscanf",
			comments: `Sscanf scans the argument string, storing successive space-separated values into successive arguments as determined by the format. It returns the number of items successfully parsed. Newlines in the input must match newlines in the format.
`,
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
			name: "Sscanln",
			comments: `Sscanln is similar to Sscan, but stops scanning at a newline and after the final item there must be a newline or EOF.
`,
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
			name:     "Formatter",
			typeName: "interface",
			source:   "",
			comments: `Formatter is implemented by any value that has a Format method. The implementation controls how State and rune are interpreted, and may call Sprint(f) or Fprint(f) etc. to generate its output.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "GoStringer",
			typeName: "interface",
			source:   "",
			comments: `GoStringer is implemented by any value that has a GoString method, which defines the Go syntax for that value. The GoString method is used to print values passed as an operand to a %#v format.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "ScanState",
			typeName: "interface",
			source:   "",
			comments: `ScanState represents the scanner state passed to custom scanners. Scanners may do rune-at-a-time scanning or ask the ScanState to discover the next space-delimited token.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Scanner",
			typeName: "interface",
			source:   "",
			comments: `Scanner is implemented by any value that has a Scan method, which scans the input for the representation of a value and stores the result in the receiver, which must be a pointer to be useful. The Scan method is called for any argument to Scan, Scanf, or Scanln that implements it.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "State",
			typeName: "interface",
			source:   "",
			comments: `State represents the printer state passed to custom formatters. It provides access to the io.Writer interface plus information about the flags and options for the operand's format specifier.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Stringer",
			typeName: "interface",
			source:   "",
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
	imports: []string{
		"io",
	},
	testImports: []string{
		"bytes", "crypto/md5", "crypto/sha1", "crypto/sha256", "crypto/sha512", "encoding", "encoding/hex",
		"fmt", "hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv", "log", "testing",
	},
	constantBlocks: []testConstantBlock{}, // no exported constants in this package
	variables:      []string{},            // TODO
	errors:         []string{},            // TODO
	functions:      []testFunction{},      // no functions in this package
	types: []testType{
		{
			name:     "Hash",
			typeName: "interface",
			source:   "",
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
			source:   "",
			comments: `Hash32 is the common interface implemented by all 32-bit hash functions.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Hash64",
			typeName: "interface",
			source:   "",
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
					name:   "TypeReg",
					source: `TypeReg  = '0'`,
				},
				{
					name:   "TypeRegA",
					source: `TypeRegA = '\x00' // Deprecated: Use TypeReg instead.`,
				},
				{
					name:   "TypeLink",
					source: `TypeLink    = '1' // Hard link`,
				},
				{
					name:   "TypeSymlink",
					source: `TypeSymlink = '2' // Symbolic link`,
				},
				{
					name:   "TypeChar",
					source: `TypeChar    = '3' // Character device node`,
				},
				{
					name:   "TypeBlock",
					source: `TypeBlock   = '4' // Block device node`,
				},
				{
					name:   "TypeDir",
					source: `TypeDir     = '5' // Directory`,
				},
				{
					name:   "TypeFifo",
					source: `TypeFifo    = '6' // FIFO node`,
				},
				{
					name:   "TypeCont",
					source: `TypeCont = '7'`,
				},
				{
					name:   "TypeXHeader",
					source: `TypeXHeader = 'x'`,
				},
				{
					name:   "TypeXGlobalHeader",
					source: `TypeXGlobalHeader = 'g'`,
				},
				{
					name:   "TypeGNUSparse",
					source: `TypeGNUSparse = 'S'`,
				},
				{
					name:   "TypeGNULongName",
					source: `TypeGNULongName = 'L'`,
				},
				{
					name:   "TypeGNULongLink",
					source: `TypeGNULongLink = 'K'`,
				},
			},
		},
		{
			typeName: "Format",
			comments:
`Constants to identify various tar formats.
`,
			source: `const (
	// Deliberately hide the meaning of constants from public API.
	_ Format = (1 << iota) / 4 // Sequence of 0, 0, 1, 2, 4, 8, etc...

	// FormatUnknown indicates that the format is unknown.
	FormatUnknown

	// The format of the original Unix V7 tar tool prior to standardization.
	formatV7

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

	// Schily's tar format, which is incompatible with USTAR.
	// This does not cover STAR extensions to the PAX format; these fall under
	// the PAX format.
	formatSTAR

	formatMax
)`,
			constants: []testConstant{
				{
					name:   "FormatUnknown",
					source: `FormatUnknown`,
				},
				{
					name:   "FormatUSTAR",
					source: `FormatUSTAR`,
				},
				{
					name:   "FormatPAX",
					source: `FormatPAX`,
				},
				{
					name:   "FormatGNU",
					source: `FormatGNU`,
				},
			},
		},
	},
	variables: []string{},       // TODO
	errors:    []string{},       // TODO
	functions: []testFunction{}, // no functions in this package
	types: []testType{
		{
			name:     "Format",
			typeName: "int",
			source:   "",
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
					name:        "String",
					comments:    ``, // no comments for this method
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
					name: "FileInfo",
					comments: `FileInfo returns an fs.FileInfo for the Header.
`,
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
			comments: `Reader provides sequential access to the contents of a tar archive. Reader.Next advances to the next file in the archive (including the first), and then Reader can be treated as an io.Reader to access the file's data.
`,
			functions: []testFunction{
				{
					name: "NewReader",
					comments: `NewReader creates a new Reader reading from r.
`,
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
					name: "Next",
					comments: `Next advances to the next entry in the tar archive. The Header.Size determines how many bytes can be read for the next file. Any remaining data in the current file is automatically discarded.

io.EOF is returned at the end of the input.
`,
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
					name: "Read",
					comments: `Read reads from the current file in the tar archive. It returns (0, io.EOF) when it reaches the end of that file, until Next is called to advance to the next file.

If the current file is sparse, then the regions marked as a hole are read back as NUL-bytes.

Calling Read on special types like TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, and TypeFifo returns (0, io.EOF) regardless of what the Header.Size claims.
`,
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
			comments: `Writer provides sequential writing of a tar archive. Write.WriteHeader begins a new file with the provided Header, and then Writer can be treated as an io.Writer to supply that file's data.
`,
			functions: []testFunction{
				{
					name: "NewWriter",
					comments: `NewWriter creates a new Writer writing to w.
`,
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
					name: "Close",
					comments: `Close closes the tar archive by flushing the padding, and writing the footer. If the current file (from a prior call to WriteHeader) is not fully written, then this returns an error.
`,
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
					name: "Flush",
					comments: `Flush finishes writing the current file's block padding. The current file must be fully written before Flush can be called.

This is unnecessary as the next call to WriteHeader or Close will implicitly flush out the file's padding.
`,
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
					name: "Write",
					comments: `Write writes to the current file in the tar archive. Write returns the error ErrWriteTooLong if more than Header.Size bytes are written after WriteHeader.

Calling Write on special types like TypeLink, TypeSymlink, TypeChar, TypeBlock, TypeDir, and TypeFifo returns (0, ErrWriteTooLong) regardless of what the Header.Size claims.
`,
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
					name: "WriteHeader",
					comments: `WriteHeader writes hdr and prepares to accept the file's contents. The Header.Size determines how many bytes can be written for the next file. If the current file is not fully written, then this returns an error. This implicitly flushes any padding necessary before writing the header.
`,
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
	imports: []string{}, // no imports for this package
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
					name:   "MaxRune",
					source: `MaxRune         = '\U0010FFFF' // Maximum valid Unicode code point.`,
				},
				{
					name:   "ReplacementChar",
					source: `ReplacementChar = '\uFFFD'     // Represents invalid code points.`,
				},
				{
					name:   "MaxASCII",
					source: `MaxASCII        = '\u007F'     // maximum ASCII value.`,
				},
				{
					name:   "MaxLatin1",
					source: `MaxLatin1       = '\u00FF'     // maximum Latin-1 value.`,
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
					name:   "UpperCase",
					source: `UpperCase = iota`,
				},
				{
					name:   "LowerCase",
					source: `LowerCase`,
				},
				{
					name:   "TitleCase",
					source: `TitleCase`,
				},
				{
					name:   "MaxCase",
					source: `MaxCase`,
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
					name:   "UpperLower",
					source: `UpperLower = MaxRune + 1 // (Cannot be a valid delta.)`,
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
					name:   "Version",
					source: `Version = "13.0.0"`,
				},
			},
		},
	},
	variables:      []string{},            // TODO
	errors:         []string{},            // TODO
	functions: []testFunction{
		{
			name: "In",
			comments: `In reports whether the rune is a member of one of the ranges.
`,
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
			name: "Is",
			comments: `Is reports whether the rune is in the specified table of ranges.
`,
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
			name: "IsControl",
			comments: `IsControl reports whether the rune is a control character. The C (Other) Unicode category includes more code points such as surrogates; use Is(C, r) to test for them.
`,
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
			name: "IsDigit",
			comments: `IsDigit reports whether the rune is a decimal digit.
`,
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
			name: "IsGraphic",
			comments: `IsGraphic reports whether the rune is defined as a Graphic by Unicode. Such characters include letters, marks, numbers, punctuation, symbols, and spaces, from categories L, M, N, P, S, Zs.
`,
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
			name: "IsLetter",
			comments: `IsLetter reports whether the rune is a letter (category L).
`,
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
			name: "IsLower",
			comments: `IsLower reports whether the rune is a lower case letter.
`,
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
			name: "IsMark",
			comments: `IsMark reports whether the rune is a mark character (category M).
`,
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
			name: "IsNumber",
			comments: `IsNumber reports whether the rune is a number (category N).
`,
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
			name: "IsOneOf",
			comments: `IsOneOf reports whether the rune is a member of one of the ranges. The function "In" provides a nicer signature and should be used in preference to IsOneOf.
`,
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
			name: "IsPrint",
			comments: `IsPrint reports whether the rune is defined as printable by Go. Such characters include letters, marks, numbers, punctuation, symbols, and the ASCII space character, from categories L, M, N, P, S and the ASCII space character. This categorization is the same as IsGraphic except that the only spacing character is ASCII space, U+0020.
`,
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
			name: "IsPunct",
			comments: `IsPunct reports whether the rune is a Unicode punctuation character (category P).
`,
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
			name: "IsSpace",
			comments: `IsSpace reports whether the rune is a space character as defined by Unicode's White Space property; in the Latin-1 space this is

	'\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).

Other definitions of spacing characters are set by category Z and property Pattern_White_Space.
`,
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
			name: "IsSymbol",
			comments: `IsSymbol reports whether the rune is a symbolic character.
`,
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
			name: "IsTitle",
			comments: `IsTitle reports whether the rune is a title case letter.
`,
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
			name: "IsUpper",
			comments: `IsUpper reports whether the rune is an upper case letter.
`,
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
			name: "To",
			comments: `To maps the rune to the specified case: UpperCase, LowerCase, or TitleCase.
`,
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
			name: "ToLower",
			comments: `ToLower maps the rune to lower case.
`,
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
			name: "ToTitle",
			comments: `ToTitle maps the rune to title case.
`,
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
			name: "ToUpper",
			comments: `ToUpper maps the rune to upper case.
`,
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
			name:     "CaseRange",
			typeName: "struct",
			source:   "",
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
			source:   "",
			comments: `Range16 represents of a range of 16-bit Unicode code points. The range runs from Lo to Hi inclusive and has the specified stride.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Range32",
			typeName: "struct",
			source:   "",
			comments: `Range32 represents of a range of Unicode code points and is used when one or more of the values will not fit in 16 bits. The range runs from Lo to Hi inclusive and has the specified stride. Lo and Hi must always be >= 1<<16.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "RangeTable",
			typeName: "struct",
			source:   "",
			comments: `RangeTable defines a set of Unicode code points by listing the ranges of code points within the set. The ranges are listed in two slices to save space: a slice of 16-bit ranges and a slice of 32-bit ranges. The two slices must be in sorted order and non-overlapping. Also, R32 should contain only values >= 0x10000 (1<<16).
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "SpecialCase",
			typeName: "[]CaseRange",
			source:   "",
			comments: `SpecialCase represents language-specific case mappings such as Turkish. Methods of SpecialCase customize (by overriding) the standard mappings.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name: "ToLower",
					comments: `ToLower maps the rune to lower case giving priority to the special mapping.
`,
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
					name: "ToTitle",
					comments: `ToTitle maps the rune to title case giving priority to the special mapping.
`,
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
					name: "ToUpper",
					comments: `ToUpper maps the rune to upper case giving priority to the special mapping.
`,
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
					name:   "DefaultRPCPath",
					source: `DefaultRPCPath   = "/_goRPC_"`,
				},
				{
					name:   "DefaultDebugPath",
					source: `DefaultDebugPath = "/debug/rpc"`,
				},
			},
		},
	},
	variables:      []string{},            // TODO
	errors:         []string{},            // TODO
	functions: []testFunction{
		{
			name: "Accept",
			comments: `Accept accepts connections on the listener and serves requests to DefaultServer for each incoming connection. Accept blocks; the caller typically invokes it in a go statement.
`,
			inputs: []testParameter{
				{
					name:     "lis",
					typeName: "net.Listener",
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
			name: "RegisterName",
			comments: `RegisterName is like Register but uses the provided name for the type instead of the receiver's concrete type.
`,
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
			name: "ServeCodec",
			comments: `ServeCodec is like ServeConn but uses the specified codec to decode requests and encode responses.
`,
			inputs: []testParameter{
				{
					name:     "codec",
					typeName: "ServerCodec",
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
					name:     "conn",
					typeName: "io.ReadWriteCloser",
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
			name:     "Call",
			typeName: "struct",
			source:   "",
			comments: `Call represents an active RPC.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Client",
			typeName: "struct",
			source:   "",
			comments: `Client represents an RPC Client. There may be multiple outstanding Calls associated with a single Client, and a Client may be used by multiple goroutines simultaneously.
`,
			functions: []testFunction{
				{
					name: "Dial",
					comments: `Dial connects to an RPC server at the specified network address.
`,
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
					name: "DialHTTP",
					comments: `DialHTTP connects to an HTTP RPC server at the specified network address listening on the default HTTP RPC path.
`,
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
					name: "DialHTTPPath",
					comments: `DialHTTPPath connects to an HTTP RPC server at the specified network address and path.
`,
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
					name: "NewClient",
					comments: `NewClient returns a new Client to handle requests to the set of services at the other end of the connection. It adds a buffer to the write side of the connection so the header and payload are sent as a unit.

The read and write halves of the connection are serialized independently, so no interlocking is required. However each half may be accessed concurrently so the implementation of conn should protect against concurrent reads or concurrent writes.
`,
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
					name: "NewClientWithCodec",
					comments: `NewClientWithCodec is like NewClient but uses the specified codec to encode requests and decode responses.
`,
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
					name: "Call",
					comments: `Call invokes the named function, waits for it to complete, and returns its error status.
`,
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
					name: "Close",
					comments: `Close calls the underlying codec's Close method. If the connection is already shutting down, ErrShutdown is returned.
`,
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
					name: "Go",
					comments: `Go invokes the function asynchronously. It returns the Call structure representing the invocation. The done channel will signal when the call is complete by returning the same Call object. If done is nil, Go will allocate a new channel. If non-nil, done must be buffered or Go will deliberately crash.
`,
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
			name:     "ClientCodec",
			typeName: "interface",
			source:   "",
			comments: `A ClientCodec implements writing of RPC requests and reading of RPC responses for the client side of an RPC session. The client calls WriteRequest to write a request to the connection and calls ReadResponseHeader and ReadResponseBody in pairs to read responses. The client calls Close when finished with the connection. ReadResponseBody may be called with a nil argument to force the body of the response to be read and then discarded. See NewClient's comment for information about concurrent access.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Request",
			typeName: "struct",
			source:   "",
			comments: `Request is a header written before every RPC call. It is used internally but documented here as an aid to debugging, such as when analyzing network traffic.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Response",
			typeName: "struct",
			source:   "",
			comments: `Response is a header written before every RPC return. It is used internally but documented here as an aid to debugging, such as when analyzing network traffic.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "Server",
			typeName: "struct",
			source:   "",
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
							name:     "",
							typeName: "*Server",
						},
					},
				},
			},
			methods: []testMethod{
				{
					name: "Accept",
					comments: `Accept accepts connections on the listener and serves requests for each incoming connection. Accept blocks until the listener returns a non-nil error. The caller typically invokes Accept in a go statement.
`,
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
					name: "HandleHTTP",
					comments: `HandleHTTP registers an HTTP handler for RPC messages on rpcPath, and a debugging handler on debugPath. It is still necessary to invoke http.Serve(), typically in a go statement.
`,
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
					name: "Register",
					comments: `Register publishes in the server the set of methods of the receiver value that satisfy the following conditions:

	- exported method of exported type
	- two arguments, both of exported type
	- the second argument is a pointer
	- one return value, of type error

It returns an error if the receiver is not an exported type or has no suitable methods. It also logs the error using package log. The client accesses each method using a string of the form "Type.Method", where Type is the receiver's concrete type.
`,
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
					name: "RegisterName",
					comments: `RegisterName is like Register but uses the provided name for the type instead of the receiver's concrete type.
`,
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
					name: "ServeCodec",
					comments: `ServeCodec is like ServeConn but uses the specified codec to decode requests and encode responses.
`,
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
					name: "ServeConn",
					comments: `ServeConn runs the server on a single connection. ServeConn blocks, serving the connection until the client hangs up. The caller typically invokes ServeConn in a go statement. ServeConn uses the gob wire format (see package gob) on the connection. To use an alternate codec, use ServeCodec. See NewClient's comment for information about concurrent access.
`,
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
					name: "ServeHTTP",
					comments: `ServeHTTP implements an http.Handler that answers RPC requests.
`,
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
					name: "ServeRequest",
					comments: `ServeRequest is like ServeCodec but synchronously serves a single request. It does not close the codec upon completion.
`,
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
			name:     "ServerCodec",
			typeName: "interface",
			source:   "",
			comments: `A ServerCodec implements reading of RPC requests and writing of RPC responses for the server side of an RPC session. The server calls ReadRequestHeader and ReadRequestBody in pairs to read requests from the connection, and it calls WriteResponse to write a response back. The server calls Close when finished with the connection. ReadRequestBody may be called with a nil argument to force the body of the request to be read and discarded. See NewClient's comment for information about concurrent access.
`,
			functions: []testFunction{}, // no functions for this type
			methods:   []testMethod{},   // no methods for this type
		},
		{
			name:     "ServerError",
			typeName: "string",
			source:   "",
			comments: `ServerError represents an error that has been returned from the remote side of the RPC connection.
`,
			functions: []testFunction{}, // no functions for this type
			methods: []testMethod{
				{
					name:        "Error",
					comments:    ``, // no comments for this method
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
}
