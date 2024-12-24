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
	"github.com/wavinamayola/user-management/internal/storage"
)

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	u := &User{store: mockStore}

	tests := []struct {
		name         string
		userID       string
		mockError    error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful deletion",
			userID:       "1",
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedBody: `"user deleted successfully"`,
		},
		{
			name:         "User not found",
			userID:       "2",
			mockError:    storage.ErrNoRowsAffected,
			expectedCode: http.StatusNotFound,
			expectedBody: `"user doesn't exist"`,
		},
		{
			name:         "Invalid ID format",
			userID:       "abc",
			mockError:    nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid id format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/users/"+test.userID, nil)
			w := httptest.NewRecorder()
			vars := map[string]string{"id": test.userID}
			req = mux.SetURLVars(req, vars)

			userID, err := strconv.Atoi(test.userID)
			if err == nil {
				mockStore.EXPECT().DeleteUser(userID).Return(test.mockError).AnyTimes()
			}

			u.Delete(w, req)
			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.expectedCode, res.StatusCode)
			assert.Contains(t, string(body), test.expectedBody)
		})
	}
}
