package online

import "context"

const onlineUserKey = "online_users"

func AddUserToRedis(username string) error {
	res := RedisClient.SAdd(context.Background(), onlineUserKey, username)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func DeleteUserFromRedis(username string) error {
	res := RedisClient.SRem(context.Background(), onlineUserKey, username)

	return res.Err()
}

func IsUserOnline(username string) bool {
	res := RedisClient.SIsMember(context.Background(), onlineUserKey, username)
	return res.Val()
}

func GetTotalOnlineUsers() int {
	return len(RedisClient.SMembers(context.Background(), onlineUserKey).Val())
}
