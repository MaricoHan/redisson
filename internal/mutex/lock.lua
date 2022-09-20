if (redis.call('exists',KEYS[1]) == 0) then
    redis.call('set',KEYS[1],ARGV[1]);
    redis.call('pexpire',KEYS[1],ARGV[2])
    return nil
end
return redis.call('pttl',KEYS[1])

