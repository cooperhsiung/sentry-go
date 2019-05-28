package alert

type Task struct {
	Project string
	Type    string
	Field   string
	Every   string
	Offset  string
	Period  string
	Rules   []Rule
	Title   string
	Message string
	Post    string
	Mail    []string
	Wechat  []string
}

type Rule struct {
	Expr  string
	Level string
}

type Expr struct {
	Method    string
	Value     string
	Operator  string
	Threshold float64
}

type Conf struct {
	Name  string
	Tasks []Task
}
