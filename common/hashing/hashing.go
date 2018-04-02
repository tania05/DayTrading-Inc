package hashing

import "hash/fnv"

func ModuloHash(userId string, max int) int {
	digest := fnv.New32a()
	digest.Write([]byte(userId))
	hash := digest.Sum32()
	hashInt := int(hash)
	if hashInt < 0 {
		hashInt = hashInt * -1
	}
	index := hashInt % max
	return index
}

