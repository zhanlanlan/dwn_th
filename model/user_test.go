package model

import (
	"context"
	"log"
	"testing"
)

func TestCreateUser(t *testing.T) {
	InitDB()

	user := User{
		UserName: "lianghao",
		PassWord: "fhwqiofwhqiofqwio",
	}

	err := CreateUser(context.Background(), &user)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("%+v\n", user)
}

func TestDeleteUser(t *testing.T) {
	InitDB()

	userName := "lianghao"
	err := DeleteUser(context.Background(), userName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// log.Printf("%+v\n", user)
}

func TestUpdateUser(t *testing.T) {
	InitDB()

	userName := "tangpengfei"
	password := "aaaaaaaaaaaa"
	err := UpdateUserPassWord(context.Background(), userName, password)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// log.Printf("%+v\n", user)
}

func TestRune(t *testing.T) {
	InitDB()

	var u User

	err := DB.Table("t_user").Where("user_name = ?", "fuck").Find(&u).Error
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("%+v \n", u)

}

// func TestGenerateToken(t *testing.T) {
// 	t := jwt.NewWithClaims(jwt.SigningMethodES256, Claim{
// 		User: user,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(time.Hour).Unix(),
// 		},
// 	})

// 	return t.SignedString([]byte("fuck"))

// }
