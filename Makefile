CUE   := cue
KINDS := Component Line Knot Rigging Rig Source

.PHONY: check validate conformance rules schema site serve look fmt

## Shape, enums and grammar, for the whole dataset in one pass.
## The validate package keys every record by its own id, so records collide
## into one value instead of at the root, and kind picks its own schema.
validate:
	@$(CUE) vet ./validate/ data/*/*.yaml -l '"records"' -l 'id'
	@echo "$$(ls data/*/*.yaml | wc -l | tr -d ' ') records valid"

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

## Cross-record rules CUE cannot express: file paths, graph walks, warnings.
rules:
	@go run ./cmd/rig

## Published JSON Schema, one per kind, for consumers who do not use CUE.
## Not OpenAPI: that describes APIs, and this is a dataset (ADR 0002).
schema:
	@mkdir -p schema
	@for k in $(KINDS); do \
	  $(CUE) def ./cue/ --out jsonschema -e "#$$k" > "schema/$$(echo $$k | tr A-Z a-z).json" || exit 1; \
	done
	@echo "$(words $(KINDS)) schemas written to schema/"

## Build the static site. Records in, HTML and SVG out, nothing runs after.
site:
	@go run ./cmd/site

serve: site
	@echo "http://localhost:8080"; cd site && python3 -m http.server 8080

## Render rope test cases and screenshot them. Visual work gets looked at.
look:
	@mkdir -p .look
	@LOOK_OUT=.look/look.html go run ./cmd/look
	@"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless \
	  --disable-gpu --screenshot=.look/look.png --window-size=700,1080 \
	  --hide-scrollbars --default-background-color=282a36 \
	  "file://$$PWD/.look/look.html" >/dev/null 2>&1
	@echo ".look/look.png"

fmt:
	@$(CUE) fmt ./cue/ ./validate/
	@gofmt -w ./internal ./cmd
	@yamllint data/ conformance/ 2>/dev/null || true

check: validate conformance rules
