package main

import (
	"encoding/hex"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	generator "github.com/ledex/passcard-generator/internal"
	"github.com/ledex/passcard-generator/model"
	"github.com/ledex/passcard-generator/views"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getRandomPassCard)
	mux.HandleFunc("GET /{pci}", getPassCardFromId)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func getRandomPassCard(w http.ResponseWriter, r *http.Request) {
	v := 1
	csa, _ := hex.DecodeString("01")
	seed, _ := uuid.NewRandom()

	pci := model.PassCardIdentifier{
		v,
		csa[0],
		seed,
	}

	http.Redirect(w, r, "http://"+r.Host+r.URL.String()+pci.String(), http.StatusContinue)
}

func getPassCardFromId(w http.ResponseWriter, r *http.Request) {
	pci, err := model.FromString(r.PathValue("pci"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	cols, err := strconv.Atoi(r.URL.Query().Get("cols"))
	if err != nil {
		cols = 28
	}

	rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
	if err != nil {
		rows = 20
	}

	pc, err := generator.GeneratePassCard(*pci, rows, cols)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	views.PasswordCardView(*pc).Render(r.Context(), w)
}
