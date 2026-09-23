package cron

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func ExpiryDateCronjob(ctx context.Context, wg *sync.WaitGroup) {
	for {
		log.Println("Running expiry date cronjob...")
		wg.Add(1)
		err := service.CheckExpiredItems(ctx)
		if err != nil {
			log.Printf("Expiry date cronjob error - %v", err.Error())
		}
		log.Println("Expiry date cronjob done!")
		wg.Done()
		next := nextTime()
		select {
		case <-time.After(time.Until(next)):

		case <-ctx.Done():
			return
		}
	}
}

func nextTime() time.Time {
	t := time.Now()
	nextTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	return nextTime
}
