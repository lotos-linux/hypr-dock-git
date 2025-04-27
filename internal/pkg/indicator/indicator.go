package indicator

import (
	"errors"
	"fmt"
	"hypr-dock/internal/pkg/utils"
	"hypr-dock/internal/settings"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type IndicatorFile struct {
	Number    int
	Extension string
	FullName  string
}

func New(instances int, settings settings.Settings) (*gtk.Image, error) {
	available, err := GetAvailable(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to get indicators: %w", err)
	}
	if len(available) < 2 {
		return nil, errors.New("need at least 2 indicator files (0.svg and 1.svg)")
	}

	selected := available[0]
	for _, file := range available {
		if file.Number > instances {
			break
		}
		selected = file
	}

	path := filepath.Join(settings.CurrentThemeDir, "point", selected.FullName)
	return utils.CreateImageWidthScale(path, settings.IconSize, 0.56)
}

func GetAvailable(settings settings.Settings) ([]IndicatorFile, error) {
	dirPath := filepath.Join(settings.CurrentThemeDir, "point")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("indicator directory not found: %w", err)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read indicator directory: %w", err)
	}

	var files []IndicatorFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if !isSupportExt(ext) {
			continue
		}

		baseName := name[:len(name)-len(ext)]
		if num, err := strconv.Atoi(baseName); err == nil {
			files = append(files, IndicatorFile{
				Number:    num,
				Extension: ext,
				FullName:  name,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Number < files[j].Number
	})

	return files, nil
}

func isSupportExt(ext string) bool {
	supportList := []string{".svg", ".jpg", ".png"}
	return slices.Contains(supportList, ext)
}
