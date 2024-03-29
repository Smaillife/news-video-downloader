#!/bin/sh
#
# chkconfig: 2345  80 50
# description: myrandservice is for testing how to write service in Linux 
#              
# processname: myrandservice
# 

ret=0
PROC=cctvbnews
BINPATH="${DIR}/bin/"${PROC}


start() {
    # check fdb status
                echo "start ${PROC}..."
		start-stop-daemon --start --background --exec ${BINPATH}
                ret=$?
} 

stop() {
    echo "stop ${PROC}..."
		start-stop-daemon --stop --name ${PROC}
                ret=$?
} 

status() {
                local result
                echo "check status of ${PROC}..."
                #lines=$( ps -ef | grep ${PROC} | grep -v grep  |  )
                #echo $lines
                result=$( ps -ef | grep ${PROC} | grep -v grep | wc -l )
                #echo $result
                if [ $result -gt 0 ] ; then
                                echo "${PROC} is up"
                                ret=0
                else
                                echo "${PROC} is down"
                                ret=1
                fi
                echo "check status of ${PROC}...done."
} 

# See how we were called.
case "$1" in
  start)
        stop
        start
        ;;
  stop)
        stop
        ;;
  restart)
        stop
	start
        ;;
  status)
        status 
        ;;
  *)
        echo $"Usage: $0 {start|stop|status|restart}"
        exit 1
esac

if [ $ret -eq 0 ];then
    echo  "SUCCESS"
else
    echo  "FAIL"
fi

exit $ret
