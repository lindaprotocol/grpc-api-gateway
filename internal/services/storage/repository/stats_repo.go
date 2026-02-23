package repository

import (
	"encoding/json"
	"time"

	
	"github.com/lindaprotocol/grpc-api-gateway/pkg/lindapb"
	"encoding/json"
	"github.com/lindaprotocol/grpc-api-gateway/internal/models"
	"gorm.io/gorm"
)

// StatsRepository struct: Repository for stats operations
type StatsRepository struct {
	db *gorm.DB
}

// NewStatsRepository function: Creates a new stats repository
func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// SaveStatistic function: Saves a statistic
func (r *StatsRepository) SaveStatistic(statType string, value interface{}, timestamp int64) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	stat := &models.Statistic{
		Type:      statType,
		Value:     data,
		Timestamp: timestamp,
	}

	return r.db.Save(stat).Error
}

// GetStatistic function: Retrieves the latest statistic of a type
func (r *StatsRepository) GetStatistic(statType string) (json.RawMessage, error) {
	var stat models.Statistic
	err := r.db.Where("type = ?", statType).
		Order("timestamp DESC").
		First(&stat).Error
	if err != nil {
		return nil, err
	}
	return stat.Value, nil
}

// GetStatisticsInRange function: Retrieves statistics in a time range
func (r *StatsRepository) GetStatisticsInRange(statType string, fromTime, toTime int64) ([]json.RawMessage, error) {
	var stats []models.Statistic
	err := r.db.Where("type = ? AND timestamp >= ? AND timestamp <= ?",
		statType, fromTime, toTime).
		Order("timestamp ASC").
		Find(&stats).Error
	if err != nil {
		return nil, err
	}

	var results []json.RawMessage
	for _, stat := range stats {
		results = append(results, stat.Value)
	}
	return results, nil
}

// GetHomepageBundle function: Retrieves homepage bundle data
func (r *StatsRepository) GetHomepageBundle() (*models.HomepageBundleResponse, error) {
	var bundle models.HomepageBundleResponse
	
	// Get total blocks
	r.db.Raw("SELECT COUNT(*) FROM blocks").Scan(&bundle.TotalBlocks)
	
	// Get total transactions
	r.db.Raw("SELECT COUNT(*) FROM transactions").Scan(&bundle.TotalTransactions)
	
	// Get total accounts
	r.db.Raw("SELECT COUNT(*) FROM accounts").Scan(&bundle.TotalAccounts)
	
	// Get recent blocks
	r.db.Raw("SELECT * FROM blocks ORDER BY number DESC LIMIT 10").Scan(&bundle.RecentBlocks)
	
	// Get recent transactions
	r.db.Raw("SELECT * FROM transactions ORDER BY block_timestamp DESC LIMIT 10").Scan(&bundle.RecentTransactions)
	
	return &bundle, nil
}

// GetOverview function: Retrieves overview statistics
func (r *StatsRepository) GetOverview(statType string) (json.RawMessage, error) {
	return r.GetStatistic("overview_" + statType)
}

// GetEnergyStatistic function: Retrieves energy statistics
func (r *StatsRepository) GetEnergyStatistic(address string, fromTime, toTime int64) (*models.EnergyStatisticResponse, error) {
	var stats models.EnergyStatisticResponse
	stats.Address = address
	
	// Query from database
	var daily []models.DailyStat
	r.db.Raw(`
		SELECT DATE(to_timestamp(timestamp)) as date, SUM(value) as value
		FROM statistics
		WHERE type = 'energy_usage' AND timestamp >= ? AND timestamp <= ?
		GROUP BY DATE(to_timestamp(timestamp))
		ORDER BY date
	`, fromTime, toTime).Scan(&daily)
	
	stats.Daily = daily
	return &stats, nil
}

// GetTriggerStatistic function: Retrieves trigger statistics
func (r *StatsRepository) GetTriggerStatistic(contract string, fromTime, toTime int64) (*models.TriggerStatisticResponse, error) {
	var stats models.TriggerStatisticResponse
	stats.Contract = contract
	
	var daily []models.DailyStat
	r.db.Raw(`
		SELECT DATE(to_timestamp(timestamp)) as date, COUNT(*) as value
		FROM transactions
		WHERE contract_address = ? AND block_timestamp >= ? AND block_timestamp <= ?
		GROUP BY DATE(to_timestamp(block_timestamp))
		ORDER BY date
	`, contract, fromTime, toTime).Scan(&daily)
	
	stats.Daily = daily
	stats.Count = int64(len(daily))
	
	return &stats, nil
}

