package data

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/way365/bazo-block-explorer/utilities"
	"strings"
	"time"
)

const (
	host   = "127.0.0.1" // Docker will resolve the correct IP for psql.
	port   = 5432
	dbname = "blockexplorerdb"
)

var sqlStatement string
var db *sql.DB
var err error

var name string
var password string

func connectToDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, name, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func SetupDB(username string, userpassword string) {
	//arguments for database login gets stored in memory
	name = username
	password = userpassword

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, name, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Setting up Database...")
	dropTables()
	createTables()
	fmt.Println("Setup Complete")
}

func dropTables() {
	connectToDB()
	defer db.Close()
	fmt.Println("Dropping Tables...")
	sqlStatement := `drop table if exists blocks;
                   drop table if exists fundstx;
                   drop table if exists acctx;
                   drop table if exists configtx;
                   drop table if exists staketx;
                   drop table if exists updatetx;
                   drop table if exists aggtx;
                   drop table if exists accounts;
                   drop table if exists parameters;
                   drop table if exists stats;
                   drop table if exists txhistory;`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	fmt.Println("Dropped Tables")

}

func createTables() {
	connectToDB()
	defer db.Close()
	fmt.Println("Creating Tables...")

	sqlStatement1 := `create table blocks (
                    header bit(8),
                    hash char(64) primary key,
                    prevHash char(64) not null,
                    nonce char(16),
                    timestamp bigint not null,
                    timestring varchar(100) not null,
                    merkleRoot char(64) not null,
                    beneficiary char(64) not null,
                    seed varchar(100),
                    hashedseed varchar(100),
                    height smallint,
                    nrFundsTx smallint not null,
                    nrAccTx smallint not null,
                    nrConfigTx smallint not null,
                    nrStakeTx smallint,
                    nrUpdateTx smallint,
                    fundsTxData varchar(100)[],
                    accTxData varchar(100)[],
                    configTxData varchar(100)[],
                    stakeTxData varchar(100)[],
                    updateTxData varchar(100)[],
                    nrUpdates smallint not null
                    );

                    create table fundstx (
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64) not null,
                    amount bigint not null,
                    fee bigint not null,
                    txcount int not null,
                    sender char(64) not null,
                    recipient char(64) not null,
                    timestamp bigint not null,
                    signature char(128) not null,
                    data varchar(512) default null
                    );
                    
                    create table updatetx (
                    hash char(64) primary key,
                    blockhash char(64) not null,
                    fee bigint not null,
                    toUpdateHash char(64),
                    toUpdateData varchar(512) default null,
                    issuer char(64) not null,
                    timestamp bigint not null,
                    signature char(128) not null,
                    data varchar(512) default null
                    );
                    
                    create table aggtx (
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64) not null,
                    amount bigint,
                    fee bigint,
                    sender char(128) not null,
                    recipient char(128) not null,
                    merkleRoot char(64) not null,
                    aggTxData varchar(100)[],
                    timestamp bigint not null
                    );`

	sqlStatement2 := `create table acctx(
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64),
                    issuer char(64) not null,
                    fee bigint not null,
                    pubkey char(128) not null,
                    timestamp bigint not null,
                    signature char(128) not null,
                    data varchar(512) default null
                    );

                    create table configtx(
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64),
                    id int not null,
                    payload bigint not null,
                    fee bigint not null,
                    txcount int not null,
                    timestamp bigint not null,
                    signature char(128) not null
                    );

                    create table staketx(
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64),
                    timestamp bigint not null,
                    fee bigint not null,
                    account char(64) not null,
                    isstaking boolean,
                    signature char(128) not null
                    );`

	sqlStatement3 := `create table accounts(
                    hash char(64) primary key,
                    address char(128),
                    balance bigint not null,
                    txcount int not null,
                    isstaking boolean
                    );

                    create table parameters(
                    timestamp bigint not null,
                    blocksize int not null,
                    diffinterval int not null,
                    minfee int not null,
                    blockinterval int not null,
                    blockreward int not null,
                    stakingmin int,
                    waitingmin int,
                    accepancetimediff int,
                    slashingwindowsize int,
                    slashingreward int
                    );

                    create table stats(
                    totalsupply bigint,
                    nraccounts bigint,
                    timestamp bigint
                    );

                    create table txhistory(
                    timestring varchar(100) not null,
                    timestamp bigint not null,
                    txnumber bigint DEFAULT 0
                    );`

	sqlStatement4 := `create index timestamp_idx on blocks (timestamp);
                      create index sender_idx on fundstx (sender);
                      create index recipient_idx on fundstx (recipient);
                      create index address_idx on accounts (address);`

	_, err := db.Exec(sqlStatement1)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(sqlStatement2)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(sqlStatement3)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(sqlStatement4)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created Tables Successfully")
}

