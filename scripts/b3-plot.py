"""Plot class weights and predicted densities from the checked B3 inventory."""
from collections import Counter
from fractions import Fraction
import json
from pathlib import Path

import matplotlib
matplotlib.use("Agg")
from matplotlib import pyplot as plt

directory = Path("artifacts/b3")
groups = json.loads((directory / "certificates.json").read_text(encoding="utf-8"))["groups"]
colors = ["#315a84", "#006e64", "#8b4a80", "#9b501c", "#515765"]
figure, (classes, subsets) = plt.subplots(2, 1, figsize=(8.6, 5.7), layout="constrained")
for i, group in enumerate(groups):
    total = len(group["product"])
    at = Fraction(0)
    for color, members in zip(colors, group["classes"]):
        width = Fraction(len(members), total)
        classes.barh(i, float(width), left=float(at), color=color, height=0.62,
                     edgecolor="white", linewidth=1.5)
        classes.text(float(at + width / 2), i, str(len(members)), ha="center", va="center",
                     color="white", fontweight="semibold")
        at += width
classes.set_yticks(range(len(groups)), [f'{g["group"]}  (order {len(g["product"])})' for g in groups])
classes.invert_yaxis()
classes.set_xlim(0, 1)
classes.set_xlabel("Relative class size")
classes.set_title("Conjugacy classes · labels show element counts", loc="left")

s4 = groups[-1]
frequencies = Counter(Fraction(s["expected_density"]) for s in s4["stable_subsets"])
subsets.bar(list(map(float, frequencies)), list(frequencies.values()), width=0.027, color="#315a84")
subsets.set_xticks([0, 0.25, 0.5, 0.75, 1], ["0", "1/4", "1/2", "3/4", "1"])
subsets.set_yticks(range(1, max(frequencies.values()) + 1))
subsets.set_xlabel("Predicted density of a class union")
subsets.set_ylabel("Number of subsets")
subsets.set_title("S4 · 32 stable subsets, including the empty and full sets", loc="left")
subsets.grid(axis="y", alpha=0.25)
figure.savefig(directory / "class-density.png", dpi=160)
