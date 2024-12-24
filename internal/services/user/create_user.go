package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/wavinamayola/user-management/internal/models"
	"github.com/wavinamayola/user-management/internal/utils"
)

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	log.Print("received create user request")

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("failed to decode user request: %+v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := validate.Struct(user); err != nil {
		log.Printf("failed to validate user request: %+v", err)
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	id, err := u.store.CreateUser(user)
	if err != nil {
		log.Printf("failed to create user: %+v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	user.ID = int(id)
	respBody := models.Response{
		Data:    user,
		Message: "user created successfully",
	}
	utils.RespondWithJSON(w, http.StatusCreated, respBody)
	log.Print("user created successfully")
}
