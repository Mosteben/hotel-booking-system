package model

type SearchRequest struct {
	Name        string  `form:"name"`
	City        string  `form:"city"`
	Country     string  `form:"country"`
	Stars       int     `form:"stars"`
	MinPrice    float64 `form:"min_price"`
	MaxPrice    float64 `form:"max_price"`
	RoomType    string  `form:"room_type"`
	MinCapacity int     `form:"min_capacity"`
}
