package common

type Member struct {
	Name string
}

func CompareMembers(currentMembers, desiredMembers []Member) (missing, extra []Member) {
  // Vytvorime mapy pro rychle vyhledavani
	currentSet := make(map[string]struct{})
	desiredSet := make(map[string]struct{})

	// Naplnime mapy
	for _, m := range currentMembers {
		currentSet[m.Name] = struct{}{}
	}
	for _, m := range desiredMembers {
		desiredSet[m.Name] = struct{}{}
	}
  
	// Najdeme chybejici cleny
	for _, m := range desiredMembers {
		if _, exists := currentSet[m.Name]; !exists {
			missing = append(missing, m)
		}
	}

	// Najdeme prebyvajici cleny
	for _, m := range currentMembers {
		if _, exists := desiredSet[m.Name]; !exists {
			extra = append(extra, m)
		}
	}

  return missing, extra
}