func ReturnOneBlock(UrlHash string) utilities.Block {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, prevHash, timestamp, timestring, merkleRoot, beneficiary, nrFundsTx, nrAccTx, nrConfigTx, nrUpdateTx, fundsTxData, accTxData, configTxData, updateTxData FROM blocks WHERE hash = $1;`

	var returnedblock utilities.Block
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err := row.Scan(
		&returnedblock.Hash,
		&returnedblock.PrevHash,
		&returnedblock.Timestamp,
		&returnedblock.TimeString,
		&returnedblock.MerkleRoot,
		&returnedblock.Beneficiary,
		&returnedblock.NrFundsTx,
		&returnedblock.NrAccTx,
		&returnedblock.NrConfigTx,
		&returnedblock.NrUpdateTx,
		&returnedblock.FundsTxDataString,
		&returnedblock.AccTxDataString,
		&returnedblock.ConfigTxDataString,
		&returnedblock.UpdateTxDataString,
	); err {
	case sql.ErrNoRows:
	case nil:
		//convert tx-datastring into slice for all types of transactions
		if len(returnedblock.FundsTxDataString.String) > 0 {
			returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString.String[1:len(returnedblock.FundsTxDataString.String)-1], ",")
		}
		if len(returnedblock.AccTxDataString.String) > 0 {
			returnedblock.AccTxData = strings.Split(returnedblock.AccTxDataString.String[1:len(returnedblock.AccTxDataString.String)-1], ",")
		}
		if len(returnedblock.ConfigTxDataString.String) > 0 {
			returnedblock.ConfigTxData = strings.Split(returnedblock.ConfigTxDataString.String[1:len(returnedblock.ConfigTxDataString.String)-1], ",")
		}
		if len(returnedblock.UpdateTxDataString.String) > 0 {
			returnedblock.UpdateTxData = strings.Split(returnedblock.UpdateTxDataString.String[1:len(returnedblock.UpdateTxDataString.String)-1], ",")
		}
		return returnedblock
	default:
		panic(err)
	}

	return returnedblock
}

func ReturnAllBlocks(UrlHash string) []utilities.Block {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, timestamp, timestring, beneficiary, nrFundsTx, nrAccTx, nrConfigTx, nrUpdateTx FROM blocks ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Block, 0)
	for rows.Next() {
		var returnedrow utilities.Block
		err = rows.Scan(
			&returnedrow.Hash,
			&returnedrow.Timestamp,
			&returnedrow.TimeString,
			&returnedrow.Beneficiary,
			&returnedrow.NrFundsTx,
			&returnedrow.NrAccTx,
			&returnedrow.NrConfigTx,
			&returnedrow.NrUpdateTx,
		)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneFundsTx(UrlHash string) utilities.Fundstx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, amount, fee, txcount, sender, recipient, signature, data FROM fundstx WHERE hash = $1;`
	var returnedrow utilities.Fundstx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(
		&returnedrow.Hash,
		&returnedrow.BlockHash,
		&returnedrow.Amount,
		&returnedrow.Fee,
		&returnedrow.TxCount,
		&returnedrow.From,
		&returnedrow.To,
		&returnedrow.Signature,
		&returnedrow.Data,
	); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}
	return returnedrow
}

