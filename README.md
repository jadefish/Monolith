# Monolith

Hackintosh configuration and kexts.

## Relevant Hardware

* **CPU**: Intel Core i9-9900k (Coffee Lake)
    * iGPU: Intel UHD 630
* **Motherboard**: ASUS ROG Strix Z370-I Gaming
    * Audio: Realtek ALC1220 ("SupremeFX S1220A")
* **GPU**: MSI RX Vega 64 AIR Boost

## Usage

1. Create a `configuration/serials.yml` file with the following key-value pairs:
    * `MLB`
	* `SystemSerialNumber`
	* `SystemUUID`
2. Run `configuration/make_config`
