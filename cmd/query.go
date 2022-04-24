package cmd

import (
	"fmt"
	"log"
	"time"

	remote "github.com/bartmika/bleve-server/pkg/rpc_client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	queryFilename string
	querySearch   string
)

func init() {
	queryCmd.Flags().StringVarP(&queryFilename, "filename", "a", "1", "")
	queryCmd.Flags().StringVarP(&querySearch, "search", "b", "Some random data", "")
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Perform full-text search of indexed data",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: doQuery,
}

func doQuery(cmd *cobra.Command, args []string) {
	addr := viper.GetString("appAddress")
	s := remote.New(addr, 3, 15*time.Second)
	uuids, err := s.Query(queryFilename, querySearch)
	if err != nil {
		log.Fatal("doQuery err:", err)
	}
	fmt.Println("UUIDs:", uuids)
}
