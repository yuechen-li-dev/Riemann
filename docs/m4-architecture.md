# RIEMANN-M4 — Weil criterion and test-function IR

## Result and trust boundary

M4 crosses one deliberate IR boundary:

```text
Semantic / zero-location IR
    RH: ForAll rho in NontrivialZeros(zeta), Re(rho)=1/2
                         |
                         | trusted equivalence theorem
                         v
Functional / test-function IR
    ForAll f in WeilNice, Q_W(f) >= 0
```

The compiler represents this equivalence; it does not certify either open
statement. The criterion and explicit formula are trusted imported mathematics.
Function-class restriction, admissibility checking, transform matching,
provenance retention, and information-loss diagnostics are compiler derived.
The Oct result below is experimental evidence only.

## Selected criterion and normalization

The selected standard source is Jeffrey C. Lagarias, *The Riemann Hypothesis:
Arithmetic and Geometry*, §3, Theorems 3.1 and 3.2. The original source is
André Weil, *Sur les “formules explicites” de la théorie des nombres premiers*,
Communications du Séminaire Mathématique de l'Université de Lund, Tome
Supplémentaire (1952), 252–265.

Lagarias uses functions on the multiplicative group `(0,infinity)` and

```text
M[f](s) = integral_0^infinity f(x) x^s dx/x.
(f*g)(x) = integral_0^infinity f(x/y) g(y) dy/y.
tilde(f)(x) = x^-1 f(x^-1),
M[tilde(f)](s) = M[f](1-s).
```

“Nice” means complex-valued, piecewise `C^2`, compactly supported, and assigned
the average of the left and right limits at every discontinuity. This is the
typed `WeilNice` class. Its regularity supplies the Mellin holomorphy and
vertical decay used by the explicit formula; M4 records those conditions on
the zero aggregate.

For `h = f * tilde(conjugate(f))`, define

```text
Q_W(f) = W^(1)(h)
       = sum_rho M[f](rho) conjugate(M[f](1-conjugate(rho))),
```

where `rho` ranges over the zeros of completed `xi`, equivalently the
nontrivial zeros of `zeta`, counted with multiplicity and in the symmetric
limiting order supplied by the explicit formula. Lagarias Theorem 3.2 states:

```text
RH iff Q_W(f) >= 0 for every nice test function f.
```

On RH, `1-conjugate(rho)=rho`, so each summand is `|M[f](rho)|^2`. M4 imports
this equivalence as `weil-lagarias.positivity-criterion`; it does not prove it.

## Explicit formula and decomposition

Lagarias Theorem 3.1 is imported separately. For a nice `g`,

```text
Wspec(g) = M[g](1) - sum_rho M[g](rho) + M[g](0),

Warith(g) = W_infinity(g) + sum_(p prime) W_p(g),

W_p(g) = log(p) sum_(n>=1) (g(p^n) + tilde(g)(p^n)),

W_infinity(g)
  = (EulerGamma + log(pi)) g(1)
    + integral_1^infinity
        [g(x)+tilde(g)(x)-2g(1)/x] x dx/(x^2-1).
```

The imported identity is `Wspec(g)=Warith(g)`. Consequently,

```text
Q_W(f) = M[h](1) + M[h](0)
         - sum_p log(p) sum_(n>=1)(h(p^n)+tilde(h)(p^n))
         - W_infinity(h).
```

The IR labels the zero aggregate as `zero_side` and endpoint, prime-power, and
archimedean terms as `explicit_formula_side`. They are theorem-linked
representations, not terms silently added together. Prime powers and their
`log(p)`/von-Mangoldt weight are distinct from the real place.

## New semantic objects

- `TransformConvention` records kernel `x^s`, measure `dx/x`, variables, and
  normalization.
- `TestFunction` records symbol, constructor, declared class, attributes,
  transform, and parameters.
- `FunctionClass` distinguishes `WeilNice` from an enumerated finite family and
  has canonical, order-independent identity.
