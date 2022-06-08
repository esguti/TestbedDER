# Testbed

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This repository contains a testbed to perform and detect Person in The Middle (PiTM) attack on a Distributed Energy Resources (DER) system that uses [SunSpec](https://sunspec.org/) Modbus specification.

The testbed is composed of 4 virtual nodes:
 - sunspec-battery: Lithium-ion Battery Bank SunSpec model number 803
 - sunspec-hmi: Human Machine Interface (HMI)
 - sunspec-kali: kali linux image to perform de attack
 - sunspec-snort: kali linux image to performa de detection

In this attack, the communication between the HMI (*sunspec-HMI*) and the Battery (*sunspec-battery*) is intercepted by the attacker (*sunspec-kali*). Messages sent from the battery to the HMI are modified by replacing the original temperature values with fake values created by the attacker. Then, the Intrusion Detection System (*sunspec-snort*) detects the attack raising an alarm.


The SunSpec client/server code is a derivative work of the code from TRICERA-enery https://github.com/TRICERA-energy/sunspec, which is licensed Apache-2.0.

# Usage

1. Start the containers

```bash
docker-compose build
docker-compose up
```
2. Open the hmi browser page

```bash
firefox http://localhost:8080/
```

3. Connect to *sunspec-snort* and start the detection

```bash
docker exec -ti sunspec-snort bash
snort -i br-07f3d23ed18d -c /etc/snort/snort.conf -A console
```

NOTE: *br-07f3d23ed18d* is the network interface

4. Connect to *sunspec-kali* container and execute the attack

```bash
docker exec -ti sunspec-kali bash
./start_injection.sh
```
Now, you will display the new temperature value (-10) in the browser page.

It will be also displayed alerts in the IDS console

5. Connect to *sunspec-kali* container and stop the attack

```bash
docker exec -ti sunspec-kali bash
./stop_injection.sh
```
The real temperature value transmitted (30) will be again displayed in the browser page
