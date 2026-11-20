# A4 · Global gain and capacity audit

**Outcome: INCONCLUSIVE on the global estimate.** The submitted construction gives a signed label at most one row, and each diagram column can use that row at most once. The exact replay covers 72 pairs of diagrams and 126 supported monomial and label-order cases. The two signs can reach the same destination without overrunning its rows.

The proof also needs a single pair of weights whose combined loss over primes and lift stages is small. Averaging the **total loss** can select such a pair. We checked that the prime and stage sums leave a positive power of the size parameter to spare for every sufficiently small positive $\varepsilon$, provided the submitted interpolation rank estimate holds.

That proviso is the stopping point. On a whole box of weights, the claimed average fraction of missing rank is at most $C_m |TT|N_1/x^{13}$. The rank must be obtained by integral vectors that work at the required precision and isolate competing lattice edges. The source supplies a theorem with those quantifiers. A4 neither derived the rank estimate for these lattices nor found a legal instance violating it. Thus the global budget fits conditionally, and A4 found no mathematical defect or independent proof of RH or Dirichlet GRH.

**Recommendation: switch verification strategy** to independent kernel replay of the pinned submission against the intended public statements and permitted axioms. The remaining arithmetic input is a uniform integral rank theorem. Extending the small determinant experiment cannot decide it. A4 performed source authentication and exact arithmetic replay only.

## 1. Independent capacity check

[The target and call graph](a4-global-capacity-target.md) give the general support argument and the quantifiers. For a finite falsification attempt fix $p=5,d=2,N_0=5,N_1=6,a=5,T=\{0\}$. Use the diagrams $(t)$ at letter zero and $(u,v)$ at letter one, $0\le t\le5$, $0\le v\le u$, $u+v\le5$. The rank threshold is met at equality, $2\log_5(5)+2=5-1$. Use the identity row operator and the first output columns of each height. The product of its minors is exactly

$$
 z_{0,0,0}^{\,t}\,z_{1,0,0}^{\,u-v}
 (z_{1,0,0}z_{1,1,1}-z_{1,0,1}z_{1,1,0})^{v}.
$$

The four non-top exponents are $(u-h,h,h,v-h)$ for $0\le h\le v$ and the coefficient is $(-1)^h{v\choose h}$. They give row masses $u,v$ and top mass $t$. Every other coordinate has exponent zero. An odd shift reaches letter one only. The positive label qualifies exactly at $t=5$, the negative exactly at $u+v=5$. The test includes both possible enumerations when both qualify. The dual $(1,0)$ is feasible and optimal since $D=(0,u+v)$ forces both flows to zero. No odd shift from zero reaches a top. With $\gamma=5/24$, the geometric series gives $\mathrm{topTerm}=\gamma t$ if $u=0$ and zero otherwise. The negative capacity is $W=u\,1_{u+v=5}$. Take top radius $e=1/5$ when $t=5$ and $u>0$, zero otherwise, and zero at the other letter.

Writing $P,H$ for masses on the positive and negative labels respectively, the replay compares

$$
 \gamma(u+v)+\mathrm{topTerm}-\gamma W
 \ \le\ \gamma(t+u+v)-t e-P(1-5e)-H/5 .
$$

All comparisons use integers in units of $1/120$. At $(t,u,v)=(5,3,2)$ the two label orders charge $(3,2)$ and $(2,3)$. Both spend five units in total. The signed capacity of three is attained and the top row is unavailable for either charge. At zero size the support is the constant monomial. Legal empty label fibres consume zero. The Oct experiment counts eligible shapes and assignments. Go multiplies the minors, verifies every coefficient against the binomial expression and checks the bounds. These are valid instances of the **assembly proposition**. The identity operator illustrates its determinant input. No global host, generic arithmetic frame or admissible interpolation lattice is inferred from it. **FINITE_VERIFIED**. [Implementation](../audit/claudegrh/capacity.go) · [Replay fixture and assertions](../audit/claudegrh/capacity_test.go) · [Oct experiment](../experiments/a4_capacity.octest)

## 2. Units and losses from the orbit to the prime sum

For a served second place of characteristic $p$, the A3 bound enters as
$$
 \mathrm{lowS}\le\sum_{j=1}^{m}\sum_{r\in\mathrm{reps}(p)}
 1_{p\le X''}\big(F_{\rm Orb}(j,p,r)-\mathrm{errF}(j,p,r)\big)
 -\sum_{j=0}^{m-1}dG(j).
