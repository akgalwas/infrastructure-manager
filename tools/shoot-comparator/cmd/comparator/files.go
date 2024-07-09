package comparator

import (
	"fmt"
	"github.com/kyma-project/infrastructure-manager/tools/shoot-comparator/internal/files"
	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:     "files",
	Aliases: []string{"f"},
	Short:   "Compare files",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		firstFile := args[0]
		secondFile := args[1]

		fmt.Printf("Comparing files: %s and %s \n", firstFile, secondFile)
		equal, matcherErrorMessage, err := files.CompareFiles(firstFile, secondFile)
		if err != nil {
			fmt.Printf("Failed to compare files: %s", err.Error())
			return
		}

		if equal {
			fmt.Println("Shoot files are equal")
		} else {
			fmt.Println("Shoot files are NOT equal")
			fmt.Println(matcherErrorMessage)
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
}