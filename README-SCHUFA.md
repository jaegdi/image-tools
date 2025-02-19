# image-tools

## Installation

Download it from

- [image-tools  for linux](https://artifactory-pro.sf-rz.de:8443/ui/native/scpas-bin-dev-local/istag_and_image_management/image-tools-linux/))
- [image-tools .exe for windows](https://artifactory-pro.sf-rz.de:8443/ui/native/scpas-bin-dev-local/istag_and_image_management/image-tools-windows/))

and store it somewhere in your PATH. It is a statically linked go program and no installation is neccessary.

## DESCRIPTION

image-tools  reports image or istag details for a application family (eg. app-group, fpc, aps, ssp)
             or
             generate a bash script as output to delete istags

- For __existing Is, IsTags and Images__ it operates cluster and family specific.
    That means it works for __one cluster__ like
    'cid, int, ppr, vpt or pro' and for __one family__ like 'app-group, sps, fpc, aps, ...'
    The cluster must be defined by the mandatory parametter '-cluster=[cid|int|ppr|vpt|pro]'
    The family must be defined by the mandatory parameter '-family=[aps|fpc|app-group|ssp]

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
the parameter 'family=...'. For type 'used' (also included in type 'all') from all clusters.


- __generate delete script for istags__.  Define parameter __-delete__ to switch in delete mode (generate a script as output for deleting istags)
  - further define the details with the additional parameters:
    - -snapshot
    - -nonbuild (app-group specific)
    - -delfilter=regexp-string
    - -isname=string
    - -tagname=string
    - -istagname=string
    - -namespace=string

## Usage

    execute image-tools with parameter -h to get help and examples
### Command

    ./image-tools  -family=... -cluster=... -all|-image|-is|-istag|-used [output format (default json)] [filter (default none)]

the parameters can be specified in any order

#### Parameters

##### Define source, both parameters are mandatory

    -family=string    family name, eg. app-group, aps, ssp or fpc

    -cluster=string   shortname of cluster, eg.
      cluster,  ppr-scp0 or pro-scp0
      cluster, ppr-scp0, vpt-scp0, pro-scp0

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

Report all information for family app-group in cluster cid as json
(which is the default output format)

    ./image-tools  -cluster=cluster-short -family=app-groupü -all

or as table

        ./image-tools  -cluster=cluster-short -family=app-group -all -table

or csv in different files for each type of information

    ./image-tools  -cluster=cluster -family=app-group -all -csvfile=prefix

writes the output to different files 'prefix-type' in current directory

Report only __used__ istags for family app-group as pretty printed table
(the output is paginated to fit your screen size and piped to
    the pager define in the environment variable \$PAGER/\%PAGER\%.
If $PAGER is not set, it try to use 'more')

        ./image-tools  -family=app-group -used -table

or json

        ./image-tools  -family=app-group -used

or yaml

    ./image-tools  -family=app-group -used -yaml

or csv

    ./image-tools  -family=app-group -used -csv

Report istags with tag=latest for family app-group in cluster cid as yaml report

    ./image-tools  -cluster=cluster -family=aps -istag -yaml -tagname=latest

Report ImageStreams for family aps in cluster int as yaml report

    ./image-tools  -cluster=cluster -family=aps -is -yaml

Report ImageStreams with name=webcode-service for family app-group in cluster cid as table report

    ./image-tools  -cluster=cluster -family=app-group -is -isname=webcode-service -table

Delete: Generate a shell script as output to delete old istags(60 days, the default) for family app-group in cluster cid
    and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service and zahlungsstoerung-service

        image-tools -family=app-group -cluster=cluster -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'

Delete: To use the script output to really delete the istags, you can use the following line (you must be an openshift admin):

        image-tools -family=app-group -cluster=cluster -delete -snapshot -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'|xargs -n 1 -I{} bash -c "{}"

Delete: To only generate a script to delete old snapshot istags:

        image-tools -family=app-group -cluster=cluster -delete -snapshot

Delete: To delete all not used images of family 'aps' in cluster cid

        image-tools -family=aps -cluster=cluster -delete  -minage=0 -delpattern='.'

Delete: To delete all hybris istags of family app-group older than 45 days

        image-tools -family=app-group -cluster=cluster -delete -isname=hybris -minage=45

!!!  To check, which images are not deleted, because they are in use you can check the logfile, which is created in the current directory
```bash
 cat image-tools.log  |  grep logUsedIstags:  |  sort -k5  |  column -t
```
