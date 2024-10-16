package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"my-lime/api"
	"my-lime/internal/config"
	"my-lime/internal/repository"
	"my-lime/internal/service"

	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}
	db, err := repository.NewPostgresDB(cfg.DB_CONNECTION_URL)
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	redisClient := repository.NewRedisClient(cfg.REDIS_URL)

	ethClient, err := service.NewEthereumClient(cfg.ETH_NODE_URL)
	if err != nil {
		log.Fatalf("Failed to initialize Ethereum client: %s", err.Error())
	}

	rlpData := "f90110b842307866633262336236646233386135316462336239636239356465323962373139646538646562393936333036323665346234623939646630353666666237663265b842307834383630336637616466663766626663326131306232326136373130333331656536386632653464316364373361353834643537633838323164663739333536b842307863626339323065376262383963626362353430613436396131363232366266313035373832353238336162386561633366343564303038313165656638613634b842307836643630346666633634346132383266636138636238653737386531653366383234356438626431643439333236653330313661336338373862613063626264"

	// Decode the hex string to bytes
	rlpBytes, err := hex.DecodeString(rlpData)
	if err != nil {
		log.Fatalf("Failed to decode hex string: %v", err)
	}

	// Define a variable to hold the decoded data
	var decodedData []interface{}

	// Use RLP decoding from go-ethereum package
	err = rlp.Decode(bytes.NewReader(rlpBytes), &decodedData)
	if err != nil {
		log.Fatalf("Failed to decode RLP data: %v", err)
	}

	// Print the decoded data
	for i, data := range decodedData {
		fmt.Printf("Element %d: %s\n", i, data)
	}

	services := service.NewService(repos, redisClient, cfg, ethClient)
	api := api.NewApi(services)

	r := api.NewRouter(services)
	api.StartServer(r, cfg.API_PORT)
}
