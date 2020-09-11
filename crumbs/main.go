package main

import (
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
	maxFileSize int64 = 512 * 1000 // 512 Kb
	banner            = `
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
	version = "0.3.0"

	flagVertical   bool
	flagWrapLim    uint
	flagImagesPath string
	flagImagesType string
)

func main() {
	configureFlags()

	entry, err := readEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	cfg := gv.RenderConfig{
		WrapTextLimit:  flagWrapLim,
		VerticalLayout: flagVertical,
	}
	if err = gv.Render(os.Stdout, entry, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func readInput() ([]byte, error) {
	limit := maxFileSize
	args := flag.Args()
	if len(args) == 0 {
		return readFileObject(os.Stdin, limit)
	}
	return readFile(args[0], limit)
}

func readEntry() (*crumbs.Entry, error) {
	src, err := readInput()
	if err != nil {
		return nil, err
	}
	text := string(src)
	lines := strings.SplitAfter(text, "\n")
	return crumbs.ParseLines(lines, flagImagesPath, flagImagesType)
}

func readFileObject(r io.Reader, limit int64) ([]byte, error) {
	lr := io.LimitReader(r, limit)
	return ioutil.ReadAll(lr)
}

func readFile(name string, limit int64) ([]byte, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return readFileObject(r, limit)
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
	flag.CommandLine.UintVar(&flagWrapLim, "lim", 28, "wraps each line within this width in characters")

	flag.CommandLine.StringVar(&flagImagesPath, "images-path", "", "folder in which to look for image files")
	flag.CommandLine.StringVar(&flagImagesType, "images-type", "", "images file extension [png,jpg,svg]")

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
