// Copyright 2015 Alex Browne.  All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

// File scripts.go contains code related to parsing
// lua scripts in the scripts file.

// This file has been automatically generated by go generate,
// which calls scripts/generate.go. Do not edit it directly!

package jobs

import (
	"github.com/garyburd/redigo/redis"
)

var (
	
	addJobToSetScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- add_job_to_set represents a lua script that takes the following arguments:
-- 	1) The id of the job
--    2) The name of a sorted set
--    3) The score the inserted job should have in the sorted set
-- It first checks if the job exists in the database (has not been destroyed)
-- and then adds it to the sorted set with the given score.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

local jobId = ARGV[1]
local setName = ARGV[2]
local score = ARGV[3]
local jobKey = 'jobs:' .. jobId
-- Make sure the job hasn't already been destroyed
local exists = redis.call('EXISTS', jobKey)
if exists ~= 1 then
	return
end
redis.call('ZADD', setName, score, jobId)`)
	destroyJobScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- destroy_job is a lua script that takes the following arguments:
-- 	1) The id of the job to destroy
-- It then removes all traces of the job in the database by doing the following:
-- 	1) Removes the job from the status set (which it determines with an HGET call)
-- 	2) Removes the job from the time index
-- 	3) Removes the main hash for the job

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign args to variables for easy reference
local jobId = ARGV[1]
local jobKey = 'jobs:' .. jobId
-- Remove the job from the status set
local status = redis.call('HGET', jobKey, 'status')
if status ~= '' then
	local statusSet = 'jobs:' .. status
	redis.call('ZREM', statusSet, jobId)
end
-- Remove the job from the time index
redis.call('ZREM', 'jobs:time', jobId)
-- Remove the main hash for the job
redis.call('DEL', jobKey)`)
	getJobsByIdsScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- get_jobs_by_ids is a lua script that takes the following arguments:
