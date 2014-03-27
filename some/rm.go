package some

import (
	"fmt"
	"github.com/laher/someutils"
	"github.com/laher/uggo"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	someutils.RegisterPipable(func() someutils.PipableCliUtil { return NewRm() })
}

// SomeRm represents and performs a `rm` invocation
type SomeRm struct {
	IsRecursive bool
	fileGlobs   []string
}

// Name() returns the name of the util
func (rm *SomeRm) Name() string {
	return "rm"
}

// ParseFlags parses flags from a commandline []string
func (rm *SomeRm) ParseFlags(call []string, errPipe io.Writer) error {
	flagSet := uggo.NewFlagSetDefault("rm", "[options] [files...]", someutils.VERSION)
	flagSet.SetOutput(errPipe)

	flagSet.BoolVar(&rm.IsRecursive, "r", false, "Recurse into directories")

	err := flagSet.Parse(call[1:])
	if err != nil {
		fmt.Fprintf(errPipe, "Flag error:  %v\n\n", err.Error())
		flagSet.Usage()
		return err
	}

	if flagSet.ProcessHelpOrVersion() {
		return nil
	}

	rm.fileGlobs = flagSet.Args()
	return nil
}

// Exec actually performs the rm
func (rm *SomeRm) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	for _, fileGlob := range rm.fileGlobs {
		files, err := filepath.Glob(fileGlob)
		if err != nil {
			return err
		}
		for _, file := range files {
			e := delete(file, rm.IsRecursive)
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func delete(file string, recursive bool) error {
	fi, e := os.Stat(file)
	if e != nil {
		return e
	}
	if fi.IsDir() && recursive {
		e := deleteDir(file)
		if e != nil {
			return e
		}
	} else if fi.IsDir() {
		//do nothing
		return fmt.Errorf("'%s' is a directory. Use -r", file)
	}
	return os.Remove(file)
}

func deleteDir(dir string) error {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return e
	}
	for _, file := range files {
		e = delete(filepath.Join(dir, file.Name()), true)
		if e != nil {
			return e
		}
	}
	return nil
}

// Factory for *SomeRm
func NewRm() *SomeRm {
	return new(SomeRm)
}

// Fluent factory for *SomeRm
func Rm(args ...string) *SomeRm {
	rm := NewRm()
	rm.fileGlobs = args
	return rm
}

// CLI invocation for *SomeRm
func RmCli(call []string) error {
	rm := NewRm()
	inPipe, outPipe, errPipe := someutils.StdPipes()
	err := rm.ParseFlags(call, errPipe)
	if err != nil {
		return err
	}
	return rm.Exec(inPipe, outPipe, errPipe)
}