package playground

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dhl1402/covidscript/cmd/api/svc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func NewHandler(s Service, router *httprouter.Router, logger *zap.SugaredLogger) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(svc.EncodeError),
	}
	router.Handler(http.MethodPost, "/api/:ver/interpret", kithttp.NewServer(
		makeInterpretEndpoint(s, svc.NewLoggerMiddleware(logger, "Interpret")),
		decodeInterpretRequest,
		encodeInterpretResponse,
		options...,
	))
	return router
}

func decodeInterpretRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := InterpretRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeInterpretResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(r)
	return nil
}
