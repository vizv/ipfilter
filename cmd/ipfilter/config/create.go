package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagForce bool

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create ipfilter config file.",
	Long:  `Create a new ipfilter config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.SafeWriteConfig()

		if err == nil {
			// new config created
			log.Infof("config created.")
			return
		}

		if !flagForce {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
				// will not update existing config
				log.Errorf("config exists")
			} else {
				// failed to create new config
				log.Errorf("failed to create config: %+v", err)
			}
			return
		}

		if err := viper.WriteConfig(); err == nil {
			// existing config updated
			log.Infof("config updated.")
			return
		}

		// failed to update exiting config
		log.Errorf("failed to update config: %+v", err)
	},
}

func init() {
	CreateCmd.Flags().BoolVarP(&flagForce, "force", "f", false, "Overwrite config file even if already exists. (default: false)")
}
