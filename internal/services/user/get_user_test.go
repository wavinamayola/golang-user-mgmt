package user

import (
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

func Test_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	u := &User{store: mockStore}

	tests := []struct {
		name           string
		userID         string
		mockReturnUser models.User
		mockReturnErr  error
		expectedCode   int
		expectedBody   string
	}{
		{
			name:   "Successful user retrieval",
			userID: "1",
			mockReturnUser: models.User{
				ID:        1,
				Username:  testUsername,
				FirstName: testFirstName,
				LastName:  testLastName,
				Email:     testEmail,
				Age:       testAge,
			},
			mockReturnErr: nil,
			expectedCode:  http.StatusOK,
			expectedBody:  `"user details"`,
		},
		{
			name:           "User not found",
			userID:         "2",
			mockReturnUser: models.User{},
			mockReturnErr:  storage.ErrNotFound,
			expectedCode:   http.StatusNotFound,
			expectedBody:   `"user not found"`,
		},
		{
			name:           "Invalid ID format",
			userID:         "abc",
			mockReturnUser: models.User{},
			mockReturnErr:  nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody:   "invalid id format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/"+test.userID, nil)
			w := httptest.NewRecorder()
			vars := map[string]string{"id": test.userID}
			req = mux.SetURLVars(req, vars)

			userID, err := strconv.Atoi(test.userID)
			if err == nil {
				mockStore.EXPECT().GetUser(userID).Return(test.mockReturnUser, test.mockReturnErr).AnyTimes()
			}

			u.Get(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expectedCode, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(body), test.expectedBody)
		})
	}
}
