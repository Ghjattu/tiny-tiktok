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

type UserCache struct {
	ID              int64  `redis:"id"`
	Name            string `redis:"name"`
	FollowCount     int64  `mapstructure:"follow_count" redis:"follow_count"`
	FollowerCount   int64  `mapstructure:"follower_count" redis:"follower_count"`
	Avatar          string `redis:"avatar"`
	BackgroundImage string `mapstructure:"background_image" redis:"background_image"`
	Signature       string `redis:"signature"`
	TotalFavorited  int64  `mapstructure:"total_favorited" redis:"total_favorited"`
	WorkCount       int64  `mapstructure:"work_count" redis:"work_count"`
	FavoriteCount   int64  `mapstructure:"favorite_count" redis:"favorite_count"`
}
