package database

type permission uint64

const (
	NO_PERMISSIONS     permission = 0
	ADMIN              permission = (1 << (iota))
	USER_READ          permission = (1 << (iota))
	USER_WRITE         permission = (1 << (iota))
	ADD_POINTS_READ    permission = (1 << (iota))
	ADD_POINTS_WRITE   permission = (1 << (iota))
	SPENT_POINTS_READ  permission = (1 << (iota))
	SPENT_POINTS_WRITE permission = (1 << (iota))
	PARTICIPENT_READ   permission = (1 << (iota))
	PARTICIPENT_WRITE  permission = (1 << (iota))
	CATEGORY_READ      permission = (1 << (iota))
	CATEGORY_WRITE     permission = (1 << (iota))
	EVENT_READ         permission = (1 << (iota))
	EVENT_WRITE        permission = (1 << (iota))
)

func Set(value permission, flag permission) permission {
	return value | flag
}

func Unset(value permission, flag permission) permission {
	return value & ^flag
}
