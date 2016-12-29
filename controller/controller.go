package controller

import (
	"encoding/json"
	"etc-pool-admin/rpc"
	"etc-pool-admin/storage"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"net/http"
	"strings"
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
var store = sessions.NewCookieStore([]byte("secret-session-store"))

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
	signedToken, err := token.SignedString([]byte(salt))
	if err != nil {
		seelog.Info("login signed token error:", err)
	}
	//可以不再使用cookie
	http.SetCookie(res, &http.Cookie{
		Name:    "Auth",
		Value:   signedToken,
		Expires: time.Now().Add(time.Hour * 24),
	})
	res.Header().Set("Access-Control-Expose-Headers", "Json-Web-Token")
	res.Header().Set("Json-Web-Token", signedToken)
	res.WriteHeader(http.StatusOK)
	session, _ := store.Get(req, signedToken)
	session.Save(req, res)
}

func PoolChartData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
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
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
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
	cookie, err := req.Cookie("Auth")
	if err != nil {
		return false, err
	}
	if cookie.Expires.Unix() > t.Unix() {
		return false, nil
	}
	str := cookie.String()
	cookiestr := strings.Split(str, "Auth=")
	session, _ := store.Get(req, str)
	if session.IsNew {
		session.Options.MaxAge = -1
		return false, nil
	}

	token, err := jwt.Parse(cookiestr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(salt), nil
	})
	if claims, ok := token.Claims.(jwt.StandardClaims); ok && token.Valid {
		if claims.ExpiresAt < t.Unix() {
			return true, nil
		}
	} else {
		return false, err
	}
	return false, nil
}
