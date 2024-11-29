package entities

// Player represents the plyable charcter
type Player struct {
	*Sprite
	Health uint // the health can never be negative
}
