package internal

import (
	"flag"
	"os"
)

// parseFlags if a variable function that parses the command line arguments.
// This is a variable so that it can be replaced in tests to avoid side effects.
var ParseFlags = func(f *flag.FlagSet, args []string) error {
	return f.Parse(args)
}

// Options is a struct that implements a parser for the program options.
type Options struct{}

// parse is a method that parses the command line arguments and returns the
// appropriate command to run (if any).
func (o *Options) Parse() (interface{ Run(*Options) int }, error) {
	var (
		of string     = "test-report.md"
		rm reportMode = rmFailedTests
		rt string     = "Test Report"

		opts = struct {
			h, help    bool
			f, full    bool
			o, output  string
			s, summary bool
			t, title   string
		}{}
	)

	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			return showVersion{}, nil
		}

		flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flags.BoolVar(&opts.f, "f", false, "complete test report")
		flags.BoolVar(&opts.full, "full", false, "")
		flags.BoolVar(&opts.h, "h", false, "display help on usage")
		flags.BoolVar(&opts.help, "help", false, "")
		flags.StringVar(&opts.o, "o", "", "output filename")
		flags.StringVar(&opts.output, "output", "", "")
		flags.BoolVar(&opts.s, "s", false, "summary only")
		flags.BoolVar(&opts.summary, "summary", false, "")
		flags.StringVar(&opts.t, "t", "", "report title")
		flags.StringVar(&opts.title, "title", "", "")
		if err := ParseFlags(flags, os.Args[1:]); err != nil {
			return nil, err
		}

		switch {
		case opts.o != "":
			of = opts.o
		case opts.output != "":
			of = opts.output
		}

		switch {
		case opts.t != "":
			rt = opts.t
		case opts.title != "":
			rt = opts.title
		}

		switch {
		case opts.f || opts.full:
			rm = rmAllTests
		case opts.s || opts.summary:
			rm = rmSummaryOnly
		}
	}

	switch {
	case opts.h || opts.help:
		return showUsage{}, nil

	default:
		return generateReport{
			filename: of,
			title:    rt,
			mode:     rm,
			parser:   &parser{},
		}, nil
	}
}