$$
The representative is taken at the second embedding. The formula assumes `L.FDef`, $dG(0)=0$ and $dG(j)\le v_p(D\mathrm{Gam})$. At stage $j=i+1$, writing $d_o$ for orbit length, $T_0$ for its top count, $M$ for the terminal rank and $W$ for the orbit weight sum, the loss is
$$
 \gamma T_0(b+\eta_{sz}+1+3S_c)+c_6 W/p
 +d_o N_i N_j(\log_p(4d_o M a)+1)
 +K_{dl1}d_o M^2\log_p(2a)+K_{dl2}T_0\log_p^2(2a).
$$
This accounts for logarithmic coefficient losses on **every** retained stage and orbit. The radius capacity in §1 already enters the local bound. [Stage sum](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LowTermLeSecondRealPlaceBoundsStatement.lean#L11) · [Five losses and constants](a3-coefficient-bound-target.md)

The family's definitions use $s_w(y)=2\,\ell_w e_w y$, with normalized place weight $\ell_w$ and ramification index $e_w$. They keep $F_o=\sum_{j,w}s_w(F_{jw})$, $E_e(u)=\sum_{j,w}s_w(E_{ujw})$ and $G_a=\sum_w s_w(\sum_jg_{jw}-\tau_w)$ separately, with $g=0$ at dropped primes. The combined signed contribution is therefore $F_o+G_a-E_e(u)$, including the single $\tau_w$ charge at each place. The separate interpolation loss is $IO=\sum_w\ell_we_w\,io_w$. The ties identify these values, ranks, slopes, scale, host $q,d$, stage count and selected weights at the point for **each** translate $u_0$. Field and conductor dependence enters before the scalar loss is bounded. [Family definitions](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/LocalRingCardinality/DataParams₂.lean#L203) · [Ties](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalSchurBounds/TieGlobal.lean#L39)

At the arithmetic comparison the nominal extra gain is
$$
 ASK=\sum_{j,t}\Big(\sum_{p\in P_s}2\log p\, V_{jtp}+dT_{jt}\Big),
 \qquad G_a=2\sum_w\ell'_w\big(\sum_j g_{jw}-\tau_w\big).
$$
Here $V$ denotes the source's von Mangoldt class sum, and the weights in this prime decomposition obey $\sum_{w\mid p}\ell'_w=\log p$. Gains and charges must match the underlying prime on support and vanish off it. Grouping gives exactly $ASK-G_a=\sum_p2\log p[\sum_{j,t}V_{jtp}-\sum_jg_{jp}+\tau_p]+\sum_{j,t}dT_{jt}$. Consequently the global test compares a **deficit** to the allowed loss. Crediting $F_o$ to $G_a$ again would double count. [Prime splitting identity](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/SecondKindBucket.lean#L13) · [Use on the family](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/LocalToGlobal/GlobalTheory/LossRequirementsOfSmallShares.lean#L129)

### Explicit budget

Set $m=|J|$, $T=|TT|\ge1$, $U=mTb$, $s_g=100K x^{4\varepsilon-1}$. Each row below is an **additive** charge bounding $ASK-G_a$. The final line bounds the separate prime-log interpolation loss. The letters denote fields of the **same** parameter record $P$. Natural logarithms are used.

| Charge | Upper bound |
|---|---|
| Floors | $2mCfX''(N^2\log a+2(\log a)^2)$ |
| Coefficient/prime cost | $2 C S_k \log(2a)\,mT X''$ |
| Low and high prime terms | $2mT\log(4)\,p_0Q_1$ and $2mT\log(4)\,X'Q_2$ |
| Orbit term | $2mTC(\log x)a[(L_a/k+1/K_0)(\log X'+\log 4)+2(\log(p_0-1)+1)/(p_0-1)]$ |
| Served primes with missing tops | $4amT x^{-2+\varepsilon}$ |
| Order exclusions | $4mTa C_{ord}J_o[\log(2^{e_2}\max(Y,X''))]^2/2^{e_2}$ |
| Exceptional band | $4mTa I_b(4\rho+1/p_0)\log(2X')$ |
| High range | $2mT\log(2X'')(a\,\mathrm{th}^2+w''\mathrm{th})$ |
| Host discriminant parameter | $2mTa\log(d)/d$ |
| Second-class primes | $2amT x^{-9}$ |
| Interpolation, separate | $amT x^{-9}$ |

Specifically $b=\exp(x^\varepsilon),a=xb,X'=\exp(x^{\varepsilon/3}),b'=a/X',w''=b'/(xn),X''=a/w''$, $f=2n$, $N=N_m$, $S_k=S_c+k+\eta$, $\delta_1=x^{-\varepsilon/8}$, $Q_1=C r_G b'\log x$ and $Q_2=C[b'/x+\eta+S_c+(w''+t)\log(x)/\delta_1]+bo$. Here $p_0$ is rounded up from the family's cutoff, $\mathrm{th}=1/(kx^{13})$, $K_0=k^2x^{78}$, $Y=X'(1+C_r\rho_B)$, $J_o=\lfloor x^{\varepsilon/4}\rfloor$, $I_b=\lfloor x^{\varepsilon/3}\rfloor$, $d$ and $e_2$ come from the host and ranks. This $Y$ is the order-cutoff parameter in $P$ and differs from the rank theorem's $Y=x^{13}$. [Parameter substitution](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/LocalToGlobal/GlobalTheory/SuppSat.lean#L109) · [Derived terms](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/WindowSettingSatisfiable.lean#L61) · [Scales](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/ErrorFormsWeakWindowsFalse.lean#L87) · [Loss components](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/WindowSettingSatisfiable.lean#L132) · [Discriminant version](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/PrimeClassLossHypothesesSatisfiable.lean#L60)

Under `WindowsD`, `LargeX`, the prime-class restrictions, nonempty tops, unit rule and the two weighted prime bounds below, the source obtains
$$
 ASK-G_a\le(2\log(4)C+13s_g)U,\qquad \sum_{p\in\mathrm{Cand}}\log p\,io_p\le s_gU .
$$
All eleven deficit terms precede this simplification. The family supplies the common tail by taking the maximum of the window, top, bad-prime, gain, identification and joint-bound thresholds. There $0\le C\le K$ and $s_g\le1$, so $C_{loss}=3|K|+13$ works with $e_L=0$. An inequality on abstract admissible scalar records alone would leave open whether the actual family gives such a record. The source constructs $P$ from the actual $(F_m,x,i)$ and applies the prime-class theorem to that $P$. [Precise window hypotheses](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/PrimeClassLossHypothesesSatisfiable.lean#L12) · [Joint comparison](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/ChainSlopeWindows/DiscriminantSettingSatisfied.lean#L186) · [Actual record and threshold](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/LocalToGlobal/GlobalTheory/LossRequirementsOfSmallShares.lean#L86)

## 3. Does one pair pay all primes and all stages?

There are two distinct prime controls. The first is $\sum_{p\in PR}(\log p/p)\#\mathrm{badTops}_p\le T x^{-2+\varepsilon}$. The design supplies a bound with right side $4C_b x^{-2}\delta_1 n\sum_{p\in PR}\log p/p$, subsequently coarsened on the family tail. The second, for the chosen weights, is
$$
 \sum_{p\in\mathrm{Cand}}\frac{\log p}{p}
 [2\,1_{\Sigma_p>x^{-\varepsilon}}+4\Sigma_p]\le x^{-9},
 \qquad \Sigma_p=m\sum_{v\in\mathrm{goodReps}_p}\kappa_v(a,b).
$$
The stage count is **inside** the comparison with the threshold. Equality to the threshold incurs zero indicator charge. Candidate primes have the required residue class, $p<X'$, order at least $x^{\varepsilon/4}$ and the band exclusion when $p\ge p_0(k)$. The other prime classes have separate charges in the budget. [Bad tops and host characters](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalSchurBounds/Chosen.lean#L31) · [Threshold test](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalSchurBounds/WeightVectorSmallSharesSeven.lean#L17) · [Candidates and stage factor](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/IsResidueClassZero.lean#L279)

An independent calculation checks the passage from the local mean to this simultaneous budget. For $p<p_0(k)$ the source bounds $\kappa_{II}$ **pointwise** by $N_v x^{-12}$. Above the cutoff it needs the local **mean** $\le C_m N_v N_1/x^{13}$ and eventually $C_mN_1\le x$. Restriction to the local top coordinates has equal-sized fibres in the two product boxes, so the local mean is also the mean of the pullback. The representatives partition the top counts, giving $\sum_{v\ \rm good}N_v\le T$. For a uniform pair in the nonempty host boxes, linearity yields $\mathbb E\Sigma_p\le mTx^{-12}$, without independence of the prime losses. For nonnegative $z$, $1_{z>x^{-\varepsilon}}\le x^\varepsilon z$. With $Z=\sum_p(\log p/p)(1_{\Sigma_p>x^{-\varepsilon}}+\Sigma_p)$, the required left side is at most $4Z$. Thus
$$
 \mathbb EZ\le (x^\varepsilon+1)mT x^{-12}\,x^{\varepsilon/3}(2+x^{\varepsilon/3}).
$$
A pair attaining at most the mean exists in these boxes. The host transport described in the target is needed before applying this observation to the final indices. [Local-to-host averaging and prime ranges](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/LogarithmicSums/SlopeEstimateOfLemmaBSlopeLine.lean#L199) · [Sum and nonempty-box selection](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/LogarithmicSums/SlopeEstimateOfLemmaBSlopeLine.lean#L435)

For $x>1$, $m\le x^{1/2+\varepsilon/100}$ and $T\le2x^{2+8\varepsilon+\varepsilon/100}$, bound the first factor by $2x^\varepsilon$ and the last by $3x^{2\varepsilon/3}$. The result is $12x^{-19/2+(1453/150)\varepsilon}$. Hence $4\mathbb EZ\le x^{-9}$ whenever $x^{1/2-(1453/150)\varepsilon}\ge48$. Even at the upper endpoint $\varepsilon=1/52$ the remaining exponent is $2447/7800>0$. Every smaller positive $\varepsilon$ has room for these stage and prime factors. Go also saturates both size bounds at $x=2^{300q}$, $\varepsilon=1/q$ for $q=53,100,1000$. Clearing $x^{-12}$ turns the comparison into an exact integer inequality against $x^3/4$. These are finite instances of the scalar budget. Additional onset conditions, including the bounds on $r_G$, $C_mN_1$ and the box sizes, can depend on $\varepsilon$. This is a **CONDITIONAL** quantitative check. [Exact scalar inequality](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/LogarithmicSums/MeanNum.lean#L124) · [Simultaneous numerical thresholds](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/LogarithmicSums/MeanNum.lean#L216)

## 4. What it would give at the analytic endpoint

With $T_2$ now denoting the second top-set **cardinality**, the tied estimate is $\mathrm{topExcess}\le x^{(e_D+e_L+3)\varepsilon} B$ for each covered translate in the level argument, where
$$
 B=T a(1+m^2/k)+m\left(2nm b x^\varepsilon+
 \frac{T_2 k b x^\varepsilon}{n(n x^{-\varepsilon})}\right).
$$
Here $\mathrm{topExcess}$ itself is the uniform average over host translates of the sum over tops and stages of $ASK-\Re(MAIN)$. The growth step also uses the host identities, thin-scale and character bounds, and the size windows. Examples of joint restrictions are $xm\le k\le n^2x^{1/2-\varepsilon}$, $m\ge x^{1/2-\varepsilon}$, $a\le2xb$, $T\ge x^{-\varepsilon/8}n/2$ and $T_2/n\ge x^{-\varepsilon}$. The two restrictions on $k$ require $n\ge\sqrt{x}$, compatible with the actual window $x^2\le2^{e_2}\le n$. The cost conditions include $2f(x+1)\theta\le1/x$ and $6f(x+1)(\eta+1)\le b/x$ in the growth data. The actual growth call uses $\theta=\exp(-c_6x^{\varepsilon/8})$ here. The prime-share threshold in §3 is $x^{-\varepsilon}$. Using the latter in the cost window would produce a spurious overrun. The three host character bounds are at modulus parameter $8qd$, including the logarithmic derivative bound. The thin estimate requires $y\ge\exp(x^{\varepsilon/3})/3$. [Excess](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/GrowthBoundGivesRiemannHypothesis.lean#L50) · [Bracket](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/ErrorFormsWeakWindowsFalse.lean#L133) · [Joined hypotheses](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/SChain.lean#L38) · [Actual cost window](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalSchurBounds/LossesSmallShares.lean#L279) · [Host character bounds](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalSchurBounds/Chosen.lean#L38) · [Actual call](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AutomorphicArithmetic/PrimitiveCosets/CoreGm.lean#L48)

The level argument uses $e_D=2,e_L=0$. The bracket costs $3\varepsilon$ beyond $e_D+e_L$, the growth reduction costs another $3\varepsilon$, and $x^\varepsilon\ge1+4o$ absorbs the final fixed-level factor ($o=\Lambda.o$). It obtains $\Lambda.G_m(x)\le x^{1/2+9\varepsilon}$ for every sufficiently large $x$ **on the family's windows**. For each sufficiently small positive $\varepsilon$ the construction chooses a family after the required constant thresholds. For every lower cutoff $y$ there is a $Y\ge y$ with the whole interval $[Y,Y^{\mathrm{la}}]$ inside the window set. This is stronger than success at isolated scales, and does not assert an eventual full ray. The separate analytic `FormBAt` implication uses precisely such windows, choosing $\varepsilon$ after fixing a putative zero with real part above $1/2$. Conditionally the exponent beats $2/3$ when $\varepsilon<1/54$, and can approach $1/2$ as required. The rank bound alone would not establish the analytic inputs or the independent division, slope and character estimates. [Family order](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/ParamsHeightRankConstant.lean#L187) · [Level factor](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AutomorphicArithmetic/PrimitiveCosets/CoreGm.lean#L21) · [Window definition](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/RamificationFrobenius/TotientGaloisClasses/Window.lean#L7) · [Zero-free implication](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/RamificationFrobenius/TotientGaloisClasses/TopWin.lean#L68)

## 5. Attempt to falsify the remaining inequality

Fix a legal large-prime orbit and its weight box. Suppose for a fraction $f$ of weight pairs there exists an allowed field, dominant weight and target lattice in which **every** edge-isolated interpolation system has allowable rank at most $(1-\delta)\dim V$, with $\dim V>0$ and $\delta>0$. Then every admissible defect for each of those pairs is at least $\delta$, and taking the infimum preserves that lower bound. Necessarily
$$
 f\delta\le C_m |TT|N_1/x^{13}.
$$
This gives a precise falsification criterion, including the density and rank loss that have to be exhibited. For a claim existential in the uniform constant, a sequence defeating every such constant (or independent control of the actual constant) is needed. Choosing zero vectors gives a poor interpolation system. It supplies no upper bound on the best rank. In fact the empty system makes $\kappa=1$ admissible, yielding only $0\le\kappa_I\le1$. Even a bad pair at one prime would need its weight-box mass included. The small $p=5$ determinant fixtures fail the large-prime requirement $p\ge k^2x^{26}$ as $x$ grows. Neither attack gives an original counterexample. [Empty system and infimum](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/DefsInst.lean#L38) · [Precision and competing-edge conditions](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/UnramifiedQuadraticSetting.lean#L205) · [Mean and hypotheses](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/SmallOrderPrimes.lean#L181)

The unresolved statement is the rank theorem in the [target](a4-global-capacity-target.md), with its stated shifts, mean and uniform constant. The source supplies it via integral projection, branching and saturation theorems. An independent proof or counterexample would require controlling integral invariant factors and edge isolation for the varying base weights and lattices. The finite tests supply none of those invariant factors. This meets stopping rule 3. There is no demonstrated quantitative gap, exact counterexample or new trusted edge in the compiler graph.

## Replay

From the repository root, with Go and Oct available:

~~~sh
go test ./audit/claudegrh -run A4 -count=1 -v
go run ./cmd/a1audit -replay audit/claudegrh/testdata/averaging.json
go test ./audit/claudegrh -run 'TestSelectionCertificates|A3' -count=1 -v
go test ./compiler -run TestM25ReportedJUsesComponentObjectiveLowerEndpoint -count=1
go test ./...
go test -race ./...
go vet ./...
~~~

Set `A1_LEAN_ARCHIVE` to the downloaded [original archive](https://codeload.github.com/Deicyde/riemann-hypothesis/tar.gz/a172c464fd4bcfc4f53e394dcba45d742b1409a6) and run `go test ./audit/claudegrh -run TestPinnedSourceUses -count=1 -v`. The expected SHA-256 is `2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75`. This checks the archive and proof-body call sites. For each of the four A1–A4 files in `experiments/` run `oct test <absolute-path-to-file> --execution interpreted` and repeat with `compiled`. All four passed in both modes with zero compiled fallback. A1 and the two A2 negative fixtures retained **REDUCED_MODEL_COUNTEREXAMPLE**. A3's five exact jet instances retained **FINITE_VERIFIED**. The full Go suite, full race suite, vet, archive replay and focused M25 regression passed. The M25 objective regression checks $J_*=2-c_*$.
