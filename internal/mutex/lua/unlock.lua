-- KEYS[1] 锁名
-- KEYS[2] 发布订阅的channel
-- ARGV[1] 协程唯一标识：客户端标识+协程ID
-- ARGV[2] 解锁时发布的消息
if redis.call('exists',KEYS[1]) == 1 then
    if (redis.call('get',KEYS[1]) == ARGV[1]) then
        redis.call('del',KEYS[1])
    else
        return 0
    end
end
redis.call('publish',KEYS[2],ARGV[2])
return 1