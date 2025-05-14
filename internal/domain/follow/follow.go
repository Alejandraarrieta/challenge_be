package domain

type Follow struct {
	ID         uint64 // ID de la relación
	FollowerID uint64 // El que sigue
	FolloweeID uint64 // A quién sigue
	CreatedAt  string // Fecha de creación
}
