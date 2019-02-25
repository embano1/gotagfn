package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	sigCh      = make(chan os.Signal)
)

func init() {
	// open reusable connection to vCenter
	// not checking env variables here as faastagger.New would throw error when connecting to VC
	vCenterURL = os.Getenv("VC")
	vcUser = os.Getenv("VC_USER")
	vcPass = os.Getenv("VC_PASS")
	tagID = os.Getenv("TAG_URN")

	if os.Getenv("INSECURE") == "true" {
		insecure = true
	}

	tagger, err = faastagger.New(ctx, nil, vCenterURL, vcUser, vcPass, insecure)
	if err != nil {
		log.Fatalf("could not get tags: %v", err)
	}

	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	go handleSignal(ctx, sigCh)
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
	// TODO: change to ManagedObjectReference
	if event.ManagedObjectReference == nil {
		return handler.Response{
			Body:       nil,
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("managedobjectreference must not be nil")
	}
	ref := event.ManagedObjectReference
	err = tagger.TagVM(ctx, ref, tagID)
	if err != nil {
		return handler.Response{
			Body:       nil,
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	log.Printf("successfully tagged VM %v with tag %s", event.ManagedObjectReference, tagID)

	return handler.Response{
		Body:       []byte(fmt.Sprintf("successfully tagged VM %v with tag %s", event.ManagedObjectReference, tagID)),
		StatusCode: http.StatusOK,
	}, err

}

func handleSignal(ctx context.Context, sigCh <-chan os.Signal) {
	s := <-sigCh
	log.Printf("got signal: %v", s)
	err := tagger.Close(ctx)
	if err != nil {
		log.Printf("could not close connection to vCenter: %v", err)
	}
}
