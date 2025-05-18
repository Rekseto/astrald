package objects

import (
	"encoding/json"
	"github.com/cryptopunkscc/astrald/astral"
	"io"
)

type SearchResult struct {
	ObjectID *astral.ObjectID
}

// astral

func (SearchResult) ObjectType() string { return "mod.objects.search_result" }

func (sr SearchResult) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Struct(sr).WriteTo(w)
}

func (sr *SearchResult) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Struct(sr).ReadFrom(r)
}

// json

func (sr SearchResult) MarshalJSON() ([]byte, error) {
	type alias SearchResult
	return json.Marshal(alias(sr))
}

func (sr *SearchResult) UnmarshalJSON(bytes []byte) error {
	type alias SearchResult
	return json.Unmarshal(bytes, (*alias)(sr))
}

// text

func (sr SearchResult) MarshalText() (text []byte, err error) {
	return sr.ObjectID.MarshalText()
}

func (sr *SearchResult) UnmarshalText(text []byte) error {
	sr.ObjectID = &astral.ObjectID{}
	return sr.ObjectID.UnmarshalText(text)
}

// other

func (sr SearchResult) String() string {
	return sr.ObjectID.String()
}

func init() {
	_ = astral.DefaultBlueprints.Add(&SearchResult{})
}
