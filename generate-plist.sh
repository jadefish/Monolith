#!/usr/bin/env bash
set -Eeuo pipefail

# return EXIT_SUCCESS if key $1 exists.
exists() {
	local key="${1}"

	plutil -type "${key}" -- "${plist_file}" > /dev/null
}

# if key $1 already exists, replace it with $@.
# otherwise, insert it.
setval() {
	local key="${1}"
	shift

	if ! exists "${key}" || [[ $* =~ '-append' ]]; then
		plutil -insert "${key}" "$@" -- "${plist_file}"
	else
		plutil -replace "${key}" "$@" -- "${plist_file}"
	fi
}

# remove key $1 if it exists.
delete() {
	local key="${1}"

	if exists "${key}"; then
		plutil -remove "${key}" -- "${plist_file}"
	fi
}

show_help() {
	cat <<-'EOF'
	Usage: generate-plist FILE [OPTION]...
	Patch an OpenCore plist FILE.

	  -h, --help   display this help and exit
	      --debug  include debugging configuration:
	               Target, DisplayLevel, AppleDebug, ApplePanic, DisableWatchDog

	  --mlb STRING      use STRING as the value for PlatformInfo.Generic.MLB
	  --rom DATA        use DATA as the value for PlatformInfo.Generic.ROM
	  --serial STRING   use STRING as the value for PlatformInfo.Generic.SystemSerialNumber
	  --uuid STRING     use STRING as the value for PlatformInfo.Generic.SystemUUID
	  --product STRING  use STRING as the value for PlatformInfo.Generic.SystemProductName

	  The following options can be specified multiple times and in any order:
	  --ssdt FILE    include a .aml compiled SSDT FILE
	  --kext FILE    include a .kext kernel extension FILE
	  --tool FILE    include a .efi EFI tool FILE
	  --driver FILE  include a .efi drive FILE

	  SSDT AMLs, kextts, tools, and drivers are added in the order provided. If
	  one kext must be loaded before another, ensure it is passed earlier in
	  the list of kext arguments.
	EOF
}

err() {
	echo "${1}" >&2
}

