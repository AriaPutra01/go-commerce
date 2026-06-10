package constant

type Role string
type Permission string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleAdmin      Role = "ADMIN"
	RoleUser       Role = "USER"
)

const (
	ProductRead    Permission = "product:read"
	ProductCreate  Permission = "product:create"
	ProductUpdate  Permission = "product:update"
	ProductDelete  Permission = "product:delete"
	OrderRead      Permission = "order:read"
	OrderReadAll   Permission = "order:read_all"
	OrderCreate    Permission = "order:create"
	OrderUpdate    Permission = "order:update"
	OrderDelete    Permission = "order:delete"
	CartRead       Permission = "cart:read"
	CartUpdate     Permission = "cart:update"
	RatingRead     Permission = "rating:read"
	RatingCreate   Permission = "rating:create"
	RatingDelete   Permission = "rating:delete"
	DiscountRead   Permission = "discount:read"
	DiscountManage Permission = "discount:manage"
	WishlistRead   Permission = "wishlist:read"
	WishlistUpdate Permission = "wishlist:update"
	NewsRead       Permission = "news:read"
	NewsManage     Permission = "news:manage"
	AdminManage    Permission = "admin:manage"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleUser:
		return true
	}
	return false
}

func GetPermissionByRole(role Role) []Permission {
	perms := map[Role][]Permission{
		RoleSuperAdmin: {
			AdminManage,
			ProductRead, ProductCreate, ProductUpdate, ProductDelete,
			OrderReadAll, OrderUpdate, OrderDelete,
			DiscountRead, DiscountManage,
			NewsRead, NewsManage,
			RatingRead, RatingDelete,
		},
		RoleAdmin: {
			ProductRead, ProductCreate, ProductUpdate,
			OrderReadAll, OrderUpdate,
			DiscountRead, DiscountManage,
			NewsRead, NewsManage,
			RatingRead, RatingDelete,
		},
		RoleUser: {
			ProductRead,
			DiscountRead,
			NewsRead,
			CartRead, CartUpdate,
			WishlistRead, WishlistUpdate,
			OrderRead, OrderCreate,
			RatingRead, RatingCreate,
		},
	}

	return perms[role]
}
