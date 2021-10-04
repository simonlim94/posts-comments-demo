package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/simonlim94/posts-comments-demo/lib"
	"github.com/simonlim94/posts-comments-demo/model"
	"github.com/simonlim94/posts-comments-demo/util"
)

type FilterRelationship string

const (
	FilterRelationshipAnd = "and"
	FilterRelationshipOr  = "or"
)

type TopPost struct {
	PostID                uint32 `json:"post_id"`
	PostTitle             string `json:"post_title"`
	PostBody              string `json:"post_body"`
	TotalNumberOfComments uint32 `json:"total_number_of_comments"`
}

type GetCommentsByFilterRequest struct {
	Filters            []Filter           `json:"filters"`
	FilterRelationship FilterRelationship `json:"filterRelationship,omitempty"`
}

type Filter struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func GetTopPostsByComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(util.NotFoundErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	comments, err := lib.GetComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(util.InternalServerErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	// store hash map of post id with value of comment ids
	postsByComments := make(map[uint32]uint32)
	for _, comment := range comments {
		_, ok := postsByComments[comment.PostID]

		if ok {
			postsByComments[comment.PostID]++
		} else {
			postsByComments[comment.PostID] = 1
		}
	}

	posts, err := lib.GetPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(util.InternalServerErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	result := make([]TopPost, 0)
	for _, post := range posts {
		result = append(result, TopPost{
			PostID:                post.ID,
			PostTitle:             post.Title,
			PostBody:              post.Body,
			TotalNumberOfComments: postsByComments[post.ID],
		})
	}

	// sort top posts by descending order of number of comments
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalNumberOfComments > result[j].TotalNumberOfComments
	})

	if err := util.WriteOKResponse(w, result); err != nil {
		log.Printf("failed to encode response body, err: %v", err)
	}
}

func GetCommentsByFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(util.NotFoundErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	input := &GetCommentsByFilterRequest{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(util.BadRequestErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	if input.FilterRelationship == "" {
		input.FilterRelationship = FilterRelationshipAnd
	}

	comments, err := lib.GetComments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(util.InternalServerErrorResponseBody); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	result, err := CommentFiltering(comments, input.Filters, input.FilterRelationship)
	if err != nil {
		errResp := util.BadRequestErrorResponseBody
		errResp.Error = err.Error()

		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errResp); err != nil {
			log.Printf("failed to encode response body, err: %v", err)
		}
		return
	}

	if err := util.WriteOKResponse(w, result); err != nil {
		log.Printf("failed to encode response body, err: %v", err)
	}
}

func CommentFiltering(comments []model.Comment, filters []Filter, filterRelationship FilterRelationship) ([]model.Comment, error) {
	var filteredResults []model.Comment
	for i, filter := range filters {
		// skip first round of assigned back to filtered result
		if i != 0 && filterRelationship == FilterRelationshipAnd {
			comments = filteredResults
		}

		for _, comment := range comments {
			switch filter.Field {
			case "postId":
				value, ok := filter.Value.(float64)
				if !ok {
					return nil, fmt.Errorf("invalid data type for %q is provided", filter.Field)
				}

				if comment.PostID == uint32(value) {
					filteredResults = append(filteredResults, comment)
				}
			case "id":
				value, ok := filter.Value.(float64)
				if !ok {
					return nil, fmt.Errorf("invalid data type for %q is provided", filter.Field)
				}

				if comment.ID == uint32(value) {
					filteredResults = append(filteredResults, comment)
				}
			case "name":
				_, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("invalid data type for %q is provided", filter.Field)
				}

				if comment.Name == filter.Value {
					filteredResults = append(filteredResults, comment)
				}
			case "email":
				_, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("invalid data type for %q is provided", filter.Field)
				}

				if comment.Email == filter.Value {
					filteredResults = append(filteredResults, comment)
				}
			case "body":
				_, ok := filter.Value.(string)
				if !ok {
					return nil, fmt.Errorf("invalid data type for %q is provided", filter.Field)
				}

				if comment.Body == filter.Value {
					filteredResults = append(filteredResults, comment)
				}
			default:
				return nil, fmt.Errorf("invalid field %q is provided", filter.Field)
			}
		}
	}

	return filteredResults, nil
}
