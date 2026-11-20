#!/usr/bin/env bash
# Run with the pinned official Lean bin directory on PATH. The second argument
# is a new scratch project directory, writable and outside the tracked files.
set -euo pipefail
archive=$(realpath "$1")
mkdir "$2"
cd "$2"
exec > >(tee verification.log) 2>&1
trap 'rc=$?; printf "exit %d\n" "$rc" | tee exit-code.txt' EXIT
set -x
revision=a172c464fd4bcfc4f53e394dcba45d742b1409a6
mathlib=0df444a360eaa60ab8c11dca51a86af692955474
printf '%s  %s\n' 2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75 "$archive" | sha256sum -c -
base=Formalism/IncomingBounds
proof=$base/FlowDegreeBounds/MaxFlowValueAttained.lean
mkdir -p "$base/FlowDegreeBounds"
for file in "$base/IncomingBound.lean" "$proof"; do
    tar -xOzf "$archive" "riemann-hypothesis-$revision/src/$file" > "$file"
done
sha256sum "$base/IncomingBound.lean" "$proof" > sources.sha256
cat sources.sha256
printf '%s\n' \
    "201476c18ccbb46a0022e76b63901c7feea2a565b8f9802910d9c36285791ba6  $base/IncomingBound.lean" \
    "18d881da800499e504baa5eb485c4cc800837fa4f14a87fa217b504ea4a903b9  $proof" | sha256sum -c -
grep '^import ' "$base/IncomingBound.lean" "$proof"
printf '%s\n' leanprover/lean4:v4.33.1 > lean-toolchain
cat > lakefile.toml <<EOF
name = "b2"
moreLeanArgs = ["--tstack=400000"]
[leanOptions]
autoImplicit = false
relaxedAutoImplicit = false
maxSynthPendingDepth = 3
pp.unicode.fun = true
maxHeartbeats = 4000000
[[require]]
name = "mathlib"
git = "https://github.com/leanprover-community/mathlib4"
rev = "$mathlib"
[[lean_lib]]
name = "Formalism"
EOF
lean --version | tee lean-version.txt
grep -F 'version 4.33.1' lean-version.txt
lake --version
lake update
test "$(git -C .lake/packages/mathlib rev-parse HEAD)" = "$mathlib"
git -C .lake/packages/mathlib rev-parse HEAD
lake exe cache get Mathlib
lake build Formalism.IncomingBounds.FlowDegreeBounds.MaxFlowValueAttained
cat > Check.lean <<'EOF'
import Formalism.IncomingBounds.FlowDegreeBounds.MaxFlowValueAttained

open Formalism.IncomingBounds

#check FlowAttained
#check LemmaFa
#print FlowAttained
#print LemmaFa
-- Compare independently written proposition interfaces by definitional equality.
example : FlowAttained =
    (∀ (b n : ℕ) [NeZero n], b ≥ 3 → n ≥ 2 → Even n →
      ∀ cap : ZMod n → ℝ, 0 ≤ cap →
        ∃ flow, IsFlow b n cap flow ∧ flowVal n flow = maxFlow b n cap) := rfl
example : LemmaFa =
    (∀ (b n : ℕ) [NeZero n], b ≥ 3 → n ≥ 2 → Even n →
      ∀ cap : ZMod n → ℝ, 0 ≤ cap →
        let value := fun price : ZMod n → ℝ => ∑ v, cap v * (gam b - price v)
        (∃ price, IsDual b n price ∧ IsDualAll b n price ∧ Fnum b n cap = value price) ∧
        IsGreatest {y | ∃ price, IsDual b n price ∧ y = value price} (Fnum b n cap)) := rfl
example : FlowAttained := maxFlowValueAttained
example : LemmaFa := strongDuality
#check maxFlowValueAttained
#check strongDuality
#print axioms maxFlowValueAttained
#print axioms strongDuality
EOF
lake env lean -DwarningAsError=true Check.lean | tee checked.txt
python3 - <<'PY'
import pathlib, re
text = pathlib.Path("checked.txt").read_text()
allowed = {"propext", "Classical.choice", "Quot.sound"}
for name in ("maxFlowValueAttained", "strongDuality"):
    full = "Formalism.IncomingBounds." + name
    m = re.search(r"'" + re.escape(full) + r"' depends on axioms: \[([^\]]*)\]", text)
    if m is None:
        if f"'{full}' does not depend on any axioms" not in text:
            raise SystemExit(f"Missing axiom report for {full}")
    elif set(m[1].replace("\n", "").split(", ")) - allowed:
        raise SystemExit(f"Unexpected axioms for {full}: {m[1]}")
    print("Axiom policy passed for " + full)
PY
lake env leanchecker -v Formalism.IncomingBounds | tee kernel-replay.txt
