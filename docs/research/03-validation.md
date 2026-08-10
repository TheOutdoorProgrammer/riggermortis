# Validation Research

Research notes for **riggermortis**.
Central question: *"Is it possible to validate our rigs are correct, instead of just thinking they are?"*

Short answer, up front, so nobody has to read to the end: **partially, and less than you want.**
A meaningful slice of knot *topology* is machine-provable.
Almost nothing about *rigs* is provable; rigs are checkable only against rules a human wrote.
And essentially nothing about *strength* or *"this knot holds"* is provable by software at all.

The point of this document is section 6, which draws that line explicitly.
Sections 1 to 5 are the evidence for it.

---

## 1. Topological Validation

### 1.1 The framing problem: a fishing knot is an OPEN curve

This is the single most important finding in this document, and it changes what "validation" can even mean.

Classical knot theory (Jones, Alexander, HOMFLY, unknot recognition, the whole toolbox) is defined on **closed loops**, embeddings of S¹ in S³.
A fishing knot is not a closed loop.
It is an open arc with two distinguished endpoints: the standing line running off to the rod, and the tag end.
Frequently it also passes through a hook eye, a swivel, or a ring, which is a second topological object.

Open arcs have **no knot type** in the classical sense.
Any open arc can be continuously deformed to a straight line, which is precisely why a badly tied knot slips.
So you cannot naively "compute the Jones polynomial of a Palomar knot" and get a well-defined answer.
Anyone who tells you otherwise is skipping a step.

There are three legitimate ways to get a well-defined topological answer, and picking one is a design decision riggermortis has to make and then document forever:

1. **Closure schemes.**
   Join the two endpoints (to each other, to a point at infinity, or to a random point on a large enclosing sphere) and compute the classical invariant of the resulting closed loop.
   Because the answer depends on closure direction, serious tools sample many closures and report a **distribution** of knot types ("87% trefoil, 13% unknot"), not a single answer.
   This is what pyknotid's `OpenKnot` and Topoly do.
   Reassuringly, Millett, Dobay and Stasiak (*Macromolecules* 38:601, 2005) found that for open random walks up to several hundred segments, simple **direct end-to-end closure captured the dominant knot type** of the multi-closure distribution in all but a handful of cases, and the exceptions were exactly the configurations the distribution already flagged as low-confidence.
   That is empirical reassurance, not a theorem.
2. **Knotoids** (Turaev 2012; Goundaroulis, Dorier, Stasiak et al.).
   A knotoid is the honest mathematical object for an open curve with two endpoints in a plane or on a sphere.
   Knotoids have their own well-defined invariants: the Turaev loop bracket, the knotoid Jones polynomial, the arrow polynomial.
   This is the mathematically correct object for a fishing knot, and it is what **Knoto-ID** computes.
3. **Fix the boundary.**
   Treat the standing end and tag end as pinned to the boundary of a ball and classify the arc rel endpoints.
   In practice this is equivalent to (2), and it is the framing most physical-knot papers use.

**Recommendation:** if riggermortis wants an invariant that is *stable under editing* (the regression-test use case), use either knotoid invariants (path 2) or a **fixed, documented, deterministic** closure (path 1 with a pinned closure direction).
A randomly sampled closure distribution is a bad regression test.
It is stochastic and it will flake in CI.

### 1.2 Can we detect that an authored knot is actually the unknot?

**Yes for the closed-loop case, and this is the strongest genuinely-provable check available to the project.**

Unknot recognition is a *decided* problem.
It is known to lie in NP and in co-NP, and practical exact algorithms exist.
**Regina** implements Burton and Ozlen's unknot recognition via normal surface theory; it is exact, not heuristic, though worst-case exponential.
At fishing-knot sizes (roughly 3 to 12 crossings) that worst case is irrelevant; these run instantly.

Cheaper, weaker, and still useful in CI:

- **Reidemeister simplification to zero crossings.**
  Spherogram's `simplify(mode='global')` and Regina's exhaustive simplification.
  If a diagram simplifies to 0 crossings, it **is** the unknot, and a positive result here is a proof.
  Failure to simplify proves nothing.
- **Invariant mismatch against the unknot.**
  If the Jones polynomial is not 1, the diagram is provably **not** the unknot.
  This is sound in one direction only: it can prove "knotted", it cannot prove "unknotted", because the Jones unknotting conjecture is still open.

