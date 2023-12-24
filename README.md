# SkyIsYours
A Scanner and Bruter armed with arts of work distribution

## Features
```
smart scanning
fast brute force speed (17 ssh brute per second)
fast scanning speed
load distribution
fault tolerant sockets
interprocess mode
```

## The Sky architecture
```mermaid
graph TD;
    SkyScanner-->Skybruter1;
    SkyScanner-->Skybruter2;
    SkyScanner-->Skybruter3;
    SkyScanner-->Skybruter4;
```
## BENCHMARK

system information
```
CPU: 12th Gen Intel i7-12700H (20) @ 2.688GHz
Cores: 14
OS: Arch
arch: Amd64
RAM: 16gb
bandwidth: 100mbps
```


### Benchmark 1000 pass dictionary attack

![Chart (1)](https://github.com/polymaster3313/SkyIsYours/assets/93959737/9091ed91-da20-4c66-85ab-1777bcbfc607)
