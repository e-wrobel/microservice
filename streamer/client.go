package streamer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/e-wrobel/microservice/storage"
)

type Client struct {
	fileHandler        *os.File
	createJSONEndpoint string
	channel            chan []byte
}

// New is constructor for the *Client
func New(fileName, createJSONEndpoint string) *Client {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error to read file: %v, %v", fileName, err.Error())
	}

	return &Client{
		fileHandler:        f,
		createJSONEndpoint: createJSONEndpoint,
		channel:            make(chan []byte),
	}
}

// SendJSONFile is intended to parse available ports.json file and send in to
// the API endpoint
func (s *Client) SendJSONFile() error {
	defer s.fileHandler.Close()
	go sendData(s.channel, s.createJSONEndpoint)

	scanner := bufio.NewScanner(s.fileHandler)
	var indexName string
	var jsonDetails []byte
	var indexFound bool

	// Let's read desired file, line by line
	for scanner.Scan() {
		line := scanner.Bytes()
		// Check if it's main, opening '{'
		if len(line) > 0 && line[0] == 123 {
			continue
		}

		// Check if it is closing '  }' for entity
		if len(line) > 2 && line[2] == 125 {
			err := prepareData(indexName, jsonDetails, s.channel)
			switch err {
			case ErrUnmarshalling:
				log.Printf("We were unable to unmarshal data, check input file for index: %v", indexName)
			case ErrMarshaling:
				log.Printf("We were unable to marshal file, current body: %v", jsonDetails)
			case nil:
				log.Printf("Data has been prepared")
			default:
				log.Fatalf("Unable to prepare data: %v", err)
			}

			// Reset objects representing json entity
			indexName = ""
			indexFound = false
			jsonDetails = []byte{}
			continue
		}

		// Identifier in PortEntity
		stringLine := string(line)
		if (stringLine == strings.ToUpper(stringLine)) && !indexFound {
			c, _, _ := strings.Cut(stringLine, ":")
			res := strings.Trim(c, "\" ")
			if isIndex(res) {
				indexName = res
				indexFound = true
				continue
			}
		}
		jsonDetails = append(
			jsonDetails,
			line...,
		)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// prepareData is going to marshal json entity and send to the channel
func prepareData(indexName string, body []byte, c chan []byte) error {
	var fullBody []byte

	fullBody = append(fullBody, '{')
	fullBody = append(fullBody, body...)
	fullBody = append(fullBody, '}')

	jsonDetails := &storage.PortDetails{}
	if err := json.Unmarshal(fullBody, jsonDetails); err != nil {
		return ErrUnmarshalling
	}

	data := &storage.PortEntity{
		Identifier:  indexName,
		PortDetails: jsonDetails,
	}

	res, err := json.Marshal(data)
	if err != nil {
		return ErrMarshaling
	}
	c <- res

	return nil
}

// sendData is intended to send received data via channel and send it to API
func sendData(c chan []byte, endpoint string) {
	for {
		data := <-c
		reader := bytes.NewReader(data)
		_, err := sendRequest(endpoint, http.MethodPost, reader)
		if err != nil {
			log.Printf("Unable to send data to api server: %v", err)
		}
	}
}

// isIndex is helper function to check is given line from json file represents index
func isIndex(lineFromFile string) bool {
	var found bool
	for _, r := range lineFromFile {
		if r >= 'A' || r <= 'Z' {
			found = true
			break
		}
	}

	return found
}
