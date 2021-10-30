// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func currentLocation() *Room {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	return &Room{
		path: path,
	}
}

// Room is a dungeon room that represents a filesystem directory.
type Room struct {
	path string
}

func printHeader(str string, out io.Writer) {
	fmt.Fprintf(out, "%s\n", str)
	fmt.Fprintf(out, "%s\n\n", strings.Repeat("=", len(str)))
}

// Description returns a description of the Room.
func (r *Room) Description() string {
	desc := "You stand in a dusty dungeon chamber, not much different from the rest. The\n" +
		"room is full of cobwebs and everything is coated in an undisturbed layer of\n" +
		"dust.\n"

	pathDesc := "On one wall there is a metal plaque rusted with age. It reads:"

	return fmt.Sprintf("%s\n%s\n\t%s\n\n", desc, pathDesc, r.path)
}

// Display writes the information about the Room to the io.Writer.
func (r *Room) Display(out io.Writer) {
	shortName := filepath.Base(r.path)

	printHeader(shortName, out)

	fmt.Fprintf(out, r.Description())

	files := []string{}
	dirs := []string{}

	_ = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		if path == r.path {
			return nil // don't print current dir
		}

		// don't print contents of child directories
		if filepath.Join(r.path, info.Name()) != path {
			return filepath.SkipDir
		}

		if info.IsDir() {
			if os.IsPermission(err) {
				dirs = append(dirs, fmt.Sprintf("%s (locked)", info.Name()))
				return filepath.SkipDir
			}

			if err != nil {
				return filepath.SkipDir
			}

			dirs = append(dirs, info.Name())
		} else {
			files = append(files, info.Name())
		}
		return nil
	})

	if len(files) > 0 {
		fmt.Fprint(out, "Items:\n")
		for _, file := range files {
			fmt.Fprintf(out, "\t%s\n", file)
		}
	}

	fmt.Fprint(out, "Exits:\n")
	for _, dir := range dirs {
		fmt.Fprintf(out, "\t%s\n", dir)
	}
}
