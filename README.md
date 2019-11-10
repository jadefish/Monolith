# Monolith

Hackintosh configuration and files.

## Hardware

* **CPU**: Intel Core i9-9900k (Coffee Lake)
    * iGPU: Intel UHD 630
* **Motherboard**: ASUS ROG Strix Z370-I Gaming
    * Audio: Realtek ALC1220 ("SupremeFX S1220A")
* **GPU**: MSI RX Vega 64 AIR Boost

## Software

* **Bootloader**: [OpenCore](https://github.com/acidanthera/opencorepkg)
* **OS**: macOS 10.15 Catalina

## Usage

1. Create a `configuration/serials.yml` file with the following key-value pairs:
    * `MLB`
	* `SystemSerialNumber`
	* `SystemUUID`
2. Run `configuration/make_config`
3. Compile SSDT DSL files via [MaciASL.app](https://github.com/acidanthera/MaciASL)
4. Copy `EFI` directory, generated `config.plist` file(s), and compiled SSDT
   `aml` files to your EFI system partition
