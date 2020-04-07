package cmd

import (
	client "github.com/import-yuefeng/BGPParser/cmd/client"
	"github.com/spf13/cobra"
	// log "github.com/sirupsen/logrus"
)

var (
	raw  string
	data string
	WC   int
)

func init() {
	runCmd.Flags().StringVarP(&raw, "rawFilePath", "r", "", "bgp raw data path")
	runCmd.Flags().StringVarP(&data, "filePath", "d", "", "bgp data path")
	runCmd.Flags().IntVarP(&WC, "parserWC", "w", 1, "parse worker number")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "parse",
	Short: "parser bgpData",
	Long:  `parser bgpData`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewClient(":2048")
		if raw != "" {
			c.AddRawParse(raw)
		} else if data != "" {
			c.AddBGPParse(data)
		}
	},
}