-- 	1) The key of a sorted set of some job ids
-- The script then gets all the data for those job ids from their respective
-- hashes in the database. It returns an array of arrays where each element
-- contains the fields for a particular job, and the jobs are sorted by
-- priority.
-- Here's an example response:
-- [
-- 	[
-- 		"id", "afj9afjpa30",
-- 		"data", [34, 67, 34, 23, 56, 67, 78, 79],
-- 		"type", "emailJob",
-- 		"time", 1234567,
-- 		"freq", 0,
-- 		"priority", 100,
-- 		"retries", 0,
-- 		"status", "executing",
-- 		"started", 0,
-- 		"finished", 0,
-- 	],
-- 	[
-- 		"id", "E8v2ovkdaIw",
-- 		"data", [46, 43, 12, 08, 34, 45, 57, 43],
-- 		"type", "emailJob",
-- 		"time", 1234568,
-- 		"freq", 0,
-- 		"priority", 95,
-- 		"retries", 0,
-- 		"status", "executing",
-- 		"started", 0,
-- 		"finished", 0,
-- 	]
-- ]

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign keys to variables for easy access
local setKey = ARGV[1]
-- Get all the ids from the set name
local jobIds = redis.call('ZREVRANGE', setKey, 0, -1)
local allJobs = {}
if #jobIds > 0 then
	-- Iterate over the ids and find each job
	for i, jobId in ipairs(jobIds) do
		local jobKey = 'jobs:' .. jobId
		local jobFields = redis.call('HGETALL', jobKey)
		-- Add the id itself to the fields
		jobFields[#jobFields+1] = 'id'
		jobFields[#jobFields+1] = jobId
		-- Add the field values to allJobs
		allJobs[#allJobs+1] = jobFields
	end
end
return allJobs`)
	popNextJobsScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- pop_next_jobs is a lua script that takes the following arguments:
-- 	1) The maximum number of jobs to pop and return
-- 	2) The current unix time UTC with nanosecond precision
-- The script gets the next available jobs from the queued set which are
-- ready based on their time parameter. Then it adds those jobs to the
-- executing set, sets their status to executing, and removes them from the
-- queued set. It returns an array of arrays where each element contains the
-- fields for a particular job, and the jobs are sorted by priority.
-- Here's an example response:
-- [
-- 	[
-- 		"id", "afj9afjpa30",
-- 		"data", [34, 67, 34, 23, 56, 67, 78, 79],
-- 		"type", "emailJob",
-- 		"time", 1234567,
-- 		"freq", 0,
-- 		"priority", 100,
-- 		"retries", 0,
-- 		"status", "executing",
-- 		"started", 0,
-- 		"finished", 0,
-- 	],
-- 	[
-- 		"id", "E8v2ovkdaIw",
-- 		"data", [46, 43, 12, 08, 34, 45, 57, 43],
-- 		"type", "emailJob",
-- 		"time", 1234568,
-- 		"freq", 0,
-- 		"priority", 95,
-- 		"retries", 0,
-- 		"status", "executing",
-- 		"started", 0,
-- 		"finished", 0,
-- 	]
-- ]

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign args to variables for easy reference
local n = ARGV[1]
local currentTime = ARGV[2]
local poolId = ARGV[3]
-- Copy the time index set to a new temporary set
redis.call('ZUNIONSTORE', 'jobs:temp', 1, 'jobs:time')
-- Trim the new temporary set we just created to leave only the jobs which have a time
-- parameter in the past
redis.call('ZREMRANGEBYSCORE', 'jobs:temp', currentTime, '+inf')
-- Intersect the jobs which are ready based on their time with those in the
-- queued set. Use the weights parameter to set the scores entirely based on the
-- queued set, effectively sorting the jobs by priority. Store the results in the
-- temporary set.
redis.call('ZINTERSTORE', 'jobs:temp', 2, 'jobs:queued', 'jobs:temp', 'WEIGHTS', 1, 0)
-- Trim the temp set, so it contains only the first n jobs ordered by
-- priority
redis.call('ZREMRANGEBYRANK', 'jobs:temp', 0, -n - 1)
-- Get all job ids from the temp set
local jobIds = redis.call('ZREVRANGE', 'jobs:temp', 0, -1)
local allJobs = {}
if #jobIds > 0 then
	-- Add job ids to the executing set
	redis.call('ZUNIONSTORE', 'jobs:executing', 2, 'jobs:executing', 'jobs:temp')
	-- Now we are ready to construct our response.
	for i, jobId in ipairs(jobIds) do
		local jobKey = 'jobs:' .. jobId
		-- Remove the job from the queued set
		redis.call('ZREM', 'jobs:queued', jobId)
		-- Set the poolId field for the job
		redis.call('HSET', jobKey, 'poolId', poolId)
		-- Set the job status to executing
		redis.call('HSET', jobKey, 'status', 'executing')
		-- Get the fields from its main hash
		local jobFields = redis.call('HGETALL', jobKey)
		-- Add the id itself to the fields
		jobFields[#jobFields+1] = 'id'
		jobFields[#jobFields+1] = jobId
		-- Add the field values to allJobs
		allJobs[#allJobs+1] = jobFields
	end
end
-- Delete the temporary set
redis.call('DEL', 'jobs:temp')
-- Return all the fields for all the jobs
return allJobs`)
	purgeStalePoolScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- purge_stale_pool is a lua script which takes the following arguments:
-- 	1) The id of the stale pool to purge
-- It then does the following:
-- 	1) Removes the pool id from the set of active pools
-- 	2) Iterates through each job in the executing set and finds any jobs which
-- 		have a poolId field equal to the id of the stale pool
-- 	3) If it finds any such jobs, it removes them from the executing set and
-- 		adds them to the queued so that they will be retried

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign args to variables for easy reference
local stalePoolId = ARGV[1]
-- Check if the stale pool is in the set of active pools first
local isActive = redis.call('SISMEMBER', 'pools:active', stalePoolId)
if isActive then
	-- Remove the stale pool from the set of active pools
	redis.call('SREM', 'pools:active', stalePoolId)
	-- Get all the jobs in the executing set
	local jobIds = redis.call('ZRANGE', 'jobs:executing', 0, -1)
	for i, jobId in ipairs(jobIds) do
		local jobKey = 'jobs:' .. jobId
		-- Check the poolId field
		-- If the poolId is equal to the stale id, then this job is stuck
		-- in the executing set even though no worker is actually executing it
		local poolId = redis.call('HGET', jobKey, 'poolId')
		if poolId == stalePoolId then
			local jobPriority = redis.call('HGET', jobKey, 'priority')
			-- Move the job into the queued set
			redis.call('ZADD', 'jobs:queued', jobPriority, jobId)
			-- Remove the job from the executing set
			redis.call('ZREM', 'jobs:executing', jobId)
			-- Set the job status to queued and the pool id to blank
			redis.call('HMSET', jobKey, 'status', 'queued', 'poolId', '')
		end
	end
