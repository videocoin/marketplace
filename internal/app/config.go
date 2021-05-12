package app

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	Addr       string `envconfig:"ADDR" default:"0.0.0.0:8088"`
	DBURI      string `envconfig:"DBURI" default:"host=127.0.0.1 port=5432 dbname=marketplace sslmode=disable"`
	AuthSecret string `envconfig:"AUTH_SECRET" default:"secret"`

	EnableTranscoding     bool   `envconfig:"ENABLE_TRANSCODING" default:"false"`
	GCPBucket             string `envconfig:"GCP_BUCKET" required:"false"`
	GCPProject            string `envconfig:"GCP_PROJECT" required:"false"`
	GCPRegion             string `envconfig:"GCP_REGION" required:"false"`
	GCPPubSubTopic        string `envconfig:"GCP_PUBSUB_TOPIC" required:"false"`
	GCPPubSubSubscription string `envconfig:"GCP_PUBSUB_SUB" required:"false"`

	TextileAuthKey       string `envconfig:"TEXTILE_AUTH_KEY"`
	TextileAuthSecret    string `envconfig:"TEXTILE_AUTH_SECRET"`
	TextileThreadID      string `envconfig:"TEXTILE_THREAD_ID"`
	TextileBucketRootKey string `envconfig:"TEXTILE_BUCKET_ROOT_KEY"`

	BlockchainURL          string `envconfig:"BLOCKCHAIN_URL" default:"http://localhost:8545"`
	ERC1155ContractAddress string `envconfig:"ERC1155_CONTRACT_ADDRESS"`
	ERC1155ContractKeyFile string `envconfig:"ERC1155_CONTRACT_KEY"`
	ERC1155ContractKeyPass string `envconfig:"ERC1155_CONTRACT_KEY_PASS"`
}