func ReturnAllFundsTx(UrlHash string) []utilities.Fundstx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, amount, fee, txcount, sender, recipient, signature, data FROM fundstx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Fundstx, 0)
	for rows.Next() {
		var returnedrow utilities.Fundstx
		err = rows.Scan(
			&returnedrow.Hash,
			&returnedrow.Amount,
			&returnedrow.Fee,
			&returnedrow.TxCount,
			&returnedrow.From,
			&returnedrow.To,
			&returnedrow.Signature,
			&returnedrow.Data,
		)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneAccTx(UrlHash string) utilities.Acctx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, issuer, fee, pubkey, signature, data FROM acctx WHERE hash = $1;`
	var returnedrow utilities.Acctx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(
		&returnedrow.Hash,
		&returnedrow.BlockHash,
		&returnedrow.Issuer,
		&returnedrow.Fee,
		&returnedrow.PubKey,
		&returnedrow.Signature,
		&returnedrow.Data,
	); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func ReturnAllAccTx(UrlHash string) []utilities.Acctx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, issuer, fee, pubkey, data FROM acctx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Acctx, 0)
	for rows.Next() {
		var returnedrow utilities.Acctx
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Issuer, &returnedrow.Fee, &returnedrow.PubKey, &returnedrow.Data)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneUpdateTx(UrlHash string) utilities.Updatetx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, issuer, toupdatehash, toupdatedata, signature, data FROM updatetx WHERE hash = $1;`
	var returnedrow utilities.Updatetx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(
		&returnedrow.Hash,
		&returnedrow.BlockHash,
		&returnedrow.Issuer,
		&returnedrow.ToUpdateHash,
		&returnedrow.ToUpdateData,
		&returnedrow.Signature,
		&returnedrow.Data,
	); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func ReturnAllAggTx(UrlHash string) []utilities.Aggtx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, sender, recipient, amount FROM aggtx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Aggtx, 0)
	for rows.Next() {
		var returnedrow utilities.Aggtx
		err = rows.Scan(&returnedrow.Hash, &returnedrow.From, &returnedrow.To, &returnedrow.Amount)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneAggTx(UrlHash string) utilities.Aggtx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, amount, fee, sender, recipient, merkleRoot, aggTxData, timestamp FROM aggtx WHERE hash = $1;`
	var returnedrow utilities.Aggtx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(
		&returnedrow.Hash,
		&returnedrow.BlockHash,
		&returnedrow.Amount,
		&returnedrow.Fee,
		&returnedrow.From,
		&returnedrow.To,
		&returnedrow.MerkleRoot,
		&returnedrow.AggTxDataString,
		&returnedrow.Timestamp,
	); err {
	case sql.ErrNoRows:
	case nil:
		//convert tx-datastring into slice for all types of transactions
		if len(returnedrow.AggTxDataString.String) > 0 {
			returnedrow.AggTxData = strings.Split(returnedrow.AggTxDataString.String[1:len(returnedrow.AggTxDataString.String)-1], ",")
		}
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func ReturnAllUpdateTx(UrlHash string) []utilities.Updatetx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, issuer, toupdatehash, toupdatedata, data FROM updatetx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Updatetx, 0)
	for rows.Next() {
		var returnedrow utilities.Updatetx
		err = rows.Scan(
			&returnedrow.Hash,
			&returnedrow.Issuer,
			&returnedrow.ToUpdateHash,
			&returnedrow.ToUpdateData,
			&returnedrow.Data,
		)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneConfigTx(UrlHash string) utilities.Configtx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, id, payload, fee, txcount, signature FROM configtx WHERE hash = $1;`
	var returnedrow utilities.Configtx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.Signature); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func ReturnAllConfigTx(UrlHash string) []utilities.Configtx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, id, payload, fee, txcount FROM configtx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Configtx, 0)
	for rows.Next() {
		var returnedrow utilities.Configtx
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnOneStakeTx(UrlHash string) utilities.Staketx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, blockhash, fee, account, isstaking, signature FROM staketx WHERE hash = $1;`
	var returnedrow utilities.Staketx
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Fee, &returnedrow.Account, &returnedrow.IsStaking, &returnedrow.Signature); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func ReturnAllStakeTx(UrlHash string) []utilities.Staketx {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, account, isstaking, fee FROM staketx ORDER BY timestamp DESC LIMIT 100`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Staketx, 0)
	for rows.Next() {
		var returnedrow utilities.Staketx
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Account, &returnedrow.IsStaking, &returnedrow.Fee)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func ReturnBlocksAndTransactions(UrlHash string) utilities.Blocksandtx {
	var returnedBlocksAndTxs utilities.Blocksandtx
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, timestamp, timestring, beneficiary, nrFundsTx, nrAccTx, nrConfigTx, nrUpdateTx, nrUpdates FROM blocks ORDER BY timestamp DESC LIMIT 6`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedblocks := make([]utilities.Block, 0)
	for rows.Next() {
		var returnedrow utilities.Block
		err = rows.Scan(
			&returnedrow.Hash,
			&returnedrow.Timestamp,
			&returnedrow.TimeString,
			&returnedrow.Beneficiary,
			&returnedrow.NrFundsTx,
			&returnedrow.NrAccTx,
			&returnedrow.NrConfigTx,
			&returnedrow.NrUpdateTx,
			&returnedrow.NrUpdates,
		)
		//returnedrow.Timestamp = returnedrow.Timestamp[:19]
		if err != nil {
			panic(err)
		}
		returnedblocks = append(returnedblocks, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx ORDER BY timestamp DESC LIMIT 6`
	rows, err = db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedtxs := make([]utilities.Fundstx, 0)
	for rows.Next() {
		var returnedrow utilities.Fundstx
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
		if err != nil {
			panic(err)
		}
		returnedtxs = append(returnedtxs, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	returnedBlocksAndTxs.Blocks = returnedblocks
	returnedBlocksAndTxs.Txs = returnedtxs

	return returnedBlocksAndTxs
}

func WriteBlock(block utilities.Block) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO blocks 
    (hash, prevhash, timestamp, timestring, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, 
     nrupdatetx, fundstxdata, acctxdata, configtxdata, updatetxdata, nrupdates)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err = db.Exec(sqlStatement,
		block.Hash,
		block.PrevHash,
		block.Timestamp,
		block.TimeString,
		block.MerkleRoot,
		block.Beneficiary,
		block.NrFundsTx,
		block.NrAccTx,
		block.NrConfigTx,
		block.NrUpdateTx,
		pq.Array(block.FundsTxData),
		pq.Array(block.AccTxData),
		pq.Array(block.ConfigTxData),
		pq.Array(block.UpdateTxData),
		block.NrUpdates,
	)
	if err != nil {
		panic(err)
	}
}

func WriteFundsTx(tx utilities.Fundstx) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO fundstx (hash, blockhash, amount, fee, txcount, sender, recipient, timestamp, signature, data)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = db.Exec(sqlStatement,
		tx.Hash,
		tx.BlockHash,
		tx.Amount,
		tx.Fee,
		tx.TxCount,
		tx.From,
		tx.To,
		tx.Timestamp,
		tx.Signature,
		tx.Data,
	)
	if err != nil {
		panic(err)
	}
}

