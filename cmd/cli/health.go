package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"usual_store/internal/driver"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health check commands",
	Long:  "Check the health of various application components",
}

var healthCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run health checks on all services",
	Long:  "Check the status of database, API server, and other services",
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")

		fmt.Printf("🏥 Running health checks...\n\n")

		// Check database
		checkDatabase()

		if all {
			// Check API server
			checkAPIServer()

			// Check Jaeger (if enabled)
			checkJaeger()
		}

		fmt.Printf("\n✓ Health check completed!\n")
	},
}

var healthStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show application statistics",
	Long:  "Display statistics about users, widgets, orders, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")
		db, err := driver.OpenDB(dsn)
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}
		defer db.Close()

		fmt.Printf("📊 Application Statistics\n\n")

		// Count users
		var userCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
			fmt.Printf("Error getting user count: %v\n", err)
		} else {
			fmt.Printf("👥 Users: %d\n", userCount)
		}

		// Count widgets
		var widgetCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM widgets").Scan(&widgetCount); err != nil {
			fmt.Printf("Error getting widget count: %v\n", err)
		} else {
			fmt.Printf("📦 Widgets: %d\n", widgetCount)
		}

		// Count orders
		var orderCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount); err != nil {
			fmt.Printf("Error getting order count: %v\n", err)
		} else {
			fmt.Printf("🛒 Orders: %d\n", orderCount)
		}

		// Count transactions
		var transactionCount int
		if err := db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&transactionCount); err != nil {
			fmt.Printf("Error getting transaction count: %v\n", err)
		} else {
			fmt.Printf("💳 Transactions: %d\n", transactionCount)
		}

		// Latest order
		var latestOrderTime time.Time
		err = db.QueryRow("SELECT MAX(created_at) FROM orders").Scan(&latestOrderTime)
		if err == nil && !latestOrderTime.IsZero() {
			fmt.Printf("📅 Latest Order: %s\n", latestOrderTime.Format("2006-01-02 15:04:05"))
		}

		fmt.Println()
	},
}

func init() {
	healthCheckCmd.Flags().Bool("all", false, "Check all services (API, Jaeger, etc.)")

	healthCmd.AddCommand(healthCheckCmd)
	healthCmd.AddCommand(healthStatsCmd)
}

func checkDatabase() {
	dsn := getEnvOrDefault("DATABASE_DSN", "postgres://postgres:password@localhost:5433/usualstore?sslmode=disable")

	fmt.Printf("🔍 Checking database connection...\n")
	db, err := driver.OpenDB(dsn)
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("  ❌ Ping failed: %v\n", err)
		return
	}

	fmt.Printf("  ✓ Database is healthy\n")
}

func checkAPIServer() {
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "4001"
	}

	url := fmt.Sprintf("http://localhost:%s/api/widgets", apiPort)
	fmt.Printf("🔍 Checking API server (%s)...\n", url)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ API server is healthy (status: %d)\n", resp.StatusCode)
	} else {
		fmt.Printf("  ⚠️  API server returned status: %d\n", resp.StatusCode)
	}
}

func checkJaeger() {
	if os.Getenv("OTEL_ENABLED") != "true" {
		fmt.Printf("🔍 Jaeger: Skipped (OpenTelemetry not enabled)\n")
		return
	}

	url := "http://localhost:16686"
	fmt.Printf("🔍 Checking Jaeger UI (%s)...\n", url)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("  ❌ Failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("  ✓ Jaeger UI is accessible\n")
	} else {
		fmt.Printf("  ⚠️  Jaeger UI returned status: %d\n", resp.StatusCode)
	}
}
