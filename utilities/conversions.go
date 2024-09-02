package utilities

import (
	"fmt"
	"github.com/way365/bazo-miner/protocol"
	"time"
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
	convertedBlock.NrUpdateTx = unconvertedBlock.NrUpdateTx
	convertedBlock.NrUpdates = unconvertedBlock.NrUpdates

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
	for _, hash := range unconvertedBlock.UpdateTxData {
		convertedTxHash = fmt.Sprintf("%x", hash)
		convertedBlock.UpdateTxData = append(convertedBlock.UpdateTxData, convertedTxHash)
	}

	return convertedBlock
}

func ConvertFundsTransaction(unconvertedTx *protocol.FundsTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Fundstx {
	var convertedTx Fundstx

	//convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
	convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
	convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
	convertedTx.Amount = unconvertedTx.Amount
	convertedTx.Fee = unconvertedTx.Fee
	convertedTx.TxCount = unconvertedTx.TxCnt
	convertedTx.From = fmt.Sprintf("%x", unconvertedTx.From)
	convertedTx.To = fmt.Sprintf("%x", unconvertedTx.To)
	convertedTx.Timestamp = blockTimestamp
	convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig1)
	convertedTx.Data = fmt.Sprintf("%s", unconvertedTx.Data)

	return convertedTx
}

func ConvertAccTransaction(unconvertedTx *protocol.AccTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Acctx {
	var convertedTx Acctx

	//convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
	convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
	convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
	convertedTx.Fee = unconvertedTx.Fee
	convertedTx.Issuer = fmt.Sprintf("%x", unconvertedTx.Issuer)
	convertedTx.PubKey = fmt.Sprintf("%x", unconvertedTx.PubKey)
	convertedTx.Timestamp = blockTimestamp
	convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)
	convertedTx.Data = fmt.Sprintf("%s", unconvertedTx.Data)

	return convertedTx
}

func ConvertConfigTransaction(unconvertedTx *protocol.ConfigTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Configtx {
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
	convertedTx.Timestamp = blockTimestamp
	convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

	return convertedTx
}

func ConvertStakeTransaction(unconvertedTx *protocol.StakeTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Staketx {
	var convertedTx Staketx

	convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
	convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
	convertedTx.Timestamp = blockTimestamp
	convertedTx.Fee = unconvertedTx.Fee
	convertedTx.Account = fmt.Sprintf("%x", unconvertedTx.Account)
	convertedTx.IsStaking = unconvertedTx.IsStaking
	convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

	return convertedTx
}

func ConvertUpdateTransaction(unconvertedTx *protocol.UpdateTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Updatetx {
	var convertedTx Updatetx

	convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
	convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
	convertedTx.Timestamp = blockTimestamp
	convertedTx.Fee = unconvertedTx.Fee
	convertedTx.ToUpdateHash = fmt.Sprintf("%x", unconvertedTx.TxToUpdateHash)
	convertedTx.ToUpdateData = fmt.Sprintf("%s", unconvertedTx.TxToUpdateData)
	convertedTx.Issuer = fmt.Sprintf("%x", unconvertedTx.Issuer)
	convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)
	convertedTx.Data = fmt.Sprintf("%s", unconvertedTx.Data)

	return convertedTx
}

func ConvertAggTransaction(unconvertedTx *protocol.AggTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte, blockTimestamp int64) Aggtx {
	var convertedTx Aggtx

	convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
	convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
	convertedTx.Amount = unconvertedTx.Amount
	convertedTx.Fee = unconvertedTx.Fee
	convertedTx.From = fmt.Sprintf("%x", unconvertedTx.From)
	convertedTx.To = fmt.Sprintf("%x", unconvertedTx.To)
	convertedTx.Timestamp = blockTimestamp
	convertedTx.MerkleRoot = fmt.Sprintf("%x", unconvertedTx.MerkleRoot)

	for _, hash := range unconvertedTx.AggregatedTxSlice {
		convertedTxHash := fmt.Sprintf("%x", hash)
		convertedTx.AggTxData = append(convertedTx.AggTxData, convertedTxHash)
	}

	return convertedTx
}
