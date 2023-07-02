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

* **Bootloader**: [OpenCore](https://github.com/acidanthera/opencorepkg) 0.9.3
* **OS**: macOS 13 Ventura


## Requirements

* plutil
* [iasl](https://github.com/RehabMan/Intel-iasl)
* [ocvalidate](https://github.com/acidanthera/OpenCorePkg/tree/master/Utilities/ocvalidate#readme)


## Usage

1. Obtain an OpenCore `Sample.plist` file via an OpenCore debug or release build.
2. Download all required kexts and copy into `./kext`.
3. Download all required drivers and copy into `./driver`.
4. Download all required tools and copy into `./tool`.
5. Copy all required SSDTs from the downloaded OpenCore build (and wherever else) and place into ``./ssdt`.
6. Patch the sample plist as `config.plist` via `make` or the following:
    ```bash
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

Alternatively, MLB, ROM, serial number, UUID, and the system product string can be provided via environment variables if not overridden with arguments:

```bash
$ env
MLB=my-mlb-value
ROM=base64-rom-data
SERIAL=my-serial-number
UUID=my-uuid
PRODUCT=system-product-string

$ ./generate-plist.sh Sample.plist \
    --ssdt ssdt/*.aml \
    --kext kext/Lilu.kext kext/*.kext \
    --tool tool/*.efi \
    --driver driver/*.efi > config.plist
```
