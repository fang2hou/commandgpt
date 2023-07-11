package resource

type Advice struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Advices struct {
	Advices []Advice `json:"advices"`
}
