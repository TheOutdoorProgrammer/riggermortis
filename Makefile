CUE := cue

# Directory determines kind, which is rule 7. Explicit pairs rather than
# derived, because BSD sed cannot upper-case and silent breakage here would
# make the validator report success on nothing.
PAIRS := components:Component lines:Line knots:Knot riggings:Rigging rigs:Rig sources:Source

.PHONY: check validate conformance fmt

## Validate every record against the kind its directory implies.
validate:
	@fail=0; total=0; \
	for p in $(PAIRS); do \
	  d=$${p%%:*}; k=$${p##*:}; \
	  for f in data/$$d/*.yaml; do \
	    [ -e "$$f" ] || continue; total=$$((total+1)); \
	    out=$$($(CUE) vet ./cue/ "$$f" -d "#$$k" 2>&1); \
	    if [ -n "$$out" ]; then \
	      fail=$$((fail+1)); echo "FAIL $$f"; echo "$$out" | sed 's/^/     /'; \
	    fi; \
	  done; \
	done; \
	echo "$$((total-fail))/$$total records valid"; \
	[ $$fail -eq 0 ]

## Every fixture under conformance/invalid MUST be rejected.
## One that validates is a hole in the spec, not a passing test.
conformance:
	@fail=0; total=0; \
	while read -r f k; do \
	  case "$$f" in ''|\#*) continue;; esac; \
	  total=$$((total+1)); \
	  if $(CUE) vet ./cue/ "conformance/invalid/$$f" -d "#$$k" >/dev/null 2>&1; then \
	    fail=$$((fail+1)); echo "HOLE conformance/invalid/$$f validated but must not"; \
	  fi; \
	done < conformance/manifest.txt; \
	echo "$$((total-fail))/$$total invalid fixtures correctly rejected"; \
	[ $$fail -eq 0 ]

fmt:
	@$(CUE) fmt ./cue/
	@yamllint data/ conformance/ 2>/dev/null || true

check: validate conformance
