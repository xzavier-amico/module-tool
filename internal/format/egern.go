package format

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Egern YAML structures (minimal)
type egernMeta struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Author      string `yaml:"author,omitempty"`
	Homepage    string `yaml:"homepage,omitempty"`
	Icon        string `yaml:"icon,omitempty"`
	OpenURL     string `yaml:"open_url,omitempty"`
	Manual      string `yaml:"manual,omitempty"`
}

type egernRule struct {
	Domain         *egernMatch `yaml:"domain,omitempty"`
	DomainSuffix   *egernMatch `yaml:"domain_suffix,omitempty"`
	DomainKeyword  *egernMatch `yaml:"domain_keyword,omitempty"`
	DomainRegex    *egernMatch `yaml:"domain_regex,omitempty"`
	DomainWildcard *egernMatch `yaml:"domain_wildcard,omitempty"`
	GeoIP          *egernMatch `yaml:"geoip,omitempty"`
	IPCIDR         *egernMatch `yaml:"ip_cidr,omitempty"`
	IPCIDR6        *egernMatch `yaml:"ip_cidr6,omitempty"`
	ASN            *egernMatch `yaml:"asn,omitempty"`
	URLRegex       *egernMatch `yaml:"url_regex,omitempty"`
	UserAgent      *egernMatch `yaml:"user_agent,omitempty"`
	DestPort       *egernMatch `yaml:"dest_port,omitempty"`
	Protocol       *egernMatch `yaml:"protocol,omitempty"`
	Default        *egernMatch `yaml:"default,omitempty"`
}

type egernMatch struct {
	Match      string `yaml:"match"`
	Policy     string `yaml:"policy"`
	NoResolve  bool   `yaml:"no_resolve,omitempty"`
}

type egernURLRewrite struct {
	Match      string `yaml:"match"`
	Location   string `yaml:"location,omitempty"`
	StatusCode int    `yaml:"status_code,omitempty"`
}

type egernMITM struct {
	Hostnames []string `yaml:"hostnames,omitempty"`
}

type egernScripting struct {
	HTTPRequest  *egernScript `yaml:"http_request,omitempty"`
	HTTPResponse *egernScript `yaml:"http_response,omitempty"`
}

type egernScript struct {
	Name      string `yaml:"name,omitempty"`
	Match     string `yaml:"match"`
	ScriptURL string `yaml:"script_url,omitempty"`
	Timeout   int    `yaml:"timeout,omitempty"`
}

type egernModule struct {
	egernMeta
	Rules       []map[string]any `yaml:"rules,omitempty"`
	URLRewrites []map[string]any `yaml:"url_rewrites,omitempty"`
	MITM        *egernMITM       `yaml:"mitm,omitempty"`
	Scriptings  []map[string]any `yaml:"scriptings,omitempty"`
}

// ParseEgern parses an Egern module YAML into IR.
func ParseEgern(content string) *Module {
	m := NewModule()
	m.SourceFormat = Egern
	m.Raw = content
	var data map[string]any
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return m
	}
	// metadata
	for _, k := range []string{"name", "description", "author", "homepage", "icon", "open_url", "openUrl", "manual", "version", "category"} {
		if v, ok := data[k]; ok {
			m.Meta[strings.ToLower(k)] = fmt.Sprint(v)
		}
	}
	// normalize open_url key
	if v, ok := m.Meta["openurl"]; ok {
		m.Meta["open_url"] = v
		delete(m.Meta, "openurl")
	}
	// rules -> Rule lines (Surge style)
	if rules, ok := data["rules"].([]any); ok {
		var lines []string
		for _, r := range rules {
			rm, _ := r.(map[string]any)
			for typ, val := range rm {
				vv, _ := val.(map[string]any)
				match := fmt.Sprint(vv["match"])
				policy := fmt.Sprint(vv["policy"])
				line := strings.ToUpper(strings.ReplaceAll(typ, "_", "-")) + "," + match + "," + policy
				// map type names
				line = strings.ReplaceAll(line, "DOMAIN-SUFFIX", "DOMAIN-SUFFIX")
				lines = append(lines, line)
			}
		}
		m.AddSection("Rule", lines)
	}
	// url_rewrites
	if urs, ok := data["url_rewrites"].([]any); ok {
		var lines []string
		for _, u := range urs {
			um, _ := u.(map[string]any)
			match := fmt.Sprint(um["match"])
			loc := fmt.Sprint(um["location"])
			sc := ""
			if v, ok := um["status_code"]; ok {
				sc = fmt.Sprint(v)
			}
			if loc == "<nil>" {
				loc = ""
			}
			line := match
			if loc != "" {
				line += " " + loc
			}
			if sc != "" && sc != "<nil>" && sc != "0" {
				line += " " + sc
			}
			lines = append(lines, line)
		}
		m.AddSection("Url Rewrite", lines)
	}
	// mitm
	if mitm, ok := data["mitm"].(map[string]any); ok {
		if hosts, ok := mitm["hostnames"].([]any); ok {
			var hs []string
			for _, h := range hosts {
				hs = append(hs, fmt.Sprint(h))
			}
			m.AddSection("MITM", []string{"hostname = " + strings.Join(hs, ", ")})
		}
	}
	// scriptings
	if scs, ok := data["scriptings"].([]any); ok {
		var lines []string
		for _, s := range scs {
			sm, _ := s.(map[string]any)
			for typ, val := range sm {
				vv, _ := val.(map[string]any)
				match := fmt.Sprint(vv["match"])
				url := fmt.Sprint(vv["script_url"])
				kind := "http-request"
				if strings.Contains(typ, "response") {
					kind = "http-response"
				}
				if strings.Contains(typ, "cron") {
					kind = "cron"
				}
				line := kind + " " + match + " script-path=" + url
				lines = append(lines, line)
			}
		}
		m.AddSection("Script", lines)
	}
	return m
}

