package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"etc-pool-admin/storage"
	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	queryByAccount = "account"
	queryByAddress = "walletAddress"
	queryByPhone   = "phone"
	queryByEmail   = "email"
)

type userAllInfo struct {
	BasicInfo   storage.UserInfo
	LastRevenue int64
	AllRevenue  int64
	OfflineTime string
}

//QueryUsers get givin pattern correspodding users and info
func QueryUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", "GET")
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	param := mux.Vars(req)["value"]

	querykey := regexpParam(param)
	basicInfo := storage.GetUserInfo(querykey, param, Conf.Mongo)
	if len(basicInfo) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	allUserInfo := make([]userAllInfo, len(basicInfo))
	for idx, basic := range basicInfo {
		allUserInfo[idx].BasicInfo = basic
		allUserInfo[idx].LastRevenue, allUserInfo[idx].AllRevenue = Backend.GetRevenue(basic.Wallet)
		allUserInfo[idx].OfflineTime = storage.GetOnlineStat(basic.Wallet, Conf.Mongo)
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(allUserInfo)
	if err != nil {
		seelog.Error("query users error, when serializing response data:", err)
	}
}

func SetUserInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	modify := storage.UserInfo{}
	body, _ := ioutil.ReadAll(req.Body)

	err = json.Unmarshal(body, &modify)
	if err != nil {
		seelog.Info("decode requst body error:", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// decoder := json.NewDecoder(req.Body)
	// err = decoder.Decode(&modify)
	// if err != nil {
	// 	seelog.Info("decode requst body error:", err)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	err = storage.SetUserInfo(modify, Conf.Mongo)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func Fee(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodPut+","+http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	if req.Method == http.MethodPut {
		setFee(res, req)
		return
	}

	fees := storage.GetFee(Conf.Mongo)

	if fees.Norm == 0 {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(fees)
	if err != nil {
		seelog.Info("serializing error:", err)
	}
}

//SetFee set common users or vip users fee
func setFee(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		seelog.Info("read request body error:", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	fees := storage.FeeInfo{}
	err = json.Unmarshal(body, &fees)
	if err != nil {
		seelog.Info("decode request body error, ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.SetFee(fees, Conf.Mongo)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

//GetAddressChartData get 24 housrs hashrate from redis
func GetAddressChartData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}
	address := mux.Vars(req)["address"]
	address = strings.ToLower(address)
	chartdata, err := Backend.GetAccountChartValues(address)
	if err != nil {
		seelog.Info("get wallet address chart data error:", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(chartdata)
	if err != nil {
		seelog.Error("query chart data error, when serializing response data:", err)
	}
}

//GetAddressStaticData get revenue data include matured revenue, immature revenue, last paid, total paid
func GetAddressStaticData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	address := mux.Vars(req)["address"]
	address = strings.ToLower(strings.TrimSpace(address))
	data := Backend.GetWalletRevenue(address)
	if data == nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(data)
	if err != nil {
		seelog.Info("serializing response data error:", err)
	}
}

//GetMinersInfo get miner state about online offline and its hashrate data
func GetMinersInfo(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	address := mux.Vars(req)["address"]
	address = strings.ToLower(strings.TrimSpace(address))
	stats, err := Backend.GetWorkersStats(address)
	if err != nil {
		seelog.Info("cannot get worker stats")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
	err = json.NewEncoder(res).Encode(stats)
	if err != nil {
		seelog.Info("serializing response data error:", err)
	}
}

func QueryPaymentHistory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}
	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	address := mux.Vars(req)["address"]
	begin := mux.Vars(req)["begintime"]
	end := mux.Vars(req)["endtime"]
	address = strings.ToLower(strings.TrimSpace(address))
	btime, err := strconv.ParseInt(begin, 10, 64)
	if err != nil {
		seelog.Info("convert string time to int error:", err)
	}
	etime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		seelog.Info("convert string time to int error:", err)
	}

	res.WriteHeader(http.StatusOK)
	payments := Backend.GetPaymentHistory(address, btime, etime)
	if payments == nil {
		seelog.Info("no corresponding payments data")
	}
	err = json.NewEncoder(res).Encode(payments)
	if err != nil {
		seelog.Info("serializing payments history response data error:", err)
	}
}

func ExportsPayments(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/csv")
	res.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method == http.MethodOptions {
		res.Header().Set("Access-Control-Allow-Headers", "Json-Web-Token")
		res.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		res.WriteHeader(http.StatusOK)
		return
	}

	pass, err := validate(req)
	if err != nil || pass == false {
		seelog.Info("validate err:", err)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	address := mux.Vars(req)["address"]
	address = strings.ToLower(strings.TrimSpace(address))
	res.Header().Set("Access-Control-Expose-Headers", "Payments-Filename")
	res.Header().Set("Payments-Filename", address+".csv")
	begin := mux.Vars(req)["begintime"]
	end := mux.Vars(req)["endtime"]
	btime, err := strconv.ParseInt(begin, 10, 64)
	if err != nil {
		seelog.Info("convert string time to int error:", err)
	}
	etime, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		seelog.Info("convert string time to int error:", err)
	}

	res.Header().Set("Content-Disposition", "attachment;filename="+address+".csv")

	payments := Backend.GetPaymentHistory(address, btime, etime) //begintime and end time is time stamp, dont use string
	if payments == nil {
		seelog.Info("no corresponding payments data")
	}

	b := bytes.Buffer{}
	writer := csv.NewWriter(&b)

	writer.Write([]string{"date", "tx", "amount"})
	for _, payment := range payments {
		stamp := time.Unix(payment["timestamp"].(int64), 0).Format("2006-01-02 15:04:05")
		tx := payment["tx"].(string)
		amount := strconv.FormatFloat(float64(payment["amount"].(int64))/1e9, 'f', 4, 32)
		writer.Write([]string{stamp, tx, amount})
	}

	res.Header().Set("contentlength", strconv.Itoa(b.Len()))
	writer.Flush()

	res.WriteHeader(http.StatusOK)
	res.Write(b.Bytes())
	//TODO if file too big , fix this bug
}

func regexpParam(param string) string {
	//regexp wallet address
	regAddr := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	regPhone := regexp.MustCompile(`^1[3|4|5|8][0-9]{9}$`)
	regEmail := regexp.MustCompile(`^[a-z_0-9.-]{1,64}@([a-z0-9-]{1,200}.){1,5}[a-z]{1,6}$`)

	if regAddr.MatchString(param) {
		return queryByAddress
	} else if regPhone.MatchString(param) {
		return queryByPhone
	} else if regEmail.MatchString(param) {
		return queryByEmail
	}

	return queryByAccount
}
