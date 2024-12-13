package gatherer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/copyleftdev/5l4pp3r/internal/compression"
	"github.com/copyleftdev/5l4pp3r/internal/config"
	"github.com/copyleftdev/5l4pp3r/internal/model"
)

// GatherConfigFiles finds configuration files and compresses them.
func GatherConfigFiles(cfg *config.Config, algo string, level int) ([]*model.ConfigFile, error) {
	home := os.Getenv("HOME")
	if home == "" {
		exeDir, err := os.Executable()
		if err == nil {
			home = filepath.Dir(exeDir)
		} else {
			home = "/tmp"
		}
	}

	xdgConfigHome := cfg.Gather.XDGConfigHome
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(home, ".config")
	}

	xdgConfigDirs := cfg.Gather.XDGConfigDirs
	if xdgConfigDirs == "" {
		xdgConfigDirs = "/etc/xdg"
	}

	dirs := make(map[string]struct{})
	if dirExists(xdgConfigHome) {
		dirs[xdgConfigHome] = struct{}{}
	}

	for _, d := range strings.Split(xdgConfigDirs, ":") {
		if d != "" && dirExists(d) {
			dirs[d] = struct{}{}
		}
	}

	if dirExists(cfg.Gather.SystemConfigDir) {
		dirs[cfg.Gather.SystemConfigDir] = struct{}{}
	}

	var results []*model.ConfigFile
	compressor, err := compression.NewCompressor(algo, level)
	if err != nil {
		return nil, fmt.Errorf("failed to create compressor: %w", err)
	}

	for d := range dirs {
		err := filepath.Walk(d, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return nil
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			data, readErr := readFileSafely(path)
			if readErr != nil {
				return nil
			}

			compData, compErr := compressor.Compress(data)
			if compErr != nil {
				return nil
			}

			modTime := info.ModTime()
			permissions := info.Mode().String()
			size := info.Size()

			results = append(results, &model.ConfigFile{
				FilePath:     path,
				Size:         size,
				Permissions:  permissions,
				ModifiedTime: modTime,
				Data:         compData,
			})
			return nil
		})
		if err != nil {
			continue
		}
	}

	return results, nil
}

func dirExists(d string) bool {
	st, err := os.Stat(d)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func readFileSafely(path string) ([]byte, error) {
	const maxSize = 50 * 1024 * 1024 // 50MB
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Size() > maxSize {
		return nil, fmt.Errorf("file too large: %s", path)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}
