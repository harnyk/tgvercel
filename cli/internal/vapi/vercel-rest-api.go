package vapi

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	tokenEnvName = "VERCEL_TOKEN"
	vercelApiUrl = "https://api.vercel.com"
)

type Client struct {
	options Options
	rest    *resty.Client
}

func NewClient() *Client {
	return NewClientWithOptions(Options{})
}

func NewClientWithOptions(options Options) *Client {
	restyClient := resty.New().SetBaseURL(vercelApiUrl)

	return &Client{
		options: options,
		rest:    restyClient,
	}
}

func (c *Client) SetEnv(projectID, key, value string, target Target) error {
	// /v10/projects/{projectIdOrName}/env?upsert=true"

	token := c.options.Token

	payload := map[string]any{
		"key":   key,
		"value": value,
		"type":  "encrypted",
		"target": []string{
			string(target),
		},
		"comment": "Created and used by tgvercel",
	}

	resp, err := c.rest.R().
		SetAuthToken(token).
		SetBody(payload).
		Post(fmt.Sprintf("/v10/projects/%s/env?upsert=true", projectID))

	if err != nil {
		return fmt.Errorf("failed to set env: %w", err)
	}
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return fmt.Errorf("failed to set env: %s", resp.String())
	}

	return nil
}

func (c *Client) GetEnv(projectID, key string, target Target) (string, error) {
	encryptedEnvDescriptor, err := c.getEnvDescriptor(projectID, key, target)
	if err != nil {
		return "", fmt.Errorf("failed to get env descriptor for key %s: %w", key, err)
	}

	decryptedEnv, err := c.getDecryptedEnv(projectID, encryptedEnvDescriptor.ID)
	if err != nil {
		return "", fmt.Errorf("failed to get decrypted env for key %s: %w", key, err)
	}

	return decryptedEnv, nil
}

func (c *Client) GetDeployment(idOrUrl string) (*Deployment, error) {
	// /v13/deployments/{idOrUrl}

	token := c.options.Token

	var deployment *Deployment
	resp, err := c.rest.R().
		SetAuthToken(token).
		SetResult(&deployment).
		Get(fmt.Sprintf("/v13/deployments/%s", idOrUrl))

	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get deployment: %s", resp.Status())
	}

	return deployment, nil
}

func (c *Client) getDecryptedEnv(projectID, envId string) (string, error) {
	// /v1/projects/{idOrName}/env/{id}

	token := c.options.Token

	var envDescriptor *EnvDescriptor
	resp, err := c.rest.R().
		SetAuthToken(token).
		SetResult(&envDescriptor).
		Get(fmt.Sprintf("/v1/projects/%s/env/%s", projectID, envId))

	if err != nil {
		return "", fmt.Errorf("failed to get env descriptor: %w", err)
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("failed to get env descriptor: %s", resp.Status())
	}

	return envDescriptor.Value, nil
}

func (c *Client) getEnvDescriptor(projectID, key string, target Target) (*EnvDescriptor, error) {
	// /v9/projects/{PROJECT_ID}/env

	token := c.options.Token

	var envs EnvsResponse

	resp, err := c.rest.R().
		SetAuthToken(token).
		SetResult(&envs).
		Get(fmt.Sprintf("/v9/projects/%s/env", projectID))

	if err != nil {
		return nil, fmt.Errorf("failed to get env descriptor: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get env descriptor: %s", resp.Status())
	}

	var envDescriptor *EnvDescriptor
	for _, env := range envs.Envs {
		if env.Key == key {
			envDescriptor = &env
			break
		}
	}

	if envDescriptor == nil {
		return nil, fmt.Errorf("env %s not found", key)
	}

	targets := envDescriptor.Target
	if len(targets) == 0 {
		return nil, fmt.Errorf("env %s has no target", key)
	}

	for _, envTarget := range targets {
		if envTarget == target {
			return envDescriptor, nil
		}
	}

	return nil, fmt.Errorf("env %s not found in target %s", key, target)
}
