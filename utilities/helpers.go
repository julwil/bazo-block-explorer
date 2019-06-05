package utilities

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/sha3"
	"time"
)

func ExtractParameters(tx Configtx, currentParams Systemparams) Systemparams {
	currentParams.Timestamp = time.Now().Unix()

	switch tx.Id {
	case 1:
		currentParams.BlockSize = tx.Payload
	case 2:
		currentParams.DiffInterval = tx.Payload
	case 3:
		currentParams.MinFee = tx.Payload
	case 4:
		currentParams.BlockInterval = tx.Payload
	case 5:
		currentParams.BlockReward = tx.Payload
	case 6:
		currentParams.StakingMin = tx.Payload
	case 7:
		currentParams.WaitingMin = tx.Payload
	case 8:
		currentParams.AcceptanceTimeDiff = tx.Payload
	case 9:
		currentParams.SlashingWindowSize = tx.Payload
	case 10:
		currentParams.SlashingReward = tx.Payload
	default:
		return currentParams
	}
	return currentParams
}

func SerializeHashContent(data interface{}) (hash [32]byte) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, data)

	return sha3.Sum256(buf.Bytes())
}
