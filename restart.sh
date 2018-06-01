PID=`ps -ef|grep main|grep -v grep|awk '{print $2}'`
if (($PID > 0)); then
        kill $PID
fi
nohup ./main &
