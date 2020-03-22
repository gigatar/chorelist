package controllers

import (
	"chorelist/chore-service/daos"
	"net/http"
)

// ChoreController data type.
type ChoreController struct {
	dao daos.ChoreDAO
}

func (c *ChoreController) ListChores(w http.ResponseWriter, r *http.Request) {

}