// RenderEgern renders IR to Egern module YAML.
func RenderEgern(m *Module) string {
	out := map[string]any{}
	// metadata
	if v, ok := m.Meta["name"]; ok {
		out["name"] = v
	}
	if v, ok := m.Meta["desc"]; ok {
		out["description"] = v
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
		out["open_url"] = v
	}
	if v, ok := m.Meta["manual"]; ok {
		out["manual"] = v
	}

	// rules
	var rules []map[string]any
	for _, line := range m.GetSection("Rule") {
		parts := splitCSV(line)
		if len(parts) < 3 {
			continue
		}
		typ := strings.ToLower(strings.ReplaceAll(parts[0], "-", "_"))
		// map surge type to egern type
		typ = strings.ToLower(typ)
		typ = strings.TrimPrefix(typ, "domain")
		// basic mapping
		egType := surgeTypeToEgern(parts[0])
		rule := map[string]any{
			egType: map[string]any{
				"match":  parts[1],
				"policy": parts[2],
			},
		}
		rules = append(rules, rule)
	}
	if len(rules) > 0 {
		out["rules"] = rules
	}

	// url rewrites
	var urs []map[string]any
	for _, line := range m.GetSection("Url Rewrite") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ur := map[string]any{"match": fields[0]}
		if len(fields) >= 2 {
			ur["location"] = fields[1]
		}
		if len(fields) >= 3 {
			// status code
			sc := 0
			fmt.Sscan(fields[2], &sc)
			if sc != 0 {
				ur["status_code"] = sc
			}
		}
		urs = append(urs, ur)
	}
	if len(urs) > 0 {
		out["url_rewrites"] = urs
	}

	// mitm
	hosts := parseMITMHosts(m.GetSection("MITM"))
	if len(hosts) > 0 {
		out["mitm"] = map[string]any{"hostnames": hosts}
	}

	// scriptings
	var scs []map[string]any
	for _, line := range m.GetSection("Script") {
		// http-request pattern script-path=url
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		kind := fields[0]
		match := fields[1]
		url := ""
		for _, f := range fields[2:] {
			if strings.HasPrefix(f, "script-path=") {
				url = strings.TrimPrefix(f, "script-path=")
			}
		}
		entry := map[string]any{
			"match":      match,
			"script_url": url,
		}
		key := "http_request"
		if strings.Contains(kind, "response") {
			key = "http_response"
		}
		if strings.Contains(kind, "cron") {
			key = "cron"
		}
		scs = append(scs, map[string]any{key: entry})
	}
	if len(scs) > 0 {
		out["scriptings"] = scs
	}

	b, _ := yaml.Marshal(out)
	return string(b)
}

func splitCSV(s string) []string {
	var parts []string
	cur := ""
	for _, c := range s {
		if c == ',' {
			parts = append(parts, strings.TrimSpace(cur))
			cur = ""
		} else {
			cur += string(c)
		}
	}
	if strings.TrimSpace(cur) != "" {
		parts = append(parts, strings.TrimSpace(cur))
	}
	return parts
}

func surgeTypeToEgern(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	t = strings.ReplaceAll(t, "-", "_")
	switch t {
	case "domain":
		return "domain"
	case "domain_suffix":
		return "domain_suffix"
	case "domain_keyword":
		return "domain_keyword"
	case "domain_regex":
		return "domain_regex"
	case "domain_wildcard":
		return "domain_wildcard"
	case "ip_cidr":
		return "ip_cidr"
	case "ip_cidr6":
		return "ip_cidr6"
	case "geoip":
		return "geoip"
	case "user_agent":
		return "user_agent"
	case "url_regex":
		return "url_regex"
	case "final":
		return "default"
	default:
		return t
	}
}

func parseMITMHosts(lines []string) []string {
	var hosts []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if idx := strings.Index(l, "="); idx >= 0 {
			l = l[idx+1:]
		}
		for _, part := range strings.Split(l, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				hosts = append(hosts, part)
			}
		}
	}
	return hosts
}

// fix unused delete
func init() {}
