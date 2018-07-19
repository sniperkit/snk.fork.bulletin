// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//NO TESTS

package printer

import (
	"fmt"
	"io"
	"runtime/debug"

	"os"

	"github.com/fatih/color"
)

// Printer for customized printing
type Printer struct {
	errout io.Writer
	debug  bool
}

// New return Printer Obj
func New(stderr io.Writer) *Printer {
	color.Output = stderr
	return &Printer{
		errout: stderr,
		debug:  false,
	}
}

// Println Print a message to Stderr
func (p *Printer) Println(args ...interface{}) {
	// #nosec Errors unhandled
	fmt.Fprintln(p.errout, args...)
}

// Printf Print a message to Stderr
func (p *Printer) Printf(format string, args ...interface{}) {
	// #nosec Errors unhandled
	fmt.Fprintf(p.errout, format, args...)
}

// Warn Print a warning to Stderr
func (p *Printer) Warn(format string, args ...interface{}) {
	// #nosec Errors unhandled
	color.Yellow(format, args...)
}

// Error Print an error message to Stderr
func (p *Printer) Error(format string, args ...interface{}) {
	color.Red(format, args...)
	if p.debug {
		debug.PrintStack()
	}
}

// Fatal Print a error message and and stack trace to Stderr and exit.
// Prints stack if debug flag is passed
func (p *Printer) Fatal(format string, args ...interface{}) {
	p.Error(format, args...)
	os.Exit(1)
}

// VerboseInfo print a VerboseInfo message
func (p *Printer) VerboseInfo(format string, args ...interface{}) {
	// #nosec Errors unhandled
	color.New(color.Faint).Printf(format, args...)
}

// SetOutput set output for printer
func (p *Printer) SetOutput(stderr io.Writer) {
	p.errout = stderr
	color.Output = stderr
}

// GetOutput get stdout for printer
func (p *Printer) GetOutput() io.Writer {
	return p.errout
}

// SetDebug set debug bool
func (p *Printer) SetDebug(debug bool) {
	p.debug = debug
}
