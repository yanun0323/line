package linepay

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yanun0323/line/internal/env"
)

func TestPayment(t *testing.T) {
	env := env.LoadEnv(t, "../")
	lp := NewLinePay(env.Payment.IsProduction, env.Payment.ChannelID, env.Payment.ChannelSecret)
	tsID := 0
	{
		res, err := lp.RequestPayment(
			context.Background(),
			70,
			strconv.FormatInt(time.Now().UnixMilli(), 10),
			"product_123456789",
			"monthly plan",
			"https://localhost:3000/icon.png",
			1,
			70,
			70,
			"https://localhost:3000/payment/success",
			"https://localhost:3000/payment/cancel",
		)
		require.NoError(t, err)
		require.NotEmpty(t, res)
		t.Logf("res: %+v", res)

		tsID = res.Info.TransactionID
	}

	{
		res, err := lp.ConfirmPayment(context.Background(), strconv.Itoa(tsID), 50)
		require.NoError(t, err)
		require.NotEmpty(t, res)
		t.Logf("res: %+v", res)
	}
}
