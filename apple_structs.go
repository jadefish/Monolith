package main

// KextPropertyList contains selected property list entries for the Info.plist
// file contained within a kernel extension (.kext) bundle.
type KextPropertyList struct {
	CFBundleExecutable     string
	OSBundleLibrariesX8664 *map[string]string `plist:"OSBundleLibraries_x86_64,omitempty"`
	OSBundleLibraries      map[string]string
}

// libkernVersion returns the minimum kernel version required for this kext to
// be loaded.
//
// If the kext does not provide such info, "8.0.0" is returned.
func (kpl KextPropertyList) libkernVersion() string {
	if kpl.OSBundleLibrariesX8664 != nil {
		if value, ok := (*kpl.OSBundleLibrariesX8664)["com.apple.kpi.libkern"]; ok {
			return value
		}
	} else if value, ok := kpl.OSBundleLibraries["com.apple.kpi.libkern"]; ok {
		return value
	}

	return "8.0.0"
}
