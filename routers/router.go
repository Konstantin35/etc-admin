package router

import (
	"github.com/gorilla/mux"
	"manageserver/controller"
)

var Routes *mux.Router

func init() {
	Routes = mux.NewRouter()
	Routes.HandleFunc("/login", controller.Login).Methods("POST", "OPTIONS")
	Routes.HandleFunc("/", controller.MainPage).Methods("GET")
	Routes.HandleFunc("/main/poolchart", controller.PoolChartData).Methods("GET")
	Routes.HandleFunc("/main/statistic", controller.StatisticData).Methods("GET")
	Routes.HandleFunc("/management/users", controller.UserManage)
}
