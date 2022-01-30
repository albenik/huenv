package reflector

func SortWithDeps(keys []string, envs map[string]*Target) []string {
	deps := collectDeps(envs)

	start := 0
loop:
	for {
		for i := start; i < len(keys); i++ {
			k := keys[i]
			e := envs[k]
			depKeys, ok := deps[e.Field.Name]
			if !ok {
				continue
			}

			depIdx := i
			for _, dk := range depKeys {
				for j, kk := range keys {
					if kk != dk {
						continue
					}
					if j <= depIdx {
						keys = strsliceMove(keys, j, depIdx)
						depIdx = j
					}
					break
				}
			}

			start = i + 1
			continue loop
		}
		break loop
	}

	return keys
}

func collectDeps(envs map[string]*Target) map[string][]string {
	deps := make(map[string][]string)

	for k, v := range envs {
		if cond := reqIfCond(v.Condition); cond != nil {
			deps[cond.Target.Name] = append(deps[cond.Target.Name], k)
		}
	}

	return deps
}

func reqIfCond(i interface{}) *ConditionRequireIf {
	switch c := i.(type) {
	case *ConditionRequireIf:
		return c
	case *ConditionRequireIfCombined:
		return c.First
	default:
		return nil
	}
}

func strsliceInsert(array []string, idx int, value string) []string {
	return append(array[:idx], append([]string{value}, array[idx:]...)...)
}

func strsliceRemove(array []string, idx int) []string {
	return append(array[:idx], array[idx+1:]...)
}

func strsliceMove(array []string, idx int, from int) []string {
	val := array[from]
	return strsliceInsert(strsliceRemove(array, from), idx, val)
}