main() {
	plist_file="${1:-""}" # warning: global (used in functions above)

	if [[ $plist_file = '--help' || $plist_file = '-h' || ${#@} -lt 1 ]]; then
		show_help
		return 0
	fi

	if [[ ! -f "${plist_file}" ]]; then
		err "plist file ${plist_file} does not exist."
		return 1
	fi

	tmp="$(mktemp)"
	cat "${plist_file}" > "${tmp}"
	plist_file="${tmp}"
	shift

	local ssdts=()
	local kexts=()
	local drivers=()
	local tools=()
	local mlb="${MLB:-""}"
	local serial_number="${SERIAL_NUMBER:-""}"
	local uuid="${UUID:-""}"
	local rom="${ROM:-""}"
	local product="${PRODUCT:-""}"
	local debug=false

	# Need to swap IFS to ensure arguments with spaces are not split to words
	# when checking $items for duplicates in the *) case below.
	OLD_IFS="${IFS}"
	IFS="|"

	local type=""
	while (("$#")); do
		local arg="${1}"

		case "${arg}" in
		--debug)
			err "** Enabling debug options"
			debug=true
			;;

		# secrets:
		--mlb) shift; mlb="${1}";;
		--serial_number) shift; serial_number="${1}";;
		--uuid) shift; uuid="${1}";;
		--rom) shift; rom="${1}";;
		--product) shift; product="${1}";;

		# begin adding files of this type:
		--ssdt)   type="ssdt";;
		--kext)   type="kext";;
		--tool)   type="tool";;
		--driver) type="driver";;

		# file names:
		*)
			file="$(basename "${arg}")"
			case "${type}" in
			ssdt)
				# shellcheck disable=SC2076
				# https://stackoverflow.com/a/61551944
				if [[ ! "${IFS}${ssdts[*]+"${ssdts[*]}"}${IFS}" =~ "${IFS}${file}${IFS}" ]]; then
					ssdts+=("${file}")
					err "SSDT: ${file}"
				fi
				;;
			kext)
				# shellcheck disable=SC2076
				if [[ ! "${IFS}${kexts[*]+"${kexts[*]}"}${IFS}" =~ "${IFS}${file}${IFS}" ]]; then
					kexts+=("${file}")
					err "kext: ${file}"
				fi
				;;
			tool)
				# shellcheck disable=SC2076
				if [[ ! "${IFS}${tools[*]+"${tools[*]}"}${IFS}" =~ "${IFS}${file}${IFS}" ]]; then
					tools+=("${file}")
					err "tool: ${file}"
				fi
				;;
			driver)
				# shellcheck disable=SC2076
				if [[ ! "${IFS}${drivers[*]+"${drivers[*]}"}${IFS}" =~ "${IFS}${file}${IFS}" ]]; then
					drivers+=("${file}")
					err "driver: ${file}"
				fi
				;;
			esac
		esac

		shift
	done
	IFS="${OLD_IFS}"

	delete "#WARNING - 1"
	delete "#WARNING - 2"
	delete "#WARNING - 3"
	delete "#WARNING - 4"

	## ACPI:
	setval 'ACPI.Add' -array
	setval 'ACPI.Delete' -array
	setval 'ACPI.Patch' -array

	for ((i = 0; i < ${#ssdts[@]}; i++)) do
		local file="${ssdts[$i]}"

		setval "ACPI.Add.$i" -dictionary
		setval "ACPI.Add.$i.Enabled" -bool true
		setval "ACPI.Add.$i.Path" -string "${file}"
		setval "ACPI.Add.$i.Comment" -string ""
	done

	## Booter:
	setval 'Booter.MmioWhitelist' -array
	setval 'Booter.Patch' -array
	setval 'Booter.Quirks.DisableSingleUser' -bool true
	setval 'Booter.Quirks.ProvideCustomSlide' -bool false
	setval 'Booter.Quirks.EnableSafeModeSlide' -bool false

	## DeviceProperties:
	setval 'DeviceProperties.Add' -dictionary
	setval 'DeviceProperties.Add.PciRoot(0x0)/Pci(0x1f,0x3)' -dictionary
	setval 'DeviceProperties.Add.PciRoot(0x0)/Pci(0x1f,0x3).layout-id' -integer 21
	setval 'DeviceProperties.Delete' -dictionary

	## Kernel:
	setval 'Kernel.Add' -array
	setval 'Kernel.Block' -array
	setval 'Kernel.Force' -array
	setval 'Kernel.Patch' -array

	for ((i = 0; i < ${#kexts[@]}; i++)) do
		local file="${kexts[$i]}"
		local plist_path='Contents/Info.plist'
		local executable_path="Contents/MacOS/${file%.*}"
		local has_executable=true

		if [[ -d $file ]]; then
			# kext actually exists; check if it's plist-only:
			has_executable="$(test -f "${file}/${executable_path}")"
		fi

		setval "Kernel.Add.$i" -dictionary
		setval "Kernel.Add.$i.Arch" -string 'Any'
		setval "Kernel.Add.$i.BundlePath" -string "${file}"
		setval "Kernel.Add.$i.Comment" -string ''
		setval "Kernel.Add.$i.Enabled" -bool true

		if [ $has_executable = true ]; then
			setval "Kernel.Add.$i.ExecutablePath" -string "${executable_path}"
		fi

		setval "Kernel.Add.$i.PlistPath" -string "${plist_path}"
		setval "Kernel.Add.$i.MinKernel" -string "8.0.0"
		setval "Kernel.Add.$i.MaxKernel" -string ""
	done

	setval 'Kernel.Quirks.DisableIoMapper' -bool true
	setval 'Kernel.Quirks.ExtendBTFeatureFlags' -bool true
	setval 'Kernel.Quirks.PanicNoKextDump' -bool true

	## Misc:
	setval 'Misc.BlessOverride' -array
	setval 'Misc.Boot.HideAuxiliary' -bool false
	setval 'Misc.Boot.LauncherOption' -string 'Full'
	setval 'Misc.Boot.PickerMode' -string 'External'
	# USE_VOLUME_ICON | USE_POINTER_CONTROL | USE_FLAVOUR_ICON
	setval 'Misc.Boot.PickerAttributes' -integer "$((16#91))"
	setval 'Misc.Boot.PollAppleHotKeys' -bool true
	setval 'Misc.Boot.Timeout' -integer 10
	setval 'Misc.Entries' -array
	setval 'Misc.Security.AllowSetDefault' -bool true
	# FILE_SYSTEM_LOCK | DEVICE_LOCK
	# ALLOW_FS_APFS | ALLOW_FS_ESP
	# ALLOW_DEVICE_SATA | ALLOW_DEVICE_NVME | ALLOW_DEVICE_USB
	setval 'Misc.Security.ScanPolicy' -integer $((16#1290503))
	setval 'Misc.Security.Vault' -string 'Optional' # TODO

	# Debugging:
	local target=0
	local display_level=0
	if [ $debug = true ]; then
		target=$((16#47)) # enable | console | file
		display_level=$((16#80000042)) # DEBUG_WARN | DEBUG_INFO | DEBUG_ERROR
	fi
	setval 'Misc.Debug.AppleDebug' -bool $debug
	setval 'Misc.Debug.ApplePanic' -bool $debug
	setval 'Misc.Debug.DisableWatchDog' -bool $debug
	setval 'Misc.Debug.DisplayLevel' -integer $display_level
	setval 'Misc.Debug.SysReport' -bool $debug
	setval 'Misc.Debug.Target' -integer $target
	setval 'Misc.Security.AllowNvramReset' -bool $debug
	setval 'Misc.Security.DmgLoading' -string 'Disabled'

	setval 'Misc.Tools' -array
	for ((i = 0; i < ${#tools[@]}; i++)) do
		local file="${tools[$i]}"
		local name="${file%.*}"

		setval "Misc.Tools.$i" -dictionary
		setval "Misc.Tools.$i.Arguments" -string ""
		setval "Misc.Tools.$i.Auxiliary" -bool true
		setval "Misc.Tools.$i.Comment" -string ""
		setval "Misc.Tools.$i.Enabled" -bool true
		setval "Misc.Tools.$i.Flavour" -string "Auto"
		setval "Misc.Tools.$i.Name" -string "${name}"
		setval "Misc.Tools.$i.Path" -string "${file}"
		setval "Misc.Tools.$i.RealPath" -bool false
		setval "Misc.Tools.$i.TextMode" -bool false
	done

	## NVRAM:
	local boot_args=('agdpmod=pikera')
	if [ $debug = true ]; then
		boot_args+=('-v')
		boot_args+=('keepsyms')
		boot_args+=('debug=0x122')
		boot_args+=('-liludbg')
		boot_args+=('-alcdbg')
	fi
	setval "NVRAM.Add.7C436110-AB2A-4BBB-A880-FE41995C9F82.boot-args" -string "${boot_args[*]}"
	setval 'NVRAM.LegacySchema' -dictionary

	## PlatformInfo:
	setval 'PlatformInfo.Generic.AdviseFeatures' -bool true
	setval 'PlatformInfo.Generic.SystemProductName' -string "${product}"
	setval 'PlatformInfo.Generic.MLB' -string "${mlb}"
	setval 'PlatformInfo.Generic.ROM' -data "${rom}"
	setval 'PlatformInfo.Generic.SystemSerialNumber' -string "${serial_number}"
	setval 'PlatformInfo.Generic.SystemUUID' -string "${uuid}"

	## UEFI:
	setval 'UEFI.Audio.AudioCodec' -integer 0
	setval 'UEFI.Audio.AudioDevice' -string 'PciRoot(0x0)/Pci(0x1f,0x3)'
	setval 'UEFI.Audio.AudioSupport' -bool true
	setval 'UEFI.Audio.SetupDelay' -integer $((500*1000)) # 500 ms

	setval 'UEFI.Drivers' -array
	for ((i = 0; i < ${#drivers[@]}; i++)) do
		local file="${drivers[$i]}"

		setval "UEFI.Drivers.$i" -dictionary
		setval "UEFI.Drivers.$i.Arguments" -string ""
		setval "UEFI.Drivers.$i.Comment" -string ""
		setval "UEFI.Drivers.$i.Enabled" -bool true
		setval "UEFI.Drivers.$i.Path" -string "${file}"
	done

	setval 'UEFI.ReservedMemory' -array

	err 'lint: '
	plutil -lint "${plist_file}" >&2
	err 'validate: '
	ocvalidate "${plist_file}" >&2
	cat "${plist_file}"
}

main "$@"
