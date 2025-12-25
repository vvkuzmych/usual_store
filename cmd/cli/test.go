package main

import (
	"database/sql"
	"fmt"
	"log"
	"usual_store/internal/driver"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Testing utilities",
	Long:  "Commands for testing and development",
}

var testGenerateDataCmd = &cobra.Command{
	Use:   "generate-data",
	Short: "Generate comprehensive test data",
	Long:  "Generate a complete set of test data including users, widgets, and orders",
	Run: func(cmd *cobra.Command, args []string) {
		scenario, _ := cmd.Flags().GetString("scenario")

		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		fmt.Printf("🧪 Generating test data for scenario: %s\n\n", scenario)

		switch scenario {
		case "minimal":
			generateMinimalData(db)
		case "checkout":
			generateCheckoutScenario(db)
		case "subscriptions":
			generateSubscriptionScenario(db)
		default:
			generateMinimalData(db)
		}

		fmt.Printf("\n✓ Test data generation completed!\n")
	},
}

var testResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset database to initial state",
	Long:  "Clear all data and restore to initial seed state",
	Run: func(cmd *cobra.Command, args []string) {
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			log.Fatal("⚠️  This will reset the entire database! Use --confirm to proceed")
		}

		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		fmt.Printf("🔄 Resetting database...\n\n")

		// Clear all tables (except keep admin user)
		tables := []string{"orders", "transactions", "tokens", "customers"}
		for _, table := range tables {
			query := fmt.Sprintf("DELETE FROM %s", table)
			_, err := db.Exec(query)
			if err != nil {
				log.Printf("Warning: Failed to clear %s: %v", table, err)
			} else {
				fmt.Printf("✓ Cleared %s\n", table)
			}
		}

		// Reset users (keep only admin)
		_, err = db.Exec("DELETE FROM users WHERE email != 'admin@example.com'")
		if err != nil {
			log.Printf("Warning: Failed to reset users: %v", err)
		} else {
			fmt.Printf("✓ Reset users (kept admin)\n")
		}

		fmt.Printf("\n✓ Database reset completed!\n")
	},
}

func init() {
	// Generate data flags
	testGenerateDataCmd.Flags().String("scenario", "minimal", "Test scenario: minimal, checkout, subscriptions")

	// Reset flags
	testResetCmd.Flags().Bool("confirm", false, "Confirm database reset (required)")

	// Add subcommands
	testCmd.AddCommand(testGenerateDataCmd)
	testCmd.AddCommand(testResetCmd)
}

// Helper functions for different test scenarios
func generateMinimalData(db *sql.DB) {
	fmt.Println("Creating minimal test data...")

	// Create a test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass123"), bcrypt.DefaultCost)
	query := `INSERT INTO users (first_name, last_name, email, password) 
			  VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING RETURNING id`
	var userID int
	err := db.QueryRow(query, "Test", "User", "test@example.com", string(hashedPassword)).Scan(&userID)
	if err != nil {
		log.Printf("Test user might already exist")
	} else {
		fmt.Printf("✓ Created test user (test@example.com / testpass123)\n")
	}

	// Create test widgets
	widgets := []struct {
		name  string
		price int
	}{
		{"Test Product 1", 999},
		{"Test Product 2", 1999},
		{"Test Subscription", 2999},
	}

	for _, w := range widgets {
		query := `INSERT INTO widgets (name, description, price, inventory_level, image) 
				  VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Exec(query, w.name, "Test description", w.price, 100, "/static/widget.png")
		if err != nil {
			log.Printf("Warning: Failed to create widget: %v", err)
		}
	}
	fmt.Printf("✓ Created %d test widgets\n", len(widgets))
}

func generateCheckoutScenario(db *sql.DB) {
	fmt.Println("Creating checkout scenario test data...")
	generateMinimalData(db)
	fmt.Println("✓ Checkout scenario ready - use test@example.com / testpass123 to login")
}

func generateSubscriptionScenario(db *sql.DB) {
	fmt.Println("Creating subscription scenario test data...")
	generateMinimalData(db)

	// Create subscription products
	query := `INSERT INTO widgets (name, description, price, inventory_level, image, is_recurring, plan_id) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	db.Exec(query, "Monthly Plan", "Test monthly subscription", 999, 999, "/static/plan.png", true, "price_test_monthly")
	fmt.Println("✓ Subscription scenario ready")
}
