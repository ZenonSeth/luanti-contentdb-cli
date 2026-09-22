package content

import (
	"path/filepath"
	"testing"
)

func TestScanGames(t *testing.T) {
	dir := t.TempDir()

	// well-formed game
	writeFile(t, filepath.Join(dir, "minetest_game", "game.conf"),
		"title = Minetest Game\nauthor = Luanti\ndescription = A basic game\nrelease = 38214\n")

	// deprecated `name` instead of `title`
	writeFile(t, filepath.Join(dir, "oldstyle", "game.conf"),
		"name = Old Style Game\n")

	// no game.conf at all
	writeFile(t, filepath.Join(dir, "bare_game", "init.lua"), "-- nothing")

	games, err := ScanGames(dir)
	if err != nil {
		t.Fatal(err)
	}

	byID := make(map[string]Game)
	for _, g := range games {
		byID[g.ID] = g
	}

	if len(games) != 3 {
		t.Fatalf("expected 3 entries, got %d: %+v", len(games), games)
	}

	mtg := byID["minetest_game"]
	if !mtg.ConfOK || mtg.Title != "Minetest Game" || mtg.Author != "Luanti" || mtg.Release != 38214 {
		t.Errorf("minetest_game: unexpected result %+v", mtg)
	}

	old := byID["oldstyle"]
	if !old.ConfOK || old.Title != "Old Style Game" {
		t.Errorf("oldstyle: expected title from deprecated name field, got %+v", old)
	}

	bare := byID["bare_game"]
	if bare.ConfOK || bare.Title != "bare_game" {
		t.Errorf("bare_game: expected fallback to folder name, got %+v", bare)
	}
}
