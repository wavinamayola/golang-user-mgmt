package user

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock_store "github.com/wavinamayola/user-management/internal/mock/mock_store"
	"github.com/wavinamayola/user-management/internal/models"
	"github.com/wavinamayola/user-management/internal/storage"
)

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	u := &User{store: mockStore}

	tests := []struct {
		name          string
		userID        string
		input         models.UserRequest
		mockReturnErr error
		expectedCode  int
		expectedBody  string
	}{
		{
			name:   "Successful update",
			userID: "1",
			input: models.UserRequest{
				Username:  testUsername,
				FirstName: testFirstName,
				LastName:  testLastName,
				Email:     testEmail,
				Age:       testAge,
			},
			mockReturnErr: nil,
			expectedCode:  http.StatusOK,
			expectedBody:  `"user updated successfully"`,
		},
		{
			name:   "User not found",
			userID: "2",
			input: models.UserRequest{
				Username:  testUsername,
				FirstName: testFirstName,
				LastName:  testLastName,
				Email:     testEmail,
				Age:       testAge,
			},
			mockReturnErr: storage.ErrNoRowsAffected,
			expectedCode:  http.StatusNotFound,
			expectedBody:  `"user doesn't exist"`,
		},
		{
			name:   "Validation error",
			userID: "1",
			input: models.UserRequest{
				Email: "invalid-email",
			},
			mockReturnErr: nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  `"validation errors"`,
		},
		{
			name:          "Invalid ID format",
			userID:        "abc",
			input:         models.UserRequest{},
			mockReturnErr: nil,
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "invalid id format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.input)

			req := httptest.NewRequest(http.MethodPut, "/users/"+test.userID, bytes.NewReader(body))
			w := httptest.NewRecorder()
			vars := map[string]string{"id": test.userID}
			req = mux.SetURLVars(req, vars)

			userID, err := strconv.Atoi(test.userID)
			if err == nil {
				mockStore.EXPECT().UpdateUser(userID, test.input).Return(test.mockReturnErr).AnyTimes()
			}

			u.Update(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expectedCode, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(respBody), test.expectedBody)
		})
	}
}
