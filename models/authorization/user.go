package authorization

type User struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
}