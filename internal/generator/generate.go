package generator

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xzavier-amico/module-tool/internal/config"
	"github.com/xzavier-amico/module-tool/internal/converter"
)

func fetchSource(src string) (string, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := http.Get(src)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	b, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Generate(cfg *config.Config, baseDir string) error {
	srcDir := filepath.Join(baseDir, cfg.SrcDir)
	distDir := filepath.Join(baseDir, cfg.DistDir)
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(distDir, 0755)

	for _, job := range cfg.Jobs {
		if job.Name == "" {
			return fmt.Errorf("job missing name")
		}
		content, err := fetchSource(job.Source)
		if err != nil {
			return fmt.Errorf("fetch %s: %w", job.Source, err)
		}
		// save source
		srcExt := config.Ext(job.From)
		srcPath := filepath.Join(srcDir, job.Name+srcExt)
		os.WriteFile(srcPath, []byte(content), 0644)

		for _, target := range job.Targets {
			out := converter.Convert(job.From, target, content)
			// apply overrides
			m := converter.Parse(job.From, content)
			if job.Name != "" {
				m.Meta["name"] = job.Name
			}
			if job.Desc != "" {
				m.Meta["desc"] = job.Desc
			}
			if job.Author != "" {
				m.Meta["author"] = job.Author
			}
			if job.Icon != "" {
				m.Meta["icon"] = job.Icon
			}
			// re-render with overrides
			out = converter.Render(target, m)
			outDir := distDir
			if job.OutDir != "" {
				outDir = filepath.Join(baseDir, job.OutDir)
			}
			os.MkdirAll(outDir, 0755)
			outPath := filepath.Join(outDir, job.Name+"."+target+config.Ext(target))
			// avoid double ext: name.target.ext? We'll use name.ext only
			outPath = filepath.Join(outDir, job.Name+config.Ext(target))
			if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
				return err
			}
			fmt.Printf("generated %s -> %s\n", job.Name, outPath)
		}
	}
	return nil
}
