package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xzavier-amico/module-tool/internal/config"
	"github.com/xzavier-amico/module-tool/internal/converter"
	"github.com/xzavier-amico/module-tool/internal/format"
	"github.com/xzavier-amico/module-tool/internal/generator"
	"github.com/xzavier-amico/module-tool/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "convert":
		fs := flag.NewFlagSet("convert", flag.ExitOnError)
		from := fs.String("from", "", "source format")
		to := fs.String("to", "", "target format")
		in := fs.String("input", "", "input file (or - for stdin)")
		out := fs.String("output", "", "output file (or - for stdout)")
		fs.Parse(os.Args[2:])
		if !format.Valid(*from) || !format.Valid(*to) {
			fmt.Println("invalid format")
			os.Exit(1)
		}
		var content string
		if *in == "" || *in == "-" {
			b, _ := os.ReadFile("/dev/stdin")
			content = string(b)
		} else {
			b, err := os.ReadFile(*in)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			content = string(b)
		}
		res := converter.Convert(*from, *to, content)
		if *out == "" || *out == "-" {
			fmt.Print(res)
		} else {
			os.WriteFile(*out, []byte(res), 0644)
			fmt.Println("written", *out)
		}
	case "generate":
		fs := flag.NewFlagSet("generate", flag.ExitOnError)
		cfgPath := fs.String("config", "config.yaml", "config file")
		base := fs.String("base", ".", "base directory")
		fs.Parse(os.Args[2:])
		cfg, err := config.Load(*cfgPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := generator.Generate(cfg, *base); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "serve":
		fs := flag.NewFlagSet("serve", flag.ExitOnError)
		port := fs.String("port", "8080", "port")
		fs.Parse(os.Args[2:])
		server.Start(*port)
	case "formats":
		fmt.Println(strings.Join(format.All, ", "))
	default:
		usage()
	}
}

func usage() {
	fmt.Println("module-tool: convert network modules between formats")
	fmt.Println("commands: convert, generate, serve, formats")
}
