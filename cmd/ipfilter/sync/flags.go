package sync

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Interval() time.Duration {
	interval := viper.GetString("sync.interval")

	normalizedDuration := interval
	if _, err := strconv.ParseInt(normalizedDuration, 0, 64); err == nil {
		normalizedDuration += "s"
	}

	updateInterval, err := time.ParseDuration(normalizedDuration)
	if err != nil || updateInterval < 0 {
		log.WithField("interval", interval).Warnf("failed to parse update interval, use default interval - %s", DEFAULT_UPDATE_INTERVAL)
		updateInterval, _ = time.ParseDuration(DEFAULT_UPDATE_INTERVAL)
	}

	return updateInterval
}

func WebUIURL() *url.URL {
	webUIURL := viper.GetString("sync.webui-url")
	username := viper.GetString("sync.username")
	password := viper.GetString("sync.password")

	if webUIURL == "" {
		return nil
	}

	normalizedURL := webUIURL
	if !strings.HasPrefix(strings.ToLower(normalizedURL), "http://") && !strings.HasPrefix(strings.ToLower(normalizedURL), "https://") {
		normalizedURL = "http://" + normalizedURL
	}

	parsedURL, err := url.ParseRequestURI(normalizedURL)
	if err != nil {
		log.WithField("url", webUIURL).Warnf("invalid WebUI URL, disable notify")
		return nil
	}

	if username != "" && password != "" {
		parsedURL.User = url.UserPassword(username, password)
	}

	parsedURL.Path = ""

	return parsedURL
}
