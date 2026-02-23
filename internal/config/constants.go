// Copyright (c) 2015-2021 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

// Config value separator
const (
	ValueSeparator = ","
)

// Top level common ENVs
const (
	EnvAccessKey    = "S3_ACCESS_KEY"
	EnvSecretKey    = "S3_SECRET_KEY"
	EnvRootUser     = "S3_ROOT_USER"
	EnvRootPassword = "S3_ROOT_PASSWORD"

	// Legacy files
	EnvAccessKeyFile = "S3_ACCESS_KEY_FILE"
	EnvSecretKeyFile = "S3_SECRET_KEY_FILE"

	// Current files
	EnvRootUserFile     = "S3_ROOT_USER_FILE"
	EnvRootPasswordFile = "S3_ROOT_PASSWORD_FILE"

	// Set all config environment variables from 'config.env'
	// if necessary. Overrides all previous settings and also
	// overrides all environment values passed from
	// 'podman run -e ENV=value'
	EnvConfigEnvFile = "S3_CONFIG_ENV_FILE"

	EnvBrowser    = "S3_BROWSER"
	EnvDomain     = "S3_DOMAIN"
	EnvPublicIPs  = "S3_PUBLIC_IPS"
	EnvFSOSync    = "S3_FS_OSYNC"
	EnvArgs       = "S3_ARGS"
	EnvVolumes    = "S3_VOLUMES"
	EnvDNSWebhook = "S3_DNS_WEBHOOK_ENDPOINT"

	EnvSiteName   = "S3_SITE_NAME"
	EnvSiteRegion = "S3_SITE_REGION"

	EnvMinIOSubnetLicense = "S3_SUBNET_LICENSE" // Deprecated Dec 2021
	EnvMinIOSubnetAPIKey  = "S3_SUBNET_API_KEY"
	EnvMinIOSubnetProxy   = "S3_SUBNET_PROXY"

	EnvMinIOCallhomeEnable    = "S3_CALLHOME_ENABLE"
	EnvMinIOCallhomeFrequency = "S3_CALLHOME_FREQUENCY"

	EnvMinIOServerURL             = "S3_SERVER_URL"
	EnvBrowserRedirect            = "S3_BROWSER_REDIRECT" // On by default
	EnvBrowserRedirectURL         = "S3_BROWSER_REDIRECT_URL"
	EnvRootDriveThresholdSize     = "S3_ROOTDRIVE_THRESHOLD_SIZE"
	EnvRootDiskThresholdSize      = "S3_ROOTDISK_THRESHOLD_SIZE" // Deprecated Sep 2023
	EnvBrowserLoginAnimation      = "S3_BROWSER_LOGIN_ANIMATION"
	EnvBrowserSessionDuration     = "S3_BROWSER_SESSION_DURATION" // Deprecated after November 2023
	EnvMinioStsDuration           = "S3_STS_DURATION"
	EnvMinIOLogQueryURL           = "S3_LOG_QUERY_URL"
	EnvMinIOLogQueryAuthToken     = "S3_LOG_QUERY_AUTH_TOKEN"
	EnvMinIOPrometheusURL         = "S3_PROMETHEUS_URL"
	EnvMinIOPrometheusJobID       = "S3_PROMETHEUS_JOB_ID"
	EnvMinIOPrometheusExtraLabels = "S3_PROMETHEUS_EXTRA_LABELS"
	EnvMinIOPrometheusAuthToken   = "S3_PROMETHEUS_AUTH_TOKEN"
	EnvConsoleDebugLogLevel       = "S3_CONSOLE_DEBUG_LOGLEVEL"

	EnvUpdate = "S3_UPDATE"

	EnvEndpoints  = "S3_ENDPOINTS"   // legacy
	EnvWorm       = "S3_WORM"        // legacy
	EnvRegion     = "S3_REGION"      // legacy
	EnvRegionName = "S3_REGION_NAME" // legacy

)

// Expiration Token durations
// These values are used to validate the expiration time range from
// either the exp claim or MINI_STS_DURATION value
const (
	MinExpiration = 900
	MaxExpiration = 31536000
)