**The honest statement:** we can *prove* a diagram is knotted (nontrivial invariant), and we can *prove* it is unknotted (successful simplification to 0 crossings, or Regina's exact recognizer).
Anything in between is "we failed to decide", which at fishing-knot sizes essentially never happens.

There is a second, cheaper class of check worth naming separately: **diagram well-formedness**.
A PD code or Gauss code either does or does not describe a realizable planar diagram.
Every crossing must have four arc labels, arcs must pair, the underlying 4-valent graph must be planar, orientations must be consistent.
This is fully decidable and the libraries will reject a malformed code outright.
It catches transcription errors, which is probably the single most common authoring failure in a project like this.

### 1.3 Can invariants act as regression tests?

**Yes. This is the highest-value, lowest-risk use of knot theory in this project, and it is the thing I would build first.**

The pattern is a golden-file snapshot test:

```yaml
knot_id: palomar
pd_code: [[1,5,2,4],[3,1,4,6],[5,3,6,2]]   # authored data, hand-maintained
# everything below is GENERATED, committed, and diffed in CI
invariants:
  library: regina
  library_version: "7.3.1"
  closure: end_to_end        # pinned and documented, never random
  crossing_number: 3
  writhe: -3
  jones: "-t^-4 + t^-3 + t^-1"
  homfly: "..."
  alexander: "t - 1 + t^-1"
```

If someone edits the crossing sequence and the Jones polynomial changes, CI fails loudly with **"your edit silently changed which knot this is."**
That is a real, hard, automatic correctness signal, and it is exactly the class of bug that is otherwise invisible in review.

Caveats that must be documented or the test quietly lies:

- **Invariants are not complete.**
  Distinct knots can share a Jones polynomial; the smallest known pairs occur around 10 crossings and infinite families exist.
  Equal invariants means *probably* the same knot, not *certainly*.
  Combining Jones plus HOMFLY plus Alexander plus crossing number makes a collision vanishingly unlikely at fishing-knot sizes, but it is still not a proof.
- **Chirality matters and Alexander cannot see it.**
  The Jones polynomial distinguishes a knot from its mirror image; the Alexander polynomial does not.
  Fishing knots care enormously about which way you wrap.
  Prefer Jones and HOMFLY for regression purposes and treat Alexander as a weak supplement.
- **Normalisation conventions differ between libraries.**
  A HOMFLY from Regina and a HOMFLY from Sage can differ by variable substitution or sign.
  Pin one library, pin its version, store its exact output string, and record both in the snapshot (as above).

### 1.4 Library survey

| Library | Language | License | Capability | CI-viable? |
| --- | --- | --- | --- | --- |
| **pyknotid** | Python, optional Cython | MIT | Space curves and standard diagram notations. Alexander polynomial, Vassiliev degree-2, Gauss codes, self-linking, writhe. **`OpenKnot` class explicitly handles open curves** via sphere closures and closure at infinity, plus virtual-versus-classical projection checks; returns *distributions* (`alexander_polynomials`, `alexander_fractions`). Jones and HOMFLY are **not** documented on `OpenKnot`. | **Yes, with reservations.** `pip install pyknotid`, pure-Python fallback path, no system deps. **But the last PyPI release was April 2018 and third-party trackers classify maintenance as "Inactive."** Treat as unmaintained: vendor it or pin hard. |
| **Regina** (`regina-normal`) | C++ core, first-class Python bindings | GPL-2.0-or-later | Strongest exact toolkit for this job. **Burton and Ozlen unknot recognition**, Reidemeister moves, exhaustive simplification, Jones, HOMFLY-PT, Alexander (knots only), Kauffman bracket, arrow and affine-index polynomials for virtual links, link group and extended link group. Docs warn several polynomials are exponential-time on large links, which is irrelevant at 3 to 12 crossings. Alexander is explicitly noted as polynomial-time but weak. | **Yes, best in class, but not a one-line pip.** Distributed as OS packages (Debian, Ubuntu, Homebrew, conda) and as a SageMath spkg. A **PyPI package prepared by Culler, Dunfield and Goerner exists but is scoped to use inside SageMath** and excludes the GUI and CLI. In CI, use a pinned container or conda, not bare pip. Actively maintained. |
| **SnapPy / Spherogram** | Python wheels over C | Open source, GPL-compatible (verify per file) | Spherogram builds links from PD and DT codes; `simplify()` including a global Reidemeister mode; writhe, determinant, signature, linking number, braid words, alternating and planarity checks, connected sum, mirror. `exterior()` gives the link complement for hyperbolic invariants. | **Yes for the diagram layer.** `pip install snappy` ships prebuilt wheels including Python 3.13. **Critical limitation: the Jones polynomial, Alexander polynomial, Seifert matrix, knot group and knot Floer homology all require running inside SageMath.** The pip-only subset is diagram surgery plus numeric invariants, which guts its value as a polynomial-based regression oracle. |
| **Knoto-ID** | C++ | GPL-2.0-or-later | **Purpose-built for open curves.** Knotoid invariants: Jones polynomial for knotoids on the sphere, Turaev loop bracket in the plane, arrow polynomial, loop arrow polynomial. Consumes PDB, raw xyz coordinates, extended Gauss code and PD code. Produces "fingerprint matrices" summarising entanglement of all subchains. | **Yes, but you build it yourself.** CMake plus Boost (graph, random, regex). No package-manager distribution found. Repository looks dormant (few commits, no recent activity). Best consumed as a pinned Docker layer. **Mathematically the most correct tool for this project's actual object.** |
| **Topoly** | Python over a C++ core | UNCERTAIN | Reported to compute Alexander, Jones, Conway, HOMFLY, Yamada, Kauffman, BLM/Ho and GLN, and to handle **open chains via probabilistic closure schemes**. From the Sułkowska lab (protein topology). | **UNCERTAIN.** `topoly` exists on PyPI but the project page failed to load during this research; wheel availability and Python version support are unverified. Worth 20 minutes of hands-on evaluation because on paper it is the closest single-package fit (open chains plus a full invariant set plus pip). |
| **SageMath knots module** | Python inside Sage | GPL | `sage.knots.link` gives Jones, Alexander, HOMFLY, Seifert matrix, braid words. `sage.knots.knotinfo` gives programmatic access to the KnotInfo and LinkInfo databases, meaning **published invariant values for every named knot**. Superb as a cross-check oracle for "is our Palomar encoding actually knot 3_1?" | **Marginal.** Sage is a multi-gigabyte dependency. Viable only as a separate, occasionally-run Docker job using the `sagemath/sagemath` image, never in the hot CI path. |
| **KnotJob** | Java, GUI-oriented | UNCERTAIN | Khovanov homology, odd Khovanov, knight-move homology, Rasmussen s-invariant. Vastly beyond what this project needs. | **No.** Desktop application, not a library. Skip. |
| **KnotTheory.jl** | Julia | UNCERTAIN | PD data structures, Alexander, Jones, Conway, HOMFLY-PT, Seifert theory, Reidemeister simplification, braid interop. | Possible, but adds a whole Julia toolchain for no capability Regina lacks. |
| **knotkit** | C++ | UNCERTAIN | Khovanov homology and friends. Research-grade. | No. |
| **JavaScript / Rust / Go** | n/a | n/a | **No production-grade knot invariant library exists in any of the three.** Searching turns up only student projects, thesis code, and one-off Jones implementations. | **This is a real architectural constraint.** If riggermortis is a Go or JS codebase, the topology validator must be a separate Python or C++ process, most likely a pinned container invoked from CI, not an in-process library call. Plan for that now rather than discovering it later. |

**If forced to pick:** Regina in a pinned Docker image for the closed-loop layer (unknot recognition plus Jones/HOMFLY regression snapshots), plus Knoto-ID or pyknotid `OpenKnot` for the open-curve layer.
`pip install snappy` is the cheapest thing that does something useful without a container, but its polynomial invariants live behind Sage, which removes most of the regression-test value.

Sources: [pyknotid](https://github.com/SPOCKnots/pyknotid), [pyknotid OpenKnot docs](https://pyknotid.readthedocs.io/en/latest/sources/spacecurves/openknot.html), [pyknotid maintenance status](https://snyk.io/advisor/python/pyknotid), [Regina](https://regina-normal.github.io/), [Regina link analysis](https://regina-normal.github.io/docs/link-analysis.html), [Regina feature set](https://regina-normal.github.io/docs/featureset.html), [regina-python man page](https://regina-normal.github.io/docs/man-regina-python.html), [Regina in Sage](https://doc.sagemath.org/html/en/reference/spkg/regina.html), [Spherogram](https://snappy.computop.org/spherogram.html), [SnapPy install](https://snappy.computop.org/installing.html), [Knoto-ID](https://github.com/sib-swiss/Knoto-ID), [Topoly on PyPI](https://pypi.org/project/topoly/), [KnotTheory.jl](https://github.com/hyperpolymath/KnotTheory.jl), [knotkit](https://github.com/cseed/knotkit), [Burton, HOMFLY-PT is FPT](https://arxiv.org/pdf/1712.05776).

---

## 2. Physics Simulation

### 2.1 The state of the art is real and it is not ours

**Discrete Elastic Rods** (Bergou et al., SIGGRAPH 2008, ACM TOG 27(3):63, doi 10.1145/1360612.1360662; extended to viscous threads in 2010, doi 10.1145/1778765.1778853) is the standard discretisation for simulating thin elastic filaments with bending and twisting.
DER on its own does not handle self-contact, which is the entire physics of a knot.

The contact piece is the **Implicit Contact Model (IMC)** of Choi, Tong, Jawed et al., published as "Implicit Contact Model for Discrete Elastic Rods in Knot Tying" (*J. Appl. Mech.* 88(5):051010, 2021) and extended in "A fully implicit method for robust frictional contact handling in elastic rods" (*Extreme Mechanics Letters*, 2022).
IMC builds a twice-differentiable analytical contact potential so contact and friction forces have both a gradient and a Hessian, which is what makes contact-rich scenarios like knot tightening actually converge.

**Chasing the "successor project" lead:** `QuantuMope/imc-der` is the original IMC knot-tying codebase, and its own README now redirects users to **DisMech**, stating the DisMech codebase contains many improvements over the outdated `imc-der` repository.

- **DisMech** (`StructuresComp/dismech-rods`): C++ with pybind11 Python bindings, **GPL-3.0**, roughly 168 commits, moderately active.
  Simulates rods as chains of cylinders, uses FCL for broadphase and narrowphase collision, implements IMC with per-limb friction coefficients, and supports implicit integration.
  Dependencies are heavy: Eigen 3.4.0, FCL, libccd, SymEngine, Intel oneMKL, OpenGL/GLUT, pybind11.
  **The oneMKL dependency is an x86 concern for arm64 CI and arm64 dev machines. UNCERTAIN whether it is strictly required or swappable; verify before committing to this path.**
- **MAT-DiSMech** (arXiv:2504.17186): the same framework in MATLAB, covering rods, shells and soft robots.
  MATLAB licensing makes it a non-starter for an open-source project's CI.
- **rod-contact-sim** (`StructuresComp/rod-contact-sim`): the sibling DER-plus-IMC codebase, aimed at flagella bundling rather than knots.

### 2.2 Patil, Sandt, Kolle and Dunkel, *Science* 367:71-75 (2020)

This is the paper worth reading, and it is more modest than the headlines suggested.

**What they actually did.**
They combined optomechanical experiments (photonic fibres that change colour under strain, so you can literally see where stress localises inside a tied knot) with simulations and theory.
Exploiting an analogy with long-range ferromagnetic spin systems, they derived **simple topological counting rules** based on **twist charge, crossing number and handedness** to predict mechanical stability.

**What they actually claim.**
Read the abstract closely: *"simple topological counting rules to predict the **relative** mechanical stability of knots and tangles, in agreement with simulations and experiments for **commonly used climbing and sailing bends**."*
The deliverable is a **ranking**, not a number.
It explains why the reef knot holds and the granny slips, in a family of structurally similar bends.
It does not output "this knot retains 87% of line strength."

**Limits that matter to riggermortis.**

- The validation set is climbing and sailing **bends**, meaning two ropes joined end to end, in relatively thick, high-friction, textile rope.
  A fishing knot is usually a *hitch* onto a metal eye, in a slick monofilament or a slippery, near-frictionless braid, at a diameter ratio nothing like a climbing rope.
  The counting rules are a physically motivated heuristic, not a law, and their transfer to that regime is **unvalidated**.
- The rules rank knots against each other within a comparable family.
  They say nothing about absolute load, and nothing about the material-dependent failure that dominates fishing (line cutting itself under the first tight wrap).
- Friction is a first-class input to whether a knot holds, and the coefficient of friction of a given fishing line against itself, wet, under load, is not a published number.

**Code availability: yes.**
Patil released the knot simulation code as `vppatil28/knot_simulation` on Zenodo, doi [10.5281/zenodo.3528928](https://doi.org/10.5281/zenodo.3528928), described as "Code to simulate 1-tangles and 2-tangles together with a selection of initial configurations", licensed "other-open".
So the method is genuinely implementable, and the counting rules in particular are cheap to implement because they are combinatorial, not dynamical.

### 2.3 Verdict

**Academically real. Not useful as a validator for this project, with one exception.**

Full DER-plus-IMC simulation is out for riggermortis:

- It is a GPL-3.0 C++ stack with Eigen, FCL, libccd, SymEngine and oneMKL, which is a serious CI liability, likely arm64-hostile, and slow.
- Its outputs are entirely determined by parameters nobody has measured for fishing line: bending and twisting stiffness, self-friction coefficient, contact stiffness.
  Garbage parameters give confident, wrong answers, which is worse than no answer.
- A simulation that says "this knot slipped" is a statement about the model, not about the knot.
  Presenting it as validation would be exactly the overselling this document exists to prevent.

**The one exception worth building:** the Patil counting rules themselves.
Twist charge, crossing number and handedness are combinatorial quantities computable directly from the crossing sequence riggermortis already stores, with no simulation, no dependencies and no parameters.
That makes them a legitimate **Tier B heuristic**: cheap, published, citable, and honest if labelled as "predicted relative stability rank, per Patil et al. 2020, validated on climbing bends and not on fishing knots."
Shipping that with the caveat attached is defensible.
Shipping it as "verified to hold" is not.

Sources: [imc-der](https://github.com/QuantuMope/imc-der), [DisMech](https://github.com/StructuresComp/dismech-rods), [rod-contact-sim](https://github.com/StructuresComp/rod-contact-sim), [MAT-DiSMech](https://arxiv.org/html/2504.17186v2), [IMC in J. Appl. Mech.](https://asmedigitalcollection.asme.org/appliedmechanics/article/88/5/051010/1099667/Implicit-Contact-Model-for-Discrete-Elastic-Rods), [fully implicit frictional contact](https://www.sciencedirect.com/science/article/abs/pii/S2352431622002000), [Patil et al., Science 2020](https://www.science.org/doi/10.1126/science.aaz0135), [SIAM News summary](https://www.siam.org/publications/siam-news/articles/topological-knot-mechanics/).

---

## 3. Empirical Strength Data

### 3.1 The best-documented source is from climbing, and it is honest about being thin

**Moyer, Tusting and Harmston, "Comparative Testing of High Strength Cord", ITRS 2000** (reprinted in *Nylon Highway* 49, 2004) is the paper a prior run retrieved, and it is a good example of careful work at small scale.

Methodology, as stated by the authors:

- Machine: an 11,000 lb SATEC Apex 11 EMF universal test machine at Black Diamond.
- Sample size: **five samples per material per configuration, results averaged.**
- Conditions: not formally conditioned, but all tests at 29% ±4% relative humidity and 71 °F ±6.
- Pull rates and fixtures "consistent with CEN standards", though the specific rate is not given in the text.
- Knots tested: figure-eight knots, and loops tied with double fisherman's, triple fisherman's and water knots.
- **No stated protocol for dressing or setting the knots.** This is a significant omission, because dressing is a known first-order variable.

The authors' own caveat is the most valuable line in the paper: an average breaking strength is not a good quantity for deciding whether a component is strong enough, and **five samples are not sufficient to determine a meaningful statistical minimum**, so averages are presented instead.
That is exactly the right disclosure, and it is also a warning: the best-documented knot-strength paper in this space explicitly says its own numbers cannot support the engineering claim people want to make with them.

**Transferability to fishing line: poor, and this needs stating plainly.**
Kernmantle climbing rope, accessory cord and webbing are high-friction, multi-strand, large-diameter textiles.
Fishing line is a slick monofilament, a slicker fluorocarbon, or a UHMWPE braid whose whole selling point is low friction and near-zero stretch.
The failure modes differ: climbing knots typically fail by the rope cutting itself at the first tight bend; braid frequently fails by *slipping* instead.
Numbers do not carry across, and neither do knot rankings.

### 3.2 The meta-analysis that should calibrate everyone's expectations

**Thomas Evans, "A Review of Knot Strength Testing" (SAR3, 2016)** is a data-mining review of the entire published knot-testing literature, and it is devastating in the most useful way.

Findings, from the paper itself:

- **114 sources, more than 1,440 individual tests** aggregated.
- **Sample size distribution is catastrophic.** Of the tests reporting a sample size: **n=1 for 383 tests**, n=2 for 56, n=3 for 80, n=4 for 24, n=5 for 74, n=6 for 14, and above that a literal handful (one test at n=8, one at n=10, two at n=11, one at n=12).
  The single most common experimental design in published knot testing is **one pull**.
- **Residual knot strengths cluster between roughly 45% and 85%** of unknotted strength, and **the per-knot ranges overlap each other heavily**.
  Look at the reported ranges: bowline 41.8-70.7%, figure-8-on-a-bight 64.8-86.3%, butterfly loop-to-end 60.7-80.6%.
  Those intervals are wider than the differences between the knots.
- **Wet versus dry shows no consistent pattern** across 37 comparisons; some wet samples were stronger, some weaker.
- **Pull rate does appear to matter** (faster pull, lower strength) but only four comparisons existed, so the relationship cannot be quantified.
- Evans states directly that the studies use different methods and materials and **their results cannot be compared directly without incurring error**, and that readers should take the results "with a grain of salt".
- **He explicitly excluded fishing line from the review.** Articles on knots in materials other than rope, cord or webbing were excluded, naming fishing line specifically.

So the most rigorous meta-analysis available deliberately does not cover our material.

### 3.3 IGFA: a real standard, but for line, not knots

The IGFA maintains genuinely rigorous procedures, and they are worth citing, but note carefully what they measure.

- Line-class records are classified by the breaking strength of the **first 5 m (16.5 ft) of line preceding the double line, leader or hook**, which must be a single homogeneous piece.
- Testing is a **wet test**: samples soaked in water for two hours before testing, **repeated five times**, with the average break used as the line's breaking strength.
- Testing is on IGFA's own **Instron 5543** tensile machine.
- Practical consequence worth encoding as domain knowledge: **braided line consistently over-tests its stated breaking strength**, so anglers deliberately buy under-rated braid to stay inside a record class.

This is a standard for **line**, not for **knots**.
IGFA is a citable authority for line-class definitions, leader rules and legality of terminal tackle, which is genuinely useful to riggermortis for rig-legality checks.
It is not a source of knot-efficiency numbers.

### 3.4 Quality assessment of fishing-media knot tests: mostly not rigorous, and you should say so

Asked directly, the honest assessment is:

**There is no peer-reviewed literature on fishing knot strength testing methodology.**
Searching for it returns climbing and rescue literature, patents, and enthusiast content.
The most-cited fishing datasets are magazine and YouTube series (for example the "Knot War" video series, and *Sport Fishing* magazine's line and knot tests).

The specific problems, which riggermortis should state on any page that reproduces such numbers:

1. **Sample sizes of 1 to 3 pulls.** Even IGFA's knot comparisons run three knots per line type. With a 45-85% spread and overlapping per-knot ranges, n=3 cannot separate two knots.
2. **Dressing and setting are uncontrolled and dominate the result.** Reported figures for the same Palomar in fluorocarbon range from 90-98% when properly seated and lubricated to 60-70% when improperly seated. That variation is larger than the difference between most knots, and it is entirely a function of the tier's technique, not the knot's topology.
3. **No standardised pull rate**, despite the Evans data suggesting rate affects the result.
4. **No published raw data, no standard deviations, no failure-location reporting** in most cases. A single mean with no dispersion is not a measurement.
5. **Commercial conflict of interest** is common: many tests are run by line or tackle vendors, or by channels sponsored by them.
6. **Line lot variation is unaccounted for.** Mono strength varies between spools and degrades with UV and age.

**Recommendation:** riggermortis should treat all knot-strength percentages as **cited claims attributed to a named source with a stated methodology and sample size**, never as project-asserted facts, and never as a single number.
If a strength field exists at all, it should carry `n`, `line_type`, `line_diameter`, `wet/dry`, `pull_rate`, `source` and `date`, and the UI should show a range rather than a point estimate.
A schema that cannot represent "we don't know" will get filled with numbers somebody made up.

Sources: [Moyer, Tusting & Harmston, ITRS 2000 (HTML)](https://verticalsection.caves.org/nh/49/cthsc/cthsc.html), [same, PDF](https://user.xmission.com/~tmoyer/testing/High_Strength_Cord.pdf), [Evans, A Review of Knot Strength Testing (2016)](http://paci.com.au/downloads_public/knots/Knot-Testing_Thomas-Evans_2016.pdf), [IGFA International Angling Rules](https://igfa.org/international-angling-rules/), [IGFA 2024 rules PDF](https://igfa.org/wp-content/uploads/2024/06/IGFA2024_RULES-REGS_062424.pdf), [IGFA world record requirements](https://igfa.org/world-record-requirements/), [Sport Fishing line test](https://www.sportfishingmag.com/gallery/gear/2014/11/line-test/), [knots.fish strength test guide](https://knots.fish/guides/fishing-knot-strength-test/).

---

## 4. Rig Constraint Validation

### 4.1 The honest starting point

A rig has **no mathematical ground truth.**
There is no theorem that says a Carolina rig is correct and a broken Carolina rig is not.
Nothing here will ever be *proved*.

What *can* be done is precisely what electronics and CAD do: **encode the rig as a typed graph and check it against a rule set that humans wrote and can argue about.**
That is not proof.
It is conformance checking, and conformance checking is genuinely valuable as long as nobody confuses the two.

### 4.2 The analogy to borrow: EDA electrical rule checking (ERC), with a degree-of-freedom borrow from CAD

**Chosen analogy: ERC/DRC from electronic design automation, specifically the KiCad model.**

**Why it wins.**

1. **The data model is already identical.**
   A schematic is a netlist: parts with **typed pins** joined by **nets**.
   A rig is a graph of components (hook, swivel, bead, egg sinker, bobber stop, split ring, leader, main line) joined at **typed terminations** (loop, eye, ring, crimp, knot, slide).
   You do not have to invent a data model; you have to name your pin types.
2. **The check catalogue is already the check catalogue we want.**
   ERC's core is a **pin-type conflict matrix**: KiCad types pins as Input, Output, Bidirectional, Tri-state, Passive, Free, Unspecified, Power Input, Power Output, Open Collector, Open Emitter, Not Connected, and flags forbidden pairings (Output driving Output, Power Output shorted to Power Output) as errors, and dubious ones (Unspecified against anything, Input with no driver) as warnings.
   That is exactly the shape of "this knot type may not be tied to that hardware", "braid may not be tied with a clinch", "this termination has nothing on the other side".
3. **The severity and waiver culture is exactly what an open, contributor-driven project needs.**
   KiCad lets the user **edit the conflict matrix** and downgrade a rule from error to warning, and supports per-rule severity plus explicit exclusions.
   riggermortis will absolutely have legitimate rigs that violate a general rule, and a validator with no waiver mechanism gets disabled within a month.
   Bake in `severity: error | warning | info` and a per-rig `waiver: { rule: R012, reason: "..." }` from day one.
4. **DRC gives the second, dimensional layer for free.**
   PCB DRC checks clearance, minimum width, annular ring, board-edge distance: numeric geometric constraints against a rule table.
   That maps directly onto line diameter versus hook eye inner diameter, split-ring gap versus line diameter, leader length versus rod length, weight rating versus rod rating.
   Two layers, ERC for connectivity and DRC for dimensions, is the right decomposition.

**The one thing ERC does not cover, and where CAD comes in.**
The "sliding weight with nothing to stop it" problem is a **kinematic** fault, not a connectivity fault.
Borrow the CAD assembly notion of **degrees of freedom**: a component is under-constrained when it retains more freedom than intended, over-constrained when constraints conflict, and fully constrained when the remaining motion is exactly what the designer wanted.
Concretely for a rig, model each sliding component as travelling along an interval of line and require an explicit **stop** at each end of the permitted interval, where a stop is a knot, bead-plus-knot, bobber stop, crimp, swivel, or the rod tip.
Then "egg sinker with no swivel or bead below it" becomes a hard, decidable graph check: *the traversal interval of a slider is unbounded in the terminal direction*.
Note the subtlety that makes this honest: an unbounded slider is sometimes **intentional** (a free-sliding float, a slip sinker deliberately allowed to run to the swivel), so this is a rule with a legitimate waiver, not a law.
Also borrow **interference detection** in spirit only: "can this component physically pass over that one" reduces to comparing a bore or gap dimension against an outer diameter, which is DRC arithmetic, not a geometry kernel.

**Why the other candidates lose.**

- **Chemistry (SMILES/InChI validity, RDKit sanitization).**
  Rejected as the primary analogy, but **keep its central lesson**, which is the single best framing device in this whole document.
  RDKit sanitization checks valence, kekulization, aromaticity perception, conjugation, hybridization, chirality cleanup and hydrogen counts, in effect asking whether the molecule can be written as an octet-complete Lewis structure.
  Passing sanitization is the standard "chemical validity" test in cheminformatics, and it is well documented that **a molecule can sanitize cleanly and still be nonsense**: unstable, unsynthesizable, or physically implausible (see the PoseBusters work on docking poses that are formally valid and physically impossible).
  That is precisely the trap riggermortis must avoid: a rig that passes every structural check and is still a bad rig.
  The domain machinery does not transfer, though, because chemistry's rules come from conservation laws (valence is not negotiable) while every rig rule is a human convention.
- **CAD assembly constraint solvers.**
  Right for the kinematic subset, wrong as the whole framework: they presume precise geometry, mate constraints and a geometric constraint solver, and riggermortis will never have a metrically accurate 3D model of a rig. Borrow DOF reasoning, not the solver.
- **Recipe and ingredient ontologies.**
  Mostly taxonomy and substitution graphs with no real constraint engine, no severity model, no waiver culture. Nothing to borrow.
- **LEGO part compatibility.**
  The stud/anti-stud typed-connector idea is a genuinely nice miniature of "what can attach to what", and it is worth stealing as intuition, but the surrounding tooling is shallow and there is no notion of rule severity or design-rule review. Not a source of engineering practice.

### 4.3 A concrete starter rule catalogue

Structured the way a `.kicad_dru` file or an ERC matrix would be.
Each rule gets an ID, a severity, a machine-checkable predicate over the rig graph, and a waiver path.

```yaml
# Connectivity rules (ERC layer): decidable over the typed rig graph
R001  error    Every component terminal is connected, or explicitly marked no_connect.
R002  error    Referential integrity: every referenced knot, component and material exists in the catalogue.
R003  error    Exactly one path from rod/reel to each terminal hook. No orphan subgraphs.
R004  error    Termination type compatibility matrix: (knot_type x line_type x hardware_eye_type) -> allowed | warn | forbidden.
                e.g. clinch x braid = forbidden;  palomar x braid = allowed;  san_diego_jam x fluoro = allowed.
R005  error    A slider component has a bounded traversal interval: a stop exists in BOTH directions.
                Waivable when the free run is intentional (slip float, running sinker to swivel).
R006  warning  A knot appears in a load path it is not rated for (loop knot used as a terminal connection, etc.).
R007  warning  Line-to-line junction connects incompatible materials without an appropriate joining knot
                (braid-to-fluoro requires FG / Alberto / double uni, not a blood knot).
R008  info     Component count in a load path exceeds N; each junction is a failure point.

# Dimensional rules (DRC layer): arithmetic against catalogue specs
R101  error    line_diameter <= hook_eye_inner_diameter, with a doubled-line multiplier where the knot doubles.
R102  error    A sliding component's bore >= line_diameter, and its bore < the outer diameter of its stop.
                (This is the "can it pass over that" check, reduced to arithmetic.)
R103  warning  leader_length <= rod_length - guide_clearance, or the rig is uncastable on a spinning setup.
R104  warning  total_terminal_weight within the rod's stated lure-weight rating.
R105  warning  hook_gap vs intended bait size, per a curated table.

# Ordering / topology-of-assembly rules
R201  error    Ordering along the main line is a total order and matches the declared sequence
                (bead before knot, stop above float, sinker above swivel).
R202  error    No component is declared both above and below another (cycle detection in the ordering relation).
```

Every one of R001-R005 and R201-R202 is a **graph theorem** once the model is fixed: they are provable relative to the model, and they are the rules worth building first.
R004, R007, R101-R105 are **table lookups**: fully automatic, but their correctness is entirely the correctness of the human-authored table.
That distinction is the whole of section 6.

### 4.4 The thing this can never catch

ERC does not prove a circuit works.
It proves the schematic obeys the rules.
A rig validator will be identical: it can prove the rig is well-formed and rule-conformant, and it can never prove the rig catches fish, holds a big one, or is the right choice for the situation.
Write that sentence into the project README before writing the validator, or the validator's green checkmark will be read as a claim it is not making.

Sources: [KiCad ERC pin conflict matrix and severity](https://kicad-sch-api.readthedocs.io/en/latest/ERC_PRD.html), [KiCad ERC user guide](https://kicad-sch-api.readthedocs.io/en/latest/ERC_USER_GUIDE.html), [FOSDEM 2026, "A Love Letter to KiCad ERC"](https://fosdem.org/2026/events/attachments/3999GS-kicad-erc-love-letter/slides/267701/presentat_lyne5sk.pdf), [RDKit Book, molecular sanitization](https://www.rdkit.org/docs/RDKit_Book.html#molecular-sanitization), [PoseBusters: valid but physically impossible](https://arxiv.org/pdf/2308.05777), [CAD constraint generation and design intent](https://arxiv.org/html/2504.13178v1), [geometric constraint solving workbench](https://www.cad-journal.net/files/vol_5/CAD_5(1-4)_2008_471-482.pdf).

---

## 5. Provenance Patterns

For everything that cannot be machine-checked, correctness becomes an **editorial** property, and serious reference projects solve it with schema, not with prose.
Three patterns are worth copying more or less wholesale.

### 5.1 Wikidata: references and ranks at the statement level

The pattern that matters is that **the citation attaches to the individual claim, not to the page.**
A Wikidata statement is `(subject, property, value)` plus **qualifiers** (context: as-of date, measurement conditions, determination method) plus **references** (where this specific claim came from) plus a **rank**.

Ranks are the piece most projects miss:

- **normal** is the default.
- **preferred** marks the value that should be shown when several compete, for example the most recent or most authoritative measurement.
- **deprecated** marks a statement **known to be wrong or outdated**, and critically, it is **kept, not deleted**.

And the reason is itself a structured field, not a comment: `P2241 reason for deprecated rank` (which Wikidata's constraints make mandatory when deprecating) and `P7452 reason for preferred rank`.

That gives you exactly the record riggermortis needs for "this field, from this source, disputed by that source": multiple statements for the same property, each with its own reference, ranked, with the losing one retained and machine-readably explained.
Deleting a wrong claim destroys the evidence that the question was ever contested, and guarantees somebody re-adds it in two years.

### 5.2 OpenStreetMap: the verifiability principle

OSM's rule is blunt and it is the right instinct for riggermortis: a tag is verifiable **if and only if independent users observing the same feature would make the same observation every time**.
Everything an editor records should be demonstrable as true or false by another mapper.
Subjective judgements are explicitly flagged as problematic for exactly this reason, and for anything not observable on the ground (a planned road, for instance) OSM requires the source to be specified on the object.

Ported to riggermortis, this is a **triage question to run over every field in the schema**: could two independent, competent anglers tie this knot from this record and get the same object?
"Crossing sequence" passes.
"Best knot for tarpon" fails, and needs to become either a cited claim or an explicitly-labelled opinion.
Running that test across the schema, before building anything, will find more problems than any validator.

### 5.3 GBIF and Darwin Core: names versus things, and disputed authority

The taxonomic model solves a problem riggermortis definitely has: **the same knot has several names, the same name refers to several knots, and different authorities disagree.**

The mechanism:

- A stable identifier for the concept, separate from any name (GBIF taxon keys).
- Every name is its own record carrying `dwc:taxonomicStatus` (accepted, synonym, and so on) from a controlled vocabulary, plus `dwc:nomenclaturalStatus` for the nomenclatural rationale.
- **Synonyms are first-class records** that point at the accepted record via `dwc:acceptedNameUsageID`, rather than being deleted or merged away.
- Exactly one name is accepted at a time; the rest remain queryable.
- Authorship and the authority citation travel with the name.
- The backbone is assembled by merging sources **in a declared priority order**, and conflicts are resolved by that priority rather than ad hoc.

That is the pattern for "the Uni Knot, the Duncan Loop and the Grinner are the same knot", and for "source A calls this the Improved Clinch, source B disagrees".

### 5.4 What the record actually looks like

Combining all three into one concrete shape:

```yaml
knot:
  id: KNOT-0042                      # stable, opaque, never reused (GBIF taxon key pattern)
  canonical_name: "Uni Knot"
  names:                             # names are records, not strings (Darwin Core pattern)
    - name: "Uni Knot"
      status: accepted
      source: SRC-IGFA-2024
    - name: "Duncan Loop"
      status: synonym
      accepted_name_id: KNOT-0042
      source: SRC-ASHLEY-1944
    - name: "Grinner"
      status: synonym
      accepted_name_id: KNOT-0042
      region: GB
      source: SRC-ANGTIMES-2019

  topology:                          # machine-verifiable, Tier A (see section 6)
    pd_code: [[1,5,2,4],[3,1,4,6],[5,3,6,2]]
    closure: end_to_end
    verified_by:
      library: regina
      version: "7.3.1"
      checks: [well_formed, not_unknot, invariants_snapshot]
      run_at: 2026-08-10

  claims:                            # everything else is a ranked, referenced statement (Wikidata pattern)
    - property: knot_strength_pct
      value: {min: 90, max: 95}
      rank: preferred
      reason_for_preferred_rank: largest_sample_size
      qualifiers:
        line_type: braid
        line_diameter_lb: 20
        condition: wet
        sample_size: 30
        pull_rate: "unstated"
      references:
        - source: SRC-KNOTWAR-2021
          retrieved: 2026-08-10
          methodology_grade: C     # project's own rubric; see section 3.4
    - property: knot_strength_pct
      value: 98
      rank: deprecated
      reason_for_deprecated_rank: single_sample_no_dispersion   # kept, not deleted (P2241 pattern)
      qualifiers: {line_type: braid, sample_size: 1}
      references:
        - source: SRC-VENDORBLOG-2018
          conflict_of_interest: vendor_published

  verifiability:                     # OSM pattern, applied per field
    crossing_sequence: on_the_ground     # two anglers reproduce it identically
    strength_pct: cited_only             # not reproducible without a tensile tester
    best_for_species: opinion            # explicitly labelled, never presented as fact
```

The three load-bearing ideas: **the citation lives on the claim**, **wrong claims are demoted with a machine-readable reason rather than deleted**, and **every field is explicitly labelled by how it can be verified at all**.

Sources: [Wikidata Help:Ranking](https://www.wikidata.org/wiki/Help:Ranking), [Wikidata data model](https://www.wikidata.org/wiki/Wikidata:Data_model), [P2241 reason for deprecated rank](https://www.wikidata.org/wiki/Property:P2241), [P7452 reason for preferred rank](https://www.wikidata.org/wiki/Property:P7452), [OSM Verifiability](https://wiki.openstreetmap.org/wiki/Verifiability), [OSM Good practice](https://wiki.openstreetmap.org/wiki/Good_practice), [GBIF backbone taxonomy explained](https://data-blog.gbif.org/post/gbif-backbone-taxonomy/), [GBIF checklist best practices](https://ipt.gbif.org/manual/en/ipt/latest/best-practices-checklists), [GBIF on synonymy](https://docs.gbif.org/course-data-use/en/synonymy.html).

---

## 6. What Can Actually Be Machine-Validated

This is the section that matters.
Three tiers, stated plainly, with no hedging in either direction.

### Tier A. Provably checkable by machine

These are mathematical or graph-theoretic facts.
A passing result is a proof, not an opinion, and no human review can overturn it.
Note the scope carefully: A1 to A5 are proofs about the *topology*, A6 to A8 are proofs about the *encoding relative to the model*.
Neither is a proof about a physical knot in physical line.

| # | Check | What it proves | How |
| --- | --- | --- | --- |
| A1 | Diagram well-formedness | The crossing sequence describes a realizable planar diagram: four arcs per crossing, arcs pair, underlying 4-valent graph is planar, orientations consistent | Regina or Spherogram reject a malformed PD/Gauss code outright |
| A2 | The authored knot is the unknot | Exact decision. The thing falls apart | Regina's Burton-Ozlen unknot recognition (exact, exponential worst case, instant at 3-12 crossings); or simplification to 0 crossings, which is a constructive proof |
| A3 | The authored knot is NOT the unknot | Exact, one-directional proof | Any nontrivial invariant. Jones ≠ 1 proves knottedness |
| A4 | An edit changed the knot type | Exact, one-directional proof of *difference* | Committed invariant snapshot (Jones + HOMFLY + crossing number + writhe) diffed in CI. **Inequality proves change. Equality does NOT prove sameness** |
| A5 | An edit flipped handedness / mirrored the knot | Provable when Jones differs between the knot and its mirror | Jones polynomial (Alexander is blind to chirality; do not use it for this) |
| A6 | Rig graph structural integrity | Connectivity, no orphan subgraphs, no cycles in the ordering relation, no dangling terminations, referential integrity of every catalogue reference | Plain graph algorithms over the rig model |
| A7 | Unbounded slider detection | A component declared as sliding has an unbounded traversal interval. **The "egg sinker with nothing to stop it" bug is decidable** | Interval-bounds check on the ordering relation. Waivable, because free-running is sometimes intentional |
| A8 | Schema and vocabulary conformance | Every enum value is in its controlled vocabulary, every unit is dimensionally consistent, every claim carries a reference | JSON Schema / SHACL, plus a "no uncited claim" gate |

**A4 is the single highest-value thing to build.**
It converts "somebody edited a crossing sequence and nobody noticed the knot changed" from an invisible, permanent data corruption into a red CI run.
Nothing else in this document buys that much correctness for that little work.

### Tier B. Heuristically checkable

Automatic, useful, and defensible, but the answer is only as good as a human-authored table, a model parameter, or a published heuristic outside its validation domain.
Every one of these should ship with its provenance visible in the UI, not buried.

| # | Check | Why it is only heuristic |
| --- | --- | --- |
| B1 | Knot-to-line-type compatibility (clinch on braid, blood knot on braid-to-fluoro) | A curated lookup table. Machine-enforced, human-authored, and wrong wherever the table is wrong |
| B2 | Dimensional fits: line diameter versus hook eye, slider bore versus stop diameter, leader versus rod length | Arithmetic against manufacturer specs. Tolerances, real eye geometry and doubled-line multipliers make it approximate |
| B3 | Relative stability ranking from topological counting rules (twist charge, crossing number, handedness) | Patil et al. 2020 is real, published, and computable straight from the crossing sequence with zero dependencies. **But it was validated on climbing and sailing bends in high-friction rope, never on fishing knots in mono, fluoro or braid.** Ship the rank, ship the caveat |
| B4 | Physics simulation of tightening (DER + IMC via DisMech) | Output is entirely determined by bending stiffness, twisting stiffness, contact stiffness and self-friction coefficient, none of which are published for fishing line. A simulation result is a statement about the model, not the knot. Also a heavy GPL-3.0 C++ stack, likely arm64-hostile in CI |
| B5 | Knot typing of an open curve via closure sampling | Heuristic by construction: the answer depends on closure direction and is reported as a distribution. Millett/Stasiak found direct end-to-end closure agrees with the dominant type in the overwhelming majority of cases, which is reassurance, not proof |
| B6 | Duplicate / alias detection ("is this new knot just the Uni?") | Matching invariants is strong evidence and the right first pass, but invariants are incomplete, so equality is not identity |
| B7 | Castability, rod-rating, hook-gap-to-bait rules | Engineering rules of thumb with real disagreement among competent anglers |

### Tier C. Only verifiable by citation or human expert

No amount of software will decide these.
Attempting to make them look machine-validated is the failure mode this whole document exists to prevent.

| # | Claim | Why software cannot touch it |
| --- | --- | --- |
| C1 | Breaking strength, absolute or relative ("Palomar = 95% on braid") | No computational path. And the underlying literature is thin: the most-cited climbing meta-analysis found n=1 was the most common sample size across 1,440+ tests, ranges of 45-85% overlapping between knots, and it **explicitly excluded fishing line** |
| C2 | **Whether the tying instructions actually produce the encoded knot** | **The biggest silent gap in the entire project.** You can verify the encoded topology perfectly and still ship step-by-step instructions and illustrations that produce a different knot. The encoding and the tutorial are two artifacts and only a human (or computer vision on a physically tied knot) closes the loop between them. Every knot record needs a human sign-off field for exactly this |
| C3 | Dressing and setting technique | The difference between roughly 60% and roughly 98% retention on the same knot. Entirely outside any topological encoding, and larger than the difference between most knots |
| C4 | Fitness for purpose: species, conditions, tackle, "best knot for X" | A value judgement. Label it opinion (OSM verifiability test) or make it a cited claim with an attributed author |
| C5 | Legality and ethics: hook types, terminal-tackle rules, IGFA record compliance, jurisdictional regs | Machine-checkable **only** against a curated, dated, jurisdiction-scoped rule set that a human must maintain. That is provenance work, not validation |
| C6 | Naming and nomenclature: what "Improved Clinch" refers to | Pure authority problem. Solve it exactly the way taxonomy does (section 5.3), not by picking a winner |
| C7 | Safety-critical claims of any kind | Do not make them. Cite them or omit them |

### The line, in one paragraph

Software can prove that **the knot you encoded is the knot you meant to encode**, that **it is not secretly the unknot**, that **nobody's edit silently changed it**, and that **the rig graph is structurally coherent and rule-conformant**.
Software can *suggest*, from curated tables and published heuristics, that a knot suits a line type and that one bend is likely more stable than another.
Software can prove **nothing** about whether a knot holds a fish, how strong it is, or whether your instructions teach it correctly.
That last category is not a gap to close with more computation; it is a permanent editorial responsibility, and the right response is the provenance schema in section 5, not a better algorithm.

The most valuable thing riggermortis can do is not build the fanciest validator.
It is to **label every field in the data model with which tier it belongs to**, and to make that label visible to readers.
A project that says "this part is proved, this part is a rule we wrote, this part is somebody's cited claim" is more trustworthy than one with a green checkmark on everything.
