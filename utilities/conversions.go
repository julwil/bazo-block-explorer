package utilities

import (
  "fmt"
  "time"
  "github.com/bazo-blockchain/bazo-miner/protocol"
)

func ConvertBlock(unconvertedBlock *protocol.Block) Block {
  var convertedBlock Block
  var convertedTxHash string

  //convertedBlock.Header = fmt.Sprintf("%x", unconvertedBlock.Header)
  convertedBlock.Hash = fmt.Sprintf("%x", unconvertedBlock.Hash)
  convertedBlock.PrevHash = fmt.Sprintf("%x", unconvertedBlock.PrevHash)
  //convertedBlock.Nonce = fmt.Sprintf("%x", unconvertedBlock.Nonce)
  convertedBlock.Timestamp = unconvertedBlock.Timestamp
  convertedBlock.TimeString = time.Unix(unconvertedBlock.Timestamp, 0).Format("02 Jan 2006 15:04")
  convertedBlock.MerkleRoot = fmt.Sprintf("%x", unconvertedBlock.MerkleRoot)
  convertedBlock.Beneficiary = fmt.Sprintf("%x", unconvertedBlock.Beneficiary)
  convertedBlock.NrFundsTx = unconvertedBlock.NrFundsTx
  convertedBlock.NrAccTx = unconvertedBlock.NrAccTx
  convertedBlock.NrConfigTx = unconvertedBlock.NrConfigTx

  for _, hash := range unconvertedBlock.FundsTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.FundsTxData = append(convertedBlock.FundsTxData, convertedTxHash)
  }
  for _, hash := range unconvertedBlock.AccTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.AccTxData = append(convertedBlock.AccTxData, convertedTxHash)
  }
  for _, hash := range unconvertedBlock.ConfigTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.ConfigTxData = append(convertedBlock.ConfigTxData, convertedTxHash)
  }

  return convertedBlock
}

func ConvertFundsTransaction(unconvertedTx *protocol.FundsTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) Fundstx {
  var convertedTx Fundstx

  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Amount = unconvertedTx.Amount
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.TxCount = unconvertedTx.TxCnt
  convertedTx.From = fmt.Sprintf("%x", unconvertedTx.From)
  convertedTx.To = fmt.Sprintf("%x", unconvertedTx.To)
  convertedTx.Timestamp = time.Now().Unix()
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig1)

  return convertedTx
}

func ConvertAccTransaction(unconvertedTx *protocol.AccTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) Acctx {
  var convertedTx Acctx

  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.Issuer = fmt.Sprintf("%x", unconvertedTx.Issuer)
  convertedTx.PubKey = fmt.Sprintf("%x", unconvertedTx.PubKey)
  convertedTx.Timestamp = time.Now().Unix()
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}

func ConvertConfigTransaction(unconvertedTx *protocol.ConfigTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) Configtx {
  var convertedTx Configtx

  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Id = unconvertedTx.Id
  convertedTx.Fee = unconvertedTx.Fee
  if unconvertedTx.Payload > 10000000 {
    convertedTx.Payload = 10000000
  } else {
    convertedTx.Payload = unconvertedTx.Payload
  }
  convertedTx.TxCount = unconvertedTx.TxCnt
  convertedTx.Timestamp = time.Now().Unix()
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}

func ConvertStakeTransaction(unconvertedTx *protocol.StakeTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) Staketx  {
  var convertedTx Staketx

  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Timestamp = time.Now().Unix()
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.Account = fmt.Sprintf("%x", unconvertedTx.Account)
  convertedTx.IsStaking = unconvertedTx.IsStaking
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}
