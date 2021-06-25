package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"net/http"
	"os"
	"strings"
)

var (
	port        string
	name        string
	eth         *ethclient.Client
)


func ConnectionToGeth(url string) error {
	var err error
	eth, err = ethclient.Dial(url)
	return err
}


func CurrentBlock() (uint64, uint64) {
	block, err := eth.BlockByNumber(context.TODO(), nil)
	if err != nil {
		fmt.Printf("Error fetching current block height: %v\n", err)
		return 0, 0
	}

	return block.NumberU64(), block.Time()
}

func ToEther(o *big.Int) *big.Float {
	pul, int := big.NewFloat(0), big.NewFloat(0)
	int.SetInt(o)
	pul.Mul(big.NewFloat(0.000000000000000001), int)
	return pul
}

func MetricsHttp(w http.ResponseWriter, r *http.Request) {
	var allOut []string
	block_number, timestamp := CurrentBlock()
	allOut = append(allOut, fmt.Sprintf("geth_block_number{geth_name=\"%v\"} %v", name, block_number))
	allOut = append(allOut, fmt.Sprintf("geth_block_timestamp{geth_name=\"%v\"} %v", name, timestamp))
	fmt.Fprintln(w, strings.Join(allOut, "\n"))
	// https://github.com/ethereum/go-ethereum/blob/master/core/types/block.go
}

func main() {
	gethRPC := os.Getenv("GETH_RPC")
	port = os.Getenv("PORT")
	name = os.Getenv("NAME")

	err := ConnectionToGeth(gethRPC)
	if err != nil {
		panic(err)
	}

	block, _ := CurrentBlock()

	fmt.Printf("nGeth exporter listen on port %v \nGeth server: %v \nStart at block #%v\n", port, gethRPC, block)
	http.HandleFunc("/metrics", MetricsHttp)
	panic(http.ListenAndServe("0.0.0.0:"+port, nil))
}
