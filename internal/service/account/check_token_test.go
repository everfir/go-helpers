package account_test

import (
	"context"
	"testing"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/internal/service/account"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/zeebo/assert"
)

func TestCheckToken(t *testing.T) {
	// 注入一些数据
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.BusinessKey, "test")
	ctx = context.WithValue(ctx, consts.VersionKey, "1.0.0")
	ctx = context.WithValue(ctx, consts.PlatformKey, "ios")
	ctx = context.WithValue(ctx, consts.DeviceKey, "1234567890")
	ctx = context.WithValue(ctx, consts.AppTypeKey, "app")

	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOjUwNSwiYXZhdGFyIjoiaHR0cHM6Ly90aGlyZHd4LnFsb2dvLmNuL21tb3Blbi92aV8zMi9QT2dFd2g0bUlITzRuaWJIMEtsTUVDTmpqR3hRVXEyNFpFYUdUNHBvQzZpY1JpY2NWR0tTeVh3aWJjUHE0QldtaWFJR3VHMWljd3hhUVg2Z3JDOVZlbVpvSjhyZy8xMzIiLCJjaGFubmVsIjoid2VhcHAiLCJjdGltZSI6MTczNzQ3NTIwMCwiZW1haWwiOiJva2Z2dTZ4bHoxdzYxT0Y3SENjMnNyV1lSdzJRIiwiZXhwIjoxNzM5Njc2OTgxLCJuaWNrbmFtZSI6IuW-ruS_oeeUqOaItyIsInBob25lX251bSI6IiIsInJvbGUiOjEsInVzZXJuYW1lIjoib2tmdnU2eGx6MXc2MU9GN0hDYzJzcldZUncyUSIsInZpcF9leHBpcmVfdGltZSI6NDI5NDk2NzI5NX0.XzHAnJvW3uEMo4eH3PUL5eVT7MCHOCf9g6qlPc3nhak`
	resp, err := account.CheckToken(ctx, token)
	assert.NoError(t, err)
	logger.Info(ctx, "check token", field.Any("resp", resp))
}
