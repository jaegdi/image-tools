# report-istags

## DESCRIPTION

istag-mgmt reports image date for a application family (eg. pkp, fpc, aps, ssp)

- For __existing Is, IsTags and Images__ it operates cluster and family specific. 
    That means it works for __one cluster__ like
    'cid, int, ppr, vpt or pro' and for __one family__ like 'pkp, sps, fpc, aps, ...'
    The cluster must be defined by the mandatory parametter '-cluster=[cid|int|ppr|vpt|pro]'
    The family must be defined by the mandatory parameter '-family=[aps|fpc|pkp|ssp]

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

## Usage

### Command

    ./report-istags -family=... -cluster=... -all|-image|-is|-istag|-used [output format (default json)] [filter (default none)] 

the parameters can be specified in any order

#### Parameters

##### Define source, both parameters are mandatory

    -family=string    family name, eg. pkp, aps, ssp or fpc
    
    -cluster=string   shortname of cluster, eg. cid,int, ppr or pro

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

##### Filter

    -namespace=string namespace to look for istags
    
    -isname=string    filter output of one imageStream as json, eg. -is=wvv-service
    
    -istagname=string filter output of one imageStreamTag
    
    -shaname=string   filter output of a Image with this SHA
    
    -tagname=string   filter output all istags with this Tag

##### Output Format

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

    -noproxy          disable the usage of a proxy for OpenShift API requests

    -socks5=string    enable socks5 usage. E.g. -socks5=localhost:65022

## EXAMPLES

Report all information for family pkp in cluster cid as json(which is the default output format)

    ./report-istags -cluster=cid -family=pkp -all

Report only used istags for family pkp as pretty printed table.
(the output is paginated to fit your screen size, use it piped to less. Then you can go up or down with the page key)

    ./report-istags -cluster=cid -family=pkp -used -table | less

Report istags for family aps in cluster int as yaml report

    ./report-istags -cluster=int -family=aps -istag -yaml
