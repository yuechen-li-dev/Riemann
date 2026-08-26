# M21: certified contact-adapted atomic ray

M21 changes only the one-parameter family selected by M20. Starting from the
M18 exact witness

```text
c0 = 1 + pi/4 - 2/pi
w0 = 1/pi - 4/pi^2 + 8/pi^3,
```

it sets `r2=4`, `w2=m`, `c=c0+(pi/2)m`, and
`w1=w0+(2/pi)m`. No generic two-atom search, third atom, exterior density, or
Aetheris integration is introduced.

## Exact response and fixed contact

With the repository's plus-sign Fourier convention,

```text
F_m(t) = F_0(t) + m H(t),
H(t) = -pi sin(t)/t + (4/pi) cos(t) + 2 cos(4t).
```

At `a=pi/2`, `H(a)=-2+2=0` and `H'(a)=4/pi-4/pi=0`.
Therefore `F_m(a)=F_m'(a)=0` for every legal `m`. The actual curvature
identity is

```text
H''(a) = -(30+16/pi^2),
F_m''(a) = F_0''(a) - m(30+16/pi^2).
```

This resolves the sign ambiguity in the request's displayed
`F_0''-mH''`: under the declared decomposition the plus-sign chain rule is
used, and `H''` itself is negative. The certified witness below has
`F_m''(a)>1/200`, so the original contact remains a strict quadratic
tangency.

`H` is even and entire after the removable value at zero is filled in. It is
not periodic because of its sinc term, while
`H(t)=(4/pi)cos(t)+2cos(4t)+O(1/|t|)`. It has a double zero at `pi/2` and
infinitely many asymptotically periodic sign regions. The root inventory in
the report is explicitly numerical evidence, not a theorem input.

## Envelope and second contact

On every region where `H(t)<0`, define `M(t)=F_0(t)/(-H(t))`. Then
`F_m(t)>=0` exactly when `m<=M(t)`, and hence `m_*=inf_{H<0} M(t)` (with the
limiting quotient at the built-in double contact included). Direct
differentiation gives

```text
M'(t) = (-F_0'(t)H(t)+F_0(t)H'(t))/H(t)^2.
```

At `m=M(t)`, this vanishes exactly when `F_m'(t)=0`. The Oct experiment finds

```text
t1  ~= 1.36474957003247
m*  ~= 0.00387226734364153
c   ~= 1.15486093434963
F'' ~= 0.00604771325263.
```

Thus the numerical ceiling has a second strict quadratic contact before the
original contact flattens. This is an active-set/equioscillation pattern, but
M21 does not promote it to a general alternation theorem.

## Certified witness and ceiling bracket

The simple rational `m_safe=19/5000` is certified on the whole line. The exact
consequence is

```text
c_* >= 1+pi/4-2/pi+19*pi/10000 > 1+pi/4-2/pi,
J_* <= 1-pi/4+2/pi-19*pi/10000.
```

Directed displays are `1.15474741<c_safe<1.15474742` and
`0.84525258<J_safe_upper<0.84525259`. The compatible primal guard remains
`c_*<1651/1250` and is checked before theorem propagation.

The proof partition is:

- `[151/100,17/10]`: equality-aware `NonnegativeCellWithCertifiedContact`;
  exact degree-100 curvature samples plus a certified third-derivative loss
  prove strict convexity around `pi/2`.
- `[0,13/10]`, `[13/10,151/100]`, and `[17/10,4]`: targeted degree-100
  rational sine/cosine Taylor cells and the exact `|F'|<9/10` bound.
- `[4,infinity)`: the atomic tail bound
  `F>=3/4-(c-1)/2-2w1-2m>3/10`.
- Evenness supplies negative frequencies; endpoints overlap, so there is no
  frequency gap.

An independent directed evaluation proves that `m=3873/1000000` fails at
`t=5459/4000`. Therefore

```text
19/5000 <= m_* < 3873/1000000.
```

This bracket is rigorous but non-tight. The exact ray ceiling and a certified
unique-root interval for the numerical `t1` remain unresolved.

## Oct evidence and artifacts

`experiments/m21_contact_adapted_ray.octest` passes five interpreted Facts and
five compiled Facts with zero fallback. It reconstructs the formulas, solves
the envelope stationarity equation, checks perturbations above and below the
candidate, and scans 200,001 points on `[0,2000]`. Three `[Artifact]` functions
produce same-scale M18 baseline and optimized-ray contact plots plus the
envelope plot under `artifacts/m21/`. These Float results do not enter the
theorem path.

The generic skill mentions older hosted `oct_check`/`oct_build` calls, but the
installed Oct repository documents the current real path as `oct test` and
`oct artifact`; M21 uses that authoritative CLI path.

## Focused literature screen

After certification, a limited primary/preprint search checked
[Chirre--Gonçalves--de Laat](https://arxiv.org/abs/1810.08843),
[Carneiro--Milinovich--Ramos](https://arxiv.org/abs/2310.01913),
[Carneiro--Chandee--Chirre--Milinovich](https://arxiv.org/abs/2108.09258), and
[Das--Ismoilov--Ramos](https://arxiv.org/abs/2502.05106), together with exact
constant and radius-4 searches. No stated match for this ray, witness, or
two-contact constant was found. This is not a novelty or priority claim.

RH remains unresolved.

The one proposed next milestone is M22: add exactly one contact-motion degree
of freedom selected to release the new `t1` tangency while retaining the
two-atom support assumptions. M21 does not begin that milestone.
