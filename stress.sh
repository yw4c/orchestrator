iso8601=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
body=body-$iso8601.json

echo 'GET http://35.244.227.118/booking/sync' | \
    vegeta attack -rate 1000 -duration 20s -header=""Host":"orchestrator.com.tw"" -timeout 10m | \
    tee ./reports/results-$iso8601.bin | vegeta encode | \
    jaggr @count=rps \
          hist\[100,200,300,400,500\]:code \
          p25,p50,p95:latency \
          sum:bytes_in \
          sum:bytes_out | \
    jplot rps+code.hist.100+code.hist.200+code.hist.300+code.hist.400+code.hist.500 \
          latency.p95+latency.p50+latency.p25 \
          bytes_in.sum+bytes_out.sum

plot=plot-$iso8601.html
cat ./reports/results-$iso8601.bin | vegeta report
cat ./reports/results-$iso8601.bin | vegeta plot > ./reports/$plot
rm -rf ./reports/$body