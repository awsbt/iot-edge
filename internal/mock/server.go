/*
 * Copyright 2020 ForgeRock AS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock

import (
	"encoding/json"
	"fmt"
	"github.com/ForgeRock/iot-edge/pkg/things/callback"
	"github.com/ForgeRock/iot-edge/pkg/things/payload"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CookieName         = "iPlanetDirectoryPro"
	SimpleTestRealm    = "testRealm"
	SimpleTestAuthTree = "testTree"
)

func SimpleClientAuthResponse(payload *payload.Authenticate, name string) {
	payload.Callbacks = []callback.Callback{
		{
			Type:   callback.TypeNameCallback,
			Output: []callback.Entry{{Value: "simple-thing"}},
			Input:  nil,
		},
	}
}

// Server mocks the endpoints of AM used by iot edge
type Server struct {
	ServerInfoHandler   http.HandlerFunc
	AuthenticateHandler http.HandlerFunc
}

// processAuthentication mocks a simple auth tree
// On each call, the server responds by asking for the Thing's name
// On the 2nd call and onwards, the server appends the Thing's name to the incoming auth id
// If the Thing's name does not match the embedded name, then the authentication fails
// On the 4th successful call, the authentication succeeds
func processAuthentication(authenticatePayload payload.Authenticate) (reply payload.Authenticate, err error) {
	stdCB := []callback.Callback{
		{Type: callback.TypeNameCallback, Input: []callback.Entry{{}}, Output: []callback.Entry{{}}},
	}
	if authenticatePayload.AuthId == "" {
		reply.AuthId = uniuri.New()
		reply.Callbacks = stdCB
		return reply, nil
	}
	name := ""
	for _, cb := range authenticatePayload.Callbacks {
		if cb.Type == callback.TypeNameCallback && len(cb.Input) > 0 && cb.Input[0].Value != "" {
			name = cb.Input[0].Value
			break
		}
	}
	if name == "" {
		return reply, fmt.Errorf("no name provided")
	}
	count := 0
	token := ""
	for count, token = range strings.Split(authenticatePayload.AuthId, ".") {
		if count == 0 {
			continue
		}
		if token != name {
			return reply, fmt.Errorf("malformed token %s\n", authenticatePayload.AuthId)
		}
	}
	if count < 2 {
		reply.AuthId += authenticatePayload.AuthId + "." + name
		reply.Callbacks = stdCB
		return reply, nil
	}
	return payload.Authenticate{TokenId: "12345"}, nil
}

// NewSimpleServer creates a test server that does the minimum to serve the iot endpoints
// see processAuthentication for the authentication workflow
func NewSimpleServer() Server {
	return Server{
		ServerInfoHandler: func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte(fmt.Sprintf(`{"cookieName":"%s"}`, CookieName)))
		},
		AuthenticateHandler: func(writer http.ResponseWriter, request *http.Request) {
			// check that the query is correct
			if realm, ok := request.URL.Query()["realm"]; !ok || len(realm) != 1 || realm[0] != SimpleTestRealm {
				http.Error(writer, "incorrect realm query", http.StatusBadRequest)
			}
			if tree, ok := request.URL.Query()["authIndexValue"]; !ok || len(tree) != 1 || tree[0] != SimpleTestAuthTree {
				http.Error(writer, "incorrect auth tree query", http.StatusBadRequest)
			}
			if authType, ok := request.URL.Query()["authIndexType"]; !ok || len(authType) != 1 || authType[0] != "service" {
				http.Error(writer, "incorrect auth type query", http.StatusBadRequest)
			}
			// expect that the username has been provided
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				http.Error(writer, "unable to read request body", http.StatusBadRequest)
			}
			var payload payload.Authenticate
			if err := json.Unmarshal(body, &payload); err != nil {
				http.Error(writer, "unable to decode request body", http.StatusBadRequest)
			}
			reply, err := processAuthentication(payload)
			if err != nil {
				fmt.Println(err)
				http.Error(writer, err.Error(), http.StatusBadRequest)
			}
			replyBytes, err := json.Marshal(reply)
			if err != nil {
				http.Error(writer, "unable to marshall response", http.StatusInternalServerError)
			}
			writer.Write(replyBytes)
		},
	}
}

// Start the test server
func (s Server) Start(addr string) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/json/serverinfo/*", s.ServerInfoHandler)
	router.HandleFunc("/json/authenticate", s.AuthenticateHandler)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		server.ListenAndServe()
	}()
	return server
}
