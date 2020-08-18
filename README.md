# urlshortener
bluefield
Create URL shortener server

Implement HTTP server that can generate shortened URLs.

The requests to shortened URLs should be redirected to their original URL (status 302) or
return 404 for unknown URLs.

Simple HTML form should be served on the index page where users can input URL and
retrieve the shortened version from server.

All of the implemented HTTP handlers should have unit tests.
(optional) All shortened URLs should be persisted locally to a file using simple storage
methods (SQLite, BoltDB, CSV..).

(optional) The redirect requests should be cached in memory for certain amount of time.