func WriteAccTx(tx utilities.Acctx) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO acctx (hash, blockhash, fee, issuer, pubkey, timestamp, signature, data)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(sqlStatement,
		tx.Hash,
		tx.BlockHash,
		tx.Fee,
		tx.Issuer,
		tx.PubKey,
		tx.Timestamp,
		tx.Signature,
		tx.Data,
	)
	if err != nil {
		panic(err)
	}
}

func WriteConfigTx(tx utilities.Configtx) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO configtx (hash, blockhash, id, payload, fee, txcount, timestamp, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Id, tx.Payload, tx.Fee, tx.TxCount, tx.Timestamp, tx.Signature)
	if err != nil {
		panic(err)
	}
}

func WriteStakeTx(tx utilities.Staketx) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO staketx (hash, blockhash, timestamp, fee, account, isstaking, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Timestamp, tx.Fee, tx.Account, tx.IsStaking, tx.Signature)
	if err != nil {
		panic(err)
	}
}

func WriteUpdateTx(tx utilities.Updatetx) {
	connectToDB()
	defer db.Close()

	sqlStatement = `
    INSERT INTO updatetx (hash, blockhash, fee, toUpdateHash, toUpdateData, issuer, timestamp, signature, data)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.Exec(sqlStatement,
		tx.Hash,
		tx.BlockHash,
		tx.Fee,
		tx.ToUpdateHash,
		tx.ToUpdateData,
		tx.Issuer,
		tx.Timestamp,
		tx.Signature,
		tx.Data,
	)
	if err != nil {
		panic(err)
	}
}

func WriteAggTx(tx utilities.Aggtx) {
	connectToDB()
	defer db.Close()

	//sqlStatement = `
	//INSERT INTO aggtx (hash, blockhash, amount, fee, sender, recipient, timestamp)
	//VALUES ($1, $2, $3, $4, $5, $6, $7)`
	//
	sqlStatement = `
    INSERT INTO aggtx (hash, blockhash, amount, fee, sender, recipient, merkleRoot, aggTxData, timestamp)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = db.Exec(sqlStatement,
		tx.Hash,
		tx.BlockHash,
		tx.Amount,
		tx.Fee,
		tx.From,
		tx.To,
		tx.MerkleRoot,
		pq.Array(tx.AggTxData),
		tx.Timestamp,
	)
	if err != nil {
		panic(err)
	}
}

func UpdateTxToUpdate(tx utilities.Updatetx) {
	connectToDB()
	defer db.Close()

	sqlStatement := `UPDATE fundstx SET data = $1 WHERE hash = $2;`
	_, err = db.Exec(sqlStatement, tx.ToUpdateData, tx.ToUpdateHash)
	if err != nil {
		panic(err)
	}

	sqlStatement = `UPDATE acctx SET data = $1 WHERE hash = $2;`
	_, err = db.Exec(sqlStatement, tx.ToUpdateData, tx.ToUpdateHash)
	if err != nil {
		panic(err)
	}

	sqlStatement = `UPDATE updatetx SET data = $1 WHERE hash = $2;`
	_, err = db.Exec(sqlStatement, tx.ToUpdateData, tx.ToUpdateHash)
	if err != nil {
		panic(err)
	}
}

func checkEmptyDB() bool {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT CASE WHEN EXISTS (SELECT * FROM blocks LIMIT 1) THEN 1 ELSE 0 END`
	var notEmpty bool
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&notEmpty); err {
	case sql.ErrNoRows:
		fmt.Printf("No rows returned!")
	case nil:
		return notEmpty
	default:
		panic(err)
	}
	return true
}

