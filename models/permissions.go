package models

type permission uint64

const (
	NO_PERMISSION      permission = 0
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
