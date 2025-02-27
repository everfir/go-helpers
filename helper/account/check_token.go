package account

import (
	"context"
	"fmt"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/internal/service/account"
)

func CheckToken(ctx context.Context, token string) (nctx context.Context, err error) {
	if token == "" {
		err = fmt.Errorf("[go-helper] Invalid Token. should not be empty")
		return ctx, err
	}

	// 访问 check token 接口
	accountInfo, err := account.CheckToken(ctx, token)
	if err != nil {
		err = fmt.Errorf("[go-helper] call account service failed:%w", err)
		return ctx, err
	}

	if !accountInfo.Valid {
		err = fmt.Errorf("[go-helper] token has expired")
		return ctx, err
	}

	// 将用户信息存储到 Context 中
	nctx = context.WithValue(ctx, consts.AccountInfoKey, accountInfo.AccountInfo)
	return nctx, nil
}
