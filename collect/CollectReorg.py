import re
import matplotlib.pyplot as plt
import numpy as np
from collections import defaultdict
import argparse

parser = argparse.ArgumentParser(description="Parse log file and plot TPS vs TX Rate.")
parser.add_argument('logfile_path', type=str, help="Path to the log file")

args = parser.parse_args()
logfile_path = args.logfile_path

data = defaultdict(lambda: defaultdict(lambda: 0))

attack_pattern = r"attack-([a-z\-]+) (Vanilla|Modified) version reorg event info:"
event_pattern = r'time=".*?" level=info msg="reorg event" .* totalReorgDepth=(\d+)'

current_attack = None
current_version = None

with open(logfile_path, "r") as file:
    for line in file:
        attack_match = re.search(attack_pattern, line)
        if attack_match:
            current_attack = attack_match.group(1)
            current_version = attack_match.group(2).lower()
            if current_version and current_attack:
                data[current_version][current_attack] = 0
            continue
        
        event_match = re.search(event_pattern, line)
        if event_match and current_attack and current_version:
            total_reorg_depth = int(event_match.group(1))
            data[current_version][current_attack] = total_reorg_depth

print("Final Data Extracted:")
for version, attacks in data.items():
    for attack, total_reorg_depth in attacks.items():
        print(f"Version: {version}, Attack: {attack}, Total Reorg Depth: {total_reorg_depth}")

ordered_attacks = ["exante", "sandwich", "unrealized", "withholding", "staircase"]

vanilla_values = [data["vanilla"].get(attack, 0) for attack in ordered_attacks]
modified_values = [data["modified"].get(attack, 0) for attack in ordered_attacks]

x = np.arange(len(ordered_attacks))
width = 0.4

fig, ax = plt.subplots(figsize=(8, 6))
bars1 = ax.bar(x - width / 2, vanilla_values, width, label="vanilla protocol", color='#fbcccb', edgecolor='#f54747')
bars2 = ax.bar(x + width / 2, modified_values, width, label="modified protocol", color='#ccccff', edgecolor='#6e6eff')

for bar in bars1:
        ax.text(bar.get_x() + bar.get_width() / 2, bar.get_height() + 0.5, f'{int(bar.get_height())}', 
                ha='center', va='bottom', fontsize=10)

for bar in bars2:
        ax.text(bar.get_x() + bar.get_width() / 2, bar.get_height() + 0.5, f'{int(bar.get_height())}', 
                ha='center', va='bottom', fontsize=10)


ax.set_xlabel("Attacks", fontsize=12)
ax.set_ylabel("Number of Reorganized Blocks", fontsize=12)
ax.set_xticks(x)
ax.set_xticklabels(ordered_attacks)
ax.legend(loc="upper right", fontsize=10)
ax.set_ylim(0, max(vanilla_values) + 10)

ax.grid(axis="y", linestyle="--", alpha=0.7)

plt.tight_layout()
plt.savefig("reorg.png", dpi=300)
plt.show()