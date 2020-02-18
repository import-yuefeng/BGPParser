package cmd

import (
	client "github.com/import-yuefeng/BGPParser/cmd/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ip string

func init() {
	addCmd.Flags().StringVarP(&ip, "ip", "i", "", "ip address")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "search",
	Short: "search ip info",
	Long:  `search ip info`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("search ip info")
		if ip != "" {
			c := client.NewClient(":2048")
			c.Search(ip)
		}
	},
}
