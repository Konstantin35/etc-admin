package router

import (
	"etc-pool-admin/controller"
	"github.com/gorilla/mux"
)

var Routes *mux.Router

func init() {
	Routes = mux.NewRouter()
	Routes.HandleFunc("/login", controller.Login).Methods("POST")
	Routes.HandleFunc("/main/poolchart", controller.PoolChartData).Methods("GET")
	Routes.HandleFunc("/main/statistic", controller.StatisticData).Methods("GET")
	Routes.HandleFunc("/management/users", controller.UserManage)
}
