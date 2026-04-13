package model

type Gift struct {
	ID          int64
	WishlistID  int64
	Name        string
	Description string
	Link        string
	Priority    int //1..5
	Booked      bool
}
