package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dataprofileCmd)
	dataprofileCmd.AddCommand(createDataProfileCmd)
	dataprofileCmd.AddCommand(getDataProfileCmd)
	dataprofileCmd.AddCommand(listDataProfileCmd)
	dataprofileCmd.AddCommand(deleteDataProfileCmd)
}

var dataprofileCmd = &cobra.Command{
	Use:   "dataprofile",
	Short: "Manage DataProfiles",
}

var createDataProfileCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DataProfile",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement create dataprofile
	},
}

var getDataProfileCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a DataProfile",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement get dataprofile
	},
}

var listDataProfileCmd = &cobra.Command{
	Use:   "list",
	Short: "List DataProfiles",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement list dataprofiles
	},
}

var deleteDataProfileCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DataProfile",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement delete dataprofile
	},
}
