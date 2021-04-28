package app

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr    string `envconfig:"RPC_ADDR" default:"0.0.0.0:8888"`
	GWAddr     string `envconfig:"GW_ADDR" default:"0.0.0.0:8088"`
	DBURI      string `envconfig:"DBURI" default:"host=127.0.0.1 port=5432 dbname=marketplace sslmode=disable"`
	AuthSecret string `envconfig:"AUTH_SECRET"`

	GCPBucket             string `envconfig:"GCP_BUCKET"`
	GCPProject            string `envconfig:"GCP_PROJECT"`
	GCPRegion             string `envconfig:"GCP_REGION"`
	GCPPubSubTopic        string `envconfig:"GCP_PUBSUB_TOPIC"`
	GCPPubSubSubscription string `envconfig:"GCP_PUBSUB_SUB"`

	TextileAuthKey       string `envconfig:"TEXTILE_AUTH_KEY"`
	TextileAuthSecret    string `envconfig:"TEXTILE_AUTH_SECRET"`
	TextileThreadID      string `envconfig:"TEXTILE_THREAD_ID"`
	TextileBucketRootKey string `envconfig:"TEXTILE_BUCKET_ROOT_KEY"`
}
