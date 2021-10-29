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

func NewFileExplorer(in io.Reader, out io.Writer) *FileExplorer {
	return &FileExplorer{
		in:  bufio.NewReader(in),
		out: out,
	}
}

type FileExplorer struct {
	in  *bufio.Reader
	out io.Writer
}

func (fe *FileExplorer) Run() error {
	fe.print("Welcome to Filesystem Explorer!\n\n")

	for {
		cur, err := filepath.Abs(".")
		if err != nil {
			return err
		}
		fe.print("%s\n", cur)

		fe.print("Exits:\n")
		if err = filepath.Walk(cur, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == cur {
				return nil // don't print current dir
			}
			// don't print contents of child directories
			if filepath.Join(cur, info.Name()) != path {
				return filepath.SkipDir
			}
			if info.IsDir() {
				fe.print("\t%s\n", info.Name())
			}
			return nil
		}); err != nil {
			return err
		}

		if err := fe.promptCommand(); err != nil {
			return err
		}
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

func (fe *FileExplorer) processCommand(input string) error {
	tokens := strings.Split(input, " ")

	command := tokens[0]
	switch command {
	case "quit", "q":
		fe.print("Goodbye!\n")
		os.Exit(0)
	default:
		return invalidCommandError(command)
	}
	return nil
}

func invalidCommandError(cmd string) error {
	return fmt.Errorf("invalid command %q, try again", cmd)
}

func (fe *FileExplorer) prompt(req string) (string, error) {
	fe.print("%s> ", req)
	resp, err := fe.in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}
