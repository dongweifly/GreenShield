#!/usr/bin/env bash
PROG_NAME=$0
SERVICE_NAME=$1
ACTION=$2
ENV=$3

BIN_NOHUP=$(which nohup)

APP_HOME="/opt/green_shield"
APP_PORT=8088

HEALTH_CHECK_URL=http://127.0.0.1:${APP_PORT}/actuator/health
PID_FILE=${APP_HOME}/pid    # 应用的pid会保存到这个文件中
APP_START_TIMEOUT=10        # 健康检查的次数

usage() {
    echo "Usage: $PROG_NAME green_shield {start|stop|online|offline|restart} [dev|pre|prod]"
    exit 2
}

function health_check() {
    exptime=0
    echo "checking ${HEALTH_CHECK_URL}"
    while true
    do
        status_code=`/usr/bin/curl -L -o /dev/null --connect-timeout 5 -s -w %{http_code}  ${HEALTH_CHECK_URL}`
        if [[ x${status_code} != x200 ]];then
            sleep 1
            ((exptime++))
            echo -n -e "\rWait app to pass health check: $exptime..."
        else
            echo -n -e "check ${HEALTH_CHECK_URL} success\n"

            exit 0
            # 健康检查成功，PID写入到PID文件中
            echo $! > ${PID_FILE}
        fi
        if [[ ${exptime} -gt ${APP_START_TIMEOUT} ]]; then
            echo
            echo 'app start failed'
            exit 1
        fi
    done
}

function stop() {

        if [[ -f "${PID_FILE}" ]]; then
            echo "kill -9"
            /bin/kill -9 `cat ${PID_FILE}`
            rm ${PID_FILE}
        else
            echo "killall"
            /bin/killall green_shield
#            echo "pid file ${PID_FILE} does not exist, do noting"
        fi
}

function start() {
  case "$ENV" in
  dev)
    (set -x; ${BIN_NOHUP} ${APP_HOME}/green_shield -c ${APP_HOME}/conf/conf_dev.toml -s ${APP_HOME}/sensitive-words > /dev/null 2>&1&)
    ;;
  pre)
    (set -x; ${BIN_NOHUP} ${APP_HOME}/green_shield -c ${APP_HOME}/conf/conf_pre.toml -s ${APP_HOME}/sensitive-words > /dev/null 2>&1 &)
    ;;
  prod)
    (set -x; ${BIN_NOHUP} ${APP_HOME}/green_shield -c ${APP_HOME}/conf/conf_prod.toml -s ${APP_HOME}/sensitive-words > /dev/null 2>&1 &)
    ;;
  *)
    usage
    ;;
  esac

  health_check
}

case "$ACTION" in
    start)
        start
    ;;
    stop)
        stop
    ;;
    offline)
        offline
    ;;
    restart)
        stop
        start
    ;;
    *)
        usage
    ;;
esac

exit 0


