package api

import (
	"hash/fnv"
)

type Tag uint32

const (
	CERN = iota
	Experimental
	Theory
	Stats
)

var TagLabels = map[Tag]string{
	CERN:         "CERN",
	Experimental: "Experimental",
	Theory:       "Theory",
	Stats:        "Stats",
}

type Jargon struct {
	LabelId     uint32
	Label       string
	Title       string
	Description string
	Tags        []Tag
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func MakeJargon(label string, title string, description string, tags []Tag) Jargon {
	hashedLabel := hash(label)

	return Jargon{LabelId: hashedLabel, Label: label, Title: title, Description: description, Tags: tags}
}