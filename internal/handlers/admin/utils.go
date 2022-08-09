package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/utils"

	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/types"
)

var (
	ErrAbilityNotPermitted = pkg.NewError(fiber.StatusUnauthorized, "admin.ability.not_permitted")
)

func (h Handler) adminAuthorize(c *fiber.Ctx, permission types.AbilityAdminPermission, resource string) {
	currentUser := c.Locals("current_member").(*models.Member)

	resources := h.abilities.AdminPermissions[types.AbilityRole(currentUser.Role)][types.AbilityAdminPermissionManage]
	resources = append(resources, h.abilities.AdminPermissions[types.AbilityRole(currentUser.Role)][permission]...)

	if !utils.Contains(resources, resource) {
		panic(ErrAbilityNotPermitted)
	}
}
