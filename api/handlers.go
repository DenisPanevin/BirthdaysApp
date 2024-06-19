package main

import (
	"birthdays/models"
	"encoding/json"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"time"
)

var (
	incorrectEmailOrPassword = errors.New("Incorrect Email or Password")
	notAuth                  = errors.New("not authenticated")
)

const (
	SessionName               = "BrthdaysSesh"
	contextKeyUser contextKey = iota
)

type contextKey int8

func (a *Application) customErr(w http.ResponseWriter, r *http.Request, code int, err error) {
	a.costumRespond(w, r, code, map[string]string{"error:": err.Error()})
}
func (a *Application) costumRespond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (a *Application) sayHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hey btchs")
	}
}

func (a *Application) showUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]
		println(userID)
		i, err := strconv.Atoi(userID)

		U := &models.User{}
		U, err = a.storage.FindById(i)
		if err != nil {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}

		U.ClearFields()
		a.costumRespond(w, r, http.StatusOK, U)
	}
}
func (a *Application) InsertUser() http.HandlerFunc {
	type input struct {
		Email    string
		Password string
		Name     string
		Date     string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		in := &input{}
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			a.customErr(w, r, http.StatusBadRequest, err)
			return
		}
		if err := validation.Validate(&in.Date, validation.Required, validation.Date("2006-01-02")); err != nil {
			a.customErr(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		date, _ := time.Parse("2006-01-02", in.Date)
		u := models.User{
			Email:    in.Email,
			Password: in.Password,
			Name:     in.Name,
			Date:     date,
		}
		if err := u.ValidateUser(); err != nil {
			a.customErr(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := a.storage.CreateUser(&u); err != nil {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}
		a.costumRespond(w, r, http.StatusCreated, nil)
		//http.Redirect(w, r, fmt.Sprintf("/users/&id=%s", u.Id), http.StatusSeeOther)

	}
}
func (a *Application) CreateSession() http.HandlerFunc {
	type input struct {
		Email    string
		Password string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		in := &input{}
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			a.customErr(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := a.storage.FindByEmail(in.Email)
		if err != nil || !u.ComparePassword(in.Password) {

			a.customErr(w, r, http.StatusBadRequest, incorrectEmailOrPassword)
			return
		}
		session, err := a.sessionStore.Get(r, SessionName)
		if err != nil {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Values["UserId"] = u.Id
		if err := session.Save(r, w); err != nil {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}
		a.costumRespond(w, r, http.StatusOK, nil)

	}
}
func (a *Application) LogOut() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sessionStore.Get(r, SessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete(session.Values, "UserId")

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.costumRespond(w, r, http.StatusOK, "LogOut")

	}
}
func (a *Application) WhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.costumRespond(w, r, http.StatusOK, r.Context().Value(contextKeyUser))
	}
}

func (a *Application) subscribe() http.HandlerFunc {
	type input struct {
		Email string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		in := &input{}
		if err := json.NewDecoder(r.Body).Decode(in); err != nil {
			a.customErr(w, r, http.StatusBadRequest, err)
			return
		}
		err := validation.ValidateStruct(
			in,
			validation.Field(&in.Email, validation.Required, is.Email),
		)
		if err != nil {
			a.customErr(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		u := r.Context().Value(contextKeyUser).(*models.User)

		if err := a.storage.SubscribeTo(u, in.Email); err != nil {
			a.customErr(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		a.costumRespond(w, r, http.StatusOK, nil)

	}
}
func (a *Application) getUserSubscriptions() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		u := r.Context().Value(contextKeyUser).(*models.User)

		err, ids := a.storage.GetUserSubscriptions(u)
		if err != nil {
			a.customErr(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		a.costumRespond(w, r, http.StatusOK, ids)

	}
}
