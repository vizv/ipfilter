package ipfilter

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vizv/ipfilter/utils/files"
	"github.com/vizv/ipfilter/utils/iprange"
	"github.com/vizv/ipfilter/utils/parser"
)

var flagOutput string

var MergeCmd = &cobra.Command{
	Use:   "merge IPFILTER_DAT_FILE...",
	Short: "Merge ipfilter.dat files.",
	Long:  `Merge rules from multiple ipfilter.dat files, and generate a single ipfilter.dat file.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("collecting rules...")
		intervals := iprange.Intervals{}
		filesCount := 0
		rulesCount := 0
		for _, file := range files.GlobFiles(args) {
			log.Infof(`collecting rules from "%s"...`, file)
			for rule := range parser.ParseIPFilterDatFile(file) {
				from, to := rule[0], rule[1]
				log.WithFields(log.Fields{"from": from, "to": to}).Tracef("read rule")
				intervals.Append(from, to)
				rulesCount += 1
			}
			filesCount += 1
		}
		log.Infof("%d rules collected from %d files.", rulesCount, filesCount)

		log.Infof("merging rules...")
		intervals = intervals.Merge()
		mergedCount := len(intervals)
		log.Infof("merged to %d rules.", mergedCount)

		outputFilename := flagOutput
		log.Infof(`saving rules to "%s"...`, outputFilename)
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer outputFile.Close()

		for _, interval := range intervals {
			from, to := interval.From, interval.To
			log.WithFields(log.Fields{"from": from, "to": to}).Tracef("write rule")
			fmt.Fprintf(outputFile, "%s - %s , 0 , \n", from, to)
		}
		log.Infof(`merged rules saved to "%s".`, outputFilename)
	},
}

func init() {
	MergeCmd.Flags().StringVarP(&flagOutput, "output", "o", "ipfilter.dat", "Output path for merged ipfilter.dat. (default: ipfilter.dat)")
}
