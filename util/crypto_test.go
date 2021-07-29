package util

import (
	"testing"

	"gitlab.pri.ibanyu.com/middleware/dbinjection/config"
	"gitlab.pri.ibanyu.com/middleware/dbinjection/util/logger"
)

func TestCryPto(t *testing.T) {
	logger.InitLog(".", "test.log", "debug")
	config.InitConfig("../config/config.yml")

	f := "TestCryPto"
	pwd := "aaaaaa"

	cryptoPwd, err := AesCrypto([]byte(pwd))
	if err != nil {
		t.Errorf("%s crypto err: %s", f, err.Error())
		t.FailNow()
	}
	deCryptoPwd, err := AesDeCrypto(cryptoPwd)
	if err != nil {
		t.Errorf("%s decrypto err: %s", f, err.Error())
		t.FailNow()
	}

	if pwd != string(deCryptoPwd) {
		t.Errorf("%s failed, not equal", f)
		t.FailNow()
	}
}
