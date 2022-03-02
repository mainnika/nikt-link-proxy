package data

// Meta holds meta data for Link
type Meta struct {
	Ref string
	SID string
}

// Apply applies meta data to input and returns it as output
func (m Meta) Apply(in map[string][]string) (out map[string][]string) {

	if in == nil {
		in = map[string][]string{}
	}

	out = in

	if m.Ref != "" {
		out["ref"] = append(out["ref"], m.Ref)
	}
	if m.SID != "" {
		out["sid"] = append(out["sid"], m.SID)
	}

	return
}
