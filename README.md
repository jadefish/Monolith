# Monolith

Hackintosh configuration and files.

## Hardware

* **CPU**: Intel Core i9-9900k (Coffee Lake)
    * iGPU: Intel UHD 630
* **Motherboard**: ASUS ROG Strix Z370-I Gaming
    * Audio: Realtek ALC1220 ("SupremeFX S1220A")
    * 1x USB-C connector, with switch (type 9)
* **GPU**: AMD Radeon RX 6800 XT


## Unsupported

* Wi-Fi
* Bluetooth


## Software

* **Bootloader**: [OpenCore](https://github.com/acidanthera/opencorepkg)
* **OS**: macOS 12 Monterey


## Requirements

* plutil
* [iasl](https://github.com/RehabMan/Intel-iasl)
* [ocvalidate](https://github.com/acidanthera/OpenCorePkg/tree/master/Utilities/ocvalidate#readme)

## Usage

1. Obtain an OpenCore `Sample.plist` file via an OpenCore debug or release build.
2. Download all required kexts.
3. Compile all required SSDT DSL files via `iasl`.
4. Patch the sample plist as `config.plist`:

```
$ ./generate-plist.sh Sample.plist \
    --mlb my-mlb-value \
    --rom base64-rom-data \
    --serial my-serial-number \
    --uuid my-uuid \
    --product system-product-string \
    --ssdt ssdt/*.aml \
    --kext kext/Lilu.kext kext/*.kext \
    --tool tool/*.efi \
    --driver driver/*.efi > config.plist
```
