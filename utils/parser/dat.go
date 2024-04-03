package parser

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ParseIPFilterDatFile(filename string) <-chan [2]string {
	ch := make(chan [2]string)

	go func() {
		defer close(ch)

		file, err := os.Open(filename)
		if err != nil {
			log.WithField("filename", filename).Warnf("failed to open file: %+v", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			parts := strings.SplitN(line, ",", 2)
			ips := strings.SplitN(parts[0], "-", 3)

			from := strings.TrimSpace(ips[0])
			to := strings.TrimSpace(ips[1])

			ch <- [2]string{from, to}
		}

		if err := scanner.Err(); err != nil {
			log.WithField("filename", filename).Warnf("failed to read file: %+v", err)
			return
		}
	}()

	return ch
}
