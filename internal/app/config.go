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

	StorageBackend string `envconfig:"STORAGE_BACKEND" required:"true" default:"textile"`

	TextileAuthKey       string `envconfig:"TEXTILE_AUTH_KEY" required:"false"`
	TextileAuthSecret    string `envconfig:"TEXTILE_AUTH_SECRET" required:"false"`
	TextileThreadID      string `envconfig:"TEXTILE_THREAD_ID" required:"false"`
	TextileBucketRootKey string `envconfig:"TEXTILE_BUCKET_ROOT_KEY" required:"false"`

	NftStorageApiKey string `envconfig:"NFTSTORAGE_API_KEY" required:"false"`

	BlockchainURL         string `envconfig:"BLOCKCHAIN_URL" default:"http://localhost:8545"`
	ERC721ContractAddress string `envconfig:"ERC721_CONTRACT_ADDRESS"`
	ERC721ContractKeyFile string `envconfig:"ERC721_CONTRACT_KEY"`
	ERC721ContractKeyPass string `envconfig:"ERC721_CONTRACT_KEY_PASS"`
}
