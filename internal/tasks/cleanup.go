package tasks

import (
	"context"
	"time"

	"shortener/internal/logger"
	"shortener/internal/service"
)

// Run runs a background task to clean up the storage periodically.
//
// It uses a ticker to schedule the cleanup at the specified interval.
func Run(ctx context.Context, store service.URLStorage, log *logger.Log, interval time.Duration) {
	log.Debug("starting storage cleanup task", "period", interval)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		select {
		case <-ctx.Done():
			log.Debug("Stopping ticker..")
			ticker.Stop()
			return
		default:
			urls, err := store.Cleanup(ctx)
			if err != nil {
				log.Err("failed cleanup storage", err)
			}
			if len(urls) > 0 {
				log.Info("The following url IDs has been deleted from the storage", "URLs", urls)
			} else {
				log.Info("Nothing to delete. Going to sleep", "time", interval)
			}
		}
	}
}