func UpdateAccountData(tx utilities.Fundstx) {
	connectToDB()
	defer db.Close()

	sqlStatement := `INSERT INTO accounts (hash, balance, txcount)
                    VALUES ($1, $2 * -1, $3)
                    ON CONFLICT (hash) DO UPDATE SET balance = accounts.balance - $2, txcount = accounts.txcount + 1 WHERE accounts.hash = $1`

	totalAmount := tx.Amount + tx.Fee
	_, err = db.Exec(sqlStatement, tx.From, totalAmount, 1)
	if err != nil {
		panic(err)
	}

	sqlStatement = `INSERT INTO accounts (hash, balance, txcount)
                    VALUES ($1, $2, $3)
                    ON CONFLICT (hash) DO UPDATE SET balance = accounts.balance + $2 WHERE accounts.hash = $1`

	_, err = db.Exec(sqlStatement, tx.To, tx.Amount, 0)
	if err != nil {
		panic(err)
	}
}

func WriteAccountWithAddress(tx utilities.Acctx, accountHash string) {
	connectToDB()
	defer db.Close()

	sqlStatement := `INSERT INTO accounts (hash, address, balance, txcount)
                    VALUES ($1, $2, $3, $4)
                    ON CONFLICT (hash) DO UPDATE SET address = $2`

	_, err = db.Exec(sqlStatement, accountHash, tx.PubKey, 0, 0)
	if err != nil {
		panic(err)
	}
}

func UpdateAccountIsStaking(tx utilities.Staketx) {
	return
}

func ReturnOneAccount(UrlHash string) utilities.Accountwithtxs {
	var returnedData utilities.Accountwithtxs

	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, address, balance, txcount FROM accounts WHERE hash = $1 OR address = $1`
	var returnedaccount utilities.Account
	row := db.QueryRow(sqlStatement, UrlHash)
	switch err = row.Scan(&returnedaccount.Hash, &returnedaccount.Address, &returnedaccount.Balance, &returnedaccount.TxCount); err {
	case sql.ErrNoRows:
		returnedData.Account = returnedaccount
		return returnedData
	case nil:
		sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient, data FROM fundstx WHERE sender = $1 OR recipient = $1`
		rows, err := db.Query(sqlStatement, returnedaccount.Hash)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		returnedrows := make([]utilities.Fundstx, 0)
		for rows.Next() {
			var returnedrow utilities.Fundstx
			err = rows.Scan(
				&returnedrow.Hash,
				&returnedrow.Amount,
				&returnedrow.Fee,
				&returnedrow.TxCount,
				&returnedrow.From,
				&returnedrow.To,
				&returnedrow.Data,
			)
			if err != nil {
				panic(err)
			}
			returnedrows = append(returnedrows, returnedrow)
		}
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		returnedData.Account = returnedaccount
		returnedData.Txs = returnedrows
		return returnedData
	default:
		panic(err)
	}
	return returnedData
}

func ReturnTopAccounts(UrlHash string) []utilities.Account {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, address, balance, txcount FROM accounts ORDER BY balance DESC LIMIT 20`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Account, 0)
	for rows.Next() {
		var returnedrow utilities.Account
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Address, &returnedrow.Balance, &returnedrow.TxCount)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return returnedrows
}

func WriteParameters(parameters utilities.Systemparams) {
	connectToDB()
	defer db.Close()

	sqlStatement := `INSERT INTO parameters (blocksize, diffinterval, minfee, blockinterval, blockreward, stakingmin, waitingmin, accepancetimediff, slashingwindowsize, slashingreward, timestamp)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = db.Exec(sqlStatement, parameters.BlockSize, parameters.DiffInterval, parameters.MinFee, parameters.BlockInterval, parameters.BlockReward, parameters.StakingMin, parameters.WaitingMin, parameters.AcceptanceTimeDiff, parameters.SlashingWindowSize, parameters.SlashingReward, parameters.Timestamp)
	if err != nil {
		panic(err)
	}
}

