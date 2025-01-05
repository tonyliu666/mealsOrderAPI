package database

type Client struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f *Client) Save() error {
	_, err := o.Insert(f)
	return err
}
func (f *Client) Read() error {
	err := o.Read(f)
	return err
}
