// Package databases provides database connection management for MongoDB and Redis.
package databases

func InitDB() (bool, error) {
	// Connect to MongoDB
	if err := ConnectMongoDB(); err != nil {
		return false, err
	}

	return true, nil
}

func CloseDB() (bool, error) {
	// Disconnect from MongoDB
	if err := DisconnectMongoDB(); err != nil {
		return false, err
	}

	return true, nil
}
