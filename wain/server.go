package wain

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	*http.Server
	S3 map[string]*S3Connection
}

func CreateHTTPServer(config *Config) (*Server, error) {

	s3, err := CreateS3(config)
	if err != nil {
		return nil, err
	}

	server := &Server{
		&http.Server{Addr: fmt.Sprintf(":%d", config.Port)},
		s3,
	}

	rtr := mux.NewRouter()
	for _, url := range config.Urls {
		rtr.HandleFunc(url.Pattern, server.handleProcessing(url)).Methods("GET")
	}

	server.Handler = rtr
	return server, nil
}

func (s *Server) handleProcessing(url ConfigUrl) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		options := ResizeOptions{
			Width:  StringToInt(params["width"]),
			Height: StringToInt(params["height"]),
			Format: params["format"],
			Params: params,
		}
		imageBytes, err := HandleProcessing(url, s.S3, options)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("%s", err)))
		} else {
			//TODO w.Header().Set("Content-Type", "")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageBytes)))
			w.WriteHeader(http.StatusOK)
			w.Write(imageBytes)
		}
	}
}
