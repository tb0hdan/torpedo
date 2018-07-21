# [Back to main doc](../README.md)

# TorpedoBot Remote Plugin Execution (TRPE)

TRPE allows writing plugins in any language as long as content is returned via HTTP API.
Sample application is available at `tools/trpe_server.py`

Architecture is as follows:

![TRPE](https://raw.githubusercontent.com/tb0hdan/torpedo/master/doc/TRPE.png)


Bot should be launched with `-trpe_host` switch, i.e.:

`bin/torpedobot -trpe_host http://localhost:5000/trpe`


TRPE server:

`tools/trpe_server.py`
