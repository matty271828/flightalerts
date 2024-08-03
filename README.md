# FlightAlerts

FlightAlerts is a multi module repository which contains various modules involved in monitoring
cheap flight prices. 

## gf-alertsetter

gf-alertsetter is a python module used to automate opening a chrome browswer and setting google flight alerts 
for a list of routes defined in `routes.py`. There is a hard limit of 100 flight alerts per email address, 
so multiple email addresses are defined to alert to in `gmail_accounts.txt`.

## gf-flightsparser

gf-flightsparser is a Go module used to automate reading flight alert emails and output their contents to a google
sheets file. There are jobs in the module which will control: 

* Extracting email contents and outputting to google sheets
* Sorting through the flights in google sheets for flights workable matching flights - the best deals
will then be stored in a table which is reachable via API. 

A web server is included which will provide API access to the found flights. 
