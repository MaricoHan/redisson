-- KEYS[1] 锁名
-- ARGV[1] 协程唯一标识：客户端标识+协程ID
-- ARGV[2] 过期时间
local t = redis.call('type',KEYS[1])["ok"]
if t == "string" then
    return redis.call('pttl',KEYS[1])
else
    redis.call('hincrby',KEYS[1],ARGV[1],1)
    redis.call('pexpire',KEYS[1],ARGV[2])
    return nil
end
