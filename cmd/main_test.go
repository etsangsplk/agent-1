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

package main

import (
	"os"
	"testing"
	"time"

	"github.com/optimizely/agent/config"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func assertRoot(t *testing.T, actual *config.AgentConfig) {
	assert.Equal(t, "0.1.0", actual.Version)
	assert.Equal(t, "Optimizely Inc.", actual.Author)
	assert.Equal(t, "optimizely", actual.Name)
	assert.Equal(t, []string{"ddd", "eee", "fff"}, actual.SDKKeys)
}

func assertServer(t *testing.T, actual config.ServerConfig) {
	assert.Equal(t, 5*time.Second, actual.ReadTimeout)
	assert.Equal(t, 10*time.Second, actual.WriteTimeout)
	assert.Equal(t, "keyfile", actual.KeyFile)
	assert.Equal(t, "certfile", actual.CertFile)
}

func assertLog(t *testing.T, actual config.LogConfig) {
	assert.True(t, actual.Pretty)
	assert.Equal(t, "debug", actual.Level)
}

func assertAdmin(t *testing.T, actual config.AdminConfig) {
	assert.Equal(t, "3002", actual.Port)
}

func assertAdminAuth(t *testing.T, actual config.ServiceAuthConfig) {
	assert.Equal(t, 30*time.Minute, actual.TTL)
	assert.Equal(t, "efgh", actual.HMACSecret)
	assert.Equal(t, config.OAuthClientCredentials{
		ID:     "clientid2",
		Secret: "clientsecret2",
	}, actual.Clients[0])
}

func assertAPI(t *testing.T, actual config.APIConfig) {
	assert.Equal(t, 100, actual.MaxConns)
	assert.Equal(t, "3000", actual.Port)
}

func assertAPIAuth(t *testing.T, actual config.ServiceAuthConfig) {
	assert.Equal(t, 30*time.Minute, actual.TTL)
	assert.Equal(t, "abcd", actual.HMACSecret)
	assert.Equal(t, config.OAuthClientCredentials{
		ID:     "clientid1",
		Secret: "clientsecret1",
	}, actual.Clients[0])
}

func assertWebhook(t *testing.T, actual config.WebhookConfig) {
	assert.Equal(t, "3001", actual.Port)
	assert.Equal(t, "secret-10000", actual.Projects[10000].Secret)
	assert.Equal(t, []string{"aaa", "bbb", "ccc"}, actual.Projects[10000].SDKKeys)
	assert.True(t, actual.Projects[10000].SkipSignatureCheck)
	assert.Equal(t, "secret-20000", actual.Projects[20000].Secret)
	assert.Equal(t, []string{"xxx", "yyy", "zzz"}, actual.Projects[20000].SDKKeys)
	assert.False(t, actual.Projects[20000].SkipSignatureCheck)
}

func TestViperYaml(t *testing.T) {
	v := viper.New()
	v.Set("config.filename", "./testdata/default.yaml")

	actual := loadConfig(v)

	assertRoot(t, actual)
	assertServer(t, actual.Server)
	assertLog(t, actual.Log)
	assertAdmin(t, actual.Admin)
	assertAdminAuth(t, actual.Admin.Auth)
	assertAPI(t, actual.API)
	assertAPIAuth(t, actual.API.Auth)
	assertWebhook(t, actual.Webhook)
}

func TestViperProps(t *testing.T) {
	v := viper.New()

	v.Set("version", "0.1.0")
	v.Set("author", "Optimizely Inc.")
	v.Set("name", "optimizely")
	v.Set("sdkkeys", []string{"ddd", "eee", "fff"})

	v.Set("server.readtimeout", 5*time.Second)
	v.Set("server.writetimeout", 10*time.Second)
	v.Set("server.certfile", "certfile")
	v.Set("server.keyfile", "keyfile")

	v.Set("log.pretty", true)
	v.Set("log.level", "debug")

	v.Set("admin.port", "3002")
	v.Set("admin.auth.ttl", "30m")
	v.Set("admin.auth.hmacsecret", "efgh")
	v.Set("admin.auth.clients", []map[string]string{
		{
			"id":     "clientid2",
			"secret": "clientsecret2",
		},
	})

	v.Set("api.maxconns", 100)
	v.Set("api.port", "3000")
	v.Set("api.auth.ttl", "30m")
	v.Set("api.auth.hmacsecret", "abcd")
	v.Set("api.auth.clients", []map[string]string{
		{
			"id":     "clientid1",
			"secret": "clientsecret1",
		},
	})

	v.Set("webhook.port", "3001")
	v.Set("webhook.projects.10000.secret", "secret-10000")
	v.Set("webhook.projects.10000.sdkkeys", []string{"aaa", "bbb", "ccc"})
	v.Set("webhook.projects.10000.skipsignaturecheck", true)
	v.Set("webhook.projects.20000.secret", "secret-20000")
	v.Set("webhook.projects.20000.sdkkeys", []string{"xxx", "yyy", "zzz"})
	v.Set("webhook.projects.20000.skipsignaturecheck", false)

	assert.NoError(t, initConfig(v))
	actual := loadConfig(v)

	assertRoot(t, actual)
	assertServer(t, actual.Server)
	assertLog(t, actual.Log)
	assertAdmin(t, actual.Admin)
	assertAdminAuth(t, actual.Admin.Auth)
	assertAPI(t, actual.API)
	assertAPIAuth(t, actual.API.Auth)
	assertWebhook(t, actual.Webhook)
}

func TestViperEnv(t *testing.T) {
	_ = os.Setenv("OPTIMIZELY_VERSION", "0.1.0")
	_ = os.Setenv("OPTIMIZELY_AUTHOR", "Optimizely Inc.")
	_ = os.Setenv("OPTIMIZELY_NAME", "optimizely")
	_ = os.Setenv("OPTIMIZELY_SDKKEYS", "ddd,eee,fff")

	_ = os.Setenv("OPTIMIZELY_SERVER_READTIMEOUT", "5s")
	_ = os.Setenv("OPTIMIZELY_SERVER_WRITETIMEOUT", "10s")
	_ = os.Setenv("OPTIMIZELY_SERVER_CERTFILE", "certfile")
	_ = os.Setenv("OPTIMIZELY_SERVER_KEYFILE", "keyfile")

	_ = os.Setenv("OPTIMIZELY_LOG_PRETTY", "true")
	_ = os.Setenv("OPTIMIZELY_LOG_LEVEL", "debug")

	_ = os.Setenv("OPTIMIZELY_ADMIN_PORT", "3002")

	_ = os.Setenv("OPTIMIZELY_API_MAXCONNS", "100")
	_ = os.Setenv("OPTIMIZELY_API_PORT", "3000")

	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PORT", "3001")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_10000_SECRET", "secret-10000")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_10000_SDKKEYS", "aaa,bbb,ccc")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_10000_SKIPSIGNATURECHECK", "true")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_20000_SECRET", "secret-20000")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_20000_SDKKEYS", "xxx,yyy,zzz")
	_ = os.Setenv("OPTIMIZELY_WEBHOOK_PROJECTS_20000_SKIPSIGNATURECHECK", "false")

	v := viper.New()
	assert.NoError(t, initConfig(v))
	actual := loadConfig(v)

	assertRoot(t, actual)
	assertServer(t, actual.Server)
	assertLog(t, actual.Log)
	assertAdmin(t, actual.Admin)
	assertAPI(t, actual.API)
	//assertWebhook(t, actual.Webhook) // Maps don't appear to be supported
}
