package extloadtest

import "math/rand"

func getRandomIndices(count2Change int, maxLen int) []int {
  indiciesToChange := make([]int, 0, count2Change)
  for len(indiciesToChange) < count2Change {
    index := rand.Intn(maxLen)
    if !contains(indiciesToChange, index) {
      indiciesToChange = append(indiciesToChange, index)
    }
  }
  return indiciesToChange
}

func contains(s []int, e int) bool {
  for _, a := range s {
    if a == e {
      return true
    }
  }
  return false
}
