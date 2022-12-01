package helpers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Print logs values to standard error.
func Print(v ...interface{}) {
	logger := log.New(os.Stderr, debugPrefix(), 0)
	logger.Print(v...)
}

// Printf logs values to standard error with the provided format.
func Printf(format string, v ...interface{}) {
	logger := log.New(os.Stderr, debugPrefix(), 0)
	logger.Printf(format, v...)
}

// Println logs values to standard error, then starts a new line.
func Println(v ...interface{}) {
	logger := log.New(os.Stderr, debugPrefix(), 0)
	logger.Println(v...)
}

func debugPrefix() string {
	programCounter, fileAbsolutePath, line, _ := runtime.Caller(2)

	packageName := packageFromFunc(runtime.FuncForPC(programCounter))
	_, fileName := filepath.Split(fileAbsolutePath)
	fullFileName := filepath.Join(packageName, fileName)

	return fmt.Sprintf("[%s:%d] ", fullFileName, line)
}

func packageFromFunc(f *runtime.Func) string {
	splitName := strings.Split(f.Name(), ".")
	packageName := strings.Join(splitName[:len(splitName)-1], ".")
	packageName = strings.TrimPrefix(packageName, "github.com/busser/adventofcode/")
	return packageName
}
