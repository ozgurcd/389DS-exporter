# Prometheus exporter for 389-DS LDAP Server 

- https://directory.fedoraproject.org/
- https://directory.fedoraproject.org/docs/389ds/development/performance-diagnostic.html
- https://ltb-project.org/documentation/check_ldap_monitor_389ds.html
- https://access.redhat.com/documentation/en-us/red_hat_directory_server/11/html-single/performance_tuning_guide/index
- https://docs.oracle.com/cd/E19729-01/816-5597-10/dsstats.htm
- https://docs.oracle.com/cd/E20295_01/html/821-1222/baail.html#SUNWDSEEREFbaail

# The directory server 389 export there metric on ldap cn=monitor base (and generaly are accesible anonymously):
```
# ldapsearch -H ldap://localhost:389 -x -s sub -b "cn=monitor" "(objectclass=*)"
# extended LDIF
#
# LDAPv3
# base <cn=monitor> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#

# monitor
dn: cn=monitor
objectClass: top
objectClass: extensibleObject
cn: monitor
version: 389-Directory/1.3
threads: 24
currentconnections: 311
totalconnections: 423687
currentconnectionsatmaxthreads: 0
maxthreadsperconnhits: 0
dtablesize: 4096
readwaiters: 0
opsinitiated: 2661473
opscompleted: 2661472
entriessent: 2090672
bytessent: 2120635985
currenttime: 20220919133733Z
starttime: 20220918211529Z
nbackends: 1
backendmonitordn: cn=monitor,cn=userRoot,cn=ldbm database,cn=plugins,cn=config

# counters, monitor
dn: cn=counters,cn=monitor
objectClass: top
objectClass: extensibleObject
cn: counters

# snmp, monitor
dn: cn=snmp,cn=monitor
objectClass: top
objectClass: extensibleObject
cn: snmp
anonymousbinds: 64
unauthbinds: 66
simpleauthbinds: 370138
strongauthbinds: 0
bindsecurityerrors: 1267
inops: 2661473
readops: 0
compareops: 18
addentryops: 185
removeentryops: 39
modifyentryops: 333086
modifyrdnops: 0
listops: 0
searchops: 1457101
onelevelsearchops: 777376
wholesubtreesearchops: 582662
referrals: 0
chainings: 0
securityerrors: 1011
errors: 1509
connections: 311
connectionseq: 423687
connectionsinmaxthreads: 0
connectionsmaxthreadscount: 0
bytesrecv: 0
bytessent: 2120635985
entriesreturned: 2090672
referralsreturned: 0
masterentries: 0
copyentries: 0
cacheentries: 0
cachehits: 0
slavehits: 0

# search result
search: 2
result: 0 Success

# numResponses: 4
# numEntries: 3
```

This exporter request ldap cn=Monitor tree to export the metric in prometheus format.

# To build the exporter:
```
go build
```

# Exporter usage 
```
usage: 389DS-exporter [<flags>]

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9313"
                             Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                             Path under which to expose metrics.
      --ldap.ServerFQDN="localhost"
                             FQDN of the target LDAP server
      --ldap.ServerPort=389  Port to connect on LDAP server
      --version              Show application version.

```

By default the exporter listen on `http://0.0.0.0:9313/metrics`.

# Start as systemd service

Copy 389DS-exporter to /usr/local/bin.

Modify 389DS-exporter flags in 389DS-exporter.service.

Enable and Start the exporter.
```
# systemctl enable $PWD/389DS-exporter.service
# systemctl start 389DS-exporter.service
```

