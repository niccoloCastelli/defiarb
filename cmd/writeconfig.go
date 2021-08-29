package cmd

import (
	"github.com/niccoloCastelli/defiarb/config"
	"github.com/spf13/cobra"
)

var writeconfigCmd = &cobra.Command{
	Use:   "writeconfig",
	Short: "Write example config",
	Long:  `Write example config file. writeconfig [FILE_PATH (default ./config.json)]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			args = []string{config.GetDefaultConfigLocation()}
		}
		for _, arg := range args {
			if err := config.WriteConfig(arg); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeconfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// writeconfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// writeconfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