func ReturnNewestParameters() utilities.Systemparams {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT blocksize, diffinterval, minfee, blockinterval, blockreward, stakingmin, waitingmin, accepancetimediff, slashingwindowsize, slashingreward, timestamp FROM parameters ORDER BY timestamp DESC LIMIT 1`
	var returnedrow utilities.Systemparams
	row := db.QueryRow(sqlStatement)
	switch err = row.Scan(&returnedrow.BlockSize, &returnedrow.DiffInterval, &returnedrow.MinFee, &returnedrow.BlockInterval, &returnedrow.BlockReward, &returnedrow.StakingMin, &returnedrow.WaitingMin, &returnedrow.AcceptanceTimeDiff, &returnedrow.SlashingWindowSize, &returnedrow.SlashingReward, &returnedrow.Timestamp); err {
	case sql.ErrNoRows:
	case nil:
		return returnedrow
	default:
		panic(err)
	}

	return returnedrow
}

func RemoveRootFromDB() {
	connectToDB()
	defer db.Close()

	sqlStatement := `DELETE FROM accounts WHERE address IS NULL;`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

func UpdateTotals() {
	connectToDB()
	defer db.Close()
	sqlStatement = `SELECT SUM(balance) FROM accounts`

	var totalSupply sql.NullInt64
	row := db.QueryRow(sqlStatement)
	switch err = row.Scan(&totalSupply); err {
	case sql.ErrNoRows:
	case nil:
	default:
		panic(err)
	}
	if totalSupply.Valid == false {
		return
	}

	sqlStatement = `SELECT COUNT(hash) FROM accounts`

	var totalAccounts sql.NullInt64

	row = db.QueryRow(sqlStatement)
	switch err = row.Scan(&totalAccounts); err {
	case sql.ErrNoRows:
	case nil:
	default:
		panic(err)
	}
	if totalAccounts.Valid == false {
		return
	}
	//totalaccounts = accountsTemp

	sqlStatement = `INSERT INTO stats (totalsupply, nraccounts, timestamp) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, totalSupply, totalAccounts, time.Now().Unix())
	if err != nil {
		panic(err)
	}
}

func ReturnTotals() utilities.Stats {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT totalsupply, nraccounts FROM stats ORDER BY timestamp DESC LIMIT 1`

	var stats utilities.Stats
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&stats.TotalSupply, &stats.TotalNrAccounts); err {
	case sql.ErrNoRows:
	case nil:
		return stats
	default:
		panic(err)
	}
	return stats
}

func Return14Hours() []utilities.Serie {
	connectToDB()
	defer db.Close()

	sqlStatement := `SELECT hash, timestamp, nrfundstx, nracctx, nrconfigtx FROM blocks WHERE timestamp > $1;`
	twelveHoursAgo := time.Now().Unix() - 46800

	rows, err := db.Query(sqlStatement, twelveHoursAgo)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	returnedrows := make([]utilities.Block, 0)
	for rows.Next() {
		var returnedrow utilities.Block
		err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.NrFundsTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
		if err != nil {
			panic(err)
		}
		returnedrows = append(returnedrows, returnedrow)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	if len(returnedrows) == 0 {
		var emptySeries []utilities.Serie
		return emptySeries
	}

	var series []utilities.Serie
	currentHourTxs := 0
	timeThreshold := time.Unix(returnedrows[0].Timestamp, 0).Add(time.Duration(-1) * time.Hour)

	for _, block := range returnedrows {
		if block.Timestamp > timeThreshold.Unix() {
			currentHourTxs = currentHourTxs + int(block.NrFundsTx) + int(block.NrAccTx) + int(block.NrConfigTx)
		} else {
			series = append(series, utilities.Serie{timeThreshold.Format("15:04"), currentHourTxs})

			timeThreshold = timeThreshold.Add(time.Duration(-3600) * time.Second)
			currentHourTxs = int(block.NrFundsTx) + int(block.NrAccTx) + int(block.NrConfigTx)
		}
	}
	series = append(series, utilities.Serie{timeThreshold.Format("15:04"), currentHourTxs})

	return series
}
