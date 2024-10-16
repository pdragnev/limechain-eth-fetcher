package models

type Transaction struct {
	ID                string  `db:"id" json:"-"`
	TransactionHash   string  `db:"transaction_hash" json:"transactionHash"`
	TransactionStatus uint64  `db:"transaction_status" json:"transactionStatus"`
	BlockHash         string  `db:"block_hash" json:"blockHash"`
	BlockNumber       uint64  `db:"block_number" json:"blockNumber"`
	From              string  `db:"from_address" json:"from"`
	To                *string `db:"to_address,omitempty" json:"to,omitempty"`
	ContractAddress   *string `db:"contract_address,omitempty" json:"contractAddress,omitempty"`
	LogsCount         int     `db:"logs_count" json:"logsCount"`
	Input             string  `db:"input" json:"input"`
	Value             string  `db:"value" json:"value"`
}

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}
