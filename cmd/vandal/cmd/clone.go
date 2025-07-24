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
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.AddCommand(createCloneCmd)
	cloneCmd.AddCommand(getCloneCmd)
	cloneCmd.AddCommand(listCloneCmd)
	cloneCmd.AddCommand(deleteCloneCmd)
	cloneCmd.AddCommand(statusCloneCmd)
	cloneCmd.AddCommand(connectionCloneCmd)
	createCloneCmd.Flags().StringP("filename", "f", "", "Filename of the Clone to create")
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Manage Clones",
}

var createCloneCmd = &cobra.Command{
	Use:   "create -f [filename]",
	Short: "Create a Clone from a YAML file",
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

		var dc vandalv1alpha1.DataClone
		if err := yaml.Unmarshal(yamlFile, &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := c.Create(context.Background(), &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("DataClone %s created\n", dc.Name)
	},
}

var getCloneCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a Clone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dc vandalv1alpha1.DataClone
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Name: %s\n", dc.Name)
		// TODO: Print more details
	},
}

var listCloneCmd = &cobra.Command{
	Use:   "list",
	Short: "List Clones",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var dataClones vandalv1alpha1.DataCloneList
		if err := c.List(context.Background(), &dataClones, &client.ListOptions{}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, dc := range dataClones.Items {
			fmt.Printf("Name: %s\n", dc.Name)
		}
	},
}

var deleteCloneCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a Clone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dc vandalv1alpha1.DataClone
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := c.Delete(context.Background(), &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("DataClone %s deleted\n", name)
	},
}

var statusCloneCmd = &cobra.Command{
	Use:   "status [name]",
	Short: "Get the status of a Clone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dc vandalv1alpha1.DataClone
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Status: %s\n", dc.Status.Phase)
	},
}

var connectionCloneCmd = &cobra.Command{
	Use:   "connection [name]",
	Short: "Get the connection info for a Clone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := client.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name := args[0]
		var dc vandalv1alpha1.DataClone
		if err := c.Get(context.Background(), client.ObjectKey{Name: name}, &dc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if dc.Status.DatabaseConnection != nil {
			fmt.Printf("Host: %s\n", dc.Status.DatabaseConnection.Host)
			fmt.Printf("Port: %d\n", dc.Status.DatabaseConnection.Port)
			fmt.Printf("User: %s\n", dc.Status.DatabaseConnection.User)
			fmt.Printf("Password: %s\n", dc.Status.DatabaseConnection.Password)
		} else {
			fmt.Println("Connection info not available")
		}
	},
}
