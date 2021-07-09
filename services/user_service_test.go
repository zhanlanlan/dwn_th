package services

import (
	"context"
	"dwn_th/model"
	"dwn_th/proto"
	"log"
	"testing"
)

func TestCreateUser(t *testing.T) {
	model.InitDB()

	err := CreateUser(context.Background(), proto.CreateUserREQ{
		UserName: "唐鹏飞1",
		PassWord: "123456",
	})

	log.Printf("%+v \n", err)
}
