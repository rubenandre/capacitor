package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/gimlet-io/capacitor/pkg/flux"
	"github.com/gimlet-io/capacitor/pkg/streaming"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var gitRepositoryResource = schema.GroupVersionResource{
	Group:    "source.toolkit.fluxcd.io",
	Version:  "v1",
	Resource: "gitrepositories",
}

func GitRepositoryController(
	dynamicClient *dynamic.DynamicClient,
	clientHub *streaming.ClientHub,
) *Controller {
	return NewDynamicController(
		"gitrepositories.source.toolkit.fluxcd.io",
		dynamicClient,
		gitRepositoryResource,
		func(informerEvent Event, objectMeta metav1.ObjectMeta, obj interface{}) error {
			switch informerEvent.EventType {
			case "create":
				fallthrough
			case "update":
				fallthrough
			case "delete":
				fmt.Printf("Changes in %s\n", objectMeta.Name)
				fluxState, err := flux.GetFluxState(dynamicClient)
				if err != nil {
					panic(err.Error())
				}
				fluxStateBytes, err := json.Marshal(fluxState)
				if err != nil {
					panic(err.Error())
				}
				clientHub.Broadcast <- fluxStateBytes
			}
			return nil
		})
}
