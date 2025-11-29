package migrate

import (
	"embed"
	"fmt"
	"log"
	"sort"

	"gorm.io/gorm"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// AutoMigrate runs database migrations in order
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto-migration...")

	// Check if database needs initialization by checking for master_users table
	var tableExists bool
	err := db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'master_users')").Scan(&tableExists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database is initialized: %w", err)
	}

	if tableExists {
		log.Println("Database already initialized, skipping migration")
		return nil
	}

	log.Println("Database not initialized, running schema creation...")

	// Read and execute migration files in order
	migrations := []struct {
		name string
		path string
	}{
		{"schema", "sql/db.sql"},
		{"admin_user", "sql/init_admin_user.sql"},
	}

	for _, migration := range migrations {
		log.Printf("Running migration: %s", migration.name)

		content, err := sqlFiles.ReadFile(migration.path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration.path, err)
		}

		// Execute the SQL file
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.name, err)
		}

		log.Printf("Migration %s completed successfully", migration.name)
	}

	log.Println("Database auto-migration completed successfully")
	return nil
}

// RunMigrations executes all SQL migrations in the sql directory
// This is a more flexible version that automatically discovers and runs all .sql files
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migrations...")

	// Read all SQL files from the embedded directory
	entries, err := sqlFiles.ReadDir("sql")
	if err != nil {
		return fmt.Errorf("failed to read sql directory: %w", err)
	}

	// Sort files alphabetically to ensure consistent execution order
	var sqlFileNames []string
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".sql" {
			sqlFileNames = append(sqlFileNames, entry.Name())
		}
	}
	sort.Strings(sqlFileNames)

	// Execute each SQL file
	for _, fileName := range sqlFileNames {
		log.Printf("Executing migration: %s", fileName)

		content, err := sqlFiles.ReadFile("sql/" + fileName)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		log.Printf("Migration %s completed successfully", fileName)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// CheckDatabaseHealth verifies database connection and basic functionality
func CheckDatabaseHealth(db *gorm.DB) error {
	// Check database connection
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database health check passed")
	return nil
}
