package dtos

type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Body        string   `json:"body" binding:"required"`
		TagList     []string `json:"tagList,omitempty"`
	} `json:"article" binding:"required"`
}

type UpdateArticleRequest struct {
	Article struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Body        *string `json:"body,omitempty"`
	} `json:"article" binding:"required"`
}

type ArticleResponse struct {
	Article ArticleData `json:"article"`
}

type ArticlesResponse struct {
	Articles      []ArticleData `json:"articles"`
	ArticlesCount int           `json:"articlesCount"`
}

type ArticleData struct {
	Slug           string      `json:"slug"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Body           string      `json:"body"`
	TagList        []string    `json:"tagList"`
	CreatedAt      string      `json:"createdAt"`
	UpdatedAt      string      `json:"updatedAt"`
	Favorited      bool        `json:"favorited"`
	FavoritesCount int         `json:"favoritesCount"`
	Author         ProfileData `json:"author"`
}
