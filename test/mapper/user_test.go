package mapper_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/OJ-lab/oj-lab-services/service/mapper"
	"github.com/OJ-lab/oj-lab-services/service/model"
)

func TestUserMapper(t *testing.T) {
	user := model.User{
		Account:  "test",
		Password: func() *string { s := "test"; return &s }(),
		Roles:    []*model.Role{{Name: "admin"}},
	}
	err := mapper.CreateUser(user)
	if err != nil {
		t.Error(err)
	}

	dbUser, err := mapper.GetUser(user.Account)
	if err != nil {
		t.Error(err)
	}
	userJson, err := json.MarshalIndent(dbUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(userJson))

	dbPublicUser, err := mapper.GetPublicUser(user.Account)
	if err != nil {
		t.Error(err)
	}
	publicUserJson, err := json.MarshalIndent(dbPublicUser, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", string(publicUserJson))
}
