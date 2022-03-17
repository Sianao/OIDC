package models

import "time"

type Request struct {
	Scope        string `json:"scope" form:"scope"`
	ResponseType string `json:"response_type" form:"response_type"`
	ClientId     string `json:"client_id,omitempty" form:"client_id"`
	State        string `json:"state,omitempty" form:"state"`
	RedirectUrl  string `json:"redirect_url,omitempty" form:"redirect_url"`
}
type Admin struct {
	Name     string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}
type Token struct {
	AccessToken string `json:"access_Token,omitempty"`
	TokenType   string `json:"token_Type,omitempty"`
	Scope       string `json:"scope,omitempty"`
}
type HubBasicInfo struct {
	Login            string `json:"login"`
	ID               int    `json:"id"`
	NodeID           string `json:"node_id"`
	AvatarURL        string `json:"avatar_url"`
	GravatarID       string `json:"gravatar_id"`
	URL              string `json:"url"`
	HTMLURL          string `json:"html_url"`
	FollowersURL     string `json:"followers_url"`
	FollowingURL     string `json:"following_url"`
	GistsURL         string `json:"gists_url"`
	StarredURL       string `json:"starred_url"`
	SubscriptionsURL string `json:"subscriptions_
url"`
	OrganizationsURL string `json:"organizations_url"`
	ReposURL         string `json:"repos_url"`
	EvenTsURL        string `json:"even
ts_url"`
	ReceivedEventsURL string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           interface{} `json:"company"`
	Blog              string      `json:"blog"`
	Location          interface{} `json:"location"`
	Email             interface{} `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               interface{} `json:"bio"`
	TwitterUsername   interface{} `json:"twitter_username"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}
