-- KEYS[1] 锁名
-- ARGV[1] 过期时间
-- ARGV[2] 客户端协程唯一标识
local t = redis.call('type',KEYS[1])["ok"]
if t =="string" then
    if redis.call('get',KEYS[1])==ARGV[2] then
        return redis.call('pexpire',KEYS[1],ARGV[1])
    end
    return 0
elseif t == "hash" then
    if redis.call('hexists',KEYS[1],ARGV[2])==0 then
        return 0
    end
    return redis.call('pexpire',KEYS[1],ARGV[1])
else
    return 0
end
