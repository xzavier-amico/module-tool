package format

import "strings"

// textSectionNames maps canonical section -> output section name per format.
func textSectionNames(format string) map[string]string {
	switch format {
	case Surge:
		return map[string]string{
			"General":        "General",
			"Rule":           "Rule",
			"Url Rewrite":    "URL Rewrite",
			"Header Rewrite": "Header Rewrite",
			"Body Rewrite":   "Body Rewrite",
			"Map Local":      "Map Local",
			"Script":         "Script",
			"MITM":           "MITM",
			"Host":           "Host",
			"Panel":          "Panel",
			"Task":           "Task",
			"DNS":            "General",
		}
	case Loon:
		return map[string]string{
			"General":        "General",
			"Rule":           "Rule",
			"Url Rewrite":    "Url Rewrite",
			"Header Rewrite": "Header Rewrite",
			"Body Rewrite":   "Body Rewrite",
			"Map Local":      "Map Local",
			"Script":         "Script",
			"MITM":           "MITM",
			"Host":           "Host",
			"Panel":          "Panel",
			"Task":           "Script",
			"DNS":            "General",
		}
	case QX:
		return map[string]string{
			"Rule":           "filter_local",
			"Url Rewrite":    "rewrite_local",
			"Header Rewrite": "rewrite_local",
			"Body Rewrite":   "rewrite_local",
			"Script":         "rewrite_local",
			"MITM":           "mitm",
			"Host":           "host",
			"Task":           "task_local",
			"Panel":          "rewrite_local",
		}
	}
	return map[string]string{}
}

// RenderText renders a module to Surge/Loon/QX text.
func RenderText(m *Module, format string) string {
	var b strings.Builder
	// metadata
	metaKeys := []string{"name", "desc", "author", "homepage", "icon", "open_url", "manual", "version", "category"}
	for _, k := range metaKeys {
		if v, ok := m.Meta[k]; ok && strings.TrimSpace(v) != "" {
			switch format {
			case QX:
				// QX has no standard metadata header; emit as comments
				b.WriteString("#! " + k + "=" + v + "\n")
			default:
				b.WriteString("#!" + k + "=" + v + "\n")
			}
		}
	}
	names := textSectionNames(format)
	for _, sec := range m.SectionOrder {
		outName, ok := names[sec]
		if !ok {
			// pass through unknown sections with original name
			outName = sec
		}
		// skip empty DNS merged into General? We'll still render if lines
		lines := m.Sections[sec]
		if len(lines) == 0 {
			continue
		}
		// QX: skip sections that have no mapping? mapping exists for common
		b.WriteString("\n[" + outName + "]\n")
		for _, l := range lines {
			b.WriteString(l + "\n")
		}
	}
	return b.String()
}
