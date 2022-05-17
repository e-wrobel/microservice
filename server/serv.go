package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/e-wrobel/microservice/storage"
)

type Server struct {
	Storage    *storage.Storage
	HTTPServer *http.Server
}

func New(dbFile, listen string) (*Server, error) {
	db, err := storage.New(dbFile)
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate storage in server: %v", err)
	}

	httpServer := &http.Server{
		Addr: listen,
	}

	return &Server{
		Storage:    db,
		HTTPServer: httpServer,
	}, nil
}

func (s *Server) Start() {
	http.HandleFunc("/", s.createHandler())
	log.Printf("Listening on %v", s.HTTPServer.Addr)
	log.Fatal(s.HTTPServer.ListenAndServe())
}

func (s *Server) Shutdown(ctx context.Context) error {
	defer s.Storage.DB.Close()
	err := s.HTTPServer.Shutdown(ctx)

	return err
}

func (s *Server) createHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		obj := &storage.PortEntity{}
		buf := &bytes.Buffer{}
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			msg := fmt.Sprintf("Unable to read body: %v", err)
			json.NewEncoder(w).Encode(msg)
			return
		}
		err = json.Unmarshal(buf.Bytes(), obj)
		if err != nil {
			msg := fmt.Sprintf("Unable to get json: %v", err)
			json.NewEncoder(w).Encode(msg)
			return
		}

		// Finally, add item to the database
		identifier, details, err := portEntityToBytes(obj)
		if err != nil {
			msg := fmt.Sprintf("Unable to prepare data to storage: %v", err)
			json.NewEncoder(w).Encode(msg)
			return
		}

		err = s.Storage.Create(identifier, details)
		if err != nil {
			msg := fmt.Sprintf("Unable to add json to storage: %v", err)
			json.NewEncoder(w).Encode(msg)
			return
		}

		json.NewEncoder(w).Encode(obj)
	}
}

// portEntityToBytes is helper function to convert *PortEntity for database
func portEntityToBytes(obj *storage.PortEntity) (string, []byte, error) {
	if obj.Identifier == "" {
		return "", nil, ErrEmptyIdentifier
	}
	details, err := json.Marshal(obj.PortDetails)
	if err != nil {
		return "", nil, fmt.Errorf("unable to convert *PortDetails: %v", err)
	}

	return obj.Identifier, details, nil
}
