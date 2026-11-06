# image-tools

## Installation

Download it from

- [image-tools  for linux](https://artifactory-pro.corporate.de:8443/ui/native/crpas-bin-dev-local/istag_and_image_management/image-tools-linux/))
- [image-tools .exe for windows](https://artifactory-pro.corporate.de:8443/ui/native/crpas-bin-dev-local/istag_and_image_management/image-tools-windows/))

and store it somewhere in your PATH. It is a statically linked go program and no installation is neccessary.

## Cluster configuration

At runtime, `clusterconfig.json` must be next to the executable, including in
server mode. Generate it in the repository root with
`make config DOMAIN=<cluster-domain>` using authenticated
OpenShift tools and `jq`, then place it next to the binary if necessary. The
file contains credentials; keep it private and do not commit it. An existing
`clusterconfig.json` in the same format can also be used. Remove any older
generated `ocrequest/config-clusters.go` before building.

Dynamic family configuration requires the `BITBUCKET_TOKEN` environment
variable. Use `-statcfg` to select the static family configuration instead.

## DESCRIPTION

image-tools  reports image or istag details for a application appgroup (eg. appgroup)
             or
             generate a bash script as output to delete istags

- For __existing Is, IsTags and Images__ it operates cluster and appgroup specific.
    That means it works for __one cluster__ and for __one family__ like 'appgroup, ...'
    The cluster must be defined by the mandatory parametter '-cluster=[clustername]'
    The appgroup must be defined by the mandatory parameter '-appgroup=[appgroup]

- For __used IsTags__ it looks in __all clusters__ and reports the istags used by any

    - deploymentconfig,
    - job,
    - cronjob or
    - pod

    of __all namespaces that belong to the application family__.

- __Generate reports__ about imagestreamtags, imagestreams and images or all together. The content, what
    should be collected from the cluster for the content of the report can be defined by __one or more__ of
    the mandatory parameter '-is', '-istag', '-image', '-used' or '-all'.

- __Variable report output format__. The output format of the report can be chosen ba on of:
  - -json,
  - -yaml,
  - -csv,
  - -table or - tabgroup (table with grouped rows for identical content).
    Output as table or tabgroup is best used when piped into less and is am

- __filter data for reports__. Define parameter -isname=..., -istagname=..., tagname=... or -shaname=...

- _(not yet implemented) delete istags based on filters like 'older than n days' and/or 'istag name like pattern' (not yet implemented)_

For this reports the data is collected from the oc cluster defined by parameter '-cluster=...' and
the parameter 'appgroup=...'. For type 'used' (also included in type 'all') from all clusters.


- __generate delete script for istags__.  Define parameter __-delete__ to switch in delete mode (generate a script as output for deleting istags)
  - further define the details with the additional parameters:
    - -snapshot
    - -nonbuild (appgroup specific)
    - -delfilter=regexp-string
    - -isname=string
    - -tagname=string
    - -istagname=string
    - -namespace=string

## Usage

    execute image-tools with parameter -h to get help and examples
### Command

    ./image-tools  -appgroup=... -cluster=... -all|-image|-is|-istag|-used [output format (default json)] [filter (default none)]

the parameters can be specified in any order

#### Parameters

##### Define source, both parameters are mandatory

    -appgroup=string    appgroup name, eg. appgroup,

    -cluster=string   shortname of clustername

##### Type of objects to collect and report. One of them is mandatory

    -all              collect and output all
                      _imageStreams_, _imageStreamTags_, _image's__
                      of the given cluster and the
                      _used-istags_ from all clusters.

    -image            collect and output of Image's

    -istag            collect and output of imageStreamTags

    -is               collect and output of imageStreams

    -used             collect and output used imageStreams imageStreamTags
                      and Image's from all clusters

#### Options

    -delete           activate the delete mode. The same data as for reports is collected from the clusters, but
                      it is used to generate a script as output for istag deleting instead a report.

    -delpattern=regexp-string    optional, can be combined with the following filter parameters

##### Filter (usable for reports and delete mode)

    -namespace=string namespace to look for istags

    -isname=string    filter output of one imageStream as json, eg. -is=wvv-service

    -istagname=string filter output of one imageStreamTag

    -shaname=string   filter output of a Image with this SHA

    -tagname=string   filter output all istags with this Tag

##### Output Format (for reports)

    -json             defines JSON as the output format for the reported data.
                      This is the DEFAULT

    -yaml             defines YAML as the output format for the reported data

    -csv              defines CSV as the output format for the reported data

    -csvfile=string   define the common filename-part for the output files in csv format.
                      For every type of openshift objects a separate file is generated
                      with the following names schema: '<common-filename>-<type>.csv'.

    -table            defines formatted ASCI TABLE as the output format for the
                      reported data

    -tabgroup         defines formatted ASCII TABLE WITH GROUPED ROWS as the
                      output format for the reported data.

##### Other Options

    -occlient         use oc client instead of api call for cluster communication

    -statcfg          us statically defined config for families instead of dynamic generated config based on files of config-tool

    -noproxy          disable the usage of a proxy for OpenShift API requests

    -socks5=string    enable socks5 usage. E.g. -socks5=localhost:65022, this is the default. To disable socks5 set -socks5=no

## EXAMPLES

Report all information for appgroup appgroup in cluster cicd as json
(which is the default output format)

    ./image-tools  -cluster=cluster-short -appgroup=app-groupü -all

or as table

        ./image-tools  -cluster=cluster-short -appgroup=appgroup -all -table

or csv in different files for each type of information

    ./image-tools  -cluster=cluster -appgroup=appgroup -all -csvfile=prefix

writes the output to different files 'prefix-type' in current directory

Report only __used__ istags for appgroup appgroup as pretty printed table
(the output is paginated to fit your screen size and piped to
    the pager define in the environment variable \$PAGER\/\%PAGER\%.
If $PAGER is not set, it try to use 'more')

        ./image-tools  -appgroup=appgroup -used -table

or json

        ./image-tools  -appgroup=appgroup -used

or yaml

    ./image-tools  -appgroup=appgroup -used -yaml

or csv

    ./image-tools  -appgroup=appgroup -used -csv

Report istags with tag=latest for appgroup appgroup in cluster cicd as yaml report

    ./image-tools  -cluster=cluster -appgroup=appgroup -istag -yaml -tagname=latest

Report ImageStreams for appgroup appgroup in cluster int as yaml report

    ./image-tools  -cluster=cluster -appgroup=appgroup -is -yaml

Report ImageStreams with name=webcode-service for appgroup appgroup in cluster cicd as table report

    ./image-tools  -cluster=cluster -appgroup=appgroup -is -isname=webcode-service -table

Delete: Generate a shell script as output to delete old istags(60 days, the default) for appgroup appgroup in cluster cicd
    and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service and zahlungsstoerung-service

        image-tools -appgroup=appgroup -cluster=cluster -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'

Delete: To use the script output to really delete the istags, you can use the following line (you must be an openshift admin):

        image-tools -appgroup=appgroup -cluster=cluster -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'|xargs -n 1 -I{} bash -c "{}"

Delete: To only generate a script to delete old snapshot istags:

        image-tools -appgroup=appgroup -cluster=cluster -delete -snapshot

Delete: To delete all not used images of appgroup 'appgroup' in cluster cicd

        image-tools -appgroup=appgroup -cluster=cluster -delete  -minage=0 -delpattern='.'

Delete: To delete all hybris istags of appgroup appgroup older than 45 days

        image-tools -appgroup=appgroup -cluster=cluster -delete -isname=hybris -minage=45

!!!  To check, which images are not deleted, because they are in use you can check the logfile, which is created in the current directory

```bash
 cat image-tools.log  |  grep logUsedIstags:  |  sort -k5  |  column -t
```
