package env

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type Env struct {
	ChannelID          string                    `json:"channelID"`
	ChannelSecret      string                    `json:"channelSecret"`
	ChannelAccessToken string                    `json:"channelAccessToken"`
	Payment            EnvPaymentLineLinePayment `json:"payment"`
	License            EnvLicenseEnvLicense      `json:"license"`
}

type EnvPaymentLineLinePayment struct {
	IsProduction  bool   `json:"isProduction"`
	BusinessNo    string `json:"businessNo"`
	MerchantID    string `json:"merchantID"`
	ChannelID     string `json:"channelID"`
	ChannelSecret string `json:"channelSecret"`
}

type EnvLicenseEnvLicense struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Price    int    `json:"price"`
}

func (env *Env) Validate() error {
	if env == nil {
		return errors.New("env is nil")
	}

	if len(env.ChannelID) == 0 {
		return errors.New("env.ChannelID is empty")
	}

	if len(env.ChannelSecret) == 0 {
		return errors.New("env.ChannelSecret is empty")
	}

	if len(env.ChannelAccessToken) == 0 {
		return errors.New("env.ChannelAccessToken is empty")
	}

	if len(env.Payment.BusinessNo) == 0 {
		return errors.New("env.Payment.BusinessNo is empty")
	}

	if len(env.Payment.MerchantID) == 0 {
		return errors.New("env.Payment.MerchantID is empty")
	}

	if len(env.Payment.ChannelID) == 0 {
		return errors.New("env.Payment.ChannelID is empty")
	}

	if len(env.Payment.ChannelSecret) == 0 {
		return errors.New("env.Payment.ChannelSecret is empty")
	}

	if len(env.License.ID) == 0 {
		return errors.New("env.License.ID is empty")
	}

	if len(env.License.Password) == 0 {
		return errors.New("env.License.Password is empty")
	}

	if env.License.Price == 0 {
		return errors.New("env.License.Price is 0")
	}

	return nil
}

func LoadEnv(t *testing.T, envRelativePath string) *Env {
	t.Helper()
	jsonFile, err := os.Open(filepath.Join(envRelativePath, "env", "env.json"))
	require.NoError(t, err)
	defer jsonFile.Close()

	env := &Env{}
	require.NoError(t, json.NewDecoder(jsonFile).Decode(env))

	return env
}

func TestLoadEnv(t *testing.T) {
	LoadEnv(t, "../")
}
