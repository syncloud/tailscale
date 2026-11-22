package installer

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/syncloud/golib/config"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/golib/platform"
	"go.uber.org/zap"
)

const (
	App           = "tailscale"
	UserName      = "tailscale"
	HeadscalePort = 18731
	MetricsPort   = 18732
	GrpcPort      = 18733
	DerpStunPort  = 18734
	StunPort      = 3478
)

type Installer struct {
	appDir         string
	commonDir      string
	dataDir        string
	configDir      string
	platformClient *platform.Client
	logger         *zap.Logger
}

func New(logger *zap.Logger) *Installer {
	appDir := fmt.Sprintf("/snap/%s/current", App)
	dataDir := fmt.Sprintf("/var/snap/%s/current", App)
	return &Installer{
		appDir:         appDir,
		commonDir:      fmt.Sprintf("/var/snap/%s/common", App),
		dataDir:        dataDir,
		configDir:      path.Join(dataDir, "config"),
		platformClient: platform.New(),
		logger:         logger,
	}
}

func (i *Installer) Install() error {
	if err := linux.CreateUser(UserName); err != nil {
		return err
	}
	return i.createDirs()
}

func (i *Installer) Configure() error {
	if err := i.createDirs(); err != nil {
		return err
	}
	if err := i.registerOIDC(); err != nil {
		return err
	}
	if err := i.UpdateConfigs(); err != nil {
		return err
	}
	return i.writeVersion()
}

func (i *Installer) PreRefresh() error {
	return nil
}

func (i *Installer) PostRefresh() error {
	err := os.Remove(path.Join(i.dataDir, "version"))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return i.UpdateConfigs()
}

func (i *Installer) StorageChange() error {
	return nil
}

func (i *Installer) AccessChange() error {
	if err := i.UpdateConfigs(); err != nil {
		return err
	}
	return i.platformClient.RestartService(fmt.Sprintf("%s.headscale", App))
}

func (i *Installer) BackupPreStop() error {
	return nil
}

func (i *Installer) RestorePreStart() error {
	return nil
}

func (i *Installer) RestorePostStart() error {
	return i.UpdateConfigs()
}

func (i *Installer) UpdateConfigs() error {
	if err := linux.CreateUser(UserName); err != nil {
		return err
	}
	if err := i.createDirs(); err != nil {
		return err
	}

	domain, err := i.platformClient.GetDeviceDomainName()
	if err != nil {
		return err
	}
	url, err := i.platformClient.GetAppUrl(App)
	if err != nil {
		return err
	}
	appDomain, err := i.platformClient.GetAppDomainName(App)
	if err != nil {
		return err
	}
	authUrl, err := i.platformClient.GetAppUrl("auth")
	if err != nil {
		return err
	}

	variables := Variables{
		App:           App,
		AppDir:        i.appDir,
		CommonDir:     i.commonDir,
		DataDir:       i.dataDir,
		Domain:        domain,
		AppDomain:     appDomain,
		Hostname:      strings.Split(domain, ".")[0],
		Url:           url,
		AuthUrl:       authUrl,
		StunHost:      StunHost(domain),
		StunPort:      StunPort,
		HeadscalePort: HeadscalePort,
		MetricsPort:   MetricsPort,
		GrpcPort:      GrpcPort,
		DerpStunPort:  DerpStunPort,
	}
	if err := config.Generate(path.Join(i.appDir, "config"), i.configDir, variables); err != nil {
		return err
	}
	if err := i.initRecords(); err != nil {
		return err
	}

	if err := linux.Chown(i.dataDir, UserName); err != nil {
		return err
	}
	return linux.Chown(i.commonDir, UserName)
}

func (i *Installer) createDirs() error {
	for _, dir := range []string{i.configDir, path.Join(i.dataDir, "headscale"), path.Join(i.dataDir, "tailscaled")} {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return err
		}
	}
	return nil
}

func (i *Installer) initRecords() error {
	records := path.Join(i.dataDir, "extra-records.json")
	if _, err := os.Stat(records); err == nil {
		return nil
	}
	return os.WriteFile(records, []byte("[]"), 0644)
}

func (i *Installer) registerOIDC() error {
	secret, err := i.platformClient.RegisterOIDCClient(App, []string{"/oidc/callback"}, true, "client_secret_basic")
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(i.configDir, "oidc.secret"), []byte(secret), 0600)
}

func (i *Installer) writeVersion() error {
	version, err := os.ReadFile(path.Join(i.appDir, "version"))
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(i.dataDir, "version"), version, 0644)
}
