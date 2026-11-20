#!/usr/bin/env bash
# Usage: bash scripts/b3-lean.sh ARCHIVE NEW_SCRATCH_DIR conditional|degree-one|original [B2_PROJECT]
set -euo pipefail
scripts=$(cd -- "$(dirname -- "$0")" && pwd)
archive=$(realpath "$1")
mode=$3
reuse=${4:-}
if [[ -n "$reuse" ]]; then reuse=$(realpath "$reuse"); fi
mkdir "$2"
cd "$2"
exec > >(tee verification.log) 2>&1
trap 'rc=$?; printf "exit %d\n" "$rc" | tee exit-code.txt' EXIT
set -x
mathlib=0df444a360eaa60ab8c11dca51a86af692955474
printf '%s\n' leanprover/lean4:v4.33.1 > lean-toolchain
cat > lakefile.toml <<EOF
name = "b3"
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
name = "GroundTheory"
[[lean_lib]]
name = "Development"
[[lean_lib]]
name = "B3"
EOF
lean --version | tee lean-version.txt
grep -F 'version 4.33.1' lean-version.txt
if [[ -n "$reuse" ]]; then
    mkdir .lake
    ln -s "$reuse/.lake/packages" .lake/packages
    python3 - "$reuse/lake-manifest.json" <<'PY'
import json, pathlib, sys
manifest = json.loads(pathlib.Path(sys.argv[1]).read_text())
manifest["name"] = "b3"
pathlib.Path("lake-manifest.json").write_text(json.dumps(manifest))
PY
else
    lake update
    lake exe cache get Mathlib
fi
test "$(git -C .lake/packages/mathlib rev-parse HEAD)" = "$mathlib"
git -C .lake/packages/mathlib rev-parse HEAD
case "$mode" in
conditional)
    cp "$scripts/B3.lean" .
    lake build B3
    lake env lean B3.lean 2>&1 | tee checked.txt
    lake env leanchecker -v B3 | tee kernel-replay.txt
    ;;
degree-one)
    target=GroundTheory.GlobalFields.TateThesis.DegreeOneDensity
    python3 "$scripts/b3-source.py" "$archive" imports.json --extract . --extract-target "$target"
    lake build "$target"
    cat > Check.lean <<'EOF'
import GroundTheory.GlobalFields.TateThesis.DegreeOneDensity
open GroundTheory.Global
#check hasDirichletDensity_inter_degreeOnePrimes_iff
#check tendsto_primeIdealZetaSum_univ_div_univ
#print axioms hasDirichletDensity_inter_degreeOnePrimes_iff
#print axioms tendsto_primeIdealZetaSum_univ_div_univ
EOF
    lake env lean Check.lean | tee checked.txt
    lake env leanchecker -v "$target" | tee kernel-replay.txt
    ;;
original)
    python3 "$scripts/b3-source.py" "$archive" imports.json --extract .
    printf 'Pinned authored closure summary:\n'
    python3 -c 'import json; print(json.load(open("imports.json"))["closure_sizes"])'
    timeout "${B3_BUILD_SECONDS:-240}" lake build GroundTheory.GlobalFields.TateThesis.Chebotarev
    cat > Check.lean <<'EOF'
import GroundTheory.GlobalFields.TateThesis.Chebotarev
open GroundTheory.Global
#check GroundTheory.Global.HasDirichletDensity.union
#check hasDirichletDensity_deuringSet
#check hasDirichletDensity_chebotarevSet_isConj
#check hasDirichletDensity_chebotarevSet
#print axioms GroundTheory.Global.HasDirichletDensity.union
#print axioms hasDirichletDensity_deuringSet
#print axioms hasDirichletDensity_chebotarevSet_isConj
#print axioms hasDirichletDensity_chebotarevSet
EOF
    lake env lean Check.lean | tee checked.txt
    lake env leanchecker -v GroundTheory.GlobalFields.TateThesis.Chebotarev | tee kernel-replay.txt
    ;;
*) exit 2 ;;
esac
python3 - <<'PY'
import pathlib, re
text = pathlib.Path("checked.txt").read_text()
reports = re.findall(r"'[^']+' depends on axioms: \[([^\]]*)\]", text)
if not reports or any(
    set(report.replace(",", " ").split()) - {"propext", "Classical.choice", "Quot.sound"}
    for report in reports
):
    raise SystemExit("Missing axiom reports or unexpected axioms")
PY
