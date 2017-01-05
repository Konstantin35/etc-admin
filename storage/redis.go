package storage

import (
	"fmt"
	"github.com/cihub/seelog"
	"gopkg.in/redis.v3"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type RedisConfig struct {
	Endpoint string `json:"endpoint"`
	Password string `json:"password"`
	Database int64  `json:"database"`
	PoolSize int    `json:"poolSize"`
}

type RedisClient struct {
	client *redis.Client
	prefix string
}

type Miner struct {
	LastBeat  int64 `json:"lastBeat"`
	HR        int64 `json:"hr"`
	Offline   bool  `json:"offline"`
	startedAt int64
}

type Worker struct {
	Miner
	TotalHR int64 `json:"hr2"`
}

const (
	smallWindow string = "30m"
	largeWindow string = "3h"
)

func (r *RedisClient) formatKey(args ...interface{}) string {
	return join(r.prefix, join(args...))
}

func join(args ...interface{}) string {
	s := make([]string, len(args))
	for i, v := range args {
		switch v.(type) {
		case string:
			s[i] = v.(string)
		case int64:
			s[i] = strconv.FormatInt(v.(int64), 10)
		case uint64:
			s[i] = strconv.FormatUint(v.(uint64), 10)
		case float64:
			s[i] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		case bool:
			if v.(bool) {
				s[i] = "1"
			} else {
				s[i] = "0"
			}
		case *big.Int:
			n := v.(*big.Int)
			if n != nil {
				s[i] = n.String()
			} else {
				s[i] = "0"
			}
		default:
			panic("Invalid type specified for conversion")
		}
	}
	return strings.Join(s, ":")
}

//NewRedisClient init a client for reuse
func NewRedisClient(cfg *RedisConfig, prefix string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.Database,
		PoolSize: cfg.PoolSize,
	})
	return &RedisClient{client: client, prefix: prefix}
}

func (r *RedisClient) GetPoolChartData() ([]map[string]interface{}, error) {
	timenow := time.Now()

	var chartdata []map[string]interface{}
	cmd := r.client.LRange(r.formatKey("chartdata", "pool", "totalhash"), 0, -1)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return chartdata, cmd.Err()
	}
	idx := 0
	stringArray, _ := cmd.Result()
	for index, subvalue := range stringArray {
		substr := strings.Split(subvalue, "=")
		t, _ := time.Parse("2006-01-02 15:04:05", substr[0])
		duration := timenow.Unix() - t.Unix()
		if duration > 604800 { //大于7天
			idx = index
			break
		}
		if index == 167 {
			idx = index + 1
		}
	}

	//get now interfer time
	origin := time.Now().Local().Unix()
	j := 0
	for i := 0; i < 168; i++ {
		temp := make(map[string]interface{})
		timestamp := time.Unix(origin, 0).Format("2006-01-02 15:00:00")
		temp["tempstamp"] = timestamp
		if j < idx {
			substr := strings.Split(stringArray[j], "=")
			hashrate, _ := strconv.ParseInt(substr[1], 10, 64)
			temp["value"] = hashrate
		} else {
			temp["value"] = int64(0)
		}
		origin = origin - 3600
		j++
		chartdata = append(chartdata, temp)
	}

	return chartdata, nil
}

func (r *RedisClient) GetMainStatisticData() map[string]interface{} {
	value := make(map[string]interface{})
	stats, err := r.client.HGetAllMap(r.formatKey("pooldata", "statistic")).Result()
	if err != nil && err != redis.Nil {
		return nil
	}
	num, _ := strconv.ParseInt(stats["minersTotal"], 10, 32)
	hash, _ := strconv.ParseInt(stats["hashrate"], 10, 64)
	value["minersTotal"] = num
	value["hashrate"] = hash
	//get pool account address rest etc/eth coin value

	return value
}

