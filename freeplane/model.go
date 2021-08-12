package freeplane

type Map struct {
	Node struct {
		Text  string `xml:"TEXT,attr"`
		Nodes []Node `xml:"node"`
	} `xml:"node"`
}

type Node struct {
	Text  string `xml:"TEXT,attr"`
	Link  string `xml:"LINK,attr"`
	Nodes []Node `xml:"node"`
}
