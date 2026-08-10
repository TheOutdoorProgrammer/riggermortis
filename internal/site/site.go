// SPDX-License-Identifier: Apache-2.0

// Package site generates the static site.
//
// Everything is emitted at build time. There is no server, no client-side
// geometry, and the stepper works with scripting disabled.
package site

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/theoutdoorprogrammer/riggermortis/internal/render"
	"github.com/theoutdoorprogrammer/riggermortis/internal/spec"
)

//go:embed *.html *.css *.js
var assets embed.FS

type Step struct {
	N        int
	Prose    string
	Notation string
	SVG      template.HTML
}

type Cord struct {
	Name   string
	Colour string
}

type knotView struct {
	*spec.Record
	Steps []Step
	Cords []Cord
}

type indexEntry struct {
	Title, Href, Role, Status string
	Steps                     int
	Caution                   bool
}

func Build(s *spec.Set, out string) error {
	tmpl, err := template.ParseFS(assets, "*.html")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(out, "knots"), 0o755); err != nil {
		return err
	}
	for _, name := range []string{"style.css", "step.js"} {
		b, err := assets.ReadFile(name)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(out, name), b, 0o644); err != nil {
			return err
		}
	}

	var index []indexEntry
	for _, r := range s.All {
		if r.Kind != "knot" {
			continue
		}
		v := knotView{Record: r}
		for i, st := range r.Stages {
			step := Step{N: st.ID, Prose: st.Prose, Notation: st.Notation}
			if step.N == 0 {
				step.N = i + 1
			}
			if r.Geometry != nil && i < len(r.Geometry.Stages) {
				step.SVG = template.HTML(render.Stage(r.Geometry, i)) //nolint:gosec // generated markup, no user input
			}
			v.Steps = append(v.Steps, step)
		}
		if r.Geometry != nil {
			for _, c := range r.Geometry.Cords {
				v.Cords = append(v.Cords, Cord{Name: c, Colour: render.Colour(r.Geometry, c)})
			}
		}

		slug := strings.TrimPrefix(r.ID, "knot.")
		path := filepath.Join(out, "knots", slug+".html")
		if err := writePage(tmpl, path, "knot", v, title(r)); err != nil {
			return err
		}
		index = append(index, indexEntry{
			Title:   title(r),
			Href:    "knots/" + slug + ".html",
			Role:    r.Role,
			Status:  r.Validation.Status,
			Steps:   len(v.Steps),
			Caution: hasDoNotUse(r),
		})
	}
	sort.Slice(index, func(i, j int) bool { return index[i].Title < index[j].Title })

	return writePage(tmpl, filepath.Join(out, "index.html"), "index", index, "riggermortis")
}

func writePage(t *template.Template, path, block string, data any, pageTitle string) error {
	var body strings.Builder
	if err := t.ExecuteTemplate(&body, block, data); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	depth := ""
	if strings.Contains(path, string(filepath.Separator)+"knots"+string(filepath.Separator)) {
		depth = "../"
	}
	return t.ExecuteTemplate(f, "layout", map[string]any{
		"Title": pageTitle,
		"Body":  template.HTML(body.String()), //nolint:gosec // generated markup
		"Root":  depth,
	})
}

func title(r *spec.Record) string {
	if r.Names != nil && r.Names.Canonical != "" {
		return r.Names.Canonical
	}
	return r.ID
}

func hasDoNotUse(r *spec.Record) bool {
	for _, c := range r.Cautions {
		if c.Severity == "do-not-use" {
			return true
		}
	}
	return false
}
