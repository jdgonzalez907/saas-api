package domain

import "time"

type PaginatedPosts struct {
	posts           []*Post
	nextPublishedAt *time.Time
	nextID          *int64
}

func NewPaginatedPosts(posts []*Post, nextPublishedAt *time.Time, nextID *int64) PaginatedPosts {
	return PaginatedPosts{
		posts:           posts,
		nextPublishedAt: nextPublishedAt,
		nextID:          nextID,
	}
}

func (p PaginatedPosts) Posts() []*Post {
	return p.posts
}

func (p PaginatedPosts) NextPublishedAt() *time.Time {
	return p.nextPublishedAt
}

func (p PaginatedPosts) NextID() *int64 {
	return p.nextID
}

type PaginatedPostsDTO struct {
	Posts           []PostDTO  `json:"posts"`
	NextPublishedAt *time.Time `json:"next_published_at"`
	NextID          *int64     `json:"next_id"`
}

func (p PaginatedPosts) ToDTO() *PaginatedPostsDTO {
	postsDTO := make([]PostDTO, len(p.posts))
	for i, post := range p.posts {
		postsDTO[i] = *post.ToDTO()
	}
	return &PaginatedPostsDTO{
		Posts:           postsDTO,
		NextPublishedAt: p.nextPublishedAt,
		NextID:          p.nextID,
	}
}