//GetRevenue get given user last payment and total payment amount
func (r *RedisClient) GetRevenue(wallet string) (last int64, total int64) {
	tx := r.client.Multi()
	defer tx.Close()

	cmds, err := tx.Exec(func() error {
		tx.HGetAllMap(r.formatKey("miners", wallet))
		tx.ZRevRangeWithScores(r.formatKey("payments", wallet), 0, 1)
		return nil
	})

	if err != nil && err != redis.Nil {
		seelog.Info("get revenue error:", err)
	} else {
		result, _ := cmds[0].(*redis.StringStringMapCmd).Result()
		total, _ = strconv.ParseInt(result["paid"], 10, 64)
		payments := convertPaymentsResults(cmds[1].(*redis.ZSliceCmd))
		last = payments[0]["amount"].(int64)
	}

	return
}

func (r *RedisClient) GetWalletRevenue(wallet string) map[string]interface{} {
	tx := r.client.Multi()
	defer tx.Close()

	var retvalue = make(map[string]interface{})
	cmds, err := tx.Exec(func() error {
		tx.HGetAllMap(r.formatKey("miners", wallet))
		tx.ZRevRangeWithScores(r.formatKey("payments", wallet), 0, 1)
		return nil
	})

	if err != nil && err != redis.Nil {
		seelog.Info("get revenue error:", err)
		return nil
	}
	result, _ := cmds[0].(*redis.StringStringMapCmd).Result()
	retvalue = convertStringMap(result)
	payments := convertPaymentsResults(cmds[1].(*redis.ZSliceCmd))
	retvalue["lastPaid"] = payments[0]["amount"].(int64)

	return retvalue
}

//GetAccountChartValues get account 24h hashrate for chart, timestamp from recent integer time clock
func (r *RedisClient) GetAccountChartValues(login string) ([]map[string]interface{}, error) {
	timenow := time.Now()

	var chartdata []map[string]interface{}
	cmd := r.client.LRange(r.formatKey("chart", "hashrate", login), 0, -1)
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		return chartdata, cmd.Err()
	}
	idx := 0
	stringArray, _ := cmd.Result()
	for index, subvalue := range stringArray {
		substr := strings.Split(subvalue, "=")
		t, _ := time.Parse("2006-01-02 15:04:05", substr[0])
		duration := timenow.Unix() - t.Unix()
		if duration > 86400 {
			idx = index
			break
		}
		if index == 23 {
			idx = index + 1
		}
	}

	//get now interfer time
	origin := time.Now().Local().Unix()
	j := 0
	for i := 0; i < 24; i++ {
		temp := make(map[string]interface{})
		timestamp := time.Unix(origin, 0).Format("2006-01-02 15:00:00")
		temp["tempstamp"] = timestamp
		if j < idx {
			substr := strings.Split(stringArray[j], "=")
			hashrate, _ := strconv.ParseInt(substr[1], 10, 64)
			temp["value"] = hashrate
		} else {
			temp["value"] = int64(0)
		}
		origin = origin - 3600
		j++
		chartdata = append(chartdata, temp)
	}

	return chartdata, nil
}

func (r *RedisClient) GetWorkersStats(address string) (map[string]interface{}, error) {
	small, _ := time.ParseDuration(smallWindow)
	large, _ := time.ParseDuration(largeWindow)
	stats, err := r.collectWorkerStats(small, large, address)
	if err != nil {
		seelog.Error("get worker stats error:", err)
		return nil, err
	}
	return stats, nil
}

//GetPaymentHistory get a time duration payment records, if you dont give any date, it will return all payments records
func (r *RedisClient) GetPaymentHistory(address string, begin, end int64) []map[string]interface{} {
	// payments := r.getPayments(address)
	tx := r.client.Multi()
	defer tx.Close()

	cmds, err := tx.Exec(func() error {
		tx.ZRevRangeWithScores(r.formatKey("payments", address), 0, -1)
		return nil
	})

	if err != nil && err != redis.Nil {
		seelog.Info("get paymens error:", err)
		return nil
	}
	payments := convertPaymentsResults(cmds[0].(*redis.ZSliceCmd))
	seelog.Info("raw payments data:", payments)

	if begin == -1 && end == -1 { //get all records
		return payments
	}
	length := len(payments)
	if payments[0]["timestamp"].(int64) > end || payments[length-1]["timestamp"].(int64) < begin {
		return nil
	}

	index := 0
	for idx, val := range payments {
		stamp := val["timestamp"].(int64)
		if stamp >= begin && stamp <= end {
			continue
		}
		index = idx
		break
	}
	return payments[0:index]
}

