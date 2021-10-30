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
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NewFileExplorer creates a new FileExplorer session with a given input and
// output stream.
func NewFileExplorer(in io.Reader, out io.Writer) *FileExplorer {
	return &FileExplorer{
		in:  bufio.NewReader(in),
		out: out,
	}
}

// FileExplorer controls an interactive filesystem exploration session.
type FileExplorer struct {
	in  *bufio.Reader
	out io.Writer
}

// Run runs the interactive session loop.
func (fe *FileExplorer) Run() error {
	fe.print("Welcome to Filesystem Explorer!\n\n")

	var prev Room
	for {
		room := currentLocation()
		if prev == nil || room.Name() != prev.Name() {
			prev = room
			DisplayRoom(room, fe.out)
		}

		if err := fe.promptCommand(); err != nil {
			return err
		}
	}
}

func currentLocation() *DirRoom {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	return &DirRoom{
		path: path,
	}
}

func errIsRetryable(err error) bool {
	return err != nil && strings.Contains(err.Error(), "try again")
}

func (fe *FileExplorer) print(format string, args ...interface{}) {
	fmt.Fprintf(fe.out, format, args...)
}

func (fe *FileExplorer) promptCommand() error {
	for {
		cmd, err := fe.prompt("Enter command")
		if err != nil {
			return err
		}
		if cmd == "" {
			continue
		}

		err = fe.processCommand(cmd)
		if err == nil {
			break
		}
		if errIsRetryable(err) {
			fe.print("%s\n", err.Error())
			continue
		}

		return err
	}
	return nil
}

func (fe *FileExplorer) prompt(req string) (string, error) {
	fe.print("%s> ", req)
	resp, err := fe.in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

// processCommand parses the requested command and either executes it or returns
// an error.
func (fe *FileExplorer) processCommand(input string) error {
	tokens := strings.Split(input, " ")

	command := tokens[0]
	switch command {
	case "help", "h":
		fe.printHelp()
	case "quit", "q":
		fe.print("Goodbye!\n")
		os.Exit(0)
	case "go", "g":
		if len(tokens) == 1 {
			return tryAgainError("Where do you want to go?")
		}
		dest := strings.Join(tokens[1:], " ")
		if err := os.Chdir(filepath.Join(".", dest)); err != nil {
			if os.IsNotExist(err) {
				return tryAgainError("There is no door to %q here", dest)
			}
			if os.IsPermission(err) {
				return tryAgainError("The door is tightly bolted")
			}
			return tryAgainError("You can't go to %q from here", dest)
		}
	default:
		return tryAgainError("%q is not a valid command", command)
	}
	return nil
}

func (fe *FileExplorer) printHelp() {
	fe.print("Commands:\n")
	fe.print("  [h]elp\t\tDisplay this help screen\n")
	fe.print("  [q]uit\t\tQuit the program\n")
	fe.print("  [g]o [exit-name]\tGo through an exit from a room\n")
	fe.print("\n")
}

func tryAgainError(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s (try again)", msg)
}
