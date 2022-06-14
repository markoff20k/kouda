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
