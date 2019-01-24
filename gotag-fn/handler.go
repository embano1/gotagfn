package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/embano1/faastagger"
	"github.com/openfaas-incubator/go-function-sdk"
)

var (
	tagger     *faastagger.Client
	err        error
	ctx        = context.Background()
	vCenterURL string
	vcUser     string
	vcPass     string
	tagID      string
	insecure   bool
)

func init() {
	// open reusable connection to vCenter
	// not checking env variables here as faastagger.New would throw error when connecting to VC
	vCenterURL = os.Getenv("VCURL")
	vcUser = os.Getenv("VCUSER")
	vcPass = os.Getenv("VCPASS")
	tagID = os.Getenv("TAGURN")
	if os.Getenv("INSECURE") == "true" {
		insecure = true
	}
	tagger, err = faastagger.New(ctx, nil, vCenterURL, vcUser, vcPass, insecure)
	if err != nil {
		log.Fatalf("could not get tags: %v", err)
	}
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	// verify request body
	var event faastagger.InbountEvent
	err = json.Unmarshal(req.Body, &event)
	if err != nil {
		return handler.Response{
			Body:       nil,
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// did we get a ManagedObjectReference?
	if event.MoRef == nil {
		return handler.Response{
			Body:       nil,
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("managedobjectreference must not be nil")
	}
	ref := event.MoRef
	err = tagger.TagVM(ctx, ref, tagID)
	if err != nil {
		return handler.Response{
			Body:       nil,
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	log.Printf("successfully tagged VM %v with tag %s", event.MoRef, tagID)
	err = tagger.Close(ctx)
	if err != nil {
		log.Printf("could not close connection: %v", err)
	}

	return handler.Response{
		Body:       []byte(fmt.Sprintf("successfully tagged VM %v with tag %s", event.MoRef, tagID)),
		StatusCode: http.StatusOK,
	}, err

}
