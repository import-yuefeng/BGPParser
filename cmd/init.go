package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)


var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init bgpParser",
	Long:  `init bgpParser`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("init bgpParser")
	},
  }
  
func init() {
	rootCmd.AddCommand(initCmd)
}
