package utilities

import (
  "bytes"
  "time"
  "encoding/binary"
  "golang.org/x/crypto/sha3"
  "github.com/bazo-blockchain/bazo-miner/protocol"
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
  default:
    return currentParams
  }
  return currentParams
}

func invertBlockArray(array []*protocol.Block) []*protocol.Block {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}

	return array
}

func SerializeHashContent(data interface{}) (hash [32]byte) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, data)

	return sha3.Sum256(buf.Bytes())
}
