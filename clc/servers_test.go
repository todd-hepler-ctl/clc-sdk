package clc_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetServer(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	resource := getServerResource(assert, name)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.ServerService{client(ms.URL)}
	server, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, server.Name)
}

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	r := postServerResponse(assert)
	ms := mockServer(r)
	defer ms.Close()

	service := clc.ServerService{client(ms.URL)}
	server := clc.Server{
		Name:           "va1testserver01",
		CPU:            1,
		MemoryGB:       1,
		GroupID:        "group",
		SourceServerID: "UBUNTU",
		Type:           "standard",
	}
	_, err := service.Create(server)

	assert.Nil(err)
}

func getServerResource(assert *assert.Assertions, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path != "/servers/test/"+name {
			assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
		}

		server := &clc.ServerResponse{Name: name}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(server)
	}
}

func postServerResponse(assert *assert.Assertions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			assert.Fail("POST server method should be POST", r.Method)
		}

		if r.URL.Path != "/servers/test" {
			assert.Fail("POST server hitting wrong endpoint", r.URL.Path)
		}

		server := &clc.Server{}
		err := json.NewDecoder(r.Body).Decode(server)
		if err != nil {
			assert.Fail("Failed to serialize server", err)
		}

		if !server.Valid() {
			assert.Fail("Server missing required fields", server)
		}

		create := &clc.ServerCreateResponse{
			Server:   server.Name,
			IsQueued: true,
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(create)
	}
}