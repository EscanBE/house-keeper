package files

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

func Test_listFilesWithinDir(t *testing.T) {
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		dir  string
		want []string
	}{
		{
			dir: "../../constants",
			want: []string{
				"/house-keeper/constants/constants.go",
				"/house-keeper/constants/varcons.go",
			},
		},
		{
			dir: "../../config",
			want: []string{
				"/house-keeper/config/config.go",
				"/house-keeper/config/context.go",
			},
		},
		{
			dir: "../..",
			want: []string{
				"/house-keeper/.gitignore",
				"/house-keeper/go.mod",
				"/house-keeper/go.sum",
				"/house-keeper/Makefile",
				"/house-keeper/README.md",
			},
		},
	}

	curDir, err := os.Getwd()
	if err != nil {
		panic("failed to get current directory")
	}
	fmt.Printf("Current working directory: %s\n", curDir)

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			got := listFilesWithinDir(tt.dir)
			sort.Slice(got, func(i, j int) bool {
				return strings.Compare(got[i], got[j]) < 0
			})

			want := tt.want
			sort.Slice(want, func(i, j int) bool {
				return strings.Compare(want[i], want[j]) < 0
			})

			var failed bool
			if len(got) != len(want) {
				failed = true
			} else {
				for i := 0; i < len(got); i++ {
					if strings.HasSuffix(got[i], want[i]) {
						continue
					}

					fmt.Println(got[i], "does not ends with", want[i])
					failed = true
					break
				}
			}

			if failed {
				fmt.Printf("Fail: Expected %d, got %d\n", len(want), len(got))
				fmt.Println("- Expected:")
				for i, s := range want {
					fmt.Printf("\t+ [%d] %s\n", i, s)
				}
				fmt.Println("- Got:")
				for i, s := range got {
					fmt.Printf("\t+ [%d] %s\n", i, s)
				}
				t.FailNow()
			}
		})
	}
}
