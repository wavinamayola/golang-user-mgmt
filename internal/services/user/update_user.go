package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wavinamayola/user-management/internal/models"
	"github.com/wavinamayola/user-management/internal/storage"
	"github.com/wavinamayola/user-management/internal/utils"
)

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	log.Print("received update user request")
	id := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("invalid id format: %+v", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid id format, should be integer"))
		return
	}

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

	if err := u.store.UpdateUser(userID, user); err != nil {
		if err == storage.ErrNoRowsAffected {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("user doesn't exist"))
			return
		}
		log.Printf("failed to update user: %+v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	user.ID = userID
	respBody := models.Response{
		Data:    user,
		Message: "user updated successfully",
	}
	utils.RespondWithJSON(w, http.StatusOK, respBody)
	log.Print("user updated successfully")
}
