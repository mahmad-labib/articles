package validation

import (
	"fmt"

	"github.com/mahmad-labib/articles/internal/domain"
)

type CommentValidator struct{}

func NewCommentValidator() *CommentValidator {
	return &CommentValidator{}
}

func (v *CommentValidator) Validate(comment *domain.Comment) error {
	if comment.ArticleID == 0 {
		return fmt.Errorf("%v comment must be associated with an article", Validation)
	}
	if comment.AuthorID == 0 {
		return fmt.Errorf("%v comment must be associated with a user", Validation)
	}
	if len(comment.Body) == 0 {
		return fmt.Errorf("%v comment body cannot be empty", Validation)
	}
	return nil
}
