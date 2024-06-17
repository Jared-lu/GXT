local val = redis.call('GET', KEYS[1])
-- key存在，redis返回nil回复，对应的lua类型取值为false
if val == false then
    -- 没有加锁, 成功返回 OK
    redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2])
    return 1
elseif val == ARGV[1] then
    -- 上次加锁成功，重新设置过期时间，设置成功返回1，失败返回0（这个发生的概率很小）
    return redis.call('EXPIRE',KEYS[1],ARGV[2])
else
    -- 锁被人拿着
    return 2
end