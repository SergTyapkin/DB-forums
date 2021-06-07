package handlers

import (
	. "DB-forums/models"
	. "DB-forums/server/DB_requests"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"net/http"
	"strconv"
)

func PostsCreate(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug_or_id"]
	id, err := strconv.Atoi(slug)
	useId := true
	if err != nil {
		useId = false
	}

	var posts []Post
	err = json.NewDecoder(request.Body).Decode(&posts)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	} // раскодировали запрос

	var thread Thread
	if useId {
		thread, err = SELECTThread_id(id)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with id: " + string(id)))
			return
		} // не нашлась ветка
	} else {
		thread, err = SELECTThread_slug(slug)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with slug: " + slug))
			return
		} // не нашлась ветка
	}

	// Перенесено в триггер
	/*
		if post.Parent != 0 {
			parentPost, err := SELECTPost_id(post.Parent)
			if err != nil {
				response.WriteHeader(http.StatusConflict)
				response.Write(toMessage("Can't find post with id: " + string(post.Parent)))
				return
			} // родительского поста не нашлось
			if parentPost.Thread != post.Thread {
				response.WriteHeader(http.StatusConflict)
				response.Write(toMessage("Parent post is in another thread: " + string(parentPost.Thread)))
				return
			} // родительский пост в другой ветке
		} */

	// Перенесено в триггер
	/* user, err := SELECTUser_nickname(post.Author)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find user with id: " + post.Author))
		return
	} // пользователя не нашлось */

	posts, err = INSERTPosts(posts, thread.Id, thread.Forum)
	if err != nil {
		println(err.Error())
		if pgErr, ok := err.(pgx.PgError); ok {
			if pgErr.Code == "00228" {
				response.WriteHeader(http.StatusConflict)
				response.Write(toMessage("Parent post is in another thread"))
			} else {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Can't find user with current id or Can't find post with current id"))
			}
		}
		return
	}

	postsLen := len(posts)
	if postsLen == 0 { // посты не выбрались
		response.WriteHeader(http.StatusCreated)
		response.Write([]byte("[]"))
		return
	}

	// посты добавились
	_, err = AddForumPosts_slug(thread.Forum, postsLen)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
		return
	} // ошибка в запросе в БД

	body, err := json.Marshal(posts)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusCreated)
	response.Write(body)
}

func ThreadDetails(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug_or_id"]
	id, err := strconv.Atoi(slug)
	useId := true
	if err != nil {
		useId = false
	}

	if request.Method == "GET" {
		var thread Thread
		if useId {
			thread, err = SELECTThread_id(id)
			if err != nil {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Can't find thread with id: " + string(id)))
				return
			} // не нашлась ветка
		} else {
			thread, err = SELECTThread_slug(slug)
			if err != nil {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Can't find thread with slug: " + slug))
				return
			} // не нашлась ветка
		}
		// ветка выбралась
		body, err := json.Marshal(thread)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	} else { // request.Method = "POST". Обновляем ветку
		var thread Thread
		err = json.NewDecoder(request.Body).Decode(&thread)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			response.Write(toMessage("Bad request"))
			return
		} // раскодировали запрос

		if useId {
			thread, err = UPDATEThread_id(id, thread.Title, thread.Message)
			if err != nil {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Can't find thread with id: " + string(id)))
				return
			} // не нашлась ветка
		} else {
			thread, err = UPDATEThread_slug(slug, thread.Title, thread.Message)
			if err != nil {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Can't find thread with slug: " + slug))
				return
			} // не нашлась ветка
		}

		// ветка обновилась
		body, err := json.Marshal(thread)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(toMessage("Can't marshal JSON file"))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	}
}

func ThreadPosts(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug_or_id"]
	id, err := strconv.Atoi(slug)
	useId := true
	if err != nil {
		useId = false
	}

	query := request.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 100
	}
	since := query.Get("since")
	sort := query.Get("sort")
	desc, err := strconv.ParseBool(query.Get("desc"))
	if err != nil {
		desc = false
	}

	var structs []Post
	if useId {
		// Без этого SELECT выдаст не ошибку, а []
		_, err := SELECTThread_id(id)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with id: " + string(id)))
			return
		} // ветка не нашлась
	} else {
		tmpThread, err := SELECTThread_slug(slug)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with slug: " + slug))
			return
		} // ветка не нашлась
		id = tmpThread.Id
	}

	structs, err = SELECTThreadPosts_id(id, limit, since, sort, desc)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find thread with id: " + string(id)))
		return
	} // ошибка в запросе в БД

	if len(structs) == 0 { // посты не выбрались
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
	// нашли и выдали посты
}

func VoteCreate(response http.ResponseWriter, request *http.Request) {
	slug := mux.Vars(request)["slug_or_id"]
	id, err := strconv.Atoi(slug)
	useId := true
	if err != nil {
		useId = false
	}

	var vote Vote
	err = json.NewDecoder(request.Body).Decode(&vote)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write(toMessage("Bad request"))
		return
	} // раскодировали запрос

	var thread Thread
	if useId {
		thread, err = SELECTThread_id(id)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with id: " + string(id)))
			return
		} // не нашлась ветка
	} else {
		thread, err = SELECTThread_slug(slug)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find thread with slug: " + slug))
			return
		} // не нашлась ветка
		id = thread.Id
	}
	vote.Thread = id

	// Не нужно. INSERT ниже выдаст эту ошибку
	/* _, err = SELECTUser_nickname(vote.Nickname)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		response.Write(toMessage("Can't find user with id: " + vote.Nickname))
		return
	} // пользователя не нашлось */

	totalResult := vote.Result
	insertedVote, err := INSERTVote(vote)
	if err != nil {
		insertedVote, err = SELECTVote_nickname_thread(vote.Nickname, id)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Can't find user with id: " + vote.Nickname + " or thread with slug: " + slug))
			return
		}
		if insertedVote.Result != vote.Result {
			totalResult = vote.Result - insertedVote.Result
			insertedVote, err = UPDATEVote_nickname_thread(vote.Nickname, id, vote.Result)
			if err != nil {
				response.WriteHeader(http.StatusNotFound)
				response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
				return
			}
		} else {
			totalResult = 0
		}
	} // ошибка в запросе в БД

	// голос добавился
	if totalResult != 0 {
		thread, err = UPDATEThreadVotes_id(id, totalResult)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			response.Write(toMessage("Invalid DB request. Error: " + err.Error()))
			return
		} // ошибка в запросе в БД
	}
	body, err := json.Marshal(thread)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(toMessage("Can't marshal JSON file"))
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(body)
}
