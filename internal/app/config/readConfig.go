package config

import (
	"os"

	"github.com/IceCodeNew/telesend/pkg/fsHelper"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	k        = koanf.New(".")
	TSConfig *TelesendConfig
)

type TelesendConfig struct {
	// bot token is supposed to be set in the environment variable
	BotToken string
	DBPath   string `koanf:"db_path"`
	Verbose  bool   `koanf:"verbose"`
}

func (tsConfig *TelesendConfig) sanityCheck() error {
	// req, err := http.NewRequest(http.MethodGet,
	// 	fmt.Sprintf("https://api.telegram.org/bot%s/getMe", tsConfig.BotToken), nil)
	// var resp *http.Response

	// _, _, err = lo.AttemptWhileWithDelay(3, time.Second*10,
	// 	func(int, time.Duration) (error, bool) {
	// 		// BE AWARE as the resp is NOT GUARANTEED to be non-nil
	// 		resp, err = httpHelper.HttpReqHelper(req, tsConfig.Verbose)
	// 		if err != nil {
	// 			return err, true
	// 		}
	// 		return nil, false
	// 	})
	// if err != nil {
	// 	return fmt.Errorf("FATAL: the bot token failed check after 3 attempts, the last error was:\n %v", err)
	// }
	// defer func() {
	// 	_ = resp.Body.Close()
	// }()

	return fsHelper.CreateDir(tsConfig.DBPath)
}

func ReadConfig() error {
	telesendConfPath := "telesend.json"
	if _path, found := os.LookupEnv("TELESEND_CONF_PATH"); found {
		if regularFile, err := fsHelper.IsRegularFile(_path); err != nil {
			return err
		} else if regularFile {
			telesendConfPath = _path
		}
	}

	// default configuration MUST exist
	if err := k.Load(file.Provider(telesendConfPath), json.Parser()); err != nil {
		return err
	}

	k.Unmarshal("", &TSConfig)
	TSConfig.BotToken = os.Getenv("TELESEND_BotToken")
	return TSConfig.sanityCheck()
}
