package dtos

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment" binding:"required"`
}

type CommentResponse struct {
	Comment CommentData `json:"comment"`
}

type CommentsResponse struct {
	Comments []CommentData `json:"comments"`
}

type CommentData struct {
	ID        uint        `json:"id"`
	Body      string      `json:"body"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Author    ProfileData `json:"author"`
}
