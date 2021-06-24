package transport

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterHttpService(app app.App) {

	http.HandleFunc("/pipeline", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var request model.PipelineMessage
		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Fatal(err)
		}

		app.ServiceStorage.Pipes.GetJobInfo(request, 1)
	})

	log.Fatal(http.ListenAndServe(":1212", nil))
}
