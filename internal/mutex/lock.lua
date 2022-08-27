-- KEYS[1] 锁名
-- ARGV[1] 协程唯一标识：客户端标识+协程ID
-- ARGV[2] 过期时间
if (redis.call('exists',KEYS[1]) == 0) then
    redis.call('set',KEYS[1],ARGV[1]);
    redis.call('pexpire',KEYS[1],ARGV[2])
    return nil
end
return redis.call('pttl',KEYS[1])