# Monolith

Hackintosh configuration and kexts.

## Relevant Hardware

* **CPU**: Intel Core i9-9900k
    * iGPU: Intel UHD 630
* **Motherboard**: Gigabyte Z390 I Aorus Pro WiFi
    * Audio: Realtek ALC1220
* **GPU**: MSI RX Vega 64 AIR Boost

## Usage

1. Create a `serials.yaml` file with the following key-value pairs:
    * `MLB`
    * `BoardSerialNumber`
    * `SerialNumber`
    * `SmUUID`
2. Run `make_config`
