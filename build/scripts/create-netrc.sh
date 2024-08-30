#!/usr/bin/env bash
PASSWORDFILE=$1
MACHINE_HOSTNAME=$2
echo "Creating password file $PASSWORDFILE for Host $MACHINE_HOSTNAME"
echo -n "machine ${MACHINE_HOSTNAME} login " > $PASSWORDFILE
if [ -f /username ];then
    cat /username >> $PASSWORDFILE
    echo -n " password " >> $PASSWORDFILE
    cat /password >> $PASSWORDFILE
    rm /username /password
else
    cat /pipeline-secrets/username >> $PASSWORDFILE
    echo -n " password " >> $PASSWORDFILE
    cat /pipeline-secrets/password >> $PASSWORDFILE
    rm /pipeline-secrets/username /pipeline-secrets/password
fi
if [ -f .arti-machine ]; then
    MACHINE_HOSTNAME=$(cat .arti-machine)
    echo "Adding artifactory secret to file $PASSWORDFILE for Host $MACHINE_HOSTNAME"
    echo "" >> $PASSWORDFILE
    echo "" >> $PASSWORDFILE
    echo -n "machine ${MACHINE_HOSTNAME} login " >> $PASSWORDFILE
    cat .arti-username >> $PASSWORDFILE
    echo -n " password " >> $PASSWORDFILE
    cat .arti-password >> $PASSWORDFILE
    rm .arti-username .arti-password .arti-machine
fi
sed -i -e 's/^\(.*\)\n/\1/g' $PASSWORDFILE
chmod 600 $PASSWORDFILE
