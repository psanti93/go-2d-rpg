package entities

type Enemy struct {
	*Sprite       // struct embedding
	FollowsPlayer bool
}
