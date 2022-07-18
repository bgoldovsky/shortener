package auth

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mockUrls "github.com/bgoldovsky/shortener/internal/app/services/auth/mocks"
)

func TestService_SignUp(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		token  string
		err    error
	}{
		{
			name:   "success",
			userID: "qwerty",
			token:  "1234567890",
		},
		{
			name:   "success",
			userID: "qwerty",
			token:  "",
			err:    errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		genMock := mockUrls.NewMockgenerator(ctrl)
		genMock.EXPECT().RandomString(int64(idLength)).Return(tt.userID, nil)

		hasherMock := mockUrls.NewMockhasher(ctrl)
		hasherMock.EXPECT().Sign(tt.userID).Return(tt.token, tt.err)

		s := NewService(genMock, hasherMock)
		actUserID, actToken, err := s.SignUp()

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.userID, actUserID)
		assert.Equal(t, tt.token, actToken)
	}
}

func TestService_SignIn(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		userID string
		err    error
	}{
		{
			name:   "success",
			token:  "1234567890",
			userID: "qwerty",
		},
		{
			name:   "success",
			token:  "1234567890",
			userID: "",
			err:    errors.New("test err"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		genMock := mockUrls.NewMockgenerator(ctrl)

		hasherMock := mockUrls.NewMockhasher(ctrl)
		hasherMock.EXPECT().Validate(tt.token, int64(idLength)).Return(tt.userID, tt.err)

		s := NewService(genMock, hasherMock)
		act, err := s.SignIn(tt.token)

		assert.Equal(t, tt.err, err)
		assert.Equal(t, tt.userID, act)
	}
}
