package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNaming() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNamingCreate,
		ReadContext:   resourceNamingRead,
		UpdateContext: resourceNamingUpdate,
		DeleteContext: resourceNamingDelete,
		Schema: map[string]*schema.Schema{
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application name for the hostname",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The role associated with the hostname",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment (e.g. dev, prod)",
			},
			"generated_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated hostname",
			},
		},
	}
}

func resourceNamingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Retrieve input variables from Terraform configuration
	application := d.Get("application").(string)
	role := d.Get("role").(string)
	environment := d.Get("environment").(string)

	// Prepare the JSON payload for the new endpoint
	payload := map[string]interface{}{
		"application": application,
		"role":        role,
		"environment": environment,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create the HTTP POST request to the /generate-name endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:32001/generate-name", bytes.NewBuffer(jsonData))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Verify that the API call was successful
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the API response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return diag.FromErr(err)
	}

	// Extract the generated hostname from the response
	generatedName, ok := responseData["hostname"].(string)
	if !ok {
		return diag.Errorf("unexpected response format")
	}

	// Set the generated hostname in Terraform state
	if err := d.Set("generated_name", generatedName); err != nil {
		return diag.FromErr(err)
	}

	// Use the generated hostname as the unique resource ID
	d.SetId(generatedName)

	return diags
}

func resourceNamingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// In this implementation, the generated hostname remains static.
	var diags diag.Diagnostics
	return diags
}

func resourceNamingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For updates, we re-create the resource (i.e. regenerate the hostname)
	return resourceNamingCreate(ctx, d, m)
}

func resourceNamingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// In this example, we simply remove the resource from state.
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}
