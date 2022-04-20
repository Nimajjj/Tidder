package mySQL

import (
	"time"
)

// SQL STRUCTURES //

type Accounts struct {
	Id             int    `json:"id_account"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"hashed_password"`
	BirthDate      string `json:"birth_date"`
	CreationDate   string `json:"creation_date"`
	Karma          int    `json:"karma"`
	ProfilePicture string `json:"profile_picture"`
	StudentId      string `json:"student_id"`
}

type GlobalRoles struct {
	Id   int    `json:"id_global_role"`
	Name string `json:"name"`
}

type GlobalAccess struct {
	Id   int    `json:"id_global_access"`
	Name string `json:"name"`
}

type GlobalRolesManagement struct {
	IdUser         int `json:"id_user"`
	IdGlobalRole   int `json:"id_global_role"`
	IdGlobalAccess int `json:"id_global_access"`
}

type Comments struct {
	Id           int    `json:"id_comment"`
	Content      string `json:"content"`
	CreationDate string `json:"creation_date"`
	Upvotes      int    `json:"upvotes"`
	Downvotes    int    `json:"downvotes"`
	Redacted     bool   `json:"redacted"`
	IdAuthor     int    `json:"id_author"`
	ResponseToId int    `json:"response_to_id"`
	IdPost       int    `json:"id_post"`

	Score int
}

type HasSubjectRole struct {
	IdAccount     int `json:"id_account"`
	IdSubject     int `json:"id_subject"`
	IdSubjectRole int `json:"id_subject_role"`
}

type IsBan struct {
	IdAccount int `json:"id_account"`
	IdSubject int `json:"id_subject"`
}

type Posts struct {
	Id           int    `json:"idpost"`
	Title        string `json:"title"`
	MediaUrl     string `json:"media_url"`
	Content      string `json:"content"`
	CreationDate string `json:"creation_date"`
	Upvotes      int    `json:"upvotes"`
	Downvotes    int    `json:"downvotes"`
	Nsfw         bool   `json:"nsfw"`
	Redacted     bool   `json:"redacted"`
	Pinned       bool   `json:"pinned"`
	IdSubject    int    `json:"id_subject"`
	IdAuthor     int    `json:"id_author"`

	Score            int
	NumberOfComments int
}

type SubjectAccess struct {
	Id         int  `json:"id_subject_access"`
	PinPost    bool `json:"pin_post"`
	RemovePost bool `json:"remove_post"`
	BanUser    bool `json:"ban_user"`
	CreateRole bool `json:"create_role"`
	GiveRole   bool `json:"give_role"`
	DeleteRole bool `json:"delete_role"`
}

type SubjectRoles struct {
	Id              int    `json:"id_subject_roles"`
	Name            string `json:"pin_post"`
	IdSubject       int    `json:"id_subject"`
	IdSubjectAccess int    `json:"id_subject_access"`
}

type Subject struct {
	Id             int    `json:"id_subject"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	Nsfw           bool   `json:"nsfw"`
	IdOwner        int    `json:"id_owner"`
	Infos          string `json:"infos"`
	Banner         string `json:"banner"`
}

type SubscribeToSubject struct {
	IdAccount int    `json:"id_account"`
	IdSubject string `json:"id_subject"`
}

type Errors struct {
	Signup          string
	Signin          string
	CreateSubtidder string
	CreatePost      string
}

// ViewData

type MasterVD struct {
	SubtidderVD   SubtidderViewData
	SearchVD      SearchViewData
	CreatePostsVD CreatePostsVD
	Account       Accounts
	Errors        Errors

	Connected     bool
}

type SubtidderViewData struct {
	Sub        Subject
	Posts      []map[Posts]Accounts
	Subscribed bool
}

type SearchViewData struct {
	Subjects map[Subject]int
}

type CreatePostsVD struct {
	SubscribedSubjects []Subject
}

func (viewData *MasterVD) ClearErrors() {
	(*viewData).Errors.Signup = ""
	(*viewData).Errors.Signin = ""
	(*viewData).Errors.CreateSubtidder = ""
	(*viewData).Errors.CreatePost = ""
}

// COOKIES
type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}
