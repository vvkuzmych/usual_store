package main

import (
	"fmt"
	"log"
	"os"
	"usual_store/internal/driver"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  "Manage users in the Usual Store application",
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  "Create a new user with the specified details",
	Run: func(cmd *cobra.Command, args []string) {
		firstName, _ := cmd.Flags().GetString("first-name")
		lastName, _ := cmd.Flags().GetString("last-name")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")
		isAdmin, _ := cmd.Flags().GetBool("admin")

		if firstName == "" || lastName == "" || email == "" || password == "" {
			log.Fatal("All fields (first-name, last-name, email, password) are required")
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}

		// Connect to database
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		// Insert user
		query := `INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
		var userID int
		err = db.QueryRow(query, firstName, lastName, email, string(hashedPassword)).Scan(&userID)
		if err != nil {
			log.Fatalf("Error creating user: %v", err)
		}

		fmt.Printf("✓ User created successfully!\n")
		fmt.Printf("  ID: %d\n", userID)
		fmt.Printf("  Name: %s %s\n", firstName, lastName)
		fmt.Printf("  Email: %s\n", email)
		if isAdmin {
			fmt.Printf("  Role: Admin\n")
		}
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  "Display a list of all users in the system",
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to database
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		// Query users
		query := `SELECT id, first_name, last_name, email, created_at FROM users ORDER BY created_at DESC`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatalf("Error fetching users: %v", err)
		}
		defer rows.Close()

		var users []struct {
			ID        int
			FirstName string
			LastName  string
			Email     string
			CreatedAt string
		}

		for rows.Next() {
			var user struct {
				ID        int
				FirstName string
				LastName  string
				Email     string
				CreatedAt string
			}
			err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)
			if err != nil {
				log.Printf("Error scanning user: %v", err)
				continue
			}
			users = append(users, user)
		}

		fmt.Printf("\n📋 Total Users: %d\n\n", len(users))
		fmt.Printf("%-5s %-20s %-30s %-20s\n", "ID", "Name", "Email", "Created")
		fmt.Println("─────────────────────────────────────────────────────────────────────────────")

		for _, user := range users {
			fmt.Printf("%-5d %-20s %-30s %-20s\n",
				user.ID,
				user.FirstName+" "+user.LastName,
				user.Email,
				user.CreatedAt,
			)
		}
		fmt.Println()
	},
}

var userResetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset a user's password",
	Long:  "Reset the password for a specific user by email",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		if email == "" || password == "" {
			log.Fatal("Both email and password are required")
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}

		// Connect to database
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		// Update password
		query := `UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE email = $2`
		result, err := db.Exec(query, string(hashedPassword), email)
		if err != nil {
			log.Fatalf("Error resetting password: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			log.Fatalf("User with email %s not found", email)
		}

		fmt.Printf("✓ Password reset successfully for %s\n", email)
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user",
	Long:  "Delete a user by email",
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		confirm, _ := cmd.Flags().GetBool("confirm")

		if email == "" {
			log.Fatal("Email is required")
		}

		if !confirm {
			log.Fatal("Please confirm deletion with --confirm flag")
		}

		// Connect to database
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		// Delete user
		query := `DELETE FROM users WHERE email = $1`
		result, err := db.Exec(query, email)
		if err != nil {
			log.Fatalf("Error deleting user: %v", err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			log.Fatalf("User with email %s not found", email)
		}

		fmt.Printf("✓ User %s deleted successfully\n", email)
	},
}

func init() {
	// Create command
	userCreateCmd.Flags().String("first-name", "", "User's first name (required)")
	userCreateCmd.Flags().String("last-name", "", "User's last name (required)")
	userCreateCmd.Flags().String("email", "", "User's email address (required)")
	userCreateCmd.Flags().String("password", "", "User's password (required)")
	userCreateCmd.Flags().Bool("admin", false, "Create as admin user")

	// Reset password command
	userResetPasswordCmd.Flags().String("email", "", "User's email address (required)")
	userResetPasswordCmd.Flags().String("password", "", "New password (required)")

	// Delete command
	userDeleteCmd.Flags().String("email", "", "User's email address (required)")
	userDeleteCmd.Flags().Bool("confirm", false, "Confirm deletion (required)")

	// Add subcommands
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userResetPasswordCmd)
	userCmd.AddCommand(userDeleteCmd)
}

// Helper function
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
