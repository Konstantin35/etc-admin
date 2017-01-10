package router

import (
	"etc-pool-admin/controller"
	"github.com/gorilla/mux"
	"net/http"
)

var Routes *mux.Router

func init() {
	Routes = mux.NewRouter()
	Routes.NotFoundHandler = http.HandlerFunc(controller.NotFound)
	Routes.HandleFunc("/login", controller.Login).Methods("POST")
	Routes.HandleFunc("/main/poolchart", controller.PoolChartData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/main/statistic", controller.StatisticData).Methods("GET", "OPTIONS")
	/****************routers for user manage********************/
	Routes.HandleFunc("/user/query/{value}", controller.QueryUsers).Methods("GET", "OPTIONS") //query by wallet address or login account or email or phone number
	Routes.HandleFunc("/user/info/settings", controller.SetUserInfo).Methods("PUT", "OPTIONS")
	Routes.HandleFunc("/user/info/common/fees", controller.Fee).Methods("GET", "OPTIONS", "PUT") //set common user and vip user pool fee
	Routes.HandleFunc("/user/data/chart/{address:0x[0-9a-fA-F]{40}}", controller.GetAddressChartData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/data/statistic/{address:0x[0-9a-fA-F]{40}}", controller.GetAddressStaticData).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/info/miners/{address:0x[0-9a-fA-F]{40}}", controller.GetMinersInfo).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/history/payments/{address:0x[0-9a-fA-F]{40}}/{begintime}/{endtime}", controller.QueryPaymentHistory).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/user/history/payments/exports/{address:0x[0-9a-fA-F]{40}}/{begintime}/{endtime}", controller.ExportsPayments).Methods("GET", "OPTIONS")
	/***************routers for pool manage*********************/
	Routes.HandleFunc("/pool/balance", controller.GetPoolBalance).Methods("GET", "OPTIONS")
	Routes.HandleFunc("/pool/account/transaction", controller.SendTransaction).Methods("POST", "OPTIONS")
	Routes.HandleFunc("/pool/history/payments/{begintime}/{endtime}", controller.PoolPaymentHistory).Methods("GET", "OPTIONS")
}
