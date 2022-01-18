package config

import (
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

type infrastructure struct {
	filename string
	log      logger.Logger
}

type (
	CfgServer struct {
		Port string `env:"SERVER_PORT"`
	}

	CfgServices struct {
		HMACSecret   string `env:"SECRET_HMAC"`
		PasswordSalt string `env:"PASSWORD_SALT"`
	}

	CfgPostgres struct {
		DriverName    string `yaml:"driverName"`
		ConnectFormat string `yaml:"postgresConnectFormat"`
		Username      string `env:"POSTGRES_USER"`
		Password      string `env:"POSTGRES_PASSWORD"`
		DBName        string `env:"POSTGRES_DBNAME"`
		SSLMode       string `env:"POSTGRES_SSLMODE"`
	}

	Endpoints struct {
		AuthEndpoints struct {
			Auth          string `yaml:"auth" json:"Auth"`
			SignIn        string `yaml:"signIn" json:"SignIn"`
			SignUp        string `yaml:"signUp" json:"SignUp"`
			SignOut       string `yaml:"signOut" json:"SignOut"`
			NewPassword   string `yaml:"newPassword" json:"NewPassword"`
			DeleteAccount string `yaml:"deleteAccount" json:"DeleteAccount"`
		} `yaml:"authEndpoints" json:"AuthEndpoints"`

		ApiEndpoints struct {
			Api        string `yaml:"api" json:"Api"`
			CreateChar string `yaml:"createChar" json:"CreateChar"`
			GetAllChar string `yaml:"getAllChar" json:"GetAllChar"`
			SelectChar string `yaml:"selectChar" json:"SelectChar"`
			UpdChar    string `yaml:"updChar" json:"UpdChar"`
			DelChar    string `yaml:"delChar" json:"DelChar"`
			DelAllChar string `yaml:"delAllChar" json:"DelAllChar"`
		} `yaml:"apiEndpoints" json:"ApiEndpoints"`

		ActionEndpoints struct {
			Action        string `yaml:"action" json:"Action"`
			CharacterMenu string `yaml:"characterMenu" json:"CharacterMenu"`
			ActionScene   string `yaml:"actionScene" json:"ActionScene"`
			Restock       string `yaml:"restock"`
		} `yaml:"actionEndpoints" json:"ActionEndpoints"`
	}
)

func NewInfrastructureConfig(infrastructureFilename string, log logger.Logger) *infrastructure {
	return &infrastructure{
		filename: infrastructureFilename,
		log:      log,
	}
}

func (i infrastructure) InitInfrastructureConfigs(logMsg LogMsg) (*CfgServer, *CfgServices, *CfgPostgres, Endpoints, error) {
	i.log.Infof(logMsg.Format, i.log.PackageAndFileNames(), logMsg.Init)

	var (
		cfgServer   CfgServer
		cfgServices CfgServices
		cfgPostgres CfgPostgres
		endpoints   Endpoints
	)

	re := regexp.MustCompile(`^(.*` + "ArchaicReverie" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		i.log.Errorf(logMsg.Format, i.log.PackageAndFileNames(), err.Error())
	}

	readEnv(i.log, logMsg,
		&cfgServer,
		&cfgServices,
		&cfgPostgres)

	unmarshalYaml(i.filename, i.log,
		&cfgPostgres,
		&endpoints)

	i.log.Infof(logMsg.Format, i.log.PackageAndFileNames(), logMsg.InitOk)

	return &cfgServer, &cfgServices, &cfgPostgres, endpoints, nil
}
