package score

import (
	"github.com/pocketbase/pocketbase"
	"log"
	"math/rand"
	"time"
)

func ScoreUpdater(app *pocketbase.PocketBase, userID string, rate int, userLevelID string) {
	go func() {
		for range 30 {
			record, err := app.FindFirstRecordByData("counter", "user", userID)
			if err != nil {
				log.Printf(err.Error())
			}
			record.Set("count+", rate)
			err = app.Save(record)
			if err != nil {
				log.Printf(err.Error())
			}
			randDuration := (rand.Float64() + 0.5) * float64(time.Second)
			time.Sleep(time.Duration(randDuration))
		}
		userLevel, findErr := app.FindRecordById("user_levels", userLevelID)
		if findErr != nil {
			log.Printf(findErr.Error())
		}
		userLevel.Set("count-", 1)
		userErr := app.Save(userLevel)
		if userErr != nil {
			log.Printf(userErr.Error())
		}
	}()
}
