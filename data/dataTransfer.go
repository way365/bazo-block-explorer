package data

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/way365/bazo-block-explorer/utilities"
	"github.com/way365/bazo-miner/miner"
	"github.com/way365/bazo-miner/p2p"
	"github.com/way365/bazo-miner/protocol"
	"log"
	"net"
	"time"
)

var newestBlock *protocol.Block
var logger *log.Logger
var block1 *protocol.Block

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func RunDB() {

	saveInitialParameters()

	for loadAllBlocks() == false {
		time.Sleep(time.Second * 120)
	}
	for 0 < 1 {
		time.Sleep(time.Second * 3)
		RefreshState()
	}
}

func saveInitialParameters() {
	var convertedParameters utilities.Systemparams
	parameters := miner.NewDefaultParameters()

	convertedParameters.Timestamp = time.Now().Unix()
	convertedParameters.BlockSize = parameters.BlockSize
	convertedParameters.DiffInterval = parameters.DiffInterval
	convertedParameters.MinFee = parameters.FeeMinimum
	convertedParameters.BlockInterval = parameters.BlockInterval
	convertedParameters.BlockReward = parameters.DiffInterval
	convertedParameters.StakingMin = parameters.StakingMinimum
	convertedParameters.WaitingMin = parameters.WaitingMinimum
	convertedParameters.AcceptanceTimeDiff = parameters.AcceptedTimeDiff
	convertedParameters.SlashingWindowSize = parameters.SlashingWindowSize
	convertedParameters.SlashingReward = parameters.SlashReward

	WriteParameters(convertedParameters)
}

func loadAllBlocks() bool {
	defer timeTrack(time.Now(), "Copying Database")

	//request newest block with argument nil
	block := reqBlock(nil)
	var emptyBlock *protocol.Block
	if block == emptyBlock {
		//connection to miner failed, will retry after interval in RunDB()
		return false
	}
	fmt.Printf("Copying Data...")
	//newestBlock is stored for the next iteration of RefreshState() to check if it changed
	newestBlock = block
	SaveBlockAndTransactions(block)

	//query the previous block recursively.
	for block.Hash != [32]byte{} {
		var queryHash [2 * miner.BLOCKHASH_SIZE]byte
		copy(queryHash[:32], block.PrevHash[:])
		copy(queryHash[32:], block.PrevHashWithoutTx[:])

		block = reqBlock(queryHash[:])
		SaveBlockAndTransactions(block)
	}
	//remove root account from database, since its balance makes no sense
	RemoveRootFromDB()
	UpdateTotals()
	fmt.Println("All Blocks Loaded!")
	return true
}

func RefreshState() {
	fmt.Println("Refreshing State...")

	//request newest block with argument nil
	block := reqBlock(nil)
	var emptyBlock *protocol.Block
	if block == emptyBlock {
		//connection to miner failed, will retry after interval in RunDB()
		return
	}
	prevHash := block.PrevHash
	tempBlock := block

	if block.Hash == newestBlock.Hash {
		//No new Blocks
		RemoveRootFromDB()
		UpdateTotals()
		return

	} else if prevHash == newestBlock.Hash {
		//One new Block
		SaveBlockAndTransactions(block)
		newestBlock = block

		RemoveRootFromDB()
		UpdateTotals()
		return

	} else if block.Hash != newestBlock.Hash {
		//Multiple new Blocks
		SaveBlockAndTransactions(block)
	}

	for block.PrevHash != newestBlock.Hash {
		var queryHash [2 * miner.BLOCKHASH_SIZE]byte
		copy(queryHash[:32], block.PrevHash[:])
		copy(queryHash[32:], block.PrevHashWithoutTx[:])

		block = reqBlock(queryHash[:])
		SaveBlockAndTransactions(block)
	}
	newestBlock = tempBlock

	RemoveRootFromDB()
	UpdateTotals()
}

func Connect(connectionString string) (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", connectionString)

	if err != nil {
		fmt.Println("Could not connect to a miner!")
		log.Println(err)
		return conn, err
	}
	conn.SetDeadline(time.Now().Add(20 * time.Second))

	return conn, err
}

func reqBlock(blockHash []byte) (block *protocol.Block) {
	//request data using modified code from bazo's p2p messaging system
	conn, err := Connect("127.0.0.1:8000") // Docker will resolve the correct IP for bazo-miner.
	if err != nil {
		var emptyBlock *protocol.Block
		return emptyBlock
	}
	packet := p2p.BuildPacket(p2p.BLOCK_REQ, blockHash[:])
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		return
	}
	if header.TypeID == p2p.BLOCK_RES {
		block = block.Decode(payload)
	}
	conn.Close()

	return block
}

