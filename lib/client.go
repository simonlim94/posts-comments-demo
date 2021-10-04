package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/simonlim94/posts-comments-demo/model"
)

const (
	baseURL         = "https://jsonplaceholder.typicode.com"
	getCommentsPath = "comments"
	getPostsPath    = "posts"
)

func GetComments() ([]model.Comment, error) {
	req, err := prepareRequest(baseURL, getCommentsPath)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http response body")
	}

	if resp.StatusCode == http.StatusOK {
		getCommentsResp := []model.Comment{}
		if err := json.Unmarshal(body, &getCommentsResp); err != nil {
			return nil, errors.Wrap(err, "failed to decode get comments response")
		}

		return getCommentsResp, nil
	}

	return nil, errors.Errorf("failed to get comments from client, resp status: %d, err: %s", resp.StatusCode, string(body))
}

func GetPosts() ([]model.Post, error) {
	req, err := prepareRequest(baseURL, getPostsPath)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http response body")
	}

	if resp.StatusCode == http.StatusOK {
		getPostsResp := []model.Post{}
		if err := json.Unmarshal(body, &getPostsResp); err != nil {
			return nil, errors.Wrap(err, "failed to decode get posts response")
		}

		return getPostsResp, nil
	}

	return nil, errors.Errorf("failed to get posts from client, resp status: %d, err: %s", resp.StatusCode, string(body))
}

func GetPost(id uint32) (*model.Post, error) {
	req, err := prepareRequest(baseURL, fmt.Sprintf("%s/%d", getPostsPath, id))
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http response body")
	}

	if resp.StatusCode == http.StatusOK {
		getPostResp := &model.Post{}
		if err := json.Unmarshal(body, getPostResp); err != nil {
			return nil, errors.Wrap(err, "failed to decode get post response")
		}

		return getPostResp, nil
	}

	return nil, errors.Errorf("failed to get post from client, resp status: %d, err: %s", resp.StatusCode, string(body))
}

func prepareRequest(host string, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", host, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
