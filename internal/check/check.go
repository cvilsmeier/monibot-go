package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	checkVersion() // versions in README.md and version.go must match
	checkMarkers()
	checkExampleSnippets()
}

func checkVersion() {
	// README.md
	readmeVersion := cutout(readFile("README.md"), "### v", "\n")
	// version.go
	versionVersion := cutout(readFile("version.go"), "\"v", "\"")
	if readmeVersion != versionVersion {
		log.Print("README.md                   ", readmeVersion)
		log.Print("internal/version/version.go ", versionVersion)
		log.Fatal("version mismatch")
	}
}

func checkMarkers() {
	filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		panicIf(err)
		if entry.IsDir() || strings.HasPrefix(path, ".git/") {
			return nil
		}
		text := readFile(path)
		if strings.Contains(text, "cv"+"vvvvv") || strings.Contains(text, "FIX"+"ME") {
			log.Fatalf("found cv"+"vvvvvv/FIX"+"ME marker in %s", path)
		}
		return nil
	})
}

func checkExampleSnippets() {
	// example/main.go
	mainSnippet := cutout(normalize(readFile("example/main.go")), "func main() {", "\n}")
	// README.md
	readmeSnippet := cutout(normalize(readFile("README.md")), "func main() {", "\n}")
	if readmeSnippet != mainSnippet {
		log.Print("want")
		log.Print(mainSnippet)
		log.Print("--------")
		log.Print("have")
		log.Print(readmeSnippet)
		log.Print("--------")
		log.Fatal("README.md example mismatch")
	}
	// doc.go
	docSnippet := cutout(normalize(readFile("doc.go")), "func main() {", "\n}")
	if docSnippet != mainSnippet {
		log.Print("want")
		log.Print(mainSnippet)
		log.Print("--------")
		log.Print("have")
		log.Print(docSnippet)
		log.Print("--------")
		log.Fatal("doc.go example mismatch")
	}
}

func normalize(s string) string {
	s = replace(s, "\r", "")
	s = replace(s, "\t", " ")
	s = replace(s, "  ", " ")
	s = replace(s, "\n ", "\n")
	return strings.TrimSpace(s)
}

func replace(str, old, new string) string {
	var i int
	for strings.Contains(str, old) {
		str = strings.ReplaceAll(str, old, new)
		i++
		if i > 100 {
			panic("endless loop")
		}
	}
	return str
}

func readFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func cutout(s, from, to string) string {
	rest := s
	_, rest, found := strings.Cut(rest, from)
	if !found {
		log.Fatalf("cannot cutout from %q in %q", from, s)
	}
	rest, _, found = strings.Cut(rest, to)
	if !found {
		log.Fatalf("cannot cutout to %q in %q", to, s)
	}
	return rest
}

func panicIf(err any) {
	if err != nil {
		panic(err)
	}
}