func reqTx(txType uint8, txHash [32]byte) interface{} {
	//request data using modified code from bazo's p2p messaging system
	conn, _ := Connect("127.0.0.1:8000") // Docker will resolve the correct IP for bazo-miner.
	packet := p2p.BuildPacket(txType, txHash[:])
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		panic(err)
	}
	defer conn.Close()

	switch header.TypeID {
	case p2p.ACCTX_RES:
		var accTx *protocol.AccTx
		accTx = accTx.Decode(payload)
		return accTx
	case p2p.CONFIGTX_RES:
		var configTx *protocol.ConfigTx
		configTx = configTx.Decode(payload)
		return configTx
	case p2p.FUNDSTX_RES:
		var fundsTx *protocol.FundsTx
		fundsTx = fundsTx.Decode(payload)
		return fundsTx
	case p2p.STAKETX_RES:
		var stakeTx *protocol.StakeTx
		stakeTx = stakeTx.Decode(payload)
		return stakeTx
	case p2p.UPDATETX_RES:
		var updateTx *protocol.UpdateTx
		updateTx = updateTx.Decode(payload)
		return updateTx
	case p2p.AGGTX_RES:
		var aggTx *protocol.AggTx
		aggTx = aggTx.Decode(payload)
		return aggTx
	default:
		panic(err)
	}
}

func rcvData(c net.Conn) (header *p2p.Header, payload []byte, err error) {
	//request data using modified code from bazo's p2p messaging system
	reader := bufio.NewReader(c)
	header, err = p2p.ReadHeader(reader)

	if err != nil {
		c.Close()
		return nil, nil, errors.New(fmt.Sprintf("Connection to aborted: (%v)\n", err))
	}
	payload = make([]byte, header.Len)

	for cnt := 0; cnt < int(header.Len); cnt++ {
		payload[cnt], err = reader.ReadByte()
		if err != nil {
			c.Close()
			return nil, nil, errors.New(fmt.Sprintf("Connection to aborted: %v\n", err))
		}
	}

	return header, payload, nil
}

func SaveBlockAndTransactions(oneBlock *protocol.Block) {
	for _, accTxHash := range oneBlock.AccTxData {
		accTx := reqTx(p2p.ACCTX_REQ, accTxHash)
		convertedTx := utilities.ConvertAccTransaction(accTx.(*protocol.AccTx), oneBlock.Hash, accTxHash, oneBlock.Timestamp)
		accountHashBytes := utilities.SerializeHashContent(accTx.(*protocol.AccTx).PubKey)
		accountHash := fmt.Sprintf("%x", accountHashBytes)

		WriteAccountWithAddress(convertedTx, accountHash)
		WriteAccTx(convertedTx)
	}

	for _, fundsTxHash := range oneBlock.FundsTxData {
		fundsTx := reqTx(p2p.FUNDSTX_REQ, fundsTxHash)
		convertedTx := utilities.ConvertFundsTransaction(fundsTx.(*protocol.FundsTx), oneBlock.Hash, fundsTxHash, oneBlock.Timestamp)

		UpdateAccountData(convertedTx)
		WriteFundsTx(convertedTx)
	}

	for _, configTxHash := range oneBlock.ConfigTxData {
		configTx := reqTx(p2p.CONFIGTX_REQ, configTxHash)
		convertedTx := utilities.ConvertConfigTransaction(configTx.(*protocol.ConfigTx), oneBlock.Hash, configTxHash, oneBlock.Timestamp)
		currentParams := ReturnNewestParameters()
		newParams := utilities.ExtractParameters(convertedTx, currentParams)

		WriteParameters(newParams)
		WriteConfigTx(convertedTx)
	}

	for _, stakeTxHash := range oneBlock.StakeTxData {
		stakeTx := reqTx(p2p.STAKETX_REQ, stakeTxHash)
		convertedTx := utilities.ConvertStakeTransaction(stakeTx.(*protocol.StakeTx), oneBlock.Hash, stakeTxHash, oneBlock.Timestamp)

		UpdateAccountIsStaking(convertedTx)
		WriteStakeTx(convertedTx)
	}

	for _, updateTxHash := range oneBlock.UpdateTxData {
		updateTx := reqTx(p2p.UPDATETX_REQ, updateTxHash)
		convertedTx := utilities.ConvertUpdateTransaction(updateTx.(*protocol.UpdateTx), oneBlock.Hash, updateTxHash, oneBlock.Timestamp)

		UpdateTxToUpdate(convertedTx)
		WriteUpdateTx(convertedTx)
	}

	for _, aggTxHash := range oneBlock.AggTxData {
		aggTx := reqTx(p2p.AGGTX_REQ, aggTxHash)
		convertedTx := utilities.ConvertAggTransaction(aggTx.(*protocol.AggTx), oneBlock.Hash, aggTxHash, oneBlock.Timestamp)
		WriteAggTx(convertedTx)
	}

	convertedBlock := utilities.ConvertBlock(oneBlock)
	WriteBlock(convertedBlock)
}
