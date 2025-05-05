package util

import "golang.org/x/crypto/bcrypt"

func HashedPassword(password string) (string, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	return string(hasedPassword) , err
}

func CheckPassword(hashedPassword, password string) (error){
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
}