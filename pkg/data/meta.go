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

	for _, pair := range [][2]string{
		{"ref", m.Ref},
		{"sid", m.SID},
	} {
		if pair[1] != "" {
			out[pair[0]] = append(out[pair[0]], pair[1])
		}
	}

	return
}
