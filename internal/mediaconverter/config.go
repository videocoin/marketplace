package mediaconverter

type GCPConfig struct {
	Bucket             string
	Project            string
	Region             string
	PubSubTopic        string
	PubSubSubscription string
}
