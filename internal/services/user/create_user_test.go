package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_store "github.com/wavinamayola/user-management/internal/mock/mock_store"
	"github.com/wavinamayola/user-management/internal/models"
)

const (
	testUsername  = "testusername"
	testFirstName = "testfirstname"
	testLastName  = "testlastname"
	testEmail     = "testemail@gmail.com"
	testAge       = 25
)

func Test_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	u := &User{store: mockStore}

	t.Run("successful user creation", func(t *testing.T) {
		input := models.UserRequest{
			Username:  testUsername,
			FirstName: testFirstName,
			LastName:  testLastName,
			Email:     testEmail,
			Age:       testAge,
		}

		mockStore.EXPECT().CreateUser(input).Return(1, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		w := httptest.NewRecorder()

		u.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respBody models.Response
		_ = json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "user created successfully", respBody.Message)

		var user models.User
		dataBytes, err := json.Marshal(respBody.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &user)
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
	})
}

func Test_Create_WithErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockStore(ctrl)
	u := &User{store: mockStore}

	tests := []struct {
		name           string
		input          models.UserRequest
		mockReturnID   int
		mockReturnErr  error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "validation failure",
			input: models.UserRequest{
				Email: "invalid-email",
			},
			mockReturnID:   0,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "validation errors",
		},
		{
			name: "storage error during user creation",
			input: models.UserRequest{
				Username:  testUsername,
				FirstName: testFirstName,
				LastName:  testLastName,
				Email:     testEmail,
				Age:       testAge,
			},
			mockReturnID:   0,
			mockReturnErr:  errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "database error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.mockReturnErr != nil {
				mockStore.EXPECT().CreateUser(test.input).Return(test.mockReturnID, test.mockReturnErr).AnyTimes()
			}

			body, _ := json.Marshal(test.input)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			u.Create(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, test.expectedStatus, resp.StatusCode)

			var respBody models.Response
			_ = json.NewDecoder(resp.Body).Decode(&respBody)
			assert.Contains(t, respBody.Message, test.expectedMsg)
		})
	}
}
