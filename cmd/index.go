package cmd

import (
	"log"
	"time"

	remote "github.com/bartmika/bleve-server/pkg/rpc_client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	indexFilename   string
	indexIdentifier string
	indexData       string
)

func init() {
	indexCmd.Flags().StringVarP(&indexFilename, "filename", "a", "1", "")
	indexCmd.Flags().StringVarP(&indexIdentifier, "identifier", "b", "id-123", "")
	indexCmd.Flags().StringVarP(&indexData, "data", "c", "Some random data", "")
	rootCmd.AddCommand(indexCmd)
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Submit data to index",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: doIndex,
}

func doIndex(cmd *cobra.Command, args []string) {
	addr := viper.GetString("appAddress")
	s := remote.New(addr, 3, 15*time.Second)
	err := s.Index(indexFilename, indexIdentifier, []byte(indexData))
	if err != nil {
		log.Fatal("doIndex err:", err)
	}
	log.Println("Index Submitted")
}
