package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/samber/lo"
	"kafka-example/score"
	"math"
)

type LevelRequest struct {
	LevelID string `json:"level_id"`
}

type LevelResponse struct {
	UserLevel *core.Record `json:"user_level"`
	Level     *core.Record `json:"level"`
	Price     int          `json:"price"`
}

func GetPrice(basePrice int, count int) int {
	inc := float64(basePrice) * math.Pow(1.3, float64(count))
	return int(math.Round(inc))
}

func HandleLevelsList(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		records, err := app.FindAllRecords(
			"user_levels",
			dbx.NewExp("user = {:user}", dbx.Params{"user": e.Auth.Id}),
		)
		if err != nil {
			return fmt.Errorf("user_levels error: %w", err)
		}
		levels, err := app.FindAllRecords(
			"levels",
		)
		if err != nil {
			return fmt.Errorf("levels error: %w", err)
		}

		idToUserLevel := lo.KeyBy(
			records, func(item *core.Record) string {
				return item.GetString("level")
			},
		)
		results := make([]LevelResponse, len(levels))
		for i, lvl := range levels {
			userLevel, ok := idToUserLevel[lvl.Id]
			count := 0
			if ok {
				count = userLevel.GetInt("count")
			}
			results[i] = LevelResponse{
				UserLevel: userLevel,
				Level:     lvl,
				Price:     GetPrice(lvl.GetInt("cost"), count),
			}
		}
		return e.JSON(
			200, results,
		)
	}
}

func HandleLevelEvent(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var levelRequest LevelRequest
		jsonErr := e.BindBody(&levelRequest)
		if jsonErr != nil {
			return jsonErr
		}

		level, userLevel, levelErr := getUserLevel(app, e.Auth.Id, levelRequest)
		if levelErr != nil {
			return fmt.Errorf("user_levels error: %w", levelErr)
		}

		price := GetPrice(level.GetInt("cost"), userLevel.GetInt("count"))
		scoreErr := addUserScore(app, e.Auth.Id, -price)
		if scoreErr != nil {
			return fmt.Errorf("update score error: %w", scoreErr)
		}

		saveErr := saveUserLevel(e.Auth.Id, userLevel, levelRequest, app)
		if saveErr != nil {
			return fmt.Errorf("save user level error: %w", saveErr)
		}
		score.ScoreUpdater(app, e.Auth.Id, level.GetInt("rate"), userLevel.Id)
		newPrice := GetPrice(level.GetInt("cost"), userLevel.GetInt("count"))
		return e.JSON(
			200, LevelResponse{
				UserLevel: userLevel,
				Level:     level,
				Price:     newPrice,
			},
		)
	}
}

func saveUserLevel(
	userID string, userLevel *core.Record, levelRequest LevelRequest, app *pocketbase.PocketBase,
) error {
	userLevel.Set("count+", 1)
	userLevel.Set("user", userID)
	userLevel.Set("level", levelRequest.LevelID)
	saveErr := app.Save(userLevel)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

func InitUserLevels(app *pocketbase.PocketBase) error {
	userLevels, userErr := app.FindAllRecords(
		"user_levels",
	)
	if userErr != nil {
		return userErr
	}

	levels, levelErr := app.FindAllRecords(
		"levels",
	)

	if levelErr != nil {
		return levelErr
	}
	levelIDtoLevel := lo.KeyBy(
		levels, func(r *core.Record) string {
			return r.GetString("id")
		},
	)
	for _, userLevel := range userLevels {
		level := levelIDtoLevel[userLevel.GetString("level")]
		for range userLevel.GetInt("count") {
			score.ScoreUpdater(app, userLevel.GetString("user"), level.GetInt("rate"), userLevel.Id)
		}
	}

	return nil
}

func getUserLevel(
	app *pocketbase.PocketBase, userID string, levelRequest LevelRequest,
) (*core.Record, *core.Record, error) {
	collection, err := app.FindCollectionByNameOrId("user_levels")
	if err != nil {
		return nil, nil, err
	}

	level, levelErr := app.FindRecordById(
		"levels", levelRequest.LevelID,
	)
	if levelErr != nil {
		return nil, nil, levelErr
	}

	userLevel, fetchErr := app.FindFirstRecordByFilter(
		collection.Name,
		"user = {:user} && level = {:level}",
		dbx.Params{"user": userID, "level": levelRequest.LevelID},
	)
	if !errors.As(fetchErr, &sql.ErrNoRows) && fetchErr != nil {
		return nil, nil, fetchErr
	}

	if userLevel == nil {
		userLevel = core.NewRecord(collection)
	}
	return level, userLevel, nil
}

func addUserScore(
	app *pocketbase.PocketBase, userID string, price int,
) error {
	userScore, scoreErr := app.FindFirstRecordByData("counter", "user", userID)
	if errors.As(scoreErr, &sql.ErrNoRows) {
		collection, err := app.FindCollectionByNameOrId("counter")
		if err != nil {
			return err
		}
		userScore = core.NewRecord(collection)
		userScore.Set("user", userID)
	} else if scoreErr != nil {
		return scoreErr
	}

	userScore.Set("count+", price)
	userSaveErr := app.Save(userScore)
	if userSaveErr != nil {
		return userSaveErr
	}
	return nil
}
