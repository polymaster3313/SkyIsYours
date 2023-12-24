# SkyIsYours
A Scanner and Bruter armed with arts of work distribution

## Features
```
smart scanning
fast brute force speed (17-21 ssh brute per second)
fast scanning speed
load distribution
fault tolerant sockets
interprocess mode
SSH
```

## The Sky architecture
### Load balancing
```mermaid
graph TD;
    SkyScanner --> |126.207.193.116| Skybruter1;
    SkyScanner --> |159.27.168.28| Skybruter2;
    SkyScanner --> |156.213.208.17| Skybruter3;
    SkyScanner --> |11.77.180.157| Skybruter4;
```

## SkyBruter
### rapid wave

Skybruter initiates a rapidwave by spawning a limited pool of up to 200 goroutines, each attempting to connect to the server with distinct passwords. In case of ratelimiting, the goroutine joins a slowwave queue. If the result indicates an invalid password or ratelimiting, the goroutine relinquishes its position, allowing other goroutines with different passwords to join. Upon discovering the correct password, no new goroutines are allowed, and slowwave is canceled. If unsuccessful, slowwave validates ratelimited passwords.

```mermaid
graph TD
    subgraph SkyBruter
        A[rapidwave] -->|Spawns Pool| B[goroutine pool 200 rate]
        B -->|1| C[password123]
        B -->|2| D[weakpassword]
        B -->|3| E[Emily]
        B -->|4| I[123456]
        B -->|...| F[...]
        B -->|200| G[balls]
        C -->|rate limited| I1[slowwave queue]
        I -->|rate limited| I1[slowwave queue]

    end

    subgraph Server
        H[Target]
        C -->|Connects to| H
        D -->|Connects to| H
        E -->|Connects to| H
        F -->|Connects to| H
        G -->|Connects to| H
        I -->|Connects to| H

    end
```

### slowwave
Slowwave is designed to validate passwords obtained from a ratelimited source by employing a traditional SSH brute-force approach. It systematically attempts each password one by one, determining their correctness and confirming their validity. The process involves accessing the slowqueue, where potential passwords are stored. Upon retrieving a password, Slowwave establishes a connection to the target. If the password is correct, it proceeds to process the result . In the event of an incorrect password, the cycle continues by reattempting passwords through the slowqueue, Until no password is found.

```mermaid
graph TD
    subgraph SkyBruter
        A[Slowwave] -->|access| B[slowqueue]
        B -->|If empty| O[No password found]
        B -->|receive password| C[password]
        C -->|connects target| D[target]
        D -->|Correct| E[Password found]
        E --> O
        D -->|Wrong| B
    end
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
