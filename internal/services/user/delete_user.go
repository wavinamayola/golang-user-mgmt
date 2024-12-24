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

func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	log.Print("received delete user request")
	id := mux.Vars(r)["id"]

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("invalid id format: %+v", err)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid id format, should be integer"))
		return
	}

	if err := u.store.DeleteUser(userID); err != nil {
		if err == storage.ErrNoRowsAffected {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("user doesn't exist"))
			return
		}
		log.Printf("failed to delete user: %+v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respBody := models.Response{
		Data:    userID,
		Message: "user deleted successfully",
	}
	utils.RespondWithJSON(w, http.StatusOK, respBody)
	log.Print("user deleted successfully")
}
