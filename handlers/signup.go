package handlers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"net/http"
)

type SignUpRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type SignUpResponse struct {
	Success bool         `json:"success"`
	User    *core.Record `json:"user,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func HandleSignUp(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var req SignUpRequest
		if err := e.BindBody(&req); err != nil {
			return e.JSON(
				http.StatusBadRequest, SignUpResponse{
					Success: false,
					Error:   err.Error(),
				},
			)
		}

		if req.Password != req.PasswordConfirm {
			return e.JSON(
				http.StatusBadRequest, SignUpResponse{
					Success: false,
					Error:   "Passwords do not match",
				},
			)
		}

		collection, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		newUser := core.NewRecord(collection)
		newUser.SetEmail(req.Email)
		newUser.SetPassword(req.Password)

		saveErr := app.Save(newUser)
		if saveErr != nil {
			return saveErr
		}

		counter, err := app.FindCollectionByNameOrId("counter")
		if err != nil {
			return err
		}

		newCounter := core.NewRecord(counter)
		newCounter.Set("user", newUser.Id)

		counterErr := app.Save(newCounter)
		if counterErr != nil {
			return counterErr
		}

		response := SignUpResponse{
			Success: true,
			User:    newUser,
		}
		return e.JSON(http.StatusOK, response)
	}
}
