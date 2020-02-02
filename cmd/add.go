package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)


var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add bgpData source",
	Long:  `add bgpData source`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("add bgpData source")
	},
  }

func init() {
	rootCmd.AddCommand(addCmd)
}
