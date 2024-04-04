package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vizv/ipfilter/cmd/ipfilter"
)

var RootCmd = &cobra.Command{
	Use:   "ipfilter",
	Short: "qBittorrent ipfilter.dat utilities",
	Long: `ipfilter can help you merge ipfilter.dat files or synchronize ipfilter.dat from remote.

TODO: example usage.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("global.verbose") {
			log.SetLevel(log.DebugLevel)
			log.Debugf("verbose mode enabled")
		}
		if viper.GetBool("global.debug") {
			log.SetLevel(log.TraceLevel)
			log.Tracef("debug mode enabled")
		}
	},
}

func init() {
	viper.SetConfigName("ipfilter")
	viper.SetConfigType("ini")
	viper.AddConfigPath(".")
	if _, err := os.Stat("ipfilter.ini"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			log.Errorf("failed to read config: %+v", err)
		}
	}

	RootCmd.PersistentFlags().BoolVarP(&ipfilter.Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("global.verbose", RootCmd.PersistentFlags().Lookup("verbose"))

	RootCmd.PersistentFlags().BoolVarP(&ipfilter.Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("global.debug", RootCmd.PersistentFlags().Lookup("debug"))

	// RootCmd.AddCommand(ipfilter.ConfigCmd)
	RootCmd.AddCommand(ipfilter.MergeCmd)
	RootCmd.AddCommand(ipfilter.SyncCmd)
}
