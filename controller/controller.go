package controller

import (
	"encoding/json"
	"errors"
	"etc-pool-admin/rpc"
	"etc-pool-admin/storage"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type Config struct {
	Coin    string `json:"coin"`
	Threads int    `json:"threads"`
	Listen  string `json:"listen"`

	Upstream rpc.Upstream `json:"upstream"`

	Mongo storage.MongoConfig `json:"mongo"`
	Redis storage.RedisConfig `json:"redis"`

	BaseAddress string `json:"address"`

	NewrelicName    string `json:"newrelicName"`
	NewrelicKey     string `json:"newrelicKey"`
	NewrelicVerbose bool   `json:"newrelicVerbose"`
	NewrelicEnabled bool   `json:"newrelicEnabled"`
}

const salt = "pool-manage-server-by-Gavin"

var Conf Config
var Backend *storage.RedisClient
var RpcClient *rpc.RPCClient

//Login get user name and pwd, return a token
func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Credentials", "true")

	req.ParseForm()

	username := req.FormValue("username")
	pwd := req.FormValue("password")
	//validate username and pwd
	checked := storage.CheckUserAdmin(username, pwd, Conf.Mongo)
	if checked == false {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    username,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	seelog.Info("login token:", token)
	signedToken, err := token.SignedString([]byte(salt))
	if err != nil {
		seelog.Info("login signed token error:", err)
	}
	seelog.Info("login signed token:", signedToken)
	res.Header().Set("Access-Control-Expose-Headers", "Json-Web-Token")
	res.Header().Set("Json-Web-Token", signedToken)
	res.WriteHeader(http.StatusOK)
}

func PoolChartData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
	if req.Method == "OPTIONS" {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", "GET")
		res.WriteHeader(http.StatusOK)
		return
	}
	chartdata := make(map[string]interface{})
	poolhashs, err := Backend.GetPoolChartData()
	if err != nil {
		seelog.Error("get pool hash error:", err)
	}
	middle := make(map[string]interface{})
	length := len(poolhashs)
	for i := 0; i < length/2; i++ {
		middle = poolhashs[i]
		poolhashs[i] = poolhashs[length-i-1]
		poolhashs[length-i-1] = middle
	}
	res.WriteHeader(http.StatusOK)
	chartdata["poolhashs"] = poolhashs
	err = json.NewEncoder(res).Encode(chartdata)
	if err != nil {
		seelog.Error("error serializing response data:", err)
	}
}

func StatisticData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err, pass)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	if req.Method == "OPTIONS" {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", "GET")
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusOK)
	//TODO change poolbalance value to /10*9
	statistis := Backend.GetMainStatisticData()
	poolbalance, err := RpcClient.GetAccountBalance(Conf.BaseAddress)
	if err != nil {
		seelog.Info("cannot get base address balance coin:", err)
	}
	statistis["poolbalance"] = poolbalance
	err = json.NewEncoder(res).Encode(statistis)
	if err != nil {
		seelog.Error("error when get main page statistic data, ", err)
	}
}

//Validate using for check token
func validate(req *http.Request) (bool, error) {
	t := time.Now()
	webtoken := req.Header.Get("Json-Web-Token")
	if webtoken == "" {
		return false, errors.New("cannot get jwt when validate")
	}
	seelog.Info("token:", webtoken)

	token, err := jwt.Parse(webtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(salt), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		seelog.Info("token valid:", token.Valid, "claim:", claims)
		if claims.VerifyExpiresAt(t.Unix(), true) {
			return true, nil
		}
	}
	return false, err
}
