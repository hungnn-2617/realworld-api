package errors

// Common error messages used across the application
const (
	// Authentication & Authorization errors
	ErrUnauthorized       = "unauthorized"
	ErrForbidden          = "forbidden"
	ErrInvalidCredentials = "invalid email or password"
	ErrInvalidToken       = "invalid token"
	ErrTokenExpired       = "token expired"
	ErrTokenRequired      = "authorization token required"

	// User errors
	ErrUserNotFound         = "user not found"
	ErrUserAlreadyExists    = "user already exists"
	ErrEmailAlreadyTaken    = "email already taken"
	ErrUsernameAlreadyTaken = "username already taken"

	// Profile errors
	ErrCannotFollowYourself = "cannot follow yourself"

	// Article errors
	ErrArticleNotFound     = "article not found"
	ErrNotAuthorizedUpdate = "not authorized to update this article"
	ErrNotAuthorizedDelete = "not authorized to delete this article"

	// Comment errors
	ErrCommentNotFound            = "comment not found"
	ErrNotAuthorizedDeleteComment = "not authorized to delete this comment"
	ErrInvalidCommentID           = "invalid comment id"

	// Favorite errors
	ErrAlreadyFavorited = "article already favorited"
	ErrNotFavorited     = "article not favorited"

	// Validation errors
	ErrValidationFailed = "validation failed"
	ErrInvalidInput     = "invalid input"

	// Server errors
	ErrInternalServer = "internal server error"
	ErrFailedToFetch  = "failed to fetch data"
	ErrFailedToCreate = "failed to create"
	ErrFailedToUpdate = "failed to update"
	ErrFailedToDelete = "failed to delete"
)

// Helper functions to check error messages
func IsArticleNotFound(err error) bool {
	return err != nil && err.Error() == ErrArticleNotFound
}

func IsUserNotFound(err error) bool {
	return err != nil && err.Error() == ErrUserNotFound
}

func IsCommentNotFound(err error) bool {
	return err != nil && err.Error() == ErrCommentNotFound
}

func IsNotAuthorizedUpdate(err error) bool {
	return err != nil && err.Error() == ErrNotAuthorizedUpdate
}

func IsNotAuthorizedDelete(err error) bool {
	return err != nil && err.Error() == ErrNotAuthorizedDelete
}

func IsNotAuthorizedDeleteComment(err error) bool {
	return err != nil && err.Error() == ErrNotAuthorizedDeleteComment
}

func IsCannotFollowYourself(err error) bool {
	return err != nil && err.Error() == ErrCannotFollowYourself
}
