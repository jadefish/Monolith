package main

// ACPIEntry is an ACPI table entry in OpenCore's "ACPI.Add" list.
type ACPIEntry struct {
	Comment string
	Enabled bool
	Path    string
}

// Kext is kernel extension entry in OpenCore's "Kernel.Add" list.
type Kext struct {
	Arch           string
	BundlePath     string
	Comment        string
	Enabled        bool
	ExecutablePath string
	MaxKernel      string
	MinKernel      string
	PlistPath      string
}

// Tool is an entry in OpenCore's "Misc.Tools" list.
type Tool struct {
	Arguments       string
	Auxiliary       bool
	Comment         string
	Enabled         bool
	Flavour         string
	FullNvramAccess bool
	Name            string
	Path            string
	RealPath        bool
	TextMode        bool
}

// Driver is a UEFI driver entry in OpenCore's "UEFI.Drivers" list.
type Driver struct {
	Comment   string
	Enabled   bool
	Path      string
	LoadEarly bool
	Arguments string
}
