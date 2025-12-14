package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Category represents the categories table structure for PostgreSQL
// Note: If your MySQL table has different column types (e.g., INT for ID, VARCHAR for name),
// you may need to create a separate MySQLCategory struct and map it to this Category struct
type Category struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:uuidv7()" json:"id"`
	Name        string     `gorm:"type:text;not null" json:"name"`
	CategoryType string    `gorm:"type:category_type;not null" json:"category_type"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	UserUUID    string     `gorm:"type:uuid;not null" json:"user_uuid"`
	Description *string    `gorm:"type:text" json:"description,omitempty"`
	Hidden      bool       `gorm:"not null;default:false" json:"hidden"`
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}

func main() {
	// Get database connection strings from environment variables
	mysqlDSN := os.Getenv("MYSQL_DSN")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	if mysqlDSN == "" {
		log.Fatal("MYSQL_DSN environment variable is required")
	}
	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN environment variable is required")
	}

	// Connect to MySQL (source)
	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Connect to PostgreSQL (destination)
	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Test connections
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		log.Fatalf("Failed to get MySQL database instance: %v", err)
	}
	defer sqlDB.Close()

	pgSQLDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get PostgreSQL database instance: %v", err)
	}
	defer pgSQLDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}
	fmt.Println("✓ Connected to MySQL")

	if err := pgSQLDB.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	fmt.Println("✓ Connected to PostgreSQL")

	// Migrate categories from MySQL to PostgreSQL
	if err := migrateCategories(mysqlDB, postgresDB); err != nil {
		log.Fatalf("Failed to migrate categories: %v", err)
	}

	fmt.Println("✓ Migration completed successfully")
}

func migrateCategories(mysqlDB, postgresDB *gorm.DB) error {
	var categories []Category

	// Read all categories from MySQL
	// Note: If your MySQL table has different column names or types, you may need to:
	// 1. Create a MySQLCategory struct with MySQL-compatible types
	// 2. Read into MySQLCategory
	// 3. Map MySQLCategory to Category
	result := mysqlDB.Table("categories").Find(&categories)
	if result.Error != nil {
		return fmt.Errorf("failed to read categories from MySQL: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		fmt.Println("No categories found in MySQL database")
		return nil
	}

	fmt.Printf("Found %d categories to migrate\n", len(categories))

	successCount := 0
	skipCount := 0

	// Insert categories into PostgreSQL
	for i, category := range categories {
		// Ensure category_type is valid (Income, Expense, or Transfer)
		if category.CategoryType != "Income" && category.CategoryType != "Expense" && category.CategoryType != "Transfer" {
			log.Printf("Warning: Skipping category %d (ID: %s) with invalid category_type: %s", i+1, category.ID, category.CategoryType)
			skipCount++
			continue
		}

		// Insert into PostgreSQL
		if err := postgresDB.Create(&category).Error; err != nil {
			// If it's a duplicate key error, skip it
			if isDuplicateKeyError(err) {
				log.Printf("Category %d (ID: %s) already exists in PostgreSQL, skipping...", i+1, category.ID)
				skipCount++
				continue
			}
			return fmt.Errorf("failed to insert category %d (ID: %s): %w", i+1, category.ID, err)
		}

		successCount++
		if (i+1)%100 == 0 {
			fmt.Printf("Migrated %d/%d categories...\n", i+1, len(categories))
		}
	}

	fmt.Printf("Migration summary: %d successful, %d skipped, %d total\n", successCount, skipCount, len(categories))
	return nil
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	// Check for common PostgreSQL duplicate key error patterns
	return strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "violates unique constraint")
}
