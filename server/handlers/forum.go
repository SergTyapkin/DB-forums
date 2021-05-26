package handlers

import (
	. "DB-forums/models"
	. "DB-forums/server/DB_requests"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ForumCreate(response http.ResponseWriter, request *http.Request) {
	var forum Forum
	err := json.NewDecoder(request.Body).Decode(&forum)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	} // раскордировали запрос

	user, err := SELECTUser_nickname(forum.User)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find user with nickname: " + forum.User))
		return
	} // не нашёлся пользователь
	forum.User = user.Nickname

	insertedForum, err := INSERTForum(forum)
	if err != nil { // если форум не добавился
		forumExisting, err := SELECTForum_slug(forum.Slug)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
			return
		} // форума такого и нет, значит ошибка в запросе в БД

		body, err := json.Marshal(forumExisting)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusConflict)
		response.Write(body)
		return
		// такой форум уже есть
	}

	//INSERTForumToUser(insertedForum.Slug, user)
	// форум добавился
	body, err := json.Marshal(insertedForum)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Write(body)
}

func ForumDetails(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug"]

	forumExisting, err := SELECTForum_slug(slug)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find forum with slug: " + slug))
		return
	} // форум не нашёлся

	body, err := json.Marshal(forumExisting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(body)
	// форум нашёлся
}

func ThreadCreate(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug"]
	var thread Thread
	err := json.NewDecoder(request.Body).Decode(&thread)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	} // раскордировали запрос
	thread.Forum = slug

	user, err := SELECTUser_nickname(thread.Author)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find user with nickname: " + thread.Author))
		return
	} // не нашёлся пользователь

	forum, err := IncrementForumThreads_slug(slug)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find forum with slug: " + slug))
		return
	} // не нашёлся форум
	thread.Forum = forum.Slug

	insertedThread, err := INSERTThread(thread)
	if err != nil { // если ветка не добавилась
		structureExisting, err := SELECTThread_slug(thread.Slug.String)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
			return
		} // ветки такой и нет, значит ошибка в запросе в БД

		body, err := json.Marshal(structureExisting)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusConflict)
		response.Write(body)
		return
		// такая ветка уже есть
	}

	INSERTForumToUser(insertedThread.Forum, user)
	// ветка добавилась
	body, err := json.Marshal(insertedThread)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Write(body)
}

func ForumUsers(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug"]

	query := request.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 100
	}
	since := query.Get("since")
	desc, err := strconv.ParseBool(query.Get("desc"))
	if err != nil {
		desc = false
	}

	_, err = SELECTForum_slug(slug)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find forum with slug: " + slug))
		return
	} // форум не нашёлся

	users, err := SELECTForumUsers(slug, limit, since, desc)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
		return
	}

	if len(users) == 0 { // пользователи не выбрались
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("[]"))
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(body)
	// нашли и выдали пользователей
}

func ForumThreads(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug"]

	query := request.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 100
	}
	since := query.Get("since")
	desc, err := strconv.ParseBool(query.Get("desc"))

	if err != nil {
		desc = false
	}

	forum, err := SELECTForum_slug(slug)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find forum forum with slug: " + slug))
		return
	} // форум не нашёлся
	slug = forum.Slug

	structs, err := SELECTForumThreads(slug, limit, since, desc)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
		return
	} // ошибка в запросе в БД

	if len(structs) == 0 { // ветки не выбрались
		response.WriteHeader(http.StatusOK)
		response.Write([]byte("[]"))
		return
	}

	body, err := json.Marshal(structs)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(body)
	// нашли и выдали ветки
}
