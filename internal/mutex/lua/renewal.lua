-- KEYS[1] 锁名
-- ARGV[1] 过期时间
return redis.call('pexpire',KEYS[1],ARGV[1])