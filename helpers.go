package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/exp/slices"
	"howett.net/plist"
)

var fullNvramAccessTools = []string{
	"CleanNvram.efi",
	"ControlMsrE2.efi",
	"OpenControl.efi",
}

var loadEarlyDrivers = []string{
	"OpenVariableRuntimeDxe.efi",
}

type Helpers struct{}

// Kext returns a Kernel.Add struct for the provided kext file.
func (Helpers) Kext(filename string) (Kext, error) {
	f, err := os.Open(filename)

	if err != nil {
		return Kext{}, fmt.Errorf("kext: %w", err)
	}

	// Determine values for "ExecutablePath" and "MinKernel":
	executablePath := ""
	minKernel := "8.0.0"
	plistPath := path.Join("Contents", "Info.plist")
	plistFile, err := os.Open(path.Join(filename, plistPath))

	if err == nil {
		var kextPlist KextPropertyList
		if err := plist.NewDecoder(plistFile).Decode(&kextPlist); err == nil {
			executablePath = path.Join(
				"Contents",
				"MacOS",
				kextPlist.CFBundleExecutable,
			)
			minKernel = kextPlist.libkernVersion()
		}
	}

	return Kext{
		Arch:           "Any", // TODO: not sure how to determine this outside of lipo and kextfind
		BundlePath:     path.Base(f.Name()),
		Comment:        "",
		Enabled:        true,
		ExecutablePath: executablePath,
		MaxKernel:      "",
		MinKernel:      minKernel,
		PlistPath:      plistPath,
	}, nil
}

// ACPI returns a ACPI.Add struct for the provided compiled .aml ACPI table
// file.
func (Helpers) ACPI(filename string) (ACPIEntry, error) {
	f, err := os.Open(filename)

	if err != nil {
		return ACPIEntry{}, fmt.Errorf("add_acpi: %w", err)
	}

	return ACPIEntry{
		Comment: "",
		Enabled: true,
		Path:    path.Base(f.Name()),
	}, nil
}

// Tool returns a Misc.Tools struct for the provided .efi tool file.
func (Helpers) Tool(filename string) (Tool, error) {
	f, err := os.Open(filename)

	if err != nil {
		return Tool{}, fmt.Errorf("tool: %w", err)
	}

	nameExt := path.Base(f.Name())

	return Tool{
		Arguments:       "",
		Auxiliary:       true, // TODO: I guess all tools being auxiliary is fine for now
		Comment:         "",
		Enabled:         true,
		Flavour:         "Auto",
		FullNvramAccess: slices.Contains(fullNvramAccessTools, nameExt),
		Name:            strings.TrimSuffix(nameExt, path.Ext(nameExt)),
		Path:            nameExt,
		RealPath:        false,
		TextMode:        false,
	}, nil
}

// Driver returns a UEFI.Drivers struct for the provided .efi driver file.
func (Helpers) Driver(filename string) (Driver, error) {
	f, err := os.Open(filename)

	if err != nil {
		return Driver{}, fmt.Errorf("driver: %w", err)
	}

	nameExt := path.Base(f.Name())

	return Driver{
		Comment:   "",
		Enabled:   true,
		Path:      nameExt,
		LoadEarly: slices.Contains(loadEarlyDrivers, nameExt),
		Arguments: "",
	}, nil
}
