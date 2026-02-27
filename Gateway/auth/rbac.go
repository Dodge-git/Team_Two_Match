package auth

func HasPermission(role string,permission string)bool{
	if role == RoleAdmin{
		return true
	}
	permissions , exist := RolePermissions[role]
	if exist == false {
		return false
	}
	for _, r := range permissions{
		if r == permission{
			return true
		}
	}
	return false
}
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
		DeleteMatch,
		CreateMatch,
		AdSport ,
		AdTeam ,
		ManageMatch,
		AdGoal,
		AdPlayer,
	},
	RoleAdmin: {
	
	},
}

