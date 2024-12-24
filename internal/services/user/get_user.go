package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wavinamayola/user-management/internal/models"
	"github.com/wavinamayola/user-management/internal/storage"
	"github.com/wavinamayola/user-management/internal/utils"
)

func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	log.Print("received get user request")
	id := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("invalid id format: %+v", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid id format, should be integer"))
		return
	}

	user, err := u.store.GetUser(userID)
	if err != nil {
		if err == storage.ErrNotFound {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err)
		}
		return
	}

	respBody := models.Response{
		Data:    user,
		Message: "user details",
	}
	utils.RespondWithJSON(w, http.StatusOK, respBody)
	log.Print("user retrieved successfully")
}
