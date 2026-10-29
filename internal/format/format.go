package format

// Format names
const (
	Surge  = "surge"
	Loon   = "loon"
	QX     = "qx"
	Egern  = "egern"
	Stash  = "stash"
)

var All = []string{Surge, Loon, QX, Egern, Stash}

func Valid(f string) bool {
	for _, v := range All {
		if v == f {
			return true
		}
	}
	return false
}

// Module is an intermediate representation shared by all formats.
type Module struct {
	Meta       map[string]string // name, desc, author, homepage, icon, open_url, manual, version, category
	Sections   map[string][]string
	SectionOrder []string
	SourceFormat string
	Raw          string
}

func NewModule() *Module {
	return &Module{
		Meta:     map[string]string{},
		Sections: map[string][]string{},
	}
}

func (m *Module) AddSection(name string, lines []string) {
	if _, ok := m.Sections[name]; !ok {
		m.SectionOrder = append(m.SectionOrder, name)
	}
	m.Sections[name] = append(m.Sections[name], lines...)
}

func (m *Module) GetSection(name string) []string {
	return m.Sections[name]
}
