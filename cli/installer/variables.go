package installer

type Variables struct {
	App           string
	AppDir        string
	CommonDir     string
	DataDir       string
	Domain        string
	AppDomain     string
	Hostname      string
	Url           string
	AuthUrl       string
	AuthSocket    string
	StunHost      string
	StunPort      int
	HeadscalePort int
	MetricsPort   int
	GrpcPort      int
	DerpStunPort  int
}
