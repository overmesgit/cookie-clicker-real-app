package handlers

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/samber/lo"
	"log"
)

type Event struct {
	Event string `json:"event"`
}

type Resp struct {
	Counter int `json:"counter"`
}

func HandleClickEvent(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var event Event
		jsonErr := e.BindBody(&event)
		if jsonErr != nil {
			return jsonErr
		}

		if event.Event == "click" {
			clickPower, powerErr := GetClickPower(app, e.Auth.Id)
			if powerErr != nil {
				return powerErr
			}
			err := addUserScore(app, e.Auth.Id, lo.Ternary(clickPower > 0, clickPower, 1))
			if err != nil {
				return err
			}
			if err := e.JSON(200, nil); err != nil {
				log.Printf("Write error: %v", err)
				return nil
			}

		}
		return nil
	}
}

func GetClickPower(app core.App, userID string) (int, error) {
	record := core.Record{}
	err := app.RecordQuery("user_upgrades").
		Select("SUM(count*rate) as count").
		Join("INNER JOIN", "upgrades", dbx.NewExp("user_upgrades.upgrade=upgrades.id")).
		Where(dbx.HashExp{"user": userID}).
		One(&record)

	if err != nil {
		return 0, err
	}

	return record.GetInt("count"), nil
}
