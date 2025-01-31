package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"

	generator "github.com/ledex/passcard-generator/internal"
	"github.com/ledex/passcard-generator/model"
	"github.com/ledex/passcard-generator/views"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", getRandomPassCard)
	mux.HandleFunc("GET /{pci}", getPassCardFromId)
	mux.HandleFunc("GET /multi", getMultipleRandomPassCards)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func getMultipleRandomPassCards(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		count = 10
	}

	cols, err := strconv.Atoi(r.URL.Query().Get("cols"))
	if err != nil {
		cols = 28
	}

	rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
	if err != nil {
		rows = 20
	}

	v, err := strconv.Atoi(r.URL.Query().Get("version"))
	if err != nil {
		v = 1
	}

	css := r.URL.Query().Get("charset")
	if css == "" {
		css = "01"
	}

	csa, _ := hex.DecodeString(css)

	pcs := make([]*model.PassCard, 0)
	for range count {
		pci, err := model.WithRandomSeed(v, csa[0])
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		pc, err := generator.GeneratePassCard(*pci, rows, cols)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		pcs = append(pcs, pc)
	}

	views.MultiView(pcs).Render(r.Context(), w)
}

func getRandomPassCard(w http.ResponseWriter, r *http.Request) {
	v := 1
	csa, _ := hex.DecodeString("01")

	pci, err := model.WithRandomSeed(v, csa[0])
	if err != nil {
		http.Error(w, err.Error(), 500)
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
