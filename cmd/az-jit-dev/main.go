package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	tea "charm.land/bubbletea/v2"
	"github.com/h2non/gock"
	"github.com/varantha/az-jit/internal/clients"
	"github.com/varantha/az-jit/internal/ui"
)

//go:embed testdata/*.json
var testdataFS embed.FS

func main() {

	defer gock.Off()

	route := "/providers/Microsoft.Authorization/roleEligibilityScheduleInstances"
	response, err := testdataFS.ReadFile("testdata/eligibleschedule.json")
	if err != nil {
		log.Fatalf("failed to read test data file: %v", err)
	}

	gock.New("https://management.azure.com").
		Get(route).
		Reply(200).
		BodyString(string(response))

	route = "/providers/Microsoft.Authorization/roleAssignmentScheduleInstances"
	response, err = testdataFS.ReadFile("testdata/activeassignments.json")
	if err != nil {
		log.Fatalf("failed to read test data file: %v", err)
	}

	gock.New("https://management.azure.com").
		Get(route).
		Reply(200).
		BodyString(string(response))

	hc := &http.Client{}
	gock.InterceptClient(hc)

	armOptions := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: hc,
		},
	}

	OCOptions := clients.OmniClientOptions{
		ARM: armOptions,
	}

	cred := &clients.CredentialMock{}
	client, err := clients.NewOmniClient(cred, OCOptions)
	if err != nil {
		log.Fatalf("failed to build omni client: %v", err)
	}

	p := tea.NewProgram(ui.InitialModel(client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
