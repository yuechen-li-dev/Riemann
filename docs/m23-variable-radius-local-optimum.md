# M23: certified variable-radius two-atom local optimum

## Outcome

M23 achieves **Success A** without adding support or changing the EP3.1 route.
It certifies a unique joint root of the three active equations and proves that
root is a strict local maximum of the contact-preserving two-atom envelope.
Separately, the explicit whole-line witness `r=7/2, m=69/10000` proves
`c_* >= 1+pi/4-2/pi+69*pi*sqrt(2)/40000` and the complementary bound on `J_*`.
The optimizer itself is not promoted to a global family or full-class theorem.

## Active system and two derivations

Put `A(r)=cos(pi*r/2)`. The system is

```text
Q1 = F0(t)+m H_r(t),
Q2 = F0'(t)+m partial_t H_r(t),
Q3 = A'(r)H_r(t)-A(r)partial_r H_r(t).
```

Along an active double-contact branch, `F=0` and `F_t=0` give
`m'=-m partial_r(H)/H`, hence
`C'=(pi*m/2)*(A'H-A H_r)/H`. Independently, Lagrange stationarity of `c`
subject to `F=F_t=0` has zero `F_t` multiplier when `F_tt!=0`; eliminating the
other multiplier gives the same `Q3`. Thus `partial_r F=0` alone would be wrong.

## Certified finite object

The initial box was `m∈[0.0069,0.0071]`, `r∈[3.50,3.52]`,
`t∈[1.78,1.81]`. A purpose-specific rational Krawczyk verifier proves exactly
one joint root in

```text
m ∈ [0.0069950800000000, 0.0069951000000000]
r ∈ [3.5078500000000000, 3.5078570000000000]
t ∈ [1.7953100000000000, 1.7953170000000000]
c ∈ [1.1566431645986454, 1.1566433083347464].
```

The interval determinant is `[0.0050240206736656,0.0051965704520062]`.
Certified curvatures and branch data are

```text
F_tt(pi/2) ∈ [0.0075885043389963, 0.0075912826299478]
F_tt(t1)   ∈ [0.0080108372335414, 0.0080878403414934]
H*F_tt     ∈ [-0.0041655541069319,-0.0041245352180551]
C''        ∈ [-0.0265117460969315,-0.0261034715593046].
```

The implicit-function theorem therefore supplies the smooth local branch and
`C'` changes from positive to negative. Positive `F_tt` also proves the movable
double zero is the only zero in its displayed neighborhood. Box-wide checks
prove parameter legality, structural preservation of `±pi/2`, and
`F(t)>0.2749669735996053` for `|t|>=4`.

The active contacts are fixed quadratic zeros `±pi/2` and movable quadratic
zeros `±t1`. A bounded continuation finds no competing nearby branch. The next
numerical minima near `9.44` and `15.54` remain about `0.571` and `0.597`; these
scans are research evidence.

## Explicit witness and whole-line replay

For `r=7/2, m=69/10000`,

```text
c = 1+pi/4-2/pi+69*pi*sqrt(2)/40000
c ∈ [1.1564423523252728,1.1564423642491166].
```

It is below the local optimizer by about `2.0088e-4` and improves M22 by about
`1.1107e-4`. The proof uses an equality-aware fixed-contact cell `[1.54,1.61]`,
adaptive degree-100 cells with step `1/200000` on `[1.7,1.9]`, remaining exact
compact cells to `4`, and the analytic tail. The exact replay passed in `310.9s`
(M22 reported about `146s`).

## Oct, literature, and architecture

The independent Oct experiment passes `5/5` interpreted Facts in `666ms` and
`5/5` compiled Facts in `612ms`, with zero fallback. Three plots were unchanged
on deterministic replay. Oct is research evidence only.

A post-certification screen of primary/preprint Fourier-optimization literature,
including Chirre--Quesada-Herrera and Das--Ismoilov--Ramos, found related
frameworks but no stated match for this optimizer, radius, contact, or active
system. See [Chirre--Goncalves--de Laat](https://arxiv.org/abs/1810.08843),
[Chirre--Quesada-Herrera](https://arxiv.org/abs/2012.07781), and
[Das--Ismoilov--Ramos](https://arxiv.org/abs/2502.05106). This is not a novelty
or priority claim and requires expert review.

Imported mathematics: the M18 fixed contact, M22 surface/Pell/Taylor/tail route,
and standard Krawczyk/implicit-function theorems. Newly derived: the two-route
Q3 identity, certified joint box, local maximum proof, and stronger witness.
The narrow interval jet beside the Laurent-in-pi IR is the main architectural
awkwardness; generalizing the IR would exceed M23.

Compiler Theory lesson: equality constraints can compile into stable invariants,
leaving the next contact and stationarity relation as a finite proof object,
provided local optimizer evidence stays separate from a globally feasible
witness. RH remains unresolved.

The one proposed next milestone is M24: certify or refute whether this local
two-atom branch controls the envelope over the full legal radius component,
without adding support or changing EP3.1. M23 does not begin it.
