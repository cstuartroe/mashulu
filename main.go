package main

import (
	"encoding/json"
	"fmt"

	"github.com/cstuartroe/mashulu/src/parse"
	"github.com/cstuartroe/mashulu/src/segment"
)

func main() {
	segments, err := segment.SegmentAndValidate("wa sela te wa sela te yo")
	if err != nil {
		fmt.Println(err)
	}
	parser := parse.New(segments)
	tree, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}
	blob, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(blob))
}
