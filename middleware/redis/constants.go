package redis

const (
	// Hash user:<user_id> save the user's information.
	UserKey = "user:"

	// Hash video:<video_id> save the video's information.
	VideoKey = "video:"
	// List favorite_videos:<user_id> save the user's favorite videos id.
	FavoriteVideosKey = "favorite_videos:"
	// List videos:author:<author_id> save the author's videos id.
	VideosByAuthorKey = "videos:author:"

	// List following:<user_id> save the user's following user id.
	FollowingKey = "following:"
	// List follower:<user_id> save the user's follower user id.
	FollowerKey = "follower:"

	// Hash comment:<comment_id> save the comment's information.
	CommentKey = "comment:"
	// List comments:video:<video_id> save the video's comments id.
	CommentsByVideoKey = "comments:video:"
)
