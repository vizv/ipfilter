package sync

import "os"

const (
	slotAFile = "ipfilter-a.dat"
	slotBFile = "ipfilter-b.dat"
)

func GetSlotFiles() (string, string) {
	slotAStat, err := os.Stat(slotAFile)
	if err != nil {
		return slotAFile, slotBFile
	}
	slotBStat, err := os.Stat(slotBFile)
	if err != nil {
		return slotBFile, slotAFile
	}

	if slotAStat.ModTime().Before(slotBStat.ModTime()) {
		return slotAFile, slotBFile
	} else {
		return slotBFile, slotAFile
	}
}
