package helpers

import "go.mongodb.org/mongo-driver/bson"

// NotDeletedFilter returns a BSON filter that matches documents where is_deleted is false or doesn't exist
// This handles cases where older documents don't have the is_deleted field
func NotDeletedFilter() bson.M {
	return bson.M{
		"$or": []bson.M{
			{"is_deleted": false},
			{"is_deleted": bson.M{"$exists": false}},
		},
	}
}

// AddNotDeletedFilter adds the not-deleted filter to an existing filter map
func AddNotDeletedFilter(filter bson.M) bson.M {
	if len(filter) == 0 {
		return NotDeletedFilter()
	}

	return bson.M{
		"$and": []bson.M{
			filter,
			NotDeletedFilter(),
		},
	}
}
