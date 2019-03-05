package controller

import (
	"blog/models"
	"blog/utils"
	"net/http"
	"strconv"
)

// GET /err?msg=
func Err(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	sess, err := utils.Session(w, r)
	data := struct {
		User models.User
		Msg  string
	}{Msg: vals.Get("msg")}
	if err != nil {
		utils.GenerateHTML(w, data, "layout", "public.navbar", "error")
	} else {
		data.User, _ = sess.User()
		utils.GenerateHTML(w, data, "layout", "private.navbar", "error")
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")
	data := struct {
		Entries  []models.Entry
		User     models.User
		Page     int
		MaxPage  int
		HasNext  bool
		HasFont  bool
		NextPage int
		FontPage int
	}{}
	num, _ := models.EntriesCount()
	if num != 0 {
		if num%models.PageCount == 0 {
			data.MaxPage = num / models.PageCount
		} else {
			data.MaxPage = num/models.PageCount + 1
		}
	} else {
		data.MaxPage = 1
		data.Page = 1
	}

	if page != "" {
		var err error
		data.Page, err = strconv.Atoi(page)
		if err != nil {
			utils.Error_message(w, r, "Cannot get entries")
			return
		}
		if data.Page > data.MaxPage {
			data.Page = data.MaxPage
		}

		if data.Page < 1 {
			data.Page = 1
			data.HasFont = false
		} else if data.Page == 1 {
			data.HasFont = false
		} else {
			data.HasFont = true
			data.FontPage = data.Page - 1
		}
		entries, err := models.EntriesPage(data.Page)
		data.Entries = entries

	} else {
		entries, err := models.Entries()
		data.Entries = entries
		data.Page = 1
		data.HasFont = false
		if err != nil {
			utils.Error_message(w, r, "Cannot get entries")
			return
		}
	}
	if data.Page >= data.MaxPage {
		data.HasNext = false
	} else {
		data.HasNext = true
		data.NextPage = data.Page + 1
	}
	sess, err := utils.Session(w, r)
	if err != nil {
		utils.GenerateHTML(w, data, "layout", "public.navbar", "index")
	} else {
		user, err := sess.User()
		data.User = user
		if err != nil {
			utils.GenerateHTML(w, data, "layout", "public.navbar", "index")
		} else {
			utils.GenerateHTML(w, data, "layout", "private.navbar", "index")
		}
	}
}

// 友情链接
func FriendLink(w http.ResponseWriter, r *http.Request) {
	sess, err := utils.Session(w, r)
	data := struct {
		User models.User
	}{}
	if err != nil {
		utils.GenerateHTML(w, data, "layout", "public.navbar", "friend")
	} else {
		data.User, _ = sess.User()
		utils.GenerateHTML(w, data, "layout", "private.navbar", "friend")
	}
}

// About Me
func AboutMe(w http.ResponseWriter, r *http.Request) {
	sess, err := utils.Session(w, r)
	data := struct {
		User models.User
	}{}
	if err != nil {
		utils.GenerateHTML(w, data, "layout", "public.navbar", "aboutme")
	} else {
		data.User, _ = sess.User()
		utils.GenerateHTML(w, data, "layout", "private.navbar", "aboutme")
	}
}
