"""Plot Go-checked certificates. Fractions are converted only for display."""
import json
from fractions import Fraction
from pathlib import Path

import matplotlib
matplotlib.use("Agg")
from matplotlib import pyplot as plt

directory = Path("artifacts/b2")
records = json.loads((directory / "certificates.json").read_text(encoding="utf-8"))
sweep = sorted(
    (int(c["capacity"][0]), c)
    for c in records["certificates"]
    if c["l"] == 3 and c["d"] == 2 and c["capacity"][1:] == ["1"]
)
a = [budget for budget, _ in sweep]
fig, (primal, dual) = plt.subplots(2, 1, sharex=True, figsize=(8.3, 6), layout="constrained")
primal.plot(a, [sum(map(Fraction, c["primal"])) for _, c in sweep],
            "o-", color="#163c62", label="Total optimal throughput")
primal.set_ylabel("Primal throughput")
colors = ["#006e64", "#ad4e17"]
for vertex, color in enumerate(colors):
    volume = [Fraction(c["primal"][vertex]) for _, c in sweep]
    primal.plot(a, volume, linestyle=":", marker=".", color=color, label=fr"Flow $x_{{{vertex}}}$")
    # Mark the selected optimal price at the sampled budgets. It need not
    # vary continuously, and an optimizer need not be unique at a breakpoint.
    dual.plot(a, [Fraction(c["dual"][vertex]) for _, c in sweep],
              linestyle="none", marker=("s", "d")[vertex], color=color, label=fr"$\sigma_{vertex}$")
primal.legend(loc="upper left")
dual.legend(loc="center right")
dual.set_ylabel("Selected dual prices")
dual.set_xlabel(r"Varied capacity $D_0$")
dual.set_xticks(a)
for ax in (primal, dual):
    ax.grid(axis="y", alpha=0.25)
fig.suptitle(r"Cyclic LP at $l=3$, $d=2$, $D=(D_0,1)$")
fig.savefig(directory / "capacity-sweep.png", dpi=160)
