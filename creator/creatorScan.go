package nessusCreator

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kkirsche/nessusControl/api"
	"net/http"
	"strings"
	"sync"
)

func (c *Creator) createScan(createScanCh chan nessusAPI.CreateScan, skipTLSVerify bool) chan nessusAPI.CreateScanResponse {
	c.debugln("createScan(): Creating scans from JSON")
	createScanResponseCh := make(chan nessusAPI.CreateScanResponse)
	wg := new(sync.WaitGroup)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
	}
	httpClient := &http.Client{Transport: transport}

	for createScanJSON := range createScanCh {
		wg.Add(1)
		go func(createScanJSON nessusAPI.CreateScan, wg *sync.WaitGroup, httpClient *http.Client) {
			marshalledJSON, err := json.Marshal(createScanJSON)
			if err != nil {
				wg.Done()
				return
			}

			createdScan, err := c.apiClient.CreateScan(httpClient, string(marshalledJSON))
			if err == nil {
				createScanResponseCh <- createdScan
			}
			wg.Done()
		}(createScanJSON, wg, httpClient)
	}

	go func(wg *sync.WaitGroup, createScanResponseCh chan nessusAPI.CreateScanResponse) {
		wg.Wait()
		close(createScanResponseCh)
	}(wg, createScanResponseCh)

	return createScanResponseCh
}

// buildCreateScanJSON is used to create a struct which can be marshalled into
// JSON and sent to a remote Nessus server.
func (c *Creator) buildCreateScanJSON(requestScanch chan RequestedScan) chan nessusAPI.CreateScan {
	c.debugln("buildCreateScanJSON(): Building JSON for Create Scan")
	createScanJSONch := make(chan nessusAPI.CreateScan)
	wg := new(sync.WaitGroup)

	for requestedScan := range requestScanch {
		wg.Add(1)
		c.debugln("buildCreateScanJSON(): Building JSON for request ID #" + requestedScan.RequestID)
		go func(requestedScan RequestedScan, wg *sync.WaitGroup, createScanJSONch chan nessusAPI.CreateScan) {
			scanNameAndDescription := fmt.Sprintf("Scan Request #%s, Method: %s", requestedScan.RequestID, requestedScan.Method)
			targets := strings.Join(requestedScan.TargetIPs, " ")
			policyID := c.scanMethodToPolicyID(requestedScan.Method)

			createScanJSONch <- nessusAPI.CreateScan{
				UUID: "ab4bacd2-05f6-425c-9d79-3ba3940ad1c24e51e1f403febe40",
				Settings: nessusAPI.CreateScanSettings{
					Name:        scanNameAndDescription,
					Description: scanNameAndDescription,
					FolderID:    "65",
					ScannerID:   "1",
					PolicyID:    policyID,
					TextTargets: targets,
					Launch:      "ONETIME",
					Enabled:     false,
					LaunchNow:   false,
					Emails:      "",
				},
			}
			wg.Done()
		}(requestedScan, wg, createScanJSONch)
	}

	go func(wg *sync.WaitGroup, createScanJSONch chan nessusAPI.CreateScan) {
		wg.Wait()
		close(createScanJSONch)
	}(wg, createScanJSONch)

	return createScanJSONch
}

// scanMethodToPolicyID takes a method which was extracted from
// processRequestedScanFile() and returns the corresponding Nessus Policy ID.
func (c *Creator) scanMethodToPolicyID(method string) string {
	c.debugln("scanMethodToPolicyID(): Determining policy ID from scan method")
	switch method {
	case "allportswithping":
		return "52"
	case "allportsnoping":
		return "53"
	case "atomic":
		return "54"
	case "pci":
		return "55"
	default:
		return "19"
	}
}