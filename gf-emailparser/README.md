# gf-emailparser

gf-emailparser is a Go project used to automate reading flight alert emails from google flights and output 
their contents to a google sheets file. There are jobs in the module which will control: 

* Extracting email contents and outputting to google sheets
* Sorting through the flights in google sheets for flights workable matching flights - the best deals
will then be stored in a table which is reachable via API. 

A web server is included which will provide API access to the found flights. 