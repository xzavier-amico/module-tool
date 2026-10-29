package format

import (
	"bufio"
	"regexp"
	"strings"
)

var sectionAlias = map[string]string{
	"general":        "General",
	"rule":           "Rule",
	"rules":          "Rule",
	"url rewrite":    "Url Rewrite",
	"url_rewrite":    "Url Rewrite",
	"urlrewrite":     "Url Rewrite",
	"rewrite":        "Url Rewrite",
	"rewrite_local":  "Url Rewrite",
	"rewrite_remote": "Url Rewrite",
	"header rewrite": "Header Rewrite",
	"header_rewrite": "Header Rewrite",
	"body rewrite":   "Body Rewrite",
	"body_rewrite":   "Body Rewrite",
	"map local":      "Map Local",
	"map_local":      "Map Local",
	"script":         "Script",
	"scripts":        "Script",
	"mitm":           "MITM",
	"host":           "Host",
	"hosts":          "Host",
	"panel":          "Panel",
	"panels":         "Panel",
	"task":           "Task",
	"task_local":     "Task",
	"tasks":          "Task",
	"dns":            "DNS",
	"filter_local":   "Rule",
	"filter_remote":  "Rule",
}

func CanonicalSection(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	key = strings.ReplaceAll(key, "-", " ")
	key = strings.ReplaceAll(key, "_", " ")
	key = strings.Join(strings.Fields(key), " ")
	if v, ok := sectionAlias[key]; ok {
		return v
	}
	// preserve original casing-ish
	return strings.TrimSpace(name)
}

var headerRe = regexp.MustCompile(`^\[(.+)\]$`)
var metaRe = regexp.MustCompile(`^#!\s*([A-Za-z0-9_\-]+)\s*=\s*(.*)$`)

// ParseText parses Surge/Loon/QX style text modules.
func ParseText(content, srcFormat string) *Module {
	m := NewModule()
	m.SourceFormat = srcFormat
	m.Raw = content
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	current := ""
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if mm := metaRe.FindStringSubmatch(trimmed); mm != nil {
			m.Meta[strings.ToLower(mm[1])] = strings.TrimSpace(mm[2])
			continue
		}
		if hm := headerRe.FindStringSubmatch(trimmed); hm != nil {
			current = CanonicalSection(hm[1])
			if _, ok := m.Sections[current]; !ok {
				m.SectionOrder = append(m.SectionOrder, current)
				m.Sections[current] = []string{}
			}
			continue
		}
		if current == "" {
			current = "General"
			if _, ok := m.Sections[current]; !ok {
				m.SectionOrder = append(m.SectionOrder, current)
				m.Sections[current] = []string{}
			}
		}
		if _, ok := m.Sections[current]; !ok {
			m.SectionOrder = append(m.SectionOrder, current)
			m.Sections[current] = []string{}
		}
		m.Sections[current] = append(m.Sections[current], trimmed)
	}
	return m
}
