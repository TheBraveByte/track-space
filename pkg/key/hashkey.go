package key

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

//HashPassword : this is to convert the user password using the hash algorithm
// to generate string of bytes of character
func HashPassword(inputPassword string) string {
	keyByte, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 12)

	if err != nil {
		log.Println(err)
		panic("cannot generate hashed Password")

	}
	return string(keyByte)
}

//VerifyPassword : for confirmation to check if the hashed password and the user
//password  matches with each other when hashed
func VerifyPassword(inputPassword, hashedPassword string) (bool, string) {
	var hashMsg string
	var validHash bool
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	hashMsg = "password hashed successfully"
	validHash = true
	if err == bcrypt.ErrMismatchedHashAndPassword {
		hashMsg = "input password and hashed password dont matches"
		validHash = false
		log.Println(err)
		return validHash, hashMsg
	}

	return validHash, hashMsg
}
