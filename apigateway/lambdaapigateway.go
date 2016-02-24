package apigateway

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

/*
This file contains functions and types to create data structures from lambda
function input received from api-gateway. You need to use the associated template
to map values into the structure
*/

type apiRequest struct {
	m map[string]string
}

// DecodeRequest takes raw json supplied to the lambda function via node.js and returns
// a struct with a GetValue func to extract data
func DecodeRequest(cmd string) (*apiRequest, error) {

	ar := &apiRequest{}
	ar.m = make(map[string]string)

	var err error
	var event map[string]json.RawMessage

	err = json.Unmarshal([]byte(cmd), &event)
	if err != nil {
		return ar, fmt.Errorf("unable to find apiRequest in input: %v\n", err)
	}

	for k, v := range event {
		if isJSON(v) {
			copyFromJSON(ar, k, v)
		} else {
			ar.m[strings.ToLower(k)] = string(v)
		}
	}
	return ar, nil
}

// copyFromJSON recurcively copys elements from the source json string
func copyFromJSON(ar *apiRequest, pk string, j json.RawMessage) {
	var tm map[string]json.RawMessage

	_ = json.Unmarshal(j, &tm)
	for k, v := range tm {
		if isJSON(v) {
			copyFromJSON(ar, pk+"."+k, v)
		} else {
			ar.m[strings.ToLower(pk+"."+k)] = string(v)
		}
	}
}

// isJSON returns true if it is given a json string
func isJSON(s json.RawMessage) bool {
	var js map[string]interface{}
	return json.Unmarshal(s, &js) == nil
}

// GetValue returns the event attribute as a string or an empty string if not found
func (ar *apiRequest) GetValue(a string) string {
	s, _ := ar.GetValueBool(a)
	return s
}

// GetValueBool returns the event attribute as a string and true or an empty string and false
// if the attribute can not be found
func (ar *apiRequest) GetValueBool(a string) (string, bool) {

	// make sure the struct is valid before scanning it
	if ar == nil {
		return "", false
	}

	for k, v := range ar.m {
		if strings.EqualFold(a, k) {
			return strings.Trim(v, "\""), true
		}
	}
	return "", false
}

// ListAttributes prints out a list of all known attributes. Debugging
func (ar *apiRequest) ListAttributes() {

	for k, v := range ar.m {
		fmt.Printf("k: %s v: %s\n", k, v)
	}
}

// GetJSON will return a JSON encoded byte array of all values
func (ar *apiRequest) GetJSON() ([]byte, error) {
	return json.MarshalIndent(ar.m, "", "\t")
}

// GetKeys returns a slice of strings containing each available value name
func (ar *apiRequest) GetKeys() []string {

	r := make([]string, len(ar.m))
	for k, _ := range ar.m {
		r = append(r, k)
	}
	return r
}

type Response struct {
	Code    string `json:"errCode"`
	Message string `json:"errMessage"`
}

// ExitOnErr will create a json error required by api-gateway and the exit
// This will be detected by api-gateway and it will send the correct
// response back to the remote client
// code = apigateway regex code to map to response value
// msg = message to return to remote client
// lamerr = error to log to Lambda logs
func (ar *apiRequest) ExitOnErr(code, msg, lamerr string) {
	r := &Response{
		Code:    code,
		Message: msg,
	}
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Printf("error with json marshal: %v\n", err)
	} else {
		//
		os.Stdout.Write(b)
	}

	if lamerr != "" {
		// if provided then log info to lambda logs
		os.Stderr.Write([]byte(lamerr))
	}

	// exit > 1 to ensure wrapper passes to api-gateway as error and is mapped to a response
	os.Exit(2)
}

// Redirect302 will return a 302 redirect to the provided location and exit
func (ar *apiRequest) Redirect302(loc string) {
	os.Stdout.Write([]byte(loc))
	// exit 1 to ensure the wrapper passes to api-gateway as a redirect
	os.Exit(1)
}

/*

 */
