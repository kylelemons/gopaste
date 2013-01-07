package server

import (
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

const maxAge = 12 * time.Hour

func expire(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	max := maxAge
	if age := r.FormValue("age"); age != "" {
		if ageDur, err := time.ParseDuration(age); err == nil {
			max = ageDur
		}
	}
	oldest := time.Now().Add(-max)

	q := datastore.NewQuery("paste").
		Filter("Pasted<=", oldest).
		KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		ctx.Errorf("get: %s", err)
		http.Error(w, "get: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.Infof("Deleting %d entitites over %v old", len(keys), max)
	if err := datastore.DeleteMulti(ctx, keys); err != nil {
		ctx.Errorf("delete: %s", err)
		http.Error(w, "delete: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")
}

func init() {
	http.HandleFunc("/admin/expire", expire)
}
