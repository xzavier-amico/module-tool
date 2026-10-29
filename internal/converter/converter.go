package converter

import (
	"github.com/xzavier-amico/module-tool/internal/format"
)

func Parse(srcFormat, content string) *format.Module {
	switch srcFormat {
	case format.Egern:
		return format.ParseEgern(content)
	case format.Stash:
		return format.ParseStash(content)
	default:
		return format.ParseText(content, srcFormat)
	}
}

func Render(dstFormat string, m *format.Module) string {
	switch dstFormat {
	case format.Egern:
		return format.RenderEgern(m)
	case format.Stash:
		return format.RenderStash(m)
	default:
		return format.RenderText(m, dstFormat)
	}
}

func Convert(srcFormat, dstFormat, content string) string {
	if srcFormat == dstFormat {
		return content
	}
	m := Parse(srcFormat, content)
	return Render(dstFormat, m)
}
