package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const defaultNotesDir = "/Users/kristakoch/Documents/notes"
const defaultRotateDir = "weekly"

func main() {
	var err error

	var notesDir string
	var rotateDir string

	flag.StringVar(&notesDir, "notesdir", defaultNotesDir, "root notes directory")
	flag.StringVar(&rotateDir, "rotatedir", defaultRotateDir, "directory within root to rotate notes to")
	flag.Parse()

	// Base current week values off of the current week's sunday,
	// which we can get by subtracting the weekday from the current
	// date. Weekdays range from 0-6, starting with sunday.
	n := time.Now()
	cwsun := n.AddDate(0, 0, -int(n.Weekday()))
	lwsun := cwsun.AddDate(0, 0, -7)

	// Rotate last week's note.
	{
		lwfn := buildFilenameForDate(lwsun)
		lwloc := fmt.Sprintf("%s/%s", notesDir, lwfn)

		var exists bool
		if exists, err = fileExists(lwloc); nil != err {
			log.Fatal(err)
		}

		if exists {
			fmt.Printf("found last week note at %s, rotating into %s \n", lwfn, rotateDir)

			cmd := exec.Command("mv", lwloc, fmt.Sprintf("%s/%s/%s", notesDir, rotateDir, lwfn))
			if err = cmd.Run(); nil != err {
				log.Fatal(err)
			}
		}
	}

	// Ensure this week's note exists.
	{
		cwfn := buildFilenameForDate(cwsun)
		cwloc := fmt.Sprintf("%s/%s", notesDir, cwfn)

		var exists bool
		if exists, err = fileExists(cwloc); nil != err {
			log.Fatal(err)
		}

		if !exists {
			fmt.Printf("no current week note found at %s, creating \n", cwfn)

			if _, err = os.Create(cwloc); nil != err {
				log.Fatal(err)
			}

			initialLines := fmt.Sprintf(
				"# %s \n### monday \n### tuesday \n### wednesday \n### thursday \n ### friday \n### weekend",
				strings.TrimSuffix(cwfn, ".md"),
			)

			// Write the first line.
			if err = os.WriteFile(
				cwloc,
				[]byte(initialLines),
				os.ModeAppend,
			); nil != err {
				log.Fatal(err)
			}
		}
	}
}

func fileExists(name string) (bool, error) {
	f, err := os.Stat(name)
	if nil != err {
		if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}
	}
	return nil != f, nil
}

func buildFilenameForDate(t time.Time) string {
	y := fmt.Sprintf("%d", t.Year())[2:]
	m := fmt.Sprintf("%02d", t.Month())
	d := fmt.Sprintf("%02d", t.Day())

	return fmt.Sprintf("%s.%s.%s.md", y, m, d)
}
