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

	ObjectStorageBucket       string `env:"OBJECT_STORAGE_BUCKET" envDefault:"zsmartex-tech"`
	ObjectStorageAccessKey    string `env:"OBJECT_STORAGE_ACCESS_KEY" envDefault:"AKIA6MAKWE2NUXH5KAG4"`
	ObjectStorageAccessSecret string `env:"OBJECT_STORAGE_ACCESS_SECRET" envDefault:"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"`
	ObjectStorageRegion       string `env:"OBJECT_STORAGE_REGION" envDefault:"us-east-1"`

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
