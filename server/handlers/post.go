package handlers

import (
	. "DB-forums/models"
	. "DB-forums/server/DB_requests"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func PostDetails(response http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	}

	if request.Method == "GET" {
		related := request.URL.Query().Get("related")

		result := map[string]interface{}{}
		post, err := SELECTPost_id(id)
		result["post"] = post
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find post with id: " + string(id)))
			return
		} // не нашёлся пост

		for _, elem := range strings.Split(related, ",") {
			switch elem {
			case "user":
				author, err := SELECTUser_nickname(post.Author)
				if err != nil {
					response.WriteHeader(http.StatusNotFound)
					response.Write(toMessage("Can't find user with nickname: " + post.Author))
					return
				}
				result["author"] = author
			case "thread":
				thread, err := SELECTThread_id(post.Thread)
				if err != nil {
					response.WriteHeader(http.StatusNotFound)
					response.Write(toMessage("Can't find thread with id: " + string(post.Thread)))
					return
				}
				result["thread"] = thread
			case "forum":
				forum, err := SELECTForum_slug(post.Forum)
				if err != nil {
					response.WriteHeader(http.StatusNotFound)
					response.Write(toMessage("Can't find forum with slug: " + post.Forum))
					return
				}
				result["forum"] = forum
			}
		}
		body, err := json.Marshal(result)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	} else { // method = POST. (изменить сообщение)
		var post Post
		err := json.NewDecoder(request.Body).Decode(&post)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(toMessage("Bad request"))
			return
		} // раскордировали запрос
		message := post.Message

		post, err = SELECTPost_id(id)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find post with id: " + string(id)))
			return
		} // не нашёлся пост

		if message != "" && message != post.Message {
			post, err = UPDATEPost_id(id, message)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
				return
			} // ошибка в запросе в БД
		}

		// пост добавился
		body, err := json.Marshal(post)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	}
}
