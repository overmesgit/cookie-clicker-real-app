package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/samber/lo"
)

type UpgradeRequest struct {
	UpgradeID string `json:"upgrade_id"`
}

type UpgradeResponse struct {
	UserUpgrade *core.Record `json:"user_upgrade"`
	Upgrade     *core.Record `json:"upgrade"`
	Price       int          `json:"price"`
}

func HandleUpgradesList(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		records, err := app.FindAllRecords(
			"user_upgrades",
			dbx.NewExp("user = {:user}", dbx.Params{"user": e.Auth.Id}),
		)
		if err != nil {
			return err
		}

		upgrades, upgErr := app.FindAllRecords(
			"upgrades",
		)
		if upgErr != nil {
			return upgErr
		}

		idToUserUpgrade := lo.KeyBy(
			records, func(item *core.Record) string {
				return item.GetString("upgrade")
			},
		)
		results := make([]UpgradeResponse, len(upgrades))
		for i, upgrade := range upgrades {
			userUpgrade, ok := idToUserUpgrade[upgrade.Id]
			count := 0
			if ok {
				count = userUpgrade.GetInt("count")
			}
			results[i] = UpgradeResponse{
				UserUpgrade: userUpgrade,
				Upgrade:     upgrade,
				Price:       GetPrice(upgrade.GetInt("rate")*10, count),
			}
		}
		return e.JSON(
			200, results,
		)
	}
}

func HandleUpgradeEvent(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var upgradeRequest UpgradeRequest
		jsonErr := e.BindBody(&upgradeRequest)
		if jsonErr != nil {
			return jsonErr
		}

		upgrade, userUpgrade, upgradeErr := getUserUpgrade(app, e.Auth.Id, upgradeRequest)
		if upgradeErr != nil {
			return fmt.Errorf("get user upgrade: %w", upgradeErr)
		}

		price := GetPrice(upgrade.GetInt("rate")*10, userUpgrade.GetInt("count"))
		scoreErr := addUserScore(app, e.Auth.Id, -price)
		if scoreErr != nil {
			return fmt.Errorf("score error: %w", scoreErr)
		}

		saveErr := saveUpgradeLevel(app, userUpgrade, e.Auth.Id, upgrade.Id)
		if saveErr != nil {
			return fmt.Errorf("save error: %w", saveErr)
		}
		newPrice := GetPrice(upgrade.GetInt("rate")*10, userUpgrade.GetInt("count"))
		return e.JSON(
			200, UpgradeResponse{
				UserUpgrade: userUpgrade,
				Upgrade:     upgrade,
				Price:       newPrice,
			},
		)
	}
}

func getUserUpgrade(
	app *pocketbase.PocketBase, userID string, upgradeRequest UpgradeRequest,
) (*core.Record, *core.Record, error) {
	collection, err := app.FindCollectionByNameOrId("user_upgrades")
	if err != nil {
		return nil, nil, err
	}

	upgrade, upgradeErr := app.FindRecordById(
		"upgrades", upgradeRequest.UpgradeID,
	)
	if upgradeErr != nil {
		return nil, nil, upgradeErr
	}

	userUpgrade, fetchErr := app.FindFirstRecordByFilter(
		collection.Name,
		"user = {:user} && upgrade = {:upgrade}",
		dbx.Params{"user": userID, "upgrade": upgradeRequest.UpgradeID},
	)
	if !errors.As(fetchErr, &sql.ErrNoRows) && fetchErr != nil {
		return nil, nil, fetchErr
	}

	if userUpgrade == nil {
		userUpgrade = core.NewRecord(collection)
	}
	return upgrade, userUpgrade, nil
}

func saveUpgradeLevel(
	app *pocketbase.PocketBase, userUpgrade *core.Record, userID string, upgradeID string,
) error {
	userUpgrade.Set("count+", 1)
	userUpgrade.Set("user", userID)
	userUpgrade.Set("upgrade", upgradeID)
	saveErr := app.Save(userUpgrade)
	if saveErr != nil {
		return saveErr
	}
	return nil
}
