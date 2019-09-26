# Monolith

Hackintosh configuration and kexts.

## Relevant Hardware

* **CPU**: Intel Core i9-9900k (Coffee Lake)
    * iGPU: Intel UHD 630
* **Motherboard**: Gigabyte Z390 I Aorus Pro WiFi
    * Audio: Realtek ALC1220
* **GPU**: MSI RX Vega 64 AIR Boost

## Usage

1. Create a `configuration/serials.yml` file with the following key-value pairs:
    * `MLB`
    * `BoardSerialNumber`
    * `SerialNumber`
    * `SmUUID`
2. Run `configuration/make_config`
