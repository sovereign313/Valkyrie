package vtypes

type ForemanC struct {
	DupProtection string `json:"DupProtection"` 
	ProtectTime string `json:"ProtectTime"` 
	Host string `json:"Host"` 
}

type LauncherC struct {
	DefaultImage string   `json:"DefaultImage"`
	Host         []string `json:"Host"`
}

type WorkerC struct { 
	ExternalPath string `json:"ExternalPath"` 
	GitRepoURL string `json:"GitRepoUrl"` 
	SSHUser string `json:"SSHUser"` 
	SSHPrivateKey string `json:"SSHPrivateKey"` 
	SSHPublicKey string `json:"SSHPublicKey"` 
}

type LoggerC struct {
	UseEventStreams  string `json:"UseEventStreams"`
	UseSecureLogging string `json:"UseSecureLogging"`
	ESHostPort       string `json:"ESHostPort"`
	LogKey           string `json:"LogKey"`
	LogFileLocation  string `json:"LogFileLocation"`
	Host             string `json:"Host"`
}

type AlerterC struct {
	EmailServer       string `json:"EmailServer"`
	FromAddress       string `json:"FromAddress"`
	TwilioAccount     string `json:"TwilioAccount"`
	TwilioToken       string `json:"TwilioToken"`
	TwilioPhoneNumber string `json:"TwilioPhoneNumber"`
	Host              string `json:"Host"`
}

type AWSC struct {
	Region        string `json:"Region"`
	SQSName       string `json:"SQSName"`
	AWSAccessKey  string `json:"AWSAccessKey"`
	AWSSecretKey  string `json:"AWSSecretKey"`
	EncryptionKey string `json:"EncryptionKey"`
}

type DispatcherC struct {
	DupProtection string `json:"DupProtection"`
	ProtectTime   string `json:"ProtectTime"`
	Host          string `json:"Host"`
}

type SQSReaderC struct {
	SleepTimeout string `json:"SleepTimeout"`
	Host         string `json:"Host"`
}

type MailReaderC struct {
	DupProtection      string `json:"DupProtection"`
	UseTLS             string `json:"UseTLS"`
	MailProtocol       string `json:"MailProtocol"`
	ProtectTime        string `json:"ProtectTime"`
	SleepTimeout       string `json:"SleepTimeout"`
	MailServer         string `json:"MailServer"`
	MailUser           string `json:"MailUser"`
	MailPassword       string `json:"MailPassword"`
	MailSubjectTrigger string `json:"MailSubjectTrigger"`
	Host               string `json:"Host"`
}

type ValkConfig struct {
	UseDocker bool
	LicenseKey string
	BusinessName string
	HostConfig []string
	ForemanConfig ForemanC
	LauncherConfig LauncherC
	WorkerConfig WorkerC
	LoggerConfig LoggerC
	AlerterConfig AlerterC
	AWSConfig AWSC
	DispatcherConfig DispatcherC
	SQSReaderConfig SQSReaderC
	MailReaderConfig MailReaderC
}


