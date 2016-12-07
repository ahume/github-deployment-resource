package resource

import (
	"fmt"
	"os"
  "runtime/debug"

	"github.com/mitchellh/colorstring"
)

func Fatal(doing string, err error) {
	Sayf(colorstring.Color("[red]error %s: %s\n"), doing, err)
  debug.PrintStack()
	os.Exit(1)
}

func Sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}
