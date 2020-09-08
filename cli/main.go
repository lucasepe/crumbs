package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/crumbs"
	"github.com/lucasepe/crumbs/gv"
)

const (
	maxFileSize = 512 * 1000 // 512 Kb
	banner      = `
    __  ____  __ __  ___ ___  ____    _____
   /  ]|    \|  |  ||   |   ||    \  / ___/ v{{VERSION}}
  /  / |  D  )  |  || _   _ ||  o  )(   \_ 
 /  /  |    /|  |  ||  \_/  ||     | \__  |
/   \_ |    \|  :  ||   |   ||  O  | /  \ |
\     ||  .  \     ||   |   ||     | \    |
 \____||__|\_|\__,_||___|___||_____|  \___| 
Crafted with passion by Luca Sepe - https://github.com/lucasepe/crumbs`
)

var (
	version = "dev"
	commit  = "none"

	flagVertical  bool
	flagImagePath string
)

func main() {
	configureFlags()

	var entry *crumbs.Entry
	var err error
	if entry, err = readFromStdIn(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	if entry == nil {
		if entry, err = readFromFile(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if entry == nil {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	opts := []gv.GraphOption{gv.Vertical(flagVertical)}
	if len(flagImagePath) > 0 {
		opts = append(opts, gv.ImagesPath(flagImagePath))
	}

	if err = gv.Render(os.Stdout, entry, opts...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func readFromStdIn() (*crumbs.Entry, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		reader := io.LimitReader(bufio.NewReader(os.Stdin), maxFileSize)

		dat, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		return crumbs.FromString(string(dat))
	}

	return nil, nil
}

func readFromFile() (*crumbs.Entry, error) {
	if len(flag.Args()) == 0 {
		return nil, nil
	}

	reader, err := os.Open(flag.Args()[0])
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	dat, err := ioutil.ReadAll(io.LimitReader(reader, maxFileSize))
	if err != nil {
		return nil, err
	}

	return crumbs.FromString(string(dat))
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		printBanner()
		fmt.Printf("Turn asterisk-indented text lines into mind maps.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <path/to/your/file.txt>\n\n", name)

		fmt.Print("EXAMPLE(s):\n\n")
		fmt.Printf("  %s agenda.txt | dot -Tpng > output.png\n", name)
		fmt.Printf("  cat agenda.txt | %s | dot -Tpng > output.png\n\n", name)

		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.BoolVar(&flagVertical, "vertical", false,
		"layout entries as vertical directed graph")
	//flag.CommandLine.StringVar(&flagImagePath, "images-path", "",
	//	"list of directories (each separated by a semicolon for Windows or a colon for all other OS) in which to look for image files ")

	flag.CommandLine.Parse(os.Args[1:])
}

func printBanner() {
	str := strings.Replace(banner, "{{VERSION}}", version, 1)
	fmt.Print(str, "\n\n")
}

func appName() string {
	return filepath.Base(os.Args[0])
}

// exitOnErr check for an error and eventually exit
func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
