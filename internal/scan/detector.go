package scan

import (
	"os"
	"path/filepath"
	"sort"
)

type Result struct {
	Path        string             `json:"path"`
	Confidence  map[string]float64 `json:"confidence"`
	Tags        []string           `json:"tags"`
	Evidence    []string           `json:"evidence"`
	Suggestions []string           `json:"suggestions,omitempty"`
}

func Detect(path string) (Result, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return Result{}, err
	}
	res := Result{Path: abs, Confidence: map[string]float64{}}

	add := func(key string, score float64) {
		if prev, ok := res.Confidence[key]; !ok || score > prev {
			res.Confidence[key] = score
		}
	}
	check := func(rel string, scoreMap map[string]float64) {
		if exists(filepath.Join(abs, rel)) {
			res.Evidence = append(res.Evidence, rel)
			for k, v := range scoreMap {
				add(k, v)
			}
		}
	}

	check("package.json", map[string]float64{"node-js": 0.95})
	check("pnpm-lock.yaml", map[string]float64{"node-js": 0.90})
	check("package-lock.json", map[string]float64{"node-js": 0.90})
	check("yarn.lock", map[string]float64{"node-js": 0.90})
	check("tsconfig.json", map[string]float64{"typescript": 0.96})
	check("next.config.js", map[string]float64{"nextjs": 0.95, "react": 0.90})
	check("next.config.mjs", map[string]float64{"nextjs": 0.95, "react": 0.90})
	check("next.config.ts", map[string]float64{"nextjs": 0.95, "react": 0.90, "typescript": 0.95})
	check("schema.prisma", map[string]float64{"prisma": 0.92})
	check("Dockerfile", map[string]float64{"docker": 0.90})

	if exists(filepath.Join(abs, "src", "App.tsx")) || exists(filepath.Join(abs, "src", "App.jsx")) {
		res.Evidence = append(res.Evidence, "src/App.tsx|jsx")
		add("react", 0.88)
	}
	if exists(filepath.Join(abs, "app", "layout.tsx")) || exists(filepath.Join(abs, "app", "page.tsx")) {
		res.Evidence = append(res.Evidence, "app/layout.tsx|page.tsx")
		add("nextjs", 0.93)
		add("react", 0.90)
	}

	res.Tags = normalizeTags(res.Confidence)
	sort.Strings(res.Evidence)
	return res, nil
}

func normalizeTags(conf map[string]float64) []string {
	tags := make([]string, 0, len(conf))
	for k, v := range conf {
		if v < 0.50 {
			continue
		}
		switch k {
		case "react", "nextjs":
			tags = append(tags, "tech:"+k)
		case "prisma", "docker":
			tags = append(tags, "tool:"+k)
		case "typescript", "node-js":
			tags = append(tags, "lang:"+k)
		default:
			tags = append(tags, "tech:"+k)
		}
	}
	sort.Strings(tags)
	return tags
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
