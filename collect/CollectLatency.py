import pandas as pd
import re
import argparse


parser = argparse.ArgumentParser(description="Process latency log and generate performance comparison CSV.")
parser.add_argument('input_filename', type=str, help="Path to the input result.log file")

args = parser.parse_args()

with open(args.input_filename, "r") as file:
    log_content = file.readlines()

pattern = r"([a-zA-Z ]+ version) (.*) cost avg= ([0-9\.]+) (ms|us)"


data = []

for line in log_content:
    match = re.search(pattern, line)
    if match:
        version = match.group(1).strip()
        test_item = match.group(2).strip()
        cost = match.group(3)
        unit = match.group(4).strip()

        cost_with_unit = f"{cost} {unit}"

        data.append([test_item, version, cost_with_unit])


df = pd.DataFrame(data, columns=["test items", "version", "cost"])


df["version"] = df["version"].apply(lambda x: "Vanilla" if "Vanilla" in x else "Modified")

df_pivot = df.pivot(index="test items", columns="version", values="cost").reset_index()

df_pivot["overhead"] = (
    (df_pivot["Modified"].str.extract(r"([0-9\.]+)")[0].astype(float) -
     df_pivot["Vanilla"].str.extract(r"([0-9\.]+)")[0].astype(float)) /
    df_pivot["Vanilla"].str.extract(r"([0-9\.]+)")[0].astype(float)) * 100

df_pivot.columns = ["test items", "vanilla", "modified", "overhead"]

output_filename = "latency.csv"

df_pivot.to_csv(output_filename, index=False)

print(df_pivot)