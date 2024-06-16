if redis.call('GET',KEYS[1]) == ARGV[1] then
    -- 是自己的锁
    return redis.call('EXPIRE',KEYS[1],ARGV[2])
else
    -- 不是自己的锁，或者没有持有锁
    return 0
end