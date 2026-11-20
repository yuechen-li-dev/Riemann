# RIEMANN-M26 — Release the Forced pi/2 Contact

## Outcome

M26 reaches **Success A** for its bounded question. Releasing only the two
algebraic equations that forced the original `pi/2` double contact strictly
improves the legal two-atom dual objective. The first atomic radius remains
fixed at `1`; no third atom, density, or second free radius is introduced.

The certified rational interior witness is

```text
r       = 350785/100000 = 3.50785
base m  = 699/100000    = 0.00699
epsilon = 9/100000      = 0.00009
w2      = m+epsilon     = 0.00708
c       in [1.1567274714890000,1.1567274714900000]
```

Its lower objective endpoint strictly exceeds M25's complete
contact-preserving optimum upper endpoint `1.1566433083347464`. Therefore the
M18 contact was useful scaffolding but had become an artificial active
constraint.

The corrected projection is

```text
J = 2-c,
J_* <= 0.8432725285110000.
```

The stale M25 value `-0.1566431645986454`, which came from subtracting from
`1`, has been removed from the report, README, and contextual helper text. A
semantic regression now exercises the authoritative `2-c` projection.

## Minimal unrestricted family

Repository normalization agrees with the milestone formula:

```text
F(c,w1,w2,r;t)
 = 1 + 2(cos(t)-1)/t^2 + 2(1-c)sin(t)/t
   + 2w1 cos(t) + 2w2 cos(rt),
```

with continuous value `2(1+w1+w2-c)` at the origin. Legal parameters retain
`c,w1,w2>=0` and `r>1`; certification supplies whole-line nonnegativity.

For `A=cos(pi r/2)` and `S=sin(pi r/2)`, substitution of

```text
c  = c0+(pi/2)mA
w1 = w0+m((2/pi)A-rS)
w2 = m
```

collects mechanically to

```text
F = F0 + m[-pi A sin(t)/t + 2((2/pi)A-rS)cos(t) + 2cos(rt)]
  = F0 + m H_r.
```

Oct and Go regress this identity at the origin and representative radii and
frequencies. Thus M22-M25 is an exact lower-dimensional subfamily, not a
renormalized lookalike.

## First variation and active set

At the M25 optimizer, the release direction

```text
(dc,dw1,dw2,dr) = (1,0,1,0)
D(t) = -2 sin(t)/t + 2 cos(rt)
```

has positive value at both old active contacts: Oct finds `D(pi/2)>0.1` and
`D(t1)>0.8`. Hence a sufficiently small positive step increases `c` while
lifting both quadratic minima.

For deterministic seed parameters `r=3.50785`, `m=0.00699`, Oct continues the
ray

```text
c  = c_CP+epsilon
w1 = w1_CP
w2 = m+epsilon
```

to a sampled cap

```text
epsilon_cap approximately 0.0000946358742382
t_new       approximately 1.523146843527.
```

The old `pi/2` and `t1` contacts are inactive; a new double contact controls
the ray boundary. These boundary decimals are active-set discovery, not
certified theorem state. The Go witness deliberately uses the rational
interior step `epsilon=0.00009`, for which the spectrum is strictly positive
on the whole line.

## Oct-first evidence and Go replay

`experiments/m26_release_forced_contact.octest` owns:

- the unrestricted/contact-preserving embedding audit;
- the first-variation experiment;
- release-ray minimization and replacement-contact discovery;
- a deterministic `100001`-point spectrum scan on `[0,40]`;
- seed publication to `artifacts/m26/m26_release_seed.txt`.

The final suite passes `4/4` interpreted in `5443 ms` and `4/4` compiled in
`751 ms`, with zero fallback. The compiler identity is `oct dev`, execution
identity is `gooct-cli`, and execution used CPU `Float`. The bounded experiment
uses no external data, GPU, or hosted execution ID. The seed SHA-256 is
`3f172eef56eb4ba26937ccd6336530f6f33be3229d8c9b4b7c7b964e8d08f48a`.

Go rationalizes the Oct seed and replays it through the authoritative certified
semantics. This is a staged discovery/certification path, not a claim of two
independent implementations. Go proves:

- parameter legality;
- positive values at the released `pi/2` and old movable-contact windows;
- an analytic origin lower bound on `[0,1/2]`;
- centered degree-80 rational Taylor cells on `[1/2,4]`, refined around the
  newly discovered danger band;
- the analytic atomic tail on `[4,infinity)`, with persisted lower bound above
  `0.274`.

The explicit replay completed in `288.06 s`:

```powershell
go run ./cmd/m26replay
# or
go test -tags=slow ./semantic -run M26Slow -count=1
```

Ordinary tests validate the persisted manifest and do not repeat the
multi-minute rational cover. Oct's generated seed is replay input, not an
independent certificate.

## Theorem boundary

M26 proves existence of a legal unrestricted fixed-first-radius two-atom
witness strictly better than the globally optimized contact-preserving
component. It does **not** certify the global optimum of the unrestricted
family. It also does not free the first radius.

The next coherent question is to determine and certify the controlling active
set of this unrestricted fixed-first-radius family. Freeing the first radius
belongs after that boundary, unless the active-set mathematics proves the two
operations cannot be separated.
