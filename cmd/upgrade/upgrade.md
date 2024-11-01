# Upgrade

This document tracks opon questions and discussions related to upgrading
v1 skupper sites to v2.

## Decisions

1. Upgrade tool will be given kube access to skupper sites.

## Questions

1. Is it correct for the tool to create AccessGrants / AccessTokens?  Or should the tool create Links?

   * The plus of creating Links is: the cost can be added to the v2 Link CR. 
   * TODO: how to save the cost in the AccessToken CR.  
   * When creating an AccessGrant, these fields are hardcoded. 
     I did not see where to retrieve these values from a v1 site.
```
            DefaultRedemptionsAllowed := 10
            DefaultExpirationWindow := "1h"
```

2. Is it ok for the tool to assume that v1 site names are unique?

   * When the tool creates AccessGrant / AccessToken, names the token <originating-sitename>-to-<target-sitename>.
     My thought was: after V2 is up and running, a 2nd script can match AccessGrants / AccessTokens using the name.

3. I have open questions on how to name listeners, connectors when there are multiple targets / ports.

4. At the moment, the code is calling ReadSiteConfig from the v2 branch.  
   At some point, this should be updated to call ReadSiteConfig from the v1 branch. 

## TODOs

