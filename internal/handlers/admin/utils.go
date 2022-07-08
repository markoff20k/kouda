package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/arrutil"
	"github.com/zsmartex/pkg/v2"

	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/types"
)

var (
	ErrAbilityNotPermitted = pkg.NewError(fiber.StatusUnauthorized, "admin.ability.not_permitted")
)

func (h Handler) adminAuthorize(c *fiber.Ctx, permission types.AbilityAdminPermission, resource string) {
	current_user := c.Locals("current_user").(*models.User)

	resources := h.abilities.AdminPermissions[types.AbilityRole(current_user.Role)][permission]

	if arrutil.NotContains(resources, resource) {
		panic(ErrAbilityNotPermitted)
	}
}
