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
	</head>
	<body>
		<h1>Webservice IMAGE-TOOL Documentation</h1>
		<p>Welcome to the webservice documentation page of the image-tool. Below you will find information on how to use the webservice.</p>
		<h2>Endpoints</h2>
		<h3>GET /</h3>
		<p>This endpoint show this documentation</p>
		<h3>GET /execute</h3>
		<p>This endpoint executes a command based on the provided query parameters.</p>
		<p><strong>Query Parameters:</strong></p>
		<ul>
			<li><code>family</code> (required for <code>is_tag_used</code>): The family parameter.</li>
			<li><code>kind</code>: The kind of operation to perform. Valid values are <code>is_tag_used</code>.
			<br><pre>       The default is <code>is_tag_used</code></pre></li>
			<li><code>cluster</code>: The cluster parameter (is not necessary for kind <code>is_tag_used</code>): eg. cid-scp0, ... or pro-scp0</li>
			<li><code>tagname</code> (required for <code>is_tag_used</code>): The tagname parameter.</li>
			<li><code>namespace</code>: The namespace parameter.</li>
		</ul>
		<p><strong>Responses:</strong></p>
		<ul>
			<li><code>200 OK</code>: The command was executed successfully. The response is in JSON format.
			<br><pre>       eg.: <code>{"TagIsUsed":true,"TagName":"pkp-3.19.0-build-3"}</code>
			<br>       eg.: <code>{"TagIsUsed":false,"TagName":"pkp-x-not-there"}</code></pre></li>

			<li><code>400 Bad Request</code>: Missing or invalid parameters.</li>
		</ul>
		<p>Example usage:</p>
		<pre><code>GET /execute?family=pkp&kind=is_tag_used&tagname=pkp-3.19.0-build-3</code></pre>
	</body>
	</html>
	`
	return docPage
}
