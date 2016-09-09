package mailserf

import (
	"fmt"
	"time"
)


const (

)

type Envelope struct {
	To		string
	From	string
	DateTime	time.Time
}



func (e Envelope) String() string {
		return fmt.Sprintf("{To=<%s>; From <%s>}", e.To, e.From)
}

