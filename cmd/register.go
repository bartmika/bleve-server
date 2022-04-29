package cmd

import (
	"fmt"

	"github.com/blevesearch/bleve/v2"
	"github.com/spf13/cobra"
)

var (
	filename string
)

func init() {
	registerCmd.Flags().StringVarP(&filename, "filename", "f", "", "The filename of where to save the bleve index.")
	registerCmd.MarkFlagRequired("filename")
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a bleve index",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Save everything into the db folder.
		filepath := applicationHomeDirectoryPath + "/" + filename

		// open a new index
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(filepath, mapping)
		if err != nil {
			index, err = bleve.Open(filepath)
			fmt.Println("Already registered, aborting now.")
		} else {
			fmt.Println("Registering:", filename)
		}
		index.Close()
	},
}
