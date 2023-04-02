baseURL="curl -s -u Administrator:password -X POST http://localhost:8091%s"
initializeNodeCommand="/nodes/self/controller/settings -d path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata -d index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fadata -d cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fedata -d eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fidata"

startServicesCommand="/node/controller/setupServices -d services=kv%2Cindex%2Cn1ql%2Cfts"

setBucketMemoryQuotaCommand="/pools/default -d memoryQuota=256 -d indexMemoryQuota=256 -d ftsMemoryQuota=256"

createAdminUserCommand="/settings/web -d port=8091 -d username=admin -d password=password"

newBaseURL="curl -s -u admin:password -X POST http://localhost:8091%s"
createPostAPIBucketCommand="/pools/default/buckets -d flushEnabled=1 -d name=post -d ramQuotaMB=100 -d replicaNumber=0 -d evictionPolicy=fullEviction -d bucketType=couchbase"
createTestBucketCommand="/pools/default/buckets -d flushEnabled=1 -d name=test -d ramQuotaMB=100 -d replicaNumber=0 -d evictionPolicy=fullEviction -d bucketType=couchbase"

setIndexStorageSetting="/settings/indexes -d indexerThreads=0 -d logLevel=info -d maxRollbackPoints=5 -d memorySnapshotInterval=200 -d stableSnapshotInterval=5000 -d storageMode=forestdb"

runCurlCommand() {
    local url=$1
    local command=$2
    local fullURL="$(printf "$url" "$command")"
    eval "$fullURL"
}

runCurlCommand "$baseURL" "$initializeNodeCommand"
runCurlCommand "$baseURL" "$startServicesCommand"
runCurlCommand "$baseURL" "$setBucketMemoryQuotaCommand"
runCurlCommand "$baseURL" "$createAdminUserCommand"
runCurlCommand "$newBaseURL" "$createPostAPIBucketCommand"
runCurlCommand "$newBaseURL" "$createTestBucketCommand"
runCurlCommand "$newBaseURL" "$setIndexStorageSetting"