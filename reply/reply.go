package reply

import (
	"encoding/json"
	"net/http"

	"github.com/uber-go/zap"
)

var logger = zap.New(zap.NewJSONEncoder())

func Json(w http.ResponseWriter, r *http.Request, status int, data interface{}) {

	d, err := json.Marshal(data)
	if err != nil {
		logger.Error("Marshal()",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("err", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(d); err != nil {
		logger.Error("Write()",
			zap.String("err", err.Error()),
		)
	}
}
