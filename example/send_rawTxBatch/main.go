package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	// directory of the generated code using the provided relay.proto file
	pb "github.com/BlockRazorinc/relay_example/protobuf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"google.golang.org/grpc"
)

// auth will use to verify the certificate.
type Authentication struct {
	apiKey string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"apiKey": a.apiKey}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func main() {

	// BlockRazor relay endpoint address
	blzrelayEndPoint := "ip:port"

	// auth will be used to verify the credential
	auth := Authentication{
		"your auth token",
	}

	// open gRPC connection to BlockRazor relay
	var err error
	conn, err := grpc.Dial(blzrelayEndPoint, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth), grpc.WithWriteBufferSize(0), grpc.WithInitialConnWindowSize(128*1024))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	// use the Gateway client connection interface
	client := pb.NewGatewayClient(conn)

	// create context
	ctx := context.Background()

	// replace with your address
	from_private_address1 := "6c0456……8b8003"
	from_private_address2 := "42b565……44d05c"
	to_public_address := "0x4321……3f1c66"

	// replace with your transaction data
	nonce1 := uint64(1)
	nonce2 := uint64(1)
	toAddress := common.HexToAddress(to_public_address)
	var data []byte
	gasPrice := big.NewInt(1e9)
	gasLimit := uint64(22000)
	value := big.NewInt(0)
	chainid := types.NewEIP155Signer(big.NewInt(56))

	// create new transaction
	tx1 := types.NewTransaction(nonce1, toAddress, value, gasLimit, gasPrice, data)
	tx2 := types.NewTransaction(nonce2, toAddress, value, gasLimit, gasPrice, data)

	privateKey1, err := crypto.HexToECDSA(from_private_address1)
	if err != nil {
		fmt.Println("fail to casting private key to ECDSA")
		return
	}

	privateKey2, err := crypto.HexToECDSA(from_private_address2)
	if err != nil {
		fmt.Println("fail to casting private key to ECDSA")
		return
	}

	// sign transaction by private key
	signedTx1, err := types.SignTx(tx1, chainid, privateKey1)
	if err != nil {
		fmt.Println("fail to sign transaction")
		return
	}

	signedTx2, err := types.SignTx(tx2, chainid, privateKey2)
	if err != nil {
		fmt.Println("fail to sign transaction")
		return
	}

	// use rlp to encode your transaction
	body, _ := rlp.EncodeToBytes([]types.Transaction{*signedTx1, *signedTx2})

	// encode []byte to string
	encode_txs := hex.EncodeToString(body)

	// send raw tx batch by BlockRazor
	res, err := client.SendTxs(ctx, &pb.SendTxsRequest{Transactions: encode_txs})

	if err != nil {
		fmt.Println("failed to send raw tx batch: ", err)
		return
	} else {
		fmt.Println("raw tx batch sent by BlockRazor, tx hashes are ", res.TxHashs)
	}

}
