package utils

import (
	"github.com/rohitashvadangi/identity-server/internal/proto"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CreateUsers() map[string]*proto.User {
	m := make(map[string]*proto.User)
	user1Pwd, err := HashPassword("password1")
	if err != nil {
		panic(err)
	}
	user1 := &proto.User{UserPwdHash: user1Pwd, UserName: "user1", Id: "1", FamilyName: "Test1", GivenName: "User1", Email: "User1@example.com"}
	// Add key/value
	m[user1.Id] = user1

	user2Pwd, err := HashPassword("password2")
	if err != nil {
		panic(err)
	}
	user2 := &proto.User{UserPwdHash: user2Pwd, UserName: "user2", Id: "2", FamilyName: "Test2", GivenName: "User2", Email: "User2@example.com"}
	// Add key/value
	m[user2.Id] = user2

	user3Pwd, err := HashPassword("password3")
	if err != nil {
		panic(err)
	}
	user3 := &proto.User{UserPwdHash: user3Pwd, UserName: "user3", Id: "3", FamilyName: "Test3", GivenName: "User3", Email: "User3@example.com"}
	// Add key/value
	m[user3.Id] = user3

	return m
}
