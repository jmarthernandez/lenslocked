package controllers

import (
	"fmt"
	"net/http"

	"github.com/lenslocked/views"
)

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

// New is used to render the form where user can create
// an account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create is used to process the submitted signup form
// will create a new account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	fmt.Fprintf(w, r.PostFormValue("email"))
	fmt.Fprintf(w, r.PostFormValue("password"))
}
