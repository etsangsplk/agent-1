/****************************************************************************
 * Copyright 2019-2020, Optimizely, Inc. and contributors                   *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package config contains all the configuration attributes for running Optimizely Agent
package config

import (
	"time"
)

// NewDefaultConfig returns the default configuration for Optimizely Agent
func NewDefaultConfig() *AgentConfig {

	config := AgentConfig{
		Version: "",
		Author:  "Optimizely Inc.",
		Name:    "optimizely",

		Admin: AdminConfig{
			Auth: ServiceAuthConfig{
				Clients:    make([]OAuthClientCredentials, 0),
				HMACSecret: "",
				TTL:        0,
			},
			Port: "8088",
		},
		API: APIConfig{
			Auth: ServiceAuthConfig{
				Clients:    make([]OAuthClientCredentials, 0),
				HMACSecret: "",
				TTL:        0,
			},
			MaxConns: 0,
			Port:     "8080",
		},
		Log: LogConfig{
			Pretty: false,
			Level:  "info",
		},
		Processor: ProcessorConfig{
			BatchSize:     10,
			QueueSize:     1000,
			FlushInterval: 30 * time.Second,
		},
		Server: ServerConfig{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			CertFile:     "",
			KeyFile:      "",
		},
		Webhook: WebhookConfig{
			Port: "8085",
		},
	}

	return &config
}

// AgentConfig is the top level configuration struct
type AgentConfig struct {
	Version string `json:"version"`
	Author  string `json:"author"`
	Name    string `json:"name"`

	SDKKeys []string `yaml:"sdkkeys" json:"sdkkeys"`

	Admin     AdminConfig     `json:"admin"`
	API       APIConfig       `json:"api"`
	Log       LogConfig       `json:"log"`
	Processor ProcessorConfig `json:"processor"`
	Server    ServerConfig    `json:"server"`
	Webhook   WebhookConfig   `json:"webhook"`
}

// ProcessorConfig holds the configuration options for the Optimizely Event Processor.
type ProcessorConfig struct {
	BatchSize     int           `json:"batchSize" default:"10"`
	QueueSize     int           `json:"queueSize" default:"1000"`
	FlushInterval time.Duration `json:"flushInterval" default:"30s"`
}

// LogConfig holds the log configuration
type LogConfig struct {
	Pretty bool   `json:"pretty"`
	Level  string `json:"level"`
}

// ServerConfig holds the global http server configs
type ServerConfig struct {
	ReadTimeout  time.Duration `json:"readtimeout"`
	WriteTimeout time.Duration `json:"writetimeout"`
	CertFile     string        `json:"certfile"`
	KeyFile      string        `json:"keyfile"`
}

// APIConfig holds the REST API configuration
type APIConfig struct {
	Auth     ServiceAuthConfig `json:"-"`
	MaxConns int               `json:"maxconns"`
	Port     string            `json:"port"`
}

// AdminConfig holds the configuration for the admin web interface
type AdminConfig struct {
	Auth ServiceAuthConfig `json:"-"`
	Port string            `json:"port"`
}

// WebhookConfig holds configuration for Optimizely Webhooks
type WebhookConfig struct {
	Port     string                   `json:"port"`
	Projects map[int64]WebhookProject `json:"projects"`
}

// WebhookProject holds the configuration for a single Project webhook
type WebhookProject struct {
	SDKKeys            []string `json:"sdkKeys"`
	Secret             string   `json:"-"`
	SkipSignatureCheck bool     `json:"skipSignatureCheck" default:"false"`
}

// OAuthClientCredentials are used for issuing access tokens
type OAuthClientCredentials struct {
	ID     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

// ServiceAuthConfig holds the authentication configuration for a particular service
type ServiceAuthConfig struct {
	Clients    []OAuthClientCredentials `yaml:"clients" json:"-"`
	HMACSecret string                   `yaml:"hmacSecret" json:"-"`
	TTL        time.Duration            `yaml:"ttl" json:"-"`
}
