package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testStore(t *testing.T) (*store, string) {
	t.Helper()
	rootPath := t.TempDir()
	root, err := os.OpenRoot(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { root.Close() })
	return &store{roots: []*os.Root{root}}, rootPath
}

func writeSkill(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "SKILL.md"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestListAndReadSkills(t *testing.T) {
	skills, root := testStore(t)
	writeSkill(t, root, "second", "Second description\n\nSecond body")
	writeSkill(t, root, "first", "---\ndescription: First description\n---\nFirst body")
	list, err := skills.list()
	if err != nil {
		t.Fatal(err)
	}
	if list != "- first: First description\n- second: Second description\n" {
		t.Fatalf("unexpected list: %q", list)
	}
	content, err := skills.read("first")
	if err != nil {
		t.Fatal(err)
	}
	if content != "First body" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestRejectsInvalidNamesAndEscapes(t *testing.T) {
	skills, root := testStore(t)
	outside := filepath.Join(t.TempDir(), "SKILL.md")
	if err := os.WriteFile(outside, []byte("secret"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := skills.read("../outside"); err == nil {
		t.Fatal("expected invalid name error")
	}
	link := filepath.Join(root, "linked")
	if err := os.Symlink(filepath.Dir(outside), link); err != nil {
		t.Fatal(err)
	}
	if _, err := skills.read("linked"); err == nil {
		t.Fatal("expected symlink escape error")
	}
}

func TestParseFallback(t *testing.T) {
	description, content := parse("\nTitle\nBody")
	if description != "Title" || content != "\nTitle\nBody" {
		t.Fatalf("unexpected parse result: %q %q", description, content)
	}
}

func TestValidateName(t *testing.T) {
	for _, name := range []string{"", "Upper", "bad_name", "-bad", "bad-", "bad--name"} {
		if validateName(name) == nil {
			t.Fatalf("accepted invalid name: %q", name)
		}
	}
	if err := validateName("bigquery-metrics"); err != nil {
		t.Fatal(err)
	}
}

func TestDefaultRootsAndProjectPrecedence(t *testing.T) {
	project := t.TempDir()
	home := t.TempDir()
	t.Chdir(project)
	t.Setenv("HOME", home)
	t.Setenv("SKILLX_ROOT", "")
	writeSkill(t, filepath.Join(project, ".agents", "skills"), "shared", "Project skill")
	writeSkill(t, filepath.Join(home, ".agents", "skills"), "shared", "Home skill")
	writeSkill(t, filepath.Join(home, ".agents", "skills"), "home-only", "Home only")
	skills, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	defer skills.close()
	list, err := skills.list()
	if err != nil {
		t.Fatal(err)
	}
	if list != "- home-only: Home only\n- shared: Project skill\n" {
		t.Fatalf("unexpected list: %q", list)
	}
	content, err := skills.read("shared")
	if err != nil {
		t.Fatal(err)
	}
	if content != "Project skill" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestMissingDefaultRootsAreEmpty(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("HOME", t.TempDir())
	t.Setenv("SKILLX_ROOT", "")
	skills, err := openStore()
	if err != nil {
		t.Fatal(err)
	}
	defer skills.close()
	list, err := skills.list()
	if err != nil {
		t.Fatal(err)
	}
	if list != "No skills installed." {
		t.Fatalf("unexpected list: %q", list)
	}
}

func TestDescriptors(t *testing.T) {
	if len(descriptors) != 2 || !strings.Contains(descriptors[0], `"name":"skills"`) || !strings.Contains(descriptors[1], `"name":"skill"`) {
		t.Fatal("invalid descriptors")
	}
}
