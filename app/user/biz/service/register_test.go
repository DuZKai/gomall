package service

import (
	"context"
	"github.com/joho/godotenv"
	"gomall/app/user/biz/dal/mysql"
	user "gomall/rpc_gen/kitex_gen/user"
	"testing"
)

func TestRegister_Run(t *testing.T) {
	godotenv.Load("../../.env")
	mysql.Init()
	ctx := context.Background()
	s := NewRegisterService(ctx)
	// init req and assert value

	req := &user.RegisterReq{
		Email:           "1demo@damin.com",
		Password:        "FJODIAFUFJO",
		PasswordConfirm: "FJODIAFUFJO",
	}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

}
