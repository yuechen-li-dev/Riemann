"""Authenticate and inventory the pinned Chebotarev import closure."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import tarfile

REVISION = "a172c464fd4bcfc4f53e394dcba45d742b1409a6"
DIGEST = "2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75"
PREFIX = f"riemann-hypothesis-{REVISION}/src/"
BASE = "GroundTheory.GlobalFields.TateThesis."
ROOT = BASE + "Chebotarev"
IMPORT = re.compile(rb"^(?:(?:public|private|meta|all) )*import[ \t]+([^ \t\r\n]+)", re.M)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("archive", type=Path)
    parser.add_argument("manifest", type=Path)
    parser.add_argument("--extract", type=Path)
    parser.add_argument("--extract-target", default=ROOT)
    args = parser.parse_args()
    with args.archive.open("rb") as f:
        if hashlib.file_digest(f, "sha256").hexdigest() != DIGEST:
            raise SystemExit("Pinned archive hash mismatch")
    graph = {}
    with tarfile.open(args.archive, "r|gz") as tar:
        for member in tar:
            if not member.isfile() or not member.name.startswith(PREFIX) or not member.name.endswith(".lean"):
                continue
            module = member.name[len(PREFIX):-5].replace("/", ".")
            source = tar.extractfile(member)
            graph[module] = [m.decode() for m in IMPORT.findall(source.read(16384))]

    def closure(root):
        seen, external, pending = set(), set(), [root]
        while pending:
            module = pending.pop()
            if module in seen:
                continue
            if module not in graph:
                external.add(module)
                continue
            seen.add(module)
            pending.extend(graph[module])
        return sorted(seen), sorted(external)

    modules, external = closure(ROOT)
    if external != ["Mathlib"] or not 1 <= len(modules) <= 400:
        raise SystemExit(f"Unexpected closure: {len(modules)} modules, external={external}")
    if args.extract_target not in modules:
        raise SystemExit("Extraction target is outside the pinned closure")
    extracted = set(closure(args.extract_target)[0])
    wanted = {PREFIX + m.replace(".", "/") + ".lean": m for m in modules}
    records = {}
    declarations = {}
    with tarfile.open(args.archive, "r|gz") as tar:
        for member in tar:
            module = wanted.get(member.name)
            if module is None:
                continue
            data = tar.extractfile(member).read()
            imports = [m.decode() for m in IMPORT.findall(data)]
            if imports != graph[module]:
                raise SystemExit(f"Header scan omitted imports of {module}")
            path = member.name[len(PREFIX):]
            records[module] = {"path": "src/" + path, "sha256": hashlib.sha256(data).hexdigest(), "imports": imports}
            if module == ROOT:
                names = ("HasDirichletDensity.union", "hasDirichletDensity_deuringSet",
                         "hasDirichletDensity_chebotarevSet_isConj", "hasDirichletDensity_chebotarevSet")
                for name in names:
                    match = re.search(
                        rb"^theorem " + re.escape(name.encode()) + rb"(?=\s)[\s\S]+?(?= := by)", data, re.M)
                    if match is None:
                        raise SystemExit(f"Declaration header missing: {name}")
                    declarations["GroundTheory.Global." + name] = {
                        "line": data[:match.start()].count(b"\n") + 1,
                        "header": match[0].decode(),
                    }
            if args.extract is not None and module in extracted:
                target = args.extract / path
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(data)
    if len(records) != len(wanted):
        raise SystemExit("Incomplete closure extraction")
    roots = [ROOT, BASE + "ChebotarevAbelian", BASE + "DegreeOneDensity", BASE + "ChebotarevCount"]
    args.manifest.write_text(json.dumps({
        "revision": REVISION, "archive_sha256": DIGEST,
        "external_imports": external,
        "closure_sizes": {m: len(closure(m)[0]) for m in roots},
        "modules": records,
        "declarations": declarations,
    }, indent=2, sort_keys=True) + "\n", encoding="utf-8")


if __name__ == "__main__":
    main()
