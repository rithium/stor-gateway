package handler

import (
	"net/http"
	"log"
	"../model"
	"bytes"
	"io/ioutil"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
)

func HandleGetData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	data := model.DataRequest{}

	err := data.FromQuery(r.URL.Query())

	log.Printf("%+v", data)

	if err != nil {
		log.Println("data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := "http://localhost:11000/data/1"

	raw := data.ToDataRequestPacketQuery()

	raw.Set("nodeId", vars["nodeId"])
	req, _ := http.NewRequest("GET", url, nil)

	req.URL.RawQuery = raw.Encode()

	log.Printf("%+v", req.URL.Query)

	log.Println(req.URL.RawQuery)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("data service:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("data-get validate:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Data validated")

	b, err := ioutil.ReadAll(resp.Body)

	log.Println(string(b))

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func HandlePostData(hub *model.Hub, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	client := &http.Client{}

	url := "http://localhost:11000/data/1"

	data := model.Data{}

	data.FromJson(r.Body)

	nodeId, _ := strconv.Atoi(vars["nodeId"])

	raw, err := data.ToDataRequestBytes(nodeId)

	hub.Broadcast <- &model.HubMessage{Key: data.ApiKey, Message: raw}

	if err != nil {
		log.Println("data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(raw))

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Println("data service:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("data-get validate:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}



	log.Println("Data validated")
	w.WriteHeader(http.StatusOK)
}
