package gui

import (
	"bytes"
	"log"
	"net/http"

	"github.com/ubavic/bas-celik/document"
)

func uploadDoc(doc document.Document) {
	data, err := doc.BuildJson()
	if err != nil {
		log.Fatalf("generating json: %v", err)
		return
	}

	posturl := application.Preferences().String("apiURL")
	if len(posturl) < 1 {
		setApiStatus("API: Nije Aktiviran")
		return
	}

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(data))
	if err != nil {
		setApiStatus("API: Greska Prilikom kreiranja dokumenta")
		return
	}
	setApiStatus("API:Slanje Podataka")

	r.Header.Add("Content-Type", "application/json")

	apiKey := application.Preferences().String("apiKey")
	if len(apiKey) > 0 {
		r.Header.Add("Authorization", apiKey)
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		setApiStatus("API: Greska Prilikom slanja")
		return
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusUnauthorized: //401
		setApiStatus("API: Nije Autorizovan")
	case http.StatusBadRequest: //400
		setApiStatus("API: Neispravan Dokument")
	case http.StatusCreated: //201
		setApiStatus("API: Dodata LK")
	case http.StatusOK: //200
		setApiStatus("API: Auzirana LK")
	case http.StatusConflict: //409
		setApiStatus("API: LK Vec postoji u bazi")
	default:
		setApiStatus("API: Greska")
	}
}
