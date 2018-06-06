RSP=`curl -s --connect-timeout 5 -m 5  "http://127.0.0.1:8002/moments/?type=1" |grep code -c`
if (($RSP == 0)); then
     ./restart.sh
fi
