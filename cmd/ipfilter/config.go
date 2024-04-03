package ipfilter

import (
	"github.com/spf13/cobra"

	"github.com/vizv/ipfilter/cmd/ipfilter/config"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage ipfilter config file.",
	Long:  `Create config file, or get/set config values.`,
}

func init() {
	ConfigCmd.AddCommand(config.CreateCmd)
}
