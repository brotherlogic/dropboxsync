package main

func diffFileList(master, new []string) []string {
	newOnes := []string{}

	for _, f := range new {
		found := false
		for _, m := range master {
			if m == f {
				found = true
			}
		}

		if !found {
			newOnes = append(newOnes, f)
		}
	}

	return newOnes
}
