// Package content scans and parses locally installed Luanti content
// (mods, and eventually games/texture packs).
package content

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Mod describes one mod folder found under a mods directory.
type Mod struct {
	Name            string
	Title           string
	Description     string
	Author          string
	Release         int // from mod.conf's `release`; 0 if absent
	Depends         []string
	OptionalDepends []string

	Dir       string // folder name, e.g. "everness"
	Path      string // full path to the mod folder
	IsModpack bool   // has modpack.conf instead of mod.conf; not scanned further yet

	// ConfOK is true only if mod.conf exists and has a usable `name`
	// field. If mod.conf is missing, empty, or has no name, Name
	// falls back to the folder name and ConfOK is false.
	ConfOK bool
}

// ScanMods scans the immediate subdirectories of dir for mods.
// Directories containing a modpack.conf are reported with
// IsModpack set, but are not scanned for nested mods (not implemented yet).
func ScanMods(dir string) ([]Mod, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var mods []Mod

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modDir := entry.Name()
		modPath := filepath.Join(dir, modDir)

		if fileExists(filepath.Join(modPath, "modpack.conf")) {
			mods = append(mods, Mod{
				Name:      modDir,
				Dir:       modDir,
				Path:      modPath,
				IsModpack: true,
			})
			continue
		}

		mods = append(mods, scanModDir(modDir, modPath))
	}

	return mods, nil
}

func scanModDir(dir, path string) Mod {
	m := Mod{
		Name: dir, // fallback, overwritten below if mod.conf has a name
		Dir:  dir,
		Path: path,
	}

	confPath := filepath.Join(path, "mod.conf")
	data, err := os.ReadFile(confPath)
	if err != nil {
		return m
	}

	conf := parseConfFile(data)

	name, ok := conf["name"]
	if ok && name != "" {
		m.Name = name
		m.ConfOK = true
	}

	m.Title = conf["title"]
	m.Description = conf["description"]
	m.Author = conf["author"]
	m.Release = parseIntField(conf, "release")
	m.Depends = splitList(conf["depends"])
	m.OptionalDepends = splitList(conf["optional_depends"])

	return m
}

// parseConfFile parses a Luanti "Settings" style file:
// one `key = value` per line
func parseConfFile(data []byte) map[string]string {
	result := make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		result[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return result
}

// parseIntField reads an integer field from a parsed conf map,
// defaulting to 0 if the key is missing or not a valid integer.
func parseIntField(conf map[string]string, key string) int {
	n, err := strconv.Atoi(conf[key])
	if err != nil {
		return 0
	}
	return n
}

// splitList splits a comma-separated mod.conf list field (depends,
// optional_depends), trimming whitespace and dropping empty entries.
func splitList(value string) []string {
	if value == "" {
		return nil
	}

	var out []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}

	return out
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
