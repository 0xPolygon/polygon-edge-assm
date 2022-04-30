package main

import (
	"Trapesys/polygon-edge-assm/aws"
	"Trapesys/polygon-edge-assm/genesis"
	"Trapesys/polygon-edge-assm/types"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)



var nodes = types.Nodes{
	Total: 0,
	Finished: make([]string, 0),
	NodeIPs: make(map[string]string),
	Node: make(map[string]types.NodeInfo),
}

var (
	logger log.Logger
	logFileFlag string
)

func main() {

	flag.StringVar(&aws.Region,"aws-region","us-west-2","set AWS region")
	flag.StringVar(&aws.BucketName,"s3-name","polygon-edge-shared","set S3 bucket name")
	flag.StringVar(&logFileFlag, "log-file","/var/log/edge-assm.log","log file location")
	flag.Parse()

	logFile, err := os.Create(logFileFlag)
	if err != nil {
		log.Println("could not set log file location")
	} else {
		logger = *log.New(logFile,"edge-assm",log.Ldate|log.Ltime)
	}
	
	r := mux.NewRouter()
	// all nodes done, start generating genesis.json /init
	r.HandleFunc("/init",handleInit).Methods("GET")
	// this node has finished init phase /node-done?name=node1&ip=10.150.1.4
	r.HandleFunc("/node-done",handleDoneNode).Methods("GET")
	// get the total number of nodes /total-nodes?total=4
	r.HandleFunc("/total-nodes",handleTotalNodes).Methods("GET")

	srv := &http.Server{
		Addr: "0.0.0.0:9001",
		Handler: r,
	}

	srv.ListenAndServe()
}

func handleTotalNodes(w http.ResponseWriter, r *http.Request) {
 	total, err := strconv.Atoi(r.URL.Query()["total"][0])
	if err != nil {
		logger.Println("could not convert string to int, %w", err)
		return
	}
	nodes.Total = total
	json.NewEncoder(w).Encode(types.Responce{Success: true, Message: "total node number set!"})
}

func handleInit(w http.ResponseWriter, r *http.Request) {
	// skip if there are no nodes registered
	if nodes.Total == 0 {
		json.NewEncoder(w).Encode(types.Responce{Success: false,Message: "there are 0 nodes registered!"})
		return
	}
	// if there no nodes are finished registered skip this function
	if len(nodes.Finished) == 0 {
		json.NewEncoder(w).Encode(types.Responce{Success: false,Message: "there are 0 nodes that have finished init phase!"})
		return
	}

	// if there are less finished nodes than registered nodes skip this function
	if !(len(nodes.Finished) == nodes.Total) {
		json.NewEncoder(w).Encode(types.Responce{Success: false,Message: "the number of finished nodes and total number of nodes doesn't match"})
		return
	}

	// get the data only if all nodes have finished
	for _,name := range nodes.Finished {
		// get network-key from ASSM
		id, err := aws.GetSecret(fmt.Sprintf("/polygon-edge/nodes/%s/network-key",name))
		if err != nil {
			logger.Println("could not fetch network key secret: "+ name + err.Error())
			return
		}

		// get validator-key from ASSM
		key, err := aws.GetSecret(fmt.Sprintf("/polygon-edge/nodes/%s/validator-key",name))
		if err != nil {
			logger.Println("coult not fetch validator key secret: ", name + err.Error())
			return
		}

		// get new node info based on private keys 
		nodeInfo, err := types.NewNodeInfo(id, key, nodes.NodeIPs[name])
		if err != nil {
			logger.Println("could not set validator and network params: %w", err)
			return
		}

		// set node info
		nodes.Node[name] = *nodeInfo
	}
	
	if err := genesis.GenerateAndStore(&nodes); err != nil {
		log.Println("genesis genrator failed: ", err)
		return
	}

	json.NewEncoder(w).Encode(types.Responce{Success: true, Message: "genesis.json file generated and stored to S3 bucket"})

	// after generating genesis.json reset this variable
	nodes = types.Nodes{
		Total: 0,
		Finished: make([]string, 0),
		NodeIPs: make(map[string]string),
		Node: make(map[string]types.NodeInfo),
	}
}

func handleDoneNode(w http.ResponseWriter, r *http.Request){
	for _, n := range nodes.Finished {
		// if we already have this node name, don't run this function
		if n == r.URL.Query()["name"][0] {
			return
		}
	}

	nodeName := r.URL.Query()["name"][0]
	nodeIP := r.URL.Query()["ip"][0]
 	nodes.Finished = append(nodes.Finished,nodeName)
	nodes.NodeIPs[nodeName] = nodeIP

	json.NewEncoder(w).Encode(types.Responce{Success: true, Message: "node registered"})
}

