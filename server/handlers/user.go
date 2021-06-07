package handlers

import (
	. "DB-forums/models"
	. "DB-forums/server/DB_requests"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func UserCreate(response http.ResponseWriter, request *http.Request) {
	nickname := mux.Vars(request)["nickname"]

	var user User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	} // раскодировали запрос
	user.Nickname = nickname

	foundedUser, err := INSERTUser(user)
	if err != nil { // если пользователь не добавился
		var usersExisting []User
		foundedUser, err = SELECTUser_nickname(nickname) // пытаемся найти по нику
		if err == nil {
			usersExisting = append(usersExisting, foundedUser)
		}
		foundedUser, err = SELECTUser_email(user.Email) // дальше по емейлу
		if err == nil && foundedUser.Nickname != nickname {
			usersExisting = append(usersExisting, foundedUser)
		}
		if len(usersExisting) == 0 {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
			return
		} // пользователя такого и нет, значит ошибка в запросе в БД

		// такой пользователь уже есть
		body, err := json.Marshal(usersExisting)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusConflict)
		response.Write(body)
		return
	}
	// пользователь добавился
	body, err := json.Marshal(foundedUser)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Write(body)
}

func UserProfile(response http.ResponseWriter, request *http.Request) {
	nickname := mux.Vars(request)["nickname"]

	if request.Method == "GET" {
		user, err := SELECTUser_nickname(nickname)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find user with nickname: " + nickname))
			return
		} // не нашёлся пользователь

		// пользователь обновился
		body, err := json.Marshal(user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	} else { // request.Method = "POST". Обновляем пользователя
		var user User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(toMessage("Bad request"))
			return
		} // раскодировали запрос
		if user.Nickname == "" {
			user.Nickname = nickname
		} else if user.Nickname != nickname {
			_, err := SELECTUser_nickname(user.Nickname)
			if err == nil {
				response.WriteHeader(http.StatusConflict)
				response.Write(toMessage("User already exists with nickname: " + user.Nickname))
				return
			} // не нашёлся пользователь
		}

		if user.Email != "" {
			tmpUser, err := SELECTUser_email(user.Email)
			if err == nil && tmpUser.Nickname != nickname {
				response.WriteHeader(http.StatusConflict)
				response.Write(toMessage("User already exists with email: " + user.Email))
				return
			} // не нашёлся пользователь
		}

		user, err = UPDATEUser_nickname(nickname, user.Nickname, user.Name, user.About, user.Email)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find user with nickname: " + nickname))
			return
		} // не нашёлся пользователь

		// пользователь обновился
		body, err := json.Marshal(user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	}
}
