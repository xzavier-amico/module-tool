package format

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseStash parses a Stash override YAML into IR.
func ParseStash(content string) *Module {
	m := NewModule()
	m.SourceFormat = Stash
	m.Raw = content
	var data map[string]any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return m
	}
	// metadata
	for _, k := range []string{"name", "desc", "description", "author", "homepage", "icon", "openUrl", "open_url", "category", "version", "date"} {
		if v, ok := data[k]; ok {
			key := strings.ToLower(k)
			if key == "description" {
				key = "desc"
			}
			if key == "openurl" {
				key = "open_url"
			}
			m.Meta[key] = fmt.Sprint(v)
		}
	}
	// http
	if http, ok := data["http"].(map[string]any); ok {
		if urs, ok := http["url-rewrite"].([]any); ok {
			var lines []string
			for _, u := range urs {
				lines = append(lines, fmt.Sprint(u))
			}
			m.AddSection("Url Rewrite", lines)
		}
		if mitm, ok := http["mitm"].([]any); ok {
			var hs []string
			for _, h := range mitm {
				hs = append(hs, fmt.Sprint(h))
			}
			m.AddSection("MITM", []string{"hostname = " + strings.Join(hs, ", ")})
		}
		if scripts, ok := http["script"].([]any); ok {
			// map script-providers
			providers := map[string]string{}
			if sp, ok := data["script-providers"].(map[string]any); ok {
				for name, v := range sp {
					vm, _ := v.(map[string]any)
					if url, ok := vm["url"]; ok {
						providers[name] = fmt.Sprint(url)
					}
				}
			}
			var lines []string
			for _, s := range scripts {
				sm, _ := s.(map[string]any)
				match := fmt.Sprint(sm["match"])
				name := fmt.Sprint(sm["name"])
				typ := fmt.Sprint(sm["type"])
				kind := "http-request"
				if typ == "response" {
					kind = "http-response"
				}
				url := providers[name]
				line := kind + " " + match + " script-path=" + url
				lines = append(lines, line)
			}
			m.AddSection("Script", lines)
		}
		// header-rewrite / body-rewrite maybe strings
		if hrs, ok := http["header-rewrite"].([]any); ok {
			var lines []string
			for _, h := range hrs {
				lines = append(lines, fmt.Sprint(h))
			}
			m.AddSection("Header Rewrite", lines)
		}
		if brs, ok := http["body-rewrite"].([]any); ok {
			var lines []string
			for _, b := range brs {
				lines = append(lines, fmt.Sprint(b))
			}
			m.AddSection("Body Rewrite", lines)
		}
	}
	return m
}

// RenderStash renders IR to a Stash override YAML.
func RenderStash(m *Module) string {
	out := map[string]any{}
	if v, ok := m.Meta["name"]; ok {
		out["name"] = v
	}
	if v, ok := m.Meta["desc"]; ok {
		out["desc"] = v
	}
	if v, ok := m.Meta["author"]; ok {
		out["author"] = v
	}
	if v, ok := m.Meta["homepage"]; ok {
		out["homepage"] = v
	}
	if v, ok := m.Meta["icon"]; ok {
		out["icon"] = v
	}
	if v, ok := m.Meta["open_url"]; ok {
		out["openUrl"] = v
	}

	http := map[string]any{}
	// url-rewrite
	if lines := m.GetSection("Url Rewrite"); len(lines) > 0 {
		http["url-rewrite"] = lines
	}
	// mitm
	hosts := parseMITMHosts(m.GetSection("MITM"))
	if len(hosts) > 0 {
		http["mitm"] = hosts
	}
	// header/body rewrite
	if lines := m.GetSection("Header Rewrite"); len(lines) > 0 {
		http["header-rewrite"] = lines
	}
	if lines := m.GetSection("Body Rewrite"); len(lines) > 0 {
		http["body-rewrite"] = lines
	}
	// script
	if lines := m.GetSection("Script"); len(lines) > 0 {
		// convert to stash script entries (basic, using url from script-path)
		var scripts []map[string]any
		providers := map[string]map[string]any{}
		for i, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}
			match := fields[1]
			url := ""
			for _, f := range fields[2:] {
				if strings.HasPrefix(f, "script-path=") {
					url = strings.TrimPrefix(f, "script-path=")
				}
			}
			name := fmt.Sprintf("script-%d", i+1)
			typ := "request"
			if strings.Contains(fields[0], "response") {
				typ = "response"
			}
			scripts = append(scripts, map[string]any{
				"match": match,
				"name":  name,
				"type":  typ,
			})
			providers[name] = map[string]any{
				"url":      url,
				"interval": 86400,
			}
		}
		if len(scripts) > 0 {
			http["script"] = scripts
			out["script-providers"] = providers
		}
	}
	if len(http) > 0 {
		out["http"] = http
	}

	b, _ := yaml.Marshal(out)
	return string(b)
}
