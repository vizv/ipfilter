package ipfilter

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vizv/ipfilter/cmd/ipfilter/sync"
	"github.com/vizv/ipfilter/utils/hash"
	"github.com/vizv/ipfilter/utils/iprange"
	"github.com/vizv/ipfilter/utils/json"
	"github.com/vizv/ipfilter/utils/parser"
	"github.com/vizv/ipfilter/utils/qb"
)

var SyncCmd = &cobra.Command{
	Use:   "sync [IPFILTER_DAT_FILE_URL...]",
	Short: "Synchronize ipfilter.dat files.",
	Long:  `Synchronize rules from multiple remote ipfilter.dat files, and optionally notify qBittorrent.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		updateInterval := sync.Interval()
		runOnce := updateInterval == 0
		log.Debugf("updateInterval: %+v", updateInterval)
		log.Debugf("runOnce: %+v", runOnce)

		cacheDir := viper.GetString("sync.cache-dir")
		outputDir := viper.GetString("sync.output-dir")
		log.Debugf("cacheDir: %+v", cacheDir)
		log.Debugf("outputDir: %+v", outputDir)

		rawDATURLs := args
		if len(rawDATURLs) == 0 {
			rawDATURLs = strings.Split(viper.GetString("sync.dat-urls"), ",")
		}
		if len(rawDATURLs) == 0 {
			rawDATURLs = []string{sync.DEFAULT_IPFILTER_DAT_FILE_URL}
		}
		datURLsWithCachePath := map[string]string{}
		for _, datURL := range rawDATURLs {
			if _, err := url.ParseRequestURI(datURL); err != nil {
				log.WithField("url", datURL).Warnf("ignore invalid filter.dat URL")
				continue
			}
			cacheFilename := fmt.Sprintf("ipfilter-%s.dat", hash.CalculateMD5([]byte(datURL)))
			cachePath := path.Join(cacheDir, cacheFilename)
			datURLsWithCachePath[datURL] = cachePath
		}
		if len(datURLsWithCachePath) == 0 {
			log.Fatalf("no valid filter.dat URL found.")
		}
		log.Debugf("rawDATURLs: %+v", rawDATURLs)
		log.Debugf("datURLs: %+v", datURLsWithCachePath)

		webUIURL := sync.WebUIURL()
		notifyQB := webUIURL != nil
		prefPath := ""
		var qbClient *qb.Client
		if notifyQB {
			if client, err := qb.NewClient(webUIURL); err != nil {
				log.WithField("url", webUIURL).Warnf("failed to create qBittorrent client, disable notifyQB.")
				qbClient = nil
				notifyQB = false
			} else {
				if prefJson, err := client.GetPreferences(); err != nil {
					log.WithField("url", webUIURL).Warnf("failed to get preferences from qBittorrent client, disable notifyQB.")
					qbClient = nil
					notifyQB = false
				} else {
					ipFilterEnabled := json.GetJsonValueBoolean(prefJson, "ip_filter_enabled")
					ipFilterPath := json.GetJsonValueString(prefJson, "ip_filter_path")
					log.Infof("Current Preferences: ip_filter_enabled = %t, ip_filter_path = %s", ipFilterEnabled, ipFilterPath)

					qbClient = client
					prefPath = ipFilterPath
				}
			}
		}

		log.Debugf("webUIURL: %+v", webUIURL)
		log.Debugf("notifyQB: %+v", notifyQB)

		firstPass := true
		isRetry := false
		for {
			if !firstPass {
				if runOnce {
					break
				}

				if isRetry {
					log.Warnf("retry in %s...", updateInterval)
					isRetry = false
				}

				time.Sleep(updateInterval)
			}
			firstPass = false

			log.Infof(`downloading ipfilter.dat files to "%s"...`, cacheDir)
			totalCount := len(datURLsWithCachePath)
			downloadedCount := 0
			updatedCount := 0
			for datURL, cachePath := range datURLsWithCachePath {
				logFields := log.Fields{"url": datURL, "cache": cachePath}

				log.Infof(`downloading "%s" to "%s"...`, datURL, cachePath)
				datBytes, err := sync.Download(datURL)
				if err != nil {
					log.WithFields(logFields).Warnf("failed to download: %v, skipping...", err)
					continue
				}
				downloadedCount += 1

				cacheBytes, err := os.ReadFile(cachePath)
				if err != nil || !bytes.Equal(datBytes, cacheBytes) {
					if err := os.MkdirAll(cacheDir, 0o755); err != nil {
						log.WithField("dir", cacheDir).Fatalf("failed to create cache directory: %v", err)
					}

					if err := os.WriteFile(cachePath, datBytes, 0o644); err != nil {
						log.WithFields(logFields).Warnf("failed to save: %v, skipping...", err)
						continue
					}
					updatedCount += 1
				}
			}
			log.Infof("%d ipfilter.dat files downloaded from %d URLs, %d files updated.", downloadedCount, totalCount, updatedCount)

			log.Infof("collecting rules...")
			intervals := iprange.Intervals{}
			rulesCount := 0
			for _, file := range datURLsWithCachePath {
				log.Infof(`collecting rules from "%s"...`, file)
				for rule := range parser.ParseIPFilterDatFile(file) {
					from, to := rule[0], rule[1]
					log.WithFields(log.Fields{"from": from, "to": to}).Tracef("read rule")
					intervals.Append(from, to)
					rulesCount += 1
				}
			}
			log.Infof("%d rules collected.", rulesCount)

			log.Infof("merging rules...")
			intervals = intervals.Merge()
			mergedCount := len(intervals)
			log.Infof("merged to %d rules.", mergedCount)

			mergedCachePath := path.Join(cacheDir, "ipfilter-merged.dat")
			log.Infof(`saving rules to "%s"...`, mergedCachePath)
			mergedCacheFile, err := os.Create(mergedCachePath)
			if err != nil {
				log.Fatalf("failed to create output file: %v", err)
			}
			defer mergedCacheFile.Close()

			for _, interval := range intervals {
				from, to := interval.From, interval.To
				log.WithFields(log.Fields{"from": from, "to": to}).Tracef("write rule")
				fmt.Fprintf(mergedCacheFile, "%s - %s , 0 , \n", from, to)
			}
			if err := mergedCacheFile.Sync(); err != nil {
				log.Warnf("failed to write merged ipfilter.dat: %+v", err)
				isRetry = true
				continue
			}
			log.Infof(`merged rules saved to "%s".`, mergedCachePath)

			log.Infof("switching slots...")
			mergedBytes, err := os.ReadFile(mergedCachePath)
			if err != nil {
				log.Warnf("failed to read merged ipfilter.dat: %+v", err)
				isRetry = true
				continue
			}

			outputFilename, currentFilename := sync.GetSlotFiles()
			outputPath, err := filepath.Abs(path.Join(outputDir, outputFilename))
			if err != nil {
				log.Warnf("error getting absolute path for ipfilter.dat to be saved: %v", err)
				isRetry = true
				continue
			}
			currentPath, err := filepath.Abs(path.Join(outputDir, currentFilename))
			if err != nil {
				log.Warnf("error getting absolute path for current ipfilter.dat: %v", err)
				isRetry = true
				continue
			}
			if outputPath == prefPath {
				outputPath, currentPath = currentPath, outputPath
			}
			log.Infof(`switching "%s" to "%s"...`, currentPath, outputPath)

			currentBytes, err := os.ReadFile(currentPath)
			if err != nil || !bytes.Equal(mergedBytes, currentBytes) {
				if err := os.WriteFile(outputPath, mergedBytes, 0o644); err != nil {
					log.Warnf("failed to save ipfilter.dat: %+v", err)
					isRetry = true
					continue
				}

				if notifyQB {
					if err := qbClient.RefreshIPFilter(outputPath); err != nil {
						log.Warnf("error refreshing IP filter: %v", err)
						isRetry = true
						continue
					} else {
						log.Infof("slot switched to %s", outputPath)
					}
				}
			} else {
				log.Infof("ipfilter.dat unchanged, switching slots cancelled.")
			}
		}
	},
}

func init() {
	SyncCmd.Flags().StringP("interval", "i", sync.DEFAULT_UPDATE_INTERVAL, fmt.Sprintf("Synchronize interval. (default: %s)", sync.DEFAULT_UPDATE_INTERVAL))
	viper.BindPFlag("sync.interval", SyncCmd.Flags().Lookup("interval"))

	SyncCmd.Flags().StringP("cache-dir", "c", "cache", "Directory to keep previously downloaded ipfilter.dat files. (default: caches)")
	viper.BindPFlag("sync.cache-dir", SyncCmd.Flags().Lookup("cache-dir"))

	SyncCmd.Flags().StringP("output-dir", "o", ".", "Directory to keep previously downloaded ipfilter.dat files. (default: .)")
	viper.BindPFlag("sync.output-dir", SyncCmd.Flags().Lookup("output-dir"))

	SyncCmd.Flags().StringP("webui-url", "w", "", "qBittorrent WebUI URL to notify the ipfilter.dat changes, leave empty to disable. (empty by default)")
	viper.BindPFlag("sync.webui-url", SyncCmd.Flags().Lookup("webui-url"))

	SyncCmd.Flags().StringP("username", "u", "admin", "Username used to authenticate with qBittorrent WebUI. (default: admin)")
	viper.BindPFlag("sync.username", SyncCmd.Flags().Lookup("username"))

	SyncCmd.Flags().StringP("password", "p", "", "Password used to authenticate with qBittorrent WebUI, leave empty to disable authentication. (empty by default)")
	viper.BindPFlag("sync.password", SyncCmd.Flags().Lookup("password"))

	viper.SetDefault("sync.dat-urls", sync.DEFAULT_IPFILTER_DAT_FILE_URL)
}