- `UniversalFunctionalStatement` carries quantifier, class, functional,
  predicate, and transform.
- `QuadraticFunctional` and `FunctionalContribution` expose both formula sides.
- `Aggregate` retains index domain, summand, convergence, transform, theorem
  lineage, and source.
- `TestFunctionAdmissibility` is a proposition; prose cannot admit a member.
- `ExplicitFormulaIdentity` separates definition from representation theorem.

No expression parser, integration engine, or general infinite-sum algebra was
introduced.

## Restriction and evidence boundaries

`FunctionClassRestriction` is the M4 binder mechanism. Given `ForAll f in A:
P(f)` and finite `B`, it checks every member has a certified
`TestFunctionAdmissibility(f,A)` claim and transform conventions agree. Only
then does it derive `ForAll f in B: P(f)`.

The directional edge carries:

```text
function_space_restriction:
  source covers every admissible test function;
  conclusion covers only a strict finite family
```

The graph propagates this loss across later IR changes. Reversing the edge, or
using the finite claim to discharge RH through the Weil equivalence, is
rejected. Approximate evidence receives both an exactness diagnostic and a
function-class coverage diagnostic.

## Oct sign-convention experiment

Path: `experiments/m4_weil_toy.octest`

Purpose: check the simplest real restriction of a zero-side reflection-pair
block. A critical-line fixed point contributes `q(a)=a^2`; a formal off-line
pair contributes `q(a,b)=2ab`, an indefinite `2x2` block.

Command (from the local Oct repository):

```text
go run ./cmd/oct test C:\Users\yuech\source\repos\Riemann\experiments\m4_weil_toy.octest --json
```

Exact setup and result:

```text
q_fixed(0)=0; q_fixed(-2)=4
q_pair(1,1)=2; q_pair(1,-1)=-2
Oct version: dev; execution identity: gooct-cli; timing: 604 ms
3 passed, 0 failed, 0 skipped
compiled cases: 3; interpreted fallbacks: 0; diagnostics: none
```

This tests sign/reflection conventions only. It does not show prescribed
transform values come from an admissible function, evaluate the zeta
prime/archimedean formula, or say anything universal. It is not evidence for
or against RH.

## Architectural answers

1. The old geometric `Domain` should not stretch over function spaces. The
   quantifier generalizes; each domain family needs typed identity and subset
   evidence.
2. M3's quantified transport becomes a reusable small pattern: element
   certificates plus a typed subset relation justify universal restriction.
3. A provenance-bearing `Aggregate` is required, but not a general evaluator.
4. The Weil functional needs no CAS when input construction, transform,
   summand, and named components are fixed structurally.
5. Zero/prime duality is a real IR boundary: theorem-linked representations of
   one functional are more honest than a flattened formula.
6. Proposition-switch duplication is now visible in validation, matching,
   cloning, and JSON projection. Redesign during M4 would exceed scope.
7. Finite restriction exposes a Hermitian block in the Oct toy, but a genuine
   matrix lowering needs polarization and evaluable entries; M4 adds no matrix
   IR.
8. Oct adds value for bounded sign, convention, and future finite-matrix work,
   never authority for universal claims.

## Compiler-theory lesson

Lowering between mathematical languages is not a text rewrite. A sound
contract carries the target binder, domain identity, operator identity,
transform convention, aggregate lineage, and theorem trust. Coverage is a
resource: restricting a function space spends coverage just as weakening a
quantifier spends information. Evidence and exactness are orthogonal, so a
typed finite experiment cannot cross the universal theorem boundary.

## Smallest sensible next milestone

Add a finite-basis bilinear/Hermitian-form lowering for a certified admissible
family, with polarization, typed coefficients, and explicit-formula component
entries, while preserving finite-coverage loss. It should not include rank,
inertia inequalities, pair correlation, or percentage-on-line arguments.
