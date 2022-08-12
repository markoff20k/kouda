package types

type ENV struct {
	Domain          string `env:"DOMAIN" envDefault:"zsmartex.tech"`
	ApplicationName string `env:"APP_NAME" envDefault:"Kouda"`

	EventAPIJWTPrivateKey string `env:"EVENT_API_JWT_PRIVATE_KEY"`

	DatabaseHost string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort int    `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseUser string `env:"DATABASE_USER" envDefault:"root"`
	DatabasePass string `env:"DATABASE_PASS" envDefault:"changeme"`
	DatabaseName string `env:"DATABASE_NAME" envDefault:"kouda"`

	ObjectStorageBucket       string `env:"OBJECT_STORAGE_BUCKET" envDefault:"zsmartex"`
	ObjectStorageAccessKey    string `env:"OBJECT_STORAGE_ACCESS_KEY" envDefault:"DO00KNY6QY9WY74X3ULM"`
	ObjectStorageAccessSecret string `env:"OBJECT_STORAGE_ACCESS_SECRET" envDefault:"1oRKb2zH8H2VcgR4goLnC6r4DNQr7TfeSWsUKUCRQbo"`
	ObjectStorageRegion       string `env:"OBJECT_STORAGE_REGION" envDefault:"sgp1"`
	ObjectStorageEnpoint      string `env:"OBJECT_STORAGE_ENPOINT" envDefault:""`
	ObjectStorageVersion      int    `env:"OBJECT_STORAGE_VERSION" envDefault:2`

	JWTPublicKey string `env:"JWT_PUBLIC_KEY"`
}

type KoudaConfig struct {
}

type AbilityRole string
type AbilityAdminPermission string

const (
	AbilityAdminPermissionRead   AbilityAdminPermission = "read"
	AbilityAdminPermissionManage AbilityAdminPermission = "manage"
)

type Abilities struct {
	Roles            []AbilityRole                                       `yaml:"roles"`
	AdminPermissions map[AbilityRole]map[AbilityAdminPermission][]string `yaml:"admin_permissions"`
}
