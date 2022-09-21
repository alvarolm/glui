package main

import (
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"sort"

	"github.com/alvarolm/glui/fuzzy"
)

var loadesPackages packages

type packages struct {
	list     []*jsonPackage
	impPaths []string
}

func loadPackages() {
	cmd := exec.Command("go", "list", "-e", "--json", "...")

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln("failed to redirect go list stdout:", err.Error())
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln("failed to run go list command:", err.Error())
	}

	defer cmd.Wait()

	dec := json.NewDecoder(out)
	for {
		var p jsonPackage
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln("failed to decode go list json output:", err.Error())
		}

		loadesPackages.list = append(loadesPackages.list, &p)
		loadesPackages.impPaths = append(loadesPackages.impPaths, p.ImportPath)
	}
}

type showItem struct {
	impPath string
	index   int
	score   float64
}

func filterShowList(query string) (matches []*showItem) {
	if len(query) == 0 {
		for i, impPath := range loadesPackages.impPaths {
			matches = append(matches, &showItem{
				impPath: impPath,
				index:   i,
				score:   0,
			})
		}
		return matches
	}

	matcher := fuzzy.NewSymbolMatcher(query)

	for i, impPath := range loadesPackages.impPaths {
		m, score := matcher.Match([]string{impPath})

		if m != -1 {
			matches = append(matches, &showItem{
				impPath: impPath,
				index:   i,
				score:   score,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool { return matches[i].score > matches[j].score })

	return matches
}
