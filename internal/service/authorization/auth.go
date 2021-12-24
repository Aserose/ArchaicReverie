package authorization

import (
	"crypto/sha1"
	"fmt"
	"github.com/Aserose/ArchaicReverie/internal/repository/model"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

const (
	empty                  = ""
	AuthServicePackageName = "AuthService"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId    int             `json:"id"`
	Character model.Character `json:"character"`
	//TODO
}

func (s serviceAuthorization) SignUp(username, password string) string {
	id, status := s.db.Postgres.UserData.Create(username, s.createPasswordHash(password))
	if status != empty {
		return status
	}
	return s.createToken(id, model.Character{})
}

func (s serviceAuthorization) SignIn(username, password string) string {
	id, status := s.db.Postgres.UserData.Check(username, s.createPasswordHash(password))
	if status != empty {
		return status
	}
	return s.createToken(id, model.Character{})
}

func (s serviceAuthorization) UpdateToken(userId int, character model.Character) string {
	log.Print("update token ", character)
	return s.createToken(userId, character)
}

func (s serviceAuthorization) createToken(userId int, character model.Character) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix()},
		userId,
		character,
	})

	tokenString, err := token.SignedString([]byte(s.cfgServices.HMACSecret))
	if err != nil {
		s.log.Errorf(s.logMsg.FormatErr, AuthServicePackageName, s.logMsg.CreateToken, err.Error())
		return empty
	}

	return tokenString
}

func (s serviceAuthorization) Verification(tokenString string) (int, model.Character, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfgServices.HMACSecret), nil
	})
	if err != nil {
		s.log.Errorf(s.logMsg.FormatErr, AuthServicePackageName, s.logMsg.ReadToken, err.Error())
		return 0, model.Character{}, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.UserId, claims.Character, nil
	} else {
		s.log.Errorf(s.logMsg.FormatErr, AuthServicePackageName, s.logMsg.ReadToken, err.Error())
		return 0, model.Character{}, err
	}
}

func (s serviceAuthorization) createPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.cfgServices.PasswordSalt)))
}
