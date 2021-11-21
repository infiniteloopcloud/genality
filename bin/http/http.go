package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/infiniteloopcloud/genality"
)

type Communicator struct {
	g genality.Descriptor
	authentication HTTPAuthentication
}

func New(authentication HTTPAuthentication, database string) (Communicator, error) {
	g, err := genality.New(genality.Opts{ConnectionString: database})
	if err != nil {
		return Communicator{}, err
	}
	return Communicator{
		g:              g,
		authentication: authentication,
	}, nil
}

type HTTPAuthentication interface {
	Check(data string) error
}

func (c Communicator) Serve() error {
	r := chi.NewRouter()
	c.routes(r)
	return nil
}

func (c Communicator) routes(r *chi.Mux) {
	r.Post("/records", c.add)
}

func (c Communicator) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Response{Error: err}.ToJSON())
		return
	}

	if err := c.g.Add(ctx, ""); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Response{Error: err}.ToJSON())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(Response{Message: "success"}.ToJSON())
	return
}
