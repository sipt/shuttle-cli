package apis

// // DumpStatus dump status structure
// type DumpStatus struct {
// 	Dump bool `json:"dump"`
// 	MITM bool `json:"mitm"`
// }

// // DumpStatusRequest dump status request structure
// type DumpStatusRequest struct {
// 	Dump bool `json:"dump"`
// 	MITM bool `json:"mitm"`
// }

// // GetDumpStatus gets stream dump status
// func (c *APIClient) GetDumpStatus() error {
// 	result := &model.Response[DumpStatus]{}
// 	resp, err := c.client.R().
// 		SetResult(result).
// 		Get("/api/stream/dump/status")

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	if result.Code != 0 {
// 		return fmt.Errorf("API error: %s", result.Message)
// 	}

// 	status := result.Data
// 	fmt.Println("Stream Dump Status:")
// 	fmt.Printf("  Dump: ")
// 	if status.Dump {
// 		green.Println("Enabled")
// 	} else {
// 		red.Println("Disabled")
// 	}
// 	fmt.Printf("  MITM: ")
// 	if status.MITM {
// 		green.Println("Enabled")
// 	} else {
// 		red.Println("Disabled")
// 	}

// 	return nil
// }

// // SetDumpStatus sets stream dump status
// func (c *APIClient) SetDumpStatus(dump, mitm bool) error {
// 	request := DumpStatusRequest{
// 		Dump: dump,
// 		MITM: mitm,
// 	}

// 	result := &model.Response[any]{}
// 	resp, err := c.client.R().
// 		SetResult(result).
// 		SetBody(request).
// 		Put("/api/stream/dump/status")

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	if result.Code == 0 {
// 		green.Println("Dump status updated successfully")
// 		fmt.Printf("  Dump: ")
// 		if dump {
// 			green.Println("Enabled")
// 		} else {
// 			red.Println("Disabled")
// 		}
// 		fmt.Printf("  MITM: ")
// 		if mitm {
// 			green.Println("Enabled")
// 		} else {
// 			red.Println("Disabled")
// 		}
// 	} else {
// 		fmt.Print("Set dump status failed: ")
// 		red.Println(result.Message)
// 	}

// 	return nil
// }

// // GenerateCACert generates CA certificate
// func (c *APIClient) GenerateCACert() error {
// 	result := &model.Response[any]{}
// 	resp, err := c.client.R().
// 		SetResult(result).
// 		Post("/api/stream/dump/cert")

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	if result.Code == 0 {
// 		green.Println("CA certificate generated successfully")
// 	} else {
// 		fmt.Print("Generate CA certificate failed: ")
// 		red.Println(result.Message)
// 	}

// 	return nil
// }

// // DownloadCACert downloads CA certificate
// func (c *APIClient) DownloadCACert(output string) error {
// 	resp, err := c.client.R().
// 		Get("/api/stream/dump/cert")

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	// Default output path
// 	if output == "" {
// 		homeDir, err := os.UserHomeDir()
// 		if err != nil {
// 			output = "shuttle-ca.crt"
// 		} else {
// 			output = filepath.Join(homeDir, "shuttle-ca.crt")
// 		}
// 	}

// 	err = os.WriteFile(output, resp.Body(), 0644)
// 	if err != nil {
// 		return fmt.Errorf("failed to write certificate: %v", err)
// 	}

// 	green.Printf("CA certificate downloaded to: %s\n", output)
// 	return nil
// }

// // GetDumpSession gets session dump data
// func (c *APIClient) GetDumpSession(id string) error {
// 	resp, err := c.client.R().
// 		Get(fmt.Sprintf("/api/stream/dump/session/%s", id))

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	// Try to parse as JSON for pretty printing
// 	var jsonData interface{}
// 	if err := json.Unmarshal(resp.Body(), &jsonData); err == nil {
// 		prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
// 		if err == nil {
// 			fmt.Println(string(prettyJSON))
// 			return nil
// 		}
// 	}

// 	// Fallback to raw output
// 	fmt.Println(string(resp.Body()))
// 	return nil
// }

// // GetDumpRequestBody gets request body data
// func (c *APIClient) GetDumpRequestBody(id string) error {
// 	resp, err := c.client.R().
// 		Get(fmt.Sprintf("/api/stream/dump/request/body/%s", id))

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	fmt.Println("Request Body:")
// 	fmt.Println(string(resp.Body()))
// 	return nil
// }

// // GetDumpResponseBody gets response body data
// func (c *APIClient) GetDumpResponseBody(id string) error {
// 	resp, err := c.client.R().
// 		Get(fmt.Sprintf("/api/stream/dump/response/body/%s", id))

// 	if err != nil {
// 		return fmt.Errorf("request failed: %v", err)
// 	}

// 	if resp.StatusCode() != 200 {
// 		return fmt.Errorf("HTTP error: %d", resp.StatusCode())
// 	}

// 	fmt.Println("Response Body:")
// 	fmt.Println(string(resp.Body()))
// 	return nil
// }
