package main

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/yvasiyarov/gorelic"
	"etc-pool-admin/controller"
	"etc-pool-admin/routers"
	"etc-pool-admin/rpc"
	"etc-pool-admin/storage"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func readConfig(cfg *controller.Config) {
	configFileName := "manage.json"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	configFileName, _ = filepath.Abs(configFileName)
	seelog.Infof("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		seelog.Critical("File error: ", err.Error())
		panic(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&cfg); err != nil {
		seelog.Critical("Config error: ", err.Error())
		panic(err.Error())
	}
}

func startNewrelic(cfg controller.Config) {
	if cfg.NewrelicEnabled {
		nr := gorelic.NewAgent()
		nr.Verbose = cfg.NewrelicVerbose
		nr.NewrelicLicense = cfg.NewrelicKey
		nr.NewrelicName = cfg.NewrelicName
		nr.Run()
	}
}

func main() {
	//recover function
	defer func() {
		if r := recover(); r != nil {
			seelog.Critical("Critical error, recover:", r)
		}
	}()

	//load log config
	logger, err := seelog.LoggerFromConfigAsFile("./logconfig.xml")
	if err != nil {
		seelog.Critical("load log config fail:", err)
		panic(err)
	}
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()

	readConfig(&controller.Conf)

	if controller.Conf.Threads > 0 {
		runtime.GOMAXPROCS(controller.Conf.Threads)
		seelog.Infof("Running with %v threads", controller.Conf.Threads)
	}
	startNewrelic(controller.Conf)

	controller.Backend = storage.NewRedisClient(&controller.Conf.Redis, controller.Conf.Coin)
	controller.RpcClient = rpc.NewRPCClient(controller.Conf.Upstream.Url, controller.Conf.Upstream.Timeout)

	//main function
	err = http.ListenAndServe(controller.Conf.Listen, router.Routes)
	if err != nil {
		seelog.Critical("listen error:", err)
		panic(err)
	}
}
