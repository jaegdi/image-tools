package ocrequest

func GetDocPage() string {
	docPage := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Webservice IMAGE-TOOL Documentation</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 40px; }
			h1 { color: #333; }
			p { margin: 10px 0; }
			code { background-color: #f4f4f4; padding: 2px 4px; border-radius: 4px; }
		</style>
		<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.11.5/css/jquery.dataTables.css">
		<script type="text/javascript" charset="utf8" src="https://code.jquery.com/jquery-3.5.1.js"></script>
		<script type="text/javascript" charset="utf8" src="https://cdn.datatables.net/1.11.5/js/jquery.dataTables.js"></script>
		<style>
			html, body {
				height: 100%;
				margin: 10px;
				padding: 0;
			}
			.dataTables_wrapper {
				height: 100%;
				display: flex;
				flex-direction: column;
			}
			.dataTables_wrapper .dataTables_filter {
				display: flex;
				justify-content: space-between;
				align-items: center;
			}
			table {
				width: 100%;
			}
			.dataTables_scrollBody {
				flex: 1 1 auto;
				overflow: auto;
			}
			.dataTables_scrollHead {
				flex: 0 0 auto;
			}
		</style>
		<script>
			$(document).ready(function() {
				var table = $('table').DataTable({
					"scrollY": "calc(100vh - 150px)",
					"scrollCollapse": true,
					"paging": false,
					"scrollX": true
				});
				$('.dataTables_filter').append(` + "`" + `
					<button onclick="downloadTableAsExcel()">Download as Excel</button>
				` + "`" + `);
			});
		</script>
		<style type="text/css">
			.left
			{
				float: left;
				margin-left: 50px;
			}

			.right
			{
				float: right;
				margin-right: 50px;
			}
			#content
			{
			}
		</style>
	</head>
	<body>
		<h1>Webservice IMAGE-TOOL Documentation</h1>

		<p style="margin-left: 25px;">Welcome to the webservice documentation page of the image-tool. Below you will find information on how to use the webservice.</p>
		<p style="margin-left: 25px;">This webservice provides a sub of the functionalilies, that the <code>cli image-tool</code> provides.</p>
		<p style="margin-left: 25px;">Therefore for more detailed queries or generation of image prune skripts unds the <code>cli image-tool</code>.</p>
		<br>
		<p style="margin-left: 25px;">There are two main function:</p>
		<br>
		<p style="margin-left: 25px;">/is-tag-used:   check if images with this tag are used somewhere in the cluster and give a json response with the answer.</p>
		<br>
		<p style="margin-left: 25px;">/query:         execute a query to find all objects and genereates a table with this objects and detail informatin as answer.</p>
		<br>

		<h2>Endpoints</h2>
		<h3>GET /</h3>
		<p>This endpoint show this documentation</p>

		<h3>GET /swagger</h3>
		<p>This endpoint provides the Swagger documentation for the webservice. <a href="/swagger">Click here</a> to access the Swagger documentation.</p>

		<br>
		<hr> <!-- Horizontale Linie eingefügt -->

		<h3>GET /is-tag-used</h3>
		<p>This endpoint checks, if the as parameter given tagname is as imagetag used in a pod somewhere in the clusters.</p>
		<h4>Query Parameters:</h4>
		<ul>
			<li><code>family</code> (required for <code>is_tag_used</code>): The family parameter.</li>
			<li><code>tagname</code> (required for <code>is_tag_used</code>): The tagname parameter.</li>
		</ul>

		<h4>Responses:</h4>
		<ul>
			<li><code>200 OK</code>: The command was executed successfully. The response is in JSON format.
				<pre>
					<br>       eg.: <code>{"TagIsUsed":true,"TagName":"pkp-3.19.0-build-3"}</code>
					<br>       eg.: <code>{"TagIsUsed":false,"TagName":"pkp-x-not-there"}</code>
				</pre>
			</li>

			<li><code>400 Bad Request</code>: Missing or invalid parameters.</li>
		</ul>

		<h4>Example usage:</h4>
		<p><a href="/is-tag-used?family=pkp&kind=is_tag_used&tagname=pkp-3.19.0-build-3">GET /is-tag-used?family=pkp&kind=is_tag_used&tagname=pkp-3.19.0-build-3</a></p>
		<h4>Example answer</h4>
		<code>{"TagIsUsed":true,"TagName":"pkp-3.19.0-build-3"}</code>
		<br>
		<hr> <!-- Horizontale Linie eingefügt -->

		<h3>GET /query</h3>
		<p>This endpoint executes a query to generate several lists as answers, depending onthe parameters.</p>
		<h4>Query Parameters:</h4>
		<table>
		<ul>
			<tr><td><li><code>family</code>:</td>  	<td>The family parameter, eg. pkp, ebs, b2c, b2b, vps, ...</li></td></tr>
			<tr><td><li><code>cluster</code>:</td> 	<td>The cluster parameter, eg. cid-scp0, ppr-scp0 ... or pro-scp0</li></td></tr>
			<tr><td><li><code>kind</code>:</td>    	<td>The kind of operation to perform. Valid values are: <code>istag, image, is, used</code>. <br>The default is <code>istag</code></pre></li></td></tr>
			<tr><td><li><code>tagname</code>:</td> 	<td>optional. This is a name or a regex pattern of an imagetag. This string is alway interpreted as regex. <br>Valid values are: <code>pkp-3.19, build-1, v1.0.*, used</code>. <br>Special chars like '.*{}[]:' must be masked with a '\', if they want explicitly set in the string.</pre></li></td></tr>
			<tr><td><li><code>namespace</code>:</td> <td>optional. List only objects from this namespace.</li></td></tr>
		</ul>
		</table>
		<h4>Responses:</h4>
		<ul>
			<li><code>200 OK</code>: The command was executed successfully. The response is in text/HTML format.
			<br>A HTML table wit all of the image details found in the cluster.
			<br>For comfort the table has some functions
			<br> - to filter
			<br> Download as Excel
			<br> sort function on every column

			<li><code>400 Bad Request</code>: Missing or invalid parameters.</li>
		</ul>

		<h4>Example usage:</h4>
		<p>Get all istags from family=scp in all clusters with any tagname</p>
		<p><a href="/query?family=scp&kind=istag&cluster=cid-scp0,ppr-scp0,vpt-scp0,pro-scp0">GET /query?family=scp&kind=istag&cluster=cid-scp0,ppr-scp0,vpt-scp0,pro-scp0</a></p>
		<p>Get all istags from family=pkp in all clusters with any tagname</p>
		<p><a href="/query?family=pkp&kind=istag&cluster=cid-scp0,ppr-scp0,vpt-scp0,pro-scp0">GET /query?family=pkp&kind=istag&cluster=cid-scp0,ppr-scp0,vpt-scp0,pro-scp0</a></p>
		<p>Get all used istags from family=pkp in all clusters with any tagname</p>
		<p><a href="/query?family=pkp&kind=used">GET /query?family=pkp&kind=used</a></p>
		<hr>

		<h4>Example usage:</h4>
		<p>Get all istags from family=pkp in cluster=cid-scp0, where the tagname contains '0.9'</p>
		<p><a href="/query?family=pkp&cluster=cid-scp0&kind=istag&tagname=0.9">GET /query?family=pkp&cluster=cid-scp0&kind=istag&tagname=0.9</a></p>

		<h4>Example answer</h4>
		<hr> <!-- Horizontale Linie eingefügt -->
		<br>
		<table class="game-of-thrones">
			<thead>
				<tr>
					<th>&nbsp;</th>
					<th>istag pkp</th>
					<th>Cluster</th>
					<th>Imagestream</th>
					<th>Tagname</th>
					<th>Namespace</th>
					<th>Date</th>
					<th>AgeInDays</th>
					<th>Image</th>
					<th>CommitAuthor</th>
					<th>CommitDate</th>
					<th>CommitId</th>
					<th>CommitRef</th>
					<th>CommitVersion</th>
					<th>IsProdImage</th>
					<th>BuildName</th>
					<th>BuildNamespace</th>
				</tr>
			</thead>
			<tbody>
			<tr>
				<td align="right">1</td>
				<td>bi-service-service:0.9.1-0</td>
				<td>cid-scp0</td>
				<td>bi-service-service</td>
				<td>0.9.1-0</td>
				<td>pkp-images</td>
				<td>2024-07-15T15:12:44Z</td>
				<td>   87</td>
				<td>sha256:a0609b629cd6d90091347b85463db9cb2f91cdcf3434957c107f98e7268182df</td>
				<td>&#34;Yannick Guth &lt;yannick.guth@schufa.de&gt;&#34;</td>
				<td>&#34;Mon, 15 Jul 2024 17:05:45 +0200&#34;</td>
				<td>3c7878055bd94f576ced21f6eb244b6c9d4ff03f</td>
				<td>PR-28</td>
				<td>0.9.1</td>
				<td>false</td>
				<td>bi-service-pr-28-1</td>
				<td>pkp-build</td>
			</tr>
			<tr>
				<td align="right">2</td>
				<td>bi-service-service:0.9.2-0</td>
				<td>cid-scp0</td>
				<td>bi-service-service</td>
				<td>0.9.2-0</td>
				<td>pkp-images</td>
				<td>2024-07-18T10:16:00Z</td>
				<td>   85</td>
				<td>sha256:bab8203c1ed7f3510454c41803a911fdfdd982d031084780f5640a61528a9453</td>
				<td>&#34;Klaus Mandola &lt;Klaus.Mandola@schufa.de&gt;&#34;</td>
				<td>&#34;Thu, 18 Jul 2024 12:08:01 +0200&#34;</td>
				<td>00f77dc03a8fb32f713bac42e31bde47f1aa1588</td>
				<td>PR-29</td>
				<td>0.9.2</td>
				<td>false</td>
				<td>bi-service-pr-29-1</td>
				<td>pkp-build</td>
			</tr>
			<tr>
				<td align="right">3</td>
				<td>bi-service-service:0.9.3-0</td>
				<td>cid-scp0</td>
				<td>bi-service-service</td>
				<td>0.9.3-0</td>
				<td>pkp-images</td>
				<td>2024-07-19T11:08:57Z</td>
				<td>   84</td>
				<td>sha256:ccfb906c328fcf7a92a2d0a5707c16aff1bbe5534710e53601e9045f9655272d</td>
				<td>&#34;Klaus Mandola &lt;Klaus.Mandola@schufa.de&gt;&#34;</td>
				<td>&#34;Fri, 19 Jul 2024 13:01:56 +0200&#34;</td>
				<td>f9e18d501b3e5a8969dbf1a9a0ec82fe4d82fbf0</td>
				<td>PR-29</td>
				<td>0.9.3</td>
				<td>false</td>
				<td>bi-service-pr-29-1</td>
				<td>pkp-build</td>
			</tr>
			<tr>
				<td align="right">4</td>
				<td>bi-service-service:0.9.4-0</td>
				<td>cid-scp0</td>
				<td>bi-service-service</td>
				<td>0.9.4-0</td>
				<td>pkp-images</td>
				<td>2024-07-19T13:29:20Z</td>
				<td>   83</td>
				<td>sha256:d94b198188e133e7fdef3b56f4926d26a2f5f1883a9b278e33aa0dd6eeefb8d6</td>
				<td>&#34;Yannick Guth &lt;yannick.guth@schufa.de&gt;&#34;</td>
				<td>&#34;Fri, 19 Jul 2024 15:21:38 +0200&#34;</td>
				<td>1005ada79194ce5771eabac0e5a9e8c75146aa0d</td>
				<td>PR-28</td>
				<td>0.9.4</td>
				<td>false</td>
				<td>bi-service-pr-28-1</td>
				<td>pkp-build</td>
			</tr>
			<tr>
				<td align="right">5</td>
				<td>bi-service-service:0.9.5-0</td>
				<td>cid-scp0</td>
				<td>bi-service-service</td>
				<td>0.9.5-0</td>
				<td>pkp-images</td>
				<td>2024-09-13T10:24:37Z</td>
				<td>   28</td>
				<td>sha256:5bb342d50b91cc86353e4a6ee223c8f81c66df3e4250b2f9247ad63bfecdfe99</td>
				<td>&#34;Jenkins &lt;grow-jenkins@schufa.de&gt;&#34;</td>
				<td>&#34;Tue, 10 Sep 2024 09:52:40 +0000&#34;</td>
				<td>8a0b25bc77bb70ffeb456a1ef71cea518d1c1343</td>
				<td>PR-31</td>
				<td>0.9.5</td>
				<td>false</td>
				<td>bi-service-pr-31-1</td>
				<td>pkp-build</td>
			</tr>
			<tr>
				<td align="right">6</td>
				<td>payment-async-adapter-service:0.9.0-EPIC-1-0</td>
				<td>cid-scp0</td>
				<td>payment-async-adapter-service</td>
				<td>0.9.0-EPIC-1-0</td>
				<td>pkp-images</td>
				<td>2024-09-05T07:16:15Z</td>
				<td>   36</td>
				<td>sha256:402c619e272df737bd57a1501313a19c1190ac1eac7e5e2814c41240ef3e06b7</td>
				<td>&#34;Gerhard Dickescheid &lt;gerhard.dickescheid@schufa.de&gt;&#34;</td>
				<td>&#34;Thu, 5 Sep 2024 08:45:09 +0200&#34;</td>
				<td>b45e232b06fa51cfc5367c2cd8f4a558d4b67374</td>
				<td>PR-2</td>
				<td>0.9.0-EPIC-1</td>
				<td>false</td>
				<td>payment-async-adapter-pr-2-1</td>
				<td>pkp-build</td>
			</tr>


			</tbody>

			<tfoot>
				<tr>
					<td>&nbsp;</td>
					<td> </td>
					<td> </td>
					<td> </td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
					<td>&nbsp;</td>
				</tr>
			</tfoot>
		</table>
		<script>
			function downloadTableAsExcel() {
				var table = document.querySelector('.dataTables_scrollBody > table');
				var html = table.outerHTML;
				var url = 'data:application/vnd.ms-excel,' + escape(html);
				var link = document.createElement('a');
				link.href = url;
				link.setAttribute('download', 'table.xls');
				link.click();
			}
		</script>
		<br>
		<hr> <!-- Horizontale Linie eingefügt -->
	</body>
	</html>
	`
	return docPage
}
