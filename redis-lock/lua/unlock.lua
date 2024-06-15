-- 检查是不是自己的锁
-- 是，就删除锁
-- 以上两个步骤要做成原子操作，因此需要使用lua脚本来实现
if redis.call('GET',KEYS[1]) == ARGV[1] then
    -- 是自己的锁
    return redis.call('del',KEYS[1])
else
    -- 不是自己的锁，或者没有持有锁
    return 0
end