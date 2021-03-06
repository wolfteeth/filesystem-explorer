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

// Room is an interface for all rooms.
type Room interface {
	Name() string
	DisplayName() string
	Description() string
	Exits() []string
	Items() []string
}

// DisplayRoom writes a depiction of the Room to the io.Writer.
func DisplayRoom(r Room, out io.Writer) {
	printHeader(r.Name(), out)

	fmt.Fprintf(out, r.Description())

	items := r.Items()
	if len(items) > 0 {
		fmt.Fprint(out, "Items:\n")
		for _, item := range items {
			fmt.Fprintf(out, "\t%s\n", item)
		}
	}

	fmt.Fprint(out, "Exits:\n")
	for _, exit := range r.Exits() {
		fmt.Fprintf(out, "\t%s\n", exit)
	}
}

func printHeader(str string, out io.Writer) {
	fmt.Fprintf(out, "%s\n", str)
	fmt.Fprintf(out, "%s\n\n", strings.Repeat("=", len(str)))
}

// DirRoom is a dungeon room that represents a filesystem directory.
type DirRoom struct {
	path string
}

// Name returns the path of the directory.
func (r *DirRoom) Name() string {
	return r.path
}

// DisplayName returns a short title to be used for display.
func (r *DirRoom) DisplayName() string {
	return filepath.Base(r.path)
}

// Description returns a description of the Room.
func (r *DirRoom) Description() string {
	desc := "You stand in a dusty dungeon chamber, not much different from the rest. The\n" +
		"room is full of cobwebs and everything is coated in an undisturbed layer of\n" +
		"dust.\n"

	pathDesc := "On one wall there is a metal plaque rusted with age. It reads:"

	return fmt.Sprintf("%s\n%s\n\t%s\n\n", desc, pathDesc, r.path)
}

// Exits returns a list of the child directories as exits.
func (r *DirRoom) Exits() []string {
	var dirs []string

	_ = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		// Do nothing with current dir
		if path == r.path {
			return nil
		}

		// don't recurse into child directories
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
		}
		return nil
	})

	return dirs
}

// Items returns the files in the directory as objects.
func (r *DirRoom) Items() []string {
	var files []string

	_ = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
		// Do nothing with current dir
		if path == r.path {
			return nil
		}

		// don't recurse into child directories
		if filepath.Join(r.path, info.Name()) != path {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return filepath.SkipDir
		} else {
			files = append(files, info.Name())
		}
		return nil
	})

	return files
}
