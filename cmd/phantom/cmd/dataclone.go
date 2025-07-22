package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(datacloneCmd)
	datacloneCmd.AddCommand(createDataCloneCmd)
	datacloneCmd.AddCommand(getDataCloneCmd)
	datacloneCmd.AddCommand(listDataCloneCmd)
	datacloneCmd.AddCommand(deleteDataCloneCmd)
}

var datacloneCmd = &cobra.Command{
	Use:   "dataclone",
	Short: "Manage DataClones",
}

var createDataCloneCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DataClone",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement create dataclone
	},
}

var getDataCloneCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a DataClone",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement get dataclone
	},
}

var listDataCloneCmd = &cobra.Command{
	Use:   "list",
	Short: "List DataClones",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement list dataclones
	},
}

var deleteDataCloneCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a DataClone",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement delete dataclone
	},
}
