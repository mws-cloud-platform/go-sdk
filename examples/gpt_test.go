package examples_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mws.cloud/go-sdk/mws"
	"io"
	"log"
	"net/http"
	"os"

	mwshttp "go.mws.cloud/go-sdk/pkg/http"
	gptclient "go.mws.cloud/go-sdk/service/gpt/client"
	gptmodel "go.mws.cloud/go-sdk/service/gpt/model"
	gptsdk "go.mws.cloud/go-sdk/service/gpt/sdk"
	gptref "go.mws.cloud/go-sdk/service/resources/references/gpt"
)

// This example demonstrates creating, infering and deleting a GPT deployment
// "example-deployment" inside the `$MWS_PROJECT` project.
//
// Before running this example, set a service account API key environment
// variable:
//
//	export MWS_SERVICE_ACCOUNT_API_KEY="your-service-account-api-key"
func Example_gptDeployment() {
	ctx := context.Background()

	apiKey := os.Getenv("MWS_SERVICE_ACCOUNT_API_KEY")
	if apiKey == "" {
		log.Panicln("MWS_SERVICE_ACCOUNT_API_KEY is not set")
	}

	// Use the default SDK loader. It will load configuration from the
	// environment variables and sensible defaults. You can override logic using
	// [mws.LoadSDKOption] options. Check the [mws.Load] and [mws.Config] for
	// more details.
	sdk, err := mws.Load(ctx)
	if err != nil {
		log.Panicln("load sdk:", err)
	}
	defer sdk.Close(ctx)

	// Create a new gpt model client using the provided SDK.
	modelClient, err := gptsdk.NewModel(ctx, sdk)
	if err != nil {
		log.Panicln("create gpt model client:", err)
	}

	// Get first model from list that takes and returns text.
	gptModels, err := modelClient.ListModels(ctx, gptclient.ListModelsRequest{
		Filter: new("spec.inputModalities.text = true AND spec.outputModalities.text = true"),
	})
	if err != nil {
		log.Panicln("list models:", err)
	}

	items := gptModels.GetItems()
	if len(items) == 0 {
		log.Panicln("no models found")
	}
	modelRef := items[0].GetMetadata().GetId().AsRef()

	// Create a new deployments client using the provided SDK.
	deploymentClient, err := gptsdk.NewDeployment(ctx, sdk)
	if err != nil {
		log.Panicln("create gpt deployment client:", err)
	}

	// Use example names for demonstration purposes.
	const deploymentName = "example-deployment"

	// Create a new deployment required for Chat Completions API. Clean them up
	// after the example run.
	deleteDeployment := createDeployment(ctx, deploymentClient, deploymentName, modelRef)
	defer deleteDeployment()

	chatCompletions(ctx, sdk, deploymentName, apiKey)
}

func chatCompletions(ctx context.Context, sdk *mws.SDK, modelName string, apiKey string) {
	endpoint, err := sdk.ServiceEndpointResolver().Resolve(ctx, "gpt")
	if err != nil {
		log.Panicln("resolve gpt service endpoint:", err)
	}

	project := sdk.DefaultProject()
	url := fmt.Sprintf("%s/projects/%s/openai/v1/chat/completions", endpoint, project)

	reqBody := map[string]any{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "user", "content": "2+2*2=?"},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Panicln("marshal request:", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		log.Panicln("create request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := mwshttp.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("send request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("read response:", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Panicf("bad status: %s, body: %s\n", resp.Status, string(body))
	}

	// Print the model response
	fmt.Println("response:", string(body))
}

func createDeployment(ctx context.Context, deploymentClient *gptsdk.Deployment, deploymentName string, modelRef *gptref.ModelRef) func() {
	deployment, err := deploymentClient.CreateDeployment(ctx, gptclient.UpsertDeploymentRequest{
		DeploymentName: deploymentName,
		Body: gptmodel.DeploymentRequest{
			Spec: gptmodel.DeploymentSpecRequest{
				IsActive: new(true),
				Model:    modelRef,
			},
		},
	}, gptclient.WithWait())
	if err != nil {
		log.Panicln("create deployment:", err)
	}
	fmt.Println("deployment created:", deployment.GetMetadata().GetId().ResourceName())
	return func() {
		err := deploymentClient.DeleteDeployment(ctx, gptclient.DeleteDeploymentRequest{
			DeploymentName: deploymentName,
		}, gptclient.WithWait())
		if err != nil {
			log.Panicln("delete deployment:", err)
		}
		fmt.Println("deployment deleted:", deployment.GetMetadata().GetId().ResourceName())
	}
}
