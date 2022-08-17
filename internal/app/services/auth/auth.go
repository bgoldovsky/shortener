//go:generate mockgen -source=auth.go -destination=mocks/mocks.go
package auth

import "github.com/sirupsen/logrus"

const idLength = 8

type generator interface {
	RandomString(n int64) (string, error)
}

type hasher interface {
	Sign(value string) (string, error)
	Validate(value string, dataLength int64) (string, error)
}

type service struct {
	hasher    hasher
	generator generator
}

func NewService(generator generator, hasher hasher) *service {
	return &service{
		generator: generator,
		hasher:    hasher,
	}
}

// SignUp Регистрирует пользователя в системе, возвращая userID и подписанный токен
func (s *service) SignUp() (string, string, error) {
	userID, err := s.generator.RandomString(idLength)
	if err != nil {
		logrus.WithError(err).WithField("userID", userID).Error("generate userID error")
		return userID, "", err
	}

	signedUserID, err := s.hasher.Sign(userID)
	if err != nil {
		logrus.WithError(err).WithField("userID", userID).Error("sign userID error")
		return userID, "", err
	}

	return userID, signedUserID, nil
}

// SignIn Аутентифицирует пользователя по токену и возвращает его userID
func (s *service) SignIn(token string) (string, error) {
	userID, err := s.hasher.Validate(token, idLength)
	if err != nil {
		logrus.WithError(err).WithField("token", token).Error("validate userID sign error")
		return "", err
	}

	return userID, nil
}
