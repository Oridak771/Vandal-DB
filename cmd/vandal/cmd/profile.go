package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	vandalv1alpha1 "github.com/Oridak771/Vandal/apis/v1alpha1"
	"github.com/Oridak771/Vandal/pkg/client"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(createProfileCmd)
	profileCmd.AddCommand(getProfileCmd)
	profileCmd.AddCommand(listProfileCmd)
	profileCmd.AddCommand(deleteProfileCmd)
	profileCmd.AddCommand(statusProfileCmd)
	createProfileCmd.Flags().StringP("filename", "f", "", "Filename of the Profile to create")
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage Profiles",
}

var createProfileCmd = &cobra.Command{
	Use:   "create -f [filename]",
	Short: "Create a Profile from a YAML file",
	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		if filename == "" {
			fmt.Println("Please provide a filename with the -f flag")
			os.Exit(1)
		}

		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dp vandalv1alpha1.DataProfile
		if err := yaml.Unmarshal(yamlFile, &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := c.Create(context.Background(), &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("DataProfile %s created\n", dp.Name)
	},
}

var getProfileCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a Profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dp vandalv1alpha1.DataProfile
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Name: %s\n", dp.Name)
		// TODO: Print more details
	},
}

var listProfileCmd = &cobra.Command{
	Use:   "list",
	Short: "List Profiles",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dataProfiles vandalv1alpha1.DataProfileList
		if err := c.List(context.Background(), &dataProfiles, &client.ListOptions{}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, dp := range dataProfiles.Items {
			fmt.Printf("Name: %s\n", dp.Name)
		}
	},
}

var deleteProfileCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a Profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dp vandalv1alpha1.DataProfile
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := c.Delete(context.Background(), &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("DataProfile %s deleted\n", name)
	},
}

var statusProfileCmd = &cobra.Command{
	Use:   "status [name]",
	Short: "Get the status of a Profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dp vandalv1alpha1.DataProfile
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dp); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Status: %s\n", dp.Status.Phase)
	},
}
