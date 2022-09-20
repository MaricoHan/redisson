-- KEYS[1] 锁名
-- ARGV[1] 过期时间
-- ARGV[2] 客户端名

if (redis.call('exists',KEYS[1]) == 0) then
    redis.call('hincrby',KEYS[1],ARGV[2],1);
    redis.call('pexpire',KEYS[1],ARGV[1]);
    return nil;
end
if (redis.call('hexist',KEYS[1],ARGV[2])==1) then
    redis.call('hincrby',KEYS[1],ARGV[2],1);
    redis.call('pexpire',KEYS[1],ARGV[1]);
    return nil;
end
return redis.call('pttl',KEYS[1]);

