# Polygon Edge secrets manager initializer

On running automated Polygon Edge deployments this API helps with creation of `genesis.json` file that the Polygon Edge server needs to run a chain.
When the private keys are saved in some secrets manager solution ( right now only AWS SSM is supported ) this API fetches these keys and converts them to network_id and validator address.

### Prequestites
* Dedicated node that will clone this repo, compile binary and run it.
* This node needs to have access to both ASSM and S3 that will hold `genesis.json` file ( instance IAM polices ).  
* All polygon edge nodes need to be able to access this node at TCP 9001 by default ( security groups ).

## How to use
The genesis creation process involves several stages, consisting of hitting the API and delivering the required data to it.  

### Total number of nodes API
If there are 4 validator nodes in total: `/total-nodes?total=4`  
Now the API knows that there are 4 nodes that we would like to initialize as validator nodes.

### Initialization done
Each validator node needs to send the following api when it finishes the `secrets init` stage: ` /node-done?name=node1&ip=10.150.1.4`   
Node sends its name, which coresponds to the name in secrets manager, and its IP address.  
Once the program receives enough calls to this API ( for 4 nodes, the program expects 4 calls to this API ), it moves to the next stage.   

### Fetch keys and generate genesis.json file
Once all validator nodes reported that they have successfully completed the `secrets init` stage, this program fetches validator secrets from the secrets store, generates `genesis.json` file and puts it in the S3 bucket.   
The API call that triggers this action is: `/init`   

Each node can be configured to send all 3 API calls.   
Once the last node hits `/init` api, the `genesis.json` file generation will start.

### Configuration options
Flags that can be set for this program are:   
* `aws-region` - sets the AWS region for the SSM. Default: `us-west-2`
* `s3-name` - sets S3 bucket name in which to place the `genesis.json` file. Default: `polygon-edge-shared`
* `log-file` - sets the log file output. Default: `/var/log/edge-assm.log`
* `genesis-log-file` - sets log file output for genesis module. Default: `/var/log/edge-assm-genesis.log`
* `chain-name` - sets chain name. Default: pulled from `polygon-edge genesis` command
* `pos` - sets PoS consensus. Default: false
* `epoch-size` - sets epoch size. Default: pulled from `polygon-edge genesis` command
* `premine` - premine accounts. For multiple accounts, separate them with `,`. Format: `<account>:<ammount>`
* `chain-id` - sets chain id. Default: pulled from `polygon-edge genesis` command
* `block-gas-limit` - sets block gas limit. Default: pulled from `polygon-edge genesis` command
* `max-validator-count` - sets maximum validator count, only for PoS consensus. Default: pulled from `polygon-edge genesis` command
* `min-validator-count` - sets minimum validator count, only for PoS consensus. Default: pulled from `polygon-edge genesis` command


