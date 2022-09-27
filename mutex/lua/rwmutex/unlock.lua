-- KEYS[1] 锁名
-- KEYS[2] 发布订阅的channel
-- ARGV[1] 协程唯一标识：客户端标识+协程ID
-- ARGV[2] 解锁时发布的消息
local t = redis.call('type',KEYS[1])["ok"]
if  t == "hash" then
    if redis.call('hexists',KEYS[1],ARGV[1]) == 0 then
        return 0
    end
    if redis.call('hincrby',KEYS[1],ARGV[1],-1) == 0 then
        redis.call('hdel',KEYS[1],ARGV[1])
        if (redis.call('hlen',KEYS[1]) > 0 )then
            return 2
        end
        redis.call('del',KEYS[1])
        redis.call('publish',KEYS[2],ARGV[2])
        return 1
    else
        return 1
    end
elseif t == "none" then
        redis.call('publish',KEYS[2],ARGV[2])
        return 1
elseif redis.call('get',KEYS[1]) == ARGV[1] then
        redis.call('del',KEYS[1])
        redis.call('publish',KEYS[2],ARGV[2])
        return 1
else
    return 0
end