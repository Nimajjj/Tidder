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
	Id         int `json:"id_subject_access"`
	CreatePost int `json:"create_post"`
	Pin        int `json:"pin_post"`
	ManagePost int `json:"manage_post"`
	BanUser    int `json:"ban_user"`
	ManageRole int `json:"manage_role"`
	ManageSub  int `json:"manage_subtidder"`
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
	CanCreatePost  int    `json:"can_create_post"`
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

type PostCreationTest struct {
	Sub        Subject
	UserRole   SubjectRoles
	UserAccess SubjectAccess
}

// ViewData

type DisplayablePost struct {
	// POST RELATED
	Id               int    `json:"idpost"`
	Title            string `json:"title"`
	MediaUrl         string `json:"media_url"`
	Content          string `json:"content"`
	CreationDate     string `json:"creation_date"`
	Nsfw             bool   `json:"nsfw"`
	Redacted         bool   `json:"redacted"`
	Pinned           bool   `json:"pinned"`
	Score            int
	NumberOfComments int
	Vote             int

	// AUTHOR
	AuthorName string

	// SUBTIDDER
	SubtidderName string
	SubtidderPP   string
}

type DisplayableComment struct {
	// COMMENT RELATED
	Id           int
	Content      string
	CreationDate string
	Upvotes      int
	Downvotes    int
	Redacted     bool
	Response     []DisplayableComment
	IdPost       int
	Vote         int

	Score int

	// AUTHOR
	AuthorName string
	AuthorPP   string
}

type RoleAccess struct {
	Role   SubjectRoles
	Access SubjectAccess
}

type AccountSubscribed struct {
	Account Accounts
	Banned  bool
	Role    SubjectRoles
}

type MasterVD struct {
	IndexVD       IndexViewData
	SubtidderVD   SubtidderViewData
	SearchVD      SearchViewData
	CreatePostsVD CreatePostsVD
	PostVD        PostVD
	ProfilePageVD ProfilePageVD
	Account       Accounts
	Errors        Errors

	Connected bool
	Page      string
}

type IndexViewData struct {
	Posts []DisplayablePost
}

type SubtidderViewData struct {
	Sub            Subject
	Posts          []DisplayablePost
	SubscribedUser []AccountSubscribed
	Subscribed     bool
	Roles          []RoleAccess

	UserRole RoleAccess
}

type SearchViewData struct {
	Subjects map[Subject]int
}

type CreatePostsVD struct {
	PostCreation []PostCreationTest
}

type PostVD struct {
	Post      DisplayablePost
	Subtidder Subject
	Comments  []DisplayableComment
}

type ProfilePageVD struct {
	Account    Accounts
	Posts      []DisplayablePost
	Subtidders []Subject
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