end
`)
	retryOrFailJobScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- retry_or_fail_job represents a lua script that takes the following arguments:
-- 	1) The id of the job to either retry or fail
-- It first checks if the job has any retries remaining. If it does,
-- then it:
-- 	1) Decrements the number of retries for the given job
-- 	2) Adds the job to the queued set
-- 	3) Removes the job from the executing set
-- 	4) Returns true
-- If the job has no retries remaining then it:
-- 	1) Adds the job to the failed set
-- 	3) Removes the job from the executing set
-- 	2) Returns false

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

-- Assign args to variables for easy reference
local jobId = ARGV[1]
local jobKey = 'jobs:' .. jobId
-- Make sure the job hasn't already been destroyed
local exists = redis.call('EXISTS', jobKey)
if exists ~= 1 then
	return 0
end
-- Check how many retries remain
local retries = redis.call('HGET', jobKey, 'retries')
local newStatus = ''
if retries == '0' then
	-- newStatus should be failed because there are no retries left
	newStatus = 'failed'
else
	-- subtract 1 from the remaining retries
	redis.call('HINCRBY', jobKey, 'retries', -1)
	-- newStatus should be queued, so the job will be retried
	newStatus = 'queued'
end
-- Get the job priority (used as score)
local jobPriority = redis.call('HGET', jobKey, 'priority')
-- Add the job to the appropriate new set
local newStatusSet = 'jobs:' .. newStatus
redis.call('ZADD', newStatusSet, jobPriority, jobId)	
-- Remove the job from the old status set
local oldStatus = redis.call('HGET', jobKey, 'status')
if ((oldStatus ~= '') and (oldStatus ~= newStatus)) then
	local oldStatusSet = 'jobs:' .. oldStatus
	redis.call('ZREM', oldStatusSet, jobId)
end
-- Set the job status in the hash
redis.call('HSET', jobKey, 'status', newStatus)
if retries == '0' then
	-- Return false to indicate the job has not been queued for retry
	-- NOTE: 0 is used to represent false because apparently
	-- false gets converted to nil
	return 0
else
	-- Return true to indicate the job has been queued for retry
	-- NOTE: 1 is used to represent true (for consistency)
	return 1
end`)
	setJobFieldScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- set_job_field represents a lua script that takes the following arguments:
-- 	1) The id of the job
--    2) The name of the field
--    3) The value to set the field to
-- It first checks if the job exists in the database (has not been destroyed)
-- and then sets the given field to the given value.

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go

local jobId = ARGV[1]
local fieldName = ARGV[2]
local fieldVal = ARGV[3]
local jobKey = 'jobs:' .. jobId
-- Make sure the job hasn't already been destroyed
local exists = redis.call('EXISTS', jobKey)
if exists ~= 1 then
	return
end
redis.call('HSET', jobKey, fieldName, fieldVal)`)
	setJobStatusScript = redis.NewScript(0, `-- Copyright 2015 Alex Browne.  All rights reserved.
-- Use of this source code is governed by the MIT
-- license, which can be found in the LICENSE file.

-- set_job_status is a lua script that takes the following arguments:
-- 	1) The id of the job
-- 	2) The new status (e.g. "queued")
-- It then does the following:
-- 	1) Adds the job to the new status set
-- 	2) Removes the job from the old status set (which it gets with an HGET call)
-- 	3) Sets the 'status' field in the main hash for the job

-- IMPORTANT: If you edit this file, you must run go generate . to rewrite ../scripts.go
	
-- Assign args to variables for easy reference
local jobId = ARGV[1]
local newStatus = ARGV[2]
local jobKey = 'jobs:' .. jobId
-- Make sure the job hasn't already been destroyed
local exists = redis.call('EXISTS', jobKey)
if exists ~= 1 then
	return
end
local newStatusSet = 'jobs:' .. newStatus
-- Add the job to the new status set
local jobPriority = redis.call('HGET', jobKey, 'priority')
redis.call('ZADD', newStatusSet, jobPriority, jobId)
-- Remove the job from the old status set
local oldStatus = redis.call('HGET', jobKey, 'status')
if ((oldStatus ~= '') and (oldStatus ~= newStatus)) then
	local oldStatusSet = 'jobs:' .. oldStatus
	redis.call('ZREM', oldStatusSet, jobId)
end
-- Set the status field
redis.call('HSET', jobKey, 'status', newStatus)`)
)
