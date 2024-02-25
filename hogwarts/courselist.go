//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	learned := make(map[string]bool)
	entered := make(map[string]bool)
	result := make([]string, 0)

	for course := range prereqs {
		var impl func(string)
		impl = func(course string) {
			if entered[course] {
				panic("Already entered")
			}

			if learned[course] {
				return
			}

			entered[course] = true

			for _, prereq := range prereqs[course] {
				impl(prereq)
			}

			learned[course] = true
			entered[course] = false

			result = append(result, course)
		}

		impl(course)
	}

	return result
}
