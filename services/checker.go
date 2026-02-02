package services

import (
	"log"
	"net/http"
	"sync"
	"time"

	"omega-home/models"
)

var statusMap sync.Map

type ServiceStatus struct {
	Online    bool   `json:"online"`
	Latency   int64  `json:"latency"` // ms
	CheckedAt string `json:"checked_at"`
}

func GetAllStatus() map[uint]ServiceStatus {
	result := make(map[uint]ServiceStatus)
	statusMap.Range(func(key, value interface{}) bool {
		result[key.(uint)] = value.(ServiceStatus)
		return true
	})
	return result
}

func StartChecker(interval int) {
	go func() {
		for {
			checkAll()
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
}

func checkAll() {
	var svcs []models.Service
	models.DB.Where("status_check = ?", true).Find(&svcs)

	var wg sync.WaitGroup
	client := &http.Client{Timeout: 5 * time.Second}

	for _, svc := range svcs {
		wg.Add(1)
		go func(s models.Service) {
			defer wg.Done()
			start := time.Now()
			resp, err := client.Get(s.URL)
			latency := time.Since(start).Milliseconds()
			online := err == nil && resp != nil && resp.StatusCode < 500
			if resp != nil {
				resp.Body.Close()
			}
			statusMap.Store(s.ID, ServiceStatus{
				Online:    online,
				Latency:   latency,
				CheckedAt: time.Now().Format(time.RFC3339),
			})
		}(svc)
	}
	wg.Wait()
	log.Printf("Status check complete: %d services checked", len(svcs))
}
