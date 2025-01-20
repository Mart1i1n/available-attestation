import re
import matplotlib.pyplot as plt
import argparse

parser = argparse.ArgumentParser(description="Parse log file and plot TPS vs TX Rate.")
parser.add_argument('logfile_path', type=str, help="Path to the log file")

args = parser.parse_args()
logfile_path = args.logfile_path

data = {
    "vanilla": {"tps": [], "tx_rate": []},
    "modified": {"tps": [], "tx_rate": []}
}

version_pattern = r"(Vanilla|Modified) version"
tps_rate_pattern = r"tps=(\d+)\s+tx rate=(\d+)"

current_version = None

with open(logfile_path, "r") as file:
    for line in file:
        version_match = re.search(version_pattern, line)
        if version_match:
            current_version = version_match.group(1).lower()

        tps_rate_match = re.search(tps_rate_pattern, line)
        if tps_rate_match and current_version:
            tps = int(tps_rate_match.group(1))
            tx_rate = int(tps_rate_match.group(2))
            data[current_version]["tps"].append(tps)
            data[current_version]["tx_rate"].append(tx_rate)

plt.figure(figsize=(8, 6))
plt.plot(data["vanilla"]["tx_rate"], data["vanilla"]["tps"], 'r-o', label="vanilla protocol")
plt.plot(data["modified"]["tx_rate"], data["modified"]["tps"], 'b-^', label="modified protocol")

plt.xlabel("Transaction Input Rate (tx/s)")
plt.ylabel("Transaction Performance (tx/s)")
plt.legend(loc="lower right")
plt.grid(True)

plt.savefig("tps", dpi=300, bbox_inches='tight')

plt.show()