// Try to convert all numeric strings to int64
func convertStringMap(m map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	var err error
	for k, v := range m {
		result[k], err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			result[k] = v
		}
	}
	return result
}

func convertPaymentsResults(raw *redis.ZSliceCmd) []map[string]interface{} {
	var result []map[string]interface{}
	for _, v := range raw.Val() {
		seelog.Info("every data:", v.Score)
		tx := make(map[string]interface{})
		tx["timestamp"] = int64(v.Score)
		fields := strings.Split(v.Member.(string), ":")
		tx["tx"] = fields[0]
		// Individual or whole payments row
		if len(fields) < 3 {
			tx["amount"], _ = strconv.ParseInt(fields[1], 10, 64)
		} else {
			tx["address"] = fields[1]
			tx["amount"], _ = strconv.ParseInt(fields[2], 10, 64)
		}
		result = append(result, tx)
	}
	return result
}

func (r *RedisClient) collectWorkerStats(sWindow, lWindow time.Duration, login string) (map[string]interface{}, error) {
	smallWindow := int64(sWindow / time.Second)
	largeWindow := int64(lWindow / time.Second)
	stats := make(map[string]interface{})

	tx := r.client.Multi()
	defer tx.Close()

	now := MakeTimestamp() / 1000

	cmds, err := tx.Exec(func() error {
		tx.ZRemRangeByScore(r.formatKey("hashrate", login), "-inf", fmt.Sprint("(", now-largeWindow))
		tx.ZRangeWithScores(r.formatKey("hashrate", login), 0, -1)
		return nil
	})

	if err != nil {
		return nil, err
	}

	totalHashrate := int64(0)
	currentHashrate := int64(0)
	online := int64(0)
	offline := int64(0)
	workers := convertWorkersStats(smallWindow, cmds[1].(*redis.ZSliceCmd))

	for id, worker := range workers {
		timeOnline := now - worker.startedAt
		if timeOnline < 600 {
			timeOnline = 600
		}

		boundary := timeOnline
		if timeOnline >= smallWindow {
			boundary = smallWindow
		}
		worker.HR = worker.HR / boundary

		boundary = timeOnline
		if timeOnline >= largeWindow {
			boundary = largeWindow
		}
		worker.TotalHR = worker.TotalHR / boundary

		if worker.LastBeat < (now - smallWindow/2) {
			worker.Offline = true
			offline++
		} else {
			online++
		}

		currentHashrate += worker.HR
		totalHashrate += worker.TotalHR
		workers[id] = worker
	}
	stats["workers"] = workers
	stats["workersTotal"] = len(workers)
	stats["workersOnline"] = online
	stats["workersOffline"] = offline
	stats["hashrate"] = totalHashrate
	stats["currentHashrate"] = currentHashrate
	return stats, nil
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func convertWorkersStats(window int64, raw *redis.ZSliceCmd) map[string]Worker {
	now := MakeTimestamp() / 1000
	workers := make(map[string]Worker)

	for _, v := range raw.Val() {
		parts := strings.Split(v.Member.(string), ":")
		share, _ := strconv.ParseInt(parts[0], 10, 64)
		id := parts[1]
		score := int64(v.Score)
		worker := workers[id]

		// Add for large window
		worker.TotalHR += share

		// Add for small window if matches
		if score >= now-window {
			worker.HR += share
		}

		if worker.LastBeat < score {
			worker.LastBeat = score
		}
		if worker.startedAt > score || worker.startedAt == 0 {
			worker.startedAt = score
		}
		workers[id] = worker
	}
	return workers
}
