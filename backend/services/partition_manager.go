package services

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type PartitionManager struct {
	db *gorm.DB
}

func NewPartitionManager(db *gorm.DB) *PartitionManager {
	return &PartitionManager{db: db}
}

// CreateMonthlyPartitions creates partitions for the next few months
func (pm *PartitionManager) CreateMonthlyPartitions() error {
	log.Println("Creating monthly partitions for stock_prices")

	now := time.Now()
	
	// Create partitions for next 6 months
	for i := 0; i < 6; i++ {
		targetMonth := now.AddDate(0, i, 0)
		err := pm.createPartitionForMonth(targetMonth)
		if err != nil {
			log.Printf("Failed to create partition for %s: %v", targetMonth.Format("2006-01"), err)
		}
	}

	return nil
}

func (pm *PartitionManager) createPartitionForMonth(targetMonth time.Time) error {
	year := targetMonth.Year()
	month := int(targetMonth.Month())
	
	// Calculate start and end dates
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)
	
	partitionName := fmt.Sprintf("stock_prices_%04d_%02d", year, month)
	
	// Check if partition already exists
	var count int64
	err := pm.db.Raw(`
		SELECT count(*) 
		FROM information_schema.tables 
		WHERE table_name = ? AND table_schema = current_schema()
	`, partitionName).Scan(&count).Error
	
	if err != nil {
		return fmt.Errorf("failed to check partition existence: %w", err)
	}
	
	if count > 0 {
		log.Printf("Partition %s already exists, skipping", partitionName)
		return nil
	}
	
	// Create partition
	createSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s PARTITION OF stock_prices
		FOR VALUES FROM ('%s') TO ('%s')
	`, partitionName, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	
	err = pm.db.Exec(createSQL).Error
	if err != nil {
		return fmt.Errorf("failed to create partition %s: %w", partitionName, err)
	}
	
	log.Printf("Successfully created partition: %s", partitionName)
	return nil
}

// CleanupOldPartitions removes partitions older than specified months
func (pm *PartitionManager) CleanupOldPartitions(monthsToKeep int) error {
	log.Printf("Cleaning up partitions older than %d months", monthsToKeep)

	cutoffDate := time.Now().AddDate(0, -monthsToKeep, 0)
	
	// Get list of old partitions
	var partitions []string
	err := pm.db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_name LIKE 'stock_prices_____%%' 
		AND table_schema = current_schema()
	`).Scan(&partitions).Error
	
	if err != nil {
		return fmt.Errorf("failed to get partition list: %w", err)
	}
	
	for _, partition := range partitions {
		// Extract date from partition name (stock_prices_YYYY_MM)
		var year, month int
		_, err := fmt.Sscanf(partition, "stock_prices_%d_%d", &year, &month)
		if err != nil {
			log.Printf("Failed to parse partition name %s: %v", partition, err)
			continue
		}
		
		partitionDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		
		if partitionDate.Before(cutoffDate) {
			err := pm.dropPartition(partition)
			if err != nil {
				log.Printf("Failed to drop partition %s: %v", partition, err)
			} else {
				log.Printf("Successfully dropped old partition: %s", partition)
			}
		}
	}
	
	return nil
}

func (pm *PartitionManager) dropPartition(partitionName string) error {
	// First, detach the partition
	detachSQL := fmt.Sprintf("ALTER TABLE stock_prices DETACH PARTITION %s", partitionName)
	err := pm.db.Exec(detachSQL).Error
	if err != nil {
		return fmt.Errorf("failed to detach partition: %w", err)
	}
	
	// Then drop the table
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", partitionName)
	err = pm.db.Exec(dropSQL).Error
	if err != nil {
		return fmt.Errorf("failed to drop partition table: %w", err)
	}
	
	return nil
}

// GetPartitionInfo returns information about existing partitions
func (pm *PartitionManager) GetPartitionInfo() ([]PartitionInfo, error) {
	var partitions []PartitionInfo
	
	err := pm.db.Raw(`
		SELECT 
			t.table_name as name,
			pg_size_pretty(pg_total_relation_size(t.table_name::regclass)) as size,
			(SELECT count(*) FROM stock_prices WHERE 
				timestamp >= (regexp_split_to_array(t.table_name, '_'))[3] || '-' || 
				(regexp_split_to_array(t.table_name, '_'))[4] || '-01'::timestamp
				AND timestamp < (regexp_split_to_array(t.table_name, '_'))[3] || '-' || 
				(regexp_split_to_array(t.table_name, '_'))[4] || '-01'::timestamp + interval '1 month'
			) as row_count
		FROM information_schema.tables t
		WHERE t.table_name LIKE 'stock_prices_____%%'
		AND t.table_schema = current_schema()
		ORDER BY t.table_name
	`).Scan(&partitions).Error
	
	return partitions, err
}

type PartitionInfo struct {
	Name     string `json:"name"`
	Size     string `json:"size"`
	RowCount int64  `json:"row_count"`
}

// ScheduledMaintenance runs partition maintenance tasks
func (pm *PartitionManager) ScheduledMaintenance() {
	log.Println("Running scheduled partition maintenance")
	
	// Create future partitions
	err := pm.CreateMonthlyPartitions()
	if err != nil {
		log.Printf("Error creating partitions: %v", err)
	}
	
	// Cleanup old partitions (keep 24 months)
	err = pm.CleanupOldPartitions(24)
	if err != nil {
		log.Printf("Error cleaning up partitions: %v", err)
	}
	
	log.Println("Partition maintenance completed")
}