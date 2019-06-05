package data

import (
	"encoding/hex"
	"fmt"
	"github.com/bazo-blockchain/bazo-block-explorer/utilities"
	"github.com/bazo-blockchain/bazo-client/network"
	"github.com/bazo-blockchain/bazo-miner/miner"
	"github.com/bazo-blockchain/bazo-miner/p2p"
	"github.com/bazo-blockchain/bazo-miner/protocol"
	"log"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func RunDB() {
	saveInitialParameters()
	network.Uptodate = true
	incomingBlocks()
}

func saveInitialParameters() {
	var convertedParameters utilities.Systemparams
	parameters := miner.NewDefaultParameters()

	convertedParameters.Timestamp = time.Now().Unix()
	convertedParameters.BlockSize = parameters.Block_size
	convertedParameters.DiffInterval = parameters.Diff_interval
	convertedParameters.MinFee = parameters.Fee_minimum
	convertedParameters.BlockInterval = parameters.Block_interval
	convertedParameters.BlockReward = parameters.Diff_interval
	convertedParameters.StakingMin = parameters.Staking_minimum
	convertedParameters.WaitingMin = parameters.Waiting_minimum
	convertedParameters.AcceptanceTimeDiff = parameters.Accepted_time_diff
	convertedParameters.SlashingWindowSize = parameters.Slashing_window_size
	convertedParameters.SlashingReward = parameters.Slash_reward

	WriteParameters(convertedParameters)
}

func incomingBlocks() {
	for {
		blockIn := <-network.BlockIn
		blockInConverted := utilities.ConvertBlock(blockIn)

		//The block does not exist in DB
		if ReturnOneBlock(blockInConverted.Hash).Hash == "" {
			if err := SaveBlockAndTransactions(blockIn, true); err != nil {
				log.Fatal(fmt.Sprintf("Saving block %x and its transactions failed: %v", blockIn.Hash[8], err))
			}
		}

		lastBlockConverted := blockInConverted
		for ReturnOneBlock(lastBlockConverted.PrevHash).Hash == "" && lastBlockConverted.Height > 1 {
			network.Uptodate = false
			hash, _ := hex.DecodeString(lastBlockConverted.PrevHash)

			var lastBlock *protocol.Block
			if lastBlock = fetchBlock(hash[:]); lastBlock == nil {
				for lastBlock == nil {
					fmt.Printf("Try to fetch block %x again\n", hash[:])
					lastBlock = fetchBlock(hash[:])
				}
			}

			lastBlockConverted = utilities.ConvertBlock(lastBlock)
			if err := SaveBlockAndTransactions(lastBlock, false); err != nil {
				log.Fatal(fmt.Sprintf("Saving block %x and its transactions failed: %v", lastBlock.Hash[8], err))
			}
		}

		network.Uptodate = true
	}
}

func fetchBlock(blockHash []byte) (block *protocol.Block) {
	var errormsg string
	if blockHash != nil {
		errormsg = fmt.Sprintf("Loading block %x failed: ", blockHash[:8])
	}

	err := network.BlockReq(blockHash[:])
	if err != nil {
		fmt.Println(errormsg + err.Error())
		return nil
	}

	blockI, err := network.Fetch(network.BlockChan)
	if err != nil {
		fmt.Println(errormsg + err.Error())
		return nil
	}

	block = blockI.(*protocol.Block)

	return block
}

func SaveBlockAndTransactions(oneBlock *protocol.Block, doPostSave bool) error {
	for _, accTxHash := range oneBlock.AccTxData {
		err := network.TxReq(p2p.ACCTX_REQ, accTxHash)
		if err != nil {
			return err
		}

		txI, err := network.Fetch(network.AccTxChan)
		if err != nil {
			return err
		}

		accTx := txI.(*protocol.AccTx)

		convertedTx := utilities.ConvertAccTransaction(accTx, oneBlock.Hash, accTxHash, oneBlock.Timestamp)
		accountHashBytes := utilities.SerializeHashContent(accTx.PubKey)
		accountHash := fmt.Sprintf("%x", accountHashBytes)

		WriteAccountWithAddress(convertedTx, accountHash)
		WriteAccTx(convertedTx)
	}

	for _, fundsTxHash := range oneBlock.FundsTxData {
		err := network.TxReq(p2p.FUNDSTX_REQ, fundsTxHash)
		if err != nil {
			return err
		}

		txI, err := network.Fetch(network.FundsTxChan)
		if err != nil {
			return err
		}

		fundsTx := txI.(*protocol.FundsTx)

		convertedTx := utilities.ConvertFundsTransaction(fundsTx, oneBlock.Hash, fundsTxHash, oneBlock.Timestamp)

		UpdateAccountDataFunds(convertedTx)
		WriteFundsTx(convertedTx)
	}

	for _, configTxHash := range oneBlock.ConfigTxData {
		err := network.TxReq(p2p.CONFIGTX_REQ, configTxHash)
		if err != nil {
			return err
		}

		txI, err := network.Fetch(network.ConfigTxChan)
		if err != nil {
			return err
		}

		configTx := txI.(*protocol.ConfigTx)

		convertedTx := utilities.ConvertConfigTransaction(configTx, oneBlock.Hash, configTxHash, oneBlock.Timestamp)
		currentParams := ReturnNewestParameters()
		newParams := utilities.ExtractParameters(convertedTx, currentParams)

		WriteParameters(newParams)
		WriteConfigTx(convertedTx)
	}

	for _, stakeTxHash := range oneBlock.StakeTxData {
		err := network.TxReq(p2p.STAKETX_REQ, stakeTxHash)
		if err != nil {
			return err
		}

		txI, err := network.Fetch(network.StakeTxChan)
		if err != nil {
			return err
		}

		stakeTx := txI.(*protocol.StakeTx)

		convertedTx := utilities.ConvertStakeTransaction(stakeTx, oneBlock.Hash, stakeTxHash, oneBlock.Timestamp)

		UpdateAccountDataStaking(convertedTx)
		WriteStakeTx(convertedTx)
	}

	convertedBlock := utilities.ConvertBlock(oneBlock)
	WriteBlock(convertedBlock)

	if doPostSave {
		postSaveBlockAndTransactions()
	}

	fmt.Printf("Block with height %v loaded and saved\n", oneBlock.Height)

	return nil
}

func postSaveBlockAndTransactions() {
	RemoveRootFromDB()
	UpdateTotals()
}
