package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "starkabi --abi",
		Short: "Generate Go bindings for Cairo contracts",
		Long: `Generate Go bindings for Cairo contracts
								Creating now random content just to test
								how this text will be displayed in the termina.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			ticker := time.NewTicker(1 * time.Second)
			go func() {
				for {
					<-ticker.C
					fmt.Print(".")
				}
			}()
			select {}
		},
	}
)

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func main() {
	rootCmd.Execute()
}
