# Available Attestation Repository

This repository contains the implementation for the paper "Available Attestation: Towards a Reorg-Resilient Solution for Ethereum Proof-of-Stake" ([eprint](https://eprint.iacr.org/2025/097.pdf)). It includes a modified Ethereum PoS protocol and the attacks analyzed in the paper.

The repository establishes two local testnets, each pre-configured with 16,384 validators:
- The **vanilla client** is compiled from Prysm [v5.0](https://github.com/prysmaticlabs/prysm/releases/tag/v5.0.0) (see `normal_prysm` for the vanilla Ethereum PoS protocol implementation).  
- The **modified client** is based on a modified version of Prysm [v5.0](https://github.com/prysmaticlabs/prysm/releases/tag/v5.0.0) (see `modified_prysm` for the modified Ethereum PoS protocol implementation).  

---

## 1. Ethical Considerations

All experiments are conducted on a local testnet (running on a single local machine). No experiments are conducted on the live Ethereum network. This repository does not uncover new vulnerabilities but instead analyzes known malicious reorganization attacks. We will not provide any additional exploit information.

---
## 2. Requirements

### 2.1. Hardware Dependencies

These experiments do not require any specialized hardware. The computer used in our setup is configured as follows:
- **CPU:** 4-core
- **RAM:** 16 GB
- **Storage:** 100 GB
- **Network Bandwidth:** 100 Mbps

---

### 2.2. Software Dependencies

#### 2.2.1. Docker

We use Docker to run our experiments. You can install Docker by following the instructions provided in the [official Docker documentation](https://docs.docker.com/engine/install/). Ensure that your Docker Engine version is at least **Docker 24**.

#### 2.2.2. Python Requirements

Python is used for data processing and plotting. Ensure that your Python version is at least **3.10**. Install the required Python packages by running the following command:

```shell
pip3 install -r requirements.txt
```

---

## 3. Run the Experiments Step by Step

We conduct three experiments for both vanilla Ethereum PoS protocol and modified Ethereum PoS protocol. 

* Reorg resilience experiments
* Throughput experiments
* Latency experiments

---


### 3.1. Build the Docker Image

After entering the repository directory (referred to as ``$HOME``), run the following script to build the Docker image:

```shell
./build.sh
```

---


### 3.2. Reorg resilience experiments 

#### 3.2.1 Run All Attacks at Once

To conduct five attacks (i.e., exante reorg attack, sandwich reorg attack, unrealized justification reorg attack, justification withholding attack, and staircase attack) simultaneously, run the command:

```shell
./runtest.sh reorg
```

The experiments will run each attack for 9000 seconds (approximately 25 hours in total) across both protocols. After completion, the results can be found in the ``$HOME`` directory. The number of reorg blocks for the modified protocol should match Figure 13 in the paper (i.e., the number of reorg blocks in the modified protocol is zero). 

#### 3.2.2. Run Individual Attacks

If one does not want to wait for experiments that long, run the following commands to conduct each attack separately. Each attack will last five hours. 

Run the modified exante reorg attack:

```shell
./runtest.sh 1
```

Run the sandwich reorg attack:

```shell
./runtest.sh 2
```

Run the unrealized justification reorg attack:

```shell
./runtest.sh 3
```

Run the justification withholding reorg attack:

```shell
./runtest.sh 4
```

Run the staircase attack:

```shell
./runtest.sh 5
```

---

### 3.3. Throughput experiment

Run the throughput experiment using the following command:

```shell
./runtest.sh tps
```

This experiment will take approximately two hours. After completion, the results can be found in the ``$HOME`` directory. The throughput of the modified protocol should match Figure 14 in the paper (i.e., the throughput of the vanilla protocol and the modified protocol are almost the same).

---

### 3.4. Latency experiment

Run the latency experiment using the following command:

```shell
./runtest.sh latency
```

This experiment will take approximately ten minutes. After completion, the results can be found in the ``$HOME`` directory. The latency of the modified protocol should match Figure 15 in the paper (i.e., the latency of the vanilla protocol and the modified protocol are almost the same).

---

### 3.5. Expected Output

The output of each experiment should look like this:


```
[+] Running 19/19
 ✔ Network Created                                                  0.1s 
 ✔ Container execute5          Started                              0.4s 
 ✔ Container execute3          Started                              0.4s 
 ✔ Container execute1          Started                              0.4s 
 ✔ Container ethmysql          Started                              0.4s 
 ✔ Container execute2          Started                              0.4s 
 ✔ Container execute4          Started                              0.2s 
 ✔ Container beacon-2          Started                              0.6s 
 ✔ Container attacker-1        Started                              0.4s 
 ✔ Container beacon-3          Started                              0.4s 
 ✔ Container beacon-1          Started                              0.6s 
 ✔ Container beacon-4          Started                              0.4s 
 ✔ Container beacon-5          Started                              0.6s 
 ✔ Container validator-4       Started                              0.6s 
 ✔ Container validator-3       Started                              0.8s 
 ✔ Container strategy          Started                              0.8s 
 ✔ Container validator-1       Started                              0.8s 
 ✔ Container validator-2       Started                              0.8s 
 ✔ Container validator-5       Started                              0.8s
```

## 4. Data Availability

The raw log files from the experiments presented in the paper are located in the ``data`` folder. The data used to generate the figures in the paper are available in the ``result`` folder. 