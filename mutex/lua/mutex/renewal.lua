-- KEYS[1] 锁名
-- ARGV[1] 过期时间
-- ARGV[2] 客户端协程唯一标识
if redis.call('get',KEYS[1])==ARGV[2] then
    return redis.call('pexpire',KEYS[1],ARGV[1])
end
return 0