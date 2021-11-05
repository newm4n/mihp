package notification

import "fmt"

type Recipient struct {
	Name  string
	Email string
}

func (r *Recipient) String() string {
	if len(r.Name) == 0 {
		return r.Email
	}
	return fmt.Sprintf("%s <%s>", r.Name, r.Email)
}
