package handlers

import (
	. "DB-forums/server/DB_requests"
	"encoding/json"
	"net/http"
)

/*
func InitDB(response http.ResponseWriter, request *http.Request) {
	str, err := InitDB_all()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage(str + ". Error: " + err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(toMessage(str))
}
*/

func ClearDB(response http.ResponseWriter, request *http.Request) {
	str, err := ClearDB_all()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage(str + ". Error: " + err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(toMessage(str))
}

/*
func DropDB(response http.ResponseWriter, request *http.Request) {
	str, err := DropDB_all()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage(str + ". Error: " + err.Error()))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(toMessage(str))
}
*/
/*---------TEST--------
func TestInitDB() {
	str, err := InitDB_all()
	if err != nil {
		fmt.Println(err)
		return
	} // ошибка в запросе в БД
	fmt.Println(str)
}

func TestDropDB() {
	str, err := DropDB_all()
	if err != nil {
		fmt.Println(err)
		return
	} // ошибка в запросе в БД
	fmt.Println(str)
}
//---------TEST--------*/

func StatusDB(response http.ResponseWriter, request *http.Request) {
	stats, err := GetStats()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Invalid DB request"))
		return
	}

	body, err := json.Marshal(stats)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(body)
}
