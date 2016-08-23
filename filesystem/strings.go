package filesystem

func SplitPath(p Path) []Path {
	paths := make([]Path, 0, 32)
	lastStart := 0
	for i, c := range p {
		if c == '/' {
			if i != 0 {
				paths = append(paths, p[lastStart:i])
			}
			lastStart = i + 1 // ignore the slash
		}
	}

	// Don't add a trailing slash, but otherwise add the last component
	if lastStart != len(p) {
		paths = append(paths, p[lastStart:])
	}
	return paths
}