// GetCallerAddressStatistic function: Retrieves caller statistics
func (r *StatsRepository) GetCallerAddressStatistic(contract string, fromTime, toTime int64) (*models.CallerAddressStatisticResponse, error) {
	var stats models.CallerAddressStatisticResponse
	stats.Contract = contract
	
	var callers []models.CallerStat
	r.db.Raw(`
		SELECT from_address as address, COUNT(*) as count
		FROM transactions
		WHERE contract_address = ? AND block_timestamp >= ? AND block_timestamp <= ?
		GROUP BY from_address
		ORDER BY count DESC
		LIMIT 100
	`, contract, fromTime, toTime).Scan(&callers)
	
	stats.Callers = callers
	return &stats, nil
}

// GetEnergyDailyStatistic function: Retrieves daily energy statistics
func (r *StatsRepository) GetEnergyDailyStatistic(fromTime, toTime int64) (*models.EnergyDailyStatisticResponse, error) {
	var stats models.EnergyDailyStatisticResponse
	
	var daily []models.EnergyDailyStat
	r.db.Raw(`
		SELECT 
			DATE(to_timestamp(timestamp)) as date,
			SUM((value->>'energy_usage')::int) as energy_usage,
			SUM((value->>'energy_fee')::int) as energy_fee
		FROM statistics
		WHERE type = 'energy_daily' AND timestamp >= ? AND timestamp <= ?
		GROUP BY DATE(to_timestamp(timestamp))
		ORDER BY date
	`, fromTime, toTime).Scan(&daily)
	
	stats.Daily = daily
	return &stats, nil
}

// GetFreezeResource function: Retrieves freeze resource information
func (r *StatsRepository) GetFreezeResource(address string, resourceType string) (*models.FreezeResourceResponse, error) {
	var response models.FreezeResourceResponse
	
	// Get frozen records
	r.db.Raw(`
		SELECT amount, expire_at, type
		FROM frozen_records
		WHERE address = ? AND (? = '' OR type = ?)
		ORDER BY expire_at
	`, address, resourceType, resourceType).Scan(&response.Frozen)
	
	// Get delegated records
	r.db.Raw(`
		SELECT from_address, to_address, amount, expire_at, type
		FROM delegated_records
		WHERE to_address = ? AND (? = '' OR type = ?)
		ORDER BY expire_at
	`, address, resourceType, resourceType).Scan(&response.Delegated)
	
	return &response, nil
}

// GetTurnover function: Retrieves turnover statistics
func (r *StatsRepository) GetTurnover(fromTime, toTime int64) (*models.TurnoverResponse, error) {
	var response models.TurnoverResponse
	
	// Get total
	r.db.Raw(`
		SELECT COALESCE(SUM(amount), 0) as total
		FROM transactions
		WHERE block_timestamp >= ? AND block_timestamp <= ?
	`, fromTime, toTime).Scan(&response.Total)
	
	// Get daily
	var daily []models.DailyTurnover
	r.db.Raw(`
		SELECT 
			DATE(to_timestamp(block_timestamp)) as date,
			SUM(amount) as turnover
		FROM transactions
		WHERE block_timestamp >= ? AND block_timestamp <= ?
		GROUP BY DATE(to_timestamp(block_timestamp))
		ORDER BY date
	`, fromTime, toTime).Scan(&daily)
	
	response.Daily = daily
	return &response, nil
}

// GetLindHolders function: Retrieves LIND holders
func (r *StatsRepository) GetLindHolders(offset, limit int) (*models.LindHolderResponse, error) {
	var response models.LindHolderResponse
	
	// Get total
	r.db.Raw("SELECT COUNT(*) FROM accounts WHERE balance > 0").Scan(&response.Total)
	
	// Get holders
	r.db.Raw(`
		SELECT address, balance, 
		ROW_NUMBER() OVER (ORDER BY balance DESC) as rank
		FROM accounts
		WHERE balance > 0
		ORDER BY balance DESC
		OFFSET ? LIMIT ?
	`, offset, limit).Scan(&response.Holders)
	
	return &response, nil
}

// GetTopAccounts function: Retrieves top accounts
func (r *StatsRepository) GetTopAccounts(limit int) ([]models.AccountResponse, error) {
    var accounts []models.AccountResponse
    err := r.db.Model(&models.Account{}).
        Order("balance DESC").
        Limit(limit).
        Find(&accounts).Error
    return accounts, err
}

// GetTopTokens function: Retrieves top tokens
func (r *StatsRepository) GetTopTokens(limit int) ([]models.LRC20TokenInfo, error) {
    var tokens []models.LRC20TokenInfo
    err := r.db.Model(&models.LRC20TokenInfo{}).
        Order("holders DESC").
        Limit(limit).
        Find(&tokens).Error
    return tokens, err
}