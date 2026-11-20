import Mathlib

open NumberField IsDedekindDomain Filter Topology Set
open scoped BigOperators

namespace B3

-- Additivity of the precise ratio in Mathlib's number-field density definition.
theorem densityUnion {K : Type*} [Field K] [NumberField K]
    {A B : Set (HeightOneSpectrum (𝓞 K))} {a b : ℝ}
    (disj : Disjoint A B) (ha : A.HasDirichletDensity a) (hb : B.HasDirichletDensity b) :
    (A ∪ B).HasDirichletDensity (a+b) := by
  unfold NumberField.Set.HasDirichletDensity at *
  refine (ha.add hb).congr' (Eventually.of_forall fun s => ?_)
  let w := fun p : HeightOneSpectrum (𝓞 K) => (p.asIdeal.absNorm : ℝ) ^ (-s)
  have norm : (Set.univ : Set (HeightOneSpectrum (𝓞 K))).primeIdealZetaSum s = ∑' p, w p := by
    change (∑' p : (Set.univ : Set (HeightOneSpectrum (𝓞 K))), w p) = ∑' p, w p
    exact tsum_univ w
  change A.primeIdealZetaSum s / Set.univ.primeIdealZetaSum s +
    B.primeIdealZetaSum s / Set.univ.primeIdealZetaSum s =
      (A ∪ B).primeIdealZetaSum s / Set.univ.primeIdealZetaSum s
  rw [norm]
  by_cases h : Summable w
  · have hs (S : Set (HeightOneSpectrum (𝓞 K))) : Summable (fun p : S => w p) :=
      h.subtype S
    have hsum : (A ∪ B).primeIdealZetaSum s = A.primeIdealZetaSum s+B.primeIdealZetaSum s :=
      Summable.tsum_union_disjoint (f := w) disj (hs A) (hs B)
    rw [hsum, add_div]
  · simp [tsum_eq_zero_of_not_summable h]

section Classes
variable {G : Type*} [Group G] [Finite G]

def elements (S : Set (ConjClasses G)) : Set G := ConjClasses.mk ⁻¹' S

omit [Finite G] in
theorem classCount (q : ConjClasses G) :
    Nat.card (elements ({q} : Set (ConjClasses G))) = Nat.card q.carrier := by
  have eq : elements ({q} : Set (ConjClasses G)) = q.carrier :=
    Set.ext fun _ => ConjClasses.mem_carrier_iff_mk_eq.symm
  rw [eq]

variable {K : Type*} [Field K] [NumberField K]

def primeFiber (tag : HeightOneSpectrum (𝓞 K) → ConjClasses G)
    (good : Set (HeightOneSpectrum (𝓞 K))) (S : Set (ConjClasses G)) :
    Set (HeightOneSpectrum (𝓞 K)) := good ∩ tag ⁻¹' S

theorem densityFiniteClasses (tag : HeightOneSpectrum (𝓞 K) → ConjClasses G)
    (good : Set (HeightOneSpectrum (𝓞 K)))
    (hclass : ∀ q, (primeFiber tag good {q}).HasDirichletDensity
      ((Nat.card q.carrier : ℝ) / Nat.card G))
    (S : Finset (ConjClasses G)) :
    (primeFiber tag good (↑S : Set (ConjClasses G))).HasDirichletDensity
      ((Nat.card (elements (↑S : Set (ConjClasses G))) : ℝ) / Nat.card G) := by
  classical
  induction S using Finset.induction_on with
  | empty => simpa [primeFiber, elements] using
      (NumberField.Set.hasDirichletDensity_empty (K := K))
  | @insert q S hq ih =>
    have hp : primeFiber tag good (↑(insert q S) : Set (ConjClasses G)) =
        primeFiber tag good {q} ∪ primeFiber tag good (↑S : Set (ConjClasses G)) := by
      ext p
      simp [primeFiber, and_or_left]
    have hpd : Disjoint (primeFiber tag good {q}) (primeFiber tag good (↑S : Set (ConjClasses G))) := by
      refine Set.disjoint_left.2 fun p hp₁ hp₂ => hq ?_
      have he : tag p = q := hp₁.2
      simpa [← he] using hp₂.2
    have hg : elements (↑(insert q S) : Set (ConjClasses G)) =
        elements {q} ∪ elements (↑S : Set (ConjClasses G)) := by
      ext g
      simp [elements]
    have hgd : Disjoint (elements ({q} : Set (ConjClasses G))) (elements (↑S : Set (ConjClasses G))) := by
      refine Set.disjoint_left.2 fun g hg₁ hg₂ => hq ?_
      have he : ConjClasses.mk g = q := hg₁
      change ConjClasses.mk g ∈ S at hg₂
      rwa [← he]
    have count : Nat.card (elements (↑(insert q S) : Set (ConjClasses G))) =
        Nat.card q.carrier + Nat.card (elements (↑S : Set (ConjClasses G))) := by
      rw [hg, ← classCount q]
      simp only [Nat.card_coe_set_eq]
      exact Set.ncard_union_eq hgd
    rw [hp, count, Nat.cast_add, add_div]
    exact densityUnion hpd (hclass q) ih

-- The class map is any choice of representatives after passing to conjugacy.
-- Actual number-field arithmetic enters solely via hclass and the choice of tag/good.
theorem densityStable (tag : HeightOneSpectrum (𝓞 K) → ConjClasses G)
    (good : Set (HeightOneSpectrum (𝓞 K)))
    (hclass : ∀ q, (primeFiber tag good {q}).HasDirichletDensity
      ((Nat.card q.carrier : ℝ) / Nat.card G))
    (C : Set G) (stable : ∀ ⦃a b⦄, IsConj a b → a ∈ C → b ∈ C) :
    (primeFiber tag good (ConjClasses.mk '' C)).HasDirichletDensity
      ((Nat.card C : ℝ) / Nat.card G) := by
  classical
  have hf : (ConjClasses.mk '' C).Finite := (Set.toFinite C).image _
  have hC : elements (ConjClasses.mk '' C) = C := by
    ext g
    constructor
    · rintro ⟨b, hb, he⟩
      exact stable (ConjClasses.mk_eq_mk_iff_isConj.1 he) hb
    · intro hg
      exact ⟨g, hg, rfl⟩
  have h := densityFiniteClasses tag good hclass hf.toFinset
  rw [Set.Finite.coe_toFinset, hC] at h
  exact h

variable {W : Type*}

def overPrimeSet (above : HeightOneSpectrum (𝓞 K) → W → Prop)
    (unramified : W → Prop) (frob : W → G) (C : Set G) :
    Set (HeightOneSpectrum (𝓞 K)) :=
  {p | ∀ w, above p w → unramified w ∧ frob w ∈ C}

-- Mirror the corpus's "every place above p" predicate. Its two structural
-- inputs let us pass to a well-defined class tag on base primes.
theorem densityAllAbove (above : HeightOneSpectrum (𝓞 K) → W → Prop)
    (unramified : W → Prop) (frob : W → G)
    (existsAbove : ∀ p, ∃ w, above p w)
    (agrees : ∀ p w z, above p w → above p z → IsConj (frob w) (frob z))
    (hclass : ∀ q : ConjClasses G, (overPrimeSet above unramified frob q.carrier).HasDirichletDensity
      ((Nat.card q.carrier : ℝ) / Nat.card G))
    (C : Set G) (stable : ∀ ⦃a b⦄, IsConj a b → a ∈ C → b ∈ C) :
    (overPrimeSet above unramified frob C).HasDirichletDensity
      ((Nat.card C : ℝ) / Nat.card G) := by
  classical
  let tag (p : HeightOneSpectrum (𝓞 K)) := ConjClasses.mk (frob (Classical.choose (existsAbove p)))
  let good : Set (HeightOneSpectrum (𝓞 K)) := {p | ∀ w, above p w → unramified w}
  have bridge (S : Set G) (hs : ∀ ⦃a b⦄, IsConj a b → a ∈ S → b ∈ S) :
      primeFiber tag good (ConjClasses.mk '' S) = overPrimeSet above unramified frob S := by
    ext p
    have chosen : above p (Classical.choose (existsAbove p)) := Classical.choose_spec (existsAbove p)
    constructor
    · rintro ⟨hu, ⟨a, ha, he⟩⟩ w hw
      exact ⟨hu w hw, hs ((ConjClasses.mk_eq_mk_iff_isConj.1 he).trans
        (agrees p _ w chosen hw)) ha⟩
    · intro hp
      exact ⟨fun w hw => (hp w hw).1, ⟨_, (hp _ chosen).2, rfl⟩⟩
  have htag (q : ConjClasses G) : (primeFiber tag good {q}).HasDirichletDensity
      ((Nat.card q.carrier : ℝ) / Nat.card G) := by
    have saturated : ∀ ⦃a b : G⦄, IsConj a b → a ∈ q.carrier → b ∈ q.carrier := by
      intro a b h ha
      exact ConjClasses.mem_carrier_iff_mk_eq.2
        ((ConjClasses.mk_eq_mk_iff_isConj.2 h).symm.trans (ConjClasses.mem_carrier_iff_mk_eq.1 ha))
    have image : ConjClasses.mk '' q.carrier = {q} := by
      ext z
      constructor
      · rintro ⟨a, ha, rfl⟩
        exact ConjClasses.mem_carrier_iff_mk_eq.1 ha
      · intro hz
        have hzq : z = q := hz
        subst z
        obtain ⟨a, ha⟩ := ConjClasses.mk_surjective q
        exact ⟨a, ConjClasses.mem_carrier_iff_mk_eq.2 ha, ha⟩
    rw [← image, bridge _ saturated]
    exact hclass q
  rw [← bridge C stable]
  exact densityStable tag good htag C stable

theorem classWeight (a : G) :
    (Nat.card (ConjClasses.mk a).carrier : ℝ) / Nat.card G =
      1 / (Nat.card (Subgroup.centralizer ({a} : Set G)) : ℝ) := by
  classical
  let : Fintype G := Fintype.ofFinite _
  let : Fintype (ConjAct G) := Fintype.ofEquiv G ConjAct.toConjAct.toEquiv
  let : Fintype (MulAction.orbit (ConjAct G) a) := Fintype.ofFinite _
  let : Fintype (MulAction.stabilizer (ConjAct G) a) := Fintype.ofFinite _
  have h := MulAction.card_orbit_mul_card_stabilizer_eq_card_group (ConjAct G) a
  simp only [← Nat.card_eq_fintype_card, ConjAct.orbit_eq_carrier_conjClasses,
    ← Subgroup.nat_card_centralizer_nat_card_stabilizer] at h
  have size : Nat.card (ConjAct G) = Nat.card G :=
    Nat.card_congr ConjAct.toConjAct.toEquiv.symm
  rw [size] at h
  have hn : (Nat.card G : ℝ) ≠ 0 := by exact_mod_cast (Nat.card_pos : 0 < Nat.card G).ne'
  have hc : (Nat.card (Subgroup.centralizer ({a} : Set G)) : ℝ) ≠ 0 := by
    exact_mod_cast (Nat.card_pos : 0 < Nat.card (Subgroup.centralizer ({a} : Set G))).ne'
  field_simp
  exact_mod_cast h

theorem densityFromCentralizers (above : HeightOneSpectrum (𝓞 K) → W → Prop)
    (unramified : W → Prop) (frob : W → G)
    (existsAbove : ∀ p, ∃ w, above p w)
    (agrees : ∀ p w z, above p w → above p z → IsConj (frob w) (frob z))
    (single : ∀ a, (overPrimeSet above unramified frob {b | IsConj a b}).HasDirichletDensity
      (1 / Nat.card (Subgroup.centralizer {a})))
    (C : Set G) (stable : ∀ ⦃a b⦄, IsConj a b → a ∈ C → b ∈ C) :
    (overPrimeSet above unramified frob C).HasDirichletDensity
      ((Nat.card C : ℝ) / Nat.card G) := by
  apply densityAllAbove above unramified frob existsAbove agrees _ C stable
  intro q
  obtain ⟨a, rfl⟩ := ConjClasses.mk_surjective q
  have carrier : (ConjClasses.mk a).carrier = {b | IsConj a b} :=
    Set.ext fun b => ConjClasses.mem_carrier_iff_mk_eq.trans
      (eq_comm.trans ConjClasses.mk_eq_mk_iff_isConj)
  rw [classWeight, carrier]
  exact single a

end Classes
end B3

#check B3.densityUnion
#check B3.classCount
#check B3.densityFiniteClasses
#check B3.densityStable
#check B3.densityAllAbove
#check B3.classWeight
#check B3.densityFromCentralizers
#print axioms B3.densityUnion
#print axioms B3.classCount
#print axioms B3.densityFiniteClasses
#print axioms B3.densityStable
#print axioms B3.densityAllAbove
#print axioms B3.classWeight
#print axioms B3.densityFromCentralizers
