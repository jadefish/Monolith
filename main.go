package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"howett.net/plist"
)

func warn(a ...any) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}

type Variables struct {
	Debug        bool
	Product      string
	MLB          string
	ROM          []byte
	SerialNumber string
	UUID         string
}

func decodePlistFile(filename string) (map[string]any, error) {
	reader, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	data := map[string]any{}
	err = plist.NewDecoder(reader).Decode(data)

	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return data, nil
}

var debug = flag.Bool("debug", false, "Include debug configuration")
var product = flag.String("product", "", "Value for SystemProductName")
var mlb = flag.String("mlb", "", "Logic board serial number (MLB)")
var rom = flag.String("rom", "", "Base64-encoded 6-byte ROM value")
var serial = flag.String("serial", "", "System serial number")
var uuid = flag.String("uuid", "", "SMBIOS unique system identifier")

func main() {
	flag.Parse()

	if flag.NArg() < 2 || flag.NFlag() < 5 {
		_, _ = fmt.Fprintln(
			os.Stderr,
			"usage: monolith --product STRING --mlb STRING --rom BASE64_STRING --serial STRING --uuid UUID [--debug] BASE_PLIST_FILE INSTRUCTION_FILE",
		)
		os.Exit(1)
	}

	plistFile := flag.Arg(0)
	instructionFile := flag.Arg(1)

	data, err := decodePlistFile(plistFile)

	if err != nil {
		warn(err)
		os.Exit(1)
	}

	f, err := os.Open(instructionFile)

	if err != nil {
		warn(err)
		os.Exit(1)
	}

	// package plist encodes []byte fields as Base64 when encoding a plist, so
	// we must hold the unencoded form of it in memory.
	plainROM, err := base64.StdEncoding.DecodeString(*rom)

	if err != nil {
		warn("rom: %w", err)
		os.Exit(1)
	}

	vars := Variables{
		Debug:        *debug,
		Product:      *product,
		MLB:          *mlb,
		ROM:          plainROM,
		SerialNumber: *serial,
		UUID:         *uuid,
	}

	// read in instructions:
	bs := bufio.NewScanner(f)
	var instructions []string

	for bs.Scan() {
		if err := bs.Err(); err != nil {
			warn(fmt.Errorf("read error: %w", err))
			os.Exit(1)
		}

		line := strings.TrimSpace(bs.Text())

		// expr supports comments, but can't parse expressions that are *only*
		// comments.
		if len(line) < 1 || strings.HasPrefix(line, "//") {
			continue
		}

		instructions = append(instructions, line)
	}

	// ... compile them...
	evaluator := Evaluator{data: data}
	env := map[string]any{
		"delete":  evaluator.delete,
		"set":     evaluator.set,
		"append":  evaluator.append,
		"helpers": Helpers{},
		"vars":    vars,
	}
	var programs []*vm.Program

	for _, instruction := range instructions {
		program, err := expr.Compile(instruction, expr.Env(env))

		if err != nil {
			warn(fmt.Errorf("compile error: %w", err))
			os.Exit(1)
		}

		programs = append(programs, program)
	}

	// ... then run 'em:
	exprVM := vm.VM{}

	for _, program := range programs {
		warn(">", program.Source.Content())
		if _, err := exprVM.Run(program, env); err != nil {
			warn(fmt.Errorf("runtime error: %w", err))
			os.Exit(1)
		}
	}

	bw := bufio.NewWriter(os.Stdout)
	encoder := plist.NewEncoder(bw)
	encoder.Indent("    ") // 4 spaces for parity with plutil
	err = encoder.Encode(data)

	if err != nil {
		warn(fmt.Errorf("unable to encode plist: %w", err))
		os.Exit(1)
	}
}
