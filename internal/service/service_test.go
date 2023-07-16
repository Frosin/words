package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	data = `{
		"bot_token": "test",
		"yadisk_token": "test2"
	}`

	testFile = "../test.json"
)

func Test_InitConfig_DefaultFileName(t *testing.T) {
	err := os.WriteFile(defaultConfigFile, []byte(data), 0777)
	assert.NoError(t, err)
	defer os.Remove(defaultConfigFile)

	sc := SConfig{}
	err = sc.initConfig("")
	assert.NoError(t, err)

	assert.NotEmpty(t, sc.BotToken)
	assert.NotEmpty(t, sc.YadiskToken)
}

func Test_InitConfig_CustomFile(t *testing.T) {
	err := os.WriteFile(testFile, []byte(data), 0777)
	assert.NoError(t, err)
	defer os.Remove(testFile)

	sc := SConfig{}
	err = sc.initConfig(testFile)
	assert.NoError(t, err)

	assert.Equal(t, "test", sc.BotToken)
	assert.NotEmpty(t, "test2", sc.YadiskToken)
}
