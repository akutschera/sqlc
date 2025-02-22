package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func parseExecCommand(path string) (string, error) {
	var exec = struct {
		Command string `json:"command"`
	}{
		Command: "generate",
	}

	execJsonPath := filepath.Join(path, "exec.json")
	if _, err := os.Stat(execJsonPath); !os.IsNotExist(err) {
		blob, err := os.ReadFile(execJsonPath)
		if err != nil {
			return "", err
		}
		if err := json.Unmarshal(blob, &exec); err != nil {
			return "", err
		}
	}

	return exec.Command, nil
}

func regenerate(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "sqlc.json") || strings.HasSuffix(path, "sqlc.yaml") {
			cwd := filepath.Dir(path)
			command, err := parseExecCommand(cwd)
			if err != nil {
				return fmt.Errorf("failed to parse exec.json: %w", err)
			}

			if command != "generate" {
				return nil
			}

			var expectFailure bool
			if _, err := os.Stat(filepath.Join(cwd, "stderr.txt")); !os.IsNotExist(err) {
				expectFailure = true
			}

			cmd := exec.Command("sqlc-dev", "generate", "--experimental")
			cmd.Dir = cwd
			out, failed := cmd.CombinedOutput()
			if failed != nil && !expectFailure {
				return fmt.Errorf("%s: sqlc-dev generate failed\n%s", cwd, out)
			}
			if expectFailure {
				if err := ioutil.WriteFile(filepath.Join(cwd, "stderr.txt"), out, 0644); err != nil {
					return fmt.Errorf("failed to update stderr.txt: %v", err)
				}
			}
		}
		return nil
	})
}

func main() {
	dirs := []string{
		filepath.Join("internal", "endtoend", "testdata"),
		filepath.Join("examples"),
	}
	for _, d := range dirs {
		if err := regenerate(d); err != nil {
			log.Fatal(err)
		}
	}
}
