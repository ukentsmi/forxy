package http

import (
	"bytes"
	"encoding/json"
	ForxyHttpApiRequest "github.com/dragoscojocaru/forxy/pkg/handler/http/api/request"
	"github.com/dragoscojocaru/forxy/pkg/handler/http/api/response"
	"github.com/dragoscojocaru/forxy/pkg/logger"
	"net/http"
)

// HTTPSequentialHandler used for testing purposes
func HTTPSequentialHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var body ForxyHttpApiRequest.ForxyBodyPayload
	err := decoder.Decode(&body)
	if err != nil {
		go logger.FileErrorLog(err)
	}

	forxyResponsePayload := response.NewForxyResponsePayload()
	for idx := range body.Requests {

		bodyReader := bytes.NewReader(body.Requests[idx].Body)

		req, err1 := http.NewRequest(body.Requests[idx].Method, body.Requests[idx].URL, bodyReader)
		if err1 != nil {
			logger.FileErrorLog(err1)
		}

		for key, value := range body.Requests[idx].Headers {
			req.Header.Set(key, value)
		}

		host, err := GetHost(body.Requests[idx].URL)
		if err != nil {
			logger.FileErrorLog(err)
		}

		client := connectionPool.GetServerConnection(host)
		resp, err2 := client.Do(req)

		if err2 != nil {
			logger.FileErrorLog(err2)
		}

		forxyResponsePayload.AddResponse(idx, *resp)

	}
	forxyPayloadWriter := response.NewForxyPayloadWriter()
	forxyPayloadWriter.JsonMarshal(w, *forxyResponsePayload)
}
