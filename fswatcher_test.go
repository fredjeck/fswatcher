package fswatcher

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestSunnyDay(t *testing.T) {

	w, e := NewFsWatcher(".", 500)
	if e != nil {
		t.Fatalf("%v", e)
	}
	w.RegisterFileExtension(".go", func(path string, info os.FileInfo) bool {
		log.Printf("Go file changed %v", path)
		return false
	})

	w.RegisterFileExtension("*", func(path string, info os.FileInfo) bool {
		log.Printf("* handler called for %v", path)
		if strings.Contains(path, ".git") {
			log.Printf("Skipping git directory")
			return true
		}
		return false
	})
	//blocker := make(chan bool)
	//	w.Watch()
	//	<-blocker
}
