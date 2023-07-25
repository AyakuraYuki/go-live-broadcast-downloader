package sequence

import (
	"go-live-broadcast-downloader/plugins/log"
	"go-live-broadcast-downloader/plugins/net/ip"
	"go-live-broadcast-downloader/plugins/sequence/snowflake"
	"os"
	"strconv"
)

const (
	ErrType = "tools-sequence"
)

var (
	sequence *snowflake.Node
)

func init() {
	var sid int64
	ssn := os.Getenv("SEQUENCE_SN")
	if ssn != "" {
		ssni, err := strconv.ParseInt(ssn, 10, 64)
		if err != nil {
			panic(err)
		}
		sid = ssni
	} else {
		// ipu := uint16(1)<<8 + uint16(177)
		ipu, err := ip.Lower16BitPrivateIP()
		if err != nil {
			log.Error(ErrType).Msgf("Lower16BitPrivateIP err:%v", err)
		}
		sid = int64(ipu)
	}
	sn, err := snowflake.NewNode(sid)
	if err != nil {
		log.Error(ErrType).Msgf("NewNode err:%v", err)
	}
	sequence = sn
}

// ID id
func ID() int64 {
	sid := sequence.Generate()
	return sid.Int64()
}
