package auth

var RolePermissions = map[string][]string{
	RoleUser: {
		GetProfile,
		DeleteProfile,
	},
	RoleFan: {
		AdComment,
		AdFavoriteSport,
		GetProfile,
		DeleteProfile,
	},
	RoleCommentator: {
		AdComment,
		AdCommentary,
		AdEvent,
		AdFavoriteSport,
		GetProfile,
		DeleteComment,
		DeleteCommentary,
		DeleteEvent,
		DeleteProfile,
	},
	RoleAdmin: {
		DeleteComment,
		DeleteCommentary,
		DeleteEvent,
		DeleteProfile,
		AdComment,
		AdCommentary,
		AdEvent,
		AdFavoriteSport,
		GetProfile,
		EditProfile,
	},
}
