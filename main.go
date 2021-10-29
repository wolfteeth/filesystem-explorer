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
	"os"
	"path/filepath"
	"strings"
)

func main() {
	out := os.Stdout
	in := bufio.NewScanner(os.Stdin)

	fmt.Fprintf(out, "Welcome to Filesystem Explorer!\n\n")

	for {
		cur, err := filepath.Abs(".")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(out, "%s\n", cur)

		fmt.Fprintf(out, "Exits:\n")
		_ = filepath.Walk(cur, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				fmt.Fprintf(out, "\t%s\n", filepath.Base(path))
			}
			return nil
		})

		fmt.Fprintf(out, "Enter command: ")
		var input strings.Builder
		for in.Scan() {
			input.WriteString(in.Text())
		}

		fmt.Fprintf(out, "got %q\n", input.String())
	}
}
