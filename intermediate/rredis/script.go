package rredis

import "github.com/go-redis/redis/v8"

var IncrScript = redis.NewScript(`if (redis.call('exists', KEYS[1]) == 1) then
				local stock = tonumber(redis.call('get', KEYS[1]));
				local num = tonumber(ARGV[1]);
				if (stock == -1) then
					return -1;
				end;
				if (stock + num >= 0) then
					return redis.call('incrby', KEYS[1], num);
				end;
				return -2;
			end;
			return -3;`)

var IncrTopLimitScript = redis.NewScript(`
				local stock = tonumber(redis.call('get', KEYS[1]));
				if type(stock) == "nil" then
					stock = 0
				end;
				local num = tonumber(ARGV[1]);
				local top = tonumber(ARGV[2]);
				if ( stock + num > top) then
					return 0;
				else
					return redis.call('incrby', KEYS[1], num);
				end;
			`)

var HIncrUnMinusScript = redis.NewScript(`if (redis.call('exists', KEYS[1]) == 1) then
				local stock = tonumber(redis.call('hget', KEYS[1], KEYS[2]));
				local num = tonumber(ARGV[1]);
				if (stock == -1) then
					return -1;
				end;
				if (stock + num >= 0) then
					return redis.call('HINCRBY', KEYS[1], KEYS[2], num);
				end;
				return -2;
			end;
			return -3;`)

var HIncrMinZeroScript = redis.NewScript(`if (redis.call('exists', KEYS[1]) == 1) then
				local stock = tonumber(redis.call('hget', KEYS[1], KEYS[2]));
				local num = tonumber(ARGV[1]);
				if (stock == -1) then
					return -1;
				end;
				if (stock + num >= 0) then
					return redis.call('HINCRBY', KEYS[1], KEYS[2], num);
				elseif (stock + num < 0) then
					return redis.call('HSET', KEYS[1], KEYS[2], 0)
				end;
				return -2;
			end;
			return -3;`)
