package model

import (
	"time"
	"io"
	"encoding/json"
	"net/url"
	"log"
)

type Data struct {
	ApiKey 		string
	Timestamp 	time.Time
	Data		map[string]interface{}
}

type DataPacket struct {
	NodeId		int
	Timestamp	time.Time
	TTL		int
	Data		map[string]interface{}
}

type DataRequest struct {
	ApiKey	string
	Start	time.Time
	End	time.Time
}

type DataRequestPacket struct {
	NodeId	int
	Start	time.Time
	End	time.Time
}

func (d *Data) FromJson(rawJson io.Reader) (error) {
	decoder := json.NewDecoder(rawJson)

	if err := decoder.Decode(d); err != nil {
		return err
	}

	return nil
}

func (d *DataRequest) FromQuery(query url.Values) (error) {
	apiKey := query.Get("apiKey")

	start, err := time.Parse(time.RFC3339, query.Get("start"))

	if err != nil {
		return err
	}

	end, err := time.Parse(time.RFC3339, query.Get("end"))

	if err != nil {
		return err
	}

	d.ApiKey 	= apiKey
	d.Start 	= start
	d.End 		= end

	return nil
}

func (d *DataRequest) ToDataRequestPacketQuery() (url.Values) {
	request := &DataRequestPacket{}

	request.NodeId = 1
	request.Start = d.Start
	request.End = d.End

	log.Printf("%+v", request)

	qs := url.Values{}

	qs.Set("nodeId", "1")
	qs.Set("start", request.Start.Format(time.RFC3339))
	qs.Set("end", request.End.Format(time.RFC3339))

	//log.Println(qs.Encode())

	return qs
}

func (d *Data) ToDataRequestBytes(nodeId int) ([]byte, error) {
	dataRequest := &DataPacket{}

	dataRequest.NodeId = nodeId
	dataRequest.Timestamp = d.Timestamp
	dataRequest.TTL = 3600
	dataRequest.Data = d.Data

	rawJson, err := json.Marshal(dataRequest)

	return rawJson, err
}