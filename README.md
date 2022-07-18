# AWS Lambda - Polygon Edge chain initializer

On running automated Polygon Edge deployments, this Lambda function helps with creation of `genesis.json` file 
that Polygon Edge server needs to run a chain.   
When the private keys are saved in AWS SSM this Lambda function fetches these keys 
and converts them to valid network_id and validator address.

### Prequestites
* EC2 instances must have permission to run AWS Lambda functions   
* EC2 instances must be able to read from S3 bucket
* EC2 instances have a service that starts the chain once `genesis.json` is found in S3

## How to use
Once basic init stage is complete and keys are stored in AWS SSM, a node will run this Lambda function
and send the data about the chain and itself, in json format.

When all nodes send their information, Lambda will create `genesis.json` file and store it in the S3 bucket.

All nodes should have a service running in the background that check if `genesis.json` exists in S3,
and if it does it will be downloaded and used to start the chain automatically.

Deploying GO on AWS Lambda [doc](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)

### JSON file format
```json
{
  "config": {
    "aws_region": "<AWS_REGION>", 
    "s3_bucket_name": "<S3_BUCKET_NAME>", 
    "s3_key_name": "<S3_CONFIG_KEY>",
	  
    "premine": "<PREMINE_ACCOUNT>:<PREMINE_AMOUNT>",
    "chain_name": "<CHAIN_NAME>",
    "chain_id": "<CHAIN_ID>",
    "pos": "<POS>",
    "epoch_size": "<EPOCH_SIZE>",
    "block_gas_limit": "<BLOCK_GAS_LIMIT>",
    "max_validator_count": "<MAX_VALIDATOR_COUNT>",
    "min_validator_count": "<MIN_VALIDATOR_COUNT>",
    "consensus": "<CONSENSUS>"
    
  },
  "nodes": {
	"total": <TOTAL_NUMBER_OF_NODES>,
	"node_info": {
	  "ssm_param_id": "<SSM_PREFIX>",
      "node_name": "<NODE_NAME>",
      "ip": "<NODE_IP>"
	}
  }
}
```

### General configuration
* `AWS_REGION` - AWS region for your resources
* `PREMINE_ACCOUNT:PREMINE_AMOUNT` - account and the amount that will receive defined amount of native currency 
* `S3_BUCKET_NAME` - the name of S3 bucket that will hold configuration and `genesis.json` file
* `S3_CONFIG_KEY` - the name of the file in S3 bucket that will hold configuration data
* `TOTAL_NUMBER_OF_NODES` - total number of validator nodes (int)
* `SSM_PREFIX` - AWS SSM Parameter Store prefix used to store secrets
* `NODE_NAME` - the name of node that will be used to differentiate stored secrets
* `NODE_IP` - ip address of the node

### Chain configuration options
All chain configuration options are well explained in the [docs](https://docs.polygon.technology/docs/edge/get-started/cli-commands#genesis-flags)  


