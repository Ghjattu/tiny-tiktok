package redis

type VideoCache struct {
	ID            int64  `redis:"id"`
	AuthorID      int64  `redis:"author_id"`
	PlayUrl       string `redis:"play_url"`
	CoverUrl      string `redis:"cover_url"`
	FavoriteCount int64  `redis:"favorite_count"`
	CommentCount  int64  `redis:"comment_count"`
	Title         string `redis:"title"`
}

type CommentCache struct {
	ID         int64  `redis:"id"`
	UserID     int64  `redis:"user_id"`
	Content    string `redis:"content"`
	CreateDate string `redis:"create_date"`
}
