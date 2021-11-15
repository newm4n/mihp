# MIHP

**MIHP** is *short* of "MIHP Is HTTP Probe". It's another attempt to create another
*Synthetic Monitoring* tool, created using Golang programming language suitable for 
making system tool applications.

Basicaly, what MIHP does is executing HTTP Call(s) toward remote HTTP *end-point* and
then analyse the HTTP server's responses. Making sure that the *end-point* performance is
consistent throughout day-to-day, hour-by-hour and minute-by-minute 24/7 operations.

There are many uses of MIHP, for webmasters, web developers, quality assurances, dev-ops, etc.
It established the workflow from testers to devops. 

## Why another synthetic monitoring tool?

Other systhetic monitoring tools is already too bloated. Too complex for many simple
testing and monitoring activities. They also often too expensive to operate in house and
the online counter-part is not robust enough to be tailored for complex monitor.

MIHP also lightweight and small, the executable only took around 20mb, which already include
the probing, the **Minion** (stand-alone probing automation) and the **Central** 
(server that organizes multiple of minions). Beat that !!

- For Web Developers : You can create a unit test to check some web application workflow, (from Login to Dashboard to Shopping Cart to Purchase, etc). Making sure that EPIC is done.
- For QA and Testers : You can gather multiple MIHP config, put them in GIT and have an entire set of regression testing, ready to execute in your deployment pipeline. No more manual and cumbersome test. All automated.
- For DevOps : Deploy your MIHP config into **Central** and have your production be tested for performance and workflow correctness. 24/7. Be alerted through email, slack, telegram. Watch performance graph and SLA performance.
- For DataCenters : Provide services and ability for your customer to be confident if they service is running as expected.

## Where to start

### Obtaining MIHP

## Using MIHP Probing

### Running your first Probe

### Configuring Probe

### Understanding ProbeContext

### Understanding Probe Request

### Chaining Probe Request

### Understanding Expression

### Chaining Requests

## MIHP Minion

### Configuring Minion

### Running Minion

### Minion Health Check

## MIHP Central

### Configuring Central

### Managing Users and Organization

### Managing your Minions 

### Assigning Probe to Minion

