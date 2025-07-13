/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/joaqu1m/goexpert-labs/modules/16/internal/repository"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new category",
	Long:  "Create a new category with the specified name and description",
	Run: func(cmd *cobra.Command, args []string) {

		db, err := GetDB()
		if err != nil {
			cmd.Println("Error connecting to the database:", err)
			return
		}
		defer db.Close()

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			cmd.Println("Error: Name flag is required")
			return
		}
		cmd.Printf("Creating category with name: %s\n", name)
		description, _ := cmd.Flags().GetString("description")
		if description == "" {
			cmd.Println("Error: Description flag is required")
			return
		}
		cmd.Printf("Description: %s\n", description)

		repo := repository.NewCategoryRepository(db)
		_, err = repo.CreateCategory(name, &description)
		if err != nil {
			cmd.Println("Error creating category:", err)
			return
		}
		cmd.Println("Category created successfully")
	},
}

func init() {
	categoryCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "Name of the category to create")
	createCmd.Flags().StringP("description", "d", "", "Description of the category")

	createCmd.MarkFlagsRequiredTogether("name", "description")
}
