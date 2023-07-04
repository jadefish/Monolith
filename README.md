# Monolith

Hackintosh configuration and files.

## Usage

1. Obtain an OpenCore `Sample.plist` file via an OpenCore debug or release build
2. Download all required kexts, drivers, and tools
3. Compile all necessary ACPI tables via `iasl` and keep the `.aml` files handy
4. Write an [instructions script](https://github.com/jadefish/Monolith/wiki/Instructions-scripts)
5. Generate the output `config.plist`:
   ```bash
   $ ./monolith --product $PRODUCT \
         --mlb $MLB \
         --rom $ROM \
         --serial $SERIAL \
         --uuid $UUID \
         Sample.plist \
         instructions > config.plist
   ```

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


## Building

Install Go 1.20 or later, then simply build:

```bash
$ go build
```
