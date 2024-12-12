package ocrequest

import (
	"flag"
	"fmt"
	"os"
)

// cmdUsage print the man page
func CmdUsage() {
	usageText := `

DESCRIPTION

  image-tool
        - generate reports about images, imagestream(is), imagestreamtags(istag)
          with information like AgeInDays, Date, Namespace, Buildtags,..
          for a application family (eg. pkp, fpc, aps, ssp)
        or
        - generate shellscript output to delete istags when parameter -delete is set.
          The -delete parameter disables the report output, instead the delete script is generated as output


  Per defaul image-tool is executed as a cmdline tool but with parameter -server it can be startet as a webservice.

  image-tool as web-service

    In serverMode image-tool looks for used images in all clusters filtered by family and tagname and returns the a
    JSON list as HTTP response

    image-tool starts a webserver to listen on port 8080.
    To request the webservice by curl, use the following pattern:

        curl "http://localhost:8080/execute?family=exampleFamily&tagname=exampleTagname" | jq

PARAMETERS

  `
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	fmt.Println(usageText)
	flag.PrintDefaults()
	fmt.Println("\nTo get the man page, set parameter '-man'")
}

func ManPage() {
	manpageText := `

MAN PAGE

  image-tool
        - Generates reports about images, image streams (is), and image stream tags (istag)
          with information such as AgeInDays, Date, Namespace, Buildtags, etc.
          for an application family (e.g., pkp, fpc, aps, ssp)
        or
        - Generates a shell script to delete istags when the -delete parameter is set.
          The -delete parameter disables the report output and instead generates the delete shell script as output.

  By default, image-tool is executed as a command-line tool, but with the -server parameter, it can be started as a web service.

  image-tool as a web service

    In server mode, image-tool looks for used images in all clusters filtered by family and tag name and returns a
    JSON list as an HTTP response.

    image-tool starts a web server that listens on port 8080.
    To request the web service using curl, use the following pattern:

        curl "http://localhost:8080/execute?family=exampleFamily&tagname=exampleTagname" | jq

  image-tool as cmdline tool

    image-tool only read information around images from the clusters and generate output. It never change or
    delete something in the clusters. Eg. for delete istags it only generates a script output, which than can
    be executed by a cluster admin to really delete the istags.

    image-tool always write a log to 'image-tool.log' in the current directory.
    Additional it writes the log messages also to STDERR. To disable the log output to STDERR use parameter -nolog
    The report or delete-script output is written to STDOUT.

  - For reporting existing Images and delete istags it operates cluster and family specific.
    For reporting used images it works over all clusters but family specific.
    It never works across different families.

      For existing images, istags or imagestreams(is) that means it works for one or more clusters like
        cid-scp0,  ppr-scp0, vpt-scp0, pro-scp0
        or more than one like 'cid-scp0,cid-scp0,ppr-scp0,ppr-scp0'

    and for one families like
        'pkp, sps, fpc, aps, ...'

    The cluster must be defined by the mandatory parameter. To run the query across multiple clusters,
    multiple clusters can be specified, separated by commas.
        '-cluster=[cid-scp0|ppr-scp0|vpt-scp0|pro-scp0|dev-scp0|dev-scp1-c1|dev-scp1-c2|pro-scp1]'

    The family must be defined by the mandatory parameter
        '-family=[aps|fpc|pkp|ssp]

  - Generate reports about imagestreamtags, imagestreams and images.
        The content of the report can be defined by the mandatory parameter
        '-output=[is|istag|image|used|unused|all]'.

  - For used images it looks in all clusters and reports the istags used
      by any deploymentconfig, job, cronjob or pod of all namespaces that
      belong to the application family.

  - Variable output format: json(default), yaml, csv, csvfile, table and tabgroup
      (table with grouped rows for identical content).
      Output as table or tabgroup is automatically piped into less (or what is defined as PAGER)

  - filter data for reports.
      By specifying one of the parameters
          -isname=..., -istagname=..., tagname=... or -shaname=...
      the report is filtered.

  - delete istags based on filters
      The idea is to delete istags by filterpatterns like
      'older than n days' and/or 'istag name like pattern'
      The image tool didn't delete the istags directly instead
      it generate a shell-script that can be executed by a cluster admin
      to delete the istag, they fit to the given filter parameters

      To switch from reporting mode to delete mode, set the praameter -delete
      But it needs further parameters:
      -snapshot        delete istags with snapshot or PR-nn in the tag name.
      -nonbuild        is specific for family pkp and delete istags fo all images, if they have no build tag.
      -delminage=int      defines the minimum age (in days) for istag to delete them. Default is 60 days.
      -delpattern=str  define a regexp pattern for istags to delete
      and
      -isname=str
      -tagname=str
      -istagname=str
      -namespace=str
      can also be used to filter istags to delete them.
      See examples in the EXAMPLES section

    For this reports the data is collected from the openshift cluster defined by
    the mandatory parameters
         '-cluster=...' and the
         'family=...'
    For type '-used' (also included in type '-all') the data is collected
    from all clusters.

    For more speed a cache is build from the first run in  /tmp/tmp-report-istags/*
    and used if not older than 5 minutes. If the cache is older or deleted, the
    data is fresh collected from the clusters.


INSTALLATION

    image-tool is a statically linked go programm and has no runtime dependencies.
    No installation is neccessary. Copy the binary into a directory, which is in the search path.
    Copy the binary from local artifactory (as long as it exists)
        - linux:  https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-develop/istag_and_image_management/image-tool-linux/image-tool
        - windows:  https://artifactory-pro.sf-rz.de:8443/artifactory/scpas-bin-develop/istag_and_image_management/image-tool-windows/image-tool.exe
    or copy it from the new artifactory SaaS
        - linux:   https://atfschufa.jfrog.io/artifactory/scptools-bin-develop/tools/image-tool/image-tool-linux/image-tool
        - windows: https://atfschufa.jfrog.io/artifactory/scptools-bin-develop/tools/image-tool/image-tool-windows/image-tool.exe
    you can also use the jf tool
        - linux:   jf rt dl scptools-bin-develop/tools/image-tool/image-tool-linux/image-tool
        - windows: jf rt dl scptools-bin-develop/tools/image-tool/image-tool-windows/image-tool.exe


EXAMPLES

  SERVERMODE
  |
  |        Start the image-tool with the parameter -server
  |        The image-tool opens a listening port on localhost:8080 to listen for requests
  |
  |           image-tool -server
  |
  |        Then a request can sendt to query if an given istag is used somewhere in the clusters
  |
  |           curl "http://localhost:8080/execute?family=exampleFamily&tagname=exampleTagname" | jq
  |
  |           eg. curl "http://localhost:8080/execute?family=pkp&tagname=pkp-3.20.0-build-1" | jq
  |
  -----------------------------------------------------------------------------------------------------------------------------


  REPORTING
  |
  |        Report all information for family pkp in cluster cid as json
  |        (which is the default output format)
  |
  |            image-tool -cluster=cid-scp0 -family=pkp -all
  |
  |            or as table
  |            image-tool -cluster=cid-scp0 -family=pkp -all -table
  |
  |            or csv in different files for each type of information
  |            image-tool -cluster=cid-scp0 -family=pkp -all -csvfile=prefix
  |            writes the output to different files 'prefix-type' in current directory
  |
  |        Report only __used__ istags for family pkp as pretty printed table
  |        (the output is paginated to fit your screen size and piped to
  |            the pager define in the environment variable $PAGER/%PAGER%.
  |            If $PAGER is not set, it try to use 'more')
  |
  |            image-tool -cluster=cid-scp0 -family=pkp -used -table
  |            or json
  |            image-tool -cluster=cid-scp0 -family=pkp -used
  |            or yaml
  |            image-tool -cluster=cid-scp0 -family=pkp -used -yaml
  |            or csv
  |            image-tool -cluster=cid-scp0 -family=pkp -used -csv
  |
  |        Report istags for family aps in cluster int as yaml report
  |
  |            image-tool -cluster=int-scp0 -family=aps -istag -yaml
  |
  |        Report ImageStreams for family aps in cluster int as yaml report
  |
  |            image-tool -cluster=int-scp0 -family=aps -is -yaml
  |
  |        Report Images for family aps in cluster int as yaml report
  |
  |            image-tool -cluster=int-scp0 -family=aps -image -yaml
  |
  |        Report combined with pc(print columns) tool
  |
  |            image-tool -socks5=localhost:65022 -family=pkp -cluster=cid-scp0,ppr-scp0,pro-scp0 -istag -csv | pc -sep=, -sortcol=4  1 5 8 6 7
  |
  -----------------------------------------------------------------------------------------------------------------------------


  DELETE
  |
  |        Generate a shell script to delete old istags(60 days, the default) for family pkp in cluster cid
  |        and all old snapshot istags and nonbuild istags and all istags of header-service, footer-service
  |        and zahlungsstoerung-service
  |
  |            image-tool -family=pkp -cluster=cid-scp0 -delete -snapshot \
  |                        -nonbuild -delpattern='(header|footer|zahlungsstoerung)-service'
  |
  |        To use the script output to really delete the istags, you can use the following line:
  |
  |            image-tool -family=pkp -cluster=cid-scp0 -delete -snapshot -nonbuild \
  |                        -delpattern='(header|footer|zahlungsstoerung)-service'      | xargs -n 1 -I{} bash -c "{}"
  |
  |        To only generate a script to delete old snapshot istags:
  |
  |            image-tool -family=pkp -cluster=cid-scp0 -delete -snapshot
  |
  |        To delete all not used images of family 'aps' in cluster cid
  |
  |            image-tool -family=aps -cluster=cid-scp0 -delete  -delminage=0 -delpattern='.'
  |
  |        To delete all hybris istags of family pkp older than 45 days
  |
  |            image-tool -family=pkp -cluster=cid-scp0 -delete -isname=hybris -delminage=45
  |
  |   !!!  To check, which images are not deleted, because they are in use somewhare in the clusters,
  |        you can check the logfile, which is created in the current directory.
  |
  |            cat image-tool.log  |  grep logUsedIstags:  |  sort -k5  |  column -t
  |
  | HINT
  |
  |        To directly delete the istags, that reportet by 'image-tool -delete ...', make shure, you are
  |        logged in into the correct cluster, because the output is executed with oc client and work on the
  |        currently logged in cluster. And append the following to the end of the image-tool - command:
  |
  |
  |            | xargs -n 1 -I{} bash -c "{}"
  |
  |        After deleting the istags, the images must removed from the registry by executing a command similar
  |        to this example:
  |
  |            oc login ..... to the cluster
  |            registry_url="$(oc -n default get route|grep docker-registry|awk '{print $2}')"
  |            oc adm prune images --registry-url=$registry_url --keep-tag-revisions=3 --keep-younger-than=60m --confirm
  |
  |        or if you have the admin-tools repo installed you can use the script
  |
  |            prune-registry-of-current-cluster.sh
  |
  -----------------------------------------------------------------------------------------------------------------------------


  CONNECTION
  |        As default the sock5 proxy to localhost:65022 is enabled becaus the api of the upper clusters is only reacheable
  |        over the sprungserver. To disable SOCKS5 set the parameter -socks5=no
  |        If your socks5 jumpserver config listens on a different port set the parameter -socks5=<host>:<port>
  |
  |        If there are problems with the connection to the clusters,
  |        there is the option to disable the use of web proxy with the
  |        parameter '-noproxy'.
  |
  |        A socks5 proxy can be the solution, eg. to run it from your notebook over VPN, then establish
  |        a socks tunnel over the sprungserver and give the
  |        parameter '-socks5=host:port' to the image-tool program.
  |
  -----------------------------------------------------------------------------------------------------------------------------
`

	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	fmt.Println(manpageText)
}
