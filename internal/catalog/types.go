package catalog

import (
	"fmt"
	"strings"
)

const (
	TypeRule  = "rule"
	TypeSkill = "skill"
)

type Ref struct {
	Raw  string
	Type string
	Slug string
}

func ParseRef(raw string) (Ref, error) {
	raw = strings.TrimSpace(raw)
	if !strings.HasPrefix(raw, "@") {
		return Ref{}, fmt.Errorf("invalid identifier %q: expected @rule/<slug> or @skill/<slug>", raw)
	}
	parts := strings.SplitN(strings.TrimPrefix(raw, "@"), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Ref{}, fmt.Errorf("invalid identifier %q: expected @rule/<slug> or @skill/<slug>", raw)
	}
	t := parts[0]
	if t != TypeRule && t != TypeSkill {
		return Ref{}, fmt.Errorf("unsupported type %q: use rule or skill", t)
	}
	return Ref{Raw: raw, Type: t, Slug: parts[1]}, nil
}

func NormalizeRef(t, slug string) string {
	return "@" + t + "/" + slug
}
