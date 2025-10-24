package git

import (
	"errors"
	"log"
	"os/exec"
)

func Commit(msg string) {
	log.Printf("Running command: git commit -m \"%s\"\n", msg)
	_, err := exec.Command("git", "commit", "-m", msg).Output()
	if err != nil {
		log.Fatalf("failed to execute git commit: %v", err)
	}
}

func Diff() (string, error) {
	log.Println("Executing git diff --cached")
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		log.Println("git diff failed:", err)
		return "", errors.New("failed to execute git diff --cached" + err.Error())
	}
	if string(diff) == "" {
		log.Println("git diff is empty")
		return "", errors.New("no changes to commit, working tree clean")
	}
	log.Println("git diff retrieved, bytes:", len(diff))
	return string(diff), nil
}
