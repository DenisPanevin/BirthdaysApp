package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
	"time"
	"unicode/utf8"
)

type User struct {
	Id           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Date         time.Time `json:"date"`
}

func NoCyrillic() func(string) error {
	return func(value string) error {
		for _, r := range value {
			if utf8.ValidRune(r) && (r >= '\u0400' && r <= '\u04FF' || r >= '\u0500' && r <= '\u052F' || r >= '\u2DE0' && r <= '\u2DFF' || r >= '\uA640' && r <= '\uA69F') {
				return fmt.Errorf("string contains Cyrillic characters")
			}
		}
		return nil
	}
}

func (u *User) ValidateUser() error {
	if err := NoCyrillic()(u.Email); err != nil {
		fmt.Println("Validation failed for string without Cyrillic characters:", err)
	} else {
		fmt.Println("Validation succeeded for string without Cyrillic characters")
	}
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(4, 100)),
		validation.Field(&u.Name, validation.Required),
	)

}

func (u *User) ClearFields() {
	u.Password = ""
	//u.Date=
}

func (u *User) ComparePassword(password string) bool {

	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
