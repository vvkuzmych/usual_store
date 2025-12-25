package main

import (
	"fmt"
	"log"
	"usual_store/internal/driver"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  "Manage database operations for the Usual Store application",
}

var dbStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check database connection status",
	Long:  "Verify that the application can connect to the database",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")

		fmt.Printf("Connecting to database...\n")
		fmt.Printf("DSN: %s\n\n", maskDSN(dsn))

		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("❌ Failed to connect: %v\n", err)
		}
		defer db.Close()

		// Ping database
		if err := db.Ping(); err != nil {
			log.Fatalf("❌ Database ping failed: %v\n", err)
		}

		fmt.Printf("✓ Database connection successful!\n\n")

		// Get database stats
		stats := db.Stats()
		fmt.Printf("📊 Database Statistics:\n")
		fmt.Printf("  Open Connections: %d\n", stats.OpenConnections)
		fmt.Printf("  In Use: %d\n", stats.InUse)
		fmt.Printf("  Idle: %d\n", stats.Idle)
		fmt.Printf("  Wait Count: %d\n", stats.WaitCount)
		fmt.Printf("  Max Open Connections: %d\n", stats.MaxOpenConnections)
	},
}

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with test data",
	Long:  "Populate the database with sample data for testing",
	Run: func(cmd *cobra.Command, args []string) {
		users, _ := cmd.Flags().GetInt("users")
		widgets, _ := cmd.Flags().GetInt("widgets")

		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		fmt.Printf("🌱 Seeding database with test data...\n\n")

		if users > 0 {
			fmt.Printf("Creating %d test users...\n", users)
			for i := 1; i <= users; i++ {
				query := `INSERT INTO users (first_name, last_name, email, password) 
						  VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING`
				_, err := db.Exec(query,
					fmt.Sprintf("User%d", i),
					fmt.Sprintf("Test%d", i),
					fmt.Sprintf("user%d@test.com", i),
					"$2a$12$VR1wDmweaF3ZTVgEHiJrNOSi8VcS4j0eamr96A/7iOe8vlum3O3/q", // "qwerty"
				)
				if err != nil {
					log.Printf("Warning: Failed to create user %d: %v", i, err)
				}
			}
			fmt.Printf("✓ %d users created\n", users)
		}

		if widgets > 0 {
			fmt.Printf("Creating %d test widgets...\n", widgets)
			for i := 1; i <= widgets; i++ {
				query := `INSERT INTO widgets (name, description, price, inventory_level, image) 
						  VALUES ($1, $2, $3, $4, $5)`
				_, err := db.Exec(query,
					fmt.Sprintf("Test Widget %d", i),
					fmt.Sprintf("Description for test widget %d", i),
					1000+i*100, // Price in cents
					10+i,       // Inventory
					"/static/widget.png",
				)
				if err != nil {
					log.Printf("Warning: Failed to create widget %d: %v", i, err)
				}
			}
			fmt.Printf("✓ %d widgets created\n", widgets)
		}

		fmt.Printf("\n✓ Database seeding completed!\n")
	},
}

var dbClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear test data from database",
	Long:  "Remove all test data (use with caution!)",
	Run: func(cmd *cobra.Command, args []string) {
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			log.Fatal("⚠️  This will delete data! Please confirm with --confirm flag")
		}

		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		fmt.Printf("🗑️  Clearing test data...\n\n")

		// Clear test users (keep admin)
		query := `DELETE FROM users WHERE email != 'admin@example.com'`
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Warning: Failed to clear users: %v", err)
		} else {
			rows, _ := result.RowsAffected()
			fmt.Printf("✓ Removed %d test users\n", rows)
		}

		// Clear test widgets (optional - keep if you want)
		// Uncomment if you want to clear widgets too
		/*
			query = `DELETE FROM widgets WHERE name LIKE 'Test Widget%'`
			result, err = db.Exec(query)
			if err != nil {
				log.Printf("Warning: Failed to clear widgets: %v", err)
			} else {
				rows, _ := result.RowsAffected()
				fmt.Printf("✓ Removed %d test widgets\n", rows)
			}
		*/

		fmt.Printf("\n✓ Test data cleared!\n")
	},
}

func init() {
	// Seed command flags
	dbSeedCmd.Flags().Int("users", 10, "Number of test users to create")
	dbSeedCmd.Flags().Int("widgets", 5, "Number of test widgets to create")

	// Clear command flags
	dbClearCmd.Flags().Bool("confirm", false, "Confirm data deletion (required)")

	// Add subcommands
	dbCmd.AddCommand(dbStatusCmd)
	dbCmd.AddCommand(dbSeedCmd)
	dbCmd.AddCommand(dbClearCmd)
}

// maskDSN masks sensitive information in the DSN
func maskDSN(dsn string) string {
	// Simple masking - in production, use a proper URL parser
	return "postgres://postgres:****@localhost:5433/usualstore?sslmode=disable"
}
